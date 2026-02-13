// Package postgres는 notification-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/notification-service/internal/service"
)

// typeToDB maps service NotificationType to DB enum (schema has measurement, health_alert, coaching, subscription, order, family, system, reminder, promotion)
var typeToDB = map[service.NotificationType]string{
	service.TypeMeasurement:  "measurement",
	service.TypeHealthAlert:  "health_alert",
	service.TypeAppointment:  "reminder",
	service.TypePrescription: "order",
	service.TypeCommunity:    "family",
	service.TypeSystem:       "system",
	service.TypePromotion:    "promotion",
	service.TypeUnknown:      "system",
}

var dbToType = map[string]service.NotificationType{
	"measurement":  service.TypeMeasurement,
	"health_alert": service.TypeHealthAlert,
	"reminder":     service.TypeAppointment,
	"order":        service.TypePrescription,
	"family":       service.TypeCommunity,
	"system":       service.TypeSystem,
	"promotion":    service.TypePromotion,
}

var channelToDB = map[service.NotificationChannel]string{
	service.ChannelPush:   "push",
	service.ChannelEmail:  "email",
	service.ChannelSMS:    "sms",
	service.ChannelInApp:  "in_app",
	service.ChannelUnknown: "in_app",
}

var dbToChannel = map[string]service.NotificationChannel{
	"push":   service.ChannelPush,
	"email":  service.ChannelEmail,
	"sms":    service.ChannelSMS,
	"in_app": service.ChannelInApp,
}

var priorityToDB = map[service.NotificationPriority]string{
	service.PriorityLow:    "low",
	service.PriorityNormal: "normal",
	service.PriorityHigh:   "high",
	service.PriorityUrgent: "urgent",
	service.PriorityUnknown: "normal",
}

var dbToPriority = map[string]service.NotificationPriority{
	"low":    service.PriorityLow,
	"normal": service.PriorityNormal,
	"high":   service.PriorityHigh,
	"urgent": service.PriorityUrgent,
}

// ============================================================================
// NotificationRepository
// ============================================================================

// NotificationRepository는 PostgreSQL 기반 NotificationRepository 구현입니다.
type NotificationRepository struct {
	pool *pgxpool.Pool
}

// NewNotificationRepository는 NotificationRepository를 생성합니다.
func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

// Save는 알림을 저장합니다.
func (r *NotificationRepository) Save(ctx context.Context, n *service.Notification) error {
	dataJSON := []byte("{}")
	if n.Data != "" {
		if json.Valid([]byte(n.Data)) {
			dataJSON = []byte(n.Data)
		} else {
			// store as {"raw": data}
			b, _ := json.Marshal(map[string]string{"raw": n.Data})
			dataJSON = b
		}
	}
	typeStr, _ := typeToDB[n.Type]
	if typeStr == "" {
		typeStr = "system"
	}
	channelStr, _ := channelToDB[n.Channel]
	if channelStr == "" {
		channelStr = "in_app"
	}
	priorityStr, _ := priorityToDB[n.Priority]
	if priorityStr == "" {
		priorityStr = "normal"
	}

	const q = `INSERT INTO notifications (id, user_id, type, channel, priority, title, body, data, is_read, created_at, read_at)
		VALUES ($1, $2, $3::notification_type, $4::notification_channel, $5::notification_priority, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, q,
		n.ID, n.UserID, typeStr, channelStr, priorityStr,
		n.Title, nullIfEmpty(n.Body), dataJSON, n.IsRead, n.CreatedAt, nullTime(n.ReadAt),
	)
	return err
}

// FindByID는 ID로 알림을 조회합니다.
func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*service.Notification, error) {
	const q = `SELECT id, user_id, type::text, channel::text, priority::text, title, COALESCE(body,''), COALESCE(data::text,'{}'), is_read, created_at, read_at
		FROM notifications WHERE id = $1`
	var n service.Notification
	var typeStr, channelStr, priorityStr string
	var dataJSON string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&n.ID, &n.UserID, &typeStr, &channelStr, &priorityStr,
		&n.Title, &n.Body, &dataJSON, &n.IsRead, &n.CreatedAt, &n.ReadAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	n.Type = dbToType[typeStr]
	if n.Type == 0 {
		n.Type = service.TypeUnknown
	}
	n.Channel = dbToChannel[channelStr]
	if n.Channel == 0 {
		n.Channel = service.ChannelInApp
	}
	n.Priority = dbToPriority[priorityStr]
	if n.Priority == 0 {
		n.Priority = service.PriorityNormal
	}
	if dataJSON != "" && dataJSON != "{}" {
		var m map[string]interface{}
		if json.Unmarshal([]byte(dataJSON), &m) == nil {
			if raw, ok := m["raw"].(string); ok {
				n.Data = raw
			} else {
				n.Data = dataJSON
			}
		} else {
			n.Data = dataJSON
		}
	}
	return &n, nil
}

// FindByUserID는 사용자 알림 목록을 조회합니다 (페이지네이션).
func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string, typeFilter service.NotificationType, unreadOnly bool, limit, offset int) ([]*service.Notification, int, error) {
	// count total
	var countQuery string
	countArgs := []interface{}{userID}
	if typeFilter != service.TypeUnknown {
		typeStr, _ := typeToDB[typeFilter]
		if typeStr == "" {
			typeStr = "system"
		}
		countQuery = `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND type = $2::notification_type`
		countArgs = append(countArgs, typeStr)
	} else {
		countQuery = `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	}
	if unreadOnly {
		countQuery += ` AND is_read = FALSE`
	}

	var total int
	err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// select with pagination
	selArgs := []interface{}{userID}
	selQuery := `SELECT id, user_id, type::text, channel::text, priority::text, title, COALESCE(body,''), COALESCE(data::text,'{}'), is_read, created_at, read_at
		FROM notifications WHERE user_id = $1`
	if typeFilter != service.TypeUnknown {
		typeStr, _ := typeToDB[typeFilter]
		if typeStr == "" {
			typeStr = "system"
		}
		selQuery += ` AND type = $2::notification_type`
		selArgs = append(selArgs, typeStr)
	}
	if unreadOnly {
		selQuery += ` AND is_read = FALSE`
	}
	selQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(selArgs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(selArgs)+2)
	selArgs = append(selArgs, limit, offset)

	rows, err := r.pool.Query(ctx, selQuery, selArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.Notification
	for rows.Next() {
		var n service.Notification
		var typeStr, channelStr, priorityStr, dataJSON string
		if err := rows.Scan(
			&n.ID, &n.UserID, &typeStr, &channelStr, &priorityStr,
			&n.Title, &n.Body, &dataJSON, &n.IsRead, &n.CreatedAt, &n.ReadAt,
		); err != nil {
			return nil, 0, err
		}
		n.Type = dbToType[typeStr]
		if n.Type == 0 {
			n.Type = service.TypeUnknown
		}
		n.Channel = dbToChannel[channelStr]
		if n.Channel == 0 {
			n.Channel = service.ChannelInApp
		}
		n.Priority = dbToPriority[priorityStr]
		if n.Priority == 0 {
			n.Priority = service.PriorityNormal
		}
		if dataJSON != "" && dataJSON != "{}" {
			var m map[string]interface{}
			if json.Unmarshal([]byte(dataJSON), &m) == nil {
				if raw, ok := m["raw"].(string); ok {
					n.Data = raw
				} else {
					n.Data = dataJSON
				}
			} else {
				n.Data = dataJSON
			}
		}
		list = append(list, &n)
	}
	return list, total, rows.Err()
}

// MarkAsRead는 알림을 읽음 처리합니다.
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	const q = `UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

// MarkAllAsRead는 사용자의 모든 알림을 읽음 처리합니다.
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) (int, error) {
	const q = `UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE user_id = $1 AND is_read = FALSE`
	res, err := r.pool.Exec(ctx, q, userID)
	if err != nil {
		return 0, err
	}
	return int(res.RowsAffected()), nil
}

// GetUnreadCount는 읽지 않은 알림 수를 반환합니다 (total 및 타입별).
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, map[string]int, error) {
	const q = `SELECT type::text, COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE GROUP BY type`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	byType := make(map[string]int)
	total := 0
	for rows.Next() {
		var t string
		var c int
		if err := rows.Scan(&t, &c); err != nil {
			return 0, nil, err
		}
		typeStr := dbToTypeStr(t)
		byType[typeStr] = c
		total += c
	}
	return total, byType, rows.Err()
}

func dbToTypeStr(db string) string {
	switch db {
	case "measurement":
		return "measurement"
	case "health_alert":
		return "health_alert"
	case "reminder":
		return "appointment"
	case "order":
		return "prescription"
	case "family":
		return "community"
	case "system":
		return "system"
	case "promotion":
		return "promotion"
	default:
		return "unknown"
	}
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}

// ============================================================================
// PreferencesRepository
// ============================================================================

// PreferencesRepository는 PostgreSQL 기반 PreferencesRepository 구현입니다.
type PreferencesRepository struct {
	pool *pgxpool.Pool
}

// NewPreferencesRepository는 PreferencesRepository를 생성합니다.
func NewPreferencesRepository(pool *pgxpool.Pool) *PreferencesRepository {
	return &PreferencesRepository{pool: pool}
}

// Save는 알림 설정을 저장합니다 (UPSERT).
func (r *PreferencesRepository) Save(ctx context.Context, pref *service.NotificationPreferences) error {
	const q = `INSERT INTO notification_preferences (user_id, push_enabled, email_enabled, sms_enabled, in_app_enabled, health_alert_enabled, coaching_enabled, promotion_enabled, quiet_hours_start, quiet_hours_end, language, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			push_enabled = EXCLUDED.push_enabled,
			email_enabled = EXCLUDED.email_enabled,
			sms_enabled = EXCLUDED.sms_enabled,
			in_app_enabled = EXCLUDED.in_app_enabled,
			health_alert_enabled = EXCLUDED.health_alert_enabled,
			coaching_enabled = EXCLUDED.coaching_enabled,
			promotion_enabled = EXCLUDED.promotion_enabled,
			quiet_hours_start = EXCLUDED.quiet_hours_start,
			quiet_hours_end = EXCLUDED.quiet_hours_end,
			language = EXCLUDED.language,
			updated_at = NOW()`
	_, err := r.pool.Exec(ctx, q,
		pref.UserID, pref.PushEnabled, pref.EmailEnabled, pref.SMSEnabled, pref.InAppEnabled,
		pref.HealthAlertEnabled, pref.CoachingEnabled, pref.PromotionEnabled,
		nullIfEmpty(pref.QuietHoursStart), nullIfEmpty(pref.QuietHoursEnd),
		nullIfEmptyOrDefault(pref.Language, "ko"),
	)
	return err
}

// FindByUserID는 사용자 알림 설정을 조회합니다.
func (r *PreferencesRepository) FindByUserID(ctx context.Context, userID string) (*service.NotificationPreferences, error) {
	const q = `SELECT user_id, push_enabled, email_enabled, sms_enabled, in_app_enabled, health_alert_enabled, coaching_enabled, promotion_enabled, COALESCE(quiet_hours_start,''), COALESCE(quiet_hours_end,''), COALESCE(language,'ko')
		FROM notification_preferences WHERE user_id = $1`
	var pref service.NotificationPreferences
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&pref.UserID, &pref.PushEnabled, &pref.EmailEnabled, &pref.SMSEnabled, &pref.InAppEnabled,
		&pref.HealthAlertEnabled, &pref.CoachingEnabled, &pref.PromotionEnabled,
		&pref.QuietHoursStart, &pref.QuietHoursEnd, &pref.Language,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pref, nil
}

func nullIfEmptyOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
