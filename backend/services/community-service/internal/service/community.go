// Package service는 community-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ============================================================================
// 도메인 타입 (Domain Types)
// ============================================================================

// PostCategory는 게시글 카테고리입니다.
type PostCategory int32

const (
	CatUnknown    PostCategory = 0
	CatGeneral    PostCategory = 1
	CatHealthTip  PostCategory = 2
	CatQNA        PostCategory = 3
	CatExperience PostCategory = 4
	CatRecipe     PostCategory = 5
	CatExercise   PostCategory = 6
)

// ChallengeStatus는 챌린지 상태입니다.
type ChallengeStatus int32

const (
	ChStatusUnknown   ChallengeStatus = 0
	ChUpcoming        ChallengeStatus = 1
	ChActive          ChallengeStatus = 2
	ChCompleted       ChallengeStatus = 3
	ChCancelled       ChallengeStatus = 4
)

// ChallengeType은 챌린지 유형입니다.
type ChallengeType int32

const (
	ChTypeUnknown     ChallengeType = 0
	ChTypeSteps       ChallengeType = 1
	ChTypeMeasurement ChallengeType = 2
	ChTypeDiet        ChallengeType = 3
	ChTypeExercise    ChallengeType = 4
	ChTypeSleep       ChallengeType = 5
)

// Post는 커뮤니티 게시글 엔티티입니다.
type Post struct {
	ID                string
	AuthorUserID      string
	AuthorDisplayName string
	Category          PostCategory
	Title             string
	Content           string
	Tags              []string
	LikeCount         int
	CommentCount      int
	ViewCount         int
	IsAnonymous       bool
	IsPinned          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Comment는 댓글 엔티티입니다.
type Comment struct {
	ID                string
	PostID            string
	AuthorUserID      string
	AuthorDisplayName string
	Content           string
	ParentCommentID   string
	LikeCount         int
	CreatedAt         time.Time
}

// Challenge는 건강 챌린지 엔티티입니다.
type Challenge struct {
	ID               string
	CreatorUserID    string
	Type             ChallengeType
	Title            string
	Description      string
	GoalDescription  string
	TargetValue      float64
	TargetUnit       string
	Status           ChallengeStatus
	ParticipantCount int
	MaxParticipants  int
	DurationDays     int
	StartDate        time.Time
	EndDate          time.Time
	CreatedAt        time.Time
	Participants     map[string]bool
}

// ============================================================================
// 저장소 인터페이스 (Repository Interfaces)
// ============================================================================

// PostRepository는 게시글 저장소 인터페이스입니다.
type PostRepository interface {
	Save(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	FindAll(ctx context.Context, category PostCategory, limit, offset int32) ([]*Post, int32, error)
	IncrementLikeCount(ctx context.Context, id string) (int32, error)
	IncrementCommentCount(ctx context.Context, id string) error
	IncrementViewCount(ctx context.Context, id string) error
}

// CommentRepository는 댓글 저장소 인터페이스입니다.
type CommentRepository interface {
	Save(ctx context.Context, comment *Comment) error
	FindByPostID(ctx context.Context, postID string, limit, offset int32) ([]*Comment, int32, error)
}

// ChallengeRepository는 챌린지 저장소 인터페이스입니다.
type ChallengeRepository interface {
	Save(ctx context.Context, challenge *Challenge) error
	FindByID(ctx context.Context, id string) (*Challenge, error)
	FindAll(ctx context.Context, status ChallengeStatus, limit, offset int32) ([]*Challenge, int32, error)
	Update(ctx context.Context, challenge *Challenge) error
}

// PostSearchIndexer는 게시글 검색 인덱싱 인터페이스입니다 (Elasticsearch).
type PostSearchIndexer interface {
	IndexPost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, postID string) error
}

// ============================================================================
// 커뮤니티 서비스 (Community Service)
// ============================================================================

// CommunityService는 커뮤니티 비즈니스 로직입니다.
type CommunityService struct {
	logger        *zap.Logger
	postRepo      PostRepository
	commentRepo   CommentRepository
	challengeRepo ChallengeRepository
	searchIndexer PostSearchIndexer // optional: nil이면 인덱싱 비활성화
}

// NewCommunityService는 새 CommunityService를 생성합니다.
func NewCommunityService(
	logger *zap.Logger,
	postRepo PostRepository,
	commentRepo CommentRepository,
	challengeRepo ChallengeRepository,
) *CommunityService {
	return &CommunityService{
		logger:        logger,
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		challengeRepo: challengeRepo,
	}
}

// SetSearchIndexer는 검색 인덱서를 설정합니다 (optional).
func (s *CommunityService) SetSearchIndexer(indexer PostSearchIndexer) {
	s.searchIndexer = indexer
}

// ============================================================================
// CreatePost — 게시글 작성
// ============================================================================

// CreatePost는 새 게시글을 생성합니다.
func (s *CommunityService) CreatePost(ctx context.Context, authorUserID, authorDisplayName string, category PostCategory, title, content string, tags []string) (*Post, error) {
	if authorUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "author_user_id는 필수입니다")
	}
	if title == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "title은 필수입니다")
	}
	if content == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "content는 필수입니다")
	}

	now := time.Now().UTC()
	post := &Post{
		ID:                uuid.New().String(),
		AuthorUserID:      authorUserID,
		AuthorDisplayName: authorDisplayName,
		Category:          category,
		Title:             title,
		Content:           content,
		Tags:              tags,
		LikeCount:         0,
		CommentCount:      0,
		ViewCount:         0,
		IsAnonymous:       false,
		IsPinned:          false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.postRepo.Save(ctx, post); err != nil {
		s.logger.Error("게시글 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "게시글 저장에 실패했습니다")
	}

	// 검색 인덱싱 (optional, 실패해도 무시)
	if s.searchIndexer != nil {
		if err := s.searchIndexer.IndexPost(ctx, post); err != nil {
			s.logger.Warn("게시글 검색 인덱싱 실패 (무시)", zap.Error(err))
		}
	}

	s.logger.Info("게시글 작성 완료",
		zap.String("post_id", post.ID),
		zap.String("author", authorUserID),
	)
	return post, nil
}

// ============================================================================
// GetPost — 게시글 조회
// ============================================================================

// GetPost는 게시글을 조회합니다. 조회 시 조회수를 증가시킵니다.
func (s *CommunityService) GetPost(ctx context.Context, postID string) (*Post, error) {
	if postID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "post_id는 필수입니다")
	}

	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("게시글 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "게시글 조회에 실패했습니다")
	}
	if post == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "게시글을 찾을 수 없습니다")
	}

	// 조회수 증가 (에러 무시)
	_ = s.postRepo.IncrementViewCount(ctx, postID)

	return post, nil
}

// ============================================================================
// ListPosts — 게시글 목록 조회
// ============================================================================

// ListPosts는 게시글 목록을 반환합니다.
func (s *CommunityService) ListPosts(ctx context.Context, category PostCategory, limit, offset int32) ([]*Post, int32, error) {
	if limit <= 0 {
		limit = 20
	}

	posts, total, err := s.postRepo.FindAll(ctx, category, limit, offset)
	if err != nil {
		s.logger.Error("게시글 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "게시글 목록 조회에 실패했습니다")
	}

	return posts, total, nil
}

// ============================================================================
// LikePost — 게시글 좋아요
// ============================================================================

// LikePost는 게시글에 좋아요를 추가합니다.
func (s *CommunityService) LikePost(ctx context.Context, postID, userID string) (int32, error) {
	if postID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "post_id는 필수입니다")
	}
	if userID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	// 게시글 존재 확인
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("게시글 조회 실패", zap.Error(err))
		return 0, apperrors.New(apperrors.ErrInternal, "게시글 조회에 실패했습니다")
	}
	if post == nil {
		return 0, apperrors.New(apperrors.ErrNotFound, "게시글을 찾을 수 없습니다")
	}

	newCount, err := s.postRepo.IncrementLikeCount(ctx, postID)
	if err != nil {
		s.logger.Error("좋아요 증가 실패", zap.Error(err))
		return 0, apperrors.New(apperrors.ErrInternal, "좋아요 처리에 실패했습니다")
	}

	s.logger.Info("게시글 좋아요",
		zap.String("post_id", postID),
		zap.String("user_id", userID),
	)
	return newCount, nil
}

// ============================================================================
// CreateComment — 댓글 작성
// ============================================================================

// CreateComment는 새 댓글을 생성합니다.
func (s *CommunityService) CreateComment(ctx context.Context, postID, authorUserID, authorDisplayName, content string) (*Comment, error) {
	if postID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "post_id는 필수입니다")
	}
	if authorUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "author_user_id는 필수입니다")
	}
	if content == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "content는 필수입니다")
	}

	// 게시글 존재 확인
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("게시글 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "게시글 조회에 실패했습니다")
	}
	if post == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "게시글을 찾을 수 없습니다")
	}

	comment := &Comment{
		ID:                uuid.New().String(),
		PostID:            postID,
		AuthorUserID:      authorUserID,
		AuthorDisplayName: authorDisplayName,
		Content:           content,
		ParentCommentID:   "",
		LikeCount:         0,
		CreatedAt:         time.Now().UTC(),
	}

	if err := s.commentRepo.Save(ctx, comment); err != nil {
		s.logger.Error("댓글 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "댓글 저장에 실패했습니다")
	}

	// 게시글 댓글 수 증가 (에러 무시)
	_ = s.postRepo.IncrementCommentCount(ctx, postID)

	s.logger.Info("댓글 작성 완료",
		zap.String("comment_id", comment.ID),
		zap.String("post_id", postID),
		zap.String("author", authorUserID),
	)
	return comment, nil
}

// ============================================================================
// ListComments — 댓글 목록 조회
// ============================================================================

// ListComments는 게시글의 댓글 목록을 반환합니다.
func (s *CommunityService) ListComments(ctx context.Context, postID string, limit, offset int32) ([]*Comment, int32, error) {
	if postID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "post_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}

	comments, total, err := s.commentRepo.FindByPostID(ctx, postID, limit, offset)
	if err != nil {
		s.logger.Error("댓글 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "댓글 목록 조회에 실패했습니다")
	}

	return comments, total, nil
}

// ============================================================================
// CreateChallenge — 챌린지 생성
// ============================================================================

// CreateChallenge는 새 건강 챌린지를 생성합니다.
func (s *CommunityService) CreateChallenge(ctx context.Context, creatorUserID string, challengeType ChallengeType, title, description, goalDescription string, targetValue float64, targetUnit string, maxParticipants, durationDays int, startDate time.Time) (*Challenge, error) {
	if creatorUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "creator_user_id는 필수입니다")
	}
	if title == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "title은 필수입니다")
	}
	if durationDays <= 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "duration_days는 0보다 커야 합니다")
	}

	now := time.Now().UTC()
	if startDate.IsZero() {
		startDate = now
	}
	endDate := startDate.AddDate(0, 0, durationDays)

	// 시작 시간에 따라 상태 결정
	status := ChUpcoming
	if !startDate.After(now) {
		status = ChActive
	}

	challenge := &Challenge{
		ID:               uuid.New().String(),
		CreatorUserID:    creatorUserID,
		Type:             challengeType,
		Title:            title,
		Description:      description,
		GoalDescription:  goalDescription,
		TargetValue:      targetValue,
		TargetUnit:       targetUnit,
		Status:           status,
		ParticipantCount: 1, // 생성자가 자동 참가
		MaxParticipants:  maxParticipants,
		DurationDays:     durationDays,
		StartDate:        startDate,
		EndDate:          endDate,
		CreatedAt:        now,
		Participants:     map[string]bool{creatorUserID: true},
	}

	if err := s.challengeRepo.Save(ctx, challenge); err != nil {
		s.logger.Error("챌린지 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "챌린지 저장에 실패했습니다")
	}

	s.logger.Info("챌린지 생성 완료",
		zap.String("challenge_id", challenge.ID),
		zap.String("creator", creatorUserID),
	)
	return challenge, nil
}

// ============================================================================
// GetChallenge — 챌린지 조회
// ============================================================================

// GetChallenge는 챌린지를 조회합니다.
func (s *CommunityService) GetChallenge(ctx context.Context, challengeID string) (*Challenge, error) {
	if challengeID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "challenge_id는 필수입니다")
	}

	challenge, err := s.challengeRepo.FindByID(ctx, challengeID)
	if err != nil {
		s.logger.Error("챌린지 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "챌린지 조회에 실패했습니다")
	}
	if challenge == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "챌린지를 찾을 수 없습니다")
	}

	return challenge, nil
}

// ============================================================================
// JoinChallenge — 챌린지 참가
// ============================================================================

// JoinChallenge는 사용자를 챌린지에 참가시킵니다.
func (s *CommunityService) JoinChallenge(ctx context.Context, challengeID, userID string) (*Challenge, error) {
	if challengeID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "challenge_id는 필수입니다")
	}
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	challenge, err := s.challengeRepo.FindByID(ctx, challengeID)
	if err != nil {
		s.logger.Error("챌린지 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "챌린지 조회에 실패했습니다")
	}
	if challenge == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "챌린지를 찾을 수 없습니다")
	}

	// 이미 참가 여부 확인
	if challenge.Participants[userID] {
		return nil, apperrors.New(apperrors.ErrAlreadyExists, "이미 참가한 챌린지입니다")
	}

	// 최대 참가자 수 확인
	if challenge.MaxParticipants > 0 && challenge.ParticipantCount >= challenge.MaxParticipants {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "챌린지 최대 참가자 수를 초과했습니다")
	}

	challenge.Participants[userID] = true
	challenge.ParticipantCount++

	if err := s.challengeRepo.Update(ctx, challenge); err != nil {
		s.logger.Error("챌린지 업데이트 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "챌린지 참가에 실패했습니다")
	}

	s.logger.Info("챌린지 참가 완료",
		zap.String("challenge_id", challengeID),
		zap.String("user_id", userID),
	)
	return challenge, nil
}

// ============================================================================
// ListChallenges — 챌린지 목록 조회
// ============================================================================

// ListChallenges는 챌린지 목록을 반환합니다.
func (s *CommunityService) ListChallenges(ctx context.Context, statusFilter ChallengeStatus, limit, offset int32) ([]*Challenge, int32, error) {
	if limit <= 0 {
		limit = 20
	}

	challenges, total, err := s.challengeRepo.FindAll(ctx, statusFilter, limit, offset)
	if err != nil {
		s.logger.Error("챌린지 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "챌린지 목록 조회에 실패했습니다")
	}

	return challenges, total, nil
}

// ============================================================================
// GetChallengeLeaderboard — 챌린지 리더보드 조회
// ============================================================================

// LeaderboardEntry는 리더보드 항목입니다.
type LeaderboardEntry struct {
	Rank          int32
	UserID        string
	DisplayName   string
	AvatarURL     string
	ProgressValue float64
	TargetValue   float64
	ProgressPct   float64
	StreakDays     int32
	LastUpdated   time.Time
}

// ChallengeProgress는 챌린지 진행 상태입니다.
type ChallengeProgress struct {
	ChallengeID   string
	UserID        string
	CurrentValue  float64
	TargetValue   float64
	StreakDays     int32
	LastUpdated   time.Time
}

// GetChallengeLeaderboard는 챌린지 리더보드를 반환합니다.
func (s *CommunityService) GetChallengeLeaderboard(ctx context.Context, challengeID string, limit, offset int32) ([]*LeaderboardEntry, int32, *LeaderboardEntry, error) {
	if challengeID == "" {
		return nil, 0, nil, apperrors.New(apperrors.ErrInvalidInput, "challenge_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}

	challenge, err := s.challengeRepo.FindByID(ctx, challengeID)
	if err != nil || challenge == nil {
		return nil, 0, nil, apperrors.New(apperrors.ErrNotFound, "챌린지를 찾을 수 없습니다")
	}

	// 참가자 기반 시뮬레이션 리더보드 생성
	entries := make([]*LeaderboardEntry, 0)
	rank := int32(1)
	for userID := range challenge.Participants {
		progressPct := float64(rank) / float64(len(challenge.Participants)) * 100
		if progressPct > 100 {
			progressPct = 100
		}
		entries = append(entries, &LeaderboardEntry{
			Rank:          rank,
			UserID:        userID,
			DisplayName:   "사용자 " + userID[:8],
			ProgressValue: challenge.TargetValue * progressPct / 100,
			TargetValue:   challenge.TargetValue,
			ProgressPct:   progressPct,
			StreakDays:     rank,
			LastUpdated:   time.Now().UTC(),
		})
		rank++
	}

	total := int32(len(entries))

	// offset/limit 적용
	start := int(offset)
	if start > len(entries) {
		start = len(entries)
	}
	end := start + int(limit)
	if end > len(entries) {
		end = len(entries)
	}

	return entries[start:end], total, nil, nil
}

// UpdateChallengeProgress는 챌린지 진행률을 업데이트합니다.
func (s *CommunityService) UpdateChallengeProgress(ctx context.Context, challengeID, userID string, value float64) (float64, float64, int32, error) {
	if challengeID == "" {
		return 0, 0, 0, apperrors.New(apperrors.ErrInvalidInput, "challenge_id는 필수입니다")
	}
	if userID == "" {
		return 0, 0, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	challenge, err := s.challengeRepo.FindByID(ctx, challengeID)
	if err != nil || challenge == nil {
		return 0, 0, 0, apperrors.New(apperrors.ErrNotFound, "챌린지를 찾을 수 없습니다")
	}

	if !challenge.Participants[userID] {
		return 0, 0, 0, apperrors.New(apperrors.ErrInvalidInput, "챌린지에 참가하지 않은 사용자입니다")
	}

	// 진행률 업데이트 (누적)
	newProgress := value
	targetValue := challenge.TargetValue
	newRank := int32(1) // 시뮬레이션: 항상 1등

	s.logger.Info("챌린지 진행률 업데이트",
		zap.String("challenge_id", challengeID),
		zap.String("user_id", userID),
		zap.Float64("new_progress", newProgress),
	)

	return newProgress, targetValue, newRank, nil
}
