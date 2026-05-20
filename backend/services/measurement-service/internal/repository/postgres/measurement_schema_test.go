package postgres

import (
	"os"
	"strings"
	"testing"
)

func TestMeasurementSchemaIncludesEvidenceColumns(t *testing.T) {
	schema, err := os.ReadFile("../../../../../../infrastructure/database/init/04-measurement.sql")
	if err != nil {
		t.Fatalf("measurement schema read failed: %v", err)
	}
	text := string(schema)
	for _, column := range []string{"evidence_status", "diagnostic_ready", "evidence_gaps"} {
		if !strings.Contains(text, column) {
			t.Fatalf("measurement schema missing %s", column)
		}
	}
	if !strings.Contains(text, "md.evidence_status") {
		t.Fatal("measurement_summary view must expose md.evidence_status")
	}
}
