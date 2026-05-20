// Package poct는 POCT1-A (CLSI POCT01-A2) 메시지 파서와 차분식 보정 엔진입니다.
//
// POCT1-A는 Point-of-Care Testing 디바이스 ↔ 게이트웨이 간 표준 메시지 형식.
// SSOT v1.1 (CLAUDE.md): Sdiff_n = S_n − α_n × R_n (alpha 0.98, 범위 0.90~1.10)
//
// **영역 경계 (SSOT 무결성)**: 본 패키지는 SW 표준·AI 코어 영역으로 한정됩니다.
// 디바이스 BoM(MCU/SiPM/Samtec CSI/PPM 등 물리 부품)은 Phase 5 산출물의 SSOT이며
// 본 패키지에서 직접 참조하지 않습니다. AFE 블록은 측정 카테고리(전기화학·광전·열·증폭 등)로
// 추상화하여 하드웨어 구현체와 분리합니다.
package poct

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// SSOT 베이스라인 — 차분식·채널·알고리즘 한정 (CLAUDE.md v2.2 일치)
//
// 주의: 디바이스 물리 사양(CSI 핀 수, 부품 명세 등)은 Phase 5 SSOT.
// 본 패키지는 측정 알고리즘 표준에만 의존합니다.
// ============================================================================

const (
	DefaultAlpha = 0.98
	AlphaMin     = 0.90
	AlphaMax     = 1.10
	MaxChannels  = 1792
)

// ============================================================================
// AFE 블록 — 측정 카테고리 추상화
//
// 본 식별자는 "어떤 측정 원리"를 가리키며, 특정 부품/벤더(SiPM, Hamamatsu 등)나
// 물리적 위치(B1∼B9 슬롯)와 결합되지 않습니다. 슬롯 매핑·부품 선정은 Phase 5 SSOT.
// ============================================================================

// AFEBlock은 측정 원리(category) 식별자입니다.
type AFEBlock string

const (
	BlockElectrochemical AFEBlock = "electrochemical" // 전기화학
	BlockImmunoassay     AFEBlock = "immunoassay"     // 면역검사 (광전 검출 일반)
	BlockColorimetry     AFEBlock = "colorimetry"     // 색측정
	BlockThermal         AFEBlock = "thermal"         // 열 측정/제어
	BlockIsothermalAmp   AFEBlock = "isothermal_amp"  // 등온증폭 (LAMP 등)
	BlockGas             AFEBlock = "gas"             // 가스 측정
	BlockImpedance       AFEBlock = "impedance"       // 임피던스
	BlockRadiation       AFEBlock = "radiation"       // 방사선
	BlockQuantitativePCR AFEBlock = "qpcr"            // 정량 PCR
)

// SampleType은 시료 유형입니다.
type SampleType string

const (
	SampleLiquid SampleType = "liquid"
	SampleGas    SampleType = "gas"
	SampleSolid  SampleType = "solid"
)

// ============================================================================
// POCT1-A 메시지 모델
// ============================================================================

// Message는 POCT1-A 표준 메시지입니다.
//
// 형식 (CLSI POCT01-A2 단순화):
//   MSH|^~\&|SENDER|FACILITY|RECEIVER|FACILITY||TIMESTAMP|MSG_TYPE|MSG_ID
//   PID|||PATIENT_ID
//   OBR|1|ORDER_ID||TEST_CODE^TEST_NAME
//   OBX|1|NM|LOINC^DESC|1|VALUE|UNIT|REF_LO-REF_HI|FLAG||F|TIMESTAMP
type Message struct {
	MessageID   string
	MessageType string // "ORU^R01" (결과 보고)
	Sender      string
	Facility    string
	Timestamp   time.Time
	PatientID   string
	OrderID     string
	TestCode    string
	TestName    string
	Observations []*Observation
}

// Observation은 단일 관측값입니다.
type Observation struct {
	SequenceNum   int
	ValueType     string // "NM" (numeric), "ST" (string), "CE" (coded)
	LOINCCode     string
	Description   string
	Value         float64
	Unit          string
	RefLow        float64
	RefHigh       float64
	Flag          string // "L" (low), "H" (high), "N" (normal), "A" (abnormal)
	Status        string // "F" (final), "P" (preliminary), "C" (corrected)
	Timestamp     time.Time

	// 만파식 확장
	AFEBlock      AFEBlock
	SampleType    SampleType
	RawSignal     float64 // S_n (보정 전)
	RefSignal     float64 // R_n
	Alpha         float64 // α_n
	CILow         float64 // 95% CI 하한
	CIHigh        float64 // 95% CI 상한
}

// IsAbnormal은 관측값이 비정상 범위에 있는지 반환합니다.
func (o *Observation) IsAbnormal() bool {
	if o.Flag == "A" || o.Flag == "L" || o.Flag == "H" {
		return true
	}
	if o.RefLow > 0 && o.Value < o.RefLow {
		return true
	}
	if o.RefHigh > 0 && o.Value > o.RefHigh {
		return true
	}
	return false
}

// ============================================================================
// POCT1-A 파서
// ============================================================================

// Parser는 POCT1-A 메시지를 파싱합니다.
type Parser struct{}

// NewParser는 새 파서를 생성합니다.
func NewParser() *Parser { return &Parser{} }

// Parse는 텍스트 메시지를 파싱합니다.
//
// 세그먼트는 "\r" 또는 "\n"으로 구분, 필드는 "|", 컴포넌트는 "^"로 구분.
func (p *Parser) Parse(raw string) (*Message, error) {
	if raw == "" {
		return nil, errors.New("empty message")
	}

	lines := splitSegments(raw)
	if len(lines) == 0 {
		return nil, errors.New("no segments found")
	}

	msg := &Message{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "MSH":
			if err := parseMSH(fields, msg); err != nil {
				return nil, fmt.Errorf("MSH: %w", err)
			}
		case "PID":
			if len(fields) > 3 {
				msg.PatientID = fields[3]
			}
		case "OBR":
			if err := parseOBR(fields, msg); err != nil {
				return nil, fmt.Errorf("OBR: %w", err)
			}
		case "OBX":
			obs, err := parseOBX(fields)
			if err != nil {
				return nil, fmt.Errorf("OBX: %w", err)
			}
			msg.Observations = append(msg.Observations, obs)
		}
	}

	if msg.MessageID == "" {
		return nil, errors.New("missing MSH segment")
	}
	return msg, nil
}

func splitSegments(raw string) []string {
	// "\r" 우선, 없으면 "\n"
	if strings.Contains(raw, "\r") {
		return strings.Split(raw, "\r")
	}
	return strings.Split(raw, "\n")
}

func parseMSH(fields []string, msg *Message) error {
	if len(fields) < 10 {
		return errors.New("insufficient MSH fields")
	}
	// MSH|^~\&|SENDER|S_FAC|RECEIVER|R_FAC||TIMESTAMP|MSG_TYPE|MSG_ID
	// idx:[0] [1]   [2]    [3]    [4]    [5]   [6]  [7]      [8]    [9]
	msg.Sender = safeIndex(fields, 2)
	msg.Facility = safeIndex(fields, 3)
	if ts := safeIndex(fields, 7); ts != "" {
		msg.Timestamp = parseTimestamp(ts)
	}
	msg.MessageType = safeIndex(fields, 8)
	msg.MessageID = safeIndex(fields, 9)
	return nil
}

func parseOBR(fields []string, msg *Message) error {
	if len(fields) < 5 {
		return errors.New("insufficient OBR fields")
	}
	msg.OrderID = safeIndex(fields, 2)
	if testField := safeIndex(fields, 4); testField != "" {
		parts := strings.Split(testField, "^")
		msg.TestCode = parts[0]
		if len(parts) > 1 {
			msg.TestName = parts[1]
		}
	}
	return nil
}

func parseOBX(fields []string) (*Observation, error) {
	if len(fields) < 6 {
		return nil, errors.New("insufficient OBX fields")
	}

	obs := &Observation{
		ValueType:   safeIndex(fields, 2),
		Description: "",
	}

	if seq := safeIndex(fields, 1); seq != "" {
		if n, err := strconv.Atoi(seq); err == nil {
			obs.SequenceNum = n
		}
	}

	if codeField := safeIndex(fields, 3); codeField != "" {
		parts := strings.Split(codeField, "^")
		obs.LOINCCode = parts[0]
		if len(parts) > 1 {
			obs.Description = parts[1]
		}
	}

	if vStr := safeIndex(fields, 5); vStr != "" {
		if v, err := strconv.ParseFloat(vStr, 64); err == nil {
			obs.Value = v
		}
	}

	obs.Unit = safeIndex(fields, 6)

	if refRange := safeIndex(fields, 7); refRange != "" {
		parts := strings.SplitN(refRange, "-", 2)
		if len(parts) == 2 {
			if lo, err := strconv.ParseFloat(parts[0], 64); err == nil {
				obs.RefLow = lo
			}
			if hi, err := strconv.ParseFloat(parts[1], 64); err == nil {
				obs.RefHigh = hi
			}
		}
	}

	obs.Flag = safeIndex(fields, 8)
	// OBX|seq|type|code|sub_id|value|unit|ref_range|flag|notes|status|...|timestamp
	// idx:[0][1] [2]  [3] [4]   [5]  [6]  [7]      [8]  [9]   [10]      [11]
	obs.Status = safeIndex(fields, 10)
	if ts := safeIndex(fields, 11); ts != "" {
		obs.Timestamp = parseTimestamp(ts)
	}

	return obs, nil
}

func safeIndex(fields []string, i int) string {
	if i < 0 || i >= len(fields) {
		return ""
	}
	return strings.TrimSpace(fields[i])
}

// parseTimestamp는 "YYYYMMDDHHMMSS" 형식을 time.Time으로 변환합니다.
func parseTimestamp(s string) time.Time {
	formats := []string{
		"20060102150405",
		"200601021504",
		"20060102",
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

// ============================================================================
// 차분식 보정 엔진 (Go 측 포팅, Rust differential과 동일 수식)
// ============================================================================

// DifferentialEngine은 차분식 보정을 수행합니다.
//
// 수식: Sdiff_n = S_n − α_n × R_n
type DifferentialEngine struct {
	mu     sync.RWMutex
	alphas map[string]float64 // 채널별 α_n
}

// NewDifferentialEngine은 새 엔진을 생성합니다.
func NewDifferentialEngine() *DifferentialEngine {
	return &DifferentialEngine{alphas: make(map[string]float64)}
}

// SetAlpha는 채널의 α 계수를 설정합니다 (자동 클램프 [0.90, 1.10]).
func (e *DifferentialEngine) SetAlpha(channel string, alpha float64) float64 {
	clamped := alpha
	if clamped < AlphaMin {
		clamped = AlphaMin
	}
	if clamped > AlphaMax {
		clamped = AlphaMax
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.alphas[channel] = clamped
	return clamped
}

// GetAlpha는 채널의 α 값을 반환합니다 (없으면 기본값 0.98).
func (e *DifferentialEngine) GetAlpha(channel string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if a, ok := e.alphas[channel]; ok {
		return a
	}
	return DefaultAlpha
}

// Correct는 차분식 보정값을 계산합니다.
func (e *DifferentialEngine) Correct(channel string, sN, rN float64) float64 {
	alpha := e.GetAlpha(channel)
	return sN - alpha*rN
}

// CorrectObservation은 관측값에 차분식을 적용합니다.
//
// obs.RawSignal과 obs.RefSignal이 설정되어 있어야 하며, 결과는 obs.Value를 갱신합니다.
func (e *DifferentialEngine) CorrectObservation(obs *Observation, channel string) {
	if obs == nil {
		return
	}
	alpha := e.GetAlpha(channel)
	obs.Alpha = alpha
	obs.Value = obs.RawSignal - alpha*obs.RefSignal
}

// ComputeCI95는 95% 신뢰구간을 계산합니다 (정규 근사).
//
// stdErr: 측정 표준오차. 신뢰구간 = value ± 1.96 × stdErr.
func ComputeCI95(value, stdErr float64) (low, high float64) {
	const z = 1.96
	margin := z * math.Abs(stdErr)
	return value - margin, value + margin
}

// ApplyCI95는 관측값에 95% CI를 적용합니다.
func ApplyCI95(obs *Observation, stdErr float64) {
	low, high := ComputeCI95(obs.Value, stdErr)
	obs.CILow = low
	obs.CIHigh = high
}

// ============================================================================
// 시료 유형 자동 분류
// ============================================================================

// SampleClassifier는 임피던스 시그니처와 카트리지 메타데이터로 시료 유형을 분류합니다.
type SampleClassifier struct{}

// NewSampleClassifier는 새 분류기를 생성합니다.
func NewSampleClassifier() *SampleClassifier { return &SampleClassifier{} }

// ClassifyByImpedance는 임피던스 값으로 시료 유형을 분류합니다.
//
// 임피던스 범위 (대략):
//   - 액체: 100Ω ~ 100kΩ
//   - 가스: > 1MΩ
//   - 고체: 100kΩ ~ 1MΩ
func (c *SampleClassifier) ClassifyByImpedance(impedanceOhm float64) SampleType {
	switch {
	case impedanceOhm < 100_000:
		return SampleLiquid
	case impedanceOhm > 1_000_000:
		return SampleGas
	default:
		return SampleSolid
	}
}

// ============================================================================
// 9블록 자가테스트
// ============================================================================

// SelfTestResult는 9블록 자가테스트 결과입니다.
type SelfTestResult struct {
	Block       AFEBlock
	Healthy     bool
	NoiseLevel  float64
	BaselineMv  float64
	ErrorCode   string
	TestedAt    time.Time
}

// SelfTestRunner는 9블록 자가테스트를 수행합니다.
type SelfTestRunner struct{}

// NewSelfTestRunner는 새 러너를 생성합니다.
func NewSelfTestRunner() *SelfTestRunner { return &SelfTestRunner{} }

// RunAll은 모든 측정 카테고리의 자가테스트 결과를 평가합니다.
//
// 운영에서는 디바이스 펌웨어로부터 진단 데이터를 받아 평가합니다.
// 카테고리 ↔ 물리 슬롯(B1∼B9) 매핑은 Phase 5 SSOT.
func (r *SelfTestRunner) RunAll(diagnostics map[AFEBlock]*BlockDiagnostic) []*SelfTestResult {
	blocks := []AFEBlock{
		BlockElectrochemical, BlockImmunoassay, BlockColorimetry,
		BlockThermal, BlockIsothermalAmp, BlockGas,
		BlockImpedance, BlockRadiation, BlockQuantitativePCR,
	}

	results := make([]*SelfTestResult, 0, len(blocks))
	now := time.Now().UTC()
	for _, b := range blocks {
		diag, ok := diagnostics[b]
		if !ok {
			results = append(results, &SelfTestResult{
				Block: b, Healthy: false, ErrorCode: "no_data", TestedAt: now,
			})
			continue
		}
		results = append(results, &SelfTestResult{
			Block:      b,
			Healthy:    diag.Healthy(),
			NoiseLevel: diag.NoiseLevel,
			BaselineMv: diag.BaselineMv,
			ErrorCode:  diag.ErrorCode,
			TestedAt:   now,
		})
	}
	return results
}

// BlockDiagnostic은 블록의 진단 정보입니다.
type BlockDiagnostic struct {
	NoiseLevel float64 // mV RMS
	BaselineMv float64
	ErrorCode  string
}

// Healthy는 진단이 정상 범위인지 반환합니다.
func (d *BlockDiagnostic) Healthy() bool {
	return d.ErrorCode == "" && d.NoiseLevel < 5.0 // 5 mV 미만
}

// AnyUnhealthy는 결과 목록 중 unhealthy 블록이 있는지 반환합니다.
func AnyUnhealthy(results []*SelfTestResult) bool {
	for _, r := range results {
		if !r.Healthy {
			return true
		}
	}
	return false
}
