import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// FoodAnalysisScreen — Sanggam Orbit 음식 칼로리 분석
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 5x 제거
// [Rule 4] theme.textTheme.* ~12x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~3x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold ~8x → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Card 4x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] padding:12→16
// ───────────────────────────────────────────────────

/// 음식 칼로리 분석 화면
///
/// 사진 선택 → AI 분석 → 영양소 결과 표시
/// REST API 연동: 이미지 업로드 → AI 서버 분석 → 결과 수신
/// 카메라/갤러리 미연결 시 시뮬레이션 모드로 폴백
class FoodAnalysisScreen extends ConsumerStatefulWidget {
  const FoodAnalysisScreen({super.key});

  @override
  ConsumerState<FoodAnalysisScreen> createState() =>
      _FoodAnalysisScreenState();
}

class _FoodAnalysisScreenState
    extends ConsumerState<FoodAnalysisScreen> {
  _AnalysisState _state = _AnalysisState.initial;
  _FoodResult? _result;
  String? _imagePath;

  /// image_picker 연동 지점 (패키지 설치 후 주석 해제)
  Future<void> _pickImage(
      {bool fromCamera = true}) async {
    _imagePath = fromCamera
        ? 'camera_simulated.jpg'
        : 'gallery_simulated.jpg';
  }

  Future<void> _startAnalysis(
      {bool fromCamera = true}) async {
    setState(() => _state = _AnalysisState.analyzing);

    await _pickImage(fromCamera: fromCamera);

    try {
      final client = ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';
      final res = await client.analyzeFoodImage(
        userId: userId,
        imagePath: _imagePath ?? '',
      );
      if (mounted) {
        setState(() {
          _state = _AnalysisState.result;
          _result = _FoodResult(
            name: res['food_name'] as String? ??
                '인식 실패',
            calories:
                (res['calories'] as num?)?.toInt() ??
                    0,
            carbs:
                (res['carbs'] as num?)?.toInt() ?? 0,
            protein:
                (res['protein'] as num?)?.toInt() ??
                    0,
            fat:
                (res['fat'] as num?)?.toInt() ?? 0,
            confidence: (res['confidence'] as num?)
                    ?.toDouble() ??
                0.0,
          );
        });
        return;
      }
    } catch (_) {
      // API 미연결 → 시뮬레이션 폴백
    }

    await Future.delayed(
        const Duration(seconds: 2));
    if (mounted) {
      setState(() {
        _state = _AnalysisState.result;
        _result = _simulateResult();
      });
    }
  }

  _FoodResult _simulateResult() {
    final rng = Random();
    final foods = [
      ('비빔밥', 550, 75, 20, 12),
      ('김치찌개', 320, 15, 22, 18),
      ('불고기', 420, 10, 35, 28),
      ('된장찌개', 280, 12, 18, 15),
      ('제육볶음', 480, 20, 30, 25),
    ];
    final food = foods[rng.nextInt(foods.length)];
    return _FoodResult(
      name: food.$1,
      calories: food.$2 + rng.nextInt(50) - 25,
      carbs: food.$3 + rng.nextInt(10) - 5,
      protein: food.$4 + rng.nextInt(5) - 2,
      fat: food.$5 + rng.nextInt(5) - 2,
      confidence: 0.75 + rng.nextDouble() * 0.2,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '음식 칼로리 분석',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: switch (_state) {
                _AnalysisState.initial =>
                  _buildInitial(),
                _AnalysisState.analyzing =>
                  _buildAnalyzing(),
                _AnalysisState.result =>
                  _buildResult(),
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInitial() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                color: SanggamTheme.primary
                    .withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.camera_alt,
                  size: 48,
                  color: SanggamTheme.primary),
            ),
            const SizedBox(height: 24),
            const Text(
              '음식 사진을 촬영하면\nAI가 칼로리를 분석합니다.',
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 32),
            Row(
              mainAxisAlignment:
                  MainAxisAlignment.center,
              children: [
                FilledButton.icon(
                  onPressed: () => _startAnalysis(
                      fromCamera: true),
                  icon:
                      const Icon(Icons.camera_alt),
                  label: const Text('카메라'),
                  style: FilledButton.styleFrom(
                    backgroundColor:
                        SanggamTheme.primary,
                    foregroundColor:
                        SanggamTheme.background,
                    shape: RoundedRectangleBorder(
                      borderRadius:
                          BorderRadius.circular(16),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                OutlinedButton.icon(
                  onPressed: () => _startAnalysis(
                      fromCamera: false),
                  icon: const Icon(
                      Icons.photo_library),
                  label: const Text('갤러리'),
                ),
              ],
            ),
            const SizedBox(height: 16),
            const Text(
              '* 이미지 선택 후 AI 서버로 분석합니다.\n  서버 미연결 시 시뮬레이션 모드로 동작합니다.',
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAnalyzing() {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const SizedBox(
            width: 80,
            height: 80,
            child: CircularProgressIndicator(
              strokeWidth: 4,
              color: SanggamTheme.primary,
            ),
          ),
          const SizedBox(height: 24),
          const Text(
            'AI가 음식을 분석하고 있습니다...',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'YOLO 객체 탐지 → 영양소 DB 매칭',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildResult() {
    final r = _result!;
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.stretch,
        children: [
          // 음식 인식 결과
          SanggamContainer(
            borderRadius: 16,
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                const Icon(Icons.restaurant,
                    size: 48,
                    color: SanggamTheme.primary),
                const SizedBox(height: 8),
                Text(
                  r.name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  '인식 정확도: ${(r.confidence * 100).toInt()}%',
                  style: const TextStyle(
                    color:
                        SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // 칼로리 큰 숫자
          Center(
            child: Column(
              children: [
                Text(
                  '${r.calories}',
                  style: const TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 40,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Text(
                  'kcal',
                  style: TextStyle(
                    color:
                        SanggamTheme.onSurfaceDim,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // 영양소 카드 3개
          Row(
            children: [
              Expanded(
                  child: _nutrientCard(
                      '탄수화물',
                      '${r.carbs}g',
                      SanggamTheme.jagaeCyan)),
              const SizedBox(width: 8),
              Expanded(
                  child: _nutrientCard(
                      '단백질',
                      '${r.protein}g',
                      SanggamTheme.error)),
              const SizedBox(width: 8),
              Expanded(
                  child: _nutrientCard(
                      '지방',
                      '${r.fat}g',
                      SanggamTheme.primary)),
            ],
          ),
          const SizedBox(height: 24),

          // AI 코멘트
          SanggamContainer(
            borderRadius: 16,
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,
              children: [
                const Row(
                  children: [
                    Icon(Icons.smart_toy,
                        size: 20,
                        color: SanggamTheme.primary),
                    SizedBox(width: 8),
                    Text(
                      'AI 코멘트',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Text(
                  _generateComment(r),
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                    height: 1.6,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // 다시 분석
          OutlinedButton.icon(
            onPressed: () => setState(() =>
                _state = _AnalysisState.initial),
            icon: const Icon(Icons.refresh),
            label: const Text('다시 분석하기'),
          ),
        ],
      ),
    );
  }

  Widget _nutrientCard(
      String label, String value, Color color) {
    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          Text(
            value,
            style: TextStyle(
              color: color,
              fontSize: 22,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  String _generateComment(_FoodResult r) {
    if (r.calories > 500) {
      return '${r.name}은(는) 칼로리가 높은 편입니다. 한 끼 권장량(600~700kcal)에 가까우므로 반찬 양을 조절하시면 좋겠습니다. 단백질 ${r.protein}g으로 적절한 수준입니다.';
    }
    return '${r.name}은(는) 균형 잡힌 한 끼 식사입니다. 칼로리 ${r.calories}kcal로 적정 범위이며, 단백질/탄수화물/지방 비율이 양호합니다.';
  }
}

enum _AnalysisState { initial, analyzing, result }

class _FoodResult {
  final String name;
  final int calories, carbs, protein, fat;
  final double confidence;
  const _FoodResult({
    required this.name,
    required this.calories,
    required this.carbs,
    required this.protein,
    required this.fat,
    required this.confidence,
  });
}
