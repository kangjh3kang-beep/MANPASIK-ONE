// Package handler는 community-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/community-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CommunityHandler는 CommunityService gRPC 서버를 구현합니다.
type CommunityHandler struct {
	v1.UnimplementedCommunityServiceServer
	svc *service.CommunityService
	log *zap.Logger
}

// NewCommunityHandler는 CommunityHandler를 생성합니다.
func NewCommunityHandler(svc *service.CommunityService, log *zap.Logger) *CommunityHandler {
	return &CommunityHandler{svc: svc, log: log}
}

// ============================================================================
// CreatePost — 게시글 작성
// ============================================================================

func (h *CommunityHandler) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.Post, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	post, err := h.svc.CreatePost(
		ctx,
		req.AuthorId,
		"", // display name resolved internally
		protoPostCategoryToService(req.Category),
		req.Title,
		req.Content,
		req.Tags,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return postToProto(post), nil
}

// ============================================================================
// GetPost — 게시글 조회
// ============================================================================

func (h *CommunityHandler) GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.Post, error) {
	if req == nil || req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id는 필수입니다")
	}

	post, err := h.svc.GetPost(ctx, req.PostId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return postToProto(post), nil
}

// ============================================================================
// ListPosts — 게시글 목록 조회
// ============================================================================

func (h *CommunityHandler) ListPosts(ctx context.Context, req *v1.ListPostsRequest) (*v1.ListPostsResponse, error) {
	if req == nil {
		req = &v1.ListPostsRequest{}
	}

	posts, total, err := h.svc.ListPosts(
		ctx,
		protoPostCategoryToService(req.Category),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoPosts := make([]*v1.Post, 0, len(posts))
	for _, p := range posts {
		protoPosts = append(protoPosts, postToProto(p))
	}

	return &v1.ListPostsResponse{
		Posts:      protoPosts,
		TotalCount: total,
	}, nil
}

// ============================================================================
// LikePost — 게시글 좋아요
// ============================================================================

func (h *CommunityHandler) LikePost(ctx context.Context, req *v1.LikePostRequest) (*v1.LikePostResponse, error) {
	if req == nil || req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id는 필수입니다")
	}

	newCount, err := h.svc.LikePost(ctx, req.PostId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.LikePostResponse{
		Success:   true,
		LikeCount: newCount,
		IsLiked:   true,
	}, nil
}

// ============================================================================
// CreateComment — 댓글 작성
// ============================================================================

func (h *CommunityHandler) CreateComment(ctx context.Context, req *v1.CreateCommentRequest) (*v1.Comment, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	comment, err := h.svc.CreateComment(
		ctx,
		req.PostId,
		req.AuthorId,
		"", // display name resolved internally
		req.Content,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return commentToProto(comment), nil
}

// ============================================================================
// ListComments — 댓글 목록 조회
// ============================================================================

func (h *CommunityHandler) ListComments(ctx context.Context, req *v1.ListCommentsRequest) (*v1.ListCommentsResponse, error) {
	if req == nil || req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post_id는 필수입니다")
	}

	comments, total, err := h.svc.ListComments(ctx, req.PostId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoComments := make([]*v1.Comment, 0, len(comments))
	for _, c := range comments {
		protoComments = append(protoComments, commentToProto(c))
	}

	return &v1.ListCommentsResponse{
		Comments:   protoComments,
		TotalCount: total,
	}, nil
}

// ============================================================================
// CreateChallenge — 챌린지 생성
// ============================================================================

func (h *CommunityHandler) CreateChallenge(ctx context.Context, req *v1.CreateChallengeRequest) (*v1.Challenge, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()
	durationDays := int(endDate.Sub(startDate).Hours() / 24)

	challenge, err := h.svc.CreateChallenge(
		ctx,
		req.CreatorId,
		protoChallengeTypeToService(req.ChallengeType),
		req.Title,
		req.Description,
		"",
		float64(req.TargetValue),
		req.Unit,
		int(req.MaxParticipants),
		durationDays,
		startDate,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return challengeToProto(challenge), nil
}

// ============================================================================
// GetChallenge — 챌린지 조회
// ============================================================================

func (h *CommunityHandler) GetChallenge(ctx context.Context, req *v1.GetChallengeRequest) (*v1.Challenge, error) {
	if req == nil || req.ChallengeId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id는 필수입니다")
	}

	challenge, err := h.svc.GetChallenge(ctx, req.ChallengeId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return challengeToProto(challenge), nil
}

// ============================================================================
// JoinChallenge — 챌린지 참가
// ============================================================================

func (h *CommunityHandler) JoinChallenge(ctx context.Context, req *v1.JoinChallengeRequest) (*v1.JoinChallengeResponse, error) {
	if req == nil || req.ChallengeId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id는 필수입니다")
	}

	challenge, err := h.svc.JoinChallenge(ctx, req.ChallengeId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.JoinChallengeResponse{
		Success:          true,
		ParticipantCount: int32(challenge.ParticipantCount),
	}, nil
}

// ============================================================================
// ListChallenges — 챌린지 목록 조회
// ============================================================================

func (h *CommunityHandler) ListChallenges(ctx context.Context, req *v1.ListChallengesRequest) (*v1.ListChallengesResponse, error) {
	if req == nil {
		req = &v1.ListChallengesRequest{}
	}

	challenges, total, err := h.svc.ListChallenges(
		ctx,
		protoChallengeStatusToService(req.StatusFilter),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoChallenges := make([]*v1.Challenge, 0, len(challenges))
	for _, ch := range challenges {
		protoChallenges = append(protoChallenges, challengeToProto(ch))
	}

	return &v1.ListChallengesResponse{
		Challenges: protoChallenges,
		TotalCount: total,
	}, nil
}

// ============================================================================
// GetChallengeLeaderboard — 챌린지 리더보드 조회
// ============================================================================

func (h *CommunityHandler) GetChallengeLeaderboard(ctx context.Context, req *v1.GetChallengeLeaderboardRequest) (*v1.GetChallengeLeaderboardResponse, error) {
	if req == nil || req.ChallengeId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id는 필수입니다")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	entries, totalParticipants, myEntry, err := h.svc.GetChallengeLeaderboard(ctx, req.ChallengeId, limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoEntries := make([]*v1.LeaderboardEntry, 0, len(entries))
	for _, e := range entries {
		protoEntries = append(protoEntries, leaderboardEntryToProto(e))
	}

	resp := &v1.GetChallengeLeaderboardResponse{
		Entries:           protoEntries,
		TotalParticipants: totalParticipants,
	}
	if myEntry != nil {
		resp.MyEntry = leaderboardEntryToProto(myEntry)
	}

	return resp, nil
}

// ============================================================================
// UpdateChallengeProgress — 챌린지 진행도 업데이트
// ============================================================================

func (h *CommunityHandler) UpdateChallengeProgress(ctx context.Context, req *v1.UpdateChallengeProgressRequest) (*v1.UpdateChallengeProgressResponse, error) {
	if req == nil || req.ChallengeId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id와 user_id는 필수입니다")
	}

	newProgress, targetValue, newRank, err := h.svc.UpdateChallengeProgress(ctx, req.ChallengeId, req.UserId, req.Value)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.UpdateChallengeProgressResponse{
		Success:     true,
		NewProgress: newProgress,
		TargetValue: targetValue,
		NewRank:     newRank,
	}, nil
}

// ============================================================================
// 헬퍼 함수 — Proto ↔ Service 변환
// ============================================================================

// --- Post 변환 ---

func postToProto(p *service.Post) *v1.Post {
	return &v1.Post{
		PostId:       p.ID,
		AuthorId:     p.AuthorUserID,
		AuthorName:   p.AuthorDisplayName,
		Category:     servicePostCategoryToProto(p.Category),
		Title:        p.Title,
		Content:      p.Content,
		Tags:         p.Tags,
		LikeCount:    int32(p.LikeCount),
		CommentCount: int32(p.CommentCount),
		ViewCount:    int32(p.ViewCount),
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}
}

// --- Comment 변환 ---

func commentToProto(c *service.Comment) *v1.Comment {
	return &v1.Comment{
		CommentId:       c.ID,
		PostId:          c.PostID,
		AuthorId:        c.AuthorUserID,
		AuthorName:      c.AuthorDisplayName,
		Content:         c.Content,
		ParentCommentId: c.ParentCommentID,
		LikeCount:       int32(c.LikeCount),
		CreatedAt:       timestamppb.New(c.CreatedAt),
	}
}

// --- Challenge 변환 ---

func challengeToProto(ch *service.Challenge) *v1.Challenge {
	return &v1.Challenge{
		ChallengeId:      ch.ID,
		CreatorId:        ch.CreatorUserID,
		Title:            ch.Title,
		Description:      ch.Description,
		ChallengeType:    serviceChallengeTypeToProto(ch.Type),
		Status:           serviceChallengeStatusToProto(ch.Status),
		TargetValue:      ch.TargetValue,
		Unit:             ch.TargetUnit,
		ParticipantCount: int32(ch.ParticipantCount),
		MaxParticipants:  int32(ch.MaxParticipants),
		StartDate:        timestamppb.New(ch.StartDate),
		EndDate:          timestamppb.New(ch.EndDate),
		CreatedAt:        timestamppb.New(ch.CreatedAt),
	}
}

// --- LeaderboardEntry 변환 ---

func leaderboardEntryToProto(e *service.LeaderboardEntry) *v1.LeaderboardEntry {
	return &v1.LeaderboardEntry{
		Rank:          e.Rank,
		UserId:        e.UserID,
		DisplayName:   e.DisplayName,
		AvatarUrl:     e.AvatarURL,
		ProgressValue: e.ProgressValue,
		TargetValue:   e.TargetValue,
		ProgressPct:   e.ProgressPct,
		StreakDays:     e.StreakDays,
		LastUpdated:   timestamppb.New(e.LastUpdated),
	}
}

// --- PostCategory 변환 ---

func protoPostCategoryToService(c v1.PostCategory) service.PostCategory {
	switch c {
	case v1.PostCategory_POST_CATEGORY_GENERAL:
		return service.CatGeneral
	case v1.PostCategory_POST_CATEGORY_HEALTH_TIP:
		return service.CatHealthTip
	case v1.PostCategory_POST_CATEGORY_QNA:
		return service.CatQNA
	case v1.PostCategory_POST_CATEGORY_EXPERIENCE:
		return service.CatExperience
	case v1.PostCategory_POST_CATEGORY_RECIPE:
		return service.CatRecipe
	case v1.PostCategory_POST_CATEGORY_EXERCISE:
		return service.CatExercise
	default:
		return service.CatUnknown
	}
}

func servicePostCategoryToProto(c service.PostCategory) v1.PostCategory {
	switch c {
	case service.CatGeneral:
		return v1.PostCategory_POST_CATEGORY_GENERAL
	case service.CatHealthTip:
		return v1.PostCategory_POST_CATEGORY_HEALTH_TIP
	case service.CatQNA:
		return v1.PostCategory_POST_CATEGORY_QNA
	case service.CatExperience:
		return v1.PostCategory_POST_CATEGORY_EXPERIENCE
	case service.CatRecipe:
		return v1.PostCategory_POST_CATEGORY_RECIPE
	case service.CatExercise:
		return v1.PostCategory_POST_CATEGORY_EXERCISE
	default:
		return v1.PostCategory_POST_CATEGORY_UNKNOWN
	}
}

// --- ChallengeStatus 변환 ---

func protoChallengeStatusToService(s v1.ChallengeStatus) service.ChallengeStatus {
	switch s {
	case v1.ChallengeStatus_CHALLENGE_STATUS_UPCOMING:
		return service.ChUpcoming
	case v1.ChallengeStatus_CHALLENGE_STATUS_ACTIVE:
		return service.ChActive
	case v1.ChallengeStatus_CHALLENGE_STATUS_COMPLETED:
		return service.ChCompleted
	case v1.ChallengeStatus_CHALLENGE_STATUS_CANCELLED:
		return service.ChCancelled
	default:
		return service.ChStatusUnknown
	}
}

func serviceChallengeStatusToProto(s service.ChallengeStatus) v1.ChallengeStatus {
	switch s {
	case service.ChUpcoming:
		return v1.ChallengeStatus_CHALLENGE_STATUS_UPCOMING
	case service.ChActive:
		return v1.ChallengeStatus_CHALLENGE_STATUS_ACTIVE
	case service.ChCompleted:
		return v1.ChallengeStatus_CHALLENGE_STATUS_COMPLETED
	case service.ChCancelled:
		return v1.ChallengeStatus_CHALLENGE_STATUS_CANCELLED
	default:
		return v1.ChallengeStatus_CHALLENGE_STATUS_UNKNOWN
	}
}

// --- ChallengeType 변환 ---

func protoChallengeTypeToService(t v1.ChallengeType) service.ChallengeType {
	switch t {
	case v1.ChallengeType_CHALLENGE_TYPE_STEPS:
		return service.ChTypeSteps
	case v1.ChallengeType_CHALLENGE_TYPE_MEASUREMENT:
		return service.ChTypeMeasurement
	case v1.ChallengeType_CHALLENGE_TYPE_DIET:
		return service.ChTypeDiet
	case v1.ChallengeType_CHALLENGE_TYPE_EXERCISE:
		return service.ChTypeExercise
	case v1.ChallengeType_CHALLENGE_TYPE_SLEEP:
		return service.ChTypeSleep
	default:
		return service.ChTypeUnknown
	}
}

func serviceChallengeTypeToProto(t service.ChallengeType) v1.ChallengeType {
	switch t {
	case service.ChTypeSteps:
		return v1.ChallengeType_CHALLENGE_TYPE_STEPS
	case service.ChTypeMeasurement:
		return v1.ChallengeType_CHALLENGE_TYPE_MEASUREMENT
	case service.ChTypeDiet:
		return v1.ChallengeType_CHALLENGE_TYPE_DIET
	case service.ChTypeExercise:
		return v1.ChallengeType_CHALLENGE_TYPE_EXERCISE
	case service.ChTypeSleep:
		return v1.ChallengeType_CHALLENGE_TYPE_SLEEP
	default:
		return v1.ChallengeType_CHALLENGE_TYPE_UNKNOWN
	}
}

// --- 에러 변환 ---

// toGRPC는 AppError를 gRPC status로 변환합니다.
func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
