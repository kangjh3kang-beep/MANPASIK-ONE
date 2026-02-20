package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerCommunityRoutes는 커뮤니티/가족/의료 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerCommunityRoutes(mux *http.ServeMux) {
	// Community – Posts
	mux.HandleFunc("GET /api/v1/posts", h.handleListPosts)
	mux.HandleFunc("POST /api/v1/posts", h.handleCreatePost)
	mux.HandleFunc("GET /api/v1/posts/{postId}", h.handleGetPost)
	mux.HandleFunc("POST /api/v1/posts/{postId}/like", h.handleLikePost)
	mux.HandleFunc("POST /api/v1/posts/{postId}/comments", h.handleCreateComment)
	mux.HandleFunc("GET /api/v1/posts/{postId}/comments", h.handleListComments)

	// Community – Challenges
	mux.HandleFunc("POST /api/v1/challenges", h.handleCreateChallenge)
	mux.HandleFunc("GET /api/v1/challenges", h.handleListChallenges)
	mux.HandleFunc("GET /api/v1/challenges/{challengeId}", h.handleGetChallenge)
	mux.HandleFunc("POST /api/v1/challenges/{challengeId}/join", h.handleJoinChallenge)
	mux.HandleFunc("GET /api/v1/challenges/{challengeId}/leaderboard", h.handleGetChallengeLeaderboard)
	mux.HandleFunc("POST /api/v1/challenges/{challengeId}/progress", h.handleUpdateChallengeProgress)

	// Family
	mux.HandleFunc("POST /api/v1/family/groups", h.handleCreateFamilyGroup)
	mux.HandleFunc("GET /api/v1/family/groups", h.handleListFamilyGroups)
	mux.HandleFunc("GET /api/v1/family/groups/{groupId}", h.handleGetFamilyGroup)
	mux.HandleFunc("GET /api/v1/family/groups/{groupId}/report", h.handleGetFamilyGroupReport)
	mux.HandleFunc("GET /api/v1/family/groups/{groupId}/members", h.handleListFamilyMembers)
	mux.HandleFunc("POST /api/v1/family/groups/{groupId}/invite", h.handleInviteMember)
	mux.HandleFunc("DELETE /api/v1/family/groups/{groupId}/members/{userId}", h.handleRemoveFamilyMember)
	mux.HandleFunc("PUT /api/v1/family/groups/{groupId}/members/{userId}/role", h.handleUpdateMemberRole)
	mux.HandleFunc("PUT /api/v1/family/groups/{groupId}/sharing", h.handleSetSharingPreferences)
	mux.HandleFunc("POST /api/v1/family/invitations/{invitationId}/respond", h.handleRespondToInvitation)
	mux.HandleFunc("POST /api/v1/family/validate-access", h.handleValidateSharingAccess)

	// Facility/Reservation
	mux.HandleFunc("GET /api/v1/facilities", h.handleSearchFacilities)
	mux.HandleFunc("GET /api/v1/facilities/{facilityId}", h.handleGetFacility)
	mux.HandleFunc("GET /api/v1/facilities/{facilityId}/doctors", h.handleListDoctorsByFacility)
	mux.HandleFunc("GET /api/v1/doctors/{doctorId}/availability", h.handleGetDoctorAvailability)
	mux.HandleFunc("POST /api/v1/reservations", h.handleCreateReservation)
	mux.HandleFunc("GET /api/v1/reservations", h.handleListReservations)
	mux.HandleFunc("GET /api/v1/reservations/{reservationId}", h.handleGetReservation)
	mux.HandleFunc("GET /api/v1/reservations/{reservationId}/slots", h.handleGetAvailableSlots)
	mux.HandleFunc("DELETE /api/v1/reservations/{reservationId}", h.handleCancelReservation)
	mux.HandleFunc("POST /api/v1/reservations/{reservationId}/doctor", h.handleSelectDoctor)

	// Family – Guardian Dashboard
	mux.HandleFunc("GET /api/v1/family/groups/{groupId}/guardian-dashboard", h.handleGetGuardianDashboard)

	// Telemedicine
	mux.HandleFunc("GET /api/v1/telemedicine/doctors", h.handleSearchDoctors)
	mux.HandleFunc("POST /api/v1/telemedicine/consultations", h.handleCreateConsultation)
	mux.HandleFunc("GET /api/v1/telemedicine/consultations", h.handleListConsultations)
	mux.HandleFunc("GET /api/v1/telemedicine/consultations/{consultationId}", h.handleGetConsultation)
	mux.HandleFunc("POST /api/v1/telemedicine/consultations/{consultationId}/video", h.handleStartVideoSession)
	mux.HandleFunc("POST /api/v1/telemedicine/consultations/{consultationId}/end-video", h.handleEndVideoSession)
	mux.HandleFunc("POST /api/v1/telemedicine/consultations/{consultationId}/rate", h.handleRateConsultation)

	// Video
	mux.HandleFunc("POST /api/v1/video/rooms", h.handleCreateVideoRoom)
	mux.HandleFunc("GET /api/v1/video/rooms/{roomId}", h.handleGetVideoRoom)
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/join", h.handleJoinVideoRoom)
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/leave", h.handleLeaveVideoRoom)
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/end", h.handleEndVideoRoom)
	mux.HandleFunc("POST /api/v1/video/rooms/{roomId}/signal", h.handleSendSignal)
	mux.HandleFunc("GET /api/v1/video/rooms/{roomId}/participants", h.handleListParticipants)
	mux.HandleFunc("GET /api/v1/video/rooms/{roomId}/stats", h.handleGetRoomStats)

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
	// Proto에 ListFamilyGroups RPC가 없으므로
	// GetSharedHealthData로 사용자 소속 그룹 데이터를 조회
	resp, err := h.family.GetSharedHealthData(r.Context(), &v1.GetSharedHealthDataRequest{
		RequesterUserId: userId,
	})
	if err != nil {
		// 그룹 없는 경우 빈 목록 반환
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"groups": []interface{}{}, "user_id": userId,
		})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
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

// ── Avatar Upload ──

func (h *RestHandler) handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	if h.user == nil {
		writeError(w, http.StatusServiceUnavailable, "user service unavailable")
		return
	}
	userId := r.PathValue("userId")
	// 멀티파트 폼에서 avatar URL 추출 (base64 or URL)
	var avatarURL string
	if r.Header.Get("Content-Type") == "application/json" {
		var body struct {
			AvatarURL string `json:"avatar_url"`
		}
		if err := readJSON(r, &body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		avatarURL = body.AvatarURL
	} else {
		avatarURL = "/static/avatars/" + userId + ".jpg"
	}
	resp, err := h.user.UpdateProfile(r.Context(), &v1.UpdateProfileRequest{
		UserId:    userId,
		AvatarUrl: avatarURL,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── AI Food/Exercise Analysis ──

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
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Community – Comments ──

func (h *RestHandler) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	postId := r.PathValue("postId")
	var body struct {
		AuthorID        string `json:"author_id"`
		Content         string `json:"content"`
		ParentCommentID string `json:"parent_comment_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.CreateComment(r.Context(), &v1.CreateCommentRequest{
		PostId:          postId,
		AuthorId:        body.AuthorID,
		Content:         body.Content,
		ParentCommentId: body.ParentCommentID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListComments(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	postId := r.PathValue("postId")
	resp, err := h.community.ListComments(r.Context(), &v1.ListCommentsRequest{
		PostId: postId,
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Community – Challenges ──

func (h *RestHandler) handleCreateChallenge(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	var body struct {
		CreatorID       string `json:"creator_id"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		ChallengeType   int32  `json:"challenge_type"`
		TargetValue     float64 `json:"target_value"`
		Unit            string `json:"unit"`
		StartDate       string `json:"start_date"`
		EndDate         string `json:"end_date"`
		MaxParticipants int32  `json:"max_participants"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.CreateChallenge(r.Context(), &v1.CreateChallengeRequest{
		CreatorId:       body.CreatorID,
		Title:           body.Title,
		Description:     body.Description,
		ChallengeType:   v1.ChallengeType(body.ChallengeType),
		TargetValue:     body.TargetValue,
		Unit:            body.Unit,
		MaxParticipants: body.MaxParticipants,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetChallenge(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	challengeId := r.PathValue("challengeId")
	resp, err := h.community.GetChallenge(r.Context(), &v1.GetChallengeRequest{
		ChallengeId: challengeId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleJoinChallenge(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	challengeId := r.PathValue("challengeId")
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.JoinChallenge(r.Context(), &v1.JoinChallengeRequest{
		ChallengeId: challengeId,
		UserId:      body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListChallenges(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	resp, err := h.community.ListChallenges(r.Context(), &v1.ListChallengesRequest{
		TypeFilter:   v1.ChallengeType(queryInt(r, "type", 0)),
		StatusFilter: v1.ChallengeStatus(queryInt(r, "status", 0)),
		Limit:        queryInt(r, "limit", 20),
		Offset:       queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetChallengeLeaderboard(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	challengeId := r.PathValue("challengeId")
	resp, err := h.community.GetChallengeLeaderboard(r.Context(), &v1.GetChallengeLeaderboardRequest{
		ChallengeId: challengeId,
		Limit:       queryInt(r, "limit", 20),
		Offset:      queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateChallengeProgress(w http.ResponseWriter, r *http.Request) {
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	challengeId := r.PathValue("challengeId")
	var body struct {
		UserID   string  `json:"user_id"`
		Value    float64 `json:"value"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.UpdateChallengeProgress(r.Context(), &v1.UpdateChallengeProgressRequest{
		ChallengeId: challengeId,
		UserId:      body.UserID,
		Value:       body.Value,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Family (new) ──

func (h *RestHandler) handleCreateFamilyGroup(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	var body struct {
		OwnerUserID string `json:"owner_user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.CreateFamilyGroup(r.Context(), &v1.CreateFamilyGroupRequest{
		OwnerUserId: body.OwnerUserID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetFamilyGroup(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	resp, err := h.family.GetFamilyGroup(r.Context(), &v1.GetFamilyGroupRequest{
		GroupId: groupId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleInviteMember(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	var body struct {
		InviterUserID string `json:"inviter_user_id"`
		InviteeEmail  string `json:"invitee_email"`
		Role          int32  `json:"role"`
		Relationship  string `json:"relationship"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.InviteMember(r.Context(), &v1.InviteMemberRequest{
		GroupId:       groupId,
		InviterUserId: body.InviterUserID,
		InviteeEmail:  body.InviteeEmail,
		Role:          v1.FamilyRole(body.Role),
		Relationship:  body.Relationship,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleRespondToInvitation(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	invitationId := r.PathValue("invitationId")
	var body struct {
		UserID string `json:"user_id"`
		Accept bool   `json:"accept"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.RespondToInvitation(r.Context(), &v1.RespondToInvitationRequest{
		InvitationId: invitationId,
		UserId:       body.UserID,
		Accept:       body.Accept,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRemoveFamilyMember(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	userId := r.PathValue("userId")
	var body struct {
		RequesterUserID string `json:"requester_user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.RemoveMember(r.Context(), &v1.RemoveMemberRequest{
		GroupId:         groupId,
		RequesterUserId: body.RequesterUserID,
		TargetUserId:    userId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	userId := r.PathValue("userId")
	var body struct {
		RequesterUserID string `json:"requester_user_id"`
		NewRole         int32  `json:"new_role"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.UpdateMemberRole(r.Context(), &v1.UpdateMemberRoleRequest{
		GroupId:         groupId,
		RequesterUserId: body.RequesterUserID,
		TargetUserId:    userId,
		NewRole:         v1.FamilyRole(body.NewRole),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListFamilyMembers(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	resp, err := h.family.ListFamilyMembers(r.Context(), &v1.ListFamilyMembersRequest{
		GroupId: groupId,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSetSharingPreferences(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	var body struct {
		UserID               string   `json:"user_id"`
		ShareMeasurements    bool     `json:"share_measurements"`
		ShareHealthRecords   bool     `json:"share_health_records"`
		SharePrescriptions   bool     `json:"share_prescriptions"`
		ShareCoaching        bool     `json:"share_coaching"`
		SharedWithUserIDs    []string `json:"shared_with_user_ids"`
		ShareHealthScore     bool     `json:"share_health_score"`
		ShareGoals           bool     `json:"share_goals"`
		ShareAlerts          bool     `json:"share_alerts"`
		MeasurementDaysLimit int32    `json:"measurement_days_limit"`
		AllowedBiomarkers    []string `json:"allowed_biomarkers"`
		RequireApproval      bool     `json:"require_approval"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.SetSharingPreferences(r.Context(), &v1.SetSharingPreferencesRequest{
		GroupId:              groupId,
		UserId:               body.UserID,
		ShareMeasurements:    body.ShareMeasurements,
		ShareHealthRecords:   body.ShareHealthRecords,
		SharePrescriptions:   body.SharePrescriptions,
		ShareCoaching:        body.ShareCoaching,
		SharedWithUserIds:    body.SharedWithUserIDs,
		ShareHealthScore:     body.ShareHealthScore,
		ShareGoals:           body.ShareGoals,
		ShareAlerts:          body.ShareAlerts,
		MeasurementDaysLimit: body.MeasurementDaysLimit,
		AllowedBiomarkers:    body.AllowedBiomarkers,
		RequireApproval:      body.RequireApproval,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleValidateSharingAccess(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	var body struct {
		GroupID         string `json:"group_id"`
		RequesterUserID string `json:"requester_user_id"`
		TargetUserID    string `json:"target_user_id"`
		Biomarker       string `json:"biomarker"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.family.ValidateSharingAccess(r.Context(), &v1.ValidateSharingAccessRequest{
		GroupId:         body.GroupID,
		RequesterUserId: body.RequesterUserID,
		TargetUserId:    body.TargetUserID,
		Biomarker:       body.Biomarker,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Telemedicine (new) ──

func (h *RestHandler) handleGetConsultation(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	consultationId := r.PathValue("consultationId")
	resp, err := h.telemedicine.GetConsultation(r.Context(), &v1.GetConsultationRequest{
		ConsultationId: consultationId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListConsultations(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	resp, err := h.telemedicine.ListConsultations(r.Context(), &v1.ListConsultationsRequest{
		UserId:       r.URL.Query().Get("user_id"),
		StatusFilter: v1.ConsultationStatus(queryInt(r, "status", 0)),
		Limit:        queryInt(r, "limit", 20),
		Offset:       queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleStartVideoSession(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	consultationId := r.PathValue("consultationId")
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.telemedicine.StartVideoSession(r.Context(), &v1.StartVideoSessionRequest{
		ConsultationId: consultationId,
		UserId:         body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleEndVideoSession(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	consultationId := r.PathValue("consultationId")
	var body struct {
		SessionID   string `json:"session_id"`
		DoctorNotes string `json:"doctor_notes"`
		Diagnosis   string `json:"diagnosis"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.telemedicine.EndVideoSession(r.Context(), &v1.EndVideoSessionRequest{
		SessionId:      body.SessionID,
		ConsultationId: consultationId,
		DoctorNotes:    body.DoctorNotes,
		Diagnosis:      body.Diagnosis,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRateConsultation(w http.ResponseWriter, r *http.Request) {
	if h.telemedicine == nil {
		writeError(w, http.StatusServiceUnavailable, "telemedicine service unavailable")
		return
	}
	consultationId := r.PathValue("consultationId")
	var body struct {
		Rating float64 `json:"rating"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.telemedicine.RateConsultation(r.Context(), &v1.RateConsultationRequest{
		ConsultationId: consultationId,
		Rating:         body.Rating,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Video (new) ──

func (h *RestHandler) handleCreateVideoRoom(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	var body struct {
		HostUserID      string `json:"host_user_id"`
		RoomType        int32  `json:"room_type"`
		Title           string `json:"title"`
		MaxParticipants int32  `json:"max_participants"`
		ScheduledAt     string `json:"scheduled_at"`
		ReservationID   string `json:"reservation_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.video.CreateRoom(r.Context(), &v1.CreateRoomRequest{
		HostUserId:      body.HostUserID,
		RoomType:        v1.RoomType(body.RoomType),
		Title:           body.Title,
		MaxParticipants: body.MaxParticipants,
		ScheduledAt:     body.ScheduledAt,
		ReservationId:   body.ReservationID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetVideoRoom(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	resp, err := h.video.GetRoom(r.Context(), &v1.GetRoomRequest{
		RoomId: roomId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleEndVideoRoom(w http.ResponseWriter, r *http.Request) {
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
	resp, err := h.video.EndRoom(r.Context(), &v1.EndRoomRequest{
		RoomId: roomId,
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSendSignal(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	var body struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		SignalType int32  `json:"signal_type"`
		Payload    string `json:"payload"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.video.SendSignal(r.Context(), &v1.SendSignalRequest{
		RoomId:     roomId,
		FromUserId: body.FromUserID,
		ToUserId:   body.ToUserID,
		SignalType: v1.SignalType(body.SignalType),
		Payload:    body.Payload,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListParticipants(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	resp, err := h.video.ListParticipants(r.Context(), &v1.ListParticipantsRequest{
		RoomId: roomId,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetRoomStats(w http.ResponseWriter, r *http.Request) {
	if h.video == nil {
		writeError(w, http.StatusServiceUnavailable, "video service unavailable")
		return
	}
	roomId := r.PathValue("roomId")
	resp, err := h.video.GetRoomStats(r.Context(), &v1.GetRoomStatsRequest{
		RoomId: roomId,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Reservation (신규 엔드포인트) ──

func (h *RestHandler) handleGetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	reservationId := r.PathValue("reservationId")
	resp, err := h.reservation.GetAvailableSlots(r.Context(), &v1.GetAvailableSlotsRequest{
		FacilityId: reservationId,
		DoctorId:   r.URL.Query().Get("doctor_id"),
		Specialty:  v1.DoctorSpecialty(queryInt(r, "specialty", 0)),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCancelReservation(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	reservationId := r.PathValue("reservationId")
	resp, err := h.reservation.CancelReservation(r.Context(), &v1.CancelReservationRequest{
		ReservationId: reservationId,
		Reason:        r.URL.Query().Get("reason"),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListDoctorsByFacility(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	facilityId := r.PathValue("facilityId")
	resp, err := h.reservation.ListDoctorsByFacility(r.Context(), &v1.ListDoctorsByFacilityRequest{
		FacilityId: facilityId,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetDoctorAvailability(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	doctorId := r.PathValue("doctorId")
	resp, err := h.reservation.GetDoctorAvailability(r.Context(), &v1.GetDoctorAvailabilityRequest{
		DoctorId: doctorId,
		Date:     r.URL.Query().Get("date"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSelectDoctor(w http.ResponseWriter, r *http.Request) {
	if h.reservation == nil {
		writeError(w, http.StatusServiceUnavailable, "reservation service unavailable")
		return
	}
	reservationId := r.PathValue("reservationId")
	var body struct {
		DoctorID string `json:"doctor_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.reservation.SelectDoctor(r.Context(), &v1.SelectDoctorRequest{
		FacilityId: reservationId,
		DoctorId:   body.DoctorID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Family – Guardian Dashboard ──

func (h *RestHandler) handleGetGuardianDashboard(w http.ResponseWriter, r *http.Request) {
	if h.family == nil {
		writeError(w, http.StatusServiceUnavailable, "family service unavailable")
		return
	}
	groupId := r.PathValue("groupId")
	userId := r.URL.Query().Get("user_id")
	resp, err := h.family.GetSharedHealthData(r.Context(), &v1.GetSharedHealthDataRequest{
		RequesterUserId: userId,
		GroupId:         groupId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
