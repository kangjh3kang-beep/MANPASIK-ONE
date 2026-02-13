// Package cache는 device-service의 Redis 캐시 데코레이터입니다.
//
// Cache-Aside 패턴 적용:
//   - 읽기: 캐시 hit → 반환, miss → DB 조회 후 캐시 저장
//   - 쓰기: DB 처리 후 관련 캐시 키 삭제 (무효화)
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/manpasik/backend/services/device-service/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	deviceByIDPrefix    = "device:id:"
	deviceListByUser    = "device:user:list:"
	deviceCountByUser   = "device:user:count:"
	deviceTTL           = 5 * time.Minute
	deviceListTTL       = 1 * time.Minute
	deviceCountTTL      = 2 * time.Minute
)

// DeviceRepository는 Redis 캐시를 적용한 DeviceRepository 래퍼입니다.
type DeviceRepository struct {
	inner service.DeviceRepository
	rdb   *redis.Client
}

// NewDeviceRepository는 캐시 DeviceRepository를 생성합니다.
func NewDeviceRepository(inner service.DeviceRepository, rdb *redis.Client) *DeviceRepository {
	return &DeviceRepository{inner: inner, rdb: rdb}
}

// Create는 디바이스를 생성하고 관련 캐시를 무효화합니다.
func (r *DeviceRepository) Create(ctx context.Context, device *service.Device) error {
	if err := r.inner.Create(ctx, device); err != nil {
		return err
	}
	// 사용자별 목록/카운트 캐시 무효화
	r.rdb.Del(ctx,
		deviceListByUser+device.UserID,
		deviceCountByUser+device.UserID,
	)
	return nil
}

// GetByID는 캐시를 확인한 후 DB에서 조회합니다.
func (r *DeviceRepository) GetByID(ctx context.Context, deviceID string) (*service.Device, error) {
	key := deviceByIDPrefix + deviceID

	// 캐시 히트 시도
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var d service.Device
		if json.Unmarshal(data, &d) == nil {
			return &d, nil
		}
	}

	// DB 조회
	device, err := r.inner.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, nil
	}

	// 캐시 저장
	if b, e := json.Marshal(device); e == nil {
		r.rdb.Set(ctx, key, b, deviceTTL)
	}
	return device, nil
}

// ListByUser는 캐시를 확인한 후 DB에서 조회합니다.
func (r *DeviceRepository) ListByUser(ctx context.Context, userID string) ([]*service.Device, error) {
	key := deviceListByUser + userID

	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var devices []*service.Device
		if json.Unmarshal(data, &devices) == nil {
			return devices, nil
		}
	}

	devices, err := r.inner.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if b, e := json.Marshal(devices); e == nil {
		r.rdb.Set(ctx, key, b, deviceListTTL)
	}
	return devices, nil
}

// UpdateStatus는 상태를 업데이트하고 캐시를 무효화합니다.
func (r *DeviceRepository) UpdateStatus(ctx context.Context, deviceID string, status service.DeviceStatus, battery int, lastSeen time.Time) error {
	if err := r.inner.UpdateStatus(ctx, deviceID, status, battery, lastSeen); err != nil {
		return err
	}
	// 해당 디바이스 캐시 무효화
	r.rdb.Del(ctx, deviceByIDPrefix+deviceID)

	// 사용자 ID를 모르므로 목록/카운트는 TTL에 의존
	// (GetByID 후 무효화하면 추가 쿼리 발생 → TTL 방식 선택)
	return nil
}

// CountByUser는 캐시를 확인한 후 DB에서 조회합니다.
func (r *DeviceRepository) CountByUser(ctx context.Context, userID string) (int, error) {
	key := deviceCountByUser + userID

	data, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var count int
		if _, e := fmt.Sscanf(data, "%d", &count); e == nil {
			return count, nil
		}
	}

	count, err := r.inner.CountByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	r.rdb.Set(ctx, key, fmt.Sprintf("%d", count), deviceCountTTL)
	return count, nil
}

// Delete는 디바이스를 삭제하고 캐시를 무효화합니다.
func (r *DeviceRepository) Delete(ctx context.Context, deviceID string) error {
	// 삭제 전에 디바이스 조회 (userID를 위해)
	device, _ := r.inner.GetByID(ctx, deviceID)

	if err := r.inner.Delete(ctx, deviceID); err != nil {
		return err
	}

	// 캐시 무효화
	keys := []string{deviceByIDPrefix + deviceID}
	if device != nil {
		keys = append(keys,
			deviceListByUser+device.UserID,
			deviceCountByUser+device.UserID,
		)
	}
	r.rdb.Del(ctx, keys...)
	return nil
}
