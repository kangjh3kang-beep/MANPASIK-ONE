import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/chat/data/chat_repository_rest.dart';
import 'package:manpasik/features/chat/domain/chat_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('ChatRepositoryRest', () {
    test('ChatRepositoryRest는 ChatRepository를 구현한다', () {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = ChatRepositoryRest(client, 'user-1');
      expect(repo, isA<ChatRepository>());
    });

    test('getSessions는 DioException 시 빈 리스트를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = ChatRepositoryRest(client, 'user-1');
      final sessions = await repo.getSessions();
      expect(sessions, isEmpty);
    });

    test('deleteSession은 예외 없이 실행된다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = ChatRepositoryRest(client, 'user-1');
      // deleteSession은 no-op이므로 예외가 발생하지 않아야 한다
      await repo.deleteSession('session-1');
    });
  });

  group('ChatSession 도메인 모델', () {
    test('ChatSession 생성 확인', () {
      final now = DateTime.now();
      final session = ChatSession(
        id: 'session-1',
        userId: 'user-1',
        title: 'AI 코칭',
        createdAt: now,
      );
      expect(session.id, 'session-1');
      expect(session.userId, 'user-1');
      expect(session.title, 'AI 코칭');
      expect(session.createdAt, now);
      expect(session.lastMessageAt, isNull);
    });

    test('ChatSession lastMessageAt 포함 생성', () {
      final now = DateTime.now();
      final lastMsg = now.add(const Duration(minutes: 5));
      final session = ChatSession(
        id: 'session-2',
        userId: 'user-1',
        title: '건강 상담',
        createdAt: now,
        lastMessageAt: lastMsg,
      );
      expect(session.lastMessageAt, lastMsg);
    });
  });
}
