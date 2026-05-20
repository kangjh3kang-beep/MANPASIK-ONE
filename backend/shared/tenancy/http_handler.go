package tenancy

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// HTTPHandler 는 tenancy 의 5개 REST 엔드포인트를 net/http ServeMux 에 등록.
//
// 엔드포인트:
//   POST   /tenancy/invitations              초대 발급 (body: {tenant_id, role, invitee_hint?})
//   POST   /tenancy/invitations/accept       초대 수락 (body: {token})
//   DELETE /tenancy/invitations/{token}      초대 취소
//   GET    /tenancy/me/memberships           내 활성 멤버십 목록
//   DELETE /tenancy/tenants/{tid}/members/{uid}  멤버 제거 (admin)
//
// 인증: 모든 엔드포인트는 ctx 에 사용자 ID (UserFromContext) 가 있어야 함.
// 게이트웨이의 TenantPropagation 미들웨어와 결합하여 사용.
type HTTPHandler struct {
	invSvc      *InvitationService
	memStore    MembershipStore
	policy      *PolicyEngine
	pathPrefix  string
}

// NewHTTPHandler 생성. policy=nil 가능 (admin 권한 검사 생략).
func NewHTTPHandler(invSvc *InvitationService, memStore MembershipStore, policy *PolicyEngine) *HTTPHandler {
	return &HTTPHandler{invSvc: invSvc, memStore: memStore, policy: policy}
}

// SetPathPrefix 는 모든 라우트에 prefix 부착 (예: "/api/v1").
func (h *HTTPHandler) SetPathPrefix(prefix string) {
	h.pathPrefix = strings.TrimRight(prefix, "/")
}

// RegisterRoutes 는 ServeMux 에 핸들러 등록.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	p := h.pathPrefix
	mux.HandleFunc(p+"/tenancy/invitations", h.handleInvitations)
	mux.HandleFunc(p+"/tenancy/invitations/accept", h.handleAccept)
	mux.HandleFunc(p+"/tenancy/invitations/", h.handleInvitationsByToken)
	mux.HandleFunc(p+"/tenancy/me/memberships", h.handleMyMemberships)
	mux.HandleFunc(p+"/tenancy/tenants/", h.handleTenantMembers)
}

// inviteRequestBody 는 POST /invitations 본문.
type inviteRequestBody struct {
	TenantID    string `json:"tenant_id"`
	Role        string `json:"role"`
	InviteeHint string `json:"invitee_hint,omitempty"`
	TTLHours    int    `json:"ttl_hours,omitempty"`
}

// acceptRequestBody 는 POST /invitations/accept 본문.
type acceptRequestBody struct {
	Token string `json:"token"`
}

// invitationResponse 는 초대 정보 JSON.
type invitationResponse struct {
	Token       string `json:"token"`
	TenantID    string `json:"tenant_id"`
	Role        string `json:"role"`
	InviteeHint string `json:"invitee_hint,omitempty"`
	Status      string `json:"status"`
	IssuedAt    string `json:"issued_at"`
	ExpiresAt   string `json:"expires_at"`
}

func (h *HTTPHandler) handleInvitations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "POST 만 지원")
		return
	}

	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "사용자 컨텍스트 없음")
		return
	}

	var body inviteRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "JSON 본문 파싱 실패")
		return
	}
	if body.TenantID == "" || body.Role == "" {
		writeErr(w, http.StatusBadRequest, "tenant_id, role 필수")
		return
	}

	ttl := time.Duration(body.TTLHours) * time.Hour
	inv, err := h.invSvc.Invite(InviteRequest{
		InviterID:   userID,
		TenantID:    TenantID(body.TenantID),
		InviteeHint: body.InviteeHint,
		Role:        TenantRole(body.Role),
		TTL:         ttl,
	})
	if err != nil {
		writeErr(w, statusFromInviteErr(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, invitationToJSON(inv))
}

func (h *HTTPHandler) handleAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "POST 만 지원")
		return
	}
	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "사용자 컨텍스트 없음")
		return
	}

	var body acceptRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "JSON 본문 파싱 실패")
		return
	}
	if body.Token == "" {
		writeErr(w, http.StatusBadRequest, "token 필수")
		return
	}

	m, err := h.invSvc.Accept(body.Token, userID)
	if err != nil {
		writeErr(w, statusFromAcceptErr(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, membershipToJSON(m))
}

// handleInvitationsByToken 처리 — DELETE /tenancy/invitations/{token}
func (h *HTTPHandler) handleInvitationsByToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErr(w, http.StatusMethodNotAllowed, "DELETE 만 지원")
		return
	}
	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "사용자 컨텍스트 없음")
		return
	}

	token := strings.TrimPrefix(r.URL.Path, h.pathPrefix+"/tenancy/invitations/")
	if token == "" || strings.Contains(token, "/") {
		writeErr(w, http.StatusBadRequest, "유효하지 않은 token 경로")
		return
	}

	if err := h.invSvc.Revoke(token, userID); err != nil {
		writeErr(w, statusFromRevokeErr(err), err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *HTTPHandler) handleMyMemberships(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "GET 만 지원")
		return
	}
	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "사용자 컨텍스트 없음")
		return
	}
	list := h.memStore.ListUserTenants(userID)
	out := make([]map[string]interface{}, 0, len(list))
	for _, m := range list {
		if m.Active {
			out = append(out, membershipToJSON(m))
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"memberships": out})
}

// handleTenantMembers 는 다음 라우트를 분기:
//   - GET    /tenancy/tenants/{tid}/members           → 멤버 목록 (admin)
//   - DELETE /tenancy/tenants/{tid}/members/{uid}     → 멤버 제거 (자신 또는 admin)
//   - PATCH  /tenancy/tenants/{tid}/members/{uid}/role → 역할 변경 (admin)
func (h *HTTPHandler) handleTenantMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, "사용자 컨텍스트 없음")
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, h.pathPrefix+"/tenancy/tenants/")
	parts := strings.Split(rest, "/")

	// GET /tenancy/tenants/{tid}/members → 멤버 목록
	if r.Method == http.MethodGet && len(parts) == 2 && parts[1] == "members" {
		h.serveTenantMemberList(w, r, userID, TenantID(parts[0]))
		return
	}

	// /tenancy/tenants/{tid}/members/{uid}[/role]
	if len(parts) < 3 || parts[1] != "members" {
		writeErr(w, http.StatusBadRequest, "경로 형식: /tenants/{tid}/members/{uid}")
		return
	}
	tid := TenantID(parts[0])
	targetUserID := parts[2]

	// PATCH /role 으로 역할 변경
	if r.Method == http.MethodPatch && len(parts) == 4 && parts[3] == "role" {
		h.handleRoleUpdate(w, r, userID, tid, targetUserID)
		return
	}

	if r.Method != http.MethodDelete {
		writeErr(w, http.StatusMethodNotAllowed, "지원하지 않는 메서드")
		return
	}

	// 권한 검사: 본인이거나 admin 권한
	if targetUserID != userID && h.policy != nil {
		res := &Resource{TenantID: tid}
		if d := h.policy.Evaluate(userID, tid, res, ActionAdmin); !d.Allowed {
			writeErr(w, http.StatusForbidden, "멤버 제거 권한 없음")
			return
		}
	}

	if err := h.memStore.Remove(targetUserID, tid); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

// serveTenantMemberList 는 GET /tenants/{tid}/members 처리.
func (h *HTTPHandler) serveTenantMemberList(w http.ResponseWriter, r *http.Request, userID string, tid TenantID) {
	// admin 권한 체크
	if h.policy != nil {
		res := &Resource{TenantID: tid}
		if d := h.policy.Evaluate(userID, tid, res, ActionAdmin); !d.Allowed {
			writeErr(w, http.StatusForbidden, "조회 권한 없음")
			return
		}
	}
	list := h.memStore.ListTenantMembers(tid)
	out := make([]map[string]interface{}, 0, len(list))
	for _, m := range list {
		out = append(out, membershipToJSON(m))
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"members": out})
}

// handleRoleUpdate 는 PATCH /tenants/{tid}/members/{uid}/role 처리.
type roleUpdateBody struct {
	Role string `json:"role"`
}

func (h *HTTPHandler) handleRoleUpdate(w http.ResponseWriter, r *http.Request, userID string, tid TenantID, targetUserID string) {
	// admin 권한 필수 (본인이라도 자신을 admin → owner 로 승격 불가; 권한 검사 일관)
	if h.policy != nil {
		res := &Resource{TenantID: tid}
		if d := h.policy.Evaluate(userID, tid, res, ActionAdmin); !d.Allowed {
			writeErr(w, http.StatusForbidden, "역할 변경 권한 없음")
			return
		}
	}

	var body roleUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "JSON 본문 파싱 실패")
		return
	}
	newRole := TenantRole(body.Role)
	if !newRole.IsKnown() {
		writeErr(w, http.StatusBadRequest, "알 수 없는 역할: "+body.Role)
		return
	}
	if err := h.memStore.UpdateRole(targetUserID, tid, newRole); err != nil {
		if err == ErrNoMembership {
			writeErr(w, http.StatusNotFound, err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	got, err := h.memStore.Get(targetUserID, tid)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"updated": true})
		return
	}
	writeJSON(w, http.StatusOK, membershipToJSON(got))
}

// invitationToJSON 변환.
func invitationToJSON(inv *Invitation) map[string]interface{} {
	return map[string]interface{}{
		"token":        inv.Token,
		"tenant_id":    string(inv.TenantID),
		"role":         string(inv.Role),
		"invitee_hint": inv.InviteeHint,
		"status":       string(inv.Status),
		"issued_at":    inv.IssuedAt.Format(time.RFC3339),
		"expires_at":   inv.ExpiresAt.Format(time.RFC3339),
	}
}

// membershipToJSON 변환.
func membershipToJSON(m *Membership) map[string]interface{} {
	return map[string]interface{}{
		"user_id":   m.UserID,
		"tenant_id": string(m.TenantID),
		"role":      string(m.Role),
		"active":    m.Active,
		"joined_at": m.JoinedAt,
	}
}

// statusFromInviteErr / AcceptErr / RevokeErr — 에러 → HTTP 상태 매핑.
func statusFromInviteErr(err error) int {
	if err == nil {
		return http.StatusOK
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "권한 없음"):
		return http.StatusForbidden
	case strings.Contains(msg, "필수") || strings.Contains(msg, "알 수 없는 역할"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func statusFromAcceptErr(err error) int {
	switch {
	case errors.Is(err, ErrInvitationNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvitationExpired):
		return http.StatusGone
	case errors.Is(err, ErrInvitationConsumed):
		return http.StatusConflict
	default:
		if strings.Contains(err.Error(), "이미 활성 멤버") {
			return http.StatusConflict
		}
		return http.StatusBadRequest
	}
}

func statusFromRevokeErr(err error) int {
	switch {
	case errors.Is(err, ErrInvitationNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvitationConsumed):
		return http.StatusConflict
	default:
		if strings.Contains(err.Error(), "권한") {
			return http.StatusForbidden
		}
		return http.StatusInternalServerError
	}
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
