package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/manpasik/backend/services/vision-service/internal/service"
)

// FoodAnalysisRepository는 인메모리 음식 분석 저장소입니다.
type FoodAnalysisRepository struct {
	mu    sync.RWMutex
	store map[string]*service.FoodAnalysis
}

// NewFoodAnalysisRepository는 인메모리 FoodAnalysisRepository를 생성합니다.
func NewFoodAnalysisRepository() *FoodAnalysisRepository {
	return &FoodAnalysisRepository{
		store: make(map[string]*service.FoodAnalysis),
	}
}

func (r *FoodAnalysisRepository) Save(_ context.Context, analysis *service.FoodAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[analysis.ID] = analysis
	return nil
}

func (r *FoodAnalysisRepository) FindByID(_ context.Context, id string) (*service.FoodAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.store[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (r *FoodAnalysisRepository) FindByUserID(_ context.Context, userID string, limit, offset int32) ([]*service.FoodAnalysis, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userAnalyses []*service.FoodAnalysis
	for _, a := range r.store {
		if a.UserID == userID {
			userAnalyses = append(userAnalyses, a)
		}
	}

	// 최신 순 정렬
	sort.Slice(userAnalyses, func(i, j int) bool {
		return userAnalyses[i].CreatedAt.After(userAnalyses[j].CreatedAt)
	})

	total := int32(len(userAnalyses))

	// 페이지네이션
	start := int(offset)
	if start >= len(userAnalyses) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(userAnalyses) {
		end = len(userAnalyses)
	}

	return userAnalyses[start:end], total, nil
}

func (r *FoodAnalysisRepository) Update(_ context.Context, analysis *service.FoodAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[analysis.ID]; !ok {
		return nil
	}
	r.store[analysis.ID] = analysis
	return nil
}
