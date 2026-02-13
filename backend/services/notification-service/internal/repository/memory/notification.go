// Package memory는 notification-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/notification-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// NotificationRepository는 알림 인메모리 저장소입니다.
type NotificationRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Notification
}

// NewNotificationRepository는 새 인메모리 알림 저장소를 생성합니다.
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		store: make(map[string]*service.Notification),
	}
}

// Save는 알림을 저장합니다.
func (r *NotificationRepository) Save(_ context.Context, n *service.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[n.ID] = n
	return nil
}

// FindByID는 ID로 알림을 조회합니다.
func (r *NotificationRepository) FindByID(_ context.Context, id string) (*service.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return n, nil
}

// FindByUserID는 사용자 알림 목록을 조회합니다.
func (r *NotificationRepository) FindByUserID(_ context.Context, userID string, typeFilter service.NotificationType, unreadOnly bool, limit, offset int) ([]*service.Notification, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Notification
	for _, n := range r.store {
		if n.UserID != userID {
			continue
		}
		if typeFilter != service.TypeUnknown && n.Type != typeFilter {
			continue
		}
		if unreadOnly && n.IsRead {
			continue
		}
		filtered = append(filtered, n)
	}

	// 시간 역순 정렬 (최신 먼저)
	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].CreatedAt.After(filtered[i].CreatedAt) {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

// MarkAsRead는 알림을 읽음 처리합니다.
func (r *NotificationRepository) MarkAsRead(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.store[id]
	if !ok {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
	return nil
}

// MarkAllAsRead는 사용자의 모든 알림을 읽음 처리합니다.
func (r *NotificationRepository) MarkAllAsRead(_ context.Context, userID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	now := time.Now()
	for _, n := range r.store {
		if n.UserID == userID && !n.IsRead {
			n.IsRead = true
			n.ReadAt = &now
			count++
		}
	}
	return count, nil
}

// GetUnreadCount는 읽지 않은 알림 수를 반환합니다.
func (r *NotificationRepository) GetUnreadCount(_ context.Context, userID string) (int, map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := 0
	byType := make(map[string]int)

	for _, n := range r.store {
		if n.UserID == userID && !n.IsRead {
			total++
			typeName := service.NotificationTypeToString(n.Type)
			byType[typeName]++
		}
	}

	return total, byType, nil
}

// PreferencesRepository는 알림 설정 인메모리 저장소입니다.
type PreferencesRepository struct {
	mu    sync.RWMutex
	store map[string]*service.NotificationPreferences
}

// NewPreferencesRepository는 새 인메모리 알림 설정 저장소를 생성합니다.
func NewPreferencesRepository() *PreferencesRepository {
	return &PreferencesRepository{
		store: make(map[string]*service.NotificationPreferences),
	}
}

// Save는 알림 설정을 저장합니다.
func (r *PreferencesRepository) Save(_ context.Context, pref *service.NotificationPreferences) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[pref.UserID] = pref
	return nil
}

// FindByUserID는 사용자 알림 설정을 조회합니다.
func (r *PreferencesRepository) FindByUserID(_ context.Context, userID string) (*service.NotificationPreferences, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pref, ok := r.store[userID]
	if !ok {
		return nil, nil // 설정 없음 (기본값 사용)
	}
	return pref, nil
}
