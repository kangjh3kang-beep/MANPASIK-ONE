package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// ============================================================================
// Health Record Service Handlers
// ============================================================================

// POST /api/v1/health-records
func (r *Router) handleCreateRecord(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID        string            `json:"user_id"`
		RecordType    int32             `json:"record_type"`
		Title         string            `json:"title"`
		Description   string            `json:"description"`
		Provider      string            `json:"provider"`
		Metadata      map[string]string `json:"metadata"`
		MeasurementID string            `json:"measurement_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.healthRecordAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "건강기록 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewHealthRecordServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateRecord(ctx, &v1.CreateHealthRecordRequest{
		UserId:        body.UserID,
		RecordType:    v1.HealthRecordType(body.RecordType),
		Title:         body.Title,
		Description:   body.Description,
		Provider:      body.Provider,
		Metadata:      body.Metadata,
		MeasurementId: body.MeasurementID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, healthRecordToMap(resp))
}

// GET /api/v1/health-records?user_id=...&type_filter=...&limit=...&offset=...
func (r *Router) handleListRecords(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	userID := q.Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	typeFilter, _ := strconv.Atoi(q.Get("type_filter"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.healthRecordAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "건강기록 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewHealthRecordServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListRecords(ctx, &v1.ListHealthRecordsRequest{
		UserId:     userID,
		TypeFilter: v1.HealthRecordType(typeFilter),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	records := make([]map[string]interface{}, 0, len(resp.GetRecords()))
	for _, rec := range resp.GetRecords() {
		records = append(records, healthRecordToMap(rec))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"records":     records,
		"total_count": resp.GetTotalCount(),
	})
}

// GET /api/v1/health-records/{recordId}
func (r *Router) handleGetRecord(w http.ResponseWriter, req *http.Request) {
	recordID := req.PathValue("recordId")
	if recordID == "" {
		writeError(w, http.StatusBadRequest, "record_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.healthRecordAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "건강기록 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewHealthRecordServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetRecord(ctx, &v1.GetHealthRecordRequest{
		RecordId: recordID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, healthRecordToMap(resp))
}

// POST /api/v1/health-records/export/fhir
func (r *Router) handleExportToFHIR(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID     string   `json:"user_id"`
		RecordIDs  []string `json:"record_ids"`
		TargetType int32    `json:"target_type"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.healthRecordAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "건강기록 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewHealthRecordServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ExportToFHIR(ctx, &v1.ExportToFHIRRequest{
		UserId:     body.UserID,
		RecordIds:  body.RecordIDs,
		TargetType: v1.FHIRResourceType(body.TargetType),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"fhir_bundle_json": resp.GetFhirBundleJson(),
		"resource_type":    resp.GetResourceType().String(),
		"resource_count":   resp.GetResourceCount(),
	})
}

func healthRecordToMap(rec *v1.HealthRecord) map[string]interface{} {
	m := map[string]interface{}{
		"record_id":      rec.GetRecordId(),
		"user_id":        rec.GetUserId(),
		"record_type":    rec.GetRecordType().String(),
		"title":          rec.GetTitle(),
		"description":    rec.GetDescription(),
		"provider":       rec.GetProvider(),
		"metadata":       rec.GetMetadata(),
		"measurement_id": rec.GetMeasurementId(),
	}
	if rec.GetCreatedAt() != nil {
		m["created_at"] = rec.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if rec.GetUpdatedAt() != nil {
		m["updated_at"] = rec.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Notification Service Handlers
// ============================================================================

// GET /api/v1/notifications?user_id=...&unread_only=...&limit=...&offset=...
func (r *Router) handleListNotifications(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	userID := q.Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	unreadOnly, _ := strconv.ParseBool(q.Get("unread_only"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.notificationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "알림 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewNotificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListNotifications(ctx, &v1.ListNotificationsRequest{
		UserId:     userID,
		UnreadOnly: unreadOnly,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	notifications := make([]map[string]interface{}, 0, len(resp.GetNotifications()))
	for _, n := range resp.GetNotifications() {
		notifications = append(notifications, notificationToMap(n))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"total_count":   resp.GetTotalCount(),
	})
}

// POST /api/v1/notifications/{notificationId}/read
func (r *Router) handleMarkAsRead(w http.ResponseWriter, req *http.Request) {
	notificationID := req.PathValue("notificationId")
	if notificationID == "" {
		writeError(w, http.StatusBadRequest, "notification_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.notificationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "알림 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewNotificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.MarkAsRead(ctx, &v1.MarkAsReadRequest{
		NotificationId: notificationID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": resp.GetSuccess(),
	})
}

// GET /api/v1/notifications/unread-count?user_id=...
func (r *Router) handleGetUnreadCount(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.notificationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "알림 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewNotificationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUnreadCount(ctx, &v1.GetUnreadCountRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count": resp.GetCount(),
	})
}

func notificationToMap(n *v1.Notification) map[string]interface{} {
	m := map[string]interface{}{
		"notification_id": n.GetNotificationId(),
		"user_id":         n.GetUserId(),
		"type":            n.GetType().String(),
		"title":           n.GetTitle(),
		"body":            n.GetBody(),
		"priority":        n.GetPriority().String(),
		"channel":         n.GetChannel().String(),
		"is_read":         n.GetIsRead(),
		"data":            n.GetData(),
		"action_url":      n.GetActionUrl(),
	}
	if n.GetCreatedAt() != nil {
		m["created_at"] = n.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if n.GetReadAt() != nil {
		m["read_at"] = n.GetReadAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Community Service Handlers
// ============================================================================

// GET /api/v1/posts?category=...&author_id=...&query=...&limit=...&offset=...
func (r *Router) handleListPosts(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	category, _ := strconv.Atoi(q.Get("category"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.communityAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "커뮤니티 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCommunityServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListPosts(ctx, &v1.ListPostsRequest{
		Category: v1.PostCategory(category),
		AuthorId: q.Get("author_id"),
		Query:    q.Get("query"),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	posts := make([]map[string]interface{}, 0, len(resp.GetPosts()))
	for _, p := range resp.GetPosts() {
		posts = append(posts, postToMap(p))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"posts":       posts,
		"total_count": resp.GetTotalCount(),
	})
}

// POST /api/v1/posts
func (r *Router) handleCreatePost(w http.ResponseWriter, req *http.Request) {
	var body struct {
		AuthorID string   `json:"author_id"`
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category int32    `json:"category"`
		Tags     []string `json:"tags"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.communityAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "커뮤니티 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCommunityServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreatePost(ctx, &v1.CreatePostRequest{
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

	writeJSON(w, http.StatusCreated, postToMap(resp))
}

// GET /api/v1/posts/{postId}
func (r *Router) handleGetPost(w http.ResponseWriter, req *http.Request) {
	postID := req.PathValue("postId")
	if postID == "" {
		writeError(w, http.StatusBadRequest, "post_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.communityAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "커뮤니티 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCommunityServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPost(ctx, &v1.GetPostRequest{
		PostId: postID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, postToMap(resp))
}

// POST /api/v1/posts/{postId}/like
func (r *Router) handleLikePost(w http.ResponseWriter, req *http.Request) {
	postID := req.PathValue("postId")
	if postID == "" {
		writeError(w, http.StatusBadRequest, "post_id가 필요합니다")
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.communityAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "커뮤니티 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCommunityServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.LikePost(ctx, &v1.LikePostRequest{
		PostId: postID,
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":    resp.GetSuccess(),
		"like_count": resp.GetLikeCount(),
		"is_liked":   resp.GetIsLiked(),
	})
}

func postToMap(p *v1.Post) map[string]interface{} {
	m := map[string]interface{}{
		"post_id":       p.GetPostId(),
		"author_id":     p.GetAuthorId(),
		"author_name":   p.GetAuthorName(),
		"title":         p.GetTitle(),
		"content":       p.GetContent(),
		"category":      p.GetCategory().String(),
		"tags":          p.GetTags(),
		"like_count":    p.GetLikeCount(),
		"comment_count": p.GetCommentCount(),
		"view_count":    p.GetViewCount(),
	}
	if p.GetCreatedAt() != nil {
		m["created_at"] = p.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if p.GetUpdatedAt() != nil {
		m["updated_at"] = p.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Admin Service Handlers
// ============================================================================

// GET /api/v1/admin/stats
func (r *Router) handleGetSystemStats(w http.ResponseWriter, req *http.Request) {
	conn, err := dialGRPC(r.adminAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "관리자 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetSystemStats(ctx, &v1.GetSystemStatsRequest{})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := map[string]interface{}{
		"total_users":          resp.GetTotalUsers(),
		"active_users":         resp.GetActiveUsers(),
		"total_devices":        resp.GetTotalDevices(),
		"total_measurements":   resp.GetTotalMeasurements(),
		"users_by_tier":        resp.GetUsersByTier(),
		"measurements_by_type": resp.GetMeasurementsByType(),
	}
	if resp.GetGeneratedAt() != nil {
		result["generated_at"] = resp.GetGeneratedAt().AsTime().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/v1/admin/users?query=...&tier_filter=...&active_only=...&limit=...&offset=...
func (r *Router) handleAdminListUsers(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	tierFilter, _ := strconv.Atoi(q.Get("tier_filter"))
	activeOnly, _ := strconv.ParseBool(q.Get("active_only"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.adminAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "관리자 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ListUsers(ctx, &v1.AdminListUsersRequest{
		Query:      q.Get("query"),
		TierFilter: v1.SubscriptionTier(tierFilter),
		ActiveOnly: activeOnly,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	users := make([]map[string]interface{}, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		um := map[string]interface{}{
			"user_id":           u.GetUserId(),
			"email":             u.GetEmail(),
			"display_name":      u.GetDisplayName(),
			"tier":              u.GetTier().String(),
			"is_active":         u.GetIsActive(),
			"device_count":      u.GetDeviceCount(),
			"measurement_count": u.GetMeasurementCount(),
		}
		if u.GetCreatedAt() != nil {
			um["created_at"] = u.GetCreatedAt().AsTime().Format(time.RFC3339)
		}
		if u.GetLastActiveAt() != nil {
			um["last_active_at"] = u.GetLastActiveAt().AsTime().Format(time.RFC3339)
		}
		users = append(users, um)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":       users,
		"total_count": resp.GetTotalCount(),
	})
}

// GET /api/v1/admin/audit-log?admin_id=...&limit=...&offset=...
func (r *Router) handleGetAuditLog(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.adminAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "관리자 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetAuditLog(ctx, &v1.GetAuditLogRequest{
		AdminId: q.Get("admin_id"),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	entries := make([]map[string]interface{}, 0, len(resp.GetEntries()))
	for _, e := range resp.GetEntries() {
		em := map[string]interface{}{
			"entry_id":      e.GetEntryId(),
			"admin_id":      e.GetAdminId(),
			"admin_email":   e.GetAdminEmail(),
			"action":        e.GetAction().String(),
			"resource_type": e.GetResourceType(),
			"resource_id":   e.GetResourceId(),
			"details":       e.GetDetails(),
			"ip_address":    e.GetIpAddress(),
		}
		if e.GetTimestamp() != nil {
			em["timestamp"] = e.GetTimestamp().AsTime().Format(time.RFC3339)
		}
		entries = append(entries, em)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"entries":     entries,
		"total_count": resp.GetTotalCount(),
	})
}
