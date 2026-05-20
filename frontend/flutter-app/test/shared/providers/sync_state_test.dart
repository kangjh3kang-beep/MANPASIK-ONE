import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/providers/sync_provider.dart';

void main() {
  group('SyncStatus enum', () {
    test('값 목록 확인 (4개)', () {
      expect(SyncStatus.values, hasLength(4));
      expect(SyncStatus.values, contains(SyncStatus.idle));
      expect(SyncStatus.values, contains(SyncStatus.syncing));
      expect(SyncStatus.values, contains(SyncStatus.error));
      expect(SyncStatus.values, contains(SyncStatus.offline));
    });
  });

  group('SyncState', () {
    test('기본값 확인', () {
      const state = SyncState();
      expect(state.status, SyncStatus.idle);
      expect(state.pendingCount, 0);
      expect(state.syncedCount, 0);
      expect(state.failedCount, 0);
      expect(state.lastError, isNull);
      expect(state.lastSyncAt, isNull);
      expect(state.hasConflicts, isFalse);
    });

    test('커스텀 상태 생성', () {
      final now = DateTime.now();
      final state = SyncState(
        status: SyncStatus.syncing,
        pendingCount: 5,
        syncedCount: 10,
        failedCount: 2,
        lastError: '네트워크 오류',
        lastSyncAt: now,
        hasConflicts: true,
      );
      expect(state.status, SyncStatus.syncing);
      expect(state.pendingCount, 5);
      expect(state.syncedCount, 10);
      expect(state.failedCount, 2);
      expect(state.lastError, '네트워크 오류');
      expect(state.lastSyncAt, now);
      expect(state.hasConflicts, isTrue);
    });

    test('copyWith로 부분 업데이트', () {
      const initial = SyncState(
        status: SyncStatus.idle,
        pendingCount: 3,
      );

      final updated = initial.copyWith(
        status: SyncStatus.syncing,
        syncedCount: 1,
      );
      expect(updated.status, SyncStatus.syncing);
      expect(updated.pendingCount, 3); // 유지
      expect(updated.syncedCount, 1); // 업데이트
    });

    test('copyWith로 hasConflicts 토글', () {
      const state = SyncState(hasConflicts: false);
      final withConflicts = state.copyWith(hasConflicts: true);
      expect(withConflicts.hasConflicts, isTrue);
    });

    test('copyWith로 에러 상태 설정', () {
      const state = SyncState();
      final errorState = state.copyWith(
        status: SyncStatus.error,
        lastError: '서버 연결 실패',
        failedCount: 1,
      );
      expect(errorState.status, SyncStatus.error);
      expect(errorState.lastError, '서버 연결 실패');
      expect(errorState.failedCount, 1);
    });

    test('copyWith 원본 불변 확인', () {
      final now = DateTime.now();
      final original = SyncState(
        status: SyncStatus.idle,
        lastSyncAt: now,
      );
      original.copyWith(status: SyncStatus.syncing);

      // 원본은 변경되지 않아야 한다
      expect(original.status, SyncStatus.idle);
      expect(original.lastSyncAt, now);
    });
  });
}
