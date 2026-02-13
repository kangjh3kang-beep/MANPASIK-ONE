//! 오프라인 동기화 모듈
//!
//! CRDT (Conflict-free Replicated Data Type) 기반 충돌 없는 데이터 동기화
//!
//! # 지원 CRDT 타입
//! - **GCounter**: 증가만 가능한 카운터 (측정 횟수 등)
//! - **LWWRegister**: Last-Writer-Wins 레지스터 (설정값, 프로필 등)
//! - **ORSet**: Observed-Remove 셋 (디바이스 목록, 카트리지 목록 등)
//!
//! # 오프라인 동기화 전략
//! 1. 모든 측정 데이터를 로컬에 먼저 저장 (오프라인 100% 동작)
//! 2. 네트워크 연결 시 CRDT 병합으로 충돌 없는 동기화
//! 3. 동기화 큐로 미전송 데이터 관리

use serde::{Deserialize, Serialize};
use std::collections::{BTreeMap, HashMap, HashSet};
use thiserror::Error;
use uuid::Uuid;

/// 동기화 관련 에러
#[derive(Debug, Error)]
pub enum SyncError {
    #[error("직렬화 실패: {0}")]
    SerializationFailed(String),

    #[error("역직렬화 실패: {0}")]
    DeserializationFailed(String),

    #[error("동기화 큐가 가득 참 (최대: {max})")]
    QueueFull { max: usize },

    #[error("네트워크 연결 없음")]
    NoConnection,

    #[error("동기화 충돌: {0}")]
    ConflictDetected(String),
}

/// 동기화 상태
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum SyncState {
    /// 동기화 완료
    Synced,
    /// 동기화 대기 중
    Pending,
    /// 동기화 중
    Syncing,
    /// 충돌 발생
    Conflicted,
    /// 에러
    Error,
}

/// 동기화 큐 항목
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SyncQueueItem {
    /// 고유 ID
    pub id: String,
    /// 작업 유형
    pub operation: SyncOperation,
    /// 직렬화된 데이터
    pub data: Vec<u8>,
    /// 생성 타임스탬프 (밀리초)
    pub timestamp: i64,
    /// 동기화 상태
    pub state: SyncState,
    /// 재시도 횟수
    pub retry_count: u32,
    /// 노드 ID (이 데이터를 생성한 디바이스)
    pub node_id: String,
}

/// 동기화 작업 유형
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub enum SyncOperation {
    /// 측정 데이터 업로드
    MeasurementUpload,
    /// 디바이스 상태 업데이트
    DeviceStatusUpdate,
    /// 사용자 설정 동기화
    SettingsSync,
    /// 카트리지 사용 기록
    CartridgeUsageLog,
    /// CRDT 상태 병합
    CrdtMerge,
}

// =============================================================================
// CRDT: G-Counter (Grow-only Counter)
// =============================================================================

/// G-Counter: 증가만 가능한 분산 카운터
///
/// 각 노드가 독립적으로 증가시킬 수 있으며,
/// 병합 시 각 노드의 최댓값을 취합니다.
///
/// 용도: 측정 횟수, 카트리지 사용 횟수 등
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GCounter {
    /// 노드별 카운터 값
    counts: BTreeMap<String, u64>,
}

impl GCounter {
    /// 새 G-Counter 생성
    pub fn new() -> Self {
        Self {
            counts: BTreeMap::new(),
        }
    }

    /// 특정 노드의 카운터 증가
    pub fn increment(&mut self, node_id: &str) {
        let count = self.counts.entry(node_id.to_string()).or_insert(0);
        *count += 1;
    }

    /// 특정 노드의 카운터를 지정값만큼 증가
    pub fn increment_by(&mut self, node_id: &str, amount: u64) {
        let count = self.counts.entry(node_id.to_string()).or_insert(0);
        *count += amount;
    }

    /// 전체 카운터 값 (모든 노드 합산)
    pub fn value(&self) -> u64 {
        self.counts.values().sum()
    }

    /// 특정 노드의 카운터 값
    pub fn node_value(&self, node_id: &str) -> u64 {
        self.counts.get(node_id).copied().unwrap_or(0)
    }

    /// 다른 G-Counter와 병합
    ///
    /// 각 노드에 대해 최댓값을 취합니다 (교환법칙, 결합법칙, 멱등성 보장)
    pub fn merge(&mut self, other: &GCounter) {
        for (node_id, &count) in &other.counts {
            let current = self.counts.entry(node_id.clone()).or_insert(0);
            *current = (*current).max(count);
        }
    }
}

impl Default for GCounter {
    fn default() -> Self {
        Self::new()
    }
}

// =============================================================================
// CRDT: LWW-Register (Last-Writer-Wins Register)
// =============================================================================

/// LWW-Register: Last-Writer-Wins 레지스터
///
/// 타임스탬프가 가장 최신인 값이 승리합니다.
///
/// 용도: 사용자 설정, 프로필 정보, 디바이스 이름 등
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LWWRegister<T: Clone + Serialize> {
    /// 현재 값
    value: T,
    /// 마지막 업데이트 타임스탬프 (밀리초)
    timestamp: i64,
    /// 업데이트한 노드 ID
    node_id: String,
}

impl<T: Clone + Serialize> LWWRegister<T> {
    /// 새 LWW-Register 생성
    pub fn new(value: T, node_id: &str) -> Self {
        Self {
            value,
            timestamp: chrono::Utc::now().timestamp_millis(),
            node_id: node_id.to_string(),
        }
    }

    /// 특정 타임스탬프로 LWW-Register 생성
    pub fn with_timestamp(value: T, node_id: &str, timestamp: i64) -> Self {
        Self {
            value,
            timestamp,
            node_id: node_id.to_string(),
        }
    }

    /// 현재 값 조회
    pub fn value(&self) -> &T {
        &self.value
    }

    /// 타임스탬프 조회
    pub fn timestamp(&self) -> i64 {
        self.timestamp
    }

    /// 노드 ID 조회
    pub fn node_id(&self) -> &str {
        &self.node_id
    }

    /// 값 업데이트 (현재 시간 사용)
    pub fn set(&mut self, value: T, node_id: &str) {
        self.value = value;
        self.timestamp = chrono::Utc::now().timestamp_millis();
        self.node_id = node_id.to_string();
    }

    /// 다른 LWW-Register와 병합 (Last-Writer-Wins)
    ///
    /// 타임스탬프가 더 최신인 쪽이 승리.
    /// 동일 타임스탬프면 노드 ID로 결정 (결정적).
    pub fn merge(&mut self, other: &LWWRegister<T>) {
        if other.timestamp > self.timestamp
            || (other.timestamp == self.timestamp && other.node_id > self.node_id)
        {
            self.value = other.value.clone();
            self.timestamp = other.timestamp;
            self.node_id = other.node_id.clone();
        }
    }
}

// =============================================================================
// CRDT: OR-Set (Observed-Remove Set)
// =============================================================================

/// OR-Set 요소 (고유 태그 포함)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
struct TaggedElement {
    /// 요소 값
    value: String,
    /// 고유 태그 (UUID)
    tag: String,
}

/// OR-Set: Observed-Remove Set
///
/// 요소를 추가/제거할 수 있으며, 동시 추가+제거 시
/// "추가가 제거를 이긴다" (add-wins) 의미론을 따릅니다.
///
/// 용도: 디바이스 목록, 즐겨찾기, 카트리지 목록 등
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ORSet {
    /// 활성 요소 (값 → 태그 집합)
    elements: HashMap<String, HashSet<String>>,
    /// 삭제된 태그 집합 (tombstone)
    tombstones: HashSet<String>,
}

impl ORSet {
    /// 새 OR-Set 생성
    pub fn new() -> Self {
        Self {
            elements: HashMap::new(),
            tombstones: HashSet::new(),
        }
    }

    /// 요소 추가
    ///
    /// 고유 태그를 생성하여 추가합니다.
    pub fn add(&mut self, value: &str) -> String {
        let tag = Uuid::new_v4().to_string();
        self.elements
            .entry(value.to_string())
            .or_default()
            .insert(tag.clone());
        tag
    }

    /// 요소 제거
    ///
    /// 해당 값의 현재 모든 태그를 tombstone에 추가합니다.
    pub fn remove(&mut self, value: &str) {
        if let Some(tags) = self.elements.remove(value) {
            for tag in tags {
                self.tombstones.insert(tag);
            }
        }
    }

    /// 요소 존재 여부 확인
    pub fn contains(&self, value: &str) -> bool {
        self.elements
            .get(value)
            .is_some_and(|tags| !tags.is_empty())
    }

    /// 모든 요소 조회
    pub fn elements(&self) -> Vec<&str> {
        self.elements
            .iter()
            .filter(|(_, tags)| !tags.is_empty())
            .map(|(k, _)| k.as_str())
            .collect()
    }

    /// 요소 수
    pub fn len(&self) -> usize {
        self.elements
            .values()
            .filter(|tags| !tags.is_empty())
            .count()
    }

    /// 비어있는지 확인
    pub fn is_empty(&self) -> bool {
        self.len() == 0
    }

    /// 다른 OR-Set과 병합
    ///
    /// 1. 두 셋의 요소를 합집합
    /// 2. tombstone에 있는 태그는 제거
    pub fn merge(&mut self, other: &ORSet) {
        // tombstone 병합
        self.tombstones.extend(other.tombstones.iter().cloned());

        // 요소 병합
        for (value, tags) in &other.elements {
            let entry = self.elements.entry(value.clone()).or_default();
            entry.extend(tags.iter().cloned());
        }

        // tombstone에 해당하는 태그 제거
        for (_value, tags) in self.elements.iter_mut() {
            tags.retain(|tag| !self.tombstones.contains(tag));
        }

        // 빈 엔트리 정리
        self.elements.retain(|_, tags| !tags.is_empty());
    }
}

impl Default for ORSet {
    fn default() -> Self {
        Self::new()
    }
}

// =============================================================================
// 동기화 매니저
// =============================================================================

/// 동기화 큐 최대 크기
const MAX_QUEUE_SIZE: usize = 10_000;

/// 최대 재시도 횟수
const MAX_RETRIES: u32 = 5;

/// 동기화 매니저
///
/// 오프라인 데이터 큐 관리 및 CRDT 기반 동기화를 담당합니다.
pub struct SyncManager {
    /// 이 디바이스의 노드 ID
    node_id: String,
    /// 동기화 대기 큐
    queue: Vec<SyncQueueItem>,
    /// 측정 횟수 카운터 (G-Counter)
    measurement_counter: GCounter,
    /// 디바이스 목록 (OR-Set)
    device_set: ORSet,
    /// 연결 상태
    is_connected: bool,
}

impl SyncManager {
    /// 새 동기화 매니저 생성
    pub fn new(node_id: &str) -> Self {
        Self {
            node_id: node_id.to_string(),
            queue: Vec::new(),
            measurement_counter: GCounter::new(),
            device_set: ORSet::new(),
            is_connected: false,
        }
    }

    /// 기본 노드 ID로 생성 (UUID)
    pub fn with_random_id() -> Self {
        Self::new(&Uuid::new_v4().to_string())
    }

    /// 노드 ID 조회
    pub fn node_id(&self) -> &str {
        &self.node_id
    }

    // =========================================================================
    // 동기화 큐 관리
    // =========================================================================

    /// 동기화 큐에 항목 추가
    pub fn enqueue(
        &mut self,
        operation: SyncOperation,
        data: Vec<u8>,
    ) -> Result<String, SyncError> {
        if self.queue.len() >= MAX_QUEUE_SIZE {
            return Err(SyncError::QueueFull {
                max: MAX_QUEUE_SIZE,
            });
        }

        let id = Uuid::new_v4().to_string();
        let item = SyncQueueItem {
            id: id.clone(),
            operation,
            data,
            timestamp: chrono::Utc::now().timestamp_millis(),
            state: SyncState::Pending,
            retry_count: 0,
            node_id: self.node_id.clone(),
        };

        self.queue.push(item);
        Ok(id)
    }

    /// 대기 중인 항목 수
    pub fn pending_count(&self) -> usize {
        self.queue
            .iter()
            .filter(|i| i.state == SyncState::Pending)
            .count()
    }

    /// 전체 큐 크기
    pub fn queue_size(&self) -> usize {
        self.queue.len()
    }

    /// 동기화 실행
    ///
    /// 대기 중인 모든 항목을 처리합니다.
    /// 실제 네트워크 전송은 외부에서 처리하며,
    /// 이 메서드는 큐 상태 관리를 담당합니다.
    pub async fn sync(&mut self) -> Result<SyncResult, SyncError> {
        if !self.is_connected {
            return Err(SyncError::NoConnection);
        }

        let mut synced = 0;
        let mut failed = 0;

        for item in self.queue.iter_mut() {
            if item.state == SyncState::Pending || item.state == SyncState::Error {
                if item.retry_count >= MAX_RETRIES {
                    item.state = SyncState::Error;
                    failed += 1;
                    continue;
                }

                item.state = SyncState::Syncing;
                // 실제 전송은 외부 콜백/트레이트로 처리
                // 여기서는 성공으로 시뮬레이션
                item.state = SyncState::Synced;
                synced += 1;
            }
        }

        // 동기화 완료된 항목 제거
        self.queue.retain(|item| item.state != SyncState::Synced);

        Ok(SyncResult {
            synced_count: synced,
            failed_count: failed,
            remaining_count: self.queue.len(),
        })
    }

    /// 연결 상태 설정
    pub fn set_connected(&mut self, connected: bool) {
        self.is_connected = connected;
    }

    /// 연결 상태 조회
    pub fn is_connected(&self) -> bool {
        self.is_connected
    }

    // =========================================================================
    // CRDT 접근
    // =========================================================================

    /// 측정 횟수 증가
    pub fn increment_measurement_count(&mut self) {
        self.measurement_counter.increment(&self.node_id);
    }

    /// 전체 측정 횟수
    pub fn total_measurement_count(&self) -> u64 {
        self.measurement_counter.value()
    }

    /// 측정 카운터 병합
    pub fn merge_measurement_counter(&mut self, other: &GCounter) {
        self.measurement_counter.merge(other);
    }

    /// 측정 카운터 참조
    pub fn measurement_counter(&self) -> &GCounter {
        &self.measurement_counter
    }

    /// 디바이스 추가
    pub fn add_device(&mut self, device_id: &str) -> String {
        self.device_set.add(device_id)
    }

    /// 디바이스 제거
    pub fn remove_device(&mut self, device_id: &str) {
        self.device_set.remove(device_id);
    }

    /// 디바이스 존재 확인
    pub fn has_device(&self, device_id: &str) -> bool {
        self.device_set.contains(device_id)
    }

    /// 디바이스 목록
    pub fn devices(&self) -> Vec<&str> {
        self.device_set.elements()
    }

    /// 디바이스 셋 병합
    pub fn merge_device_set(&mut self, other: &ORSet) {
        self.device_set.merge(other);
    }

    /// 디바이스 셋 참조
    pub fn device_set(&self) -> &ORSet {
        &self.device_set
    }

    /// 큐의 실패한 항목 재시도 예약
    pub fn retry_failed(&mut self) {
        for item in self.queue.iter_mut() {
            if item.state == SyncState::Error && item.retry_count < MAX_RETRIES {
                item.state = SyncState::Pending;
                item.retry_count += 1;
            }
        }
    }

    /// 큐 전체 비우기 (주의: 미동기화 데이터 손실)
    pub fn clear_queue(&mut self) {
        self.queue.clear();
    }
}

impl Default for SyncManager {
    fn default() -> Self {
        Self::with_random_id()
    }
}

/// 동기화 결과
#[derive(Debug, Clone)]
pub struct SyncResult {
    /// 동기화 성공 수
    pub synced_count: usize,
    /// 동기화 실패 수
    pub failed_count: usize,
    /// 큐에 남은 항목 수
    pub remaining_count: usize,
}

// =============================================================================
// 테스트
// =============================================================================

#[cfg(test)]
mod tests {
    use super::*;

    // =========================================================================
    // G-Counter 테스트
    // =========================================================================

    #[test]
    fn test_gcounter_증가_및_합산() {
        let mut counter = GCounter::new();

        counter.increment("device-A");
        counter.increment("device-A");
        counter.increment("device-B");

        assert_eq!(counter.value(), 3);
        assert_eq!(counter.node_value("device-A"), 2);
        assert_eq!(counter.node_value("device-B"), 1);
        assert_eq!(counter.node_value("device-C"), 0);
    }

    #[test]
    fn test_gcounter_병합() {
        // 디바이스 A에서
        let mut counter_a = GCounter::new();
        counter_a.increment("device-A");
        counter_a.increment("device-A");
        counter_a.increment("device-B");

        // 디바이스 B에서
        let mut counter_b = GCounter::new();
        counter_b.increment("device-A");
        counter_b.increment("device-B");
        counter_b.increment("device-B");
        counter_b.increment("device-B");

        // 병합
        counter_a.merge(&counter_b);

        // A: max(2, 1) = 2, B: max(1, 3) = 3
        assert_eq!(counter_a.node_value("device-A"), 2);
        assert_eq!(counter_a.node_value("device-B"), 3);
        assert_eq!(counter_a.value(), 5);
    }

    #[test]
    fn test_gcounter_병합_멱등성() {
        let mut counter_a = GCounter::new();
        counter_a.increment("node-1");
        counter_a.increment("node-1");

        let counter_b = counter_a.clone();

        // 같은 것 두 번 병합해도 결과 동일 (멱등성)
        counter_a.merge(&counter_b);
        assert_eq!(counter_a.value(), 2);

        counter_a.merge(&counter_b);
        assert_eq!(counter_a.value(), 2);
    }

    #[test]
    fn test_gcounter_병합_교환법칙() {
        let mut counter_a = GCounter::new();
        counter_a.increment("A");
        counter_a.increment("A");

        let mut counter_b = GCounter::new();
        counter_b.increment("B");

        // A.merge(B)
        let mut result_ab = counter_a.clone();
        result_ab.merge(&counter_b);

        // B.merge(A)
        let mut result_ba = counter_b.clone();
        result_ba.merge(&counter_a);

        // 결과 동일 (교환법칙)
        assert_eq!(result_ab.value(), result_ba.value());
    }

    // =========================================================================
    // LWW-Register 테스트
    // =========================================================================

    #[test]
    fn test_lww_register_기본_동작() {
        let reg = LWWRegister::new("초기값".to_string(), "node-1");

        assert_eq!(reg.value(), "초기값");
        assert_eq!(reg.node_id(), "node-1");
    }

    #[test]
    fn test_lww_register_업데이트() {
        let mut reg = LWWRegister::new("old".to_string(), "node-1");

        reg.set("new".to_string(), "node-1");

        assert_eq!(reg.value(), "new");
    }

    #[test]
    fn test_lww_register_병합_최신_승리() {
        let mut reg_a = LWWRegister::with_timestamp("값A".to_string(), "node-A", 100);
        let reg_b = LWWRegister::with_timestamp("값B".to_string(), "node-B", 200);

        reg_a.merge(&reg_b);

        // B가 더 최신이므로 B의 값이 승리
        assert_eq!(reg_a.value(), "값B");
    }

    #[test]
    fn test_lww_register_동일_타임스탬프_노드_id_결정() {
        let mut reg_a = LWWRegister::with_timestamp("값A".to_string(), "node-A", 100);
        let reg_b = LWWRegister::with_timestamp("값B".to_string(), "node-B", 100);

        reg_a.merge(&reg_b);

        // 같은 타임스탬프이면 노드 ID가 더 큰 쪽이 승리 (B > A)
        assert_eq!(reg_a.value(), "값B");
    }

    // =========================================================================
    // OR-Set 테스트
    // =========================================================================

    #[test]
    fn test_or_set_추가_제거() {
        let mut set = ORSet::new();

        set.add("device-1");
        set.add("device-2");

        assert!(set.contains("device-1"));
        assert!(set.contains("device-2"));
        assert_eq!(set.len(), 2);

        set.remove("device-1");
        assert!(!set.contains("device-1"));
        assert!(set.contains("device-2"));
        assert_eq!(set.len(), 1);
    }

    #[test]
    fn test_or_set_병합_동시_추가() {
        let mut set_a = ORSet::new();
        let mut set_b = ORSet::new();

        set_a.add("device-1");
        set_b.add("device-2");

        set_a.merge(&set_b);

        // 양쪽 모두 포함
        assert!(set_a.contains("device-1"));
        assert!(set_a.contains("device-2"));
    }

    #[test]
    fn test_or_set_병합_동시_추가_제거_add_wins() {
        let mut set_a = ORSet::new();
        let mut set_b = ORSet::new();

        // A에서 추가 후 B로 병합
        set_a.add("item");
        set_b.merge(&set_a);

        // B에서 제거
        set_b.remove("item");
        assert!(!set_b.contains("item"));

        // A에서 다시 추가 (새 태그)
        set_a.add("item");

        // 병합: A의 새 태그는 B의 tombstone에 없으므로 살아남음 (add-wins)
        set_b.merge(&set_a);
        assert!(set_b.contains("item"));
    }

    #[test]
    fn test_or_set_빈_셋() {
        let set = ORSet::new();
        assert!(set.is_empty());
        assert_eq!(set.len(), 0);
        assert!(!set.contains("anything"));
    }

    // =========================================================================
    // SyncManager 테스트
    // =========================================================================

    #[test]
    fn test_sync_manager_큐_추가() {
        let mut manager = SyncManager::new("device-1");

        let id = manager
            .enqueue(SyncOperation::MeasurementUpload, vec![1, 2, 3])
            .unwrap();

        assert!(!id.is_empty());
        assert_eq!(manager.pending_count(), 1);
        assert_eq!(manager.queue_size(), 1);
    }

    #[tokio::test]
    async fn test_sync_manager_동기화_실행() {
        let mut manager = SyncManager::new("device-1");
        manager.set_connected(true);

        manager
            .enqueue(SyncOperation::MeasurementUpload, vec![1, 2, 3])
            .unwrap();
        manager
            .enqueue(SyncOperation::DeviceStatusUpdate, vec![4, 5, 6])
            .unwrap();

        assert_eq!(manager.pending_count(), 2);

        let result = manager.sync().await.unwrap();

        assert_eq!(result.synced_count, 2);
        assert_eq!(result.failed_count, 0);
        assert_eq!(manager.pending_count(), 0);
    }

    #[tokio::test]
    async fn test_sync_manager_오프라인_동기화_실패() {
        let mut manager = SyncManager::new("device-1");
        // 연결 안 됨

        manager
            .enqueue(SyncOperation::MeasurementUpload, vec![1, 2, 3])
            .unwrap();

        let result = manager.sync().await;
        assert!(result.is_err());
        assert!(matches!(result.unwrap_err(), SyncError::NoConnection));

        // 큐에 여전히 있음
        assert_eq!(manager.pending_count(), 1);
    }

    #[test]
    fn test_sync_manager_측정_카운터() {
        let mut manager = SyncManager::new("device-1");

        manager.increment_measurement_count();
        manager.increment_measurement_count();
        manager.increment_measurement_count();

        assert_eq!(manager.total_measurement_count(), 3);
    }

    #[test]
    fn test_sync_manager_디바이스_관리() {
        let mut manager = SyncManager::new("device-1");

        manager.add_device("reader-001");
        manager.add_device("reader-002");

        assert!(manager.has_device("reader-001"));
        assert!(manager.has_device("reader-002"));
        assert_eq!(manager.devices().len(), 2);

        manager.remove_device("reader-001");
        assert!(!manager.has_device("reader-001"));
    }

    #[test]
    fn test_sync_manager_크로스_디바이스_병합() {
        // 디바이스 1
        let mut manager1 = SyncManager::new("device-1");
        manager1.increment_measurement_count();
        manager1.increment_measurement_count();
        manager1.add_device("reader-A");

        // 디바이스 2
        let mut manager2 = SyncManager::new("device-2");
        manager2.increment_measurement_count();
        manager2.add_device("reader-B");

        // 디바이스 1에서 디바이스 2의 상태 병합
        manager1.merge_measurement_counter(manager2.measurement_counter());
        manager1.merge_device_set(manager2.device_set());

        // 측정 횟수: device-1(2) + device-2(1) = 3
        assert_eq!(manager1.total_measurement_count(), 3);

        // 디바이스 목록: reader-A + reader-B
        assert!(manager1.has_device("reader-A"));
        assert!(manager1.has_device("reader-B"));
    }

    #[test]
    fn test_sync_manager_재시도() {
        let mut manager = SyncManager::new("device-1");

        manager
            .enqueue(SyncOperation::MeasurementUpload, vec![1, 2, 3])
            .unwrap();

        // 수동으로 에러 상태 설정 (실제로는 네트워크 실패 시)
        if let Some(item) = manager.queue.first_mut() {
            item.state = SyncState::Error;
        }

        manager.retry_failed();

        // Pending으로 복구됨
        assert_eq!(manager.pending_count(), 1);
        assert_eq!(manager.queue.first().unwrap().retry_count, 1);
    }
}
