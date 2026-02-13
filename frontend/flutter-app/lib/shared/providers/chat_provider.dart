import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';

/// 채팅 메시지 모델
class ChatMessage {
  final String role; // "user", "assistant", "system"
  final String content;
  final DateTime timestamp;

  const ChatMessage({
    required this.role,
    required this.content,
    required this.timestamp,
  });
}

/// 채팅 상태 모델
class ChatState {
  final List<ChatMessage> messages;
  final bool isLoading;
  final String? error;

  const ChatState({
    this.messages = const [],
    this.isLoading = false,
    this.error,
  });

  ChatState copyWith({
    List<ChatMessage>? messages,
    bool? isLoading,
    String? error,
  }) {
    return ChatState(
      messages: messages ?? this.messages,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// AI 건강 어시스턴트 채팅 Notifier
///
/// gRPC AIInferenceService와 연동하며, 서버 미연결 시 로컬 fallback 응답 제공.
class ChatNotifier extends StateNotifier<ChatState> {
  ChatNotifier(this._manager, this._accessTokenProvider)
      : super(const ChatState());

  final GrpcClientManager _manager;
  final String? Function() _accessTokenProvider;

  /// 사용자 메시지 전송 → AI 응답 수신
  Future<void> sendMessage(String text) async {
    if (text.trim().isEmpty) return;

    // 사용자 메시지 추가
    final userMessage = ChatMessage(
      role: 'user',
      content: text.trim(),
      timestamp: DateTime.now(),
    );
    state = state.copyWith(
      messages: [...state.messages, userMessage],
      isLoading: true,
      error: null,
    );

    try {
      final response = await _callAiService(text.trim());
      final aiMessage = ChatMessage(
        role: 'assistant',
        content: response,
        timestamp: DateTime.now(),
      );
      state = state.copyWith(
        messages: [...state.messages, aiMessage],
        isLoading: false,
      );
    } catch (e) {
      debugPrint('[ChatNotifier] AI 호출 실패: $e');
      // fallback 로컬 응답
      final fallback = _generateFallbackResponse(text.trim());
      final fallbackMessage = ChatMessage(
        role: 'assistant',
        content: fallback,
        timestamp: DateTime.now(),
      );
      state = state.copyWith(
        messages: [...state.messages, fallbackMessage],
        isLoading: false,
      );
    }
  }

  /// gRPC AI 서비스 호출 시도
  Future<String> _callAiService(String userText) async {
    final token = _accessTokenProvider();
    final interceptors = token != null
        ? [AuthInterceptor(() => token)]
        : <AuthInterceptor>[];

    final client = AIInferenceServiceClient(
      _manager.aiInferenceChannel,
      interceptors: interceptors,
    );

    // AnalyzeMeasurement를 텍스트 기반 건강 질문에 활용
    // measurementId에 사용자 질문 텍스트를 전달 (서버에서 LLM으로 처리)
    final request = AnalyzeMeasurementRequest()
      ..userId = 'chat-user'
      ..measurementId = userText;

    final result = await client.analyzeMeasurement(request);
    if (result.summary.isNotEmpty) {
      return result.summary;
    }
    throw Exception('빈 응답');
  }

  /// gRPC 미연결 시 로컬 fallback 응답 생성
  String _generateFallbackResponse(String userText) {
    final lower = userText.toLowerCase();

    if (lower.contains('혈당') || lower.contains('blood sugar') || lower.contains('glucose')) {
      return '혈당 관리에 대해 물어봐 주셨네요.\n\n'
          '일반적인 공복 혈당 정상 범위는 70~100 mg/dL입니다. '
          '식후 2시간 기준 140 mg/dL 미만이 정상이에요.\n\n'
          '⚠️ 이 정보는 일반적인 참고 사항입니다. '
          '정확한 진단은 전문 의료인과 상담해주세요.\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('혈압') || lower.contains('blood pressure')) {
      return '혈압에 대해 알려드릴게요.\n\n'
          '정상 혈압: 수축기 120mmHg 미만 / 이완기 80mmHg 미만\n'
          '주의 혈압: 120-139 / 80-89 mmHg\n'
          '고혈압: 140/90 mmHg 이상\n\n'
          '⚠️ 일반적인 참고 정보입니다. 전문의 상담을 권장합니다.\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('운동') || lower.contains('exercise') || lower.contains('workout')) {
      return '건강한 운동 습관에 대해 알려드릴게요.\n\n'
          '세계보건기구(WHO) 권장:\n'
          '• 주 150~300분 중강도 유산소 운동\n'
          '• 또는 주 75~150분 고강도 유산소 운동\n'
          '• 주 2회 이상 근력 운동\n\n'
          '꾸준한 운동은 혈당, 혈압 관리에도 큰 도움이 됩니다.\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('식단') || lower.contains('diet') || lower.contains('음식') || lower.contains('food')) {
      return '건강한 식단 관리에 대해 알려드릴게요.\n\n'
          '• 채소와 과일을 충분히 섭취하세요\n'
          '• 정제 탄수화물보다 통곡물을 선택하세요\n'
          '• 단백질을 적정량 섭취하세요\n'
          '• 가공식품과 나트륨 섭취를 줄이세요\n'
          '• 수분을 충분히 섭취하세요 (하루 1.5~2L)\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('수면') || lower.contains('sleep') || lower.contains('잠')) {
      return '건강한 수면에 대해 알려드릴게요.\n\n'
          '성인 기준 하루 7~9시간 수면이 권장됩니다.\n\n'
          '좋은 수면 습관:\n'
          '• 일정한 취침/기상 시간 유지\n'
          '• 취침 전 카페인, 알코올 피하기\n'
          '• 침실은 어둡고 시원하게\n'
          '• 취침 1시간 전 스마트폰 사용 줄이기\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    // 기본 응답
    return '안녕하세요! 건강 관련 질문에 답변해 드리겠습니다.\n\n'
        '현재 AI 서버에 연결되지 않아 상세한 분석이 어렵습니다. '
        '서버가 연결되면 더 정확하고 개인화된 건강 인사이트를 제공해 드릴 수 있어요.\n\n'
        '일반적인 건강 질문(혈당, 혈압, 운동, 식단, 수면 등)에 대해서는 '
        '기본적인 정보를 제공할 수 있으니 편하게 물어보세요!';
  }

  /// 채팅 기록 초기화
  void clearChat() {
    state = const ChatState();
  }
}

/// 채팅 상태 Provider
final chatProvider = StateNotifierProvider<ChatNotifier, ChatState>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return ChatNotifier(
    manager,
    () => ref.read(authProvider).accessToken,
  );
});
