package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/marketplace-service/internal/repository/memory"
	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

func setupTestService() *service.MarketplaceService {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	return service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
}

func TestListPartnerProducts_All(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	products, err := svc.ListPartnerProducts(ctx, "", "")
	if err != nil {
		t.Fatalf("전체 상품 조회 실패: %v", err)
	}
	if len(products) != 4 {
		t.Fatalf("시드 상품 수 불일치: got %d, want 4", len(products))
	}
}

func TestListPartnerProducts_ByPartner(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	products, err := svc.ListPartnerProducts(ctx, "partner-001", "")
	if err != nil {
		t.Fatalf("파트너별 상품 조회 실패: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("partner-001 상품 수 불일치: got %d, want 2", len(products))
	}
}

func TestListPartnerProducts_ByCategory(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	products, err := svc.ListPartnerProducts(ctx, "", "device")
	if err != nil {
		t.Fatalf("카테고리별 상품 조회 실패: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("device 카테고리 상품 수 불일치: got %d, want 2", len(products))
	}
}

func TestListPartnerProducts_ByPartnerAndCategory(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	products, err := svc.ListPartnerProducts(ctx, "partner-001", "device")
	if err != nil {
		t.Fatalf("파트너+카테고리 상품 조회 실패: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("partner-001 device 상품 수 불일치: got %d, want 2", len(products))
	}
}

func TestRegisterPartner_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	partner := &service.Partner{
		Name:         "테스트파트너",
		Description:  "테스트 파트너 설명",
		ContactEmail: "test@partner.com",
	}

	id, err := svc.RegisterPartner(ctx, partner)
	if err != nil {
		t.Fatalf("파트너 등록 실패: %v", err)
	}
	if id == "" {
		t.Fatal("파트너 ID가 비어 있음")
	}
}

func TestRegisterPartner_NilPartner(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.RegisterPartner(ctx, nil)
	if err == nil {
		t.Fatal("nil 파트너에 에러가 반환되어야 함")
	}
}

func TestRegisterPartner_MissingName(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	partner := &service.Partner{ContactEmail: "test@email.com"}
	_, err := svc.RegisterPartner(ctx, partner)
	if err == nil {
		t.Fatal("빈 이름에 에러가 반환되어야 함")
	}
}

func TestRegisterPartner_MissingEmail(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	partner := &service.Partner{Name: "테스트"}
	_, err := svc.RegisterPartner(ctx, partner)
	if err == nil {
		t.Fatal("빈 이메일에 에러가 반환되어야 함")
	}
}

func TestGetPartnerStats_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	stats, err := svc.GetPartnerStats(ctx, "partner-001")
	if err != nil {
		t.Fatalf("파트너 통계 조회 실패: %v", err)
	}
	if stats == nil {
		t.Fatal("통계가 nil이면 안 됨")
	}
	if stats.TotalProducts != 2 {
		t.Fatalf("TotalProducts 불일치: got %d, want 2", stats.TotalProducts)
	}
	if stats.TotalOrders != 150 {
		t.Fatalf("TotalOrders 불일치: got %d, want 150", stats.TotalOrders)
	}
}

func TestGetPartnerStats_MissingPartnerID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetPartnerStats(ctx, "")
	if err == nil {
		t.Fatal("빈 partner_id에 에러가 반환되어야 함")
	}
}

func TestGetPartnerStats_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetPartnerStats(ctx, "partner-nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 파트너에 에러가 반환되어야 함")
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	product := &service.PartnerProduct{
		ID:          "prod-001",
		PartnerID:   "partner-001",
		Name:        "업데이트된 혈압계",
		Description: "업데이트된 설명",
		Price:       99000,
		Category:    "device",
		IsActive:    true,
	}

	ok, err := svc.UpdateProduct(ctx, product)
	if err != nil {
		t.Fatalf("상품 업데이트 실패: %v", err)
	}
	if !ok {
		t.Fatal("UpdateProduct가 true를 반환해야 함")
	}
}

func TestUpdateProduct_NilProduct(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.UpdateProduct(ctx, nil)
	if err == nil {
		t.Fatal("nil 상품에 에러가 반환되어야 함")
	}
}

func TestUpdateProduct_MissingID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	product := &service.PartnerProduct{Name: "테스트"}
	_, err := svc.UpdateProduct(ctx, product)
	if err == nil {
		t.Fatal("빈 ID에 에러가 반환되어야 함")
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	product := &service.PartnerProduct{ID: "prod-nonexistent", Name: "없는상품"}
	_, err := svc.UpdateProduct(ctx, product)
	if err == nil {
		t.Fatal("존재하지 않는 상품에 에러가 반환되어야 함")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestRegisterPartner_DefaultPending(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	partner := &service.Partner{
		Name:         "신규파트너",
		ContactEmail: "new@partner.com",
	}
	id, err := svc.RegisterPartner(ctx, partner)
	if err != nil {
		t.Fatalf("등록 실패: %v", err)
	}

	// 파트너 상태 확인
	found, _ := partnerRepo.FindByID(ctx, id)
	if found == nil {
		t.Fatal("등록된 파트너를 찾을 수 없음")
	}
	if found.Status != "pending" {
		t.Fatalf("신규 파트너 상태가 pending이어야 함: got %s", found.Status)
	}
}

func TestApprovePartner_Success(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	// pending 파트너 등록
	partner := &service.Partner{Name: "테스트", ContactEmail: "t@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)

	// 승인
	err := svc.ApprovePartner(ctx, id, "admin-1")
	if err != nil {
		t.Fatalf("승인 실패: %v", err)
	}

	found, _ := partnerRepo.FindByID(ctx, id)
	if found.Status != "active" {
		t.Fatalf("승인 후 상태가 active여야 함: got %s", found.Status)
	}
}

func TestApprovePartner_AlreadyActive(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// partner-001은 시드 데이터에서 이미 active
	err := svc.ApprovePartner(ctx, "partner-001", "admin-1")
	if err == nil {
		t.Fatal("이미 active인 파트너 승인에 에러가 반환되어야 함")
	}
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	product := &service.PartnerProduct{
		PartnerID: "partner-001",
		Name:      "테스트 상품",
		Price:     0, // 가격 0
		Category:  "device",
	}
	_, err := svc.CreateProduct(ctx, product)
	if err == nil {
		t.Fatal("가격 0인 상품에 에러가 반환되어야 함")
	}

	product.Price = -1000
	_, err = svc.CreateProduct(ctx, product)
	if err == nil {
		t.Fatal("음수 가격인 상품에 에러가 반환되어야 함")
	}
}

func TestCreateProduct_PartnerNotActive(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	// pending 파트너 등록
	partner := &service.Partner{Name: "미승인파트너", ContactEmail: "p@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)

	product := &service.PartnerProduct{
		PartnerID: id,
		Name:      "신규 상품",
		Price:     50000,
		Category:  "device",
	}
	_, err := svc.CreateProduct(ctx, product)
	if err == nil {
		t.Fatal("pending 파트너로 상품 생성 시 에러가 반환되어야 함")
	}
}

func TestCalculateCommission(t *testing.T) {
	svc := setupTestService()

	tests := []struct {
		category string
		price    float64
		expected float64
	}{
		{"cartridge", 100000, 15000},
		{"device", 100000, 10000},
		{"accessory", 100000, 20000},
		{"subscription", 100000, 25000},
		{"service", 100000, 20000},
	}

	for _, tt := range tests {
		product := &service.PartnerProduct{
			Category: tt.category,
			Price:    tt.price,
		}
		commission := svc.CalculateCommission(product)
		if commission != tt.expected {
			t.Errorf("카테고리 %s 수수료 불일치: got %.0f, want %.0f", tt.category, commission, tt.expected)
		}
	}
}

// =============================================================================
// Phase F 테스트 보강
// =============================================================================

func TestSuspendPartner_Success(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	// pending → active → suspended
	partner := &service.Partner{Name: "중지대상", ContactEmail: "s@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)
	_ = svc.ApprovePartner(ctx, id, "admin-1")

	err := svc.SuspendPartner(ctx, id, "규정 위반")
	if err != nil {
		t.Fatalf("SuspendPartner 실패: %v", err)
	}
	found, _ := partnerRepo.FindByID(ctx, id)
	if found.Status != "suspended" {
		t.Errorf("상태: got %s, want suspended", found.Status)
	}
}

func TestSuspendPartner_AlreadySuspended(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	partner := &service.Partner{Name: "이미중지", ContactEmail: "s2@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)
	_ = svc.ApprovePartner(ctx, id, "admin-1")
	_ = svc.SuspendPartner(ctx, id, "첫 중지")

	err := svc.SuspendPartner(ctx, id, "재중지")
	if err == nil {
		t.Error("이미 suspended인 파트너 재중지에 에러가 반환되어야 합니다")
	}
}

func TestSuspendPartner_PendingPartner(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	partner := &service.Partner{Name: "미승인파트너", ContactEmail: "p@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)

	err := svc.SuspendPartner(ctx, id, "미승인 중지")
	if err == nil {
		t.Error("pending 파트너 중지에 에러가 반환되어야 합니다")
	}
}

func TestCreateProduct_Success(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	partner := &service.Partner{Name: "활성파트너", ContactEmail: "a@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)
	_ = svc.ApprovePartner(ctx, id, "admin-1")

	product := &service.PartnerProduct{
		PartnerID: id,
		Name:      "새로운 혈당계",
		Price:     89000,
		Category:  "device",
	}
	created, err := svc.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct 실패: %v", err)
	}
	if !created.IsActive {
		t.Error("생성된 상품은 활성 상태여야 합니다")
	}
}

func TestCreateProduct_InvalidCategory(t *testing.T) {
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)
	ctx := context.Background()

	partner := &service.Partner{Name: "카테고리파트너", ContactEmail: "c@t.com"}
	id, _ := svc.RegisterPartner(ctx, partner)
	_ = svc.ApprovePartner(ctx, id, "admin-1")

	product := &service.PartnerProduct{
		PartnerID: id,
		Name:      "잘못된카테고리상품",
		Price:     10000,
		Category:  "invalid_category",
	}
	_, err := svc.CreateProduct(ctx, product)
	if err == nil {
		t.Error("잘못된 카테고리에 에러가 반환되어야 합니다")
	}
}

func TestCalculateCommission_NilProduct(t *testing.T) {
	svc := setupTestService()
	commission := svc.CalculateCommission(nil)
	if commission != 0 {
		t.Errorf("nil 상품 수수료: got %.0f, want 0", commission)
	}
}

func TestCalculateCommission_UnknownCategory(t *testing.T) {
	svc := setupTestService()
	product := &service.PartnerProduct{
		Category: "unknown",
		Price:    100000,
	}
	commission := svc.CalculateCommission(product)
	// 기본 수수료율 20%
	if commission != 20000 {
		t.Errorf("미지 카테고리 수수료: got %.0f, want 20000", commission)
	}
}

func TestApprovePartner_EmptyID(t *testing.T) {
	svc := setupTestService()
	err := svc.ApprovePartner(context.Background(), "", "admin-1")
	if err == nil {
		t.Error("빈 partner_id에 에러가 반환되어야 합니다")
	}
}

func TestSuspendPartner_EmptyID(t *testing.T) {
	svc := setupTestService()
	err := svc.SuspendPartner(context.Background(), "", "이유")
	if err == nil {
		t.Error("빈 partner_id에 에러가 반환되어야 합니다")
	}
}
