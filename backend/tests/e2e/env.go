package e2e

import "os"

// getEnv returns the environment variable value or the fallback default.
// 기본값은 127.0.0.1 — WSL2 등에서 localhost가 IPv6(::1)로 해석되어 Docker 포트에 연결되지 않는 경우를 피함.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Service address helpers — override via environment variables.

func AuthAddr() string         { return getEnv("AUTH_SERVICE_ADDR", "127.0.0.1:50051") }
func UserAddr() string         { return getEnv("USER_SERVICE_ADDR", "127.0.0.1:50052") }
func DeviceAddr() string       { return getEnv("DEVICE_SERVICE_ADDR", "127.0.0.1:50053") }
func MeasurementAddr() string  { return getEnv("MEASUREMENT_SERVICE_ADDR", "127.0.0.1:50054") }
func SubscriptionAddr() string { return getEnv("SUBSCRIPTION_SERVICE_ADDR", "127.0.0.1:50055") }
func ShopAddr() string         { return getEnv("SHOP_SERVICE_ADDR", "127.0.0.1:50056") }
func PaymentAddr() string      { return getEnv("PAYMENT_SERVICE_ADDR", "127.0.0.1:50057") }
func AiInferenceAddr() string  { return getEnv("AI_INFERENCE_SERVICE_ADDR", "127.0.0.1:50058") }
func CartridgeAddr() string    { return getEnv("CARTRIDGE_SERVICE_ADDR", "127.0.0.1:50059") }
func CalibrationAddr() string  { return getEnv("CALIBRATION_SERVICE_ADDR", "127.0.0.1:50060") }
func CoachingAddr() string     { return getEnv("COACHING_SERVICE_ADDR", "127.0.0.1:50061") }
func FamilyAddr() string       { return getEnv("FAMILY_SERVICE_ADDR", "127.0.0.1:50063") }
func HealthRecordAddr() string { return getEnv("HEALTH_RECORD_SERVICE_ADDR", "127.0.0.1:50064") }
func CommunityAddr() string    { return getEnv("COMMUNITY_SERVICE_ADDR", "127.0.0.1:50065") }
func ReservationAddr() string  { return getEnv("RESERVATION_SERVICE_ADDR", "127.0.0.1:50066") }
func AdminAddr() string        { return getEnv("ADMIN_SERVICE_ADDR", "127.0.0.1:50067") }
func NotificationAddr() string { return getEnv("NOTIFICATION_SERVICE_ADDR", "127.0.0.1:50068") }
func PrescriptionAddr() string { return getEnv("PRESCRIPTION_SERVICE_ADDR", "127.0.0.1:50069") }
func VideoAddr() string        { return getEnv("VIDEO_SERVICE_ADDR", "127.0.0.1:50070") }
func TelemedicineAddr() string { return getEnv("TELEMEDICINE_SERVICE_ADDR", "127.0.0.1:50071") }
func VisionAddr() string       { return getEnv("VISION_SERVICE_ADDR", "127.0.0.1:50072") }
func TranslationAddr() string  { return getEnv("TRANSLATION_SERVICE_ADDR", "127.0.0.1:50073") }
func GatewayAddr() string      { return getEnv("GATEWAY_ADDR", "127.0.0.1:8090") }
