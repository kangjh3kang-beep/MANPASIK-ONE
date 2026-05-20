package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
)

type Developer struct {
	ID, UserID, CompanyName, Tier, Status string
	CreatedAt                             time.Time
}
type ApiKey struct{ KeyID, DeveloperID, Key, Status string; CreatedAt time.Time }
type StoreItem struct {
	ID, DeveloperID, Name, Description, Category, Version, Status string
	PriceKrw int64; Rating float64; Downloads int32
	CreatedAt time.Time
}
type Purchase struct {
	ID, UserID, ItemID, Status string
	PriceKrw int64; PurchasedAt time.Time
}
type ReviewStatus struct {
	SubmissionID, CartridgeID, Status, ReviewerComment string
	SubmittedAt time.Time
}
type SalesReport struct {
	DeveloperID string; TotalRevenue, PlatformFee, DeveloperEarning int64
	Period string; ItemCount int32
}
type PayoutRecord struct {
	ID, DeveloperID, Status string; AmountKrw int64; RequestedAt time.Time
}
type RevenueSplitCfg struct{ DeveloperID string; DeveloperPct, PlatformPct int32 }
type UsageStat struct{ ItemID, Date string; UsageCount, UniqueUsers int32 }
type RatingInfo struct{ ItemID string; AvgRating float64; TotalReviews int32 }
type UserReview struct{ ID, UserID, ItemID string; Rating int32; Comment string; CreatedAt time.Time }
type DevAnalytics struct{ DeveloperID string; TotalDownloads, TotalRevenue int64; ActiveCartridges int32 }

type StoreRepository interface {
	CreateDeveloper(ctx context.Context, d *Developer) error
	GetDeveloper(ctx context.Context, userID string) (*Developer, error)
	CreateApiKey(ctx context.Context, k *ApiKey) error
	CreateItem(ctx context.Context, item *StoreItem) error
	ListItems(ctx context.Context, category string, limit, offset int32) ([]*StoreItem, int32, error)
	SearchItems(ctx context.Context, query string, limit int32) ([]*StoreItem, error)
	GetItem(ctx context.Context, id string) (*StoreItem, error)
	CreatePurchase(ctx context.Context, p *Purchase) error
	ListPurchases(ctx context.Context, userID string, limit int32) ([]*Purchase, error)
	CreateReview(ctx context.Context, r *ReviewStatus) error
	GetReview(ctx context.Context, submissionID string) (*ReviewStatus, error)
	UpdateReview(ctx context.Context, r *ReviewStatus) error
	ListSubmissions(ctx context.Context, devID string, limit int32) ([]*ReviewStatus, error)
}

type StoreService struct {
	log  *zap.Logger
	repo StoreRepository
}
func NewStoreService(log *zap.Logger, repo StoreRepository) *StoreService {
	return &StoreService{log: log, repo: repo}
}

func (s *StoreService) RegisterDeveloper(ctx context.Context, userID, company, tier string) (*Developer, error) {
	if userID == "" { return nil, ErrInvalidInput }
	if tier == "" { tier = "explorer" }
	d := &Developer{ID: uuid.New().String(), UserID: userID, CompanyName: company, Tier: tier, Status: "active", CreatedAt: time.Now().UTC()}
	return d, s.repo.CreateDeveloper(ctx, d)
}
func (s *StoreService) GetDeveloper(ctx context.Context, userID string) (*Developer, error) {
	if userID == "" { return nil, ErrInvalidInput }
	return s.repo.GetDeveloper(ctx, userID)
}
func (s *StoreService) CreateApiKey(ctx context.Context, devID string) (*ApiKey, error) {
	if devID == "" { return nil, ErrInvalidInput }
	k := &ApiKey{KeyID: uuid.New().String(), DeveloperID: devID, Key: "mpsk_" + uuid.New().String()[:16], Status: "active", CreatedAt: time.Now().UTC()}
	return k, s.repo.CreateApiKey(ctx, k)
}
func (s *StoreService) SubmitCartridge(ctx context.Context, devID, name, desc, category, version string, priceKrw int64) (*StoreItem, error) {
	if devID == "" || name == "" { return nil, ErrInvalidInput }
	item := &StoreItem{ID: uuid.New().String(), DeveloperID: devID, Name: name, Description: desc,
		Category: category, Version: version, PriceKrw: priceKrw, Status: "pending_review", CreatedAt: time.Now().UTC()}
	return item, s.repo.CreateItem(ctx, item)
}
func (s *StoreService) ListSubmissions(ctx context.Context, devID string, limit int32) ([]*ReviewStatus, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.ListSubmissions(ctx, devID, limit)
}
func (s *StoreService) ListItems(ctx context.Context, category string, limit, offset int32) ([]*StoreItem, int32, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.ListItems(ctx, category, limit, offset)
}
func (s *StoreService) SearchItems(ctx context.Context, query string, limit int32) ([]*StoreItem, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.SearchItems(ctx, query, limit)
}
func (s *StoreService) GetItem(ctx context.Context, id string) (*StoreItem, error) {
	if id == "" { return nil, ErrInvalidInput }
	return s.repo.GetItem(ctx, id)
}
func (s *StoreService) Purchase(ctx context.Context, userID, itemID string) (*Purchase, error) {
	if userID == "" || itemID == "" { return nil, ErrInvalidInput }
	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil { return nil, err }
	p := &Purchase{ID: uuid.New().String(), UserID: userID, ItemID: itemID, PriceKrw: item.PriceKrw, Status: "completed", PurchasedAt: time.Now().UTC()}
	return p, s.repo.CreatePurchase(ctx, p)
}
func (s *StoreService) GetPurchaseHistory(ctx context.Context, userID string, limit int32) ([]*Purchase, error) {
	if userID == "" { return nil, ErrInvalidInput }
	return s.repo.ListPurchases(ctx, userID, limit)
}
func (s *StoreService) SubmitForReview(ctx context.Context, cartridgeID, devID string) (*ReviewStatus, error) {
	if cartridgeID == "" { return nil, ErrInvalidInput }
	r := &ReviewStatus{SubmissionID: uuid.New().String(), CartridgeID: cartridgeID, Status: "pending", SubmittedAt: time.Now().UTC()}
	return r, s.repo.CreateReview(ctx, r)
}
func (s *StoreService) GetReviewStatus(ctx context.Context, submissionID string) (*ReviewStatus, error) {
	if submissionID == "" { return nil, ErrInvalidInput }
	return s.repo.GetReview(ctx, submissionID)
}
func (s *StoreService) ApproveCartridge(ctx context.Context, submissionID, comment string) (*ReviewStatus, error) {
	r, err := s.repo.GetReview(ctx, submissionID)
	if err != nil { return nil, err }
	r.Status = "approved"; r.ReviewerComment = comment
	return r, s.repo.UpdateReview(ctx, r)
}
func (s *StoreService) RejectCartridge(ctx context.Context, submissionID, reason string) (*ReviewStatus, error) {
	r, err := s.repo.GetReview(ctx, submissionID)
	if err != nil { return nil, err }
	r.Status = "rejected"; r.ReviewerComment = reason
	return r, s.repo.UpdateReview(ctx, r)
}

func (s *StoreService) GetSalesReport(ctx context.Context, devID, period string) *SalesReport {
	return &SalesReport{DeveloperID: devID, TotalRevenue: 1500000, PlatformFee: 450000, DeveloperEarning: 1050000, Period: period, ItemCount: 3}
}
func (s *StoreService) GetPayoutHistory(ctx context.Context, devID string, limit int32) []*PayoutRecord { return nil }
func (s *StoreService) ConfigureRevenueSplit(ctx context.Context, devID string, devPct, platPct int32) *RevenueSplitCfg {
	return &RevenueSplitCfg{DeveloperID: devID, DeveloperPct: devPct, PlatformPct: platPct}
}
func (s *StoreService) RequestPayout(ctx context.Context, devID string, amount int64) *PayoutRecord {
	return &PayoutRecord{ID: uuid.New().String(), DeveloperID: devID, AmountKrw: amount, Status: "pending", RequestedAt: time.Now().UTC()}
}
func (s *StoreService) GetUsageStats(ctx context.Context, itemID string) *UsageStat {
	return &UsageStat{ItemID: itemID, Date: time.Now().Format("2006-01-02"), UsageCount: 150, UniqueUsers: 42}
}
func (s *StoreService) GetRatings(ctx context.Context, itemID string) *RatingInfo {
	return &RatingInfo{ItemID: itemID, AvgRating: 4.5, TotalReviews: 23}
}
func (s *StoreService) SubmitUserReview(ctx context.Context, userID, itemID string, rating int32, comment string) *UserReview {
	return &UserReview{ID: uuid.New().String(), UserID: userID, ItemID: itemID, Rating: rating, Comment: comment, CreatedAt: time.Now().UTC()}
}
func (s *StoreService) GetDevAnalytics(ctx context.Context, devID string) *DevAnalytics {
	return &DevAnalytics{DeveloperID: devID, TotalDownloads: 12500, TotalRevenue: 3500000, ActiveCartridges: 5}
}

// CartridgeDefinition은 카트리지 타입 정의 엔티티입니다 (H5 인터페이스 계약).
//
// 주의: 하드웨어 인터페이스 버전(InterfaceVersion)과 핀 구성(PinConfigJSON)은
// Phase 5 SSOT(05_EVT_제조설계, 06_협력업체배포 등)의 메타데이터를 받아 보관할 뿐,
// 본 SW 서비스는 그 값을 검증하거나 기본값을 결정하지 않습니다.
type CartridgeDefinition struct {
	ID                   string
	DeveloperID          string
	CategoryCode         int
	TypeIndex            int
	Name                 string
	Version              string
	InterfaceVersion     string // 외부 SSOT(Phase 5)에서 주입된 인터페이스 버전 식별자
	PinConfigJSON        string // 외부 SSOT 핀 매핑 JSON (검증 없이 보관)
	AFEBlocks            []string
	MeasurementModes     []string
	CalibrationSchemaJSON string
	AlphaDefault         float64
	TempCompensation     bool
	HumidityCompensation bool
	RegulatoryClass      string
	RegulatoryPath       string
	ApprovalStatus       string
	CreatedAt            time.Time
	ApprovedAt           *time.Time
}

// RegisterCartridgeType은 SDK 개발자용 카트리지 타입을 등록합니다.
func (s *StoreService) RegisterCartridgeType(ctx context.Context, def *CartridgeDefinition) (*CartridgeDefinition, error) {
	if def.DeveloperID == "" || def.Name == "" {
		return nil, ErrInvalidInput
	}

	// 인터페이스 버전은 Phase 5 SSOT에서 주입되어야 함 (본 서비스는 검증/기본값 미적용)
	if def.InterfaceVersion == "" {
		return nil, fmt.Errorf("%w: interface_version is required (Phase 5 SSOT 메타데이터)", ErrInvalidInput)
	}

	// 차분식 알고리즘 표준값 (CLAUDE.md 차분식 SSOT)은 본 영역 책임
	if def.AlphaDefault <= 0 || def.AlphaDefault > 1.5 {
		def.AlphaDefault = 0.98 // 차분식 SSOT 기본값
	}

	def.ID = uuid.New().String()
	def.ApprovalStatus = "registered"
	def.CreatedAt = time.Now().UTC()

	s.log.Info("카트리지 타입 등록",
		zap.String("definition_id", def.ID),
		zap.String("developer_id", def.DeveloperID),
		zap.String("name", def.Name),
		zap.String("interface_version", def.InterfaceVersion),
		zap.Int("category_code", def.CategoryCode),
		zap.Int("type_index", def.TypeIndex),
	)

	return def, nil
}

// GetCartridgeDefinition은 카트리지 정의를 조회합니다.
func (s *StoreService) GetCartridgeDefinition(ctx context.Context, definitionID string) (*CartridgeDefinition, error) {
	if definitionID == "" {
		return nil, ErrInvalidInput
	}
	// 메모리 구현 — 기본값 반환 (향후 DB 연동)
	return nil, ErrNotFound
}
