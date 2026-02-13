// AuthState 모델 단위 테스트
// - 기본 생성, copyWith, 인증 상태 확인 등

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

void main() {
  group('AuthState 모델 테스트', () {
    // 기본 생성자 테스트: 모든 필드가 기본값으로 초기화되는지 확인
    test('기본 생성 시 isAuthenticated는 false이어야 한다', () {
      const state = AuthState();
      expect(state.isAuthenticated, isFalse);
      expect(state.userId, isNull);
      expect(state.email, isNull);
      expect(state.displayName, isNull);
      expect(state.accessToken, isNull);
      expect(state.refreshToken, isNull);
    });

    // 명시적 인증 상태로 생성
    test('인증된 상태로 생성할 수 있어야 한다', () {
      const state = AuthState(
        isAuthenticated: true,
        userId: 'user-1',
        email: 'test@manpasik.com',
        displayName: '홍길동',
        accessToken: 'access-token-123',
        refreshToken: 'refresh-token-456',
      );
      expect(state.isAuthenticated, isTrue);
      expect(state.userId, 'user-1');
      expect(state.email, 'test@manpasik.com');
      expect(state.displayName, '홍길동');
      expect(state.accessToken, 'access-token-123');
      expect(state.refreshToken, 'refresh-token-456');
    });

    // copyWith으로 일부 필드만 변경
    test('copyWith으로 일부 필드만 변경할 수 있어야 한다', () {
      const original = AuthState(
        isAuthenticated: true,
        userId: 'user-1',
        email: 'old@manpasik.com',
        displayName: '기존 이름',
        accessToken: 'old-token',
        refreshToken: 'old-refresh',
      );

      final updated = original.copyWith(
        email: 'new@manpasik.com',
        displayName: '새 이름',
      );

      // 변경된 필드
      expect(updated.email, 'new@manpasik.com');
      expect(updated.displayName, '새 이름');
      // 변경되지 않은 필드
      expect(updated.isAuthenticated, isTrue);
      expect(updated.userId, 'user-1');
      expect(updated.accessToken, 'old-token');
      expect(updated.refreshToken, 'old-refresh');
    });

    // copyWith에 아무것도 전달하지 않으면 동일한 값 유지
    test('copyWith에 인자를 전달하지 않으면 기존 값이 유지된다', () {
      const state = AuthState(
        isAuthenticated: true,
        userId: 'u1',
        email: 'e@m.com',
      );

      final copied = state.copyWith();

      expect(copied.isAuthenticated, state.isAuthenticated);
      expect(copied.userId, state.userId);
      expect(copied.email, state.email);
    });
  });
}
