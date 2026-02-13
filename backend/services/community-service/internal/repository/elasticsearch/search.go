// Package elasticsearch는 community-service의 Elasticsearch 검색 인덱서 구현입니다.
//
// 인덱스: community_posts
// 공유 모듈: backend/shared/search/elasticsearch.go (ESClient)
package elasticsearch

import (
	"context"

	"github.com/manpasik/backend/services/community-service/internal/service"
	"github.com/manpasik/backend/shared/search"
)

const postIndexName = "community_posts"

// PostSearchIndexer는 Elasticsearch 기반 게시글 검색 인덱서입니다.
type PostSearchIndexer struct {
	es *search.ESClient
}

// NewPostSearchIndexer는 PostSearchIndexer를 생성합니다.
func NewPostSearchIndexer(es *search.ESClient) *PostSearchIndexer {
	return &PostSearchIndexer{es: es}
}

// EnsureIndex는 community_posts 인덱스를 생성합니다 (이미 존재하면 무시).
func (s *PostSearchIndexer) EnsureIndex(ctx context.Context) error {
	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"post_id":       map[string]string{"type": "keyword"},
				"author_user_id": map[string]string{"type": "keyword"},
				"category":       map[string]string{"type": "keyword"},
				"title":          map[string]string{"type": "text", "analyzer": "standard"},
				"content":        map[string]string{"type": "text", "analyzer": "standard"},
				"tags":           map[string]string{"type": "keyword"},
				"created_at":     map[string]string{"type": "date"},
			},
		},
	}
	return s.es.CreateIndex(ctx, postIndexName, mappings)
}

// postDoc는 ES에 인덱싱할 게시글 문서입니다.
type postDoc struct {
	PostID       string   `json:"post_id"`
	AuthorUserID string   `json:"author_user_id"`
	Category     string   `json:"category"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Tags         []string `json:"tags"`
	CreatedAt    string   `json:"created_at"`
}

func categoryString(c service.PostCategory) string {
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

// IndexPost는 게시글을 Elasticsearch에 인덱싱합니다.
func (s *PostSearchIndexer) IndexPost(ctx context.Context, post *service.Post) error {
	doc := postDoc{
		PostID:       post.ID,
		AuthorUserID: post.AuthorUserID,
		Category:     categoryString(post.Category),
		Title:        post.Title,
		Content:      post.Content,
		Tags:         post.Tags,
		CreatedAt:    post.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	return s.es.IndexDocument(ctx, postIndexName, post.ID, doc)
}

// DeletePost는 게시글을 Elasticsearch에서 삭제합니다.
func (s *PostSearchIndexer) DeletePost(ctx context.Context, postID string) error {
	return s.es.DeleteDocument(ctx, postIndexName, postID)
}
