package assay

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var ErrUnknownAssay = errors.New("unknown assay cartridge type")

type PrimaryValueSource string

const (
	PrimaryValueSourceCorrectedSignal PrimaryValueSource = "s_corrected"
)

type EvidenceStatus string

const (
	EvidenceStatusResearchOnly     EvidenceStatus = "research_only"
	EvidenceStatusAnalyticalLocked EvidenceStatus = "analytical_locked"
	EvidenceStatusClinicalLocked   EvidenceStatus = "clinical_locked"
)

type AnalyticalPerformance struct {
	LimitOfBlank        float64
	LimitOfDetection    float64
	LimitOfQuantitation float64
	LinearityMin        float64
	LinearityMax        float64
}

type AcceptanceCriteria struct {
	MaxBiasPercent  float64
	MaxCVPercent    float64
	MinSensitivity  float64
	MinSpecificity  float64
	MinSampleCount  int
	ComparatorStudy string
}

type EvidenceManifest struct {
	Status                EvidenceStatus
	LOINCCode             string
	UCUMUnit              string
	ReferenceMethod       string
	RequiredReference     string
	AnalyticalPerformance AnalyticalPerformance
	AcceptanceCriteria    AcceptanceCriteria
}

type Definition struct {
	Key                string
	Analyte            string
	Unit               string
	ExpectedChannels   int
	PrimaryValueSource PrimaryValueSource
	ConfidenceFloor    float64
	ConfidenceCeiling  float64
	Evidence           EvidenceManifest
}

type Signal struct {
	SCorrected      float64
	RawChannelCount int
}

type Result struct {
	Definition   Definition
	PrimaryValue float64
	Unit         string
	Confidence   float64
}

var definitions = map[string]Definition{
	"glucose": {
		Key:                "glucose",
		Analyte:            "glucose",
		Unit:               "mg/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.90,
		ConfidenceCeiling:  0.98,
		Evidence:           researchEvidence("15074-8", "mg/dL", "hexokinase_lab_comparator"),
	},
	"lipid_panel": {
		Key:                "lipid_panel",
		Analyte:            "lipid_panel",
		Unit:               "mg/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.88,
		ConfidenceCeiling:  0.96,
		Evidence:           researchEvidence("24331-1", "mg/dL", "enzymatic_lipid_panel_comparator"),
	},
	"hba1c": {
		Key:                "hba1c",
		Analyte:            "hemoglobin_a1c",
		Unit:               "%",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.88,
		ConfidenceCeiling:  0.96,
		Evidence:           researchEvidence("4548-4", "%", "hplc_or_ifcc_traceable_comparator"),
	},
	"uric_acid": {
		Key:                "uric_acid",
		Analyte:            "uric_acid",
		Unit:               "mg/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.88,
		ConfidenceCeiling:  0.96,
		Evidence:           researchEvidence("3084-1", "mg/dL", "uricase_lab_comparator"),
	},
	"creatinine": {
		Key:                "creatinine",
		Analyte:            "creatinine",
		Unit:               "mg/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.88,
		ConfidenceCeiling:  0.96,
		Evidence:           researchEvidence("2160-0", "mg/dL", "idms_traceable_creatinine_comparator"),
	},
	"vitamin_d": {
		Key:                "vitamin_d",
		Analyte:            "25_hydroxy_vitamin_d",
		Unit:               "ng/mL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("62292-8", "ng/mL", "lc_ms_ms_vitamin_d_comparator"),
	},
	"vitamin_b12": {
		Key:                "vitamin_b12",
		Analyte:            "vitamin_b12",
		Unit:               "pg/mL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("2132-9", "pg/mL", "immunoassay_b12_comparator"),
	},
	"ferritin": {
		Key:                "ferritin",
		Analyte:            "ferritin",
		Unit:               "ng/mL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("2276-4", "ng/mL", "immunoassay_ferritin_comparator"),
	},
	"tsh": {
		Key:                "tsh",
		Analyte:            "thyroid_stimulating_hormone",
		Unit:               "mIU/L",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("3016-3", "mIU/L", "immunoassay_tsh_comparator"),
	},
	"cortisol": {
		Key:                "cortisol",
		Analyte:            "cortisol",
		Unit:               "ug/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("2143-6", "ug/dL", "immunoassay_or_lc_ms_cortisol_comparator"),
	},
	"testosterone": {
		Key:                "testosterone",
		Analyte:            "testosterone",
		Unit:               "ng/dL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("2986-8", "ng/dL", "lc_ms_ms_testosterone_comparator"),
	},
	"estrogen": {
		Key:                "estrogen",
		Analyte:            "estrogen",
		Unit:               "pg/mL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("2243-4", "pg/mL", "estradiol_lab_comparator"),
	},
	"crp": {
		Key:                "crp",
		Analyte:            "c_reactive_protein",
		Unit:               "mg/L",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.88,
		ConfidenceCeiling:  0.96,
		Evidence:           researchEvidence("1988-5", "mg/L", "immunoturbidimetric_crp_comparator"),
	},
	"insulin": {
		Key:                "insulin",
		Analyte:            "insulin",
		Unit:               "uIU/mL",
		ExpectedChannels:   88,
		PrimaryValueSource: PrimaryValueSourceCorrectedSignal,
		ConfidenceFloor:    0.86,
		ConfidenceCeiling:  0.94,
		Evidence:           researchEvidence("20448-7", "uIU/mL", "immunoassay_insulin_comparator"),
	},
}

var aliases = map[string]string{
	"0x01":                  "glucose",
	"blood_glucose":         "glucose",
	"cartridge_glucose":     "glucose",
	"cartridge-glucose":     "glucose",
	"glucose":               "glucose",
	"0x02":                  "lipid_panel",
	"lipid":                 "lipid_panel",
	"lipid_panel":           "lipid_panel",
	"cartridge_lipid_panel": "lipid_panel",
	"0x03":                  "hba1c",
	"hemoglobin_a1c":        "hba1c",
	"hba1c":                 "hba1c",
	"0x04":                  "uric_acid",
	"uric_acid":             "uric_acid",
	"0x05":                  "creatinine",
	"creatinine":            "creatinine",
	"0x06":                  "vitamin_d",
	"vitamind":              "vitamin_d",
	"vitamin_d":             "vitamin_d",
	"0x07":                  "vitamin_b12",
	"vitaminb12":            "vitamin_b12",
	"vitamin_b12":           "vitamin_b12",
	"0x08":                  "ferritin",
	"ferritin":              "ferritin",
	"0x09":                  "tsh",
	"tsh":                   "tsh",
	"0x0a":                  "cortisol",
	"cortisol":              "cortisol",
	"0x0b":                  "testosterone",
	"testosterone":          "testosterone",
	"0x0c":                  "estrogen",
	"estrogen":              "estrogen",
	"0x0d":                  "crp",
	"c_reactive_protein":    "crp",
	"cartridge_crp":         "crp",
	"cartridge-crp":         "crp",
	"crp":                   "crp",
	"0x0e":                  "insulin",
	"insulin":               "insulin",
}

func Resolve(cartridgeType string) (Definition, error) {
	key := normalize(cartridgeType)
	canonical, ok := aliases[key]
	if !ok {
		return Definition{}, fmt.Errorf("%w: %s", ErrUnknownAssay, cartridgeType)
	}
	definition, ok := definitions[canonical]
	if !ok {
		return Definition{}, fmt.Errorf("%w: %s", ErrUnknownAssay, cartridgeType)
	}
	return definition, nil
}

func Evaluate(cartridgeType string, signal Signal) (Result, error) {
	definition, err := Resolve(cartridgeType)
	if err != nil {
		return Result{}, err
	}
	return definition.Evaluate(signal), nil
}

func (d Definition) Evaluate(signal Signal) Result {
	return Result{
		Definition:   d,
		PrimaryValue: signal.SCorrected,
		Unit:         d.Unit,
		Confidence:   d.confidence(signal.RawChannelCount),
	}
}

func (d Definition) confidence(rawChannelCount int) float64 {
	if d.ExpectedChannels <= 0 || rawChannelCount <= 0 {
		return d.ConfidenceFloor
	}
	completeness := float64(rawChannelCount) / float64(d.ExpectedChannels)
	completeness = math.Max(0, math.Min(1, completeness))
	return d.ConfidenceFloor + completeness*(d.ConfidenceCeiling-d.ConfidenceFloor)
}

func (d Definition) IsDiagnosticReady() bool {
	return len(d.EvidenceGaps()) == 0
}

func (d Definition) EvidenceGaps() []string {
	var gaps []string
	if d.Evidence.Status != EvidenceStatusClinicalLocked {
		gaps = append(gaps, "clinical_lock_required")
	}
	if strings.TrimSpace(d.Evidence.LOINCCode) == "" {
		gaps = append(gaps, "loinc_code_required")
	}
	if strings.TrimSpace(d.Evidence.UCUMUnit) == "" {
		gaps = append(gaps, "ucum_unit_required")
	}
	if strings.TrimSpace(d.Evidence.ReferenceMethod) == "" {
		gaps = append(gaps, "reference_method_required")
	}
	if d.Evidence.AnalyticalPerformance.LimitOfDetection <= 0 {
		gaps = append(gaps, "limit_of_detection_required")
	}
	if d.Evidence.AnalyticalPerformance.LimitOfQuantitation <= 0 {
		gaps = append(gaps, "limit_of_quantitation_required")
	}
	if d.Evidence.AnalyticalPerformance.LinearityMax <= d.Evidence.AnalyticalPerformance.LinearityMin {
		gaps = append(gaps, "linearity_range_required")
	}
	if d.Evidence.AcceptanceCriteria.MaxBiasPercent <= 0 {
		gaps = append(gaps, "max_bias_acceptance_required")
	}
	if d.Evidence.AcceptanceCriteria.MaxCVPercent <= 0 {
		gaps = append(gaps, "max_cv_acceptance_required")
	}
	if d.Evidence.AcceptanceCriteria.MinSensitivity <= 0 {
		gaps = append(gaps, "min_sensitivity_required")
	}
	if d.Evidence.AcceptanceCriteria.MinSpecificity <= 0 {
		gaps = append(gaps, "min_specificity_required")
	}
	if d.Evidence.AcceptanceCriteria.MinSampleCount <= 0 {
		gaps = append(gaps, "method_comparison_sample_count_required")
	}
	if strings.TrimSpace(d.Evidence.AcceptanceCriteria.ComparatorStudy) == "" {
		gaps = append(gaps, "comparator_study_required")
	}
	return gaps
}

func ClinicalLockEvidenceGaps(definitions map[string]Definition) map[string][]string {
	gapsByKey := make(map[string][]string)
	for key, definition := range definitions {
		if definition.Evidence.Status != EvidenceStatusClinicalLocked {
			continue
		}
		if gaps := definition.EvidenceGaps(); len(gaps) > 0 {
			gapsByKey[key] = gaps
		}
	}
	return gapsByKey
}

func normalize(cartridgeType string) string {
	value := strings.TrimSpace(strings.ToLower(cartridgeType))
	value = strings.ReplaceAll(value, "-", "_")
	return value
}

func researchEvidence(loincCode, ucumUnit, requiredReference string) EvidenceManifest {
	return EvidenceManifest{
		Status:            EvidenceStatusResearchOnly,
		LOINCCode:         loincCode,
		UCUMUnit:          ucumUnit,
		RequiredReference: requiredReference,
	}
}
