// Package memory는 인메모리 NLP 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

// NLPRepository는 인메모리 NLP 저장소입니다.
type NLPRepository struct {
	mu          sync.RWMutex
	queries     map[string]*service.HealthQuery        // key: query ID
	extractions map[string]*service.SymptomExtraction   // key: extraction ID
	suggestions map[string][]service.Suggestion          // key: query ID
}

// NewNLPRepository는 인메모리 NLPRepository를 생성합니다.
func NewNLPRepository() *NLPRepository {
	return &NLPRepository{
		queries:     make(map[string]*service.HealthQuery),
		extractions: make(map[string]*service.SymptomExtraction),
		suggestions: make(map[string][]service.Suggestion),
	}
}

// SaveQuery는 건강 질의를 저장합니다.
func (r *NLPRepository) SaveQuery(_ context.Context, query *service.HealthQuery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *query
	// Entities 슬라이스 깊은 복사
	if query.Entities != nil {
		cp.Entities = make([]string, len(query.Entities))
		copy(cp.Entities, query.Entities)
	}
	r.queries[query.ID] = &cp
	return nil
}

// GetQuery는 질의 ID로 건강 질의를 조회합니다.
func (r *NLPRepository) GetQuery(_ context.Context, queryID string) (*service.HealthQuery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.queries[queryID]
	if !ok {
		return nil, nil
	}
	cp := *q
	if q.Entities != nil {
		cp.Entities = make([]string, len(q.Entities))
		copy(cp.Entities, q.Entities)
	}
	return &cp, nil
}

// SaveExtraction은 증상 추출 결과를 저장합니다.
func (r *NLPRepository) SaveExtraction(_ context.Context, extraction *service.SymptomExtraction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *extraction
	// Symptoms 슬라이스 깊은 복사
	if extraction.Symptoms != nil {
		cp.Symptoms = make([]service.Symptom, len(extraction.Symptoms))
		copy(cp.Symptoms, extraction.Symptoms)
	}
	r.extractions[extraction.ID] = &cp
	return nil
}

// GetSuggestions는 질의 ID에 해당하는 제안 목록을 반환합니다.
func (r *NLPRepository) GetSuggestions(_ context.Context, queryID string) ([]service.Suggestion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	suggestions, ok := r.suggestions[queryID]
	if !ok {
		return nil, nil
	}
	// 깊은 복사
	result := make([]service.Suggestion, len(suggestions))
	copy(result, suggestions)
	return result, nil
}
