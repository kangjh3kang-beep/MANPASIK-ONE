package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/video-service/internal/repository/memory"
	"github.com/manpasik/backend/services/video-service/internal/service"
	"go.uber.org/zap"
)

func setupVideoService() *service.VideoService {
	logger := zap.NewNop()
	roomRepo := memory.NewRoomRepository()
	signalRepo := memory.NewSignalRepository()
	return service.NewVideoService(logger, roomRepo, signalRepo)
}

func TestCreateRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, err := svc.CreateRoom(ctx, "테스트 회의실", service.RoomTypeOneToOne, "user-1", 2)
	if err != nil {
		t.Fatalf("CreateRoom 실패: %v", err)
	}
	if room.ID == "" {
		t.Fatal("회의실 ID가 비어 있습니다")
	}
	if room.Status != service.RoomStatusWaiting {
		t.Fatalf("예상 상태 Waiting, 실제: %d", room.Status)
	}
	if room.MaxParticipants != 2 {
		t.Fatalf("최대 참가자 수 예상 2, 실제: %d", room.MaxParticipants)
	}
}

func TestCreateRoom_EmptyCreator(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	_, err := svc.CreateRoom(ctx, "", service.RoomTypeUnknown, "", 0)
	if err == nil {
		t.Fatal("빈 created_by에 에러가 발생해야 합니다")
	}
}

func TestCreateRoom_DefaultName(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, err := svc.CreateRoom(ctx, "", service.RoomTypeGroup, "user-1", 0)
	if err != nil {
		t.Fatalf("CreateRoom 실패: %v", err)
	}
	if room.Name == "" {
		t.Fatal("기본 이름이 생성되어야 합니다")
	}
	if room.MaxParticipants != 2 {
		t.Fatalf("기본 최대 참가자 수 예상 2, 실제: %d", room.MaxParticipants)
	}
}

func TestGetRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	created, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)

	got, err := svc.GetRoom(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRoom 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("ID 불일치: %s != %s", got.ID, created.ID)
	}
}

func TestGetRoom_NotFound(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	_, err := svc.GetRoom(ctx, "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 회의실에 에러가 발생해야 합니다")
	}
}

func TestJoinRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 4)

	updated, token, iceServers, err := svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "participant")
	if err != nil {
		t.Fatalf("JoinRoom 실패: %v", err)
	}
	if token == "" {
		t.Fatal("토큰이 비어 있습니다")
	}
	if iceServers == "" {
		t.Fatal("ICE 서버 정보가 비어 있습니다")
	}
	if len(updated.Participants) != 1 {
		t.Fatalf("참가자 수 예상 1, 실제: %d", len(updated.Participants))
	}
	if updated.Status != service.RoomStatusActive {
		t.Fatalf("첫 참가 시 상태가 Active여야 합니다, 실제: %d", updated.Status)
	}
}

func TestJoinRoom_MaxParticipants(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 1)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")

	_, _, _, err := svc.JoinRoom(ctx, room.ID, "user-3", "사용자3", "")
	if err == nil {
		t.Fatal("최대 인원 초과 시 에러가 발생해야 합니다")
	}
}

func TestLeaveRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeGroup, "user-1", 4)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")
	svc.JoinRoom(ctx, room.ID, "user-3", "사용자3", "")

	remaining, err := svc.LeaveRoom(ctx, room.ID, "user-2")
	if err != nil {
		t.Fatalf("LeaveRoom 실패: %v", err)
	}
	if remaining != 1 {
		t.Fatalf("남은 참가자 수 예상 1, 실제: %d", remaining)
	}
}

func TestLeaveRoom_NotFound(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)

	_, err := svc.LeaveRoom(ctx, room.ID, "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 참가자 퇴장에 에러가 발생해야 합니다")
	}
}

func TestEndRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")

	ended, err := svc.EndRoom(ctx, room.ID, "user-1")
	if err != nil {
		t.Fatalf("EndRoom 실패: %v", err)
	}
	if ended.Status != service.RoomStatusEnded {
		t.Fatalf("예상 상태 Ended, 실제: %d", ended.Status)
	}

	// 종료된 방에 참가 시도
	_, _, _, err = svc.JoinRoom(ctx, room.ID, "user-3", "사용자3", "")
	if err == nil {
		t.Fatal("종료된 방에 참가 시도 시 에러가 발생해야 합니다")
	}
}

func TestSendSignal(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")

	signalID, err := svc.SendSignal(ctx, room.ID, "user-2", "user-1", service.SignalTypeOffer, `{"sdp":"..."}`)
	if err != nil {
		t.Fatalf("SendSignal 실패: %v", err)
	}
	if signalID == "" {
		t.Fatal("시그널 ID가 비어 있습니다")
	}
}

func TestSendSignal_EndedRoom(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)
	svc.EndRoom(ctx, room.ID, "user-1")

	_, err := svc.SendSignal(ctx, room.ID, "user-1", "user-2", service.SignalTypeOffer, "")
	if err == nil {
		t.Fatal("종료된 방에서 시그널 전송 시 에러가 발생해야 합니다")
	}
}

func TestListParticipants(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeGroup, "user-1", 4)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")
	svc.JoinRoom(ctx, room.ID, "user-3", "사용자3", "")
	svc.LeaveRoom(ctx, room.ID, "user-2")

	participants, err := svc.ListParticipants(ctx, room.ID)
	if err != nil {
		t.Fatalf("ListParticipants 실패: %v", err)
	}
	if len(participants) != 1 {
		t.Fatalf("활성 참가자 수 예상 1, 실제: %d", len(participants))
	}
}

func TestGetRoomStats(t *testing.T) {
	svc := setupVideoService()
	ctx := context.Background()

	room, _ := svc.CreateRoom(ctx, "테스트", service.RoomTypeOneToOne, "user-1", 2)
	svc.JoinRoom(ctx, room.ID, "user-2", "사용자2", "")
	svc.SendSignal(ctx, room.ID, "user-2", "user-1", service.SignalTypeOffer, `{"sdp":"offer"}`)
	svc.SendSignal(ctx, room.ID, "user-1", "user-2", service.SignalTypeAnswer, `{"sdp":"answer"}`)

	stats, err := svc.GetRoomStats(ctx, room.ID)
	if err != nil {
		t.Fatalf("GetRoomStats 실패: %v", err)
	}
	if stats.SignalCount != 2 {
		t.Fatalf("시그널 수 예상 2, 실제: %d", stats.SignalCount)
	}
}
