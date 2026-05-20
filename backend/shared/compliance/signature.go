package compliance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// AuditSignature는 감사 로그 항목의 전자 서명을 표현합니다 (21 CFR Part 11).
// HMAC-SHA256 기반으로 무결성(integrity)과 부인 방지(non-repudiation)를 제공합니다.
type AuditSignature struct {
	EntryID   string    // 감사 로그 항목 ID
	SignerID  string    // 서명자 ID (관리자)
	Signature string    // HMAC-SHA256 서명 (hex)
	SignedAt  time.Time // 서명 시각
	Algorithm string    // 서명 알고리즘 ("HMAC-SHA256")
}

// AuditSigner는 감사 로그 항목에 전자 서명을 수행합니다.
type AuditSigner struct {
	key []byte
}

// NewAuditSigner는 HMAC 서명 키로 서명기를 생성합니다.
// 키는 최소 32바이트를 권장합니다 (NIST SP 800-107).
func NewAuditSigner(key []byte) (*AuditSigner, error) {
	if len(key) < 16 {
		return nil, fmt.Errorf("signing key must be at least 16 bytes, got %d", len(key))
	}
	return &AuditSigner{key: key}, nil
}

// SignEntry는 감사 로그 항목에 HMAC-SHA256 서명을 생성합니다.
// payload는 서명 대상 데이터(entryID + action + timestamp + adminID 등을 결합한 문자열)입니다.
func (s *AuditSigner) SignEntry(entryID, signerID, payload string) (*AuditSignature, error) {
	if entryID == "" || signerID == "" {
		return nil, fmt.Errorf("entryID and signerID are required")
	}
	if payload == "" {
		return nil, fmt.Errorf("payload is required for signing")
	}

	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))

	return &AuditSignature{
		EntryID:   entryID,
		SignerID:  signerID,
		Signature: sig,
		SignedAt:  time.Now(),
		Algorithm: "HMAC-SHA256",
	}, nil
}

// VerifySignature는 감사 로그 항목의 서명이 유효한지 검증합니다.
// 21 CFR Part 11 §11.10(e): 감사 로그는 변조 불가능해야 합니다.
func (s *AuditSigner) VerifySignature(signature *AuditSignature, payload string) (bool, error) {
	if signature == nil {
		return false, fmt.Errorf("signature is nil")
	}
	if payload == "" {
		return false, fmt.Errorf("payload is required for verification")
	}

	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature.Signature), []byte(expected)), nil
}

// BuildAuditPayload는 감사 로그 필드를 결합하여 서명 대상 문자열을 생성합니다.
// 필드 순서가 중요합니다 — 항상 동일한 순서로 결합합니다.
func BuildAuditPayload(entryID, adminID, action, resourceType, resourceID, description string, timestamp time.Time) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		entryID, adminID, action, resourceType, resourceID, description,
		timestamp.UTC().Format(time.RFC3339Nano))
}

// AuditChain은 해시 체인 기반 감사 로그 무결성을 보장합니다.
// 각 항목의 해시가 이전 항목의 해시를 포함하므로 중간 삽입/삭제/수정이 불가능합니다.
type AuditChain struct {
	lastHash string
}

// NewAuditChain은 초기 해시(genesis)로 체인을 생성합니다.
func NewAuditChain(genesisHash string) *AuditChain {
	if genesisHash == "" {
		genesisHash = "0000000000000000000000000000000000000000000000000000000000000000"
	}
	return &AuditChain{lastHash: genesisHash}
}

// AppendEntry는 새 감사 로그 항목을 체인에 추가하고 해시를 반환합니다.
// 반환된 해시는 DB에 저장되어야 합니다.
func (c *AuditChain) AppendEntry(payload string) string {
	data := c.lastHash + "|" + payload
	hash := sha256.Sum256([]byte(data))
	c.lastHash = hex.EncodeToString(hash[:])
	return c.lastHash
}

// VerifyChain은 감사 로그 체인의 무결성을 검증합니다.
// entries는 (payload, storedHash) 쌍 목록이며, 순서대로 검증합니다.
func VerifyChain(genesisHash string, entries []ChainEntry) (bool, int) {
	if genesisHash == "" {
		genesisHash = "0000000000000000000000000000000000000000000000000000000000000000"
	}

	prevHash := genesisHash
	for i, entry := range entries {
		data := prevHash + "|" + entry.Payload
		hash := sha256.Sum256([]byte(data))
		computed := hex.EncodeToString(hash[:])
		if computed != entry.StoredHash {
			return false, i // 위조된 항목 인덱스 반환
		}
		prevHash = computed
	}
	return true, -1
}

// ChainEntry는 체인 검증에 사용되는 감사 로그 항목입니다.
type ChainEntry struct {
	Payload    string // BuildAuditPayload() 결과
	StoredHash string // DB에 저장된 해시
}

// LastHash는 현재 체인의 마지막 해시를 반환합니다.
func (c *AuditChain) LastHash() string {
	return c.lastHash
}
