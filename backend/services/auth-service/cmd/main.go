// auth-service: JWT 인증/인가 마이크로서비스
//
// 포트: gRPC :50051
// 의존: PostgreSQL(선택), Redis(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 사용자 등록/로그인 (이메일 + bcrypt)
// - JWT 발급 (Access 15분 + Refresh 7일)
// - 토큰 갱신 (Refresh Token Rotation)
// - 로그아웃 (전체 토큰 철회)
// - gRPC AuthService (Register, Login, RefreshToken, Logout, ValidateToken)
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/auth-service/internal/handler"
	"github.com/manpasik/backend/services/auth-service/internal/oauth"
	"github.com/manpasik/backend/services/auth-service/internal/repository/memory"
	"github.com/manpasik/backend/services/auth-service/internal/repository/postgres"
	redisTokenRepo "github.com/manpasik/backend/services/auth-service/internal/repository/redis"
	"github.com/manpasik/backend/services/auth-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"github.com/manpasik/backend/shared/ops/secrets"
	"github.com/manpasik/backend/shared/tenancy"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "auth-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)

	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)

	tp, tpErr := observability.InitTracer(serviceName)
	if tpErr != nil {
		log.Printf("[WARN] Tracing init failed: %v", tpErr)
	} else {
		defer tp.Shutdown(context.Background())
	}

	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)
	log.Printf("[%s] gRPC port: %s", serviceName, cfg.GRPCPort)

	// UserRepository: PostgreSQL 또는 인메모리
	// DB_HOST 환경변수가 명시적으로 설정된 경우에만 PostgreSQL 사용 시도
	var userRepo service.UserRepository
	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if err != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, err)
			userRepo = memory.NewUserRepository()
		} else {
			// 실제 DB 연결 검증 (Ping)
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				userRepo = memory.NewUserRepository()
			} else {
				pingCancel()
				defer pool.Close()
				userRepo = postgres.NewUserRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		userRepo = memory.NewUserRepository()
		log.Printf("[%s] 인메모리 User 저장소 사용", serviceName)
	}

	// TokenRepository: Redis 또는 인메모리
	var tokenRepo service.TokenRepository
	if _, redisHostSet := os.LookupEnv("REDIS_HOST"); redisHostSet && cfg.Redis.Host != "" {
		redisClient := redisclient.NewClient(&redisclient.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := redisClient.Ping(pingCtx).Err(); err != nil {
			pingCancel()
			log.Printf("[%s] Redis 연결 실패, 인메모리 TokenRepo 사용: %v", serviceName, err)
			tokenRepo = memory.NewTokenRepository()
		} else {
			pingCancel()
			defer redisClient.Close()
			tokenRepo = redisTokenRepo.NewTokenRepository(redisClient)
			log.Printf("[%s] Redis TokenRepo 연결됨: %s", serviceName, cfg.Redis.Addr())
		}
	} else {
		tokenRepo = memory.NewTokenRepository()
		log.Printf("[%s] 인메모리 TokenRepo 사용", serviceName)
	}

	authSvc := service.NewAuthService(
		logger,
		userRepo,
		tokenRepo,
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)

	// Vault 시크릿 자동 회전: VAULT_ADDR 설정 시 활성화.
	// 시크릿 path JWT_SECRET_VAULT_PATH (기본 "auth/jwt") 의 변경을 감지하여
	// authSvc.SetJWTSecret 자동 호출 — 서비스 재시작 없이 키 회전.
	var vaultBootstrap *secrets.VaultBootstrap
	if vaultAddr := os.Getenv("VAULT_ADDR"); vaultAddr != "" {
		jwtPath := os.Getenv("JWT_SECRET_VAULT_PATH")
		if jwtPath == "" {
			jwtPath = "auth/jwt"
		}
		bs, vErr := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
			WatchPaths:                []string{jwtPath},
			AutoRenewToken:            true,
			TokenRenewIntervalSeconds: 1800,
			OnError: func(path string, err error) {
				log.Printf("[%s] Vault polling 에러 (%s): %v", serviceName, path, err)
			},
		})
		if vErr != nil {
			log.Printf("[%s] Vault 초기화 실패, 환경변수 JWT 사용: %v", serviceName, vErr)
		} else {
			bs.AddListener(func(_ context.Context, ev secrets.RotationEvent) error {
				if ev.NewSecret == nil || ev.NewSecret.Value == "" {
					return nil
				}
				authSvc.SetJWTSecret(ev.NewSecret.Value)
				log.Printf("[%s] JWT 시크릿 핫리로드 완료 (path=%s, version=%d)",
					serviceName, ev.Path, ev.NewSecret.Version)
				return nil
			})
			bs.Start(context.Background())
			vaultBootstrap = bs
			defer bs.Stop()
			log.Printf("[%s] Vault Bootstrap 활성화 (watch=%s)", serviceName, jwtPath)
		}
	}
	_ = vaultBootstrap

	// OAuth Verifier 초기화 (환경변수 기반)
	oauthVerifiers := make(map[string]oauth.OAuthVerifier)

	// Google OAuth
	if clientID := os.Getenv("GOOGLE_CLIENT_ID"); clientID != "" {
		oauthVerifiers["GOOGLE"] = oauth.NewGoogleVerifier(clientID)
		log.Printf("[%s] Google OAuth Verifier 활성화", serviceName)
	}

	// Apple OAuth
	if clientID := os.Getenv("APPLE_CLIENT_ID"); clientID != "" {
		oauthVerifiers["APPLE"] = oauth.NewAppleVerifier(clientID)
		log.Printf("[%s] Apple OAuth Verifier 활성화", serviceName)
	}

	// Kakao OAuth (Access Token 기반, 별도 키 불필요)
	oauthVerifiers["KAKAO"] = oauth.NewKakaoVerifier()
	log.Printf("[%s] Kakao OAuth Verifier 활성화", serviceName)

	// Naver OAuth (Access Token 기반, 별도 키 불필요)
	oauthVerifiers["NAVER"] = oauth.NewNaverVerifier()
	log.Printf("[%s] Naver OAuth Verifier 활성화", serviceName)

	authSvc.SetOAuthVerifiers(oauthVerifiers)

	// TokenValidator 어댑터 (인터셉터용)
	validator := &authTokenValidator{auth: authSvc}

	// auth 는 로그인/등록/리프레시 등 테넌트 컨텍스트 형성 *전* 의 RPC 가 다수.
	// 테넌시 검증을 면제할 메서드 추가.
	authSkip := map[string]bool{
		"/manpasik.v1.AuthService/Register":             true,
		"/manpasik.v1.AuthService/Login":                true,
		"/manpasik.v1.AuthService/RefreshToken":         true,
		"/manpasik.v1.AuthService/ValidateToken":        true,
		"/manpasik.v1.AuthService/OAuthLogin":           true,
		"/manpasik.v1.AuthService/RequestPasswordReset": true,
	}
	tenancyEngine := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	serverOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
			middleware.AuthInterceptor(validator),
		),
	}
	serverOpts = append(serverOpts, tenancy.MaybeServerOptions(tenancyEngine, authSkip)...)
	if tenancy.EnforcedFromEnv() {
		log.Printf("[%s] 멀티테넌트 RBAC 인터셉터 활성화 (auth flow skip 포함)", serviceName)
	}
	grpcServer := grpc.NewServer(serverOpts...)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	authHandler := handler.NewAuthHandler(authSvc, logger)
	v1.RegisterAuthServiceServer(grpcServer, authHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("[%s] Failed to listen: %v", serviceName, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() {
			time.Sleep(cfg.ShutdownTimeout)
			os.Exit(1)
		}()
		grpcServer.GracefulStop()
		cancel()
	}()

	// Start observability HTTP server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		metricsAddr := ":9100"
		if configured := os.Getenv("METRICS_PORT"); configured != "" {
			metricsAddr = configured
		}
		logger.Info("Metrics server starting", zap.String("addr", metricsAddr))
		if err := http.ListenAndServe(metricsAddr, mux); err != nil {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	log.Printf("[%s] gRPC server listening on %s", serviceName, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}

// authTokenValidator는 AuthService를 middleware.TokenValidator로 래핑합니다.
type authTokenValidator struct {
	auth *service.AuthService
}

func (a *authTokenValidator) ValidateToken(token string) (*middleware.TokenClaims, error) {
	claims, err := a.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return &middleware.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
