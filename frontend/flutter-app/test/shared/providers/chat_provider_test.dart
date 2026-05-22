// ChatNotifier (Provider) 단위 테스트
// - Mock 기반 sendMessage / clearChat / error state 검증

import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:mocktail/mocktail.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/shared/providers/chat_provider.dart';

class MockGrpcClientManager extends Mock implements GrpcClientManager {}

class MockManPaSikRestClient extends Mock implements ManPaSikRestClient {}

class MockDio extends Mock implements Dio {}

void main() {
  late MockGrpcClientManager mockGrpc;
  late MockManPaSikRestClient mockRest;
  late MockDio mockDio;
  late ChatNotifier notifier;

  setUp(() async {
    final dir = await Directory.systemTemp.createTemp('chat_test');
    Hive.init(dir.path);

    mockGrpc = MockGrpcClientManager();
    mockRest = MockManPaSikRestClient();
    mockDio = MockDio();
    when(() => mockRest.dio).thenReturn(mockDio);

    notifier = ChatNotifier(mockGrpc, mockRest, () => null);
  });

  group('ChatNotifier', () {
    test('sendMessage adds user message to state', () async {
      // Simulate DioException so it falls through to fallback
      when(() => mockDio.post(
            any(),
            data: any(named: 'data'),
          )).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/coaching/sessions'),
        type: DioExceptionType.connectionError,
      ));

      await notifier.sendMessage('안녕하세요');

      // First message should be from the user
      expect(notifier.state.messages.length, greaterThanOrEqualTo(1));
      final userMsg = notifier.state.messages.firstWhere(
        (m) => m.role == 'user',
      );
      expect(userMsg.content, '안녕하세요');
    });

    test('sendMessage adds AI response on success', () async {
      // Mock session creation
      when(() => mockDio.post(
            '/coaching/sessions',
            data: any(named: 'data'),
          )).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/coaching/sessions'),
            statusCode: 200,
            data: {'session_id': 'test-session'},
          ));

      // Mock coaching message response
      when(() => mockDio.post(
            '/coaching/messages',
            data: any(named: 'data'),
          )).thenAnswer((_) async => Response(
            requestOptions: RequestOptions(path: '/coaching/messages'),
            statusCode: 200,
            data: {'reply': '건강 상담 답변입니다.'},
          ));

      await notifier.sendMessage('혈당 정보 알려주세요');

      expect(notifier.state.isLoading, isFalse);
      expect(notifier.state.messages.length, 2);

      final aiMsg = notifier.state.messages.last;
      expect(aiMsg.role, 'assistant');
      expect(aiMsg.content, '건강 상담 답변입니다.');
    });

    test('sendMessage sets error state on failure', () async {
      // Make all network calls fail
      when(() => mockDio.post(
            any(),
            data: any(named: 'data'),
          )).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        type: DioExceptionType.connectionTimeout,
      ));

      await notifier.sendMessage('테스트 메시지');

      expect(notifier.state.isLoading, isFalse);
      expect(notifier.state.error, isNotNull);
      expect(notifier.state.error, contains('AI 서비스 일시적 오류'));

      // Should still have fallback message
      final fallbackMsg = notifier.state.messages
          .where((m) => m.role == 'assistant')
          .toList();
      expect(fallbackMsg, isNotEmpty);
    });

    test('clearChat resets messages and session', () async {
      // First send a message to populate state (with fallback)
      when(() => mockDio.post(
            any(),
            data: any(named: 'data'),
          )).thenThrow(DioException(
        requestOptions: RequestOptions(path: '/'),
        type: DioExceptionType.connectionError,
      ));

      await notifier.sendMessage('첫 번째 메시지');
      expect(notifier.state.messages, isNotEmpty);

      notifier.clearChat();

      expect(notifier.state.messages, isEmpty);
      expect(notifier.state.isLoading, isFalse);
      expect(notifier.state.isStreaming, isFalse);
      expect(notifier.state.error, isNull);
    });

    test('empty message is ignored', () async {
      await notifier.sendMessage('');
      expect(notifier.state.messages, isEmpty);
      expect(notifier.state.isLoading, isFalse);

      await notifier.sendMessage('   ');
      expect(notifier.state.messages, isEmpty);
    });
  });

  group('ChatMessage', () {
    test('toJson / fromJson round-trip', () {
      final msg = ChatMessage(
        id: 'test-1',
        role: 'user',
        content: '혈당 체크',
        timestamp: DateTime(2026, 5, 22, 10, 30),
        sessionId: 'sess-1',
      );

      final json = msg.toJson();
      final restored = ChatMessage.fromJson(json);

      expect(restored.id, msg.id);
      expect(restored.role, msg.role);
      expect(restored.content, msg.content);
      expect(restored.sessionId, msg.sessionId);
    });
  });

  group('ChatState', () {
    test('copyWith preserves unmodified fields', () {
      const original = ChatState(
        messages: [],
        isLoading: false,
        isStreaming: false,
      );

      final updated = original.copyWith(isLoading: true);
      expect(updated.isLoading, isTrue);
      expect(updated.isStreaming, isFalse);
      expect(updated.messages, isEmpty);
    });

    test('copyWith with error=null clears error', () {
      final withError = const ChatState().copyWith(
        error: '오류 발생',
      );
      expect(withError.error, '오류 발생');

      // ChatState.copyWith sets error to the parameter directly (not ?? pattern)
      final cleared = withError.copyWith();
      expect(cleared.error, isNull);
    });
  });
}
