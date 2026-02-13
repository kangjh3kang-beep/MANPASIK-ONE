// Package elasticsearch는 measurement-service의 Elasticsearch 검색 인덱서 구현입니다.
//
// 인덱스: measurements
// 공유 모듈: backend/shared/search/elasticsearch.go (ESClient)
package elasticsearch

import (
	"context"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/search"
)

const indexName = "measurements"

// SearchIndexer는 Elasticsearch 기반 검색 인덱서입니다.
type SearchIndexer struct {
	es *search.ESClient
}

// NewSearchIndexer는 SearchIndexer를 생성합니다.
func NewSearchIndexer(es *search.ESClient) *SearchIndexer {
	return &SearchIndexer{es: es}
}

// EnsureIndex는 measurements 인덱스를 생성합니다 (이미 존재하면 무시).
func (s *SearchIndexer) EnsureIndex(ctx context.Context) error {
	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"session_id":     map[string]string{"type": "keyword"},
				"device_id":     map[string]string{"type": "keyword"},
				"user_id":       map[string]string{"type": "keyword"},
				"cartridge_type": map[string]string{"type": "keyword"},
				"primary_value":  map[string]string{"type": "float"},
				"unit":           map[string]string{"type": "keyword"},
				"s_corrected":    map[string]string{"type": "float"},
				"confidence":     map[string]string{"type": "float"},
				"measured_at":    map[string]string{"type": "date"},
			},
		},
	}
	return s.es.CreateIndex(ctx, indexName, mappings)
}

// measurementDoc는 ES에 인덱싱할 문서 구조입니다.
type measurementDoc struct {
	SessionID     string  `json:"session_id"`
	DeviceID      string  `json:"device_id"`
	UserID        string  `json:"user_id"`
	CartridgeType string  `json:"cartridge_type"`
	PrimaryValue  float64 `json:"primary_value"`
	Unit          string  `json:"unit"`
	SCorrected    float64 `json:"s_corrected"`
	Confidence    float64 `json:"confidence"`
	MeasuredAt    string  `json:"measured_at"`
}

// IndexMeasurement는 측정 데이터를 Elasticsearch에 인덱싱합니다.
func (s *SearchIndexer) IndexMeasurement(ctx context.Context, sessionID string, data *service.MeasurementData) error {
	doc := measurementDoc{
		SessionID:     data.SessionID,
		DeviceID:      data.DeviceID,
		UserID:        data.UserID,
		CartridgeType: data.CartridgeType,
		PrimaryValue:  data.PrimaryValue,
		Unit:          data.Unit,
		SCorrected:    data.SCorrected,
		Confidence:    data.Confidence,
		MeasuredAt:    data.Time.Format("2006-01-02T15:04:05Z"),
	}
	docID := sessionID + "_" + data.Time.Format("20060102150405")
	return s.es.IndexDocument(ctx, indexName, docID, doc)
}
