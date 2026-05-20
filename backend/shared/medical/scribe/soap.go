// Package scribe는 Ambient Scribe — 음성/대화 → 구조화된 SOAP 노트 생성기입니다.
//
// P10-09b D-05 차별화의 핵심.
// 입력: ASR 트랜스크립트 또는 오디오 청크 + POCT 측정값 + CoT 결과
// 출력: SOAP (Subjective/Objective/Assessment/Plan) + FHIR DocumentReference
//
// stt 어댑터를 주입하면 오디오 직접 입력 → 트랜스크립트 자동 변환 후 SOAP 생성.
package scribe

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/medical/llm"
	"github.com/manpasik/backend/shared/medical/stt"
)

// ============================================================================
// SOAP 노트 모델
// ============================================================================

// SOAPNote는 표준 SOAP 노트입니다.
//
// S (Subjective): 주관적 — 환자 호소
// O (Objective): 객관적 — 측정값/관찰
// A (Assessment): 평가 — 진단/추론
// P (Plan): 계획 — 치료/추적
type SOAPNote struct {
	ID            string
	PatientID     string
	EncounterID   string
	Provider      string
	Date          time.Time

	Subjective  *Subjective
	Objective   *Objective
	Assessment  *Assessment
	Plan        *Plan

	Locale       string  // ko, en, ja, zh
	GeneratedBy  string  // "scribe-v1.0"
	HumanReviewed bool
	ReviewerID   string
	ReviewedAt   *time.Time

	// 신뢰도/감사
	Confidence   float64
	WordCount    int
}

// Subjective는 환자 주관 정보입니다.
type Subjective struct {
	ChiefComplaint string   // 주증상
	HistoryOfPresentIllness string // 현병력
	PastMedicalHistory []string
	FamilyHistory []string
	SocialHistory []string
	Medications []string
	Allergies []string
	ReviewOfSystems map[string]string // body_system → 양성/음성
}

// Objective는 객관적 측정/관찰 정보입니다.
type Objective struct {
	VitalSigns []*VitalSign
	PhysicalExam map[string]string // 신체 부위 → 소견
	LabResults []*LabResult
	Imaging []string
	OtherTests []string
}

// VitalSign은 생체 활력징후입니다.
type VitalSign struct {
	Type      string // "blood_pressure", "heart_rate", "temperature"
	Value     string // "130/80", "72", "36.8"
	Unit      string
	Timestamp time.Time
	Abnormal  bool
}

// LabResult는 검사 결과입니다.
type LabResult struct {
	TestCode    string  // LOINC
	TestName    string
	Value       float64
	Unit        string
	RefLow      float64
	RefHigh     float64
	Flag        string // H, L, N
	Alpha       float64
	CILow       float64
	CIHigh      float64
}

// Assessment는 진단/평가입니다.
type Assessment struct {
	PrimaryDiagnosis string
	ICD10           string
	Differentials   []string
	Reasoning       string  // CoT 요약
	Confidence      float64
	RiskLevel       string  // low/moderate/high/critical
}

// Plan은 치료/추적 계획입니다.
type Plan struct {
	Medications     []*MedicationOrder
	Procedures      []string
	FollowUp        string
	PatientEducation []string
	Referrals       []string
	NextVisit       *time.Time
}

// MedicationOrder는 처방 약물입니다.
type MedicationOrder struct {
	DrugName  string
	Dosage    string
	Frequency string
	Duration  string
	Route     string // oral, IV, topical
}

// ============================================================================
// 음성 트랜스크립트 입력
// ============================================================================

// Transcript는 ASR(Automatic Speech Recognition) 트랜스크립트입니다.
type Transcript struct {
	SessionID  string
	Segments   []*TranscriptSegment
	StartedAt  time.Time
	EndedAt    time.Time
	Duration   time.Duration
}

// TranscriptSegment는 단일 발화 세그먼트입니다.
type TranscriptSegment struct {
	Speaker    string // "doctor", "patient", "unknown"
	Text       string
	StartTime  float64 // 초 단위 오프셋
	EndTime    float64
	Confidence float64 // 0~1
}

// FilterBySpeaker는 특정 화자의 세그먼트만 반환합니다.
func (t *Transcript) FilterBySpeaker(speaker string) []*TranscriptSegment {
	var result []*TranscriptSegment
	for _, s := range t.Segments {
		if s.Speaker == speaker {
			result = append(result, s)
		}
	}
	return result
}

// PatientSpeech는 환자 발화만 합친 텍스트를 반환합니다.
func (t *Transcript) PatientSpeech() string {
	var sb strings.Builder
	for _, s := range t.FilterBySpeaker("patient") {
		sb.WriteString(s.Text + " ")
	}
	return strings.TrimSpace(sb.String())
}

// DoctorSpeech는 의사 발화만 합친 텍스트를 반환합니다.
func (t *Transcript) DoctorSpeech() string {
	var sb strings.Builder
	for _, s := range t.FilterBySpeaker("doctor") {
		sb.WriteString(s.Text + " ")
	}
	return strings.TrimSpace(sb.String())
}

// ============================================================================
// SOAP 추출 키워드
// ============================================================================

var (
	// 주증상 표지어
	chiefComplaintMarkers = []string{
		"아파", "통증", "증상", "느낌", "불편", "이상",
	}
	// 약물 표지어
	medicationMarkers = []string{
		"약을 먹", "복용", "처방받", "투약",
	}
	// 알레르기 표지어
	allergyMarkers = []string{
		"알레르기", "알러지", "두드러기",
	}
	// 활력징후 패턴 (혈압/체온/심박수)
	vitalSignKeywords = map[string]string{
		"혈압":   "blood_pressure",
		"심박수":  "heart_rate",
		"체온":   "temperature",
		"호흡수":  "respiratory_rate",
		"산소포화도": "oxygen_saturation",
	}
)

// ============================================================================
// Scribe 엔진
// ============================================================================

// Scribe는 트랜스크립트를 SOAP 노트로 변환합니다.
//
// stt 어댑터가 주입되면 GenerateFromAudio로 오디오 직접 입력 가능.
// llm 어댑터가 주입되면 GenerateReviewSummary로 의사 검토용 요약 생성 가능.
type Scribe struct {
	defaultLocale string
	sttAdapter    stt.Adapter
	llmAdapter    llm.Adapter
}

// NewScribe는 새 Scribe를 생성합니다.
func NewScribe(locale string) *Scribe {
	if locale == "" {
		locale = "ko"
	}
	return &Scribe{defaultLocale: locale}
}

// SetSTTAdapter는 음성 → 텍스트 변환 어댑터를 주입합니다.
//
// 주입 시 GenerateFromAudio로 오디오 청크 직접 입력 가능.
// 미주입 시 GenerateFromAudio는 에러 반환.
func (s *Scribe) SetSTTAdapter(adapter stt.Adapter) {
	if adapter != nil {
		s.sttAdapter = adapter
	}
}

// SetLLMAdapter는 LLM 어댑터를 주입합니다.
//
// 주입 시 GenerateReviewSummary로 의사 검토용 요약 생성 가능.
// 미주입 시 GenerateReviewSummary는 에러 반환.
func (s *Scribe) SetLLMAdapter(adapter llm.Adapter) {
	if adapter != nil {
		s.llmAdapter = adapter
	}
}

// ReviewSummary는 의사 검토용 SOAP 요약입니다.
type ReviewSummary struct {
	NoteID        string
	HighlightedRisks []string  // 빨간 깃발 (이상치, 약물 상호작용 등)
	KeyFindings   []string  // 핵심 발견
	OutstandingQuestions []string // 추가 확인 질문
	ProviderNotes string    // 의사 추가 메모용 자리
	GeneratedBy   string    // LLM provider
	GeneratedAt   time.Time
	LLMTokensUsed int
}

// GenerateReviewSummary는 SOAP 노트를 의사 검토용 요약으로 변환합니다.
//
// LLM에 SOAP 노트를 전달하여 위험 요인, 핵심 발견, 추가 확인 사항을 추출합니다.
// 의사 최종 검토 효율 향상이 목적이며 의료 결정을 자동화하지 않습니다.
func (s *Scribe) GenerateReviewSummary(ctx context.Context, note *SOAPNote) (*ReviewSummary, error) {
	if s.llmAdapter == nil {
		return nil, errors.New("llm adapter not configured (call SetLLMAdapter first)")
	}
	if note == nil {
		return nil, errors.New("note is nil")
	}

	soapText := formatSOAP(note)
	systemPrompt := `당신은 의료 검토 보조원입니다. SOAP 노트를 검토하여:
1. 빨간 깃발 위험 요인 (이상치, 약물 상호작용, 응급 상황 가능성)
2. 핵심 발견 (진단 근거, 주요 증상)
3. 추가 확인이 필요한 질문
을 추출하여 의사가 빠르게 검토할 수 있도록 정리해주세요.
의료 결정을 내리지 마세요. 진단 권고는 의사 영역입니다.`

	req := &llm.Request{
		Model:       "gpt-4", // 운영 시 환경변수
		Temperature: 0.3,     // 의료 도메인은 낮은 temperature
		MaxTokens:   1000,
		Messages: []*llm.Message{
			{Role: llm.RoleSystem, Content: systemPrompt},
			{Role: llm.RoleUser, Content: soapText},
		},
		RequireDisclaimer: true,
	}

	resp, err := s.llmAdapter.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM 요약 생성 실패: %w", err)
	}

	// 응답에서 섹션 추출 (단순 파싱)
	risks, findings, questions := parseSummaryResponse(resp.Content)

	return &ReviewSummary{
		NoteID:               note.ID,
		HighlightedRisks:     risks,
		KeyFindings:          findings,
		OutstandingQuestions: questions,
		GeneratedBy:          resp.Provider,
		GeneratedAt:          time.Now().UTC(),
		LLMTokensUsed:        resp.TotalTokens,
	}, nil
}

// parseSummaryResponse는 LLM 응답에서 3개 섹션을 단순 파싱합니다.
//
// 운영에서는 구조화된 출력 (JSON Schema 등) 사용 권장.
func parseSummaryResponse(content string) (risks, findings, questions []string) {
	risks = extractListAfterMarker(content, []string{"빨간 깃발", "위험 요인", "Red Flag"})
	findings = extractListAfterMarker(content, []string{"핵심 발견", "주요 발견", "Key Finding"})
	questions = extractListAfterMarker(content, []string{"추가 확인", "추가 질문", "Outstanding"})
	if len(risks) == 0 && len(findings) == 0 && len(questions) == 0 {
		// 파싱 실패 시 전체를 발견으로 분류
		findings = []string{strings.TrimSpace(content)}
	}
	return
}

func extractListAfterMarker(content string, markers []string) []string {
	for _, marker := range markers {
		idx := strings.Index(content, marker)
		if idx < 0 {
			continue
		}
		// marker 이후 다음 섹션 또는 끝까지
		start := idx + len(marker)
		end := len(content)
		// 다음 섹션 표지어 찾기
		for _, nextMarker := range []string{"빨간 깃발", "핵심 발견", "추가 확인", "DISCLAIMER", "[안내]"} {
			if next := strings.Index(content[start:], nextMarker); next > 0 {
				if start+next < end {
					end = start + next
				}
			}
		}
		section := content[start:end]
		// 줄 단위로 분리, "- " "1." 등 prefix 제거
		var items []string
		for _, line := range strings.Split(section, "\n") {
			cleaned := strings.TrimSpace(line)
			cleaned = strings.TrimPrefix(cleaned, "-")
			cleaned = strings.TrimPrefix(cleaned, "•")
			cleaned = strings.TrimSpace(cleaned)
			if len(cleaned) > 5 {
				items = append(items, cleaned)
			}
		}
		if len(items) > 0 {
			return items
		}
	}
	return nil
}

// GenerateOptions는 SOAP 생성 옵션입니다.
type GenerateOptions struct {
	PatientID    string
	EncounterID  string
	Provider     string
	LabResults   []*LabResult
	VitalSigns   []*VitalSign
	Diagnosis    string
	ICD10        string
	Differentials []string
	Reasoning    string
	RiskLevel    string
}

// AudioInput은 오디오 청크 직접 입력을 위한 옵션입니다.
type AudioInput struct {
	Chunks  []*stt.AudioChunk
	Speakers map[string]string // chunkIndex → "doctor" / "patient"
}

// GenerateFromAudio는 오디오 청크를 STT로 변환하여 SOAP 노트를 생성합니다.
//
// 사전 조건: SetSTTAdapter로 STT 어댑터 주입.
// 흐름: AudioChunks → STT 변환 → 통합 Transcript → Generate.
func (s *Scribe) GenerateFromAudio(
	ctx context.Context,
	audio *AudioInput,
	opts *GenerateOptions,
) (*SOAPNote, error) {
	if s.sttAdapter == nil {
		return nil, errors.New("stt adapter not configured (call SetSTTAdapter first)")
	}
	if audio == nil || len(audio.Chunks) == 0 {
		return nil, errors.New("audio chunks required")
	}

	// stt.StreamingAggregator로 청크를 누적
	agg := stt.NewStreamingAggregator(s.sttAdapter)
	for i, chunk := range audio.Chunks {
		// speaker 오버라이드
		if audio.Speakers != nil {
			if sp, ok := audio.Speakers[fmt.Sprintf("%d", i)]; ok {
				chunk.Speaker = sp
			}
		}
		if _, err := agg.AppendChunk(ctx, chunk); err != nil {
			return nil, fmt.Errorf("STT chunk %d: %w", i, err)
		}
	}

	// 누적 트랜스크립트 → scribe.Transcript로 변환
	sttTranscript := agg.FinalizeSession(audio.Chunks[0].SessionID)
	if sttTranscript == nil {
		return nil, errors.New("STT 변환 결과 없음")
	}

	// stt.Segment → scribe.TranscriptSegment 매핑
	scribeSegments := make([]*TranscriptSegment, 0, len(sttTranscript.Segments))
	for _, seg := range sttTranscript.Segments {
		scribeSegments = append(scribeSegments, &TranscriptSegment{
			Speaker:    seg.Speaker,
			Text:       seg.Text,
			StartTime:  seg.StartTime,
			EndTime:    seg.EndTime,
			Confidence: seg.Confidence,
		})
	}

	transcript := &Transcript{
		SessionID: sttTranscript.SessionID,
		Segments:  scribeSegments,
		StartedAt: time.Now().UTC().Add(-time.Duration(sttTranscript.DurationSec) * time.Second),
		EndedAt:   time.Now().UTC(),
		Duration:  time.Duration(sttTranscript.DurationSec) * time.Second,
	}

	return s.Generate(transcript, opts)
}

// Generate는 트랜스크립트와 측정값으로 SOAP 노트를 생성합니다.
func (s *Scribe) Generate(transcript *Transcript, opts *GenerateOptions) (*SOAPNote, error) {
	if transcript == nil {
		return nil, errors.New("transcript required")
	}
	if opts == nil || opts.PatientID == "" {
		return nil, errors.New("patient_id required")
	}

	now := time.Now().UTC()
	note := &SOAPNote{
		ID:          fmt.Sprintf("soap-%d", now.UnixNano()),
		PatientID:   opts.PatientID,
		EncounterID: opts.EncounterID,
		Provider:    opts.Provider,
		Date:        now,
		Locale:      s.defaultLocale,
		GeneratedBy: "scribe-v1.0",
	}

	// S (Subjective) — 환자 발화 분석
	note.Subjective = s.extractSubjective(transcript)

	// O (Objective) — 측정값 + 활력징후
	note.Objective = &Objective{
		LabResults: opts.LabResults,
		VitalSigns: opts.VitalSigns,
	}

	// A (Assessment)
	note.Assessment = &Assessment{
		PrimaryDiagnosis: opts.Diagnosis,
		ICD10:           opts.ICD10,
		Differentials:   opts.Differentials,
		Reasoning:       opts.Reasoning,
		RiskLevel:       opts.RiskLevel,
	}

	// P (Plan) — 의사 발화에서 추출
	note.Plan = s.extractPlan(transcript)

	// 신뢰도 + 단어 수
	note.Confidence = computeConfidence(transcript)
	note.WordCount = countWords(formatSOAP(note))

	return note, nil
}

func (s *Scribe) extractSubjective(transcript *Transcript) *Subjective {
	patientText := transcript.PatientSpeech()
	subj := &Subjective{
		ReviewOfSystems: make(map[string]string),
	}

	// 주증상 추출 (첫 번째 chief complaint marker 인접 문장)
	for _, marker := range chiefComplaintMarkers {
		if idx := strings.Index(patientText, marker); idx >= 0 {
			start := max(0, idx-50)
			end := min(len(patientText), idx+50)
			subj.ChiefComplaint = strings.TrimSpace(patientText[start:end])
			break
		}
	}
	if subj.ChiefComplaint == "" {
		// 첫 100자를 chief complaint로
		if len(patientText) > 100 {
			subj.ChiefComplaint = patientText[:100]
		} else {
			subj.ChiefComplaint = patientText
		}
	}

	// 현병력 — 환자 발화 전체 요약
	subj.HistoryOfPresentIllness = patientText

	// 약물 이력
	for _, m := range medicationMarkers {
		if strings.Contains(patientText, m) {
			subj.Medications = append(subj.Medications, "[음성에서 약물 언급 감지 - 검토 필요]")
			break
		}
	}

	// 알레르기
	for _, m := range allergyMarkers {
		if strings.Contains(patientText, m) {
			subj.Allergies = append(subj.Allergies, "[음성에서 알레르기 언급 감지 - 검토 필요]")
			break
		}
	}

	return subj
}

func (s *Scribe) extractPlan(transcript *Transcript) *Plan {
	doctorText := transcript.DoctorSpeech()
	plan := &Plan{}

	// 약물 처방 검출 (간단)
	if strings.Contains(doctorText, "처방") || strings.Contains(doctorText, "약을 드릴") {
		plan.Medications = []*MedicationOrder{{
			DrugName: "[음성에서 처방 언급 감지 - 약사 확인 필요]",
		}}
	}

	// 추적 진료
	if strings.Contains(doctorText, "다음 주") || strings.Contains(doctorText, "내원") {
		plan.FollowUp = "[다음 진료 일정 음성에서 감지 - 일정 확정 필요]"
	}

	// 환자 교육
	if strings.Contains(doctorText, "주의") || strings.Contains(doctorText, "관리") {
		plan.PatientEducation = []string{"[음성에서 교육 메시지 감지]"}
	}

	return plan
}

func computeConfidence(transcript *Transcript) float64 {
	if len(transcript.Segments) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range transcript.Segments {
		sum += s.Confidence
	}
	return sum / float64(len(transcript.Segments))
}

func countWords(text string) int {
	return len(strings.Fields(text))
}

// ============================================================================
// SOAP 포맷팅
// ============================================================================

// FormatSOAP는 SOAP 노트를 텍스트로 포맷팅합니다.
func FormatSOAP(note *SOAPNote) string {
	return formatSOAP(note)
}

func formatSOAP(note *SOAPNote) string {
	if note == nil {
		return ""
	}
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# SOAP 노트 — %s\n\n", note.Date.Format("2006-01-02 15:04")))
	sb.WriteString(fmt.Sprintf("**환자 ID**: %s\n", note.PatientID))
	if note.Provider != "" {
		sb.WriteString(fmt.Sprintf("**작성자**: %s\n", note.Provider))
	}
	sb.WriteString("\n")

	// S
	sb.WriteString("## S (Subjective)\n")
	if note.Subjective != nil {
		if note.Subjective.ChiefComplaint != "" {
			sb.WriteString(fmt.Sprintf("**주증상**: %s\n", note.Subjective.ChiefComplaint))
		}
		if note.Subjective.HistoryOfPresentIllness != "" {
			sb.WriteString(fmt.Sprintf("**현병력**: %s\n", note.Subjective.HistoryOfPresentIllness))
		}
		if len(note.Subjective.Medications) > 0 {
			sb.WriteString("**복용 약물**:\n")
			for _, m := range note.Subjective.Medications {
				sb.WriteString("  - " + m + "\n")
			}
		}
		if len(note.Subjective.Allergies) > 0 {
			sb.WriteString("**알레르기**:\n")
			for _, a := range note.Subjective.Allergies {
				sb.WriteString("  - " + a + "\n")
			}
		}
	}
	sb.WriteString("\n")

	// O
	sb.WriteString("## O (Objective)\n")
	if note.Objective != nil {
		if len(note.Objective.VitalSigns) > 0 {
			sb.WriteString("**활력징후**:\n")
			for _, v := range note.Objective.VitalSigns {
				flag := ""
				if v.Abnormal {
					flag = " ⚠️"
				}
				sb.WriteString(fmt.Sprintf("  - %s: %s %s%s\n", v.Type, v.Value, v.Unit, flag))
			}
		}
		if len(note.Objective.LabResults) > 0 {
			sb.WriteString("**검사 결과**:\n")
			for _, r := range note.Objective.LabResults {
				flag := ""
				if r.Flag == "H" {
					flag = " ⬆ (High)"
				} else if r.Flag == "L" {
					flag = " ⬇ (Low)"
				}
				ci := ""
				if r.CILow != 0 || r.CIHigh != 0 {
					ci = fmt.Sprintf(" [95%% CI: %.2f-%.2f]", r.CILow, r.CIHigh)
				}
				sb.WriteString(fmt.Sprintf("  - %s (%s): %.2f %s (정상: %.0f-%.0f)%s%s\n",
					r.TestName, r.TestCode, r.Value, r.Unit, r.RefLow, r.RefHigh, flag, ci))
			}
		}
	}
	sb.WriteString("\n")

	// A
	sb.WriteString("## A (Assessment)\n")
	if note.Assessment != nil {
		if note.Assessment.PrimaryDiagnosis != "" {
			sb.WriteString(fmt.Sprintf("**주 진단**: %s", note.Assessment.PrimaryDiagnosis))
			if note.Assessment.ICD10 != "" {
				sb.WriteString(fmt.Sprintf(" (ICD-10: %s)", note.Assessment.ICD10))
			}
			sb.WriteString("\n")
		}
		if len(note.Assessment.Differentials) > 0 {
			sb.WriteString("**감별진단**:\n")
			for _, d := range note.Assessment.Differentials {
				sb.WriteString("  - " + d + "\n")
			}
		}
		if note.Assessment.RiskLevel != "" {
			sb.WriteString(fmt.Sprintf("**위험도**: %s\n", note.Assessment.RiskLevel))
		}
		if note.Assessment.Reasoning != "" {
			sb.WriteString(fmt.Sprintf("**근거**: %s\n", note.Assessment.Reasoning))
		}
	}
	sb.WriteString("\n")

	// P
	sb.WriteString("## P (Plan)\n")
	if note.Plan != nil {
		if len(note.Plan.Medications) > 0 {
			sb.WriteString("**처방**:\n")
			for _, m := range note.Plan.Medications {
				sb.WriteString(fmt.Sprintf("  - %s %s %s %s\n",
					m.DrugName, m.Dosage, m.Frequency, m.Duration))
			}
		}
		if note.Plan.FollowUp != "" {
			sb.WriteString(fmt.Sprintf("**추적**: %s\n", note.Plan.FollowUp))
		}
		if len(note.Plan.PatientEducation) > 0 {
			sb.WriteString("**환자 교육**:\n")
			for _, e := range note.Plan.PatientEducation {
				sb.WriteString("  - " + e + "\n")
			}
		}
		if len(note.Plan.Referrals) > 0 {
			sb.WriteString("**진료 의뢰**:\n")
			for _, r := range note.Plan.Referrals {
				sb.WriteString("  - " + r + "\n")
			}
		}
	}

	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("_생성: %s, 신뢰도: %.0f%%, 단어 수: %d_\n",
		note.GeneratedBy, note.Confidence*100, note.WordCount))
	if note.HumanReviewed {
		sb.WriteString(fmt.Sprintf("_검토자: %s_\n", note.ReviewerID))
	} else {
		sb.WriteString("_⚠️ 의사 검토 미완료_\n")
	}

	return sb.String()
}

// ============================================================================
// FHIR DocumentReference 출력
// ============================================================================

// ToFHIRDocumentReference는 SOAP 노트를 FHIR DocumentReference로 변환합니다.
func ToFHIRDocumentReference(note *SOAPNote) (*fhir.DocumentReference, error) {
	if note == nil {
		return nil, errors.New("note is nil")
	}
	soapText := formatSOAP(note)

	return &fhir.DocumentReference{
		ResourceType: "DocumentReference",
		ID:           fmt.Sprintf("docref-%s", note.ID),
		FhirVersion:  fhir.VersionR5,
		Status:       "current",
		DocStatus:    docStatusForNote(note),
		Type: fhir.CodeableConcept{
			Coding: []fhir.Coding{{
				System:  "http://loinc.org",
				Code:    "11488-4",
				Display: "Consultation note",
			}},
		},
		Subject: fhir.Reference{Reference: "Patient/" + note.PatientID},
		Date:    note.Date,
		Description: fmt.Sprintf("화상진료 SOAP 노트 (신뢰도 %.0f%%)", note.Confidence*100),
		Content: []fhir.DocContent{{
			Attachment: fhir.Attachment{
				ContentType: "text/markdown",
				Data:        []byte(soapText),
				Title:       "SOAP Note",
				Size:        len(soapText),
			},
		}},
	}, nil
}

func docStatusForNote(note *SOAPNote) string {
	if note.HumanReviewed {
		return "final"
	}
	return "preliminary"
}

// ============================================================================
// 헬퍼
// ============================================================================

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SortVitalSignsByTime은 활력징후를 시간순 정렬합니다.
func SortVitalSignsByTime(vitals []*VitalSign) {
	sort.Slice(vitals, func(i, j int) bool {
		return vitals[i].Timestamp.Before(vitals[j].Timestamp)
	})
}
