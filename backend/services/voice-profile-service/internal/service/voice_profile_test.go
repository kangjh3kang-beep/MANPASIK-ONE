package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/voice-profile-service/internal/repository/memory"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.VoiceProfileService {
	return service.NewVoiceProfileService(zap.NewNop(), memory.NewVoiceProfileRepository())
}
func TestCreateAndGet(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, err := s.Create(ctx, "u1", "내 목소리", "ko", "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Get(ctx, p.ProfileID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ProfileName != "내 목소리" {
		t.Error("name mismatch")
	}
}
func TestListProfiles(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	s.Create(ctx, "u1", "A프로필", "ko", "")
	s.Create(ctx, "u1", "B프로필", "en", "")
	s.Create(ctx, "u2", "C프로필", "ko", "")
	list, _ := s.List(ctx, "u1")
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}
func TestDeleteProfile(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "삭제대상", "ko", "")
	if err := s.Delete(ctx, p.ProfileID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, p.ProfileID); err == nil {
		t.Error("expected error")
	}
}
func TestSynthesize(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "합성용", "ko", "")

	// processing 상태에서는 합성 불가
	_, _, err := s.SynthesizeTranslation(ctx, p.ProfileID, "안녕하세요", "en")
	if err == nil {
		t.Error("processing 상태에서 합성이 성공해서는 안 됩니다")
	}

	// ready로 변경 후 합성
	s.UpdateStatus(ctx, p.ProfileID, "ready")
	url, text, err := s.SynthesizeTranslation(ctx, p.ProfileID, "안녕하세요", "en")
	if err != nil {
		t.Fatal(err)
	}
	if url == "" || text == "" {
		t.Error("empty result")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestCreate_InvalidLanguage(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	_, err := s.Create(ctx, "u1", "프로필", "xx", "")
	if err == nil {
		t.Error("지원하지 않는 언어에 에러가 반환되어야 합니다")
	}
}

func TestCreate_ProfileLimit(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	for i := 0; i < service.MaxProfilesPerUser; i++ {
		_, err := s.Create(ctx, "u-limit", "프로필"+string(rune('A'+i)), "ko", "")
		if err != nil {
			t.Fatalf("프로필 %d 생성 실패: %v", i, err)
		}
	}

	// 11번째 → 에러
	_, err := s.Create(ctx, "u-limit", "초과프로필", "ko", "")
	if err == nil {
		t.Error("프로필 제한 초과 시 에러가 반환되어야 합니다")
	}
}

func TestUpdateStatus_Lifecycle(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "상태테스트", "ko", "")

	// processing → ready
	updated, err := s.UpdateStatus(ctx, p.ProfileID, "ready")
	if err != nil {
		t.Fatalf("ready 전환 실패: %v", err)
	}
	if updated.Status != "ready" {
		t.Errorf("Status = %s, want ready", updated.Status)
	}
	if updated.ModelURL == "" {
		t.Error("ready 상태에서 ModelURL이 설정되어야 합니다")
	}

	// ready → processing (불가)
	_, err = s.UpdateStatus(ctx, p.ProfileID, "processing")
	if err == nil {
		t.Error("ready → processing 전환에 에러가 반환되어야 합니다")
	}
}

func TestUpdateStatus_ErrorRecovery(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "에러복구", "ko", "")

	// processing → error
	updated, err := s.UpdateStatus(ctx, p.ProfileID, "error")
	if err != nil {
		t.Fatalf("error 전환 실패: %v", err)
	}
	if updated.Status != "error" {
		t.Errorf("Status = %s, want error", updated.Status)
	}

	// error → processing (재처리)
	recovered, err := s.UpdateStatus(ctx, p.ProfileID, "processing")
	if err != nil {
		t.Fatalf("재처리 전환 실패: %v", err)
	}
	if recovered.Status != "processing" {
		t.Errorf("Status = %s, want processing", recovered.Status)
	}
}

func TestUpdateName(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "원래이름", "ko", "")

	updated, err := s.UpdateName(ctx, p.ProfileID, "새이름")
	if err != nil {
		t.Fatalf("UpdateName 실패: %v", err)
	}
	if updated.ProfileName != "새이름" {
		t.Errorf("ProfileName = %s, want 새이름", updated.ProfileName)
	}

	// 빈 이름
	_, err = s.UpdateName(ctx, p.ProfileID, "")
	if err == nil {
		t.Error("빈 이름에 에러가 반환되어야 합니다")
	}
}

func TestSynthesize_InvalidLanguage(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "언어검증", "ko", "")
	s.UpdateStatus(ctx, p.ProfileID, "ready")

	_, _, err := s.SynthesizeTranslation(ctx, p.ProfileID, "테스트", "xx")
	if err == nil {
		t.Error("지원하지 않는 언어로 합성 시 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase F 추가 테스트
// ============================================================================

func TestCreate_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.Create(context.Background(), "", "프로필", "ko", "")
	if err == nil {
		t.Error("빈 userID에 에러가 반환되어야 합니다")
	}
}

func TestCreate_EmptyName(t *testing.T) {
	s := newSvc()
	_, err := s.Create(context.Background(), "u1", "", "ko", "")
	if err == nil {
		t.Error("빈 프로필 이름에 에러가 반환되어야 합니다")
	}
}

func TestCreate_DefaultLanguage(t *testing.T) {
	s := newSvc()
	p, err := s.Create(context.Background(), "u1", "기본언어", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Language != "ko" {
		t.Errorf("기본 언어 ko 기대: got %s", p.Language)
	}
}

func TestCreate_NameTooLong(t *testing.T) {
	s := newSvc()
	longName := ""
	for i := 0; i <= service.MaxProfileNameLen; i++ {
		longName += "가"
	}
	_, err := s.Create(context.Background(), "u1", longName, "ko", "")
	if err == nil {
		t.Error("긴 이름에 에러가 반환되어야 합니다")
	}
}

func TestGet_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.Get(context.Background(), "")
	if err == nil {
		t.Error("빈 profileID에 에러가 반환되어야 합니다")
	}
}

func TestList_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.List(context.Background(), "")
	if err == nil {
		t.Error("빈 userID에 에러가 반환되어야 합니다")
	}
}

func TestDelete_EmptyID(t *testing.T) {
	s := newSvc()
	err := s.Delete(context.Background(), "")
	if err == nil {
		t.Error("빈 profileID에 에러가 반환되어야 합니다")
	}
}

func TestSynthesize_EmptyText(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "합성테스트", "ko", "")
	s.UpdateStatus(ctx, p.ProfileID, "ready")

	_, _, err := s.SynthesizeTranslation(ctx, p.ProfileID, "", "en")
	if err == nil {
		t.Error("빈 텍스트로 합성 시 에러가 반환되어야 합니다")
	}
}

func TestCreate_SupportedLanguages(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	langs := []string{"ko", "en", "ja", "zh", "es", "fr", "de"}
	for _, lang := range langs {
		p, err := s.Create(ctx, "u-lang-"+lang, "프로필_"+lang, lang, "")
		if err != nil {
			t.Errorf("언어 %s 생성 실패: %v", lang, err)
		}
		if p.Language != lang {
			t.Errorf("언어 %s 설정 후 불일치: got %s", lang, p.Language)
		}
	}
}

func TestUpdateName_TooLong(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "원래이름", "ko", "")

	longName := ""
	for i := 0; i <= service.MaxProfileNameLen; i++ {
		longName += "나"
	}
	_, err := s.UpdateName(ctx, p.ProfileID, longName)
	if err == nil {
		t.Error("긴 이름으로 업데이트 시 에러가 반환되어야 합니다")
	}
}

func TestUpdateStatus_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.UpdateStatus(context.Background(), "", "ready")
	if err == nil {
		t.Error("빈 profileID에 에러가 반환되어야 합니다")
	}
}
