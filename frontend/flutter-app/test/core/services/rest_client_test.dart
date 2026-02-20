import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('ManPaSikRestClient', () {
    test('ManPaSikRestClient를 baseUrl로 생성할 수 있다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:8080/api/v1');
      expect(client, isNotNull);
    });

    test('다른 baseUrl로 인스턴스를 생성할 수 있다', () {
      final client1 = ManPaSikRestClient(baseUrl: 'http://server1:8080/api/v1');
      final client2 = ManPaSikRestClient(baseUrl: 'http://server2:9090/api/v1');
      expect(client1, isNotNull);
      expect(client2, isNotNull);
    });

    test('빈 baseUrl로도 생성 가능하다', () {
      final client = ManPaSikRestClient(baseUrl: '');
      expect(client, isNotNull);
    });
  });
}
