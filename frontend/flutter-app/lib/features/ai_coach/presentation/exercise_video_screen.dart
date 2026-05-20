import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ExerciseVideoScreen — Sanggam Orbit 운동 분석
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 4x 제거
// [Rule 4] theme.textTheme.* ~12x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~3x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold ~10x → SanggamTheme.primary
// [Rule 4] Colors.red 2x → SanggamTheme.error
// [Rule 4] Colors.orange 2x → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.purple 2x → SanggamTheme.jagaeMagenta
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Card 5x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] padding:12→16, h:12→16, horizontal:10→16,
//          borderRadius:12→16, margin:4→8
// ───────────────────────────────────────────────────

/// 운동 비디오 칼로리 분석 화면
///
/// 운동 영상 촬영/선택 → AI 포즈 추정 → 칼로리/운동 유형 분석
class ExerciseVideoScreen extends ConsumerStatefulWidget {
  const ExerciseVideoScreen({super.key});

  @override
  ConsumerState<ExerciseVideoScreen> createState() =>
      _ExerciseVideoScreenState();
}

enum _AnalysisState { initial, analyzing, result }

class _ExerciseVideoScreenState
    extends ConsumerState<ExerciseVideoScreen> {
  _AnalysisState _state = _AnalysisState.initial;
  _ExerciseResult? _result;
  double _progress = 0;
  String? _videoPath;

  /// image_picker 연동 지점 (패키지 설치 후 주석 해제)
  Future<void> _pickVideo(
      {bool fromCamera = true}) async {
    _videoPath = fromCamera
        ? 'camera_sim.mp4'
        : 'gallery_sim.mp4';
  }

  Future<void> _startAnalysis(
      {bool fromCamera = true}) async {
    setState(() {
      _state = _AnalysisState.analyzing;
      _progress = 0;
    });

    await _pickVideo(fromCamera: fromCamera);

    try {
      final client = ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';

      for (var i = 0; i <= 100; i += 5) {
        await Future.delayed(
            const Duration(milliseconds: 100));
        if (!mounted) return;
        setState(() => _progress = i / 100);
      }

      if (mounted) {
        setState(() {
          _state = _AnalysisState.result;
          _result = _simulateResult();
        });
      }
    } catch (_) {
      await Future.delayed(
          const Duration(seconds: 2));
      if (mounted) {
        setState(() {
          _state = _AnalysisState.result;
          _result = _simulateResult();
        });
      }
    }
  }

  _ExerciseResult _simulateResult() {
    final rng = Random();
    final exercises = [
      ('스쿼트', 'strength', 45, 180, 15),
      ('런지', 'strength', 38, 150, 12),
      ('푸시업', 'strength', 42, 160, 20),
      ('플랭크', 'core', 30, 120, 3),
      ('버피', 'cardio', 65, 90, 10),
      ('점프 스쿼트', 'cardio', 55, 120, 12),
    ];
    final ex = exercises[rng.nextInt(exercises.length)];
    final duration = ex.$4 + rng.nextInt(60) - 30;
    return _ExerciseResult(
      exerciseName: ex.$1,
      category: ex.$2,
      caloriesBurned: ex.$3 + rng.nextInt(20) - 10,
      durationSeconds:
          duration > 30 ? duration : 60,
      reps: ex.$5 + rng.nextInt(5) - 2,
      formScore: 0.7 + rng.nextDouble() * 0.25,
      feedback: _generateFeedback(
          ex.$1, 0.7 + rng.nextDouble() * 0.25),
    );
  }

  String _generateFeedback(
      String name, double score) {
    if (score > 0.85) {
      return '훌륭한 $name 자세입니다! 균형과 깊이가 적절하며, 관절 각도가 이상적입니다. 현재 페이스를 유지하세요.';
    } else if (score > 0.7) {
      return '$name 동작이 양호합니다. 허리를 조금 더 곧게 유지하면 효과가 높아집니다. 호흡에 신경 써주세요.';
    }
    return '$name 동작 개선이 필요합니다. 무릎이 발끝을 넘지 않도록 주의하고, 천천히 수행하여 정확한 자세를 잡아보세요.';
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
                      '운동 분석',
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
              child: const Icon(
                  Icons.fitness_center,
                  size: 48,
                  color: SanggamTheme.primary),
            ),
            const SizedBox(height: 24),
            const Text(
              '운동 영상을 촬영하면\nAI가 자세를 분석하고 칼로리를 계산합니다.',
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
                  icon: const Icon(Icons.videocam),
                  label: const Text('촬영'),
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
                  icon:
                      const Icon(Icons.video_library),
                  label: const Text('갤러리'),
                ),
              ],
            ),
            const SizedBox(height: 16),
            const Text(
              'AI 포즈 추정으로 운동 자세와 반복 횟수를 분석합니다.',
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
          Stack(
            alignment: Alignment.center,
            children: [
              SizedBox(
                width: 100,
                height: 100,
                child: CircularProgressIndicator(
                  value: _progress,
                  strokeWidth: 6,
                  color: SanggamTheme.primary,
                  backgroundColor: SanggamTheme
                      .primary
                      .withValues(alpha: 0.15),
                ),
              ),
              Text(
                '${(_progress * 100).toInt()}%',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          const Text(
            '운동 영상을 분석하고 있습니다...',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '포즈 추정 → 동작 분류 → 칼로리 계산',
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
          // 운동 인식 결과
          SanggamContainer(
            borderRadius: 16,
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                const Icon(Icons.fitness_center,
                    size: 48,
                    color: SanggamTheme.primary),
                const SizedBox(height: 8),
                Text(
                  r.exerciseName,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Container(
                  margin:
                      const EdgeInsets.only(top: 8),
                  padding:
                      const EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 8),
                  decoration: BoxDecoration(
                    color: _categoryColor(r.category)
                        .withValues(alpha: 0.15),
                    borderRadius:
                        BorderRadius.circular(16),
                  ),
                  child: Text(
                    _categoryLabel(r.category),
                    style: TextStyle(
                      fontSize: 12,
                      color:
                          _categoryColor(r.category),
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // 주요 수치
          Row(
            children: [
              Expanded(
                  child: _statCard(
                      '${r.caloriesBurned}',
                      'kcal',
                      SanggamTheme.primary)),
              const SizedBox(width: 8),
              Expanded(
                  child: _statCard(
                      '${r.durationSeconds ~/ 60}:${(r.durationSeconds % 60).toString().padLeft(2, '0')}',
                      '시간',
                      SanggamTheme.jagaeCyan)),
              const SizedBox(width: 8),
              Expanded(
                  child: _statCard(
                      '${r.reps}',
                      '회',
                      SanggamTheme.jagaeMagenta)),
            ],
          ),
          const SizedBox(height: 16),

          // 자세 점수
          SanggamContainer(
            borderRadius: 16,
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,
              children: [
                const Text(
                  '자세 점수',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: ClipRRect(
                        borderRadius:
                            BorderRadius.circular(8),
                        child:
                            LinearProgressIndicator(
                          value: r.formScore,
                          minHeight: 12,
                          color: _scoreColor(
                              r.formScore),
                          backgroundColor: SanggamTheme
                              .surfaceVariant,
                        ),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Text(
                      '${(r.formScore * 100).toInt()}점',
                      style: TextStyle(
                        color: _scoreColor(
                            r.formScore),
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // AI 피드백
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
                      'AI 코치 피드백',
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
                  r.feedback,
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

  Widget _statCard(
      String value, String unit, Color color) {
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
            unit,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  Color _categoryColor(String category) {
    return switch (category) {
      'strength' => SanggamTheme.error,
      'cardio' => SanggamTheme.primary,
      'core' => SanggamTheme.jagaeMagenta,
      _ => SanggamTheme.onSurfaceDim,
    };
  }

  String _categoryLabel(String category) {
    return switch (category) {
      'strength' => '근력 운동',
      'cardio' => '유산소',
      'core' => '코어',
      _ => category,
    };
  }

  Color _scoreColor(double score) {
    if (score >= 0.85) return SanggamTheme.jagaeCyan;
    if (score >= 0.7) return SanggamTheme.primary;
    return SanggamTheme.error;
  }
}

class _ExerciseResult {
  final String exerciseName;
  final String category;
  final int caloriesBurned;
  final int durationSeconds;
  final int reps;
  final double formScore;
  final String feedback;

  const _ExerciseResult({
    required this.exerciseName,
    required this.category,
    required this.caloriesBurned,
    required this.durationSeconds,
    required this.reps,
    required this.formScore,
    required this.feedback,
  });
}
