package service_test

import (
	"context"
	"testing"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/repository/memory"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.StoreService {
	return service.NewStoreService(zap.NewNop(), memory.NewStoreRepository())
}
func TestRegisterDeveloper(t *testing.T) {
	s := newSvc()
	d, err := s.RegisterDeveloper(context.Background(), "u1", "테스트회사", "explorer")
	if err != nil { t.Fatal(err) }
	if d.ID == "" { t.Error("empty id") }
}
func TestSubmitAndListItems(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.RegisterDeveloper(ctx, "u1", "Co", "developer")
	s.SubmitCartridge(ctx, "dev1", "카트리지A", "설명", "health", "1.0", 5000)
	s.SubmitCartridge(ctx, "dev1", "카트리지B", "설명2", "health", "1.0", 3000)
	list, total, _ := s.ListItems(ctx, "", 10, 0)
	if total != 2 || len(list) != 2 { t.Errorf("expected 2, got %d", total) }
}
func TestPurchase(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	item, _ := s.SubmitCartridge(ctx, "dev1", "X", "", "", "1.0", 1000)
	p, err := s.Purchase(ctx, "buyer1", item.ID)
	if err != nil { t.Fatal(err) }
	if p.PriceKrw != 1000 { t.Error("price mismatch") }
}
func TestReviewFlow(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	r, _ := s.SubmitForReview(ctx, "cart1", "dev1")
	if r.Status != "pending" { t.Error("expected pending") }
	approved, _ := s.ApproveCartridge(ctx, r.SubmissionID, "좋습니다")
	if approved.Status != "approved" { t.Error("expected approved") }
}
func TestCreateApiKey(t *testing.T) {
	s := newSvc()
	k, err := s.CreateApiKey(context.Background(), "dev1")
	if err != nil { t.Fatal(err) }
	if k.Key == "" || k.KeyID == "" { t.Error("empty key") }
}

// =============================================================================
// Phase F 테스트 보강
// =============================================================================

func TestRegisterDeveloper_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.RegisterDeveloper(context.Background(), "", "회사", "explorer")
	if err == nil { t.Fatal("빈 userID에 에러가 반환되어야 합니다") }
}

func TestRegisterDeveloper_DefaultTier(t *testing.T) {
	s := newSvc()
	d, err := s.RegisterDeveloper(context.Background(), "u2", "회사", "")
	if err != nil { t.Fatal(err) }
	if d.Tier != "explorer" { t.Errorf("기본 tier = %s, want explorer", d.Tier) }
}

func TestGetDeveloper(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.RegisterDeveloper(ctx, "u1", "테스트회사", "developer")
	d, err := s.GetDeveloper(ctx, "u1")
	if err != nil { t.Fatal(err) }
	if d.CompanyName != "테스트회사" { t.Errorf("CompanyName = %s", d.CompanyName) }
}

func TestGetDeveloper_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.GetDeveloper(context.Background(), "")
	if err == nil { t.Fatal("빈 userID에 에러가 반환되어야 합니다") }
}

func TestCreateApiKey_EmptyDevID(t *testing.T) {
	s := newSvc()
	_, err := s.CreateApiKey(context.Background(), "")
	if err == nil { t.Fatal("빈 devID에 에러가 반환되어야 합니다") }
}

func TestSubmitCartridge_EmptyName(t *testing.T) {
	s := newSvc()
	_, err := s.SubmitCartridge(context.Background(), "dev1", "", "", "", "1.0", 1000)
	if err == nil { t.Fatal("빈 이름에 에러가 반환되어야 합니다") }
}

func TestSubmitCartridge_EmptyDevID(t *testing.T) {
	s := newSvc()
	_, err := s.SubmitCartridge(context.Background(), "", "이름", "", "", "1.0", 1000)
	if err == nil { t.Fatal("빈 devID에 에러가 반환되어야 합니다") }
}

func TestSearchItems(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.SubmitCartridge(ctx, "dev1", "혈당 카트리지", "혈당 측정용", "health", "1.0", 5000)
	s.SubmitCartridge(ctx, "dev1", "콜레스테롤 키트", "콜레스테롤 측정", "health", "1.0", 3000)
	results, err := s.SearchItems(ctx, "혈당", 10)
	if err != nil { t.Fatal(err) }
	if len(results) != 1 { t.Errorf("검색 결과 = %d, want 1", len(results)) }
}

func TestGetItem(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	item, _ := s.SubmitCartridge(ctx, "dev1", "X", "", "", "1.0", 1000)
	got, err := s.GetItem(ctx, item.ID)
	if err != nil { t.Fatal(err) }
	if got.Name != "X" { t.Errorf("Name = %s", got.Name) }
}

func TestGetItem_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.GetItem(context.Background(), "")
	if err == nil { t.Fatal("빈 ID에 에러가 반환되어야 합니다") }
}

func TestPurchase_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.Purchase(context.Background(), "", "item1")
	if err == nil { t.Fatal("빈 userID에 에러가 반환되어야 합니다") }
}

func TestGetPurchaseHistory(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	item, _ := s.SubmitCartridge(ctx, "dev1", "X", "", "", "1.0", 1000)
	s.Purchase(ctx, "buyer1", item.ID)
	s.Purchase(ctx, "buyer1", item.ID)
	hist, err := s.GetPurchaseHistory(ctx, "buyer1", 10)
	if err != nil { t.Fatal(err) }
	if len(hist) != 2 { t.Errorf("구매 이력 = %d, want 2", len(hist)) }
}

func TestGetPurchaseHistory_EmptyUserID(t *testing.T) {
	s := newSvc()
	_, err := s.GetPurchaseHistory(context.Background(), "", 10)
	if err == nil { t.Fatal("빈 userID에 에러가 반환되어야 합니다") }
}

func TestRejectCartridge(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	r, _ := s.SubmitForReview(ctx, "cart1", "dev1")
	rejected, err := s.RejectCartridge(ctx, r.SubmissionID, "품질 미달")
	if err != nil { t.Fatal(err) }
	if rejected.Status != "rejected" { t.Errorf("Status = %s, want rejected", rejected.Status) }
	if rejected.ReviewerComment != "품질 미달" { t.Errorf("Comment = %s", rejected.ReviewerComment) }
}

func TestGetReviewStatus(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	r, _ := s.SubmitForReview(ctx, "cart1", "dev1")
	got, err := s.GetReviewStatus(ctx, r.SubmissionID)
	if err != nil { t.Fatal(err) }
	if got.Status != "pending" { t.Errorf("Status = %s", got.Status) }
}

func TestGetReviewStatus_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.GetReviewStatus(context.Background(), "")
	if err == nil { t.Fatal("빈 ID에 에러가 반환되어야 합니다") }
}

func TestSubmitForReview_EmptyCartridgeID(t *testing.T) {
	s := newSvc()
	_, err := s.SubmitForReview(context.Background(), "", "dev1")
	if err == nil { t.Fatal("빈 cartridgeID에 에러가 반환되어야 합니다") }
}

func TestListSubmissions(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.SubmitForReview(ctx, "cart1", "dev1")
	s.SubmitForReview(ctx, "cart2", "dev1")
	list, err := s.ListSubmissions(ctx, "dev1", 10)
	if err != nil { t.Fatal(err) }
	if len(list) != 2 { t.Errorf("제출 수 = %d, want 2", len(list)) }
}

func TestGetSalesReport(t *testing.T) {
	s := newSvc()
	r := s.GetSalesReport(context.Background(), "dev1", "2026-01")
	if r.DeveloperID != "dev1" { t.Errorf("DeveloperID = %s", r.DeveloperID) }
	if r.TotalRevenue <= 0 { t.Error("TotalRevenue should be > 0") }
}

func TestConfigureRevenueSplit(t *testing.T) {
	s := newSvc()
	cfg := s.ConfigureRevenueSplit(context.Background(), "dev1", 70, 30)
	if cfg.DeveloperPct != 70 || cfg.PlatformPct != 30 { t.Errorf("Split = %d/%d", cfg.DeveloperPct, cfg.PlatformPct) }
}

func TestRequestPayout(t *testing.T) {
	s := newSvc()
	p := s.RequestPayout(context.Background(), "dev1", 100000)
	if p.Status != "pending" { t.Errorf("Status = %s", p.Status) }
	if p.AmountKrw != 100000 { t.Errorf("Amount = %d", p.AmountKrw) }
}

func TestGetUsageStats(t *testing.T) {
	s := newSvc()
	st := s.GetUsageStats(context.Background(), "item1")
	if st.ItemID != "item1" { t.Errorf("ItemID = %s", st.ItemID) }
	if st.UsageCount <= 0 { t.Error("UsageCount should be > 0") }
}

func TestGetRatings(t *testing.T) {
	s := newSvc()
	r := s.GetRatings(context.Background(), "item1")
	if r.AvgRating <= 0 { t.Error("AvgRating should be > 0") }
}

func TestSubmitUserReview(t *testing.T) {
	s := newSvc()
	r := s.SubmitUserReview(context.Background(), "u1", "item1", 5, "좋아요")
	if r.Rating != 5 { t.Errorf("Rating = %d", r.Rating) }
	if r.Comment != "좋아요" { t.Errorf("Comment = %s", r.Comment) }
}

func TestGetDevAnalytics(t *testing.T) {
	s := newSvc()
	a := s.GetDevAnalytics(context.Background(), "dev1")
	if a.DeveloperID != "dev1" { t.Errorf("DeveloperID = %s", a.DeveloperID) }
	if a.TotalDownloads <= 0 { t.Error("TotalDownloads should be > 0") }
}

func TestRegisterCartridgeType(t *testing.T) {
	s := newSvc()
	def := &service.CartridgeDefinition{
		DeveloperID:      "dev1",
		Name:             "혈당 v2",
		CategoryCode:     1,
		TypeIndex:        0,
		InterfaceVersion: "phase5-v1.0", // Phase 5 SSOT에서 주입된 메타데이터
	}
	result, err := s.RegisterCartridgeType(context.Background(), def)
	if err != nil { t.Fatal(err) }
	if result.InterfaceVersion != "phase5-v1.0" {
		t.Errorf("InterfaceVersion = %s, want phase5-v1.0 (외부 SSOT 그대로 보존)", result.InterfaceVersion)
	}
	if result.AlphaDefault != 0.98 { t.Errorf("AlphaDefault = %f, want 0.98 (차분식 SSOT)", result.AlphaDefault) }
	if result.ApprovalStatus != "registered" { t.Errorf("ApprovalStatus = %s", result.ApprovalStatus) }
}

func TestRegisterCartridgeType_RequiresInterfaceVersion(t *testing.T) {
	s := newSvc()
	def := &service.CartridgeDefinition{
		DeveloperID: "dev1",
		Name:        "missing-iface",
		// InterfaceVersion 없음 → Phase 5 SSOT 메타데이터 누락
	}
	_, err := s.RegisterCartridgeType(context.Background(), def)
	if err == nil {
		t.Error("InterfaceVersion 없는 등록이 통과됨")
	}
}

func TestRegisterCartridgeType_EmptyName(t *testing.T) {
	s := newSvc()
	_, err := s.RegisterCartridgeType(context.Background(), &service.CartridgeDefinition{DeveloperID: "dev1"})
	if err == nil { t.Fatal("빈 이름에 에러가 반환되어야 합니다") }
}

func TestGetCartridgeDefinition_NotFound(t *testing.T) {
	s := newSvc()
	_, err := s.GetCartridgeDefinition(context.Background(), "nonexistent")
	if err == nil { t.Fatal("존재하지 않는 정의에 에러가 반환되어야 합니다") }
}

func TestGetCartridgeDefinition_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.GetCartridgeDefinition(context.Background(), "")
	if err == nil { t.Fatal("빈 ID에 에러가 반환되어야 합니다") }
}
