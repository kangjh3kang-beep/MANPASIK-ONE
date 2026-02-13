package memory

import (
	"context"
	"math"
	"sync"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// VectorRepository는 인메모리 벡터 저장소입니다 (개발용, 실제는 Milvus).
type VectorRepository struct {
	mu      sync.RWMutex
	vectors map[string][]float32 // key: sessionID
}

// NewVectorRepository는 인메모리 VectorRepository를 생성합니다.
func NewVectorRepository() *VectorRepository {
	return &VectorRepository{
		vectors: make(map[string][]float32),
	}
}

func (r *VectorRepository) StoreFingerprint(_ context.Context, sessionID string, vector []float32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vectors[sessionID] = vector
	return nil
}

func (r *VectorRepository) SearchSimilar(_ context.Context, vector []float32, topK int) ([]service.SimilarResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type scored struct {
		id    string
		score float32
	}
	var results []scored
	for sid, v := range r.vectors {
		score := cosineSimilarity(vector, v)
		results = append(results, scored{id: sid, score: score})
	}

	// 간단한 정렬 (TopK)
	for i := 0; i < len(results) && i < topK; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	var out []service.SimilarResult
	for i := 0; i < len(results) && i < topK; i++ {
		out = append(out, service.SimilarResult{
			SessionID: results[i].id,
			Score:     results[i].score,
			Distance:  1 - results[i].score,
		})
	}
	return out, nil
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return float32(dot / denom)
}
