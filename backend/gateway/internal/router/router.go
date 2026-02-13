package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/manpasik/backend/shared/storage"
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
	CalibrationAddr  string
	CoachingAddr     string
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
	calibrationAddr  string
	coachingAddr     string
	s3Client         *storage.S3Client // optional: nil이면 파일 업로드 비활성화
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
		calibrationAddr:  cfg.CalibrationAddr,
		coachingAddr:     cfg.CoachingAddr,
	}
	r.setupRoutes()
	return r
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
		"services":    17,
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
