package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// --- gRPC dial helpers ---

func dialGRPC(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// ============================================================================
// User Service Handlers
// ============================================================================

// GET /api/v1/users/{userId}/profile
func (r *Router) handleGetProfile(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.userAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "사용자 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProfile(ctx, &v1.GetProfileRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":           resp.GetUserId(),
		"email":             resp.GetEmail(),
		"display_name":      resp.GetDisplayName(),
		"avatar_url":        resp.GetAvatarUrl(),
		"language":          resp.GetLanguage(),
		"timezone":          resp.GetTimezone(),
		"subscription_tier": resp.GetSubscriptionTier().String(),
		"created_at":        resp.GetCreatedAt().AsTime().Format(time.RFC3339),
	})
}

// PUT /api/v1/users/{userId}/profile
func (r *Router) handleUpdateProfile(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	var body struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
		Language    string `json:"language"`
		Timezone    string `json:"timezone"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.userAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "사용자 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.UpdateProfile(ctx, &v1.UpdateProfileRequest{
		UserId:      userID,
		DisplayName: body.DisplayName,
		AvatarUrl:   body.AvatarURL,
		Language:    body.Language,
		Timezone:    body.Timezone,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":           resp.GetUserId(),
		"email":             resp.GetEmail(),
		"display_name":      resp.GetDisplayName(),
		"avatar_url":        resp.GetAvatarUrl(),
		"language":          resp.GetLanguage(),
		"timezone":          resp.GetTimezone(),
		"subscription_tier": resp.GetSubscriptionTier().String(),
		"created_at":        resp.GetCreatedAt().AsTime().Format(time.RFC3339),
	})
}

// ============================================================================
// Measurement Service Handlers
// ============================================================================

// POST /api/v1/measurements/sessions
func (r *Router) handleStartSession(w http.ResponseWriter, req *http.Request) {
	var body struct {
		DeviceID           string `json:"device_id"`
		CartridgeID        string `json:"cartridge_id"`
		UserID             string `json:"user_id"`
		CartridgeCategory  int32  `json:"cartridge_category"`
		CartridgeTypeIndex int32  `json:"cartridge_type_index"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.measurementAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "측정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewMeasurementServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.StartSession(ctx, &v1.StartSessionRequest{
		DeviceId:           body.DeviceID,
		CartridgeId:        body.CartridgeID,
		UserId:             body.UserID,
		CartridgeCategory:  body.CartridgeCategory,
		CartridgeTypeIndex: body.CartridgeTypeIndex,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"session_id": resp.GetSessionId(),
		"started_at": resp.GetStartedAt().AsTime().Format(time.RFC3339),
	})
}

// POST /api/v1/measurements/sessions/{sessionId}/end
func (r *Router) handleEndSession(w http.ResponseWriter, req *http.Request) {
	sessionID := req.PathValue("sessionId")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.measurementAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "측정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewMeasurementServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.EndSession(ctx, &v1.EndSessionRequest{
		SessionId: sessionID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":         resp.GetSessionId(),
		"total_measurements": resp.GetTotalMeasurements(),
		"ended_at":           resp.GetEndedAt().AsTime().Format(time.RFC3339),
	})
}

// GET /api/v1/measurements/history?user_id=...&limit=...&offset=...
func (r *Router) handleGetHistory(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.measurementAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "측정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewMeasurementServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetMeasurementHistory(ctx, &v1.GetHistoryRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	measurements := make([]map[string]interface{}, 0, len(resp.GetMeasurements()))
	for _, m := range resp.GetMeasurements() {
		measurements = append(measurements, map[string]interface{}{
			"session_id":     m.GetSessionId(),
			"cartridge_type": m.GetCartridgeType(),
			"primary_value":  m.GetPrimaryValue(),
			"unit":           m.GetUnit(),
			"measured_at":    m.GetMeasuredAt().AsTime().Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"measurements": measurements,
		"total_count":  resp.GetTotalCount(),
	})
}

// ============================================================================
// Device Service Handlers
// ============================================================================

// POST /api/v1/devices
func (r *Router) handleRegisterDevice(w http.ResponseWriter, req *http.Request) {
	var body struct {
		DeviceID        string `json:"device_id"`
		SerialNumber    string `json:"serial_number"`
		FirmwareVersion string `json:"firmware_version"`
		UserID          string `json:"user_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.deviceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "디바이스 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewDeviceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.RegisterDevice(ctx, &v1.RegisterDeviceRequest{
		DeviceId:        body.DeviceID,
		SerialNumber:    body.SerialNumber,
		FirmwareVersion: body.FirmwareVersion,
		UserId:          body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"device_id":          resp.GetDeviceId(),
		"registration_token": resp.GetRegistrationToken(),
		"registered_at":      resp.GetRegisteredAt().AsTime().Format(time.RFC3339),
	})
}

// GET /api/v1/devices?user_id=...
func (r *Router) handleListDevices(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.deviceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "디바이스 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewDeviceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListDevices(ctx, &v1.ListDevicesRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	devices := make([]map[string]interface{}, 0, len(resp.GetDevices()))
	for _, d := range resp.GetDevices() {
		devices = append(devices, map[string]interface{}{
			"device_id":        d.GetDeviceId(),
			"name":             d.GetName(),
			"firmware_version": d.GetFirmwareVersion(),
			"status":           d.GetStatus().String(),
			"battery_percent":  d.GetBatteryPercent(),
			"last_seen":        d.GetLastSeen().AsTime().Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"devices": devices,
	})
}

// ============================================================================
// Reservation Service Handlers
// ============================================================================

// GET /api/v1/facilities?query=...&latitude=...&longitude=...&radius_km=...&limit=...&offset=...
func (r *Router) handleSearchFacilities(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	lat, _ := strconv.ParseFloat(q.Get("latitude"), 64)
	lon, _ := strconv.ParseFloat(q.Get("longitude"), 64)
	radiusKm, _ := strconv.ParseFloat(q.Get("radius_km"), 64)
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.reservationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "예약 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewReservationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SearchFacilities(ctx, &v1.SearchFacilitiesRequest{
		Latitude:  lat,
		Longitude: lon,
		RadiusKm:  radiusKm,
		Query:     q.Get("query"),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	facilities := make([]map[string]interface{}, 0, len(resp.GetFacilities()))
	for _, f := range resp.GetFacilities() {
		facilities = append(facilities, map[string]interface{}{
			"facility_id":         f.GetFacilityId(),
			"name":                f.GetName(),
			"type":                f.GetType().String(),
			"address":             f.GetAddress(),
			"latitude":            f.GetLatitude(),
			"longitude":           f.GetLongitude(),
			"phone":               f.GetPhone(),
			"rating":              f.GetRating(),
			"is_open_now":         f.GetIsOpenNow(),
			"accepts_reservation": f.GetAcceptsReservation(),
			"operating_hours":     f.GetOperatingHours(),
			"distance_km":         f.GetDistanceKm(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"facilities":  facilities,
		"total_count": resp.GetTotalCount(),
	})
}

// GET /api/v1/facilities/{facilityId}
func (r *Router) handleGetFacility(w http.ResponseWriter, req *http.Request) {
	facilityID := req.PathValue("facilityId")
	if facilityID == "" {
		writeError(w, http.StatusBadRequest, "facility_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.reservationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "예약 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewReservationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f, err := client.GetFacility(ctx, &v1.GetFacilityRequest{
		FacilityId: facilityID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"facility_id":         f.GetFacilityId(),
		"name":                f.GetName(),
		"type":                f.GetType().String(),
		"address":             f.GetAddress(),
		"latitude":            f.GetLatitude(),
		"longitude":           f.GetLongitude(),
		"phone":               f.GetPhone(),
		"rating":              f.GetRating(),
		"is_open_now":         f.GetIsOpenNow(),
		"accepts_reservation": f.GetAcceptsReservation(),
		"operating_hours":     f.GetOperatingHours(),
		"distance_km":         f.GetDistanceKm(),
	})
}

// POST /api/v1/reservations
func (r *Router) handleCreateReservation(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID     string `json:"user_id"`
		FacilityID string `json:"facility_id"`
		SlotID     string `json:"slot_id"`
		DoctorID   string `json:"doctor_id"`
		Specialty  int32  `json:"specialty"`
		Reason     string `json:"reason"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.reservationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "예약 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewReservationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateReservation(ctx, &v1.CreateReservationRequest{
		UserId:     body.UserID,
		FacilityId: body.FacilityID,
		SlotId:     body.SlotID,
		DoctorId:   body.DoctorID,
		Specialty:  v1.DoctorSpecialty(body.Specialty),
		Reason:     body.Reason,
		Notes:      body.Notes,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, reservationToMap(resp))
}

// GET /api/v1/reservations?user_id=...&limit=...&offset=...
func (r *Router) handleListReservations(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.reservationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "예약 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewReservationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ListReservations(ctx, &v1.ListReservationsRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	reservations := make([]map[string]interface{}, 0, len(resp.GetReservations()))
	for _, rv := range resp.GetReservations() {
		reservations = append(reservations, reservationToMap(rv))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reservations": reservations,
		"total_count":  resp.GetTotalCount(),
	})
}

// GET /api/v1/reservations/{reservationId}
func (r *Router) handleGetReservation(w http.ResponseWriter, req *http.Request) {
	reservationID := req.PathValue("reservationId")
	if reservationID == "" {
		writeError(w, http.StatusBadRequest, "reservation_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.reservationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "예약 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewReservationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetReservation(ctx, &v1.GetReservationRequest{
		ReservationId: reservationID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, reservationToMap(resp))
}

func reservationToMap(rv *v1.Reservation) map[string]interface{} {
	m := map[string]interface{}{
		"reservation_id": rv.GetReservationId(),
		"user_id":        rv.GetUserId(),
		"facility_id":    rv.GetFacilityId(),
		"facility_name":  rv.GetFacilityName(),
		"doctor_id":      rv.GetDoctorId(),
		"doctor_name":    rv.GetDoctorName(),
		"specialty":      rv.GetSpecialty().String(),
		"status":         rv.GetStatus().String(),
		"reason":         rv.GetReason(),
		"notes":          rv.GetNotes(),
	}
	if rv.GetAppointmentTime() != nil {
		m["appointment_time"] = rv.GetAppointmentTime().AsTime().Format(time.RFC3339)
	}
	if rv.GetCreatedAt() != nil {
		m["created_at"] = rv.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if rv.GetUpdatedAt() != nil {
		m["updated_at"] = rv.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Prescription Service Handlers
// ============================================================================

// POST /api/v1/prescriptions/{prescriptionId}/pharmacy
func (r *Router) handleSelectPharmacy(w http.ResponseWriter, req *http.Request) {
	prescriptionID := req.PathValue("prescriptionId")
	if prescriptionID == "" {
		writeError(w, http.StatusBadRequest, "prescription_id가 필요합니다")
		return
	}

	var body struct {
		PharmacyID      string `json:"pharmacy_id"`
		PharmacyName    string `json:"pharmacy_name"`
		FulfillmentType string `json:"fulfillment_type"`
		ShippingAddress string `json:"shipping_address"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.prescriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "처방 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPrescriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SelectPharmacyAndFulfillment(ctx, &v1.SelectPharmacyRequest{
		PrescriptionId:  prescriptionID,
		PharmacyId:      body.PharmacyID,
		PharmacyName:    body.PharmacyName,
		FulfillmentType: body.FulfillmentType,
		ShippingAddress: body.ShippingAddress,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
	})
}

// POST /api/v1/prescriptions/{prescriptionId}/send
func (r *Router) handleSendToPharmacy(w http.ResponseWriter, req *http.Request) {
	prescriptionID := req.PathValue("prescriptionId")
	if prescriptionID == "" {
		writeError(w, http.StatusBadRequest, "prescription_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.prescriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "처방 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPrescriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SendPrescriptionToPharmacy(ctx, &v1.SendToPharmacyRequest{
		PrescriptionId: prescriptionID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"fulfillment_token": resp.GetFulfillmentToken(),
		"expires_at":        resp.GetExpiresAt(),
		"success":           resp.GetSuccess(),
	})
}

// GET /api/v1/prescriptions/token/{token}
func (r *Router) handleGetByToken(w http.ResponseWriter, req *http.Request) {
	token := req.PathValue("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token이 필요합니다")
		return
	}

	conn, err := dialGRPC(r.prescriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "처방 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPrescriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPrescriptionByToken(ctx, &v1.GetByTokenRequest{
		FulfillmentToken: token,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Build medications list
	meds := make([]map[string]interface{}, 0, len(resp.GetMedications()))
	for _, med := range resp.GetMedications() {
		meds = append(meds, map[string]interface{}{
			"medication_id": med.GetMedicationId(),
			"name":          med.GetName(),
			"dosage":        med.GetDosage(),
			"frequency":     med.GetFrequency(),
			"route":         med.GetRoute(),
			"duration_days": med.GetDurationDays(),
			"instructions":  med.GetInstructions(),
			"is_critical":   med.GetIsCritical(),
		})
	}

	result := map[string]interface{}{
		"prescription_id": resp.GetPrescriptionId(),
		"user_id":         resp.GetUserId(),
		"doctor_id":       resp.GetDoctorId(),
		"doctor_name":     resp.GetDoctorName(),
		"facility_id":     resp.GetFacilityId(),
		"diagnosis":       resp.GetDiagnosis(),
		"notes":           resp.GetNotes(),
		"status":          resp.GetStatus().String(),
		"medications":     meds,
	}
	if resp.GetPrescribedAt() != nil {
		result["prescribed_at"] = resp.GetPrescribedAt().AsTime().Format(time.RFC3339)
	}
	if resp.GetExpiresAt() != nil {
		result["expires_at"] = resp.GetExpiresAt().AsTime().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, result)
}
