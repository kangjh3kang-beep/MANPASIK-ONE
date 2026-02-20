import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/chat/domain/chat_repository.dart';

/// REST Gateway를 사용하는 ChatRepository 구현체
class ChatRepositoryRest implements ChatRepository {
  ChatRepositoryRest(this._client, this._userId);

  final ManPaSikRestClient _client;
  final String _userId;

  @override
  Future<List<ChatSession>> getSessions() async {
    try {
      final res = await _client.listCoachingMessages(_userId);
      final sessions = res['sessions'] as List<dynamic>? ?? [];
      return sessions.map((s) {
        final m = s as Map<String, dynamic>;
        return ChatSession(
          id: m['session_id'] as String? ?? '',
          userId: _userId,
          title: m['title'] as String? ?? 'AI 코칭',
          createdAt:
              DateTime.tryParse(m['created_at'] as String? ?? '') ??
                  DateTime.now(),
          lastMessageAt:
              m['last_message_at'] != null
                  ? DateTime.tryParse(m['last_message_at'] as String)
                  : null,
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<ChatSession> createSession({String? title}) async {
    final res = await _client.generateCoaching(userId: _userId);
    return ChatSession(
      id: res['session_id'] as String? ??
          'session-${DateTime.now().millisecondsSinceEpoch}',
      userId: _userId,
      title: title ?? 'AI 코칭',
      createdAt: DateTime.now(),
    );
  }

  @override
  Future<void> deleteSession(String sessionId) async {
    // No dedicated delete endpoint; coaching sessions are immutable
  }

  @override
  Future<String> sendMessage(String sessionId, String text) async {
    final res = await _client.streamChat(
      userId: _userId,
      message: text,
    );
    return res['response'] as String? ?? '';
  }

  @override
  Stream<String> sendMessageStream(String sessionId, String text) async* {
    // REST doesn't support true SSE; fall back to single response
    final response = await sendMessage(sessionId, text);
    yield response;
  }
}
