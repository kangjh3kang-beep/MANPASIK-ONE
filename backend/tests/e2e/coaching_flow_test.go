package e2e

import (
	"context"
	"testing"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestCoachingFlow E2E: SetHealthGoal → GetHealthGoals → GenerateCoaching → ListCoachingMessages
func TestCoachingFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, CoachingAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("coaching 서비스 연결 불가: %v", err)
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 1) SetHealthGoal — 건강 목표 설정
	goalResp, err := client.SetHealthGoal(rpcCtx, &v1.SetHealthGoalRequest{
		UserId:      "e2e-user-coach-001",
		Category:    v1.GoalCategory_GOAL_CATEGORY_BLOOD_GLUCOSE,
		MetricName:  "fasting_glucose",
		TargetValue: 100.0,
		Unit:        "mg/dL",
		Description: "공복 혈당 100 이하 유지",
	})
	if err != nil {
		t.Fatalf("SetHealthGoal 실패: %v", err)
	}
	if goalResp.GoalId == "" {
		t.Fatal("goal_id가 비어 있습니다")
	}
	t.Logf("SetHealthGoal 성공: id=%s, category=%v, target=%.1f %s",
		goalResp.GoalId, goalResp.Category, goalResp.TargetValue, goalResp.Unit)

	// 2) GetHealthGoals — 목표 목록 조회
	goalsResp, err := client.GetHealthGoals(rpcCtx, &v1.GetHealthGoalsRequest{
		UserId: "e2e-user-coach-001",
	})
	if err != nil {
		t.Fatalf("GetHealthGoals 실패: %v", err)
	}
	if len(goalsResp.Goals) == 0 {
		t.Error("목표 목록이 비어 있습니다")
	}
	t.Logf("GetHealthGoals 성공: %d개 목표", len(goalsResp.Goals))

	// 3) GenerateCoaching — AI 코칭 메시지 생성
	coachingResp, err := client.GenerateCoaching(rpcCtx, &v1.GenerateCoachingRequest{
		UserId:      "e2e-user-coach-001",
		CoachingType: v1.CoachingType_COACHING_TYPE_DAILY_TIP,
	})
	if err != nil {
		t.Logf("GenerateCoaching 실패 (AI 미연동 시 예상): %v", err)
	} else {
		if coachingResp.MessageId == "" {
			t.Error("message_id가 비어 있습니다")
		}
		t.Logf("GenerateCoaching 성공: id=%s, type=%v, title=%s",
			coachingResp.MessageId, coachingResp.CoachingType, coachingResp.Title)
	}

	// 4) ListCoachingMessages — 코칭 메시지 이력 조회
	msgsResp, err := client.ListCoachingMessages(rpcCtx, &v1.ListCoachingMessagesRequest{
		UserId: "e2e-user-coach-001",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("ListCoachingMessages 실패: %v", err)
	}
	t.Logf("ListCoachingMessages 성공: total=%d", msgsResp.TotalCount)

	// 5) GetRecommendations — 개인화 추천 조회
	recsResp, err := client.GetRecommendations(rpcCtx, &v1.GetRecommendationsRequest{
		UserId: "e2e-user-coach-001",
		Limit:  5,
	})
	if err != nil {
		t.Logf("GetRecommendations 실패 (AI 미연동 시 예상): %v", err)
	} else {
		t.Logf("GetRecommendations 성공: %d개 추천", len(recsResp.Recommendations))
	}

	t.Logf("✅ 코칭 플로우 완료: goal_id=%s", goalResp.GoalId)
}
