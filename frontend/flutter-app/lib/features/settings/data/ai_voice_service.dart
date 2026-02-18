import 'dart:async';

import 'package:flutter/services.dart';

/// AI 음성 확인 서비스
///
/// 에스컬레이션 3단계: TTS/STT 기반 자동 음성 전화로 사용자 상태 확인.
/// Platform Channel을 통해 네이티브 TTS/STT를 연동합니다:
/// - Android: TextToSpeech + SpeechRecognizer
/// - iOS: AVSpeechSynthesizer + SFSpeechRecognizer
/// 플랫폼 채널 미등록 시 시뮬레이션 폴백으로 동작합니다.
class AiVoiceService {
  AiVoiceService._();
  static final instance = AiVoiceService._();

  static const _ttsChannel = MethodChannel('com.manpasik/tts');
  static const _sttChannel = MethodChannel('com.manpasik/stt');

  bool _isCallActive = false;
  bool get isCallActive => _isCallActive;

  /// 음성 확인 시나리오 종류
  static const scenarios = {
    'health_critical': _Scenario(
      greeting: '안녕하세요. 만파식 AI 건강 비서입니다.',
      question: '건강 이상이 감지되었습니다. 괜찮으시면 "괜찮아"라고 말씀해주세요.',
      confirmKeywords: ['괜찮', '네', '예', '좋아', '양호'],
      timeout: Duration(seconds: 30),
    ),
    'fall_detected': _Scenario(
      greeting: '안녕하세요. 만파식 안전 확인 서비스입니다.',
      question: '낙상이 감지되었습니다. 도움이 필요하시면 "도와줘"라고 말씀해주세요. 괜찮으시면 "괜찮아"라고 말씀해주세요.',
      confirmKeywords: ['괜찮', '네', '예'],
      timeout: Duration(seconds: 45),
    ),
    'no_response': _Scenario(
      greeting: '안녕하세요. 만파식 AI입니다.',
      question: '일정 시간 응답이 없어 확인 전화드렸습니다. 괜찮으시면 아무 말씀이나 해주세요.',
      confirmKeywords: ['괜찮', '네', '예', '좋아', '응', '어'],
      timeout: Duration(seconds: 60),
    ),
  };

  /// AI 음성 확인 전화를 시작합니다.
  ///
  /// 1. TTS로 인사말 + 질문 재생
  /// 2. STT로 사용자 응답 대기 (타임아웃 설정)
  /// 3. 키워드 매칭으로 확인/미확인 판단
  ///
  /// Returns: true = 사용자 확인됨, false = 미확인 (다음 단계 에스컬레이션)
  Future<bool> initiateVoiceCheck(String scenarioKey, {String? phoneNumber}) async {
    final scenario = scenarios[scenarioKey];
    if (scenario == null) return false;

    _isCallActive = true;

    try {
      // 1. TTS 재생 (실제: platform channel → native TTS)
      await _speak(scenario.greeting);
      await Future.delayed(const Duration(seconds: 1));
      await _speak(scenario.question);

      // 2. STT 대기 (실제: platform channel → native STT)
      final response = await _listenForResponse(scenario.timeout);

      // 3. 키워드 매칭
      if (response != null) {
        for (final keyword in scenario.confirmKeywords) {
          if (response.contains(keyword)) {
            await _speak('확인되었습니다. 건강에 유의하세요.');
            return true;
          }
        }
        // 응답은 있으나 키워드 불일치 → 재질문
        await _speak('다시 한번 말씀해주세요. 괜찮으시면 "괜찮아"라고 말씀해주세요.');
        final retry = await _listenForResponse(const Duration(seconds: 20));
        if (retry != null) {
          for (final keyword in scenario.confirmKeywords) {
            if (retry.contains(keyword)) {
              await _speak('확인되었습니다.');
              return true;
            }
          }
        }
      }

      // 미확인 → 에스컬레이션 다음 단계
      await _speak('응답을 확인할 수 없습니다. 보호자에게 알림을 전송합니다.');
      return false;
    } finally {
      _isCallActive = false;
    }
  }

  /// TTS 출력 (네이티브 → 시뮬레이션 폴백)
  Future<void> _speak(String text) async {
    try {
      await _ttsChannel.invokeMethod('speak', {
        'text': text,
        'language': 'ko-KR',
        'rate': 0.9,
      });
      return;
    } on MissingPluginException {
      // Platform channel 미등록 → 시뮬레이션
    } on PlatformException {
      // 네이티브 TTS 실패 → 시뮬레이션
    }
    await Future.delayed(Duration(milliseconds: text.length * 50));
  }

  /// STT 입력 대기 (네이티브 → 시뮬레이션 폴백)
  Future<String?> _listenForResponse(Duration timeout) async {
    try {
      final result = await _sttChannel.invokeMethod<String>('listen', {
        'language': 'ko-KR',
        'timeoutMs': timeout.inMilliseconds,
      });
      return result;
    } on MissingPluginException {
      // Platform channel 미등록 → 시뮬레이션
    } on PlatformException {
      // 네이티브 STT 실패 → 시뮬레이션
    }
    await Future.delayed(timeout);
    return null;
  }

  /// TTS 엔진 사용 가능 여부 확인
  Future<bool> isTtsAvailable() async {
    try {
      final result = await _ttsChannel.invokeMethod<bool>('isAvailable');
      return result ?? false;
    } catch (_) {
      return false;
    }
  }

  /// STT 엔진 사용 가능 여부 확인
  Future<bool> isSttAvailable() async {
    try {
      final result = await _sttChannel.invokeMethod<bool>('isAvailable');
      return result ?? false;
    } catch (_) {
      return false;
    }
  }
}

class _Scenario {
  final String greeting;
  final String question;
  final List<String> confirmKeywords;
  final Duration timeout;

  const _Scenario({
    required this.greeting,
    required this.question,
    required this.confirmKeywords,
    required this.timeout,
  });
}
