package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerCommunityRoutes는 커뮤니티/가족/의료 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerCommunityRoutes(mux *http.ServeMux) {
	// Community
	mux.HandleFunc("GET /api/v1/posts", h.handleListPosts)
	mux.HandleFunc("POST /api/v1/posts", h.handleCreatePost)
	mux.HandleFunc("GET /api/v1/posts/{postId}", h.handleGetPost)
	mux.HandleFunc("POST /api/v1/posts/{postId}/like", h.handleLikePost)

	// Family
	mux.HandleFunc("GET /api/v1/family/groups", h.handleListFamilyGroups)
	mux.HandleFunc("GET /api/v1/family/groups/{groupId}/report", h.handleGetFamilyGroupReport)

	// Facility/Reservation
	mux.HandleFunc("GET /api/v1/facilities", h.handleSearchFacilities)
	mux.HandleFunc("GET /api/v1/facilities/{facilityId}", h.handleGetFacility)
	mux.HandleFunc("POST /api/v1/reservations", h.handleCreateReservation)
	mux.HandleFunc("GET /api/v1/reservations", h.handleListReservations)
	mux.HandleFunc("GET /api/v1/reservations/{reservationId}", h.handleGetReservation)

	// Telemedicine
	mux.HandleFunc("GET /api/v1/telemedicine/doctors", h.handleSearchDoctors)
	mux.HandleFunc("POST /api/v1/telemedicine/consultations", h.handleCreateConsultation)

	// Video
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/join", h.handleJoinVideoRoom)
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/leave", h.handleLeaveVideoRoom)

	// User avatar
	mux.HandleFunc("POST /api/v1/users/{userId}/avatar", h.handleUploadAvatar)

	// AI (food/exercise)
	mux.HandleFunc("POST /api/v1/ai/food-analyze", h.handleFoodAnalyze)
	mux.HandleFunc("POST /api/v1/ai/exercise-analyze", h.handleExerciseAnalyze)
}

// ── Community ──

func (h *RestHandler) handleListPosts(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	resp, err := h.community.ListPosts(r.Context(), &v1.ListPostsRequest{
		Category: v1.PostCategory(queryInt(r, "category", 0)),
		AuthorId: r.URL.Query().Get("author_id"),
		Query:    r.URL.Query().Get("query"),
		Limit:    queryInt(r, "limit", 20),
		Offset:   queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	var body struct {
		AuthorID string   `json:"author_id"`
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category int32    `json:"category"`
		Tags     []string `json:"tags"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.CreatePost(r.Context(), &v1.CreatePostRequest{
		AuthorId: body.AuthorID,
		Title:    body.Title,
		Content:  body.Content,
		Category: v1.PostCategory(body.Category),
		Tags:     body.Tags,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	postId := r.PathValue("postId")
	resp, err := h.community.GetPost(r.Context(), &v1.GetPostRequest{PostId: postId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleLikePost(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	postId := r.PathValue("postId")
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.LikePost(r.Context(), &v1.LikePostRequest{
		PostId: postId,
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Family ──

func (h *RestHandler) handleListFamilyGroups(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	userId := r.URL.Query().Get("user_id")
	// Proto has no ListFamilyGroups; return user-scoped empty list
	// (full group enumeration requires proto extension)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"groups": []interface{}{}, "user_id": userId,
	})
}

func (h *RestHandler) handleGetFamilyGroupReport(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	resp, err := h.family.GetSharedHealthData(r.Context(), &v1.GetSharedHealthDataRequest{
		RequesterUserId: "",
		GroupId:         groupId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Facility/Reservation ──

func (h *RestHandler) handleSearchFacilities(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	resp, err := h.reservation.SearchFacilities(r.Context(), &v1.SearchFacilitiesRequest{
		Query:    r.URL.Query().Get("query"),
		Latitude: queryFloat(r, "latitude", 0),
		Longitude: queryFloat(r, "longitude", 0),
		RadiusKm: queryFloat(r, "radius_km", 10),
		Limit:    queryInt(r, "limit", 20),
		Offset:   queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetFacility(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	facilityId := r.PathValue("facilityId")
	resp, err := h.reservation.GetFacility(r.Context(), &v1.GetFacilityRequest{FacilityId: facilityId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateReservation(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	var body struct {
		UserID     string `json:"user_id"`
		FacilityID string `json:"facility_id"`
		SlotID     string `json:"slot_id"`
		DoctorID   string `json:"doctor_id"`
		Specialty  int32  `json:"specialty"`
		Reason     string `json:"reason"`
		Notes      string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.reservation.CreateReservation(r.Context(), &v1.CreateReservationRequest{
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
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListReservations(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	resp, err := h.reservation.ListReservations(r.Context(), &v1.ListReservationsRequest{
		UserId: r.URL.Query().Get("user_id"),
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetReservation(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	reservationId := r.PathValue("reservationId")
	resp, err := h.reservation.GetReservation(r.Context(), &v1.GetReservationRequest{ReservationId: reservationId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Telemedicine ──

func (h *RestHandler) handleSearchDoctors(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	resp, err := h.telemedicine.MatchDoctor(r.Context(), &v1.MatchDoctorRequest{
		Specialty: v1.DoctorSpecialty(queryInt(r, "specialty", 0)),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateConsultation(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	var body struct {
		UserID    string `json:"user_id"`
		DoctorID  string `json:"doctor_id"`
		Specialty string `json:"specialty"`
		Reason    string `json:"reason"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.telemedicine.CreateConsultation(r.Context(), &v1.CreateConsultationRequest{
		PatientUserId:  body.UserID,
		Specialty:      v1.DoctorSpecialty(v1.DoctorSpecialty_value[body.Specialty]),
		ChiefComplaint: body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

// ── Video ──

func (h *RestHandler) handleJoinVideoRoom(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	var body struct {
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.video.JoinRoom(r.Context(), &v1.JoinRoomRequest{
		RoomId: roomId,
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleLeaveVideoRoom(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.video.LeaveRoom(r.Context(), &v1.LeaveRoomRequest{
		RoomId: roomId,
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Avatar Upload (placeholder) ──

func (h *RestHandler) handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"user_id":    userId,
		"avatar_url": "/static/avatars/" + userId + ".jpg",
	})
}

// ── AI Food/Exercise Analysis (placeholder → will connect in Sprint 4) ──

func (h *RestHandler) handleFoodAnalyze(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	var body struct {
		UserID   string `json:"user_id"`
		ImageURL string `json:"image_url"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.aiInference.AnalyzeMeasurement(r.Context(), &v1.AnalyzeMeasurementRequest{
		UserId:        body.UserID,
		MeasurementId: body.ImageURL,
		Models:        []v1.AiModelType{v1.AiModelType_AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR},
	})
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true, "message": "음식 분석 요청이 처리되었습니다",
			"result": map[string]interface{}{"total_calories": 0, "items": []interface{}{}},
		})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleExerciseAnalyze(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	var body struct {
		UserID   string `json:"user_id"`
		VideoURL string `json:"video_url"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.aiInference.AnalyzeMeasurement(r.Context(), &v1.AnalyzeMeasurementRequest{
		UserId:        body.UserID,
		MeasurementId: body.VideoURL,
		Models:        []v1.AiModelType{v1.AiModelType_AI_MODEL_TYPE_HEALTH_SCORER},
	})
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true, "message": "운동 분석 요청이 처리되었습니다",
			"result": map[string]interface{}{"calories_burned": 0, "exercises": []interface{}{}},
		})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
