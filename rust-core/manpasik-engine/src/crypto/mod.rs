//! 암호화 모듈
//!
//! AES-256-GCM 암호화/복호화, SHA-256 해시, HMAC, 해시체인, 키 유도(HKDF)
//!
//! # 보안 사양
//! - 대칭 암호화: AES-256-GCM (96비트 Nonce, 128비트 태그)
//! - 해시: SHA-256
//! - 키 유도: HKDF-SHA256
//! - 해시체인: SHA-256 기반 무결성 검증 (의료기기 IEC 62304 요구사항)

use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    AeadCore, Aes256Gcm, Nonce,
};
use ring::hmac;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use thiserror::Error;

/// 암호화 관련 에러
#[derive(Debug, Error)]
pub enum CryptoError {
    #[error("암호화 실패: {0}")]
    EncryptionFailed(String),

    #[error("복호화 실패: {0}")]
    DecryptionFailed(String),

    #[error("유효하지 않은 키 길이: expected {expected}, got {got}")]
    InvalidKeyLength { expected: usize, got: usize },

    #[error("유효하지 않은 Nonce 길이: expected {expected}, got {got}")]
    InvalidNonceLength { expected: usize, got: usize },

    #[error("해시체인 검증 실패: index {index}")]
    HashChainVerificationFailed { index: usize },

    #[error("HMAC 검증 실패")]
    HmacVerificationFailed,
}

/// AES-256-GCM 키 크기 (32바이트)
pub const AES_256_KEY_SIZE: usize = 32;

/// AES-256-GCM Nonce 크기 (12바이트)
pub const AES_GCM_NONCE_SIZE: usize = 12;

/// SHA-256 해시 크기 (32바이트)
pub const SHA256_HASH_SIZE: usize = 32;

/// 암호화된 데이터 구조
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EncryptedData {
    /// Nonce (12바이트, Base64)
    pub nonce: Vec<u8>,
    /// 암호문 + 인증 태그
    pub ciphertext: Vec<u8>,
}

/// 해시체인 항목 (의료 데이터 무결성 검증)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HashChainEntry {
    /// 단계 이름
    pub step: String,
    /// 데이터 해시
    pub data_hash: String,
    /// 이전 해시와 연결된 체인 해시
    pub chain_hash: String,
    /// 타임스탬프
    pub timestamp: i64,
}

/// 암호화 엔진
///
/// 의료 데이터 보호를 위한 암호화 기능 제공
/// - AES-256-GCM: 측정 데이터 암호화 (PHI 보호)
/// - SHA-256 해시체인: 데이터 무결성 검증 (IEC 62304)
/// - HMAC-SHA256: 메시지 인증
pub struct CryptoEngine {
    /// 해시체인 상태 (마지막 체인 해시)
    chain_state: Option<String>,
    /// 해시체인 엔트리 목록
    chain_entries: Vec<HashChainEntry>,
}

impl CryptoEngine {
    /// 새 암호화 엔진 생성
    pub fn new() -> Self {
        Self {
            chain_state: None,
            chain_entries: Vec::new(),
        }
    }

    // =========================================================================
    // AES-256-GCM 암호화/복호화
    // =========================================================================

    /// AES-256-GCM 암호화
    ///
    /// # Arguments
    /// * `plaintext` - 평문 데이터
    /// * `key` - 32바이트 AES-256 키
    ///
    /// # Returns
    /// 암호화된 데이터 (Nonce + 암호문 + 인증 태그)
    ///
    /// # Security
    /// - 랜덤 Nonce 자동 생성 (96비트)
    /// - 인증 태그 128비트 (GCM)
    pub fn encrypt_aes256(
        &self,
        plaintext: &[u8],
        key: &[u8],
    ) -> Result<EncryptedData, CryptoError> {
        // 키 길이 검증
        if key.len() != AES_256_KEY_SIZE {
            return Err(CryptoError::InvalidKeyLength {
                expected: AES_256_KEY_SIZE,
                got: key.len(),
            });
        }

        let cipher = Aes256Gcm::new_from_slice(key)
            .map_err(|e| CryptoError::EncryptionFailed(e.to_string()))?;

        // 랜덤 Nonce 생성 (96비트)
        let nonce = Aes256Gcm::generate_nonce(&mut OsRng);

        // 암호화 (인증 태그 자동 포함)
        let ciphertext = cipher
            .encrypt(&nonce, plaintext)
            .map_err(|e| CryptoError::EncryptionFailed(e.to_string()))?;

        Ok(EncryptedData {
            nonce: nonce.to_vec(),
            ciphertext,
        })
    }

    /// AES-256-GCM 복호화
    ///
    /// # Arguments
    /// * `encrypted` - 암호화된 데이터
    /// * `key` - 32바이트 AES-256 키
    ///
    /// # Returns
    /// 복호화된 평문
    pub fn decrypt_aes256(
        &self,
        encrypted: &EncryptedData,
        key: &[u8],
    ) -> Result<Vec<u8>, CryptoError> {
        // 키 길이 검증
        if key.len() != AES_256_KEY_SIZE {
            return Err(CryptoError::InvalidKeyLength {
                expected: AES_256_KEY_SIZE,
                got: key.len(),
            });
        }

        // Nonce 길이 검증
        if encrypted.nonce.len() != AES_GCM_NONCE_SIZE {
            return Err(CryptoError::InvalidNonceLength {
                expected: AES_GCM_NONCE_SIZE,
                got: encrypted.nonce.len(),
            });
        }

        let cipher = Aes256Gcm::new_from_slice(key)
            .map_err(|e| CryptoError::DecryptionFailed(e.to_string()))?;

        let nonce = Nonce::from_slice(&encrypted.nonce);

        cipher
            .decrypt(nonce, encrypted.ciphertext.as_ref())
            .map_err(|e| CryptoError::DecryptionFailed(e.to_string()))
    }

    /// AES-256 키 생성 (랜덤)
    pub fn generate_key() -> Vec<u8> {
        use aes_gcm::aead::rand_core::RngCore;
        let mut key = vec![0u8; AES_256_KEY_SIZE];
        OsRng.fill_bytes(&mut key);
        key
    }

    // =========================================================================
    // SHA-256 해싱
    // =========================================================================

    /// SHA-256 해시 (hex 문자열)
    pub fn hash_sha256(&self, data: &[u8]) -> String {
        let mut hasher = Sha256::new();
        hasher.update(data);
        format!("{:x}", hasher.finalize())
    }

    /// SHA-256 해시 (바이트 배열)
    pub fn hash_sha256_bytes(&self, data: &[u8]) -> Vec<u8> {
        let mut hasher = Sha256::new();
        hasher.update(data);
        hasher.finalize().to_vec()
    }

    // =========================================================================
    // HMAC-SHA256
    // =========================================================================

    /// HMAC-SHA256 서명 생성
    ///
    /// # Arguments
    /// * `data` - 서명할 데이터
    /// * `key` - HMAC 키
    pub fn hmac_sign(&self, data: &[u8], key: &[u8]) -> Vec<u8> {
        let signing_key = hmac::Key::new(hmac::HMAC_SHA256, key);
        let tag = hmac::sign(&signing_key, data);
        tag.as_ref().to_vec()
    }

    /// HMAC-SHA256 서명 검증
    ///
    /// # Arguments
    /// * `data` - 검증할 데이터
    /// * `key` - HMAC 키
    /// * `signature` - 검증할 서명
    pub fn hmac_verify(
        &self,
        data: &[u8],
        key: &[u8],
        signature: &[u8],
    ) -> Result<(), CryptoError> {
        let verification_key = hmac::Key::new(hmac::HMAC_SHA256, key);
        hmac::verify(&verification_key, data, signature)
            .map_err(|_| CryptoError::HmacVerificationFailed)
    }

    // =========================================================================
    // HKDF 키 유도
    // =========================================================================

    /// HKDF-SHA256 키 유도
    ///
    /// 하나의 마스터 키로부터 여러 파생 키를 생성
    ///
    /// # Arguments
    /// * `ikm` - 입력 키 재료 (Input Keying Material)
    /// * `salt` - 솔트 (선택적, None이면 빈 솔트)
    /// * `info` - 컨텍스트 정보
    /// * `output_len` - 출력 키 길이
    pub fn hkdf_derive(
        &self,
        ikm: &[u8],
        salt: Option<&[u8]>,
        info: &[u8],
        output_len: usize,
    ) -> Result<Vec<u8>, CryptoError> {
        let salt_value = salt.unwrap_or(&[]);
        let s = ring::hkdf::Salt::new(ring::hkdf::HKDF_SHA256, salt_value);
        let prk = s.extract(ikm);

        let mut output = vec![0u8; output_len];
        prk.expand(&[info], HkdfLen(output_len))
            .map_err(|_| CryptoError::EncryptionFailed("HKDF 확장 실패".to_string()))?
            .fill(&mut output)
            .map_err(|_| CryptoError::EncryptionFailed("HKDF 출력 생성 실패".to_string()))?;

        Ok(output)
    }

    // =========================================================================
    // 해시체인 (의료 데이터 무결성 - IEC 62304)
    // =========================================================================

    /// 해시체인에 새 엔트리 추가
    ///
    /// 측정 데이터의 변환 과정을 해시체인으로 기록하여
    /// 데이터 무결성을 보장합니다 (IEC 62304, ISO 14971 요구사항).
    ///
    /// # Arguments
    /// * `step` - 변환 단계 이름 (예: "raw_capture", "differential_correction")
    /// * `data` - 해당 단계의 데이터
    pub fn add_chain_entry(&mut self, step: &str, data: &[u8]) -> HashChainEntry {
        let data_hash = self.hash_sha256(data);

        // 체인 해시 = SHA256(이전_체인_해시 + 현재_데이터_해시)
        let chain_input = match &self.chain_state {
            Some(prev) => format!("{}{}", prev, data_hash),
            None => data_hash.clone(),
        };
        let chain_hash = self.hash_sha256(chain_input.as_bytes());

        let entry = HashChainEntry {
            step: step.to_string(),
            data_hash,
            chain_hash: chain_hash.clone(),
            timestamp: chrono::Utc::now().timestamp_millis(),
        };

        self.chain_state = Some(chain_hash);
        self.chain_entries.push(entry.clone());

        entry
    }

    /// 해시체인 검증
    ///
    /// 체인의 모든 엔트리가 올바른 순서로 연결되어 있는지 검증
    pub fn verify_chain(&self) -> Result<bool, CryptoError> {
        if self.chain_entries.is_empty() {
            return Ok(true);
        }

        let mut prev_chain_hash: Option<&str> = None;

        for (i, entry) in self.chain_entries.iter().enumerate() {
            let chain_input = match prev_chain_hash {
                Some(prev) => format!("{}{}", prev, entry.data_hash),
                None => entry.data_hash.clone(),
            };
            let expected_hash = self.hash_sha256(chain_input.as_bytes());

            if expected_hash != entry.chain_hash {
                return Err(CryptoError::HashChainVerificationFailed { index: i });
            }

            prev_chain_hash = Some(&entry.chain_hash);
        }

        Ok(true)
    }

    /// 해시체인 엔트리 목록 조회
    pub fn chain_entries(&self) -> &[HashChainEntry] {
        &self.chain_entries
    }

    /// 해시체인 초기화
    pub fn reset_chain(&mut self) {
        self.chain_state = None;
        self.chain_entries.clear();
    }

    /// 현재 체인 해시 조회
    pub fn current_chain_hash(&self) -> Option<&str> {
        self.chain_state.as_deref()
    }
}

impl Default for CryptoEngine {
    fn default() -> Self {
        Self::new()
    }
}

/// HKDF 출력 길이를 위한 내부 타입
struct HkdfLen(usize);

impl ring::hkdf::KeyType for HkdfLen {
    fn len(&self) -> usize {
        self.0
    }
}

// =============================================================================
// 테스트
// =============================================================================

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_aes256_gcm_암호화_복호화() {
        let engine = CryptoEngine::new();
        let key = CryptoEngine::generate_key();
        let plaintext = b"ManPaSik measurement data: glucose=120mg/dL";

        // 암호화
        let encrypted = engine.encrypt_aes256(plaintext, &key).unwrap();
        assert_ne!(encrypted.ciphertext, plaintext.to_vec());
        assert_eq!(encrypted.nonce.len(), AES_GCM_NONCE_SIZE);

        // 복호화
        let decrypted = engine.decrypt_aes256(&encrypted, &key).unwrap();
        assert_eq!(decrypted, plaintext.to_vec());
    }

    #[test]
    fn test_aes256_gcm_잘못된_키_복호화_실패() {
        let engine = CryptoEngine::new();
        let key = CryptoEngine::generate_key();
        let wrong_key = CryptoEngine::generate_key();
        let plaintext = b"sensitive medical data";

        let encrypted = engine.encrypt_aes256(plaintext, &key).unwrap();

        // 잘못된 키로 복호화 시도 → 실패
        let result = engine.decrypt_aes256(&encrypted, &wrong_key);
        assert!(result.is_err());
    }

    #[test]
    fn test_aes256_gcm_유효하지_않은_키_길이() {
        let engine = CryptoEngine::new();
        let short_key = vec![0u8; 16]; // 16바이트 (32바이트 필요)

        let result = engine.encrypt_aes256(b"test", &short_key);
        assert!(matches!(result, Err(CryptoError::InvalidKeyLength { .. })));
    }

    #[test]
    fn test_aes256_gcm_변조_감지() {
        let engine = CryptoEngine::new();
        let key = CryptoEngine::generate_key();
        let plaintext = b"critical measurement data";

        let mut encrypted = engine.encrypt_aes256(plaintext, &key).unwrap();

        // 암호문 변조
        if let Some(byte) = encrypted.ciphertext.first_mut() {
            *byte ^= 0xFF;
        }

        // 변조된 데이터 복호화 시도 → GCM 인증 태그 검증 실패
        let result = engine.decrypt_aes256(&encrypted, &key);
        assert!(result.is_err());
    }

    #[test]
    fn test_sha256_해시() {
        let engine = CryptoEngine::new();

        let hash = engine.hash_sha256(b"hello");
        assert_eq!(
            hash,
            "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
        );

        // 동일 입력 → 동일 해시
        let hash2 = engine.hash_sha256(b"hello");
        assert_eq!(hash, hash2);

        // 다른 입력 → 다른 해시
        let hash3 = engine.hash_sha256(b"world");
        assert_ne!(hash, hash3);
    }

    #[test]
    fn test_hmac_서명_검증() {
        let engine = CryptoEngine::new();
        let key = b"secret-hmac-key-for-manpasik";
        let data = b"measurement session data";

        // 서명 생성
        let signature = engine.hmac_sign(data, key);
        assert!(!signature.is_empty());

        // 서명 검증 성공
        let result = engine.hmac_verify(data, key, &signature);
        assert!(result.is_ok());

        // 데이터 변조 시 검증 실패
        let tampered_data = b"tampered measurement data";
        let result = engine.hmac_verify(tampered_data, key, &signature);
        assert!(result.is_err());
    }

    #[test]
    fn test_hkdf_키_유도() {
        let engine = CryptoEngine::new();
        let master_key = b"master-key-material";
        let salt = b"manpasik-salt";

        // 서로 다른 info로 다른 키 파생
        let key1 = engine
            .hkdf_derive(master_key, Some(salt), b"encryption-key", 32)
            .unwrap();
        let key2 = engine
            .hkdf_derive(master_key, Some(salt), b"hmac-key", 32)
            .unwrap();

        assert_eq!(key1.len(), 32);
        assert_eq!(key2.len(), 32);
        assert_ne!(key1, key2); // 다른 info → 다른 키

        // 동일 파라미터 → 동일 키 (결정적)
        let key1_again = engine
            .hkdf_derive(master_key, Some(salt), b"encryption-key", 32)
            .unwrap();
        assert_eq!(key1, key1_again);
    }

    #[test]
    fn test_해시체인_생성_및_검증() {
        let mut engine = CryptoEngine::new();

        // 측정 데이터 변환 파이프라인 시뮬레이션
        let raw_data = b"raw sensor readings: [1.0, 2.0, 3.0, 4.0]";
        let corrected_data = b"differential corrected: [0.905, 1.810, 2.715, 3.620]";
        let fingerprint_data = b"fingerprint vector: [0.12, 0.34, ...]";

        // 각 단계를 해시체인에 기록
        let entry1 = engine.add_chain_entry("raw_capture", raw_data);
        let entry2 = engine.add_chain_entry("differential_correction", corrected_data);
        let entry3 = engine.add_chain_entry("fingerprint_generation", fingerprint_data);

        // 체인 엔트리 수 확인
        assert_eq!(engine.chain_entries().len(), 3);

        // 각 엔트리의 체인 해시가 서로 다른지 확인
        assert_ne!(entry1.chain_hash, entry2.chain_hash);
        assert_ne!(entry2.chain_hash, entry3.chain_hash);

        // 체인 검증 성공
        assert!(engine.verify_chain().is_ok());
    }

    #[test]
    fn test_해시체인_빈_체인_검증() {
        let engine = CryptoEngine::new();
        assert!(engine.verify_chain().unwrap());
    }

    #[test]
    fn test_해시체인_초기화() {
        let mut engine = CryptoEngine::new();

        engine.add_chain_entry("step1", b"data1");
        engine.add_chain_entry("step2", b"data2");
        assert_eq!(engine.chain_entries().len(), 2);
        assert!(engine.current_chain_hash().is_some());

        engine.reset_chain();
        assert_eq!(engine.chain_entries().len(), 0);
        assert!(engine.current_chain_hash().is_none());
    }

    #[test]
    fn test_키_생성_유일성() {
        let key1 = CryptoEngine::generate_key();
        let key2 = CryptoEngine::generate_key();

        assert_eq!(key1.len(), AES_256_KEY_SIZE);
        assert_eq!(key2.len(), AES_256_KEY_SIZE);
        assert_ne!(key1, key2); // 매번 다른 키
    }

    #[test]
    fn test_전체_측정_데이터_암호화_파이프라인() {
        // 엔드투엔드 시나리오: 측정 → 해시체인 → 암호화 → 복호화 → 해시체인 검증
        let mut engine = CryptoEngine::new();
        let master_key = b"device-master-key-material-1234";

        // 1. 마스터 키에서 암호화 키와 HMAC 키 유도
        let enc_key = engine
            .hkdf_derive(master_key, Some(b"manpasik"), b"aes-key", 32)
            .unwrap();
        let mac_key = engine
            .hkdf_derive(master_key, Some(b"manpasik"), b"hmac-key", 32)
            .unwrap();

        // 2. 측정 데이터 생성 (시뮬레이션)
        let measurement = br#"{"glucose": 120, "unit": "mg/dL", "confidence": 0.95}"#;

        // 3. 해시체인에 기록
        engine.add_chain_entry("measurement_raw", measurement);

        // 4. 데이터 암호화 (PHI 보호)
        let encrypted = engine.encrypt_aes256(measurement, &enc_key).unwrap();

        // 5. HMAC 서명
        let signature = engine.hmac_sign(&encrypted.ciphertext, &mac_key);

        // 6. 전송 후 수신 측에서 검증
        engine
            .hmac_verify(&encrypted.ciphertext, &mac_key, &signature)
            .unwrap();

        // 7. 복호화
        let decrypted = engine.decrypt_aes256(&encrypted, &enc_key).unwrap();
        assert_eq!(decrypted, measurement.to_vec());

        // 8. 해시체인 무결성 검증
        assert!(engine.verify_chain().is_ok());
    }
}
