import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';

/// 채팅 메시지 모델
class ChatMessage {
  final String id;
  final String role; // "user", "assistant", "system"
  final String content;
  final DateTime timestamp;
  final String? sessionId;

  const ChatMessage({
    this.id = '',
    required this.role,
    required this.content,
    required this.timestamp,
    this.sessionId,
  });

  Map<String, dynamic> toJson() => {
        'id': id,
        'role': role,
        'content': content,
        'timestamp': timestamp.toIso8601String(),
        if (sessionId != null) 'session_id': sessionId,
      };

  factory ChatMessage.fromJson(Map<String, dynamic> json) => ChatMessage(
        id: json['id'] as String? ?? '',
        role: json['role'] as String? ?? 'user',
        content: json['content'] as String? ?? '',
        timestamp: DateTime.tryParse(json['timestamp'] as String? ?? '') ?? DateTime.now(),
        sessionId: json['session_id'] as String?,
      );
}

/// 채팅 상태 모델
class ChatState {
  final List<ChatMessage> messages;
  final bool isLoading;
  final bool isStreaming;
  final String streamingContent;
  final String? error;

  const ChatState({
    this.messages = const [],
    this.isLoading = false,
    this.isStreaming = false,
    this.streamingContent = '',
    this.error,
  });

  ChatState copyWith({
    List<ChatMessage>? messages,
    bool? isLoading,
    bool? isStreaming,
    String? streamingContent,
    String? error,
  }) {
    return ChatState(
      messages: messages ?? this.messages,
      isLoading: isLoading ?? this.isLoading,
      isStreaming: isStreaming ?? this.isStreaming,
      streamingContent: streamingContent ?? this.streamingContent,
      error: error,
    );
  }
}

/// AI 건강 어시스턴트 채팅 Notifier
///
/// gRPC AIInferenceService와 연동하며, 서버 미연결 시 로컬 fallback 응답 제공.
/// 웹 플랫폼에서는 REST Gateway를 통해 AI 서비스 호출.
/// Hive 기반 메시지 캐싱으로 오프라인 접근 지원.
class ChatNotifier extends StateNotifier<ChatState> {
  ChatNotifier(this._manager, this._restClient, this._accessTokenProvider)
      : super(const ChatState()) {
    _initCache();
  }

  final GrpcClientManager _manager;
  final ManPaSikRestClient _restClient;
  final String? Function() _accessTokenProvider;

  static const _cacheBoxName = 'chat_messages';
  Box<String>? _cacheBox;
  String? _sessionId;

  /// Hive 캐시 초기화 및 이전 메시지 복원
  Future<void> _initCache() async {
    try {
      _cacheBox = await Hive.openBox<String>(_cacheBoxName);
      final cached = _cacheBox?.values.toList() ?? [];
      if (cached.isNotEmpty) {
        final messages = cached
            .map((json) {
              try {
                return ChatMessage.fromJson(
                    jsonDecode(json) as Map<String, dynamic>);
              } catch (_) {
                return null;
              }
            })
            .whereType<ChatMessage>()
            .toList();
        if (messages.isNotEmpty) {
          state = state.copyWith(messages: messages);
        }
      }
    } catch (e) {
      debugPrint('[ChatNotifier] 캐시 초기화 실패: $e');
    }
  }

  /// 메시지를 Hive 캐시에 저장
  Future<void> _cacheMessages(List<ChatMessage> messages) async {
    try {
      final box = _cacheBox;
      if (box == null) return;
      await box.clear();
      for (var i = 0; i < messages.length; i++) {
        await box.put('msg_$i', jsonEncode(messages[i].toJson()));
      }
    } catch (e) {
      debugPrint('[ChatNotifier] 캐시 저장 실패: $e');
    }
  }

  /// REST 코칭 세션 생성
  Future<String> _ensureSession() async {
    if (_sessionId != null) return _sessionId!;
    try {
      final resp = await _restClient.dio.post(
        '/coaching/sessions',
        data: {'type': 'health_chat'},
      );
      final data = resp.data as Map<String, dynamic>? ?? {};
      _sessionId = data['session_id'] as String? ??
          data['id'] as String? ??
          'session-${DateTime.now().millisecondsSinceEpoch}';
      return _sessionId!;
    } on DioException {
      _sessionId = 'local-${DateTime.now().millisecondsSinceEpoch}';
      return _sessionId!;
    }
  }

  /// 사용자 메시지 전송 → AI 응답 수신
  Future<void> sendMessage(String text) async {
    if (text.trim().isEmpty) return;

    // 사용자 메시지 추가
    final userMessage = ChatMessage(
      id: 'user-${DateTime.now().millisecondsSinceEpoch}',
      role: 'user',
      content: text.trim(),
      timestamp: DateTime.now(),
      sessionId: _sessionId,
    );
    state = state.copyWith(
      messages: [...state.messages, userMessage],
      isLoading: true,
      error: null,
    );

    try {
      final response = await _callCoachingEndpoint(text.trim());
      final aiMessage = ChatMessage(
        id: 'ai-${DateTime.now().millisecondsSinceEpoch}',
        role: 'assistant',
        content: response,
        timestamp: DateTime.now(),
        sessionId: _sessionId,
      );
      state = state.copyWith(
        messages: [...state.messages, aiMessage],
        isLoading: false,
      );
      await _cacheMessages(state.messages);
    } catch (e) {
      debugPrint('[ChatNotifier] AI 호출 실패: $e');
      // fallback 로컬 응답
      final fallback = _generateFallbackResponse(text.trim());
      final fallbackMessage = ChatMessage(
        id: 'fallback-${DateTime.now().millisecondsSinceEpoch}',
        role: 'assistant',
        content: fallback,
        timestamp: DateTime.now(),
        sessionId: _sessionId,
      );
      state = state.copyWith(
        messages: [...state.messages, fallbackMessage],
        isLoading: false,
        error: 'AI 서비스 일시적 오류. 기본 정보로 응답합니다.',
      );
      await _cacheMessages(state.messages);
    }
  }

  /// REST 코칭 메시지 엔드포인트를 통한 AI 서비스 호출
  Future<String> _callCoachingEndpoint(String userText) async {
    final sessionId = await _ensureSession();
    try {
      final resp = await _restClient.dio.post(
        '/coaching/messages',
        data: {
          'session_id': sessionId,
          'content': userText,
          'role': 'user',
        },
      );
      final data = resp.data as Map<String, dynamic>? ?? {};
      final reply = data['reply'] as String? ??
          data['content'] as String? ??
          data['message'] as String? ??
          '';
      if (reply.isNotEmpty) return reply;
      // Fall through to legacy AI service
    } on DioException {
      // Coaching endpoint unavailable, fall through to legacy
    } on FormatException catch (e) {
      debugPrint('[ChatNotifier] JSON 파싱 오류: $e');
      // Fall through to legacy AI service
    } on TypeError catch (e) {
      debugPrint('[ChatNotifier] 응답 타입 오류: $e');
      // Fall through to legacy AI service
    }
    return _callAiService(userText);
  }

  /// AI 서비스 호출 시도 (웹: REST, 네이티브: gRPC)
  Future<String> _callAiService(String userText) async {
    if (kIsWeb) {
      return _callAiServiceRest(userText);
    }
    return _callAiServiceGrpc(userText);
  }

  /// REST API를 통한 AI 서비스 호출
  Future<String> _callAiServiceRest(String userText) async {
    final result = await _restClient.analyzeMeasurement(
      userId: 'chat-user',
      measurementId: userText,
    );
    final summary = result['summary'] as String? ?? '';
    if (summary.isNotEmpty) return summary;
    throw Exception('빈 응답');
  }

  /// gRPC를 통한 AI 서비스 호출
  Future<String> _callAiServiceGrpc(String userText) async {
    final token = _accessTokenProvider();
    final interceptors =
        token != null ? [AuthInterceptor(() => token)] : <AuthInterceptor>[];

    final client = AiInferenceServiceClient(
      _manager.aiInferenceChannel,
      interceptors: interceptors,
    );

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

    if (lower.contains('혈당') ||
        lower.contains('blood sugar') ||
        lower.contains('glucose')) {
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

    if (lower.contains('운동') ||
        lower.contains('exercise') ||
        lower.contains('workout')) {
      return '건강한 운동 습관에 대해 알려드릴게요.\n\n'
          '세계보건기구(WHO) 권장:\n'
          '• 주 150~300분 중강도 유산소 운동\n'
          '• 또는 주 75~150분 고강도 유산소 운동\n'
          '• 주 2회 이상 근력 운동\n\n'
          '꾸준한 운동은 혈당, 혈압 관리에도 큰 도움이 됩니다.\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('식단') ||
        lower.contains('diet') ||
        lower.contains('음식') ||
        lower.contains('food')) {
      return '건강한 식단 관리에 대해 알려드릴게요.\n\n'
          '• 채소와 과일을 충분히 섭취하세요\n'
          '• 정제 탄수화물보다 통곡물을 선택하세요\n'
          '• 단백질을 적정량 섭취하세요\n'
          '• 가공식품과 나트륨 섭취를 줄이세요\n'
          '• 수분을 충분히 섭취하세요 (하루 1.5~2L)\n\n'
          '💡 현재 AI 서버에 연결되지 않아 기본 정보를 표시하고 있습니다.';
    }

    if (lower.contains('수면') ||
        lower.contains('sleep') ||
        lower.contains('잠')) {
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

  /// 스트리밍 방식으로 메시지 전송 (C1)
  ///
  /// REST Gateway의 /coaching/messages 엔드포인트를 호출하고,
  /// 응답을 글자 단위로 스트리밍하여 실시간 표시합니다.
  Future<void> sendMessageStream(String text) async {
    if (text.trim().isEmpty) return;

    final userMessage = ChatMessage(
      id: 'user-${DateTime.now().millisecondsSinceEpoch}',
      role: 'user',
      content: text.trim(),
      timestamp: DateTime.now(),
      sessionId: _sessionId,
    );
    state = state.copyWith(
      messages: [...state.messages, userMessage],
      isStreaming: true,
      streamingContent: '',
      error: null,
    );

    try {
      // REST 코칭 엔드포인트를 통해 응답 수신 후 글자 단위 스트리밍 시뮬레이션
      final response = await _callCoachingEndpoint(text.trim());
      // 글자 단위 스트리밍 시뮬레이션
      final buffer = StringBuffer();
      for (var i = 0; i < response.length; i++) {
        buffer.write(response[i]);
        state = state.copyWith(streamingContent: buffer.toString());
        await Future.delayed(const Duration(milliseconds: 15));
      }

      final aiMessage = ChatMessage(
        id: 'ai-${DateTime.now().millisecondsSinceEpoch}',
        role: 'assistant',
        content: response,
        timestamp: DateTime.now(),
        sessionId: _sessionId,
      );
      state = state.copyWith(
        messages: [...state.messages, aiMessage],
        isStreaming: false,
        streamingContent: '',
      );
      await _cacheMessages(state.messages);
    } catch (e) {
      debugPrint('[ChatNotifier] AI 스트리밍 실패: $e');
      final fallback = _generateFallbackResponse(text.trim());

      // fallback도 스트리밍으로 표시
      final buffer = StringBuffer();
      for (var i = 0; i < fallback.length; i++) {
        buffer.write(fallback[i]);
        state = state.copyWith(streamingContent: buffer.toString());
        await Future.delayed(const Duration(milliseconds: 10));
      }

      final fallbackMessage = ChatMessage(
        id: 'fallback-${DateTime.now().millisecondsSinceEpoch}',
        role: 'assistant',
        content: fallback,
        timestamp: DateTime.now(),
        sessionId: _sessionId,
      );
      state = state.copyWith(
        messages: [...state.messages, fallbackMessage],
        isStreaming: false,
        streamingContent: '',
        error: 'AI 서비스 일시적 오류. 기본 정보로 응답합니다.',
      );
      await _cacheMessages(state.messages);
    }
  }

  /// 채팅 기록 초기화
  void clearChat() {
    state = const ChatState();
    _sessionId = null;
    _cacheBox?.clear();
  }
}

/// 채팅 상태 Provider
final chatProvider = StateNotifierProvider<ChatNotifier, ChatState>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  final restClient = ref.watch(restClientProvider);
  return ChatNotifier(
    manager,
    restClient,
    () => ref.read(authProvider).accessToken,
  );
});
