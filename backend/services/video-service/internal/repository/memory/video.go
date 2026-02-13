// Package memory는 video-service의 인메모리 저장소를 구현합니다.
package memory

import (
	"context"
	"sync"

	apperrors "github.com/manpasik/backend/shared/errors"
	"github.com/manpasik/backend/services/video-service/internal/service"
)

// RoomRepository는 인메모리 회의실 저장소입니다.
type RoomRepository struct {
	mu    sync.RWMutex
	rooms map[string]*service.Room
}

// NewRoomRepository는 RoomRepository를 생성합니다.
func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*service.Room),
	}
}

func (r *RoomRepository) Save(_ context.Context, room *service.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[room.ID] = room
	return nil
}

func (r *RoomRepository) FindByID(_ context.Context, id string) (*service.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	room, ok := r.rooms[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "회의실을 찾을 수 없습니다")
	}
	return room, nil
}

func (r *RoomRepository) Update(_ context.Context, room *service.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rooms[room.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "회의실을 찾을 수 없습니다")
	}
	r.rooms[room.ID] = room
	return nil
}

// SignalRepository는 인메모리 시그널 저장소입니다.
type SignalRepository struct {
	mu      sync.RWMutex
	signals []*service.Signal
}

// NewSignalRepository는 SignalRepository를 생성합니다.
func NewSignalRepository() *SignalRepository {
	return &SignalRepository{}
}

func (r *SignalRepository) Save(_ context.Context, signal *service.Signal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.signals = append(r.signals, signal)
	return nil
}

func (r *SignalRepository) CountByRoomID(_ context.Context, roomID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, s := range r.signals {
		if s.RoomID == roomID {
			count++
		}
	}
	return count, nil
}
