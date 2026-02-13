// Package service는 translation-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// SupportedLanguage는 지원 언어 정보입니다.
type SupportedLanguage struct {
	Code            string
	Name            string
	NativeName      string
	SupportsMedical bool
	SupportsRealtime bool
}

// TranslationRecord는 번역 이력 레코드입니다.
type TranslationRecord struct {
	ID             string
	SourceText     string
	TranslatedText string
	SourceLanguage string
	TargetLanguage string
	Confidence     float32
	UserID         string
	CreatedAt      time.Time
}

// UsageStats는 번역 사용량 통계입니다.
type UsageStats struct {
	TotalCharacters   int64
	MonthlyCharacters int64
	MonthlyLimit      int64
	TotalRequests     int
	MonthlyRequests   int
}

// TranslationRepository는 번역 이력 저장소 인터페이스입니다.
type TranslationRepository interface {
	Save(ctx context.Context, record *TranslationRecord) error
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*TranslationRecord, int, error)
}

// UsageRepository는 사용량 저장소 인터페이스입니다.
type UsageRepository interface {
	IncrementUsage(ctx context.Context, userID string, characters int) error
	GetUsage(ctx context.Context, userID string) (*UsageStats, error)
}

// TranslationService는 번역 서비스 핵심 로직입니다.
type TranslationService struct {
	log             *zap.Logger
	translationRepo TranslationRepository
	usageRepo       UsageRepository
	languages       []SupportedLanguage
	medicalTerms    map[string]map[string]string // sourceCode -> targetCode -> translation map
}

// NewTranslationService는 TranslationService를 생성합니다.
func NewTranslationService(log *zap.Logger, translationRepo TranslationRepository, usageRepo UsageRepository) *TranslationService {
	return &TranslationService{
		log:             log,
		translationRepo: translationRepo,
		usageRepo:       usageRepo,
		languages:       defaultLanguages(),
		medicalTerms:    defaultMedicalTerms(),
	}
}

// TranslateText는 텍스트를 번역합니다.
func (s *TranslationService) TranslateText(ctx context.Context, text, sourceLang, targetLang, contextHint, userID string) (string, string, float32, string, error) {
	if text == "" {
		return "", "", 0, "", apperrors.New(apperrors.ErrInvalidInput, "text는 필수입니다")
	}
	if targetLang == "" {
		return "", "", 0, "", apperrors.New(apperrors.ErrInvalidInput, "target_language는 필수입니다")
	}

	// 소스 언어 자동 감지
	if sourceLang == "" {
		sourceLang = detectLanguageCode(text)
	}

	// 의료 용어 먼저 치환, 그 외에는 시뮬레이션 번역
	translated := s.translateWithMedicalTerms(text, sourceLang, targetLang)
	confidence := float32(0.92)
	if contextHint == "medical" {
		confidence = 0.95
	}

	translationID := uuid.New().String()
	record := &TranslationRecord{
		ID:             translationID,
		SourceText:     text,
		TranslatedText: translated,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		Confidence:     confidence,
		UserID:         userID,
		CreatedAt:      time.Now(),
	}
	s.translationRepo.Save(ctx, record)

	charCount := utf8.RuneCountInString(text)
	if userID != "" {
		s.usageRepo.IncrementUsage(ctx, userID, charCount)
	}

	s.log.Info("번역 완료",
		zap.String("translation_id", translationID),
		zap.String("source_lang", sourceLang),
		zap.String("target_lang", targetLang),
		zap.Int("characters", charCount),
	)

	return translated, sourceLang, confidence, translationID, nil
}

// DetectLanguage는 텍스트의 언어를 감지합니다.
func (s *TranslationService) DetectLanguage(_ context.Context, text string) ([]DetectedLanguage, string, error) {
	if text == "" {
		return nil, "", apperrors.New(apperrors.ErrInvalidInput, "text는 필수입니다")
	}

	primary := detectLanguageCode(text)

	results := []DetectedLanguage{
		{Code: primary, Name: languageName(primary), Confidence: 0.95},
	}

	// 2순위 후보
	if primary == "ko" {
		results = append(results, DetectedLanguage{Code: "ja", Name: "Japanese", Confidence: 0.03})
	} else if primary == "en" {
		results = append(results, DetectedLanguage{Code: "fr", Name: "French", Confidence: 0.03})
	}

	return results, primary, nil
}

// DetectedLanguage는 감지된 언어 정보입니다.
type DetectedLanguage struct {
	Code       string
	Name       string
	Confidence float32
}

// ListSupportedLanguages는 지원 언어 목록을 반환합니다.
func (s *TranslationService) ListSupportedLanguages(_ context.Context) ([]SupportedLanguage, error) {
	return s.languages, nil
}

// TranslateBatch는 여러 텍스트를 일괄 번역합니다.
func (s *TranslationService) TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang, userID string) ([]string, string, int, error) {
	if len(texts) == 0 {
		return nil, "", 0, apperrors.New(apperrors.ErrInvalidInput, "texts는 필수입니다")
	}
	if targetLang == "" {
		return nil, "", 0, apperrors.New(apperrors.ErrInvalidInput, "target_language는 필수입니다")
	}

	if sourceLang == "" && len(texts) > 0 {
		sourceLang = detectLanguageCode(texts[0])
	}

	totalChars := 0
	var translated []string
	for _, text := range texts {
		t := s.translateWithMedicalTerms(text, sourceLang, targetLang)
		translated = append(translated, t)
		totalChars += utf8.RuneCountInString(text)
	}

	if userID != "" {
		s.usageRepo.IncrementUsage(ctx, userID, totalChars)
	}

	return translated, sourceLang, totalChars, nil
}

// GetTranslationHistory는 번역 이력을 조회합니다.
func (s *TranslationService) GetTranslationHistory(ctx context.Context, userID string, limit, offset int) ([]*TranslationRecord, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.translationRepo.FindByUserID(ctx, userID, limit, offset)
}

// GetTranslationUsage는 번역 사용량을 조회합니다.
func (s *TranslationService) GetTranslationUsage(ctx context.Context, userID string) (*UsageStats, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	return s.usageRepo.GetUsage(ctx, userID)
}

// translateWithMedicalTerms는 의료 용어 사전으로 번역합니다.
func (s *TranslationService) translateWithMedicalTerms(text, sourceLang, targetLang string) string {
	key := sourceLang + "->" + targetLang
	if terms, ok := s.medicalTerms[key]; ok {
		result := text
		for src, dst := range terms {
			result = strings.ReplaceAll(result, src, dst)
		}
		if result != text {
			return result
		}
	}

	// 시뮬레이션: 실제로는 외부 번역 API 호출
	return fmt.Sprintf("[%s→%s] %s", sourceLang, targetLang, text)
}

// detectLanguageCode는 간단한 언어 감지를 수행합니다.
func detectLanguageCode(text string) string {
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7A3 {
			return "ko"
		}
		if r >= 0x3040 && r <= 0x309F || r >= 0x30A0 && r <= 0x30FF {
			return "ja"
		}
		if r >= 0x4E00 && r <= 0x9FFF {
			return "zh"
		}
		if r >= 0x0900 && r <= 0x097F {
			return "hi"
		}
	}
	return "en"
}

func languageName(code string) string {
	names := map[string]string{
		"ko": "Korean", "en": "English", "ja": "Japanese", "zh": "Chinese",
		"fr": "French", "de": "German", "es": "Spanish", "hi": "Hindi",
	}
	if n, ok := names[code]; ok {
		return n
	}
	return code
}

func defaultLanguages() []SupportedLanguage {
	return []SupportedLanguage{
		{Code: "ko", Name: "Korean", NativeName: "한국어", SupportsMedical: true, SupportsRealtime: true},
		{Code: "en", Name: "English", NativeName: "English", SupportsMedical: true, SupportsRealtime: true},
		{Code: "ja", Name: "Japanese", NativeName: "日本語", SupportsMedical: true, SupportsRealtime: true},
		{Code: "zh", Name: "Chinese", NativeName: "中文", SupportsMedical: true, SupportsRealtime: true},
		{Code: "fr", Name: "French", NativeName: "Français", SupportsMedical: true, SupportsRealtime: true},
		{Code: "de", Name: "German", NativeName: "Deutsch", SupportsMedical: true, SupportsRealtime: false},
		{Code: "es", Name: "Spanish", NativeName: "Español", SupportsMedical: true, SupportsRealtime: true},
		{Code: "hi", Name: "Hindi", NativeName: "हिन्दी", SupportsMedical: false, SupportsRealtime: true},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", SupportsMedical: false, SupportsRealtime: false},
		{Code: "ar", Name: "Arabic", NativeName: "العربية", SupportsMedical: false, SupportsRealtime: false},
	}
}

func defaultMedicalTerms() map[string]map[string]string {
	return map[string]map[string]string{
		"ko->en": {
			"혈압":    "blood pressure",
			"혈당":    "blood glucose",
			"체온":    "body temperature",
			"처방전":   "prescription",
			"복용량":   "dosage",
			"부작용":   "side effect",
			"진단":    "diagnosis",
			"증상":    "symptom",
			"항생제":   "antibiotic",
			"진통제":   "analgesic",
		},
		"en->ko": {
			"blood pressure":   "혈압",
			"blood glucose":    "혈당",
			"body temperature": "체온",
			"prescription":     "처방전",
			"dosage":           "복용량",
			"side effect":      "부작용",
			"diagnosis":        "진단",
			"symptom":          "증상",
			"antibiotic":       "항생제",
			"analgesic":        "진통제",
		},
		"ko->ja": {
			"혈압":  "血圧",
			"혈당":  "血糖",
			"체온":  "体温",
			"처방전": "処方箋",
			"진단":  "診断",
		},
	}
}
