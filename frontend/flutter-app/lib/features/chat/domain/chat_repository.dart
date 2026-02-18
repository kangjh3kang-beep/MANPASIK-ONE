/// 채팅 도메인 모델 및 리포지토리
///
/// AI 건강 어시스턴트 대화 관리

/// 대화 세션
class ChatSession {
  final String id;
  final String userId;
  final String title;
  final DateTime createdAt;
  final DateTime? lastMessageAt;

  const ChatSession({
    required this.id,
    required this.userId,
    required this.title,
    required this.createdAt,
    this.lastMessageAt,
  });
}

/// 채팅 리포지토리 인터페이스
abstract class ChatRepository {
  Future<List<ChatSession>> getSessions();
  Future<ChatSession> createSession({String? title});
  Future<void> deleteSession(String sessionId);
  Future<String> sendMessage(String sessionId, String text);
  Stream<String> sendMessageStream(String sessionId, String text);
}
