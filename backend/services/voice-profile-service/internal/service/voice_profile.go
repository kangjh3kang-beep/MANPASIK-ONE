package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ============================================================================
// 에러 및 상수
// ============================================================================

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrProfileLimit    = errors.New("profile limit exceeded")
	ErrInvalidLanguage = errors.New("unsupported language")
	ErrInvalidStatus   = errors.New("invalid status transition")
)

// 지원 언어 목록
var supportedLanguages = map[string]bool{
	"ko": true, "en": true, "ja": true, "zh": true,
	"es": true, "fr": true, "de": true,
}

// 프로필 상태 상수
const (
	ProfileStatusProcessing = "processing"
	ProfileStatusReady      = "ready"
	ProfileStatusError      = "error"
)

// 사용자 당 최대 프로필 수
const MaxProfilesPerUser = 10

// 프로필 이름 길이 제한
const (
	MinProfileNameLen = 1
	MaxProfileNameLen = 50
)

// ============================================================================
// 도메인 모델
// ============================================================================

type VoiceProfile struct {
	ProfileID, UserID, ProfileName, Language, Status, ModelURL string
	CreatedAt                                                  time.Time
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

type VoiceProfileRepository interface {
	Create(ctx context.Context, p *VoiceProfile) error
	GetByID(ctx context.Context, profileID string) (*VoiceProfile, error)
	ListByUser(ctx context.Context, userID string) ([]*VoiceProfile, error)
	Update(ctx context.Context, p *VoiceProfile) error
	Delete(ctx context.Context, profileID string) error
}

// ============================================================================
// VoiceProfileService
// ============================================================================

type VoiceProfileService struct {
	log  *zap.Logger
	repo VoiceProfileRepository
}

func NewVoiceProfileService(log *zap.Logger, repo VoiceProfileRepository) *VoiceProfileService {
	return &VoiceProfileService{log: log, repo: repo}
}

func (s *VoiceProfileService) Create(ctx context.Context, userID, profileName, language, voiceSampleURL string) (*VoiceProfile, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	if profileName == "" {
		return nil, fmt.Errorf("%w: 프로필 이름은 필수입니다", ErrInvalidInput)
	}
	if len(profileName) < MinProfileNameLen || len(profileName) > MaxProfileNameLen {
		return nil, fmt.Errorf("%w: 프로필 이름은 %d~%d자여야 합니다", ErrInvalidInput, MinProfileNameLen, MaxProfileNameLen)
	}
	if language == "" {
		language = "ko"
	}
	if !supportedLanguages[language] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidLanguage, language)
	}

	// 사용자 당 프로필 수 제한 확인
	existing, _ := s.repo.ListByUser(ctx, userID)
	if len(existing) >= MaxProfilesPerUser {
		return nil, fmt.Errorf("%w: 최대 %d개 프로필", ErrProfileLimit, MaxProfilesPerUser)
	}

	p := &VoiceProfile{
		ProfileID: uuid.New().String(), UserID: userID, ProfileName: profileName,
		Language: language, Status: ProfileStatusProcessing, CreatedAt: time.Now().UTC(),
	}
	return p, s.repo.Create(ctx, p)
}

func (s *VoiceProfileService) Get(ctx context.Context, profileID string) (*VoiceProfile, error) {
	if profileID == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByID(ctx, profileID)
}

func (s *VoiceProfileService) List(ctx context.Context, userID string) ([]*VoiceProfile, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.ListByUser(ctx, userID)
}

// UpdateStatus는 프로필 상태를 변경합니다.
func (s *VoiceProfileService) UpdateStatus(ctx context.Context, profileID, newStatus string) (*VoiceProfile, error) {
	if profileID == "" {
		return nil, ErrInvalidInput
	}

	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	// 상태 전이 검증: processing → ready|error, error → processing
	valid := false
	switch p.Status {
	case ProfileStatusProcessing:
		valid = newStatus == ProfileStatusReady || newStatus == ProfileStatusError
	case ProfileStatusError:
		valid = newStatus == ProfileStatusProcessing
	}
	if !valid {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidStatus, p.Status, newStatus)
	}

	p.Status = newStatus
	if newStatus == ProfileStatusReady {
		p.ModelURL = "https://models.manpasik.com/voice/" + profileID + ".bin"
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdateName은 프로필 이름을 변경합니다.
func (s *VoiceProfileService) UpdateName(ctx context.Context, profileID, newName string) (*VoiceProfile, error) {
	if profileID == "" {
		return nil, ErrInvalidInput
	}
	if newName == "" || len(newName) < MinProfileNameLen || len(newName) > MaxProfileNameLen {
		return nil, fmt.Errorf("%w: 프로필 이름은 %d~%d자여야 합니다", ErrInvalidInput, MinProfileNameLen, MaxProfileNameLen)
	}

	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	p.ProfileName = newName
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *VoiceProfileService) Delete(ctx context.Context, profileID string) error {
	if profileID == "" {
		return ErrInvalidInput
	}
	return s.repo.Delete(ctx, profileID)
}

func (s *VoiceProfileService) SynthesizeTranslation(ctx context.Context, profileID, text, targetLang string) (string, string, error) {
	if profileID == "" || text == "" {
		return "", "", ErrInvalidInput
	}
	if targetLang != "" && !supportedLanguages[targetLang] {
		return "", "", fmt.Errorf("%w: %s", ErrInvalidLanguage, targetLang)
	}

	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return "", "", err
	}

	// ready 상태가 아니면 합성 불가
	if p.Status != ProfileStatusReady {
		return "", "", fmt.Errorf("프로필이 준비 상태가 아닙니다: %s", p.Status)
	}

	audioURL := "https://tts.manpasik.com/audio/" + uuid.New().String() + ".wav"
	translatedText := "[" + targetLang + "] " + text
	return audioURL, translatedText, nil
}
