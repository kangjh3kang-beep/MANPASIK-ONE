package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-COM-001: 커뮤니티 게시글 CRUD ───

func TestCommunityPostCRUD(t *testing.T) {
	email := uniqueEmail("comm-crud")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Community User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 게시글 작성
	postBody := map[string]interface{}{
		"title":    "혈당 관리 팁 공유합니다",
		"content":  "아침 공복 혈당 측정 후 식이 조절하니 수치가 안정됐어요.",
		"category": "health_tips",
		"tags":     []string{"혈당", "식이조절"},
	}
	createResp, createResult := apiRequest(t, "POST", "/api/v1/community/posts", postBody, token)
	defer createResp.Body.Close()
	t.Logf("1. 게시글 작성: status=%d, result=%v", createResp.StatusCode, createResult)

	// 2. 게시글 목록 조회
	listResp, listResult := apiRequest(t, "GET", "/api/v1/community/posts?limit=10", nil, token)
	defer listResp.Body.Close()
	t.Logf("2. 게시글 목록: status=%d, result=%v", listResp.StatusCode, listResult)

	// 3. 댓글 작성
	commentBody := map[string]string{
		"content": "좋은 정보 감사합니다!",
	}
	commentResp, commentResult := apiRequest(t, "POST", "/api/v1/community/posts/1/comments", commentBody, token)
	defer commentResp.Body.Close()
	t.Logf("3. 댓글 작성: status=%d, result=%v", commentResp.StatusCode, commentResult)
}

// ─── E2E-COM-002: 커뮤니티 좋아요/북마크 ───

func TestCommunityEngagement(t *testing.T) {
	email := uniqueEmail("comm-engage")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Engage User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 좋아요
	likeResp, _ := apiRequest(t, "POST", "/api/v1/community/posts/1/like", nil, token)
	defer likeResp.Body.Close()
	t.Logf("좋아요: status=%d", likeResp.StatusCode)

	// 북마크
	bookmarkResp, _ := apiRequest(t, "POST", "/api/v1/community/posts/1/bookmark", nil, token)
	defer bookmarkResp.Body.Close()
	t.Logf("북마크: status=%d", bookmarkResp.StatusCode)
}

// ─── E2E-FAM-001: 가족 그룹 생성 및 멤버 초대 ───

func TestFamilyGroupManagement(t *testing.T) {
	parentEmail := uniqueEmail("fam-parent")
	childEmail := uniqueEmail("fam-child")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "부모 사용자")
	_ = registerAndLogin(t, childEmail, "E2eTest1!@#", "자녀 사용자")
	if parentToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 가족 그룹 생성
	familyBody := map[string]string{
		"name": "우리 가족",
	}
	createResp, createResult := apiRequest(t, "POST", "/api/v1/family/groups", familyBody, parentToken)
	defer createResp.Body.Close()
	t.Logf("1. 가족 그룹 생성: status=%d, result=%v", createResp.StatusCode, createResult)

	// 2. 멤버 초대
	inviteBody := map[string]string{
		"email": childEmail,
		"role":  "child",
	}
	inviteResp, inviteResult := apiRequest(t, "POST", "/api/v1/family/groups/1/invite", inviteBody, parentToken)
	defer inviteResp.Body.Close()
	t.Logf("2. 멤버 초대: status=%d, result=%v", inviteResp.StatusCode, inviteResult)

	// 3. 가족 그룹 조회
	groupResp, groupResult := apiRequest(t, "GET", "/api/v1/family/groups", nil, parentToken)
	defer groupResp.Body.Close()
	t.Logf("3. 가족 그룹: status=%d, result=%v", groupResp.StatusCode, groupResult)
}

// ─── E2E-FAM-002: 가족 건강 데이터 공유 ───

func TestFamilyHealthDataSharing(t *testing.T) {
	parentEmail := uniqueEmail("fam-share-p")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "공유 부모")
	if parentToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 가족 건강 대시보드 조회
	dashResp, dashResult := apiRequest(t, "GET", "/api/v1/family/dashboard", nil, parentToken)
	defer dashResp.Body.Close()
	t.Logf("가족 대시보드: status=%d, result=%v", dashResp.StatusCode, dashResult)
}

// ─── E2E-FAM-003: 가족 멤버 권한 관리 ───

func TestFamilyMemberPermissions(t *testing.T) {
	parentEmail := uniqueEmail("fam-perm-p")
	childEmail := uniqueEmail("fam-perm-c")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "권한 부모")
	childToken := registerAndLogin(t, childEmail, "E2eTest1!@#", "권한 자녀")
	if parentToken == "" || childToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 자녀가 가족 그룹 삭제 시도 (실패 예상)
	delResp, _ := apiRequest(t, "DELETE", "/api/v1/family/groups/1", nil, childToken)
	defer delResp.Body.Close()
	if delResp.StatusCode == http.StatusOK {
		t.Error("가족 권한 위반: 자녀가 가족 그룹 삭제 가능")
	}
	t.Logf("자녀 그룹 삭제 차단: status=%d", delResp.StatusCode)
}
