package memory

import (
	"context"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// SearchIndexer는 인메모리(no-op) 검색 인덱서입니다.
type SearchIndexer struct{}

// NewSearchIndexer는 no-op SearchIndexer를 생성합니다.
func NewSearchIndexer() *SearchIndexer {
	return &SearchIndexer{}
}

// IndexMeasurement는 no-op 입니다.
func (s *SearchIndexer) IndexMeasurement(_ context.Context, _ string, _ *service.MeasurementData) error {
	return nil
}
