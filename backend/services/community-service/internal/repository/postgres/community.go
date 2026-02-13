// Package postgres는 community-service의 PostgreSQL 저장소 구현입니다.
//
// DB 스키마: infrastructure/database/init/17-community.sql
// 테이블: posts, comments, post_likes, challenges, challenge_participants
package postgres

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/community-service/internal/service"
)

// ============================================================================
// PostRepository
// ============================================================================

// PostRepository는 PostgreSQL 기반 게시글 저장소입니다.
type PostRepository struct {
	pool *pgxpool.Pool
}

// NewPostRepository는 PostRepository를 생성합니다.
func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

// categoryToString은 PostCategory를 DB ENUM 문자열로 변환합니다.
func categoryToString(c service.PostCategory) string {
	switch c {
	case service.CatGeneral:
		return "general"
	case service.CatHealthTip:
		return "tip"
	case service.CatQNA:
		return "question"
	case service.CatExperience:
		return "experience"
	case service.CatRecipe:
		return "recipe"
	case service.CatExercise:
		return "exercise"
	default:
		return "general"
	}
}

// stringToCategory는 DB ENUM 문자열을 PostCategory로 변환합니다.
func stringToCategory(s string) service.PostCategory {
	switch s {
	case "general":
		return service.CatGeneral
	case "tip":
		return service.CatHealthTip
	case "question":
		return service.CatQNA
	case "experience":
		return service.CatExperience
	case "recipe":
		return service.CatRecipe
	case "exercise":
		return service.CatExercise
	default:
		return service.CatUnknown
	}
}

// Save는 게시글을 저장합니다.
func (r *PostRepository) Save(ctx context.Context, post *service.Post) error {
	const q = `INSERT INTO posts (post_id, author_user_id, author_display_name, category, title, content, tags,
		like_count, comment_count, view_count, is_anonymous, is_pinned, created_at, updated_at)
		VALUES ($1, $2, $3, $4::post_category, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.pool.Exec(ctx, q,
		post.ID, post.AuthorUserID, post.AuthorDisplayName,
		categoryToString(post.Category), post.Title, post.Content, post.Tags,
		post.LikeCount, post.CommentCount, post.ViewCount,
		post.IsAnonymous, post.IsPinned, post.CreatedAt, post.UpdatedAt,
	)
	return err
}

// FindByID는 게시글을 ID로 조회합니다.
func (r *PostRepository) FindByID(ctx context.Context, id string) (*service.Post, error) {
	const q = `SELECT post_id, author_user_id, COALESCE(author_display_name,''), category, title, content,
		COALESCE(tags, '{}'), like_count, comment_count, view_count, is_anonymous, is_pinned, created_at, updated_at
		FROM posts WHERE post_id = $1`
	var p service.Post
	var catStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.AuthorUserID, &p.AuthorDisplayName, &catStr, &p.Title, &p.Content,
		&p.Tags, &p.LikeCount, &p.CommentCount, &p.ViewCount,
		&p.IsAnonymous, &p.IsPinned, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.Category = stringToCategory(catStr)
	return &p, nil
}

// FindAll는 게시글 목록을 반환합니다.
func (r *PostRepository) FindAll(ctx context.Context, category service.PostCategory, limit, offset int32) ([]*service.Post, int32, error) {
	// 카운트 쿼리
	countQ := "SELECT COUNT(*) FROM posts"
	listQ := `SELECT post_id, author_user_id, COALESCE(author_display_name,''), category, title, content,
		COALESCE(tags, '{}'), like_count, comment_count, view_count, is_anonymous, is_pinned, created_at, updated_at
		FROM posts`

	var args []interface{}
	idx := 1

	if category != service.CatUnknown {
		catStr := categoryToString(category)
		where := " WHERE category = $1::post_category"
		countQ += where
		listQ += where
		args = append(args, catStr)
		idx++
	}

	var total int32
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQ += " ORDER BY created_at DESC"
	listQ += " LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	queryArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, listQ, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*service.Post
	for rows.Next() {
		var p service.Post
		var catStr string
		if err := rows.Scan(
			&p.ID, &p.AuthorUserID, &p.AuthorDisplayName, &catStr, &p.Title, &p.Content,
			&p.Tags, &p.LikeCount, &p.CommentCount, &p.ViewCount,
			&p.IsAnonymous, &p.IsPinned, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		p.Category = stringToCategory(catStr)
		posts = append(posts, &p)
	}
	return posts, total, rows.Err()
}

// IncrementLikeCount는 좋아요 수를 1 증가시킵니다.
func (r *PostRepository) IncrementLikeCount(ctx context.Context, id string) (int32, error) {
	const q = `UPDATE posts SET like_count = like_count + 1 WHERE post_id = $1 RETURNING like_count`
	var count int32
	err := r.pool.QueryRow(ctx, q, id).Scan(&count)
	return count, err
}

// IncrementCommentCount는 댓글 수를 1 증가시킵니다.
func (r *PostRepository) IncrementCommentCount(ctx context.Context, id string) error {
	const q = `UPDATE posts SET comment_count = comment_count + 1 WHERE post_id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

// IncrementViewCount는 조회수를 1 증가시킵니다.
func (r *PostRepository) IncrementViewCount(ctx context.Context, id string) error {
	const q = `UPDATE posts SET view_count = view_count + 1 WHERE post_id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

// ============================================================================
// CommentRepository
// ============================================================================

// CommentRepository는 PostgreSQL 기반 댓글 저장소입니다.
type CommentRepository struct {
	pool *pgxpool.Pool
}

// NewCommentRepository는 CommentRepository를 생성합니다.
func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

// Save는 댓글을 저장합니다.
func (r *CommentRepository) Save(ctx context.Context, comment *service.Comment) error {
	var parentID *string
	if comment.ParentCommentID != "" {
		parentID = &comment.ParentCommentID
	}
	const q = `INSERT INTO comments (comment_id, post_id, author_user_id, author_display_name, content, parent_comment_id, like_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q,
		comment.ID, comment.PostID, comment.AuthorUserID, comment.AuthorDisplayName,
		comment.Content, parentID, comment.LikeCount, comment.CreatedAt,
	)
	return err
}

// FindByPostID는 게시글의 댓글 목록을 반환합니다.
func (r *CommentRepository) FindByPostID(ctx context.Context, postID string, limit, offset int32) ([]*service.Comment, int32, error) {
	var total int32
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM comments WHERE post_id = $1", postID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT comment_id, post_id, author_user_id, COALESCE(author_display_name,''), content,
		COALESCE(parent_comment_id,''), like_count, created_at
		FROM comments WHERE post_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*service.Comment
	for rows.Next() {
		var c service.Comment
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.AuthorUserID, &c.AuthorDisplayName,
			&c.Content, &c.ParentCommentID, &c.LikeCount, &c.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		comments = append(comments, &c)
	}
	return comments, total, rows.Err()
}

// ============================================================================
// ChallengeRepository
// ============================================================================

// ChallengeRepository는 PostgreSQL 기반 챌린지 저장소입니다.
type ChallengeRepository struct {
	pool *pgxpool.Pool
}

// NewChallengeRepository는 ChallengeRepository를 생성합니다.
func NewChallengeRepository(pool *pgxpool.Pool) *ChallengeRepository {
	return &ChallengeRepository{pool: pool}
}

func challengeStatusToString(s service.ChallengeStatus) string {
	switch s {
	case service.ChUpcoming:
		return "upcoming"
	case service.ChActive:
		return "active"
	case service.ChCompleted:
		return "completed"
	case service.ChCancelled:
		return "cancelled"
	default:
		return "upcoming"
	}
}

func stringToChallengeStatus(s string) service.ChallengeStatus {
	switch s {
	case "upcoming":
		return service.ChUpcoming
	case "active":
		return service.ChActive
	case "completed":
		return service.ChCompleted
	case "cancelled":
		return service.ChCancelled
	default:
		return service.ChStatusUnknown
	}
}

func challengeTypeToString(t service.ChallengeType) string {
	switch t {
	case service.ChTypeSteps:
		return "exercise" // steps → exercise
	case service.ChTypeMeasurement:
		return "measurement"
	case service.ChTypeDiet:
		return "nutrition"
	case service.ChTypeExercise:
		return "exercise"
	case service.ChTypeSleep:
		return "sleep"
	default:
		return "custom"
	}
}

func stringToChallengeType(s string) service.ChallengeType {
	switch s {
	case "measurement":
		return service.ChTypeMeasurement
	case "exercise":
		return service.ChTypeExercise
	case "nutrition":
		return service.ChTypeDiet
	case "sleep":
		return service.ChTypeSleep
	default:
		return service.ChTypeUnknown
	}
}

// Save는 챌린지를 저장합니다.
func (r *ChallengeRepository) Save(ctx context.Context, ch *service.Challenge) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const q = `INSERT INTO challenges (challenge_id, creator_user_id, type, title, description, goal_description,
		target_value, target_unit, status, participant_count, max_participants, duration_days, start_date, end_date, created_at)
		VALUES ($1, $2, $3::challenge_type, $4, $5, $6, $7, $8, $9::challenge_status, $10, $11, $12, $13, $14, $15)`
	_, err = tx.Exec(ctx, q,
		ch.ID, ch.CreatorUserID, challengeTypeToString(ch.Type), ch.Title, ch.Description, ch.GoalDescription,
		ch.TargetValue, ch.TargetUnit, challengeStatusToString(ch.Status),
		ch.ParticipantCount, ch.MaxParticipants, ch.DurationDays, ch.StartDate, ch.EndDate, ch.CreatedAt,
	)
	if err != nil {
		return err
	}

	// 참가자 저장
	for userID := range ch.Participants {
		_, err = tx.Exec(ctx,
			`INSERT INTO challenge_participants (challenge_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			ch.ID, userID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindByID는 챌린지를 ID로 조회합니다.
func (r *ChallengeRepository) FindByID(ctx context.Context, id string) (*service.Challenge, error) {
	const q = `SELECT challenge_id, creator_user_id, type, title, COALESCE(description,''), COALESCE(goal_description,''),
		COALESCE(target_value,0), COALESCE(target_unit,''), status, participant_count, max_participants, duration_days,
		start_date, end_date, created_at
		FROM challenges WHERE challenge_id = $1`
	var ch service.Challenge
	var typeStr, statusStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&ch.ID, &ch.CreatorUserID, &typeStr, &ch.Title, &ch.Description, &ch.GoalDescription,
		&ch.TargetValue, &ch.TargetUnit, &statusStr,
		&ch.ParticipantCount, &ch.MaxParticipants, &ch.DurationDays,
		&ch.StartDate, &ch.EndDate, &ch.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	ch.Type = stringToChallengeType(typeStr)
	ch.Status = stringToChallengeStatus(statusStr)

	// 참가자 조회
	ch.Participants = make(map[string]bool)
	rows, err := r.pool.Query(ctx, "SELECT user_id FROM challenge_participants WHERE challenge_id = $1", id)
	if err != nil {
		return &ch, nil // 참가자 조회 실패는 무시
	}
	defer rows.Close()
	for rows.Next() {
		var uid string
		if rows.Scan(&uid) == nil {
			ch.Participants[uid] = true
		}
	}

	return &ch, nil
}

// FindAll는 챌린지 목록을 반환합니다.
func (r *ChallengeRepository) FindAll(ctx context.Context, status service.ChallengeStatus, limit, offset int32) ([]*service.Challenge, int32, error) {
	countQ := "SELECT COUNT(*) FROM challenges"
	listQ := `SELECT challenge_id, creator_user_id, type, title, COALESCE(description,''), COALESCE(goal_description,''),
		COALESCE(target_value,0), COALESCE(target_unit,''), status, participant_count, max_participants, duration_days,
		start_date, end_date, created_at FROM challenges`

	var args []interface{}
	idx := 1

	if status != service.ChStatusUnknown {
		where := " WHERE status = $1::challenge_status"
		countQ += where
		listQ += where
		args = append(args, challengeStatusToString(status))
		idx++
	}

	var total int32
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQ += " ORDER BY created_at DESC LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	queryArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, listQ, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var challenges []*service.Challenge
	for rows.Next() {
		var ch service.Challenge
		var typeStr, statusStr string
		if err := rows.Scan(
			&ch.ID, &ch.CreatorUserID, &typeStr, &ch.Title, &ch.Description, &ch.GoalDescription,
			&ch.TargetValue, &ch.TargetUnit, &statusStr,
			&ch.ParticipantCount, &ch.MaxParticipants, &ch.DurationDays,
			&ch.StartDate, &ch.EndDate, &ch.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		ch.Type = stringToChallengeType(typeStr)
		ch.Status = stringToChallengeStatus(statusStr)
		ch.Participants = make(map[string]bool)
		challenges = append(challenges, &ch)
	}
	return challenges, total, rows.Err()
}

// Update는 챌린지를 업데이트합니다.
func (r *ChallengeRepository) Update(ctx context.Context, ch *service.Challenge) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const q = `UPDATE challenges SET status = $1::challenge_status, participant_count = $2, description = $3,
		goal_description = $4, target_value = $5, target_unit = $6 WHERE challenge_id = $7`
	_, err = tx.Exec(ctx, q,
		challengeStatusToString(ch.Status), ch.ParticipantCount,
		ch.Description, ch.GoalDescription, ch.TargetValue, ch.TargetUnit, ch.ID,
	)
	if err != nil {
		return err
	}

	// 참가자 UPSERT
	for userID := range ch.Participants {
		_, err = tx.Exec(ctx,
			`INSERT INTO challenge_participants (challenge_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			ch.ID, userID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ============================================================================
// 유틸
// ============================================================================

func itoa(n int) string {
	return strconv.Itoa(n)
}
