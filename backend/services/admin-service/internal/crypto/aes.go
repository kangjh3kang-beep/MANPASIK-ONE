// Package crypto는 AES-256-GCM 기반 암호화/복호화를 제공합니다.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// AESEncryptor는 AES-256-GCM 암호화기입니다.
type AESEncryptor struct {
	key []byte // 32 bytes (AES-256)
}

// NewAESEncryptor는 hex 인코딩된 키(64자)로 AESEncryptor를 생성합니다.
// 키가 비어 있으면 nil을 반환합니다 (암호화 비활성화).
func NewAESEncryptor(keyHex string) (*AESEncryptor, error) {
	if keyHex == "" {
		return nil, nil
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("암호화 키 hex 디코딩 실패: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("암호화 키는 32 bytes(64 hex chars)여야 합니다, 현재: %d bytes", len(key))
	}

	return &AESEncryptor{key: key}, nil
}

// Encrypt는 평문을 AES-256-GCM으로 암호화하여 base64 문자열을 반환합니다.
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	if e == nil {
		return plaintext, nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("AES 블록 생성 실패: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM 생성 실패: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce 생성 실패: %w", err)
	}

	// nonce + ciphertext를 합쳐서 base64 인코딩
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt는 base64 인코딩된 암호문을 복호화하여 평문을 반환합니다.
func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
	if e == nil {
		return ciphertext, nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 디코딩 실패: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("AES 블록 생성 실패: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM 생성 실패: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("암호문이 너무 짧습니다")
	}

	nonce, ct := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("복호화 실패: %w", err)
	}

	return string(plaintext), nil
}
