package memory

import (
	"context"

	"github.com/manpasik/backend/services/community-service/internal/service"
)

// PostSearchIndexer는 인메모리(no-op) 게시글 검색 인덱서입니다.
type PostSearchIndexer struct{}

// NewPostSearchIndexer는 no-op PostSearchIndexer를 생성합니다.
func NewPostSearchIndexer() *PostSearchIndexer {
	return &PostSearchIndexer{}
}

// IndexPost는 no-op 입니다.
func (s *PostSearchIndexer) IndexPost(_ context.Context, _ *service.Post) error {
	return nil
}

// DeletePost는 no-op 입니다.
func (s *PostSearchIndexer) DeletePost(_ context.Context, _ string) error {
	return nil
}
