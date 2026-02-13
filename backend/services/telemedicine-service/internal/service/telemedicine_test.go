package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/services/telemedicine-service/internal/repository/memory"
	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.TelemedicineService {
	logger := zap.NewNop()
	consultRepo := memory.NewConsultationRepository()
	doctorRepo := memory.NewDoctorRepository()
	sessionRepo := memory.NewVideoSessionRepository()
	return service.NewTelemedicineService(logger, consultRepo, doctorRepo, sessionRepo)
}

func TestCreateConsultation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	c, err := svc.CreateConsultation(ctx, "patient-1", service.SpecialtyInternal, "두통이 심합니다", "3일째 지속되는 두통")
	if err != nil {
		t.Fatalf("상담 생성 실패: %v", err)
	}
	if c.ID == "" {
		t.Fatal("상담 ID가 비어 있음")
	}
	if c.PatientUserID != "patient-1" {
		t.Fatalf("PatientUserID 불일치: got %s, want patient-1", c.PatientUserID)
	}
	if c.Status != service.StatusRequested {
		t.Fatalf("Status 불일치: got %d, want %d", c.Status, service.StatusRequested)
	}
	if c.ChiefComplaint != "두통이 심합니다" {
		t.Fatalf("ChiefComplaint 불일치: got %s", c.ChiefComplaint)
	}
	if c.Specialty != service.SpecialtyInternal {
		t.Fatalf("Specialty 불일치: got %d, want %d", c.Specialty, service.SpecialtyInternal)
	}
}

func TestCreateConsultation_MissingPatientID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreateConsultation(ctx, "", service.SpecialtyGeneral, "두통", "설명")
	if err == nil {
		t.Fatal("빈 patient_user_id에 에러가 반환되어야 함")
	}
}

func TestCreateConsultation_MissingComplaint(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreateConsultation(ctx, "patient-1", service.SpecialtyGeneral, "", "설명")
	if err == nil {
		t.Fatal("빈 chief_complaint에 에러가 반환되어야 함")
	}
}

func TestGetConsultation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateConsultation(ctx, "patient-2", service.SpecialtyCardiology, "가슴 통증", "운동 시 통증")

	got, err := svc.GetConsultation(ctx, created.ID)
	if err != nil {
		t.Fatalf("상담 조회 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("ID 불일치: got %s, want %s", got.ID, created.ID)
	}
	if got.PatientUserID != "patient-2" {
		t.Fatalf("PatientUserID 불일치: got %s", got.PatientUserID)
	}
}

func TestGetConsultation_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetConsultation(ctx, "non-existent-id")
	if err == nil {
		t.Fatal("존재하지 않는 상담에 에러가 반환되어야 함")
	}
}

func TestListConsultations_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.CreateConsultation(ctx, "patient-list", service.SpecialtyGeneral, "증상", "설명")
		if err != nil {
			t.Fatalf("상담 생성 실패: %v", err)
		}
	}

	consultations, total, err := svc.ListConsultations(ctx, "patient-list", service.StatusUnknown, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 개수 불일치: got %d, want 3", total)
	}
	if len(consultations) != 3 {
		t.Fatalf("반환 개수 불일치: got %d, want 3", len(consultations))
	}
}

func TestMatchDoctor_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	doctors, err := svc.MatchDoctor(ctx, service.SpecialtyInternal, "ko")
	if err != nil {
		t.Fatalf("의사 매칭 실패: %v", err)
	}
	if len(doctors) == 0 {
		t.Fatal("내과 의사가 최소 1명은 있어야 함")
	}
	for _, d := range doctors {
		if d.Specialty != service.SpecialtyInternal {
			t.Fatalf("전문 분야 불일치: got %d, want %d", d.Specialty, service.SpecialtyInternal)
		}
		if !d.IsAvailable {
			t.Fatal("가용하지 않은 의사가 반환됨")
		}
	}
}

func TestMatchDoctor_NoResults(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 존재하지 않는 언어로 필터링
	doctors, err := svc.MatchDoctor(ctx, service.SpecialtyInternal, "zh")
	if err != nil {
		t.Fatalf("의사 매칭 실패: %v", err)
	}
	if len(doctors) != 0 {
		t.Fatalf("결과가 없어야 함: got %d", len(doctors))
	}
}

func TestStartVideoSession_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	consultation, _ := svc.CreateConsultation(ctx, "patient-video", service.SpecialtyGeneral, "증상", "설명")

	session, err := svc.StartVideoSession(ctx, consultation.ID, "patient-video")
	if err != nil {
		t.Fatalf("비디오 세션 시작 실패: %v", err)
	}
	if session.ID == "" {
		t.Fatal("세션 ID가 비어 있음")
	}
	if session.ConsultationID != consultation.ID {
		t.Fatalf("ConsultationID 불일치: got %s, want %s", session.ConsultationID, consultation.ID)
	}
	if !strings.HasPrefix(session.RoomURL, "https://meet.manpasik.com/") {
		t.Fatalf("RoomURL 형식 불일치: got %s", session.RoomURL)
	}
	if session.Token == "" {
		t.Fatal("Token이 비어 있음")
	}
	if session.Status != service.SessionConnected {
		t.Fatalf("Status 불일치: got %d, want %d", session.Status, service.SessionConnected)
	}

	// 상담 상태가 InProgress로 변경되었는지 확인
	updated, _ := svc.GetConsultation(ctx, consultation.ID)
	if updated.Status != service.StatusInProgress {
		t.Fatalf("상담 Status 불일치: got %d, want %d", updated.Status, service.StatusInProgress)
	}
}

func TestEndVideoSession_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	consultation, _ := svc.CreateConsultation(ctx, "patient-end", service.SpecialtyGeneral, "증상", "설명")
	session, _ := svc.StartVideoSession(ctx, consultation.ID, "patient-end")

	ended, err := svc.EndVideoSession(ctx, session.ID, consultation.ID, "환자 상태 양호", "경미한 감기")
	if err != nil {
		t.Fatalf("비디오 세션 종료 실패: %v", err)
	}
	if ended.Status != service.SessionEnded {
		t.Fatalf("세션 Status 불일치: got %d, want %d", ended.Status, service.SessionEnded)
	}

	// 상담 상태가 Completed로 변경되었는지 확인
	updated, _ := svc.GetConsultation(ctx, consultation.ID)
	if updated.Status != service.StatusCompleted {
		t.Fatalf("상담 Status 불일치: got %d, want %d", updated.Status, service.StatusCompleted)
	}
	if updated.DoctorNotes != "환자 상태 양호" {
		t.Fatalf("DoctorNotes 불일치: got %s", updated.DoctorNotes)
	}
	if updated.Diagnosis != "경미한 감기" {
		t.Fatalf("Diagnosis 불일치: got %s", updated.Diagnosis)
	}
}

func TestRateConsultation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	consultation, _ := svc.CreateConsultation(ctx, "patient-rate", service.SpecialtyGeneral, "증상", "설명")

	newRating, err := svc.RateConsultation(ctx, consultation.ID, 4.5)
	if err != nil {
		t.Fatalf("평점 등록 실패: %v", err)
	}
	if newRating != 4.5 {
		t.Fatalf("평점 불일치: got %f, want 4.5", newRating)
	}

	// 상담에 평점이 반영되었는지 확인
	updated, _ := svc.GetConsultation(ctx, consultation.ID)
	if updated.Rating != 4.5 {
		t.Fatalf("상담 Rating 불일치: got %f, want 4.5", updated.Rating)
	}
}

func TestEndToEnd_ConsultationFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 1. 상담 생성
	consultation, err := svc.CreateConsultation(ctx, "patient-e2e", service.SpecialtyCardiology, "심장 두근거림", "운동 시 심장이 빨리 뜁니다")
	if err != nil {
		t.Fatalf("상담 생성 실패: %v", err)
	}
	if consultation.Status != service.StatusRequested {
		t.Fatalf("초기 상태 불일치: got %d, want %d", consultation.Status, service.StatusRequested)
	}

	// 2. 의사 매칭
	doctors, err := svc.MatchDoctor(ctx, service.SpecialtyCardiology, "ko")
	if err != nil {
		t.Fatalf("의사 매칭 실패: %v", err)
	}
	if len(doctors) == 0 {
		t.Fatal("심장내과 의사가 최소 1명은 있어야 함")
	}

	// 3. 비디오 세션 시작
	session, err := svc.StartVideoSession(ctx, consultation.ID, "patient-e2e")
	if err != nil {
		t.Fatalf("비디오 세션 시작 실패: %v", err)
	}

	// 4. 상담 상태 확인 (InProgress)
	inProgress, _ := svc.GetConsultation(ctx, consultation.ID)
	if inProgress.Status != service.StatusInProgress {
		t.Fatalf("진행 중 상태 불일치: got %d, want %d", inProgress.Status, service.StatusInProgress)
	}

	// 5. 비디오 세션 종료
	_, err = svc.EndVideoSession(ctx, session.ID, consultation.ID, "심전도 검사 권장", "심방세동 의심")
	if err != nil {
		t.Fatalf("비디오 세션 종료 실패: %v", err)
	}

	// 6. 상담 상태 확인 (Completed)
	completed, _ := svc.GetConsultation(ctx, consultation.ID)
	if completed.Status != service.StatusCompleted {
		t.Fatalf("완료 상태 불일치: got %d, want %d", completed.Status, service.StatusCompleted)
	}
	if completed.Diagnosis != "심방세동 의심" {
		t.Fatalf("진단 불일치: got %s", completed.Diagnosis)
	}

	// 7. 평점 등록
	rating, err := svc.RateConsultation(ctx, consultation.ID, 5.0)
	if err != nil {
		t.Fatalf("평점 등록 실패: %v", err)
	}
	if rating != 5.0 {
		t.Fatalf("평점 불일치: got %f, want 5.0", rating)
	}

	// 8. 목록 조회 확인
	list, total, err := svc.ListConsultations(ctx, "patient-e2e", service.StatusUnknown, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if total != 1 {
		t.Fatalf("총 개수 불일치: got %d, want 1", total)
	}
	if len(list) != 1 {
		t.Fatalf("반환 개수 불일치: got %d, want 1", len(list))
	}
	if list[0].Rating != 5.0 {
		t.Fatalf("목록 내 평점 불일치: got %f, want 5.0", list[0].Rating)
	}
}
