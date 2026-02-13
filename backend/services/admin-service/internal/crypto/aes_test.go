package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
)

func generateTestKey() string {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	return hex.EncodeToString(key)
}

func TestNewAESEncryptor_ValidKey(t *testing.T) {
	enc, err := NewAESEncryptor(generateTestKey())
	if err != nil {
		t.Fatalf("유효한 키로 생성 실패: %v", err)
	}
	if enc == nil {
		t.Fatal("enc가 nil입니다")
	}
}

func TestNewAESEncryptor_EmptyKey(t *testing.T) {
	enc, err := NewAESEncryptor("")
	if err != nil {
		t.Fatalf("빈 키에서 에러 발생: %v", err)
	}
	if enc != nil {
		t.Fatal("빈 키일 때 nil이어야 합니다")
	}
}

func TestNewAESEncryptor_InvalidHex(t *testing.T) {
	_, err := NewAESEncryptor("not-hex")
	if err == nil {
		t.Fatal("잘못된 hex에서 에러가 발생해야 합니다")
	}
}

func TestNewAESEncryptor_WrongKeyLength(t *testing.T) {
	_, err := NewAESEncryptor("abcdef0123456789") // 8 bytes
	if err == nil {
		t.Fatal("잘못된 키 길이에서 에러가 발생해야 합니다")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	keyHex := generateTestKey()
	enc, err := NewAESEncryptor(keyHex)
	if err != nil {
		t.Fatalf("생성 실패: %v", err)
	}

	testCases := []string{
		"hello world",
		"test_sk_01234567890",
		"",
		"한국어 테스트 문자열",
		"🎉 이모지 포함",
		"very-long-secret-key-that-is-used-for-testing-purposes-only-and-should-not-be-used-in-production",
	}

	for _, tc := range testCases {
		encrypted, err := enc.Encrypt(tc)
		if err != nil {
			t.Fatalf("암호화 실패 (%q): %v", tc, err)
		}

		// 암호화된 값은 원본과 달라야 함 (빈 문자열 제외)
		if tc != "" && encrypted == tc {
			t.Errorf("암호화 후에도 값이 같습니다: %q", tc)
		}

		decrypted, err := enc.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("복호화 실패 (%q): %v", tc, err)
		}

		if decrypted != tc {
			t.Errorf("왕복 실패: 원본=%q, 복호화=%q", tc, decrypted)
		}
	}
}

func TestEncryptDecrypt_DifferentCiphertexts(t *testing.T) {
	enc, _ := NewAESEncryptor(generateTestKey())
	plain := "same-plaintext"

	c1, _ := enc.Encrypt(plain)
	c2, _ := enc.Encrypt(plain)

	if c1 == c2 {
		t.Error("같은 평문이라도 nonce가 달라 암호문이 달라야 합니다")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	enc, _ := NewAESEncryptor(generateTestKey())
	_, err := enc.Decrypt("not-valid-base64!!!")
	if err == nil {
		t.Error("잘못된 base64에서 에러가 발생해야 합니다")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc, _ := NewAESEncryptor(generateTestKey())
	encrypted, _ := enc.Encrypt("secret-value")

	// 암호문 변조
	tampered := encrypted[:len(encrypted)-2] + "XX"
	_, err := enc.Decrypt(tampered)
	if err == nil {
		t.Error("변조된 암호문에서 에러가 발생해야 합니다")
	}
}

func TestNilEncryptor_Passthrough(t *testing.T) {
	var enc *AESEncryptor

	result, err := enc.Encrypt("plaintext")
	if err != nil {
		t.Fatalf("nil 암호화기 Encrypt 에러: %v", err)
	}
	if result != "plaintext" {
		t.Errorf("nil 암호화기는 원본을 그대로 반환해야 합니다: %q", result)
	}

	result, err = enc.Decrypt("plaintext")
	if err != nil {
		t.Fatalf("nil 암호화기 Decrypt 에러: %v", err)
	}
	if result != "plaintext" {
		t.Errorf("nil 암호화기는 원본을 그대로 반환해야 합니다: %q", result)
	}
}
