package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/manpasik/backend/gateway/internal/circuit"
	"github.com/manpasik/backend/shared/storage"
	"google.golang.org/grpc"
)

// Config holds all gRPC service addresses for the gateway.
type Config struct {
	AuthAddr         string
	MeasurementAddr  string
	UserAddr         string
	DeviceAddr       string
	ReservationAddr  string
	PrescriptionAddr string
	SubscriptionAddr string
	ShopAddr         string
	PaymentAddr      string
	HealthRecordAddr string
	NotificationAddr string
	CommunityAddr    string
	AdminAddr        string
	AiInferenceAddr  string
	CartridgeAddr    string
	CalibrationAddr    string
	CoachingAddr       string
	IotGatewayAddr     string
	NlpAddr            string
	MarketplaceAddr    string
	AiTrainingAddr     string
	AdAddr             string
	InventoryAddr      string
	LogisticsAddr      string
	LocationAddr       string
	ModerationAddr     string
}

// Router sets up all REST API routes and proxies to gRPC microservices.
type Router struct {
	mux              *http.ServeMux
	authAddr         string
	measurementAddr  string
	userAddr         string
	deviceAddr       string
	reservationAddr  string
	prescriptionAddr string
	subscriptionAddr string
	shopAddr         string
	paymentAddr      string
	healthRecordAddr string
	notificationAddr string
	communityAddr    string
	adminAddr        string
	aiInferenceAddr  string
	cartridgeAddr    string
	calibrationAddr    string
	coachingAddr       string
	iotGatewayAddr     string
	nlpAddr            string
	marketplaceAddr    string
	s3Client           *storage.S3Client // optional: nil이면 파일 업로드 비활성화
	breakers           map[string]*circuit.Breaker // H6: 서비스별 서킷 브레이커
}

// SetS3Client는 S3 클라이언트를 설정합니다 (optional).
func (r *Router) SetS3Client(s3 *storage.S3Client) {
	r.s3Client = s3
}

// NewRouter creates a new REST API router with gRPC backend addresses.
func NewRouter(cfg Config) *Router {
	r := &Router{
		mux:              http.NewServeMux(),
		authAddr:         cfg.AuthAddr,
		measurementAddr:  cfg.MeasurementAddr,
		userAddr:         cfg.UserAddr,
		deviceAddr:       cfg.DeviceAddr,
		reservationAddr:  cfg.ReservationAddr,
		prescriptionAddr: cfg.PrescriptionAddr,
		subscriptionAddr: cfg.SubscriptionAddr,
		shopAddr:         cfg.ShopAddr,
		paymentAddr:      cfg.PaymentAddr,
		healthRecordAddr: cfg.HealthRecordAddr,
		notificationAddr: cfg.NotificationAddr,
		communityAddr:    cfg.CommunityAddr,
		adminAddr:        cfg.AdminAddr,
		aiInferenceAddr:  cfg.AiInferenceAddr,
		cartridgeAddr:    cfg.CartridgeAddr,
		calibrationAddr:    cfg.CalibrationAddr,
		coachingAddr:       cfg.CoachingAddr,
		iotGatewayAddr:     cfg.IotGatewayAddr,
		nlpAddr:            cfg.NlpAddr,
		marketplaceAddr:    cfg.MarketplaceAddr,
		breakers:           make(map[string]*circuit.Breaker),
	}

	// H6: 서비스별 서킷 브레이커 초기화 (5회 실패 → 30초 차단)
	for _, svc := range []string{
		"auth", "user", "device", "measurement",
		"reservation", "prescription", "subscription",
		"shop", "payment", "health-record",
		"notification", "community", "admin",
		"ai-inference", "cartridge", "calibration",
		"coaching", "iot-gateway", "nlp", "marketplace",
	} {
		r.breakers[svc] = circuit.New(svc, 5, 30*time.Second)
	}

	r.setupRoutes()
	return r
}

// dialWithBreaker wraps dialGRPC with circuit breaker protection.
//
// H6: 개별 서비스 장애가 전체 게이트웨이로 전파되지 않도록 격리합니다.
// 서킷이 Open 상태이면 gRPC 연결을 시도하지 않고 즉시 에러를 반환합니다.
func (r *Router) dialWithBreaker(serviceName, addr string) (*grpcConnWithBreaker, error) {
	breaker, ok := r.breakers[serviceName]
	if !ok {
		// 브레이커 미등록 서비스 — 브레이커 없이 연결
		conn, err := dialGRPC(addr)
		if err != nil {
			return nil, err
		}
		return &grpcConnWithBreaker{conn: conn, breaker: nil}, nil
	}

	if !breaker.Allow() {
		return nil, &circuitOpenError{service: serviceName}
	}

	conn, err := dialGRPC(addr)
	if err != nil {
		breaker.RecordFailure()
		return nil, err
	}

	return &grpcConnWithBreaker{conn: conn, breaker: breaker}, nil
}

// grpcConnWithBreaker wraps a gRPC connection with its associated breaker.
type grpcConnWithBreaker struct {
	conn    *grpc.ClientConn
	breaker *circuit.Breaker
}

// Close closes the underlying gRPC connection.
func (c *grpcConnWithBreaker) Close() error {
	return c.conn.Close()
}

// RecordSuccess records a successful gRPC call to the breaker.
func (c *grpcConnWithBreaker) RecordSuccess() {
	if c.breaker != nil {
		c.breaker.RecordSuccess()
	}
}

// RecordFailure records a failed gRPC call to the breaker.
func (c *grpcConnWithBreaker) RecordFailure() {
	if c.breaker != nil {
		c.breaker.RecordFailure()
	}
}

// circuitOpenError is returned when the circuit breaker is open.
type circuitOpenError struct {
	service string
}

func (e *circuitOpenError) Error() string {
	return "서비스 일시 중단 (circuit open): " + e.service
}

// BreakerStats returns diagnostic stats for all circuit breakers.
func (r *Router) BreakerStats() []circuit.BreakerStats {
	stats := make([]circuit.BreakerStats, 0, len(r.breakers))
	for _, b := range r.breakers {
		stats = append(stats, b.Stats())
	}
	return stats
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	r.mux.ServeHTTP(w, req)
}

func (r *Router) setupRoutes() {
	// Health check endpoints
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/health/live", r.handleHealthLive)
	r.mux.HandleFunc("/health/ready", r.handleHealthReady)
	r.mux.HandleFunc("GET /api/v1/version", r.handleVersion)

	// Auth service
	r.mux.HandleFunc("POST /api/v1/auth/register", r.handleRegister)
	r.mux.HandleFunc("POST /api/v1/auth/login", r.handleLogin)
	r.mux.HandleFunc("POST /api/v1/auth/refresh", r.handleRefreshToken)
	r.mux.HandleFunc("POST /api/v1/auth/logout", r.handleLogout)

	// User service
	r.mux.HandleFunc("GET /api/v1/users/{userId}/profile", r.handleGetProfile)
	r.mux.HandleFunc("PUT /api/v1/users/{userId}/profile", r.handleUpdateProfile)

	// Measurement service
	r.mux.HandleFunc("POST /api/v1/measurements/sessions", r.handleStartSession)
	r.mux.HandleFunc("POST /api/v1/measurements/sessions/{sessionId}/end", r.handleEndSession)
	r.mux.HandleFunc("GET /api/v1/measurements/history", r.handleGetHistory)

	// Device service
	r.mux.HandleFunc("POST /api/v1/devices", r.handleRegisterDevice)
	r.mux.HandleFunc("GET /api/v1/devices", r.handleListDevices)

	// Reservation service (facilities + reservations)
	r.mux.HandleFunc("GET /api/v1/facilities", r.handleSearchFacilities)
	r.mux.HandleFunc("GET /api/v1/facilities/{facilityId}", r.handleGetFacility)
	r.mux.HandleFunc("POST /api/v1/reservations", r.handleCreateReservation)
	r.mux.HandleFunc("GET /api/v1/reservations", r.handleListReservations)
	r.mux.HandleFunc("GET /api/v1/reservations/{reservationId}", r.handleGetReservation)

	// Prescription service
	r.mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/pharmacy", r.handleSelectPharmacy)
	r.mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/send", r.handleSendToPharmacy)
	r.mux.HandleFunc("GET /api/v1/prescriptions/token/{token}", r.handleGetByToken)

	// Subscription service
	r.mux.HandleFunc("GET /api/v1/subscriptions/plans", r.handleListSubscriptionPlans)
	r.mux.HandleFunc("GET /api/v1/subscriptions/{userId}", r.handleGetSubscription)
	r.mux.HandleFunc("POST /api/v1/subscriptions", r.handleCreateSubscription)
	r.mux.HandleFunc("DELETE /api/v1/subscriptions/{subscriptionId}", r.handleCancelSubscription)

	// Shop service
	r.mux.HandleFunc("GET /api/v1/products", r.handleListProducts)
	r.mux.HandleFunc("GET /api/v1/products/{productId}", r.handleGetProduct)
	r.mux.HandleFunc("POST /api/v1/cart", r.handleAddToCart)
	r.mux.HandleFunc("GET /api/v1/cart/{userId}", r.handleGetCart)
	r.mux.HandleFunc("POST /api/v1/orders", r.handleCreateOrder)
	r.mux.HandleFunc("GET /api/v1/orders", r.handleListOrders)

	// Payment service
	r.mux.HandleFunc("POST /api/v1/payments", r.handleCreatePayment)
	r.mux.HandleFunc("POST /api/v1/payments/{paymentId}/confirm", r.handleConfirmPayment)
	r.mux.HandleFunc("GET /api/v1/payments/{paymentId}", r.handleGetPayment)

	// Health Record service
	r.mux.HandleFunc("POST /api/v1/health-records/export/fhir", r.handleExportToFHIR)
	r.mux.HandleFunc("POST /api/v1/health-records", r.handleCreateRecord)
	r.mux.HandleFunc("GET /api/v1/health-records", r.handleListRecords)
	r.mux.HandleFunc("GET /api/v1/health-records/{recordId}", r.handleGetRecord)

	// Notification service
	r.mux.HandleFunc("GET /api/v1/notifications/unread-count", r.handleGetUnreadCount)
	r.mux.HandleFunc("GET /api/v1/notifications", r.handleListNotifications)
	r.mux.HandleFunc("POST /api/v1/notifications/{notificationId}/read", r.handleMarkAsRead)

	// Community service
	r.mux.HandleFunc("GET /api/v1/posts", r.handleListPosts)
	r.mux.HandleFunc("POST /api/v1/posts", r.handleCreatePost)
	r.mux.HandleFunc("GET /api/v1/posts/{postId}", r.handleGetPost)
	r.mux.HandleFunc("POST /api/v1/posts/{postId}/like", r.handleLikePost)

	// Admin service
	r.mux.HandleFunc("GET /api/v1/admin/stats", r.handleGetSystemStats)
	r.mux.HandleFunc("GET /api/v1/admin/users", r.handleAdminListUsers)
	r.mux.HandleFunc("GET /api/v1/admin/audit-log", r.handleGetAuditLog)

	// AI Inference service
	r.mux.HandleFunc("POST /api/v1/ai/analyze", r.handleAnalyzeMeasurement)
	r.mux.HandleFunc("GET /api/v1/ai/health-score/{userId}", r.handleGetHealthScore)
	r.mux.HandleFunc("POST /api/v1/ai/predict-trend", r.handlePredictTrend)
	r.mux.HandleFunc("GET /api/v1/ai/models", r.handleListAiModels)

	// Cartridge service
	r.mux.HandleFunc("POST /api/v1/cartridges/read", r.handleReadCartridge)
	r.mux.HandleFunc("POST /api/v1/cartridges/usage", r.handleRecordCartridgeUsage)
	r.mux.HandleFunc("GET /api/v1/cartridges/types", r.handleListCartridgeCategories)
	r.mux.HandleFunc("GET /api/v1/cartridges/{cartridgeId}/remaining", r.handleGetRemainingUses)
	r.mux.HandleFunc("POST /api/v1/cartridges/validate", r.handleValidateCartridge)

	// Calibration service
	r.mux.HandleFunc("POST /api/v1/calibration/factory", r.handleRegisterFactoryCalibration)
	r.mux.HandleFunc("POST /api/v1/calibration/field", r.handlePerformFieldCalibration)
	r.mux.HandleFunc("GET /api/v1/calibration/{deviceId}/status", r.handleCheckCalibrationStatus)
	r.mux.HandleFunc("GET /api/v1/calibration/models", r.handleListCalibrationModels)

	// Coaching service
	r.mux.HandleFunc("POST /api/v1/coaching/goals", r.handleSetHealthGoal)
	r.mux.HandleFunc("GET /api/v1/coaching/goals/{userId}", r.handleGetHealthGoals)
	r.mux.HandleFunc("POST /api/v1/coaching/generate", r.handleGenerateCoaching)
	r.mux.HandleFunc("GET /api/v1/coaching/daily-report/{userId}", r.handleGenerateDailyReport)
	r.mux.HandleFunc("GET /api/v1/coaching/recommendations/{userId}", r.handleGetRecommendations)

	// IoT Gateway service
	r.mux.HandleFunc("GET /api/v1/iot/devices", r.handleListIoTDevices)
	r.mux.HandleFunc("GET /api/v1/iot/telemetry", r.handleGetIoTTelemetry)

	// NLP service
	r.mux.HandleFunc("POST /api/v1/nlp/parse", r.handleNlpParse)
	r.mux.HandleFunc("POST /api/v1/nlp/symptoms", r.handleNlpSymptoms)

	// Marketplace service
	r.mux.HandleFunc("GET /api/v1/marketplace/products", r.handleListMarketplaceProducts)
	r.mux.HandleFunc("GET /api/v1/marketplace/partners", r.handleListMarketplacePartners)

	// File upload service (S3/MinIO)
	r.setupUploadRoutes()
}

// --- Health / Version handlers ---

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "gateway",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func (r *Router) handleHealthLive(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"alive": true})
}

func (r *Router) handleHealthReady(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"ready": true})
}

func (r *Router) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"service":     "gateway",
		"version":     "1.0.0",
		"api_version": "v1",
		"services":    20,
	})
}

// --- JSON helpers ---

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
