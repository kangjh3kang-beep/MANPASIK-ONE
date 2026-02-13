// AuthResult 도메인 모델 테스트
// - success / failure 팩토리 생성, 필드 값 검증

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';

void main() {
  group('AuthResult 모델 테스트', () {
    // success 팩토리로 생성 시 모든 필드가 올바르게 설정되는지
    test('AuthResult.success 팩토리는 success=true를 반환해야 한다', () {
      final result = AuthResult.success(
        userId: 'user-1',
        email: 'test@manpasik.com',
        displayName: '테스트 사용자',
        accessToken: 'access-123',
        refreshToken: 'refresh-456',
      );

      expect(result.success, isTrue);
      expect(result.userId, 'user-1');
      expect(result.email, 'test@manpasik.com');
      expect(result.displayName, '테스트 사용자');
      expect(result.accessToken, 'access-123');
      expect(result.refreshToken, 'refresh-456');
      expect(result.errorMessage, isNull);
    });

    // failure 팩토리로 생성 시 에러 메시지가 올바른지
    test('AuthResult.failure 팩토리는 success=false와 에러 메시지를 반환해야 한다', () {
      final result = AuthResult.failure('로그인 실패');

      expect(result.success, isFalse);
      expect(result.errorMessage, '로그인 실패');
      expect(result.userId, isNull);
      expect(result.accessToken, isNull);
      expect(result.refreshToken, isNull);
    });

    // 기본 생성자로 부분 필드 지정
    test('기본 생성자로 필요한 필드만 지정할 수 있다', () {
      const result = AuthResult(
        success: true,
        userId: 'uid',
        accessToken: 'tok',
      );

      expect(result.success, isTrue);
      expect(result.userId, 'uid');
      expect(result.email, isNull);
      expect(result.displayName, isNull);
      expect(result.refreshToken, isNull);
    });
  });
}
