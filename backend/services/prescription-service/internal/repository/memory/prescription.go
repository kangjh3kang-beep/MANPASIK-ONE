// Package memory는 prescription-service의 인메모리 저장소를 구현합니다.
package memory

import (
	"context"
	"sync"
	"time"

	apperrors "github.com/manpasik/backend/shared/errors"
	"github.com/manpasik/backend/services/prescription-service/internal/service"
)

// PrescriptionRepository는 인메모리 처방전 저장소입니다.
type PrescriptionRepository struct {
	mu            sync.RWMutex
	prescriptions map[string]*service.Prescription
}

// NewPrescriptionRepository는 PrescriptionRepository를 생성합니다.
func NewPrescriptionRepository() *PrescriptionRepository {
	return &PrescriptionRepository{
		prescriptions: make(map[string]*service.Prescription),
	}
}

func (r *PrescriptionRepository) Save(_ context.Context, p *service.Prescription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prescriptions[p.ID] = p
	return nil
}

func (r *PrescriptionRepository) FindByID(_ context.Context, id string) (*service.Prescription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.prescriptions[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "처방전을 찾을 수 없습니다")
	}
	return p, nil
}

func (r *PrescriptionRepository) FindByUserID(_ context.Context, userID string, statusFilter service.PrescriptionStatus, limit, offset int) ([]*service.Prescription, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Prescription
	for _, p := range r.prescriptions {
		if p.PatientUserID != userID {
			continue
		}
		if statusFilter != service.StatusUnknown && p.Status != statusFilter {
			continue
		}
		filtered = append(filtered, p)
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

func (r *PrescriptionRepository) Update(_ context.Context, p *service.Prescription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.prescriptions[p.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "처방전을 찾을 수 없습니다")
	}
	r.prescriptions[p.ID] = p
	return nil
}

// FindByPharmacyID는 약국 ID로 처방전 목록을 조회합니다.
func (r *PrescriptionRepository) FindByPharmacyID(_ context.Context, pharmacyID string) ([]*service.Prescription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.Prescription
	for _, p := range r.prescriptions {
		if p.PharmacyID == pharmacyID {
			result = append(result, p)
		}
	}
	return result, nil
}

// FindByFulfillmentToken은 조제 토큰으로 처방전을 조회합니다.
func (r *PrescriptionRepository) FindByFulfillmentToken(_ context.Context, token string) (*service.Prescription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.prescriptions {
		if p.FulfillmentToken == token {
			return p, nil
		}
	}
	return nil, apperrors.New(apperrors.ErrNotFound, "해당 토큰의 처방전을 찾을 수 없습니다")
}

// TokenRepository는 인메모리 조제 토큰 저장소입니다.
type TokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]*service.FulfillmentToken
}

// NewTokenRepository는 TokenRepository를 생성합니다.
func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		tokens: make(map[string]*service.FulfillmentToken),
	}
}

func (r *TokenRepository) Create(_ context.Context, token *service.FulfillmentToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.Token] = token
	return nil
}

func (r *TokenRepository) GetByToken(_ context.Context, token string) (*service.FulfillmentToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ft, ok := r.tokens[token]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "토큰을 찾을 수 없습니다")
	}
	return ft, nil
}

func (r *TokenRepository) MarkUsed(_ context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ft, ok := r.tokens[token]
	if !ok {
		return apperrors.New(apperrors.ErrNotFound, "토큰을 찾을 수 없습니다")
	}
	ft.IsUsed = true
	ft.UsedAt = time.Now()
	return nil
}

// DrugInteractionRepository는 인메모리 약물 상호작용 저장소입니다.
type DrugInteractionRepository struct {
	mu           sync.RWMutex
	interactions []*service.DrugInteraction
}

// NewDrugInteractionRepository는 DrugInteractionRepository를 생성합니다 (시드 데이터 포함).
func NewDrugInteractionRepository() *DrugInteractionRepository {
	repo := &DrugInteractionRepository{
		interactions: []*service.DrugInteraction{
			{DrugA: "WARF001", DrugB: "ASPR001", Severity: service.SeverityMajor, Description: "와파린과 아스피린의 병용은 출혈 위험을 크게 증가시킵니다", Recommendation: "병용 시 INR 모니터링 강화, 출혈 징후 관찰"},
			{DrugA: "WARF001", DrugB: "IBUP001", Severity: service.SeverityMajor, Description: "와파린과 이부프로펜은 출혈 위험을 증가시킵니다", Recommendation: "NSAIDs 대신 아세트아미노펜 사용 고려"},
			{DrugA: "METF001", DrugB: "ALCO001", Severity: service.SeverityModerate, Description: "메트포르민과 알코올은 젖산산증 위험을 증가시킵니다", Recommendation: "음주 제한 권고"},
			{DrugA: "LISP001", DrugB: "ACEI001", Severity: service.SeverityModerate, Description: "리시노프릴과 ACE 억제제 병용 시 고칼륨혈증 위험", Recommendation: "혈중 칼륨 수치 모니터링"},
			{DrugA: "SSRI001", DrugB: "MAOI001", Severity: service.SeverityContraindicated, Description: "SSRI와 MAOI 병용은 세로토닌 증후군을 유발할 수 있습니다", Recommendation: "절대 병용 금지, MAOI 중단 후 최소 14일 경과 후 SSRI 투여"},
			{DrugA: "SIMV001", DrugB: "CLAR001", Severity: service.SeverityMajor, Description: "심바스타틴과 클래리스로마이신 병용 시 횡문근융해증 위험", Recommendation: "클래리스로마이신 사용 중 심바스타틴 일시 중단"},
			{DrugA: "DIGO001", DrugB: "AMIO001", Severity: service.SeverityMajor, Description: "디곡신과 아미오다론 병용 시 디곡신 혈중 농도 상승", Recommendation: "디곡신 용량 50% 감량, 혈중 농도 모니터링"},
		},
	}
	return repo
}

func (r *DrugInteractionRepository) CheckInteractions(_ context.Context, drugCodes []string) ([]*service.DrugInteraction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	codeSet := make(map[string]bool)
	for _, code := range drugCodes {
		codeSet[code] = true
	}

	var found []*service.DrugInteraction
	for _, interaction := range r.interactions {
		if codeSet[interaction.DrugA] && codeSet[interaction.DrugB] {
			found = append(found, interaction)
		}
	}

	return found, nil
}
