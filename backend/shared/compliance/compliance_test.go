package compliance_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/compliance"
)

// ============================================================================
// PHI Encryption Tests
// ============================================================================

func generateKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return key
}

func TestPHIEncryptor_EncryptDecrypt(t *testing.T) {
	enc, err := compliance.NewPHIEncryptor(generateKey())
	if err != nil {
		t.Fatal(err)
	}

	original := "홍길동"
	ciphertext, err := enc.Encrypt(original)
	if err != nil {
		t.Fatal(err)
	}
	if ciphertext == original {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted != original {
		t.Errorf("decrypted=%q, want=%q", decrypted, original)
	}
}

func TestPHIEncryptor_EmptyString(t *testing.T) {
	enc, err := compliance.NewPHIEncryptor(generateKey())
	if err != nil {
		t.Fatal(err)
	}

	ct, err := enc.Encrypt("")
	if err != nil {
		t.Fatal(err)
	}
	if ct != "" {
		t.Error("empty input should return empty output")
	}

	pt, err := enc.Decrypt("")
	if err != nil {
		t.Fatal(err)
	}
	if pt != "" {
		t.Error("empty input should return empty output")
	}
}

func TestPHIEncryptor_InvalidKey(t *testing.T) {
	_, err := compliance.NewPHIEncryptor([]byte("short"))
	if err == nil {
		t.Error("should error on invalid key length")
	}
}

func TestPHIEncryptor_WrongKey(t *testing.T) {
	enc1, _ := compliance.NewPHIEncryptor(generateKey())
	enc2, _ := compliance.NewPHIEncryptor(generateKey())

	ct, err := enc1.Encrypt("secret data")
	if err != nil {
		t.Fatal(err)
	}

	_, err = enc2.Decrypt(ct)
	if err == nil {
		t.Error("should error when decrypting with wrong key")
	}
}

func TestPHIEncryptor_EncryptDecryptFields(t *testing.T) {
	enc, _ := compliance.NewPHIEncryptor(generateKey())

	data := map[string]string{
		"name":    "홍길동",
		"email":   "hong@test.com",
		"phone":   "010-1234-5678",
		"user_id": "USR-001", // 암호화 대상 아님
	}

	fields := []string{"name", "email", "phone"}
	encrypted, err := enc.EncryptFields(data, fields)
	if err != nil {
		t.Fatal(err)
	}

	// 암호화된 필드 확인
	if encrypted["name"] == "홍길동" {
		t.Error("name should be encrypted")
	}
	if encrypted["user_id"] != "USR-001" {
		t.Error("user_id should be unchanged")
	}

	// 복호화
	decrypted, err := enc.DecryptFields(encrypted, fields)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted["name"] != "홍길동" {
		t.Errorf("name=%q, want=홍길동", decrypted["name"])
	}
	if decrypted["email"] != "hong@test.com" {
		t.Errorf("email=%q, want=hong@test.com", decrypted["email"])
	}
}

func TestPHIEncryptor_UniqueNonce(t *testing.T) {
	enc, _ := compliance.NewPHIEncryptor(generateKey())
	ct1, _ := enc.Encrypt("same data")
	ct2, _ := enc.Encrypt("same data")
	if ct1 == ct2 {
		t.Error("each encryption should produce unique ciphertext due to random nonce")
	}
}

func TestDefaultPHIFields(t *testing.T) {
	fields := compliance.DefaultPHIFields()
	if len(fields) < 5 {
		t.Errorf("expected at least 5 PHI fields, got %d", len(fields))
	}
}

// ============================================================================
// Retention Policy Tests
// ============================================================================

func TestRetentionManager_GetPolicy(t *testing.T) {
	rm := compliance.NewRetentionManager()

	policy, err := rm.GetPolicy("measurement")
	if err != nil {
		t.Fatal(err)
	}
	if policy.RetentionPeriod < 365*24*time.Hour {
		t.Error("measurement retention should be at least 1 year")
	}
	if policy.LegalBasis == "" {
		t.Error("legal basis should be specified")
	}
}

func TestRetentionManager_UnknownType(t *testing.T) {
	rm := compliance.NewRetentionManager()
	_, err := rm.GetPolicy("unknown_type")
	if err == nil {
		t.Error("should error for unknown data type")
	}
}

func TestRetentionManager_CheckRetention_NotExpired(t *testing.T) {
	rm := compliance.NewRetentionManager()
	result, err := rm.CheckRetention("measurement", "REC-001", time.Now().Add(-1*365*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsExpired {
		t.Error("1-year-old measurement should not be expired (10-year retention)")
	}
	if result.DaysRemaining <= 0 {
		t.Errorf("days remaining=%d, should be positive", result.DaysRemaining)
	}
}

func TestRetentionManager_CheckRetention_Expired(t *testing.T) {
	rm := compliance.NewRetentionManager()
	result, err := rm.CheckRetention("user_profile", "REC-002", time.Now().Add(-4*365*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsExpired {
		t.Error("4-year-old user_profile should be expired (3-year retention)")
	}
}

func TestRetentionManager_ListPolicies(t *testing.T) {
	rm := compliance.NewRetentionManager()
	policies := rm.ListPolicies()
	if len(policies) < 5 {
		t.Errorf("expected at least 5 policies, got %d", len(policies))
	}
}

// ============================================================================
// Deletion Request Tests
// ============================================================================

type fakeDeleter struct {
	dataType string
	result   *compliance.DeletionResult
	err      error
}

func (f *fakeDeleter) DeleteUserData(_ context.Context, _ string) (*compliance.DeletionResult, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}
func (f *fakeDeleter) DataType() string { return f.dataType }

func TestDeletionProcessor_FullDeletion(t *testing.T) {
	rm := compliance.NewRetentionManager()
	dp := compliance.NewDeletionProcessor(rm)

	dp.RegisterDeleter(&fakeDeleter{
		dataType: "measurement",
		result: &compliance.DeletionResult{
			DataType:     "measurement",
			RecordCount:  50,
			DeletedCount: 45,
			RetainedCount: 5,
			RetainReason:  "법적 보존 의무 (의료법 제22조)",
			Status:       compliance.DeletionPartial,
		},
	})
	dp.RegisterDeleter(&fakeDeleter{
		dataType: "user_profile",
		result: &compliance.DeletionResult{
			DataType:     "user_profile",
			RecordCount:  1,
			DeletedCount: 1,
			Status:       compliance.DeletionCompleted,
		},
	})

	req := &compliance.DeletionRequest{
		RequestID:   "DEL-001",
		UserID:      "USR-001",
		RequestedAt: time.Now(),
		Reason:      "사용자 탈퇴 요청",
	}

	err := dp.ProcessDeletion(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(req.Results) != 2 {
		t.Errorf("results=%d, want 2", len(req.Results))
	}
	if req.ProcessedAt == nil {
		t.Error("processedAt should be set")
	}
}

func TestDeletionProcessor_ScopedDeletion(t *testing.T) {
	rm := compliance.NewRetentionManager()
	dp := compliance.NewDeletionProcessor(rm)

	dp.RegisterDeleter(&fakeDeleter{
		dataType: "measurement",
		result: &compliance.DeletionResult{
			DataType:     "measurement",
			RecordCount:  10,
			DeletedCount: 10,
			Status:       compliance.DeletionCompleted,
		},
	})
	dp.RegisterDeleter(&fakeDeleter{
		dataType: "user_profile",
		result: &compliance.DeletionResult{
			DataType:     "user_profile",
			RecordCount:  1,
			DeletedCount: 1,
			Status:       compliance.DeletionCompleted,
		},
	})

	req := &compliance.DeletionRequest{
		RequestID: "DEL-002",
		UserID:    "USR-002",
		Scope:     []string{"measurement"}, // measurement만 삭제
	}

	err := dp.ProcessDeletion(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(req.Results) != 1 {
		t.Errorf("results=%d, want 1 (scoped)", len(req.Results))
	}
	if req.Status != compliance.DeletionCompleted {
		t.Errorf("status=%s, want completed", req.Status)
	}
}

func TestDeletionProcessor_NilRequest(t *testing.T) {
	rm := compliance.NewRetentionManager()
	dp := compliance.NewDeletionProcessor(rm)
	err := dp.ProcessDeletion(context.Background(), nil)
	if err == nil {
		t.Error("should error on nil request")
	}
}

func TestDeletionProcessor_EmptyUserID(t *testing.T) {
	rm := compliance.NewRetentionManager()
	dp := compliance.NewDeletionProcessor(rm)
	err := dp.ProcessDeletion(context.Background(), &compliance.DeletionRequest{})
	if err == nil {
		t.Error("should error on empty user ID")
	}
}

func TestDeletionProcessor_UnregisteredType(t *testing.T) {
	rm := compliance.NewRetentionManager()
	dp := compliance.NewDeletionProcessor(rm)

	req := &compliance.DeletionRequest{
		RequestID: "DEL-003",
		UserID:    "USR-003",
		Scope:     []string{"nonexistent"},
	}

	err := dp.ProcessDeletion(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if req.Status != compliance.DeletionPartial {
		t.Errorf("status=%s, want partial", req.Status)
	}
}

// ============================================================================
// Compliance Overview Tests
// ============================================================================

func TestGetComplianceOverview(t *testing.T) {
	overview := compliance.GetComplianceOverview()
	if len(overview) < 10 {
		t.Errorf("expected at least 10 compliance items, got %d", len(overview))
	}

	// 프레임워크 확인
	frameworks := make(map[string]bool)
	for _, s := range overview {
		frameworks[s.Framework] = true
	}
	for _, fw := range []string{"GDPR", "HIPAA", "MFDS", "FDA", "CE-IVDR"} {
		if !frameworks[fw] {
			t.Errorf("missing framework: %s", fw)
		}
	}
}

func TestGetComplianceOverview_HasNonCompliant(t *testing.T) {
	overview := compliance.GetComplianceOverview()
	hasNonCompliant := false
	for _, s := range overview {
		if s.Status == "non_compliant" {
			hasNonCompliant = true
			if len(s.Recommendations) == 0 {
				t.Errorf("non_compliant item %s/%s should have recommendations", s.Framework, s.Category)
			}
		}
	}
	if !hasNonCompliant {
		t.Error("should have at least one non_compliant item (clinical validation)")
	}
}

// ============================================================================
// 21 CFR Part 11 — 전자 서명 테스트 (Phase H)
// ============================================================================

func TestAuditSigner_SignAndVerify(t *testing.T) {
	signer, err := compliance.NewAuditSigner(generateKey())
	if err != nil {
		t.Fatal(err)
	}

	payload := compliance.BuildAuditPayload(
		"ENTRY-001", "ADMIN-001", "config_change",
		"system_config", "CFG-001", "보안 설정 변경",
		time.Now(),
	)

	sig, err := signer.SignEntry("ENTRY-001", "ADMIN-001", payload)
	if err != nil {
		t.Fatal(err)
	}
	if sig.Signature == "" {
		t.Error("signature should not be empty")
	}
	if sig.Algorithm != "HMAC-SHA256" {
		t.Errorf("algorithm=%s, want HMAC-SHA256", sig.Algorithm)
	}

	valid, err := signer.VerifySignature(sig, payload)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Error("signature should be valid")
	}
}

func TestAuditSigner_TamperedPayload(t *testing.T) {
	signer, _ := compliance.NewAuditSigner(generateKey())
	payload := compliance.BuildAuditPayload(
		"E-1", "A-1", "delete", "user", "U-1", "사용자 삭제", time.Now(),
	)

	sig, _ := signer.SignEntry("E-1", "A-1", payload)

	// 페이로드 변조
	tampered := payload + "TAMPERED"
	valid, _ := signer.VerifySignature(sig, tampered)
	if valid {
		t.Error("tampered payload should fail verification")
	}
}

func TestAuditSigner_DifferentKeys(t *testing.T) {
	signer1, _ := compliance.NewAuditSigner(generateKey())
	signer2, _ := compliance.NewAuditSigner(generateKey())

	payload := "test|payload|data"
	sig, _ := signer1.SignEntry("E-1", "A-1", payload)

	valid, _ := signer2.VerifySignature(sig, payload)
	if valid {
		t.Error("different key should fail verification")
	}
}

func TestAuditSigner_ShortKey(t *testing.T) {
	_, err := compliance.NewAuditSigner([]byte("short"))
	if err == nil {
		t.Error("should error on short key")
	}
}

func TestAuditSigner_EmptyFields(t *testing.T) {
	signer, _ := compliance.NewAuditSigner(generateKey())
	_, err := signer.SignEntry("", "A-1", "payload")
	if err == nil {
		t.Error("should error on empty entryID")
	}
	_, err = signer.SignEntry("E-1", "", "payload")
	if err == nil {
		t.Error("should error on empty signerID")
	}
	_, err = signer.SignEntry("E-1", "A-1", "")
	if err == nil {
		t.Error("should error on empty payload")
	}
}

func TestAuditSigner_NilSignature(t *testing.T) {
	signer, _ := compliance.NewAuditSigner(generateKey())
	_, err := signer.VerifySignature(nil, "payload")
	if err == nil {
		t.Error("should error on nil signature")
	}
}

// ============================================================================
// 해시 체인 감사 로그 무결성 테스트 (Phase H)
// ============================================================================

func TestAuditChain_AppendAndVerify(t *testing.T) {
	chain := compliance.NewAuditChain("")
	var entries []compliance.ChainEntry

	for i := 0; i < 5; i++ {
		payload := compliance.BuildAuditPayload(
			"E-"+string(rune('1'+i)), "A-1", "update",
			"config", "C-1", "변경",
			time.Now(),
		)
		hash := chain.AppendEntry(payload)
		entries = append(entries, compliance.ChainEntry{
			Payload:    payload,
			StoredHash: hash,
		})
	}

	valid, idx := compliance.VerifyChain("", entries)
	if !valid {
		t.Errorf("chain should be valid, failed at index %d", idx)
	}
}

func TestAuditChain_TamperDetection(t *testing.T) {
	chain := compliance.NewAuditChain("")
	var entries []compliance.ChainEntry

	for i := 0; i < 3; i++ {
		payload := "entry-" + string(rune('A'+i))
		hash := chain.AppendEntry(payload)
		entries = append(entries, compliance.ChainEntry{
			Payload:    payload,
			StoredHash: hash,
		})
	}

	// 두 번째 항목의 페이로드 변조
	entries[1].Payload = "TAMPERED"

	valid, idx := compliance.VerifyChain("", entries)
	if valid {
		t.Error("tampered chain should fail verification")
	}
	if idx != 1 {
		t.Errorf("tamper index=%d, want 1", idx)
	}
}

func TestAuditChain_EmptyChain(t *testing.T) {
	valid, idx := compliance.VerifyChain("", nil)
	if !valid {
		t.Error("empty chain should be valid")
	}
	if idx != -1 {
		t.Errorf("idx=%d, want -1", idx)
	}
}

func TestAuditChain_DeterministicHash(t *testing.T) {
	chain1 := compliance.NewAuditChain("")
	chain2 := compliance.NewAuditChain("")

	h1 := chain1.AppendEntry("same-payload")
	h2 := chain2.AppendEntry("same-payload")

	if h1 != h2 {
		t.Error("same genesis + same payload should produce same hash")
	}
}

func TestBuildAuditPayload_Deterministic(t *testing.T) {
	ts := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	p1 := compliance.BuildAuditPayload("E-1", "A-1", "update", "config", "C-1", "desc", ts)
	p2 := compliance.BuildAuditPayload("E-1", "A-1", "update", "config", "C-1", "desc", ts)
	if p1 != p2 {
		t.Error("same inputs should produce same payload")
	}
}

// ============================================================================
// ComplianceChecker 자동 검사 테스트 (Phase H)
// ============================================================================

func TestComplianceChecker_FullyCompliant(t *testing.T) {
	state := compliance.SystemState{
		PHIEncryptionEnabled:     true,
		EncryptionAlgorithm:      "AES-256-GCM",
		AuditLogEnabled:          true,
		AuditSigningEnabled:      true,
		AuditChainEnabled:        true,
		RBACEnabled:              true,
		JWTAuthEnabled:           true,
		SessionTimeout:           30 * time.Minute,
		ConsentManagementEnabled: true,
		RetentionPoliciesCount:   7,
		DeletionProcessorReady:   true,
		RegisteredDeleters:       []string{"measurement", "health_record", "user_profile"},
		SRSVersion:               "v3.0",
		RiskPlanReady:            true,
		FMEAReady:                true,
		SafetyClass:              "B",
		ClinicalTrialDone:        true,
		GMPCertified:             true,
		FDAPreSubmission:         true,
	}

	checker := compliance.NewComplianceChecker(state)
	results := checker.RunChecks()

	if len(results) < 15 {
		t.Errorf("expected at least 15 checks, got %d", len(results))
	}

	score := compliance.ComputeScore(results)
	if score.Rate < 80 {
		t.Errorf("fully compliant state should score >80%%, got %.1f%%", score.Rate)
	}
}

func TestComplianceChecker_MinimalState(t *testing.T) {
	// 모든 기능 비활성화
	checker := compliance.NewComplianceChecker(compliance.SystemState{})
	results := checker.RunChecks()

	score := compliance.ComputeScore(results)
	// predicate_device (compliant) + ce-ivdr (partial) 만 기본 점수
	if score.NonCompliant == 0 {
		t.Error("minimal state should have non_compliant items")
	}
	if score.Rate > 30 {
		t.Errorf("minimal state should score <30%%, got %.1f%%", score.Rate)
	}
}

func TestComplianceChecker_FrameworkCoverage(t *testing.T) {
	checker := compliance.NewComplianceChecker(compliance.SystemState{})
	results := checker.RunChecks()

	frameworks := make(map[string]int)
	for _, s := range results {
		frameworks[s.Framework]++
	}

	expected := []string{"GDPR", "HIPAA", "MFDS", "FDA", "CE-IVDR", "21CFR11"}
	for _, fw := range expected {
		if frameworks[fw] == 0 {
			t.Errorf("missing framework: %s", fw)
		}
	}
}

func TestComplianceChecker_Recommendations(t *testing.T) {
	// 암호화만 비활성화
	state := compliance.SystemState{
		AuditLogEnabled: true,
		RBACEnabled:     true,
		JWTAuthEnabled:  true,
	}

	checker := compliance.NewComplianceChecker(state)
	results := checker.RunChecks()

	hasRec := false
	for _, s := range results {
		if s.Status == "non_compliant" && len(s.Recommendations) > 0 {
			hasRec = true
			break
		}
	}
	if !hasRec {
		t.Error("non_compliant items should have recommendations")
	}
}

func TestComputeScore_Empty(t *testing.T) {
	score := compliance.ComputeScore(nil)
	if score.TotalItems != 0 || score.Rate != 0 {
		t.Error("empty input should return zero score")
	}
}

func TestComputeScore_Mixed(t *testing.T) {
	statuses := []compliance.ComplianceStatus{
		{Status: "compliant"},
		{Status: "compliant"},
		{Status: "partial"},
		{Status: "non_compliant"},
	}
	score := compliance.ComputeScore(statuses)
	if score.TotalItems != 4 {
		t.Errorf("total=%d, want 4", score.TotalItems)
	}
	if score.Compliant != 2 {
		t.Errorf("compliant=%d, want 2", score.Compliant)
	}
	if score.Partial != 1 {
		t.Errorf("partial=%d, want 1", score.Partial)
	}
	// (2*100 + 1*50) / 4 = 62.5
	if score.Rate != 62.5 {
		t.Errorf("rate=%.1f, want 62.5", score.Rate)
	}
}
