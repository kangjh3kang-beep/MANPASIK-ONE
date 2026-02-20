// Package service는 nlp-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// HealthQuery는 사용자의 건강 관련 자연어 질의를 파싱한 결과입니다.
type HealthQuery struct {
	ID         string
	UserID     string
	RawText    string
	Intent     string
	Entities   []string
	Confidence float64
	CreatedAt  time.Time
}

// SymptomExtraction은 텍스트에서 추출한 증상 정보입니다.
type SymptomExtraction struct {
	ID          string
	Text        string
	Symptoms    []Symptom
	ProcessedAt time.Time
}

// Symptom은 개별 증상 정보입니다.
type Symptom struct {
	Name       string
	Severity   string
	BodyPart   string
	Confidence float64
}

// Suggestion은 건강 질의에 대한 제안 항목입니다.
type Suggestion struct {
	ID       string
	QueryID  string
	Text     string
	Category string
	Priority int
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// NLPRepository는 NLP 데이터 저장소 인터페이스입니다.
type NLPRepository interface {
	SaveQuery(ctx context.Context, query *HealthQuery) error
	GetQuery(ctx context.Context, queryID string) (*HealthQuery, error)
	SaveExtraction(ctx context.Context, extraction *SymptomExtraction) error
	GetSuggestions(ctx context.Context, queryID string) ([]Suggestion, error)
}

// ============================================================================
// NLPService
// ============================================================================

// NLPService는 NLP 비즈니스 로직입니다.
type NLPService struct {
	repo NLPRepository
}

// NewNLPService는 새 NLPService를 생성합니다.
func NewNLPService(repo NLPRepository) *NLPService {
	return &NLPService{repo: repo}
}

// knownKeywords는 건강 관련 키워드 목록입니다.
var knownKeywords = []string{
	"blood sugar",
	"headache",
	"blood pressure",
	"fever",
	"pain",
	"cholesterol",
	"heart rate",
	"dizziness",
	"nausea",
	"fatigue",
}

// symptomDB는 키워드 기반 증상 매핑 테이블입니다.
var symptomDB = map[string]Symptom{
	"headache": {
		Name:       "headache",
		Severity:   "moderate",
		BodyPart:   "head",
		Confidence: 0.90,
	},
	"fever": {
		Name:       "fever",
		Severity:   "moderate",
		BodyPart:   "systemic",
		Confidence: 0.92,
	},
	"pain": {
		Name:       "pain",
		Severity:   "moderate",
		BodyPart:   "general",
		Confidence: 0.75,
	},
	"dizziness": {
		Name:       "dizziness",
		Severity:   "mild",
		BodyPart:   "head",
		Confidence: 0.85,
	},
	"nausea": {
		Name:       "nausea",
		Severity:   "mild",
		BodyPart:   "abdomen",
		Confidence: 0.88,
	},
	"fatigue": {
		Name:       "fatigue",
		Severity:   "mild",
		BodyPart:   "systemic",
		Confidence: 0.80,
	},
	"chest pain": {
		Name:       "chest pain",
		Severity:   "severe",
		BodyPart:   "chest",
		Confidence: 0.95,
	},
	"cough": {
		Name:       "cough",
		Severity:   "mild",
		BodyPart:   "respiratory",
		Confidence: 0.87,
	},
	"sore throat": {
		Name:       "sore throat",
		Severity:   "mild",
		BodyPart:   "throat",
		Confidence: 0.89,
	},
	"back pain": {
		Name:       "back pain",
		Severity:   "moderate",
		BodyPart:   "back",
		Confidence: 0.86,
	},
}

// ParseHealthQuery는 사용자의 건강 관련 텍스트를 파싱하여 의도와 엔티티를 추출합니다.
// 시뮬레이션: 텍스트에서 알려진 키워드를 검색하고 intent를 "health_inquiry"로 설정합니다.
func (s *NLPService) ParseHealthQuery(ctx context.Context, userID, text string) (*HealthQuery, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("text is required")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user_id is required")
	}

	lower := strings.ToLower(text)

	// 키워드 매칭으로 엔티티 추출
	var entities []string
	for _, kw := range knownKeywords {
		if strings.Contains(lower, kw) {
			entities = append(entities, kw)
		}
	}

	// 신뢰도 계산: 매칭된 키워드가 많을수록 높아짐
	confidence := 0.5
	if len(entities) > 0 {
		confidence = 0.7 + float64(len(entities))*0.05
		if confidence > 0.99 {
			confidence = 0.99
		}
	}

	query := &HealthQuery{
		ID:         uuid.New().String(),
		UserID:     userID,
		RawText:    text,
		Intent:     "health_inquiry",
		Entities:   entities,
		Confidence: confidence,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.SaveQuery(ctx, query); err != nil {
		return nil, err
	}

	return query, nil
}

// ExtractSymptoms는 텍스트에서 증상 키워드를 추출합니다.
// 시뮬레이션: symptomDB의 키워드와 단순 문자열 매칭을 수행합니다.
func (s *NLPService) ExtractSymptoms(ctx context.Context, text string) (*SymptomExtraction, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("text is required")
	}

	lower := strings.ToLower(text)

	var symptoms []Symptom
	for keyword, symptom := range symptomDB {
		if strings.Contains(lower, keyword) {
			symptoms = append(symptoms, symptom)
		}
	}

	extraction := &SymptomExtraction{
		ID:          uuid.New().String(),
		Text:        text,
		Symptoms:    symptoms,
		ProcessedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveExtraction(ctx, extraction); err != nil {
		return nil, err
	}

	return extraction, nil
}

// GetSuggestions는 지정된 질의에 대한 저장된 제안 목록을 반환합니다.
func (s *NLPService) GetSuggestions(ctx context.Context, queryID string) ([]Suggestion, error) {
	if strings.TrimSpace(queryID) == "" {
		return nil, errors.New("query_id is required")
	}

	return s.repo.GetSuggestions(ctx, queryID)
}
