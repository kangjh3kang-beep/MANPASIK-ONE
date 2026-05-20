import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';

/// 테스트용 가짜 클라이언트 — getMyMemberships 만 오버라이드.
///
/// 다른 메서드는 실제 호출 시 예외 발생 (테스트에서 사용하지 않음).
class _FakeClient extends ManPaSikRestClient {
  _FakeClient(this._mock) : super(baseUrl: 'http://localhost');
  final List<Map<String, dynamic>> _mock;

  @override
  Future<List<Map<String, dynamic>>> getMyMemberships() async => _mock;
}

class _FailingClient extends ManPaSikRestClient {
  _FailingClient(this._error) : super(baseUrl: 'http://localhost');
  final Object _error;

  @override
  Future<List<Map<String, dynamic>>> getMyMemberships() async => throw _error;
}

void main() {
  group('MembershipDto', () {
    test('fromJson 정상 변환', () {
      final dto = MembershipDto.fromJson({
        'user_id': 'u1',
        'tenant_id': 't1',
        'role': 'admin',
        'active': true,
        'joined_at': 1730000000,
      });
      expect(dto.userId, 'u1');
      expect(dto.tenantId, 't1');
      expect(dto.role, 'admin');
      expect(dto.active, true);
      expect(dto.joinedAt, 1730000000);
    });

    test('fromJson 누락 필드 기본값', () {
      final dto = MembershipDto.fromJson({});
      expect(dto.userId, '');
      expect(dto.tenantId, '');
      expect(dto.active, false);
      expect(dto.joinedAt, 0);
    });

    test('fromJson joined_at 문자열 처리', () {
      final dto = MembershipDto.fromJson({
        'joined_at': '1730000000',
      });
      expect(dto.joinedAt, 1730000000);
    });

    test('equality + hashCode', () {
      const a = MembershipDto(
          userId: 'u', tenantId: 't', role: 'r', active: true, joinedAt: 1);
      const b = MembershipDto(
          userId: 'u', tenantId: 't', role: 'r', active: true, joinedAt: 1);
      expect(a == b, true);
      expect(a.hashCode, b.hashCode);
    });
  });

  group('myMembershipsProvider', () {
    test('REST 응답 → DTO 리스트 변환', () async {
      final container = ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(_FakeClient([
          {
            'user_id': 'me',
            'tenant_id': 'hospA',
            'role': 'admin',
            'active': true,
            'joined_at': 1730000000,
          },
          {
            'user_id': 'me',
            'tenant_id': 'famB',
            'role': 'member',
            'active': true,
            'joined_at': 1731000000,
          },
        ])),
      ]);
      addTearDown(container.dispose);

      final list = await container.read(myMembershipsProvider.future);
      expect(list, hasLength(2));
      expect(list[0].tenantId, 'hospA');
      expect(list[1].role, 'member');
    });

    test('inactive 멤버십은 제외', () async {
      final container = ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(_FakeClient([
          {
            'user_id': 'me',
            'tenant_id': 'active-1',
            'role': 'admin',
            'active': true,
            'joined_at': 1,
          },
          {
            'user_id': 'me',
            'tenant_id': 'inactive-1',
            'role': 'admin',
            'active': false,
            'joined_at': 2,
          },
        ])),
      ]);
      addTearDown(container.dispose);

      final list = await container.read(myMembershipsProvider.future);
      expect(list, hasLength(1));
      expect(list[0].tenantId, 'active-1');
    });

    test('REST 에러 전파', () async {
      final container = ProviderContainer(overrides: [
        restClientProvider
            .overrideWithValue(_FailingClient(StateError('네트워크'))),
      ]);
      addTearDown(container.dispose);

      await expectLater(
        container.read(myMembershipsProvider.future),
        throwsA(isA<StateError>()),
      );
    });

    test('빈 응답', () async {
      final container = ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(_FakeClient([])),
      ]);
      addTearDown(container.dispose);
      final list = await container.read(myMembershipsProvider.future);
      expect(list, isEmpty);
    });
  });
}
