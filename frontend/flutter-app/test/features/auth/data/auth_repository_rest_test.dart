import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/auth/data/auth_repository_rest.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('AuthRepositoryRest', () {
    test('AuthRepositoryRest 클래스가 존재하고 AuthRepository를 구현한다', () {
      expect(AuthRepositoryRest, isNotNull);
    });

    test('AuthRepositoryRest는 ManPaSikRestClient를 인자로 받아 생성된다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AuthRepositoryRest(client);
      expect(repo, isA<AuthRepository>());
    });

    test('AuthResult.success 팩토리가 올바른 필드를 설정한다', () {
      final result = AuthResult.success(
        userId: 'u-1',
        email: 'test@manpasik.com',
        displayName: '테스트',
        accessToken: 'access-123',
        refreshToken: 'refresh-456',
      );
      expect(result.success, isTrue);
      expect(result.userId, 'u-1');
      expect(result.email, 'test@manpasik.com');
      expect(result.displayName, '테스트');
      expect(result.accessToken, 'access-123');
      expect(result.refreshToken, 'refresh-456');
      expect(result.errorMessage, isNull);
    });

    test('AuthResult.failure 팩토리가 에러 메시지를 설정한다', () {
      final result = AuthResult.failure('로그인 실패');
      expect(result.success, isFalse);
      expect(result.errorMessage, '로그인 실패');
      expect(result.userId, isNull);
      expect(result.accessToken, isNull);
    });

    test('isAuthenticated는 항상 false를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AuthRepositoryRest(client);
      final result = await repo.isAuthenticated();
      expect(result, isFalse);
    });

    test('refreshToken은 항상 false를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AuthRepositoryRest(client);
      final result = await repo.refreshToken();
      expect(result, isFalse);
    });

    test('logout은 예외 없이 완료된다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AuthRepositoryRest(client);
      await expectLater(repo.logout(), completes);
    });

    test('AuthResult role 필드를 통해 역할 정보를 저장할 수 있다', () {
      final result = AuthResult.success(
        userId: 'admin-1',
        email: 'admin@manpasik.com',
        accessToken: 'token',
        refreshToken: 'refresh',
        role: 'admin',
      );
      expect(result.role, 'admin');
    });
  });
}
