package service

import (
	"testing"

	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
)

func newTestService() *EmergencyService {
	repo := memory.NewEmergencyRepository()
	return NewEmergencyService(repo)
}

func TestReportEmergency_Success(t *testing.T) {
	svc := newTestService()

	id, err := svc.ReportEmergency(ReportEmergencyInput{
		UserID:      "user-1",
		Type:        "cardiac",
		Location:    "37.5665,126.9780",
		Description: "Chest pain reported",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty emergency ID")
	}
}

func TestReportEmergency_MissingUserID(t *testing.T) {
	svc := newTestService()

	_, err := svc.ReportEmergency(ReportEmergencyInput{
		Type: "fall",
	})
	if err == nil {
		t.Fatal("expected error for missing user_id")
	}
}

func TestReportEmergency_MissingType(t *testing.T) {
	svc := newTestService()

	_, err := svc.ReportEmergency(ReportEmergencyInput{
		UserID: "user-1",
	})
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestGetEmergencyContacts_Empty(t *testing.T) {
	svc := newTestService()

	contacts, err := svc.GetEmergencyContacts("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 0 {
		t.Fatalf("expected 0 contacts, got %d", len(contacts))
	}
}

func TestGetEmergencyContacts_WithContacts(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := NewEmergencyService(repo)

	_ = repo.AddContact(&memory.EmergencyContact{
		UserID:       "user-1",
		Name:         "Kim",
		Phone:        "010-1234-5678",
		Relationship: "spouse",
	})
	_ = repo.AddContact(&memory.EmergencyContact{
		UserID:       "user-1",
		Name:         "Lee",
		Phone:        "010-8765-4321",
		Relationship: "parent",
	})

	contacts, err := svc.GetEmergencyContacts("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}
}

func TestUpdateAndGetEmergencySettings(t *testing.T) {
	svc := newTestService()

	settings := &memory.EmergencySettings{
		UserID:              "user-1",
		AutoCall119:         true,
		EmergencyContactIDs: []string{"c-1", "c-2"},
		MedicalInfo:         "Allergic to penicillin",
	}

	err := svc.UpdateEmergencySettings(settings)
	if err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	got, err := svc.GetEmergencySettings("user-1")
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if !got.AutoCall119 {
		t.Fatal("expected AutoCall119 to be true")
	}
	if len(got.EmergencyContactIDs) != 2 {
		t.Fatalf("expected 2 contact IDs, got %d", len(got.EmergencyContactIDs))
	}
	if got.MedicalInfo != "Allergic to penicillin" {
		t.Fatalf("unexpected MedicalInfo: %s", got.MedicalInfo)
	}
}

func TestGetEmergencySettings_Default(t *testing.T) {
	svc := newTestService()

	got, err := svc.GetEmergencySettings("user-new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.AutoCall119 {
		t.Fatal("expected AutoCall119 to default to false")
	}
	if got.UserID != "user-new" {
		t.Fatalf("expected UserID user-new, got %s", got.UserID)
	}
}

func TestReportEmergency_AttachesContacts(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := NewEmergencyService(repo)

	_ = repo.AddContact(&memory.EmergencyContact{
		UserID:       "user-1",
		Name:         "Kim",
		Phone:        "010-1234-5678",
		Relationship: "spouse",
	})

	id, err := svc.ReportEmergency(ReportEmergencyInput{
		UserID:      "user-1",
		Type:        "fall",
		Location:    "37.5665,126.9780",
		Description: "Fall detected",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	emg, err := repo.GetEmergency(id)
	if err != nil {
		t.Fatalf("unexpected error getting emergency: %v", err)
	}
	if len(emg.ContactIDs) != 1 {
		t.Fatalf("expected 1 contact ID attached, got %d", len(emg.ContactIDs))
	}
}
