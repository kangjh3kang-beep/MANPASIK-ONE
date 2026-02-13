// Package postgres는 video-service의 PostgreSQL 저장소 구현입니다.
//
// DB 스키마: infrastructure/database/init/21-video.sql
// 테이블: video_rooms, video_participants, video_signals, video_room_stats
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/video-service/internal/service"
)

// ============================================================================
// RoomRepository
// ============================================================================

// RoomRepository는 PostgreSQL 기반 회의실 저장소입니다.
type RoomRepository struct {
	pool *pgxpool.Pool
}

// NewRoomRepository는 RoomRepository를 생성합니다.
func NewRoomRepository(pool *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{pool: pool}
}

func roomTypeToString(t service.RoomType) string {
	switch t {
	case service.RoomTypeOneToOne:
		return "one_to_one"
	case service.RoomTypeGroup:
		return "group"
	case service.RoomTypeWebinar:
		return "webinar"
	case service.RoomTypeConsultation:
		return "consultation"
	default:
		return "one_to_one"
	}
}

func stringToRoomType(s string) service.RoomType {
	switch s {
	case "one_to_one":
		return service.RoomTypeOneToOne
	case "group":
		return service.RoomTypeGroup
	case "webinar":
		return service.RoomTypeWebinar
	case "consultation":
		return service.RoomTypeConsultation
	default:
		return service.RoomTypeUnknown
	}
}

func roomStatusToString(s service.RoomStatus) string {
	switch s {
	case service.RoomStatusWaiting:
		return "waiting"
	case service.RoomStatusActive:
		return "active"
	case service.RoomStatusEnded:
		return "ended"
	case service.RoomStatusFailed:
		return "failed"
	default:
		return "waiting"
	}
}

func stringToRoomStatus(s string) service.RoomStatus {
	switch s {
	case "waiting":
		return service.RoomStatusWaiting
	case "active":
		return service.RoomStatusActive
	case "ended":
		return service.RoomStatusEnded
	case "failed":
		return service.RoomStatusFailed
	default:
		return service.RoomStatusUnknown
	}
}

// Save는 회의실을 저장합니다.
func (r *RoomRepository) Save(ctx context.Context, room *service.Room) error {
	const q = `INSERT INTO video_rooms (id, name, room_type, status, created_by, max_participants,
		recording_url, duration_seconds, total_bytes_transferred, created_at, started_at, ended_at)
		VALUES ($1, $2, $3::room_type, $4::room_status, $5, $6, $7, $8, $9, $10, $11, $12)`

	var startedAt, endedAt *time.Time
	if !room.StartedAt.IsZero() {
		startedAt = &room.StartedAt
	}
	if !room.EndedAt.IsZero() {
		endedAt = &room.EndedAt
	}

	_, err := r.pool.Exec(ctx, q,
		room.ID, room.Name, roomTypeToString(room.RoomType), roomStatusToString(room.Status),
		room.CreatedBy, room.MaxParticipants, room.RecordingURL,
		room.DurationSeconds, room.TotalBytesTransferred, room.CreatedAt, startedAt, endedAt,
	)
	return err
}

// FindByID는 회의실을 ID로 조회합니다.
func (r *RoomRepository) FindByID(ctx context.Context, id string) (*service.Room, error) {
	const q = `SELECT id, name, room_type, status, created_by, max_participants,
		COALESCE(recording_url,''), duration_seconds, total_bytes_transferred,
		created_at, COALESCE(started_at, '0001-01-01'), COALESCE(ended_at, '0001-01-01')
		FROM video_rooms WHERE id = $1`

	var room service.Room
	var typeStr, statusStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&room.ID, &room.Name, &typeStr, &statusStr, &room.CreatedBy,
		&room.MaxParticipants, &room.RecordingURL,
		&room.DurationSeconds, &room.TotalBytesTransferred,
		&room.CreatedAt, &room.StartedAt, &room.EndedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	room.RoomType = stringToRoomType(typeStr)
	room.Status = stringToRoomStatus(statusStr)

	// 참가자 조회
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, COALESCE(display_name,''), COALESCE(role,'participant'),
			is_audio_enabled, is_video_enabled, is_screen_sharing,
			joined_at, COALESCE(left_at, '0001-01-01')
		FROM video_participants WHERE room_id = $1 ORDER BY joined_at`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p service.Participant
			if rows.Scan(&p.UserID, &p.DisplayName, &p.Role,
				&p.IsAudioEnabled, &p.IsVideoEnabled, &p.IsScreenSharing,
				&p.JoinedAt, &p.LeftAt) == nil {
				room.Participants = append(room.Participants, &p)
			}
		}
	}
	if room.Participants == nil {
		room.Participants = make([]*service.Participant, 0)
	}

	return &room, nil
}

// Update는 회의실을 업데이트합니다.
func (r *RoomRepository) Update(ctx context.Context, room *service.Room) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var startedAt, endedAt *time.Time
	if !room.StartedAt.IsZero() {
		startedAt = &room.StartedAt
	}
	if !room.EndedAt.IsZero() {
		endedAt = &room.EndedAt
	}

	const q = `UPDATE video_rooms SET status = $1::room_status, recording_url = $2,
		duration_seconds = $3, total_bytes_transferred = $4, started_at = $5, ended_at = $6
		WHERE id = $7`
	_, err = tx.Exec(ctx, q,
		roomStatusToString(room.Status), room.RecordingURL,
		room.DurationSeconds, room.TotalBytesTransferred, startedAt, endedAt, room.ID,
	)
	if err != nil {
		return err
	}

	// 참가자 업서트
	for _, p := range room.Participants {
		var leftAt *time.Time
		if !p.LeftAt.IsZero() {
			leftAt = &p.LeftAt
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO video_participants (room_id, user_id, display_name, role,
				is_audio_enabled, is_video_enabled, is_screen_sharing, joined_at, left_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (room_id, user_id, joined_at) DO UPDATE
				SET left_at = EXCLUDED.left_at, is_audio_enabled = EXCLUDED.is_audio_enabled,
				is_video_enabled = EXCLUDED.is_video_enabled, is_screen_sharing = EXCLUDED.is_screen_sharing`,
			room.ID, p.UserID, p.DisplayName, p.Role,
			p.IsAudioEnabled, p.IsVideoEnabled, p.IsScreenSharing,
			p.JoinedAt, leftAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// ============================================================================
// SignalRepository
// ============================================================================

// SignalRepository는 PostgreSQL 기반 시그널 저장소입니다.
type SignalRepository struct {
	pool *pgxpool.Pool
}

// NewSignalRepository는 SignalRepository를 생성합니다.
func NewSignalRepository(pool *pgxpool.Pool) *SignalRepository {
	return &SignalRepository{pool: pool}
}

func signalTypeToString(t service.SignalType) string {
	switch t {
	case service.SignalTypeOffer:
		return "offer"
	case service.SignalTypeAnswer:
		return "answer"
	case service.SignalTypeICECandidate:
		return "ice_candidate"
	case service.SignalTypeRenegotiate:
		return "renegotiate"
	case service.SignalTypeMute:
		return "mute"
	case service.SignalTypeUnmute:
		return "unmute"
	default:
		return "offer"
	}
}

// Save는 시그널을 저장합니다.
func (r *SignalRepository) Save(ctx context.Context, s *service.Signal) error {
	var toUserID *string
	if s.ToUserID != "" {
		toUserID = &s.ToUserID
	}
	const q = `INSERT INTO video_signals (id, room_id, from_user_id, to_user_id, signal_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5::signal_type, $6::jsonb, $7)`
	_, err := r.pool.Exec(ctx, q,
		s.ID, s.RoomID, s.FromUserID, toUserID,
		signalTypeToString(s.Type), s.Payload, s.CreatedAt,
	)
	return err
}

// CountByRoomID는 방의 시그널 수를 반환합니다.
func (r *SignalRepository) CountByRoomID(ctx context.Context, roomID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM video_signals WHERE room_id = $1", roomID).Scan(&count)
	return count, err
}
