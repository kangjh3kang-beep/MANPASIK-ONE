// data/crdt_sync.rs

/// CRDT(Conflict-free Replicated Data Type)의 LWW(Last Writer Wins) 구현 등
/// 오프라인-우선 접근법을 위한 동기화 골격
pub struct LwwRegister<T> {
    pub value: T,
    pub timestamp: u64,
}

impl<T: Clone> LwwRegister<T> {
    pub fn new(value: T, timestamp: u64) -> Self {
        Self { value, timestamp }
    }

    /// 다른 레지스터 상태와 병합하고 더 최신의 타임스탬프 값을 우선합니다.
    pub fn merge(&mut self, other: &Self) {
        if other.timestamp > self.timestamp {
            self.value = other.value.clone();
            self.timestamp = other.timestamp;
        }
    }
}
