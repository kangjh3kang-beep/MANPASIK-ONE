// Package memory는 인메모리 커뮤니티 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/community-service/internal/service"
)

// ============================================================================
// PostRepository — 인메모리 게시글 저장소
// ============================================================================

// PostRepository는 인메모리 게시글 저장소입니다.
type PostRepository struct {
	mu    sync.RWMutex
	posts map[string]*service.Post
	order []string // 삽입 순서 유지
}

// NewPostRepository는 인메모리 PostRepository를 생성합니다.
func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make(map[string]*service.Post),
		order: make([]string, 0),
	}
}

// Save는 게시글을 저장합니다.
func (r *PostRepository) Save(_ context.Context, post *service.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := copyPost(post)
	r.posts[post.ID] = cp
	r.order = append(r.order, post.ID)
	return nil
}

// FindByID는 게시글을 ID로 조회합니다.
func (r *PostRepository) FindByID(_ context.Context, id string) (*service.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.posts[id]
	if !ok {
		return nil, nil
	}
	cp := copyPost(p)
	return cp, nil
}

// FindAll는 게시글 목록을 반환합니다. category가 0이면 전체를 반환합니다.
func (r *PostRepository) FindAll(_ context.Context, category service.PostCategory, limit, offset int32) ([]*service.Post, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Post
	// 최신순 (역순)
	for i := len(r.order) - 1; i >= 0; i-- {
		p := r.posts[r.order[i]]
		if category != service.CatUnknown && p.Category != category {
			continue
		}
		filtered = append(filtered, p)
	}

	total := int32(len(filtered))

	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	result := make([]*service.Post, 0, end-start)
	for _, p := range filtered[start:end] {
		result = append(result, copyPost(p))
	}

	return result, total, nil
}

// IncrementLikeCount는 게시글의 좋아요 수를 1 증가시킵니다.
func (r *PostRepository) IncrementLikeCount(_ context.Context, id string) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.posts[id]
	if !ok {
		return 0, nil
	}
	p.LikeCount++
	return int32(p.LikeCount), nil
}

// IncrementCommentCount는 게시글의 댓글 수를 1 증가시킵니다.
func (r *PostRepository) IncrementCommentCount(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.posts[id]; ok {
		p.CommentCount++
	}
	return nil
}

// IncrementViewCount는 게시글의 조회수를 1 증가시킵니다.
func (r *PostRepository) IncrementViewCount(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.posts[id]; ok {
		p.ViewCount++
	}
	return nil
}

// copyPost는 Post의 깊은 복사를 수행합니다.
func copyPost(src *service.Post) *service.Post {
	cp := *src
	if src.Tags != nil {
		cp.Tags = make([]string, len(src.Tags))
		copy(cp.Tags, src.Tags)
	}
	return &cp
}

// ============================================================================
// CommentRepository — 인메모리 댓글 저장소
// ============================================================================

// CommentRepository는 인메모리 댓글 저장소입니다.
type CommentRepository struct {
	mu       sync.RWMutex
	comments []*service.Comment
}

// NewCommentRepository는 인메모리 CommentRepository를 생성합니다.
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make([]*service.Comment, 0),
	}
}

// Save는 댓글을 저장합니다.
func (r *CommentRepository) Save(_ context.Context, comment *service.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *comment
	r.comments = append(r.comments, &cp)
	return nil
}

// FindByPostID는 게시글의 댓글 목록을 반환합니다.
func (r *CommentRepository) FindByPostID(_ context.Context, postID string, limit, offset int32) ([]*service.Comment, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Comment
	for _, c := range r.comments {
		if c.PostID == postID {
			filtered = append(filtered, c)
		}
	}

	total := int32(len(filtered))

	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	result := make([]*service.Comment, 0, end-start)
	for _, c := range filtered[start:end] {
		cp := *c
		result = append(result, &cp)
	}

	return result, total, nil
}

// ============================================================================
// ChallengeRepository — 인메모리 챌린지 저장소
// ============================================================================

// ChallengeRepository는 인메모리 챌린지 저장소입니다.
type ChallengeRepository struct {
	mu         sync.RWMutex
	challenges map[string]*service.Challenge
	order      []string
}

// NewChallengeRepository는 인메모리 ChallengeRepository를 생성합니다.
func NewChallengeRepository() *ChallengeRepository {
	return &ChallengeRepository{
		challenges: make(map[string]*service.Challenge),
		order:      make([]string, 0),
	}
}

// Save는 챌린지를 저장합니다.
func (r *ChallengeRepository) Save(_ context.Context, challenge *service.Challenge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := copyChallenge(challenge)
	r.challenges[challenge.ID] = cp
	r.order = append(r.order, challenge.ID)
	return nil
}

// FindByID는 챌린지를 ID로 조회합니다.
func (r *ChallengeRepository) FindByID(_ context.Context, id string) (*service.Challenge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.challenges[id]
	if !ok {
		return nil, nil
	}
	cp := copyChallenge(ch)
	return cp, nil
}

// FindAll는 챌린지 목록을 반환합니다. status가 0이면 전체를 반환합니다.
func (r *ChallengeRepository) FindAll(_ context.Context, status service.ChallengeStatus, limit, offset int32) ([]*service.Challenge, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Challenge
	for i := len(r.order) - 1; i >= 0; i-- {
		ch := r.challenges[r.order[i]]
		if status != service.ChStatusUnknown && ch.Status != status {
			continue
		}
		filtered = append(filtered, ch)
	}

	total := int32(len(filtered))

	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	result := make([]*service.Challenge, 0, end-start)
	for _, ch := range filtered[start:end] {
		result = append(result, copyChallenge(ch))
	}

	return result, total, nil
}

// Update는 챌린지를 업데이트합니다.
func (r *ChallengeRepository) Update(_ context.Context, challenge *service.Challenge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := copyChallenge(challenge)
	r.challenges[challenge.ID] = cp
	return nil
}

// copyChallenge는 Challenge의 깊은 복사를 수행합니다.
func copyChallenge(src *service.Challenge) *service.Challenge {
	cp := *src
	if src.Participants != nil {
		cp.Participants = make(map[string]bool, len(src.Participants))
		for k, v := range src.Participants {
			cp.Participants[k] = v
		}
	}
	return &cp
}
