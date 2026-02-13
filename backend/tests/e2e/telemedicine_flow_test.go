package e2e

import (
	"context"
	"testing"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestTelemedicineFlow E2E: CreateConsultation → GetConsultation → ListConsultations
func TestTelemedicineFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, TelemedicineAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("telemedicine 서비스 연결 불가: %v", err)
	}
	defer conn.Close()

	client := v1.NewTelemedicineServiceClient(conn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 1) CreateConsultation
	createResp, err := client.CreateConsultation(rpcCtx, &v1.CreateConsultationRequest{
		PatientUserId:  "e2e-user-telemed-001",
		Specialty:      v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL,
		ChiefComplaint: "두통이 지속됩니다",
		Description:    "3일 전부터 두통이 계속되고 있습니다",
	})
	if err != nil {
		t.Fatalf("CreateConsultation 실패: %v", err)
	}
	if createResp.ConsultationId == "" {
		t.Fatal("consultation_id가 비어 있습니다")
	}
	t.Logf("CreateConsultation 성공: id=%s, status=%v", createResp.ConsultationId, createResp.Status)

	// 2) GetConsultation
	getResp, err := client.GetConsultation(rpcCtx, &v1.GetConsultationRequest{
		ConsultationId: createResp.ConsultationId,
	})
	if err != nil {
		t.Fatalf("GetConsultation 실패: %v", err)
	}
	if getResp.ConsultationId != createResp.ConsultationId {
		t.Errorf("ID 불일치: got %s, want %s", getResp.ConsultationId, createResp.ConsultationId)
	}
	t.Logf("GetConsultation 성공: id=%s, complaint=%s, status=%v",
		getResp.ConsultationId, getResp.ChiefComplaint, getResp.Status)

	// 3) ListConsultations
	listResp, err := client.ListConsultations(rpcCtx, &v1.ListConsultationsRequest{
		UserId: "e2e-user-telemed-001",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("ListConsultations 실패: %v", err)
	}
	t.Logf("ListConsultations 성공: total=%d", listResp.TotalCount)

	// 4) MatchDoctor (의사 매칭)
	matchResp, err := client.MatchDoctor(rpcCtx, &v1.MatchDoctorRequest{
		Specialty: v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL,
		Language:  "ko",
	})
	if err != nil {
		t.Logf("MatchDoctor 실패 (의사 미등록 시 예상): %v", err)
	} else {
		t.Logf("MatchDoctor 성공: available=%d doctors", matchResp.TotalAvailable)
	}

	t.Logf("✅ 원격진료 플로우 완료: consultation_id=%s", createResp.ConsultationId)
}
