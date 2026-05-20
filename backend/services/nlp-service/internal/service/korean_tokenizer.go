// Package service: 한글 토크나이저 + 의도 분류
//
// 외부 라이브러리 의존 없이 정규식 기반 한글 토큰화와 의도 분류를 제공합니다.
// 의료/건강 도메인에 특화된 어휘로 구성됩니다.
package service

import (
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// ============================================================================
// 의도 카테고리
// ============================================================================

const (
	IntentQuestion   = "question"   // 질문
	IntentSymptom    = "symptom"    // 증상 신고
	IntentReservation = "reservation" // 예약/일정
	IntentEmergency  = "emergency"  // 긴급
	IntentSmalltalk  = "smalltalk"  // 잡담
	IntentGratitude  = "gratitude"  // 감사
)

// ============================================================================
// 토큰 타입
// ============================================================================

const (
	TokenWord    = "word"     // 일반 단어
	TokenNumber  = "number"   // 숫자
	TokenTime    = "time"     // 시간 표현
	TokenUnit    = "unit"     // 단위 (mg/dL, mmHg)
	TokenParticle = "particle" // 조사
)

// Token은 토큰화 결과의 단일 토큰입니다.
type Token struct {
	Text      string
	Type      string
	StartIdx  int
	EndIdx    int
}

// ============================================================================
// 한글 토크나이저
// ============================================================================

// KoreanTokenizer는 한글 텍스트를 토큰화합니다.
type KoreanTokenizer struct {
	timeRegex   *regexp.Regexp
	numberRegex *regexp.Regexp
	unitRegex   *regexp.Regexp
}

// NewKoreanTokenizer는 새 한글 토크나이저를 생성합니다.
func NewKoreanTokenizer() *KoreanTokenizer {
	return &KoreanTokenizer{
		// 시간 표현: "오늘", "어제", "오전 9시", "내일 아침", "3시간 전"
		timeRegex: regexp.MustCompile(
			`(오늘|어제|내일|모레|그제|지금|방금|아침|점심|저녁|밤|오전|오후|새벽)|(\d+(시간|일|주|개월|달|년|분)\s*(전|후|뒤))|(\d+시(\s*\d+분)?)`,
		),
		// 숫자: 정수 + 소수
		numberRegex: regexp.MustCompile(`\d+(\.\d+)?`),
		// 의료 단위
		unitRegex: regexp.MustCompile(
			`(mg/dL|mmHg|bpm|°C|kg|cm|mL|단위|밀리그램|°F|µg|mEq/L|IU|kcal|회)`,
		),
	}
}

// Tokenize는 텍스트를 토큰 목록으로 분해합니다.
func (kt *KoreanTokenizer) Tokenize(text string) []Token {
	var tokens []Token
	text = strings.TrimSpace(text)
	if text == "" {
		return tokens
	}

	// 1. 시간 표현 우선 추출
	tokens = append(tokens, kt.extractByRegex(text, kt.timeRegex, TokenTime)...)
	// 2. 단위 추출
	tokens = append(tokens, kt.extractByRegex(text, kt.unitRegex, TokenUnit)...)
	// 3. 숫자 추출 (시간/단위에 포함되지 않은 것)
	tokens = append(tokens, kt.extractByRegex(text, kt.numberRegex, TokenNumber)...)
	// 4. 어절 분리 (공백 기준 + 한글 형태소 휴리스틱)
	tokens = append(tokens, kt.extractWords(text)...)

	return dedupTokens(tokens)
}

func (kt *KoreanTokenizer) extractByRegex(text string, re *regexp.Regexp, tokenType string) []Token {
	var tokens []Token
	matches := re.FindAllStringIndex(text, -1)
	for _, m := range matches {
		tokens = append(tokens, Token{
			Text:     text[m[0]:m[1]],
			Type:     tokenType,
			StartIdx: m[0],
			EndIdx:   m[1],
		})
	}
	return tokens
}

func (kt *KoreanTokenizer) extractWords(text string) []Token {
	var tokens []Token
	fields := strings.Fields(text)
	idx := 0
	for _, f := range fields {
		// 한글 어절에서 조사 분리 (을/를/이/가/은/는/에/의/도/만)
		root, particle := splitParticle(f)
		if root != "" {
			tokens = append(tokens, Token{
				Text:     root,
				Type:     TokenWord,
				StartIdx: idx,
				EndIdx:   idx + len(root),
			})
		}
		if particle != "" {
			tokens = append(tokens, Token{
				Text:     particle,
				Type:     TokenParticle,
				StartIdx: idx + len(root),
				EndIdx:   idx + len(f),
			})
		}
		idx += len(f) + 1
	}
	return tokens
}

// splitParticle은 어절 끝의 한글 조사를 분리합니다.
func splitParticle(word string) (root, particle string) {
	particles := []string{"을", "를", "이", "가", "은", "는", "에서", "에게", "에", "의", "도", "만", "와", "과", "로", "으로"}

	for _, p := range particles {
		if strings.HasSuffix(word, p) && len(word) > len(p) {
			r := word[:len(word)-len(p)]
			// 어근의 마지막 룬이 한글이어야 영문 조사 오매칭 방지
			last, _ := utf8.DecodeLastRuneInString(r)
			if last != utf8.RuneError && isHangul(last) {
				return r, p
			}
		}
	}
	return word, ""
}

func isHangul(r rune) bool {
	return unicode.Is(unicode.Hangul, r)
}

func dedupTokens(tokens []Token) []Token {
	seen := make(map[string]bool)
	var unique []Token
	for _, t := range tokens {
		key := t.Type + ":" + t.Text
		if !seen[key] {
			seen[key] = true
			unique = append(unique, t)
		}
	}
	return unique
}

// ============================================================================
// 의도 분류기
// ============================================================================

// IntentResult는 의도 분류 결과입니다.
type IntentResult struct {
	Intent     string
	Confidence float64
	Severity   int // 0~100, 긴급도 점수
	Keywords   []string
	Tokens     []Token
}

// IntentRouter는 한글 텍스트의 의도를 분류합니다.
type IntentRouter struct {
	tokenizer *KoreanTokenizer
}

// NewIntentRouter는 새 의도 라우터를 생성합니다.
func NewIntentRouter() *IntentRouter {
	return &IntentRouter{
		tokenizer: NewKoreanTokenizer(),
	}
}

// 의도별 키워드 사전
var (
	questionKeywords = []string{"무엇", "뭐", "어떻게", "왜", "언제", "어디", "누구", "얼마", "?", "인가요", "있나요", "할까요", "되나요"}

	symptomKeywords = []string{"아파", "아프", "통증", "증상", "느낌", "두통", "어지러", "메스꺼", "기침", "발열", "열이", "구토", "설사", "변비", "피로", "무기력"}

	reservationKeywords = []string{"예약", "일정", "약속", "스케줄", "잡고", "취소", "변경", "병원", "진료", "방문"}

	emergencyKeywords = []string{"응급", "긴급", "119", "구급", "심각", "위급", "당장", "쓰러", "기절", "의식없", "호흡곤란", "심정지", "출혈", "골절"}

	gratitudeKeywords = []string{"감사", "고마워", "고맙", "땡큐", "thanks"}

	smalltalkKeywords = []string{"안녕", "반가", "잘지내", "오랜만", "날씨"}
)

// 긴급 키워드별 가중치
var emergencyWeights = map[string]int{
	"119":    100,
	"심정지":   100,
	"의식없":   95,
	"기절":    90,
	"쓰러":    85,
	"호흡곤란":  85,
	"흉통":    80,
	"출혈":    75,
	"심각":    70,
	"위급":    70,
	"긴급":    65,
	"응급":    60,
	"구급":    60,
	"골절":    55,
	"고열":    40,
	"심한":    35,
}

// Classify는 텍스트의 의도를 분류합니다.
func (ir *IntentRouter) Classify(text string) *IntentResult {
	tokens := ir.tokenizer.Tokenize(text)
	lower := strings.ToLower(text)

	// 의도별 점수 계산
	scores := map[string]int{
		IntentQuestion:    countMatches(lower, questionKeywords),
		IntentSymptom:     countMatches(lower, symptomKeywords),
		IntentReservation: countMatches(lower, reservationKeywords),
		IntentEmergency:   countMatches(lower, emergencyKeywords),
		IntentGratitude:   countMatches(lower, gratitudeKeywords),
		IntentSmalltalk:   countMatches(lower, smalltalkKeywords),
	}

	// 한국어 증상 DB 매칭 → 증상 점수 가산
	for keyword := range koreanSymptomDB {
		if strings.Contains(text, keyword) {
			scores[IntentSymptom]++
		}
	}

	// 긴급도 점수 계산
	severity := 0
	var matchedKeywords []string
	for kw, weight := range emergencyWeights {
		if strings.Contains(text, kw) {
			if weight > severity {
				severity = weight
			}
			matchedKeywords = append(matchedKeywords, kw)
		}
	}

	// 긴급 의도가 있으면 우선 반환
	if scores[IntentEmergency] > 0 || severity >= 70 {
		return &IntentResult{
			Intent:     IntentEmergency,
			Confidence: 0.95,
			Severity:   severity,
			Keywords:   matchedKeywords,
			Tokens:     tokens,
		}
	}

	// 최고 점수 의도 선택
	bestIntent := IntentSmalltalk
	bestScore := 0
	for intent, score := range scores {
		if score > bestScore {
			bestScore = score
			bestIntent = intent
		}
	}

	confidence := 0.5
	if bestScore > 0 {
		confidence = 0.6 + float64(bestScore)*0.1
		if confidence > 0.95 {
			confidence = 0.95
		}
	}

	return &IntentResult{
		Intent:     bestIntent,
		Confidence: confidence,
		Severity:   severity,
		Keywords:   matchedKeywords,
		Tokens:     tokens,
	}
}

func countMatches(text string, keywords []string) int {
	count := 0
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			count++
		}
	}
	return count
}

// ============================================================================
// 엔티티 추출
// ============================================================================

// EntityExtraction은 추출된 엔티티 목록입니다.
type EntityExtraction struct {
	Numbers     []string  // 숫자
	Times       []string  // 시간 표현
	Units       []string  // 단위
	BodyParts   []string  // 신체 부위
	Symptoms    []string  // 증상
	ExtractedAt time.Time
}

// 신체 부위 사전
var bodyParts = []string{
	"머리", "목", "어깨", "팔", "손", "가슴", "복부", "배", "허리", "다리", "발", "무릎", "엉덩이", "등", "눈", "귀", "코", "입", "이마", "관절",
}

// ExtractEntities는 텍스트에서 엔티티를 추출합니다.
func (ir *IntentRouter) ExtractEntities(text string) *EntityExtraction {
	tokens := ir.tokenizer.Tokenize(text)

	ext := &EntityExtraction{
		ExtractedAt: time.Now().UTC(),
	}

	for _, t := range tokens {
		switch t.Type {
		case TokenNumber:
			ext.Numbers = append(ext.Numbers, t.Text)
		case TokenTime:
			ext.Times = append(ext.Times, t.Text)
		case TokenUnit:
			ext.Units = append(ext.Units, t.Text)
		}
	}

	// 신체 부위
	for _, bp := range bodyParts {
		if strings.Contains(text, bp) {
			ext.BodyParts = append(ext.BodyParts, bp)
		}
	}

	// 증상
	for keyword := range koreanSymptomDB {
		if strings.Contains(text, keyword) {
			ext.Symptoms = append(ext.Symptoms, keyword)
		}
	}

	return ext
}
