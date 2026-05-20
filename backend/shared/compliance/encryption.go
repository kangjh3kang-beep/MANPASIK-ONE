// Package compliance는 규제 준수(GDPR/HIPAA/MFDS)를 위한 공통 기능을 제공합니다.
package compliance

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// PHIEncryptor는 개인건강정보(PHI) 필드를 AES-256-GCM으로 암호화/복호화합니다.
type PHIEncryptor struct {
	gcm cipher.AEAD
}

// NewPHIEncryptor는 32바이트 키로 암호화기를 생성합니다.
// 키가 nil이면 에러를 반환합니다. 환경변수 PHI_ENCRYPTION_KEY 에서 로드를 권장합니다.
func NewPHIEncryptor(key []byte) (*PHIEncryptor, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("AES-256 key must be 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	return &PHIEncryptor{gcm: gcm}, nil
}

// Encrypt는 평문을 AES-256-GCM으로 암호화하여 base64 인코딩된 문자열을 반환합니다.
// nonce는 매 호출마다 crypto/rand로 생성됩니다.
func (e *PHIEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt는 base64 인코딩된 암호문을 복호화하여 평문을 반환합니다.
func (e *PHIEncryptor) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	nonceSize := e.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptFields는 맵에서 지정된 필드들만 암호화합니다.
// PHI 필드 선별 암호화에 사용됩니다.
func (e *PHIEncryptor) EncryptFields(data map[string]string, fields []string) (map[string]string, error) {
	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = v
	}

	for _, f := range fields {
		if val, ok := result[f]; ok && val != "" {
			encrypted, err := e.Encrypt(val)
			if err != nil {
				return nil, fmt.Errorf("encrypt field %s: %w", f, err)
			}
			result[f] = encrypted
		}
	}

	return result, nil
}

// DecryptFields는 맵에서 지정된 필드들만 복호화합니다.
func (e *PHIEncryptor) DecryptFields(data map[string]string, fields []string) (map[string]string, error) {
	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = v
	}

	for _, f := range fields {
		if val, ok := result[f]; ok && val != "" {
			decrypted, err := e.Decrypt(val)
			if err != nil {
				return nil, fmt.Errorf("decrypt field %s: %w", f, err)
			}
			result[f] = decrypted
		}
	}

	return result, nil
}

// DefaultPHIFields는 HIPAA/GDPR 기준 기본 PHI 필드 목록입니다.
func DefaultPHIFields() []string {
	return []string{
		"name",
		"email",
		"phone",
		"date_of_birth",
		"address",
		"social_security_number",
		"medical_record_number",
		"health_insurance_id",
		"ip_address",
		"device_serial",
	}
}
