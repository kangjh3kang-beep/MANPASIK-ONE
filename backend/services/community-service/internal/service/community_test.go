package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/community-service/internal/repository/memory"
	"github.com/manpasik/backend/services/community-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.CommunityService {
	logger := zap.NewNop()
	return service.NewCommunityService(
		logger,
		memory.NewPostRepository(),
		memory.NewCommentRepository(),
		memory.NewChallengeRepository(),
	)
}

// ============================================================================
// TestCreatePost_Success — 게시글 작성 성공
// ============================================================================

func TestCreatePost_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(ctx, "user-1", "홍길동", service.CatGeneral, "건강 관리 팁", "매일 운동하세요!", []string{"건강", "운동"})
	if err != nil {
		t.Fatalf("CreatePost 실패: %v", err)
	}
	if post.ID == "" {
		t.Fatal("Post ID가 비어 있음")
	}
	if post.AuthorUserID != "user-1" {
		t.Errorf("AuthorUserID: got %s, want user-1", post.AuthorUserID)
	}
	if post.AuthorDisplayName != "홍길동" {
		t.Errorf("AuthorDisplayName: got %s, want 홍길동", post.AuthorDisplayName)
	}
	if post.Category != service.CatGeneral {
		t.Errorf("Category: got %d, want %d", post.Category, service.CatGeneral)
	}
	if post.Title != "건강 관리 팁" {
		t.Errorf("Title: got %s, want 건강 관리 팁", post.Title)
	}
	if post.Content != "매일 운동하세요!" {
		t.Errorf("Content: got %s, want 매일 운동하세요!", post.Content)
	}
	if len(post.Tags) != 2 {
		t.Errorf("Tags 수: got %d, want 2", len(post.Tags))
	}
	if post.LikeCount != 0 {
		t.Errorf("LikeCount: got %d, want 0", post.LikeCount)
	}
	if post.CommentCount != 0 {
		t.Errorf("CommentCount: got %d, want 0", post.CommentCount)
	}
	if post.ViewCount != 0 {
		t.Errorf("ViewCount: got %d, want 0", post.ViewCount)
	}
	if post.CreatedAt.IsZero() {
		t.Error("CreatedAt가 비어 있음")
	}
}

// ============================================================================
// TestCreatePost_MissingAuthor — 작성자 누락
// ============================================================================

func TestCreatePost_MissingAuthor(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreatePost(ctx, "", "이름", service.CatGeneral, "제목", "내용", nil)
	if err == nil {
		t.Fatal("작성자 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestCreatePost_MissingTitle — 제목 누락
// ============================================================================

func TestCreatePost_MissingTitle(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreatePost(ctx, "user-1", "이름", service.CatGeneral, "", "내용", nil)
	if err == nil {
		t.Fatal("제목 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetPost_Success — 게시글 조회 성공
// ============================================================================

func TestGetPost_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreatePost(ctx, "user-1", "홍길동", service.CatHealthTip, "혈당 관리", "식후 산책하세요", nil)

	got, err := svc.GetPost(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPost 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("Post ID: got %s, want %s", got.ID, created.ID)
	}
	if got.Title != "혈당 관리" {
		t.Errorf("Title: got %s, want 혈당 관리", got.Title)
	}
}

// ============================================================================
// TestGetPost_NotFound — 게시글 미존재
// ============================================================================

func TestGetPost_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetPost(ctx, "non-existent-id")
	if err == nil {
		t.Fatal("존재하지 않는 게시글 조회가 성공했습니다")
	}
}

// ============================================================================
// TestListPosts_Success — 게시글 목록 조회
// ============================================================================

func TestListPosts_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, _ = svc.CreatePost(ctx, "user-1", "홍길동", service.CatGeneral, "게시글1", "내용1", nil)
	_, _ = svc.CreatePost(ctx, "user-2", "김철수", service.CatQNA, "게시글2", "내용2", nil)
	_, _ = svc.CreatePost(ctx, "user-3", "이영희", service.CatHealthTip, "게시글3", "내용3", nil)

	posts, total, err := svc.ListPosts(ctx, service.CatUnknown, 20, 0)
	if err != nil {
		t.Fatalf("ListPosts 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("전체 게시글 수: got %d, want 3", total)
	}
	if len(posts) != 3 {
		t.Errorf("반환 게시글 수: got %d, want 3", len(posts))
	}

	// 카테고리 필터
	tipPosts, tipTotal, err := svc.ListPosts(ctx, service.CatHealthTip, 20, 0)
	if err != nil {
		t.Fatalf("ListPosts(Tip) 실패: %v", err)
	}
	if tipTotal != 1 {
		t.Errorf("Tip 게시글 수: got %d, want 1", tipTotal)
	}
	if len(tipPosts) != 1 {
		t.Errorf("Tip 반환 수: got %d, want 1", len(tipPosts))
	}
}

// ============================================================================
// TestLikePost_Success — 게시글 좋아요
// ============================================================================

func TestLikePost_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	post, _ := svc.CreatePost(ctx, "user-1", "홍길동", service.CatGeneral, "좋아요 테스트", "내용", nil)

	count, err := svc.LikePost(ctx, post.ID, "user-2")
	if err != nil {
		t.Fatalf("LikePost 실패: %v", err)
	}
	if count != 1 {
		t.Errorf("LikeCount: got %d, want 1", count)
	}

	count2, err := svc.LikePost(ctx, post.ID, "user-3")
	if err != nil {
		t.Fatalf("LikePost(2nd) 실패: %v", err)
	}
	if count2 != 2 {
		t.Errorf("LikeCount(2nd): got %d, want 2", count2)
	}
}

// ============================================================================
// TestCreateComment_Success — 댓글 작성 성공
// ============================================================================

func TestCreateComment_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	post, _ := svc.CreatePost(ctx, "user-1", "홍길동", service.CatGeneral, "댓글 테스트", "내용", nil)

	comment, err := svc.CreateComment(ctx, post.ID, "user-2", "김철수", "좋은 글이네요!")
	if err != nil {
		t.Fatalf("CreateComment 실패: %v", err)
	}
	if comment.ID == "" {
		t.Fatal("Comment ID가 비어 있음")
	}
	if comment.PostID != post.ID {
		t.Errorf("PostID: got %s, want %s", comment.PostID, post.ID)
	}
	if comment.AuthorUserID != "user-2" {
		t.Errorf("AuthorUserID: got %s, want user-2", comment.AuthorUserID)
	}
	if comment.Content != "좋은 글이네요!" {
		t.Errorf("Content: got %s, want 좋은 글이네요!", comment.Content)
	}
	if comment.ParentCommentID != "" {
		t.Errorf("ParentCommentID: got %s, want empty", comment.ParentCommentID)
	}
}

// ============================================================================
// TestListComments_Success — 댓글 목록 조회
// ============================================================================

func TestListComments_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	post, _ := svc.CreatePost(ctx, "user-1", "홍길동", service.CatGeneral, "댓글 목록 테스트", "내용", nil)

	_, _ = svc.CreateComment(ctx, post.ID, "user-2", "김철수", "댓글1")
	_, _ = svc.CreateComment(ctx, post.ID, "user-3", "이영희", "댓글2")
	_, _ = svc.CreateComment(ctx, post.ID, "user-4", "박지성", "댓글3")

	comments, total, err := svc.ListComments(ctx, post.ID, 20, 0)
	if err != nil {
		t.Fatalf("ListComments 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("전체 댓글 수: got %d, want 3", total)
	}
	if len(comments) != 3 {
		t.Errorf("반환 댓글 수: got %d, want 3", len(comments))
	}
}

// ============================================================================
// TestCreateChallenge_Success — 챌린지 생성 성공
// ============================================================================

func TestCreateChallenge_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	startDate := time.Now().UTC()
	challenge, err := svc.CreateChallenge(ctx, "user-1", service.ChTypeExercise, "만보 챌린지", "매일 만보 걷기", "하루 10,000보 걷기", 10000, "steps", 50, 30, startDate)
	if err != nil {
		t.Fatalf("CreateChallenge 실패: %v", err)
	}
	if challenge.ID == "" {
		t.Fatal("Challenge ID가 비어 있음")
	}
	if challenge.CreatorUserID != "user-1" {
		t.Errorf("CreatorUserID: got %s, want user-1", challenge.CreatorUserID)
	}
	if challenge.Type != service.ChTypeExercise {
		t.Errorf("Type: got %d, want %d", challenge.Type, service.ChTypeExercise)
	}
	if challenge.Title != "만보 챌린지" {
		t.Errorf("Title: got %s, want 만보 챌린지", challenge.Title)
	}
	if challenge.ParticipantCount != 1 {
		t.Errorf("ParticipantCount: got %d, want 1 (생성자 자동 참가)", challenge.ParticipantCount)
	}
	if !challenge.Participants["user-1"] {
		t.Error("생성자가 참가자 목록에 없음")
	}
	if challenge.DurationDays != 30 {
		t.Errorf("DurationDays: got %d, want 30", challenge.DurationDays)
	}
	if challenge.MaxParticipants != 50 {
		t.Errorf("MaxParticipants: got %d, want 50", challenge.MaxParticipants)
	}
	if challenge.Status != service.ChActive {
		t.Errorf("Status: got %d, want %d (Active)", challenge.Status, service.ChActive)
	}
}

// ============================================================================
// TestGetChallenge_Success — 챌린지 조회 성공
// ============================================================================

func TestGetChallenge_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateChallenge(ctx, "user-1", service.ChTypeMeasurement, "혈당 챌린지", "매일 혈당 측정", "하루 3회 측정", 3, "회", 20, 14, time.Now().UTC())

	got, err := svc.GetChallenge(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetChallenge 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("Challenge ID: got %s, want %s", got.ID, created.ID)
	}
	if got.Title != "혈당 챌린지" {
		t.Errorf("Title: got %s, want 혈당 챌린지", got.Title)
	}
}

// ============================================================================
// TestJoinChallenge_Success — 챌린지 참가 성공
// ============================================================================

func TestJoinChallenge_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	challenge, _ := svc.CreateChallenge(ctx, "creator-1", service.ChTypeExercise, "운동 챌린지", "매일 운동", "30분 운동", 30, "분", 100, 30, time.Now().UTC())

	updated, err := svc.JoinChallenge(ctx, challenge.ID, "user-2")
	if err != nil {
		t.Fatalf("JoinChallenge 실패: %v", err)
	}
	if updated.ParticipantCount != 2 {
		t.Errorf("ParticipantCount: got %d, want 2", updated.ParticipantCount)
	}
	if !updated.Participants["user-2"] {
		t.Error("user-2가 참가자 목록에 없음")
	}
	if !updated.Participants["creator-1"] {
		t.Error("creator-1이 참가자 목록에서 사라짐")
	}
}

// ============================================================================
// TestJoinChallenge_AlreadyJoined — 중복 참가
// ============================================================================

func TestJoinChallenge_AlreadyJoined(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	challenge, _ := svc.CreateChallenge(ctx, "creator-1", service.ChTypeExercise, "운동 챌린지", "매일 운동", "30분 운동", 30, "분", 100, 30, time.Now().UTC())

	// creator-1은 이미 참가한 상태
	_, err := svc.JoinChallenge(ctx, challenge.ID, "creator-1")
	if err == nil {
		t.Fatal("이미 참가한 사용자가 다시 참가할 때 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestListChallenges_Success — 챌린지 목록 조회
// ============================================================================

func TestListChallenges_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, _ = svc.CreateChallenge(ctx, "user-1", service.ChTypeExercise, "챌린지1", "설명1", "목표1", 100, "steps", 50, 30, time.Now().UTC())
	_, _ = svc.CreateChallenge(ctx, "user-2", service.ChTypeDiet, "챌린지2", "설명2", "목표2", 5, "회", 30, 14, time.Now().UTC())

	challenges, total, err := svc.ListChallenges(ctx, service.ChStatusUnknown, 20, 0)
	if err != nil {
		t.Fatalf("ListChallenges 실패: %v", err)
	}
	if total != 2 {
		t.Errorf("전체 챌린지 수: got %d, want 2", total)
	}
	if len(challenges) != 2 {
		t.Errorf("반환 챌린지 수: got %d, want 2", len(challenges))
	}
}

// ============================================================================
// TestEndToEnd_CommunityFlow — 커뮤니티 전체 흐름 E2E 테스트
// ============================================================================

func TestEndToEnd_CommunityFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 1. 게시글 작성
	post, err := svc.CreatePost(ctx, "user-e2e", "테스트 유저", service.CatExperience, "혈당 관리 경험담", "3개월간 혈당 관리한 경험을 공유합니다.", []string{"혈당", "경험"})
	if err != nil {
		t.Fatalf("E2E CreatePost 실패: %v", err)
	}

	// 2. 게시글 조회
	gotPost, err := svc.GetPost(ctx, post.ID)
	if err != nil {
		t.Fatalf("E2E GetPost 실패: %v", err)
	}
	if gotPost.Title != post.Title {
		t.Errorf("E2E Title 불일치: got %s, want %s", gotPost.Title, post.Title)
	}

	// 3. 좋아요
	likeCount, err := svc.LikePost(ctx, post.ID, "user-liker")
	if err != nil {
		t.Fatalf("E2E LikePost 실패: %v", err)
	}
	if likeCount != 1 {
		t.Errorf("E2E LikeCount: got %d, want 1", likeCount)
	}

	// 4. 댓글 작성
	comment, err := svc.CreateComment(ctx, post.ID, "user-commenter", "댓글러", "좋은 경험담 감사합니다!")
	if err != nil {
		t.Fatalf("E2E CreateComment 실패: %v", err)
	}
	if comment.PostID != post.ID {
		t.Errorf("E2E Comment PostID: got %s, want %s", comment.PostID, post.ID)
	}

	// 5. 댓글 목록 조회
	comments, commentTotal, err := svc.ListComments(ctx, post.ID, 20, 0)
	if err != nil {
		t.Fatalf("E2E ListComments 실패: %v", err)
	}
	if commentTotal != 1 {
		t.Errorf("E2E 댓글 수: got %d, want 1", commentTotal)
	}
	if len(comments) != 1 {
		t.Errorf("E2E 반환 댓글 수: got %d, want 1", len(comments))
	}

	// 6. 게시글 목록 조회
	posts, postTotal, err := svc.ListPosts(ctx, service.CatUnknown, 20, 0)
	if err != nil {
		t.Fatalf("E2E ListPosts 실패: %v", err)
	}
	if postTotal != 1 {
		t.Errorf("E2E 게시글 수: got %d, want 1", postTotal)
	}
	if len(posts) != 1 {
		t.Errorf("E2E 반환 게시글 수: got %d, want 1", len(posts))
	}

	// 7. 챌린지 생성
	challenge, err := svc.CreateChallenge(ctx, "user-e2e", service.ChTypeExercise, "E2E 만보 챌린지", "매일 만보 걷기", "10,000보", 10000, "steps", 50, 30, time.Now().UTC())
	if err != nil {
		t.Fatalf("E2E CreateChallenge 실패: %v", err)
	}

	// 8. 챌린지 조회
	gotChallenge, err := svc.GetChallenge(ctx, challenge.ID)
	if err != nil {
		t.Fatalf("E2E GetChallenge 실패: %v", err)
	}
	if gotChallenge.Title != "E2E 만보 챌린지" {
		t.Errorf("E2E Challenge Title: got %s, want E2E 만보 챌린지", gotChallenge.Title)
	}

	// 9. 챌린지 참가
	joined, err := svc.JoinChallenge(ctx, challenge.ID, "user-joiner")
	if err != nil {
		t.Fatalf("E2E JoinChallenge 실패: %v", err)
	}
	if joined.ParticipantCount != 2 {
		t.Errorf("E2E ParticipantCount: got %d, want 2", joined.ParticipantCount)
	}

	// 10. 챌린지 목록 조회
	challenges, challengeTotal, err := svc.ListChallenges(ctx, service.ChStatusUnknown, 20, 0)
	if err != nil {
		t.Fatalf("E2E ListChallenges 실패: %v", err)
	}
	if challengeTotal != 1 {
		t.Errorf("E2E 챌린지 수: got %d, want 1", challengeTotal)
	}
	if len(challenges) != 1 {
		t.Errorf("E2E 반환 챌린지 수: got %d, want 1", len(challenges))
	}
}
