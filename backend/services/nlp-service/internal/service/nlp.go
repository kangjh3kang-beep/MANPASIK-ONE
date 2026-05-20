// Package service는 nlp-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/nlp-service/internal/classifier"
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
	Urgency    string // "routine", "urgent", "emergency"
	CreatedAt  time.Time
}

// SymptomExtraction은 텍스트에서 추출한 증상 정보입니다.
type SymptomExtraction struct {
	ID          string
	Text        string
	Symptoms    []Symptom
	Urgency     string
	ProcessedAt time.Time
}

// Symptom은 개별 증상 정보입니다.
type Symptom struct {
	Name       string
	Severity   string // "mild", "moderate", "severe"
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
// 긴급도 상수
// ============================================================================

const (
	UrgencyRoutine   = "routine"
	UrgencyUrgent    = "urgent"
	UrgencyEmergency = "emergency"
)

// ============================================================================
// 한국어 증상 DB
// ============================================================================

// koreanSymptomDB는 한국어 증상 키워드 매핑 테이블입니다.
var koreanSymptomDB = map[string]Symptom{
	"두통":   {Name: "두통", Severity: "moderate", BodyPart: "head", Confidence: 0.90},
	"발열":   {Name: "발열", Severity: "moderate", BodyPart: "systemic", Confidence: 0.92},
	"어지러움": {Name: "어지러움", Severity: "mild", BodyPart: "head", Confidence: 0.85},
	"구역질":  {Name: "구역질", Severity: "mild", BodyPart: "abdomen", Confidence: 0.88},
	"흉통":   {Name: "흉통", Severity: "severe", BodyPart: "chest", Confidence: 0.95},
	"호흡곤란": {Name: "호흡곤란", Severity: "severe", BodyPart: "respiratory", Confidence: 0.94},
	"복통":   {Name: "복통", Severity: "moderate", BodyPart: "abdomen", Confidence: 0.87},
	"허리통증": {Name: "허리통증", Severity: "moderate", BodyPart: "back", Confidence: 0.86},
	"피로감":  {Name: "피로감", Severity: "mild", BodyPart: "systemic", Confidence: 0.80},
	"기침":   {Name: "기침", Severity: "mild", BodyPart: "respiratory", Confidence: 0.87},
	"인후통":  {Name: "인후통", Severity: "mild", BodyPart: "throat", Confidence: 0.86},
	"근육통":  {Name: "근육통", Severity: "mild", BodyPart: "musculoskeletal", Confidence: 0.83},
	"관절통":  {Name: "관절통", Severity: "moderate", BodyPart: "joint", Confidence: 0.85},
	"설사":   {Name: "설사", Severity: "mild", BodyPart: "abdomen", Confidence: 0.88},
	"변비":   {Name: "변비", Severity: "mild", BodyPart: "abdomen", Confidence: 0.85},
}

// koreanIntentKeywords는 한국어 건강 인텐트 키워드 매핑입니다.
var koreanIntentKeywords = map[string]string{
	"혈당":    "blood_sugar",
	"혈압":    "blood_pressure",
	"콜레스테롤": "cholesterol",
	"심박수":   "heart_rate",
	"체온":    "body_temperature",
	"산소포화도": "oxygen_saturation",
	"체중":    "weight",
	"수면":    "sleep",
	"운동":    "exercise",
	"식사":    "nutrition",
	"측정":    "measurement",
}

// ============================================================================
// 영문 증상/키워드 DB
// ============================================================================

// knownKeywords는 영문 건강 관련 키워드 목록입니다.
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

// symptomDB는 영문 키워드 기반 증상 매핑 테이블입니다.
var symptomDB = map[string]Symptom{
	"headache":    {Name: "headache", Severity: "moderate", BodyPart: "head", Confidence: 0.90},
	"fever":       {Name: "fever", Severity: "moderate", BodyPart: "systemic", Confidence: 0.92},
	"pain":        {Name: "pain", Severity: "moderate", BodyPart: "general", Confidence: 0.75},
	"dizziness":   {Name: "dizziness", Severity: "mild", BodyPart: "head", Confidence: 0.85},
	"nausea":      {Name: "nausea", Severity: "mild", BodyPart: "abdomen", Confidence: 0.88},
	"fatigue":     {Name: "fatigue", Severity: "mild", BodyPart: "systemic", Confidence: 0.80},
	"chest pain":  {Name: "chest pain", Severity: "severe", BodyPart: "chest", Confidence: 0.95},
	"cough":       {Name: "cough", Severity: "mild", BodyPart: "respiratory", Confidence: 0.87},
	"sore throat": {Name: "sore throat", Severity: "mild", BodyPart: "throat", Confidence: 0.89},
	"back pain":   {Name: "back pain", Severity: "moderate", BodyPart: "back", Confidence: 0.86},
}

// ============================================================================
// 제안 생성 규칙
// ============================================================================

// intentSuggestions는 인텐트별 기본 제안 목록입니다.
var intentSuggestions = map[string][]Suggestion{
	"blood_sugar": {
		{Text: "혈당 측정을 시작하세요", Category: "measurement", Priority: 1},
		{Text: "최근 혈당 추세를 확인하세요", Category: "analytics", Priority: 2},
	},
	"blood_pressure": {
		{Text: "혈압 측정을 시작하세요", Category: "measurement", Priority: 1},
		{Text: "혈압 관리 코칭을 확인하세요", Category: "coaching", Priority: 2},
	},
	"cholesterol": {
		{Text: "콜레스테롤 수치를 확인하세요", Category: "measurement", Priority: 1},
	},
	"measurement": {
		{Text: "새 측정을 시작하세요", Category: "measurement", Priority: 1},
	},
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
	SaveSuggestions(ctx context.Context, queryID string, suggestions []Suggestion) error
}

// ============================================================================
// NLPService
// ============================================================================

// NLPService는 NLP 비즈니스 로직입니다.
type NLPService struct {
	repo       NLPRepository
	classifier classifier.IntentClassifier // optional: nil이면 키워드 기반 분류
	router     *IntentRouter
}

// NewNLPService는 새 NLPService를 생성합니다.
func NewNLPService(repo NLPRepository) *NLPService {
	return &NLPService{
		repo:   repo,
		router: NewIntentRouter(),
	}
}

// AnalyzeKorean은 한글 텍스트를 종합 분석합니다 (토큰화 + 의도 분류 + 엔티티 추출).
func (s *NLPService) AnalyzeKorean(text string) *KoreanAnalysis {
	intentResult := s.router.Classify(text)
	entities := s.router.ExtractEntities(text)

	return &KoreanAnalysis{
		Text:       text,
		Intent:     intentResult.Intent,
		Confidence: intentResult.Confidence,
		Severity:   intentResult.Severity,
		Keywords:   intentResult.Keywords,
		Tokens:     intentResult.Tokens,
		Entities:   entities,
		AnalyzedAt: time.Now().UTC(),
	}
}

// KoreanAnalysis는 한글 텍스트 종합 분석 결과입니다.
type KoreanAnalysis struct {
	Text       string
	Intent     string
	Confidence float64
	Severity   int
	Keywords   []string
	Tokens     []Token
	Entities   *EntityExtraction
	AnalyzedAt time.Time
}

// SeverityToUrgency는 긴급도 점수를 urgency 문자열로 변환합니다.
func SeverityToUrgency(severity int) string {
	if severity >= 70 {
		return UrgencyEmergency
	}
	if severity >= 35 {
		return UrgencyUrgent
	}
	return UrgencyRoutine
}

// SetIntentClassifier는 외부 인텐트 분류기를 설정합니다 (optional).
func (s *NLPService) SetIntentClassifier(c classifier.IntentClassifier) {
	s.classifier = c
}

// DetectUrgency는 증상 목록으로 긴급도를 판단합니다.
func DetectUrgency(symptoms []Symptom) string {
	moderateCount := 0
	for _, s := range symptoms {
		if s.Severity == "severe" {
			return UrgencyEmergency
		}
		if s.Severity == "moderate" {
			moderateCount++
		}
	}
	if moderateCount >= 2 {
		return UrgencyUrgent
	}
	return UrgencyRoutine
}

// ParseHealthQuery는 사용자의 건강 관련 텍스트를 파싱하여 의도와 엔티티를 추출합니다.
// 외부 분류기가 설정되어 있으면 우선 사용하고, 실패 시 키워드 기반 폴백합니다.
func (s *NLPService) ParseHealthQuery(ctx context.Context, userID, text string) (*HealthQuery, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("text is required")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user_id is required")
	}

	// 외부 분류기 우선 시도
	if s.classifier != nil {
		result, err := s.classifier.Classify(ctx, text)
		if err != nil {
			log.Printf("[NLP] 외부 분류기 실패, 키워드 폴백 사용: %v", err)
		} else if result != nil {
			query := &HealthQuery{
				ID:         uuid.New().String(),
				UserID:     userID,
				RawText:    text,
				Intent:     result.Intent,
				Entities:   result.Entities,
				Confidence: result.Confidence,
				Urgency:    result.Urgency,
				CreatedAt:  time.Now().UTC(),
			}
			if err := s.repo.SaveQuery(ctx, query); err != nil {
				return nil, err
			}
			return query, nil
		}
	}

	// 키워드 기반 폴백
	return s.parseHealthQueryKeyword(ctx, userID, text)
}

// parseHealthQueryKeyword는 키워드 기반으로 건강 질의를 파싱합니다.
// IntentRouter를 우선 활용하여 한글 의도/엔티티/긴급도를 분석합니다.
func (s *NLPService) parseHealthQueryKeyword(ctx context.Context, userID, text string) (*HealthQuery, error) {
	lower := strings.ToLower(text)

	// IntentRouter로 한글 종합 분석
	analysis := s.router.Classify(text)
	korEntities := s.router.ExtractEntities(text)

	// 영문 키워드 매칭으로 엔티티 추출
	var entities []string
	for _, kw := range knownKeywords {
		if strings.Contains(lower, kw) {
			entities = append(entities, kw)
		}
	}

	// 한국어 인텐트 키워드 매칭
	intent := "health_inquiry"
	for keyword, intentName := range koreanIntentKeywords {
		if strings.Contains(text, keyword) {
			entities = append(entities, intentName)
			intent = fmt.Sprintf("health_inquiry.%s", intentName)
			break
		}
	}

	// IntentRouter 결과로 의도 보정 (긴급/예약/증상)
	switch analysis.Intent {
	case IntentEmergency:
		intent = "emergency"
	case IntentReservation:
		intent = "reservation"
	case IntentSymptom:
		if intent == "health_inquiry" {
			intent = "symptom_report"
		}
	}

	// 한국어 증상 키워드 매칭
	entities = append(entities, korEntities.Symptoms...)
	entities = append(entities, korEntities.BodyParts...)

	// 신뢰도 계산 (IntentRouter 우선)
	confidence := analysis.Confidence
	if len(entities) > 0 {
		confidence = (confidence + 0.7 + float64(len(entities))*0.03) / 2
		if confidence > 0.99 {
			confidence = 0.99
		}
	}

	// 긴급도: IntentRouter severity 우선, 폴백으로 증상 기반
	urgency := SeverityToUrgency(analysis.Severity)
	if urgency == UrgencyRoutine {
		symptoms := extractSymptomsFromText(text)
		urgency = DetectUrgency(symptoms)
	}

	query := &HealthQuery{
		ID:         uuid.New().String(),
		UserID:     userID,
		RawText:    text,
		Intent:     intent,
		Entities:   dedupStrings(entities),
		Confidence: confidence,
		Urgency:    urgency,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.SaveQuery(ctx, query); err != nil {
		return nil, err
	}

	return query, nil
}

// dedupStrings는 문자열 슬라이스에서 중복을 제거합니다.
func dedupStrings(in []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// ExtractSymptoms는 텍스트에서 증상 키워드를 추출합니다.
func (s *NLPService) ExtractSymptoms(ctx context.Context, text string) (*SymptomExtraction, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("text is required")
	}

	symptoms := extractSymptomsFromText(text)
	urgency := DetectUrgency(symptoms)

	extraction := &SymptomExtraction{
		ID:          uuid.New().String(),
		Text:        text,
		Symptoms:    symptoms,
		Urgency:     urgency,
		ProcessedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveExtraction(ctx, extraction); err != nil {
		return nil, err
	}

	return extraction, nil
}

// extractSymptomsFromText는 텍스트에서 영문+한국어 증상을 추출합니다.
func extractSymptomsFromText(text string) []Symptom {
	lower := strings.ToLower(text)
	var symptoms []Symptom
	seen := make(map[string]bool)

	// 영문 증상 매칭
	for keyword, symptom := range symptomDB {
		if strings.Contains(lower, keyword) && !seen[symptom.Name] {
			symptoms = append(symptoms, symptom)
			seen[symptom.Name] = true
		}
	}

	// 한국어 증상 매칭
	for keyword, symptom := range koreanSymptomDB {
		if strings.Contains(text, keyword) && !seen[symptom.Name] {
			symptoms = append(symptoms, symptom)
			seen[symptom.Name] = true
		}
	}

	return symptoms
}

// GetSuggestions는 지정된 질의에 대한 저장된 제안 목록을 반환합니다.
func (s *NLPService) GetSuggestions(ctx context.Context, queryID string) ([]Suggestion, error) {
	if strings.TrimSpace(queryID) == "" {
		return nil, errors.New("query_id is required")
	}

	return s.repo.GetSuggestions(ctx, queryID)
}

// GenerateSuggestions는 건강 질의 인텐트에 기반한 제안을 생성합니다.
func (s *NLPService) GenerateSuggestions(ctx context.Context, query *HealthQuery) ([]Suggestion, error) {
	if query == nil {
		return nil, errors.New("query is required")
	}

	var suggestions []Suggestion

	// 엔티티 기반 제안 생성
	for _, entity := range query.Entities {
		if sgs, ok := intentSuggestions[entity]; ok {
			for _, sg := range sgs {
				sg.ID = uuid.New().String()
				sg.QueryID = query.ID
				suggestions = append(suggestions, sg)
			}
		}
	}

	// 긴급 상황인 경우 응급 제안 추가
	if query.Urgency == UrgencyEmergency {
		suggestions = append(suggestions, Suggestion{
			ID:       uuid.New().String(),
			QueryID:  query.ID,
			Text:     "응급 상황이 감지되었습니다. 119에 전화하거나 가까운 응급실을 방문하세요.",
			Category: "emergency",
			Priority: 0,
		})
	}

	// 기본 제안 (제안이 없는 경우)
	if len(suggestions) == 0 {
		suggestions = append(suggestions, Suggestion{
			ID:       uuid.New().String(),
			QueryID:  query.ID,
			Text:     "건강 측정을 시작하여 현재 상태를 확인하세요",
			Category: "general",
			Priority: 3,
		})
	}

	// 저장
	if err := s.repo.SaveSuggestions(ctx, query.ID, suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}
