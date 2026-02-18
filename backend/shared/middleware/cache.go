package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// CacheStore는 캐시 저장소 인터페이스입니다.
type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// CacheConfig는 gRPC 응답 캐시 설정입니다.
type CacheConfig struct {
	// MethodTTL은 메서드별 캐시 TTL을 정의합니다.
	// 키: gRPC full method name, 값: TTL duration
	MethodTTL map[string]time.Duration

	// DefaultTTL은 MethodTTL에 없는 메서드의 기본 TTL입니다.
	// 0이면 캐시하지 않습니다.
	DefaultTTL time.Duration

	// KeyPrefix는 캐시 키 접두사입니다.
	KeyPrefix string
}

// CacheInterceptor는 gRPC 응답 캐시 인터셉터입니다.
// GET 성격의 읽기 전용 메서드만 캐싱합니다.
func CacheInterceptor(store CacheStore, config *CacheConfig) grpc.UnaryServerInterceptor {
	if config == nil {
		config = &CacheConfig{KeyPrefix: "grpc_cache:"}
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "grpc_cache:"
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 캐시 TTL 결정
		ttl := config.DefaultTTL
		if methodTTL, ok := config.MethodTTL[info.FullMethod]; ok {
			ttl = methodTTL
		}

		// TTL이 0이면 캐시 건너뜀
		if ttl == 0 {
			return handler(ctx, req)
		}

		// 쓰기 메서드(Create/Update/Delete)는 캐시 건너뜀
		if isWriteMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// 캐시 키 생성
		cacheKey := buildCacheKey(ctx, config.KeyPrefix, info.FullMethod, req)

		// 캐시에서 조회
		cached, err := store.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			// 캐시 히트: 응답 역직렬화
			resp, unmarshalErr := unmarshalCachedResponse(cached, req)
			if unmarshalErr == nil {
				return resp, nil
			}
		}

		// 캐시 미스: 핸들러 실행
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		// 응답을 캐시에 저장
		if protoMsg, ok := resp.(proto.Message); ok {
			data, marshalErr := protojson.Marshal(protoMsg)
			if marshalErr == nil {
				_ = store.Set(ctx, cacheKey, string(data), ttl)
			}
		}

		return resp, err
	}
}

// isWriteMethod는 쓰기 메서드인지 확인합니다.
func isWriteMethod(fullMethod string) bool {
	parts := strings.Split(fullMethod, "/")
	if len(parts) < 3 {
		return false
	}
	method := parts[len(parts)-1]
	writePrefixes := []string{"Create", "Update", "Delete", "Set", "Add", "Remove", "Cancel", "Send"}
	for _, prefix := range writePrefixes {
		if strings.HasPrefix(method, prefix) {
			return true
		}
	}
	return false
}

// buildCacheKey는 메서드+요청 기반 캐시 키를 생성합니다.
func buildCacheKey(ctx context.Context, prefix, method string, req interface{}) string {
	// 사용자 ID 포함 (사용자별 캐시 분리)
	userID := ""
	if uid, ok := ctx.Value(UserIDKey).(string); ok {
		userID = uid
	}

	// 요청 데이터 해시
	reqHash := ""
	if protoMsg, ok := req.(proto.Message); ok {
		data, err := protojson.Marshal(protoMsg)
		if err == nil {
			h := sha256.Sum256(data)
			reqHash = hex.EncodeToString(h[:8])
		}
	} else {
		data, err := json.Marshal(req)
		if err == nil {
			h := sha256.Sum256(data)
			reqHash = hex.EncodeToString(h[:8])
		}
	}

	return fmt.Sprintf("%s%s:%s:%s", prefix, method, userID, reqHash)
}

// unmarshalCachedResponse는 캐시된 JSON을 proto 메시지로 변환합니다.
// 참고: 실제 구현 시 각 서비스별 응답 타입 레지스트리가 필요합니다.
func unmarshalCachedResponse(cached string, _ interface{}) (interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CacheInvalidationInterceptor는 쓰기 메서드 실행 후 관련 캐시를 무효화합니다.
func CacheInvalidationInterceptor(store CacheStore, config *CacheConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		// 쓰기 메서드 성공 시 관련 캐시 패턴 무효화
		if isWriteMethod(info.FullMethod) {
			serviceName := extractServiceName(info.FullMethod)
			if serviceName != "" {
				pattern := fmt.Sprintf("%s*%s*", config.KeyPrefix, serviceName)
				_ = invalidatePattern(ctx, store, pattern)
			}
		}

		return resp, err
	}
}

// extractServiceName은 gRPC 메서드에서 서비스명을 추출합니다.
func extractServiceName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// invalidatePattern은 패턴에 매칭되는 캐시를 삭제합니다.
// 참고: 실제로는 Redis SCAN 사용 권장 (Keys는 프로덕션 비권장)
func invalidatePattern(_ context.Context, _ CacheStore, _ string) error {
	// CacheStore 인터페이스에 Keys/Del이 없으므로
	// 실제 구현 시 RedisClient를 직접 사용하거나 인터페이스 확장
	return nil
}

// NoCacheHeader는 클라이언트가 캐시를 건너뛰도록 요청할 때 사용하는 메타데이터 키입니다.
const NoCacheHeader = "x-no-cache"

// HasNoCacheHeader는 요청에 캐시 건너뛰기 헤더가 있는지 확인합니다.
func HasNoCacheHeader(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}
	values := md.Get(NoCacheHeader)
	return len(values) > 0 && values[0] == "true"
}
