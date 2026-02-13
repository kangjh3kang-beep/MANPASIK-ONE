// Package milvus는 Milvus 기반 VectorRepository 구현을 제공합니다.
package milvus

import (
	"context"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/vectordb"
)

// VectorRepository는 Milvus를 사용하는 벡터 저장소입니다.
type VectorRepository struct {
	client *vectordb.MilvusClient
}

// NewVectorRepository는 Milvus 기반 VectorRepository를 생성합니다.
func NewVectorRepository(client *vectordb.MilvusClient) *VectorRepository {
	return &VectorRepository{client: client}
}

// StoreFingerprint는 핑거프린트 벡터를 Milvus에 저장합니다.
func (r *VectorRepository) StoreFingerprint(ctx context.Context, sessionID string, vector []float32) error {
	id := uuid.New().String()
	return r.client.Insert(ctx, id, sessionID, vector)
}

// SearchSimilar는 유사한 핑거프린트를 검색합니다.
func (r *VectorRepository) SearchSimilar(ctx context.Context, vector []float32, topK int) ([]service.SimilarResult, error) {
	results, err := r.client.Search(ctx, vector, topK)
	if err != nil {
		return nil, err
	}

	similar := make([]service.SimilarResult, len(results))
	for i, res := range results {
		similar[i] = service.SimilarResult{
			SessionID: res.SessionID,
			Score:     res.Score,
			Distance:  1.0 - res.Score, // cosine distance = 1 - cosine similarity
		}
	}
	return similar, nil
}
