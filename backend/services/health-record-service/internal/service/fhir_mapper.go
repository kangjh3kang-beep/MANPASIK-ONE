package service

import (
	"encoding/json"
	"fmt"
	"time"
)

// LOINCMapping은 LOINC 코드 매핑 정보입니다.
type LOINCMapping struct {
	Code    string
	Display string
	Unit    string
}

// BiomarkerLOINCMap maps internal biomarker names to LOINC codes
var BiomarkerLOINCMap = map[string]LOINCMapping{
	"blood_glucose":             {Code: "15074-8", Display: "Glucose [Moles/volume] in Blood", Unit: "mg/dL"},
	"blood_glucose_fasting":     {Code: "1558-6", Display: "Fasting glucose [Mass/volume] in Serum or Plasma", Unit: "mg/dL"},
	"hemoglobin_a1c":            {Code: "4548-4", Display: "Hemoglobin A1c/Hemoglobin.total in Blood", Unit: "%"},
	"cholesterol_total":         {Code: "2093-3", Display: "Cholesterol [Mass/volume] in Serum or Plasma", Unit: "mg/dL"},
	"cholesterol_hdl":           {Code: "2085-9", Display: "HDL Cholesterol", Unit: "mg/dL"},
	"cholesterol_ldl":           {Code: "2089-1", Display: "LDL Cholesterol", Unit: "mg/dL"},
	"triglycerides":             {Code: "2571-8", Display: "Triglycerides", Unit: "mg/dL"},
	"blood_pressure_systolic":   {Code: "8480-6", Display: "Systolic blood pressure", Unit: "mmHg"},
	"blood_pressure_diastolic":  {Code: "8462-4", Display: "Diastolic blood pressure", Unit: "mmHg"},
	"heart_rate":                {Code: "8867-4", Display: "Heart rate", Unit: "beats/min"},
	"body_temperature":          {Code: "8310-5", Display: "Body temperature", Unit: "Cel"},
	"oxygen_saturation":         {Code: "2708-6", Display: "Oxygen saturation in Arterial blood", Unit: "%"},
	"creatinine":                {Code: "2160-0", Display: "Creatinine [Mass/volume] in Serum or Plasma", Unit: "mg/dL"},
	"uric_acid":                 {Code: "3084-1", Display: "Uric acid [Mass/volume] in Serum or Plasma", Unit: "mg/dL"},
	"cortisol":                  {Code: "2143-6", Display: "Cortisol [Mass/volume] in Serum or Plasma", Unit: "ug/dL"},
}

// MeasurementToFHIRObservation converts a measurement to FHIR R4 Observation
func MeasurementToFHIRObservation(biomarkerName string, value float64, unit string, patientRef string, measuredAt time.Time) map[string]interface{} {
	loinc, exists := BiomarkerLOINCMap[biomarkerName]
	if !exists {
		loinc = LOINCMapping{Code: "unknown", Display: biomarkerName, Unit: unit}
	}

	return map[string]interface{}{
		"resourceType": "Observation",
		"status":       "final",
		"category": []map[string]interface{}{
			{
				"coding": []map[string]interface{}{
					{"system": "http://terminology.hl7.org/CodeSystem/observation-category", "code": "laboratory"},
				},
			},
		},
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://loinc.org", "code": loinc.Code, "display": loinc.Display},
			},
			"text": loinc.Display,
		},
		"subject":           map[string]string{"reference": patientRef},
		"effectiveDateTime": measuredAt.Format(time.RFC3339),
		"valueQuantity": map[string]interface{}{
			"value":  value,
			"unit":   loinc.Unit,
			"system": "http://unitsofmeasure.org",
			"code":   loinc.Unit,
		},
	}
}

// BuildFHIRBundle creates a FHIR R4 Bundle from observations
func BuildFHIRBundle(observations []map[string]interface{}) (string, error) {
	entries := make([]map[string]interface{}, len(observations))
	for i, obs := range observations {
		entries[i] = map[string]interface{}{
			"resource": obs,
		}
	}

	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"total":        len(observations),
		"entry":        entries,
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return "", fmt.Errorf("FHIR bundle JSON 생성 실패: %w", err)
	}
	return string(data), nil
}
