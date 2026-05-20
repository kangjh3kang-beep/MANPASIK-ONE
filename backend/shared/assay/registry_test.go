package assay

import "testing"

func TestResolveKnownAssayAliases(t *testing.T) {
	tests := []struct {
		name          string
		cartridgeType string
		wantKey       string
		wantAnalyte   string
		wantUnit      string
	}{
		{
			name:          "glucose direct",
			cartridgeType: "glucose",
			wantKey:       "glucose",
			wantAnalyte:   "glucose",
			wantUnit:      "mg/dL",
		},
		{
			name:          "glucose blood alias",
			cartridgeType: "blood_glucose",
			wantKey:       "glucose",
			wantAnalyte:   "glucose",
			wantUnit:      "mg/dL",
		},
		{
			name:          "glucose cartridge prefix",
			cartridgeType: "cartridge-glucose",
			wantKey:       "glucose",
			wantAnalyte:   "glucose",
			wantUnit:      "mg/dL",
		},
		{
			name:          "crp direct",
			cartridgeType: "crp",
			wantKey:       "crp",
			wantAnalyte:   "c_reactive_protein",
			wantUnit:      "mg/L",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.cartridgeType)
			if err != nil {
				t.Fatalf("Resolve(%q) returned error: %v", tt.cartridgeType, err)
			}
			if got.Key != tt.wantKey || got.Analyte != tt.wantAnalyte || got.Unit != tt.wantUnit {
				t.Fatalf("definition mismatch: key=%q analyte=%q unit=%q", got.Key, got.Analyte, got.Unit)
			}
		})
	}
}

func TestResolveUnknownAssay(t *testing.T) {
	if _, err := Resolve("custom-research-unlocked"); err == nil {
		t.Fatal("Resolve returned nil error for unknown cartridge")
	}
}

func TestEvaluateUsesCorrectedSignalAndCompletenessConfidence(t *testing.T) {
	result, err := Evaluate("glucose", Signal{
		SCorrected:      123.4,
		RawChannelCount: 88,
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if result.PrimaryValue != 123.4 {
		t.Fatalf("PrimaryValue = %f, want 123.4", result.PrimaryValue)
	}
	if result.Unit != "mg/dL" {
		t.Fatalf("Unit = %q, want mg/dL", result.Unit)
	}
	if result.Confidence < 0.90 || result.Confidence > 0.98 {
		t.Fatalf("Confidence = %f, want bounded assay confidence", result.Confidence)
	}

	lowCompleteness, err := Evaluate("glucose", Signal{
		SCorrected:      123.4,
		RawChannelCount: 3,
	})
	if err != nil {
		t.Fatalf("Evaluate returned error for low completeness: %v", err)
	}
	if lowCompleteness.Confidence >= result.Confidence {
		t.Fatalf("low completeness confidence = %f, want below full confidence %f", lowCompleteness.Confidence, result.Confidence)
	}
}

func TestResolveIncludesEvidenceManifestWithoutDiagnosticClaim(t *testing.T) {
	definition, err := Resolve("glucose")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if definition.Evidence.LOINCCode != "15074-8" {
		t.Fatalf("LOINCCode = %q, want 15074-8", definition.Evidence.LOINCCode)
	}
	if definition.Evidence.UCUMUnit != "mg/dL" {
		t.Fatalf("UCUMUnit = %q, want mg/dL", definition.Evidence.UCUMUnit)
	}
	if definition.Evidence.Status != EvidenceStatusResearchOnly {
		t.Fatalf("Evidence status = %q, want %q", definition.Evidence.Status, EvidenceStatusResearchOnly)
	}
	if definition.IsDiagnosticReady() {
		t.Fatal("research-only assay must not be diagnostic ready")
	}
	if gaps := definition.EvidenceGaps(); len(gaps) == 0 {
		t.Fatal("research-only assay must expose evidence gaps")
	}
}

func TestClinicalLockedDefinitionRequiresCompleteEvidence(t *testing.T) {
	definition := Definition{
		Key:                "validated_glucose",
		Analyte:            "glucose",
		Unit:               "mg/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.92,
		ConfidenceCeiling:  0.99,
		Evidence: EvidenceManifest{
			Status:          EvidenceStatusClinicalLocked,
			LOINCCode:       "15074-8",
			UCUMUnit:        "mg/dL",
			ReferenceMethod: "hexokinase_lab_comparator",
			AnalyticalPerformance: AnalyticalPerformance{
				LimitOfBlank:        2.0,
				LimitOfDetection:    5.0,
				LimitOfQuantitation: 10.0,
				LinearityMin:        20.0,
				LinearityMax:        600.0,
			},
			AcceptanceCriteria: AcceptanceCriteria{
				MaxBiasPercent:  10.0,
				MaxCVPercent:    5.0,
				MinSensitivity:  0.90,
				MinSpecificity:  0.90,
				MinSampleCount:  120,
				ComparatorStudy: "prospective_method_comparison",
			},
		},
	}
	if gaps := definition.EvidenceGaps(); len(gaps) != 0 {
		t.Fatalf("complete clinical evidence gaps = %v, want none", gaps)
	}
	if !definition.IsDiagnosticReady() {
		t.Fatal("complete clinical locked definition must be diagnostic ready")
	}

	definition.Evidence.ReferenceMethod = ""
	if definition.IsDiagnosticReady() {
		t.Fatal("clinical locked definition without reference method must not be diagnostic ready")
	}
	if gaps := definition.EvidenceGaps(); len(gaps) == 0 {
		t.Fatal("missing reference method must be reported as evidence gap")
	}
}

func TestClinicalLockedDefinitionRequiresClinicalAcceptanceThresholds(t *testing.T) {
	definition := Definition{
		Key:                "validated_crp",
		Analyte:            "c_reactive_protein",
		Unit:               "mg/L",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.92,
		ConfidenceCeiling:  0.99,
		Evidence: EvidenceManifest{
			Status:          EvidenceStatusClinicalLocked,
			LOINCCode:       "1988-5",
			UCUMUnit:        "mg/L",
			ReferenceMethod: "immunoturbidimetric_crp_comparator",
			AnalyticalPerformance: AnalyticalPerformance{
				LimitOfDetection:    0.1,
				LimitOfQuantitation: 0.3,
				LinearityMin:        0.3,
				LinearityMax:        200.0,
			},
			AcceptanceCriteria: AcceptanceCriteria{
				MaxBiasPercent:  10.0,
				MaxCVPercent:    5.0,
				MinSampleCount:  120,
				ComparatorStudy: "prospective_method_comparison",
			},
		},
	}
	if definition.IsDiagnosticReady() {
		t.Fatal("clinical locked definition without sensitivity/specificity thresholds must not be diagnostic ready")
	}
	gaps := definition.EvidenceGaps()
	assertContainsGap(t, gaps, "min_sensitivity_required")
	assertContainsGap(t, gaps, "min_specificity_required")
}

func assertContainsGap(t *testing.T, gaps []string, want string) {
	t.Helper()
	for _, gap := range gaps {
		if gap == want {
			return
		}
	}
	t.Fatalf("gaps = %v, want %q", gaps, want)
}

func TestEvidenceGateReportsClinicalLockedDefinitionsWithGaps(t *testing.T) {
	definitions := map[string]Definition{
		"bad_clinical_lock": {
			Key:                "bad_clinical_lock",
			Analyte:            "glucose",
			Unit:               "mg/dL",
			ExpectedChannels:   88,
			PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
			ConfidenceFloor:    0.92,
			ConfidenceCeiling:  0.99,
			Evidence: EvidenceManifest{
				Status:    EvidenceStatusClinicalLocked,
				LOINCCode: "15074-8",
				UCUMUnit:  "mg/dL",
			},
		},
		"research_only": {
			Key:                "research_only",
			Analyte:            "glucose",
			Unit:               "mg/dL",
			ExpectedChannels:   88,
			PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
			ConfidenceFloor:    0.90,
			ConfidenceCeiling:  0.98,
			Evidence: EvidenceManifest{
				Status:    EvidenceStatusResearchOnly,
				LOINCCode: "15074-8",
				UCUMUnit:  "mg/dL",
			},
		},
	}

	gaps := ClinicalLockEvidenceGaps(definitions)
	if len(gaps) != 1 {
		t.Fatalf("clinical lock gaps count = %d, want 1: %v", len(gaps), gaps)
	}
	assertContainsGap(t, gaps["bad_clinical_lock"], "reference_method_required")
	if _, ok := gaps["research_only"]; ok {
		t.Fatalf("research-only assay must not fail clinical lock gate: %v", gaps)
	}
}

func TestEvidenceGateCurrentRegistryHasNoBrokenClinicalLocks(t *testing.T) {
	if gaps := ClinicalLockEvidenceGaps(definitions); len(gaps) != 0 {
		t.Fatalf("current registry has broken clinical locks: %v", gaps)
	}
}
