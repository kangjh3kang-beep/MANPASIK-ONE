package poct_test

import (
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/poct"
)

const sampleMessage = `MSH|^~\&|MMUP-DEVICE-001|MFG-FACILITY|MMUP-GATEWAY|HOSPITAL||20260430120000|ORU^R01|MSG-001
PID|||PATIENT-789
OBR|1|ORDER-100||5811-5^Glucose
OBX|1|NM|2345-7^Glucose|1|110.5|mg/dL|70-100|H||F|20260430120030
OBX|2|NM|2093-3^Cholesterol|1|190|mg/dL|0-200|N||F|20260430120030
OBX|3|NM|2160-0^Creatinine|1|0.9|mg/dL|0.6-1.3|N||F|20260430120030`

func TestParser_ParseHappyPath(t *testing.T) {
	p := poct.NewParser()
	msg, err := p.Parse(sampleMessage)
	if err != nil {
		t.Fatalf("Parse 실패: %v", err)
	}

	if msg.MessageID != "MSG-001" {
		t.Errorf("MessageID = %q", msg.MessageID)
	}
	if msg.PatientID != "PATIENT-789" {
		t.Errorf("PatientID = %q", msg.PatientID)
	}
	if msg.OrderID != "ORDER-100" {
		t.Errorf("OrderID = %q", msg.OrderID)
	}
	if msg.MessageType != "ORU^R01" {
		t.Errorf("MessageType = %q", msg.MessageType)
	}
	if len(msg.Observations) != 3 {
		t.Errorf("Observations = %d, want 3", len(msg.Observations))
	}
}

func TestParser_ObservationFields(t *testing.T) {
	p := poct.NewParser()
	msg, _ := p.Parse(sampleMessage)

	glucose := msg.Observations[0]
	if glucose.LOINCCode != "2345-7" {
		t.Errorf("LOINC = %q", glucose.LOINCCode)
	}
	if glucose.Description != "Glucose" {
		t.Errorf("Description = %q", glucose.Description)
	}
	if glucose.Value != 110.5 {
		t.Errorf("Value = %f", glucose.Value)
	}
	if glucose.Unit != "mg/dL" {
		t.Errorf("Unit = %q", glucose.Unit)
	}
	if glucose.RefLow != 70 || glucose.RefHigh != 100 {
		t.Errorf("Ref range = %f-%f", glucose.RefLow, glucose.RefHigh)
	}
	if glucose.Flag != "H" {
		t.Errorf("Flag = %q", glucose.Flag)
	}
}

func TestParser_EmptyMessage(t *testing.T) {
	p := poct.NewParser()
	if _, err := p.Parse(""); err == nil {
		t.Error("빈 메시지가 통과")
	}
}

func TestParser_NoMSH(t *testing.T) {
	p := poct.NewParser()
	if _, err := p.Parse("PID|||x"); err == nil {
		t.Error("MSH 없이 통과")
	}
}

func TestParser_CR_LineEnding(t *testing.T) {
	p := poct.NewParser()
	crMsg := strings.ReplaceAll(sampleMessage, "\n", "\r")
	msg, err := p.Parse(crMsg)
	if err != nil {
		t.Fatalf("CR 종료 메시지 실패: %v", err)
	}
	if len(msg.Observations) != 3 {
		t.Errorf("CR 모드: Observations = %d, want 3", len(msg.Observations))
	}
}

func TestObservation_IsAbnormal_Flag(t *testing.T) {
	cases := []struct {
		flag     string
		expected bool
	}{
		{"H", true}, {"L", true}, {"A", true},
		{"N", false}, {"", false},
	}
	for _, c := range cases {
		obs := &poct.Observation{Flag: c.flag}
		if obs.IsAbnormal() != c.expected {
			t.Errorf("Flag=%q IsAbnormal=%v, want %v", c.flag, obs.IsAbnormal(), c.expected)
		}
	}
}

func TestObservation_IsAbnormal_Range(t *testing.T) {
	high := &poct.Observation{Value: 250, RefLow: 70, RefHigh: 100}
	if !high.IsAbnormal() {
		t.Error("범위 초과가 정상으로 판정")
	}
	low := &poct.Observation{Value: 50, RefLow: 70, RefHigh: 100}
	if !low.IsAbnormal() {
		t.Error("범위 미달이 정상으로 판정")
	}
	normal := &poct.Observation{Value: 85, RefLow: 70, RefHigh: 100}
	if normal.IsAbnormal() {
		t.Error("정상 범위가 비정상으로 판정")
	}
}

func TestDifferentialEngine_DefaultAlpha(t *testing.T) {
	e := poct.NewDifferentialEngine()
	if e.GetAlpha("ch-1") != poct.DefaultAlpha {
		t.Errorf("default alpha = %f, want %f", e.GetAlpha("ch-1"), poct.DefaultAlpha)
	}
}

func TestDifferentialEngine_AlphaClamp(t *testing.T) {
	e := poct.NewDifferentialEngine()

	// 하한 클램프
	got := e.SetAlpha("ch", 0.50)
	if got != poct.AlphaMin {
		t.Errorf("0.50 → %f, want %f", got, poct.AlphaMin)
	}

	// 상한 클램프
	got = e.SetAlpha("ch", 1.50)
	if got != poct.AlphaMax {
		t.Errorf("1.50 → %f, want %f", got, poct.AlphaMax)
	}

	// 범위 내
	got = e.SetAlpha("ch", 1.05)
	if got != 1.05 {
		t.Errorf("1.05 → %f, want 1.05", got)
	}
}

func TestDifferentialEngine_Correct(t *testing.T) {
	e := poct.NewDifferentialEngine()
	// alpha 0.98로 (S=10, R=2) → 10 - 0.98*2 = 8.04
	corrected := e.Correct("ch", 10.0, 2.0)
	if corrected < 8.03 || corrected > 8.05 {
		t.Errorf("corrected = %f, want ~8.04", corrected)
	}
}

func TestDifferentialEngine_CorrectObservation(t *testing.T) {
	e := poct.NewDifferentialEngine()
	obs := &poct.Observation{
		RawSignal: 100.0,
		RefSignal: 5.0,
	}
	e.CorrectObservation(obs, "glucose-ch")

	if obs.Alpha != poct.DefaultAlpha {
		t.Errorf("Alpha = %f", obs.Alpha)
	}
	expected := 100.0 - 0.98*5.0 // 95.1
	if obs.Value < expected-0.01 || obs.Value > expected+0.01 {
		t.Errorf("Value = %f, want %f", obs.Value, expected)
	}
}

func TestComputeCI95(t *testing.T) {
	low, high := poct.ComputeCI95(100.0, 2.0)
	// margin = 1.96 * 2 = 3.92
	if low < 96.07 || low > 96.09 {
		t.Errorf("low = %f, want ~96.08", low)
	}
	if high < 103.91 || high > 103.93 {
		t.Errorf("high = %f, want ~103.92", high)
	}
}

func TestComputeCI95_Zero(t *testing.T) {
	low, high := poct.ComputeCI95(50.0, 0.0)
	if low != 50.0 || high != 50.0 {
		t.Errorf("CI(stdErr=0) = %f-%f, want 50-50", low, high)
	}
}

func TestApplyCI95(t *testing.T) {
	obs := &poct.Observation{Value: 100.0}
	poct.ApplyCI95(obs, 2.0)
	if obs.CILow == 0 || obs.CIHigh == 0 {
		t.Error("CI 미설정")
	}
	if obs.CIHigh-obs.CILow < 7.0 {
		t.Errorf("CI 폭이 너무 좁음: %f", obs.CIHigh-obs.CILow)
	}
}

func TestSampleClassifier_Liquid(t *testing.T) {
	c := poct.NewSampleClassifier()
	if c.ClassifyByImpedance(10_000) != poct.SampleLiquid {
		t.Error("저임피던스 → 액체 분류 실패")
	}
}

func TestSampleClassifier_Gas(t *testing.T) {
	c := poct.NewSampleClassifier()
	if c.ClassifyByImpedance(2_000_000) != poct.SampleGas {
		t.Error("고임피던스 → 가스 분류 실패")
	}
}

func TestSampleClassifier_Solid(t *testing.T) {
	c := poct.NewSampleClassifier()
	if c.ClassifyByImpedance(500_000) != poct.SampleSolid {
		t.Error("중간 임피던스 → 고체 분류 실패")
	}
}

func TestSelfTestRunner_AllHealthy(t *testing.T) {
	r := poct.NewSelfTestRunner()
	diagnostics := map[poct.AFEBlock]*poct.BlockDiagnostic{
		poct.BlockElectrochemical: {NoiseLevel: 1.0, BaselineMv: 0.5},
		poct.BlockImmunoassay:         {NoiseLevel: 2.0, BaselineMv: 1.0},
		poct.BlockColorimetry:     {NoiseLevel: 0.5},
		poct.BlockThermal:         {NoiseLevel: 0.3},
		poct.BlockIsothermalAmp:            {NoiseLevel: 1.5},
		poct.BlockGas:             {NoiseLevel: 0.8},
		poct.BlockImpedance:       {NoiseLevel: 1.2},
		poct.BlockRadiation:       {NoiseLevel: 0.4},
		poct.BlockQuantitativePCR:            {NoiseLevel: 1.8},
	}

	results := r.RunAll(diagnostics)
	if len(results) != 9 {
		t.Errorf("results = %d, want 9", len(results))
	}
	if poct.AnyUnhealthy(results) {
		t.Error("정상 진단인데 unhealthy 검출")
	}
}

func TestSelfTestRunner_OneFailure(t *testing.T) {
	r := poct.NewSelfTestRunner()
	diagnostics := map[poct.AFEBlock]*poct.BlockDiagnostic{
		poct.BlockImmunoassay: {NoiseLevel: 10.0, ErrorCode: "high_noise"}, // 고잡음
	}
	results := r.RunAll(diagnostics)
	if !poct.AnyUnhealthy(results) {
		t.Error("실패 진단이 healthy로 판정")
	}
}

func TestSelfTestRunner_MissingDiagnostics(t *testing.T) {
	r := poct.NewSelfTestRunner()
	results := r.RunAll(map[poct.AFEBlock]*poct.BlockDiagnostic{})

	for _, res := range results {
		if res.Healthy {
			t.Errorf("진단 누락 블록 %s가 healthy", res.Block)
		}
		if res.ErrorCode != "no_data" {
			t.Errorf("ErrorCode = %q, want no_data", res.ErrorCode)
		}
	}
}

func TestParser_TimestampParsing(t *testing.T) {
	p := poct.NewParser()
	msg, err := p.Parse(sampleMessage)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Timestamp.Year() != 2026 {
		t.Errorf("Year = %d", msg.Timestamp.Year())
	}
	if msg.Observations[0].Timestamp.Year() != 2026 {
		t.Error("OBX timestamp 미파싱")
	}
}

func TestParser_OnlyMSH(t *testing.T) {
	p := poct.NewParser()
	minimal := "MSH|^~\\&|S|F|R|F||20260430|ORU^R01|MIN-001"
	msg, err := p.Parse(minimal)
	if err != nil {
		t.Fatalf("MSH only 파싱 실패: %v", err)
	}
	if msg.MessageID != "MIN-001" {
		t.Errorf("MessageID = %q", msg.MessageID)
	}
	if len(msg.Observations) != 0 {
		t.Errorf("Observations = %d, want 0", len(msg.Observations))
	}
}
