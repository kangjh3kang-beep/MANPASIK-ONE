// Package service는 video-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// RoomStatus는 방 상태입니다.
type RoomStatus int

const (
	RoomStatusUnknown RoomStatus = iota
	RoomStatusWaiting
	RoomStatusActive
	RoomStatusEnded
	RoomStatusFailed
)

// RoomType은 방 유형입니다.
type RoomType int

const (
	RoomTypeUnknown      RoomType = iota
	RoomTypeOneToOne
	RoomTypeGroup
	RoomTypeWebinar
	RoomTypeConsultation
)

// SignalType은 시그널 유형입니다.
type SignalType int

const (
	SignalTypeUnknown      SignalType = iota
	SignalTypeOffer
	SignalTypeAnswer
	SignalTypeICECandidate
	SignalTypeRenegotiate
	SignalTypeMute
	SignalTypeUnmute
)

// Room은 화상 회의실 도메인 객체입니다.
type Room struct {
	ID                  string
	Name                string
	RoomType            RoomType
	Status              RoomStatus
	CreatedBy           string
	MaxParticipants     int
	Participants        []*Participant
	RecordingURL        string
	CreatedAt           time.Time
	StartedAt           time.Time
	EndedAt             time.Time
	DurationSeconds     int
	TotalBytesTransferred int64
	SignalCount         int
}

// Participant는 참가자 도메인 객체입니다.
type Participant struct {
	UserID         string
	DisplayName    string
	Role           string
	IsAudioEnabled bool
	IsVideoEnabled bool
	IsScreenSharing bool
	JoinedAt       time.Time
	LeftAt         time.Time
}

// Signal은 WebRTC 시그널 도메인 객체입니다.
type Signal struct {
	ID         string
	RoomID     string
	FromUserID string
	ToUserID   string
	Type       SignalType
	Payload    string
	CreatedAt  time.Time
}

// RoomRepository는 회의실 저장소 인터페이스입니다.
type RoomRepository interface {
	Save(ctx context.Context, r *Room) error
	FindByID(ctx context.Context, id string) (*Room, error)
	Update(ctx context.Context, r *Room) error
}

// SignalRepository는 시그널 저장소 인터페이스입니다.
type SignalRepository interface {
	Save(ctx context.Context, s *Signal) error
	CountByRoomID(ctx context.Context, roomID string) (int, error)
}

// VideoService는 비디오 서비스 핵심 로직입니다.
type VideoService struct {
	log        *zap.Logger
	roomRepo   RoomRepository
	signalRepo SignalRepository
}

// NewVideoService는 VideoService를 생성합니다.
func NewVideoService(log *zap.Logger, roomRepo RoomRepository, signalRepo SignalRepository) *VideoService {
	return &VideoService{
		log:        log,
		roomRepo:   roomRepo,
		signalRepo: signalRepo,
	}
}

// CreateRoom은 새 회의실을 생성합니다.
func (s *VideoService) CreateRoom(ctx context.Context, name string, roomType RoomType, createdBy string, maxParticipants int) (*Room, error) {
	if createdBy == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "created_by는 필수입니다")
	}
	if name == "" {
		name = fmt.Sprintf("Room-%s", uuid.New().String()[:8])
	}
	if maxParticipants <= 0 {
		maxParticipants = 2
	}
	if maxParticipants > 100 {
		maxParticipants = 100
	}

	room := &Room{
		ID:              uuid.New().String(),
		Name:            name,
		RoomType:        roomType,
		Status:          RoomStatusWaiting,
		CreatedBy:       createdBy,
		MaxParticipants: maxParticipants,
		Participants:    make([]*Participant, 0),
		CreatedAt:       time.Now(),
	}

	if err := s.roomRepo.Save(ctx, room); err != nil {
		return nil, fmt.Errorf("회의실 저장 실패: %w", err)
	}

	s.log.Info("회의실 생성 완료",
		zap.String("room_id", room.ID),
		zap.String("created_by", createdBy),
	)

	return room, nil
}

// GetRoom은 회의실을 조회합니다.
func (s *VideoService) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	if roomID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}
	return s.roomRepo.FindByID(ctx, roomID)
}

// JoinRoom은 참가자를 회의실에 추가합니다.
func (s *VideoService) JoinRoom(ctx context.Context, roomID, userID, displayName, role string) (*Room, string, string, error) {
	if roomID == "" {
		return nil, "", "", apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}
	if userID == "" {
		return nil, "", "", apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, "", "", err
	}

	if room.Status == RoomStatusEnded {
		return nil, "", "", apperrors.New(apperrors.ErrConflict, "종료된 회의실에 참가할 수 없습니다")
	}

	currentCount := 0
	for _, p := range room.Participants {
		if p.LeftAt.IsZero() {
			currentCount++
		}
	}
	if currentCount >= room.MaxParticipants {
		return nil, "", "", apperrors.New(apperrors.ErrConflict, "회의실 최대 인원을 초과했습니다")
	}

	if displayName == "" {
		displayName = "User-" + userID[:8]
	}
	if role == "" {
		role = "participant"
	}

	participant := &Participant{
		UserID:         userID,
		DisplayName:    displayName,
		Role:           role,
		IsAudioEnabled: true,
		IsVideoEnabled: true,
		JoinedAt:       time.Now(),
	}
	room.Participants = append(room.Participants, participant)

	if room.Status == RoomStatusWaiting {
		room.Status = RoomStatusActive
		room.StartedAt = time.Now()
	}

	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, "", "", fmt.Errorf("참가자 추가 실패: %w", err)
	}

	token := uuid.New().String()
	iceServers := `[{"urls":["stun:stun.l.google.com:19302"]},{"urls":["turn:turn.manpasik.com:3478"],"username":"manpasik","credential":"secret"}]`

	s.log.Info("참가자 입장",
		zap.String("room_id", roomID),
		zap.String("user_id", userID),
	)

	return room, token, iceServers, nil
}

// LeaveRoom은 참가자를 회의실에서 제거합니다.
func (s *VideoService) LeaveRoom(ctx context.Context, roomID, userID string) (int, error) {
	if roomID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}
	if userID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return 0, err
	}

	found := false
	remaining := 0
	for _, p := range room.Participants {
		if p.UserID == userID && p.LeftAt.IsZero() {
			p.LeftAt = time.Now()
			found = true
		}
		if p.LeftAt.IsZero() {
			remaining++
		}
	}

	if !found {
		return 0, apperrors.New(apperrors.ErrNotFound, "참가자를 찾을 수 없습니다")
	}

	if err := s.roomRepo.Update(ctx, room); err != nil {
		return 0, fmt.Errorf("참가자 제거 실패: %w", err)
	}

	s.log.Info("참가자 퇴장",
		zap.String("room_id", roomID),
		zap.String("user_id", userID),
		zap.Int("remaining", remaining),
	)

	return remaining, nil
}

// EndRoom은 회의실을 종료합니다.
func (s *VideoService) EndRoom(ctx context.Context, roomID, userID string) (*Room, error) {
	if roomID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	room.Status = RoomStatusEnded
	room.EndedAt = now
	if !room.StartedAt.IsZero() {
		room.DurationSeconds = int(now.Sub(room.StartedAt).Seconds())
	}

	// 남은 참가자 모두 퇴장 처리
	for _, p := range room.Participants {
		if p.LeftAt.IsZero() {
			p.LeftAt = now
		}
	}

	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("회의실 종료 실패: %w", err)
	}

	s.log.Info("회의실 종료",
		zap.String("room_id", roomID),
		zap.Int("duration_seconds", room.DurationSeconds),
	)

	return room, nil
}

// SendSignal은 WebRTC 시그널을 전송합니다.
func (s *VideoService) SendSignal(ctx context.Context, roomID, fromUserID, toUserID string, signalType SignalType, payload string) (string, error) {
	if roomID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}
	if fromUserID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "from_user_id는 필수입니다")
	}

	// 방 존재 확인
	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return "", err
	}
	if room.Status == RoomStatusEnded {
		return "", apperrors.New(apperrors.ErrConflict, "종료된 회의실에서 시그널을 보낼 수 없습니다")
	}

	signal := &Signal{
		ID:         uuid.New().String(),
		RoomID:     roomID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Type:       signalType,
		Payload:    payload,
		CreatedAt:  time.Now(),
	}

	if err := s.signalRepo.Save(ctx, signal); err != nil {
		return "", fmt.Errorf("시그널 저장 실패: %w", err)
	}

	room.SignalCount++
	s.roomRepo.Update(ctx, room)

	s.log.Info("시그널 전송",
		zap.String("signal_id", signal.ID),
		zap.String("room_id", roomID),
		zap.Int("signal_type", int(signalType)),
	)

	return signal.ID, nil
}

// ListParticipants는 회의실 참가자 목록을 반환합니다.
func (s *VideoService) ListParticipants(ctx context.Context, roomID string) ([]*Participant, error) {
	if roomID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	// 현재 활성 참가자만 반환
	var active []*Participant
	for _, p := range room.Participants {
		if p.LeftAt.IsZero() {
			active = append(active, p)
		}
	}

	return active, nil
}

// GetRoomStats는 회의실 통계를 반환합니다.
func (s *VideoService) GetRoomStats(ctx context.Context, roomID string) (*Room, error) {
	if roomID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "room_id는 필수입니다")
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	signalCount, _ := s.signalRepo.CountByRoomID(ctx, roomID)
	room.SignalCount = signalCount

	return room, nil
}
