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
