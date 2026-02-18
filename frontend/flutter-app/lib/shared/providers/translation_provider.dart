import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';

/// 실시간 번역 상태 (C6)
class TranslationState {
  final Map<String, String> cache;
  final bool isTranslating;
  final String? error;

  const TranslationState({
    this.cache = const {},
    this.isTranslating = false,
    this.error,
  });

  TranslationState copyWith({
    Map<String, String>? cache,
    bool? isTranslating,
    String? error,
  }) {
    return TranslationState(
      cache: cache ?? this.cache,
      isTranslating: isTranslating ?? this.isTranslating,
      error: error,
    );
  }
}

/// 실시간 번역 Notifier (C6)
///
/// REST Gateway의 translateRealtime 엔드포인트를 통해
/// 실시간 번역을 제공하며, 캐시로 중복 요청을 방지합니다.
class TranslationNotifier extends StateNotifier<TranslationState> {
  TranslationNotifier(this._restClient) : super(const TranslationState());

  final dynamic _restClient; // ManPaSikRestClient

  /// 텍스트 번역 (캐시 우선)
  Future<String> translate({
    required String text,
    required String targetLanguage,
    String sourceLanguage = 'auto',
  }) async {
    final cacheKey = '$sourceLanguage:$targetLanguage:$text';
    if (state.cache.containsKey(cacheKey)) {
      return state.cache[cacheKey]!;
    }

    state = state.copyWith(isTranslating: true, error: null);

    try {
      final result = await _restClient.translateRealtime(
        text: text,
        sourceLanguage: sourceLanguage,
        targetLanguage: targetLanguage,
      );
      final translated = result['translated_text'] as String? ??
          result['text'] as String? ??
          text;

      final newCache = Map<String, String>.from(state.cache);
      newCache[cacheKey] = translated;

      state = state.copyWith(cache: newCache, isTranslating: false);
      return translated;
    } catch (e) {
      debugPrint('[TranslationNotifier] 번역 실패: $e');
      state = state.copyWith(isTranslating: false, error: e.toString());
      return text; // 원문 반환
    }
  }

  void clearCache() {
    state = const TranslationState();
  }
}

/// 번역 Provider
final translationProvider =
    StateNotifierProvider<TranslationNotifier, TranslationState>((ref) {
  final client = ref.watch(restClientProvider);
  return TranslationNotifier(client);
});
