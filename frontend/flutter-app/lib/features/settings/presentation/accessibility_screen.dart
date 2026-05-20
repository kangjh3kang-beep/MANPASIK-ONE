import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AccessibilityScreen — Sanggam Orbit 접근성 설정
//
// [Rule 4] +sanggam_theme.dart + sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData 파라미터 1x 제거
// [Rule 4] theme.textTheme/colorScheme → 직접 TextStyle
// [Rule 4] Card → SanggamContainer
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] SimpleDialog → AlertDialog 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 4] Divider color → SanggamTheme.surfaceVariant
// ───────────────────────────────────────────────────

/// 접근성 설정 화면
class AccessibilityScreen
    extends ConsumerStatefulWidget {
  const AccessibilityScreen({super.key});

  @override
  ConsumerState<AccessibilityScreen>
      createState() =>
          _AccessibilityScreenState();
}

class _AccessibilityScreenState
    extends ConsumerState<AccessibilityScreen> {
  double _fontScale = 1.0;
  bool _highContrast = false;
  bool _screenReader = false;
  double _ttsSpeed = 1.0;
  String _ttsLanguage = 'ko';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '접근성',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: ListView(
                children: [
                  // 글꼴 크기
                  _buildSectionHeader('텍스트'),
                  Padding(
                    padding: const EdgeInsets
                        .symmetric(
                        horizontal: 16),
                    child: Row(
                      children: [
                        const Text('가',
                            style: TextStyle(
                                fontSize: 12,
                                color: Colors
                                    .white)),
                        Expanded(
                          child: Slider(
                            value: _fontScale,
                            min: 0.8,
                            max: 1.6,
                            divisions: 8,
                            activeColor:
                                SanggamTheme
                                    .primary,
                            inactiveColor:
                                SanggamTheme
                                    .surfaceVariant,
                            label:
                                '${(_fontScale * 100).toInt()}%',
                            onChanged: (v) =>
                                setState(() =>
                                    _fontScale =
                                        v),
                          ),
                        ),
                        const Text('가',
                            style: TextStyle(
                                fontSize: 24,
                                color: Colors
                                    .white)),
                      ],
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets
                        .symmetric(
                        horizontal: 16),
                    child: SanggamContainer(
                      borderRadius: 16,
                      padding:
                          const EdgeInsets.all(
                              16),
                      child: Text(
                        '미리보기: 혈당 수치가 정상 범위입니다.',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize:
                              14 * _fontScale,
                        ),
                      ),
                    ),
                  ),
                  const Divider(
                      color: SanggamTheme
                          .surfaceVariant),

                  // 시각 보조
                  _buildSectionHeader('시각 보조'),
                  SwitchListTile(
                    secondary: const Icon(
                        Icons.contrast,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '고대비 모드',
                        style: TextStyle(
                            color:
                                Colors.white)),
                    subtitle: const Text(
                        '화면 대비를 높여 가독성 향상',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    value: _highContrast,
                    activeColor:
                        SanggamTheme.primary,
                    onChanged: (v) => setState(
                        () =>
                            _highContrast = v),
                  ),
                  SwitchListTile(
                    secondary: const Icon(
                        Icons
                            .record_voice_over,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text('화면 읽기',
                        style: TextStyle(
                            color:
                                Colors.white)),
                    subtitle: const Text(
                        '화면의 텍스트를 음성으로 읽어줍니다',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    value: _screenReader,
                    activeColor:
                        SanggamTheme.primary,
                    onChanged: (v) => setState(
                        () =>
                            _screenReader = v),
                  ),
                  const Divider(
                      color: SanggamTheme
                          .surfaceVariant),

                  // 음성 설정
                  _buildSectionHeader(
                      '음성 설정 (TTS)'),
                  ListTile(
                    leading: const Icon(
                        Icons.speed,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text('읽기 속도',
                        style: TextStyle(
                            color:
                                Colors.white)),
                    subtitle: Text(
                        '${_ttsSpeed}x',
                        style: const TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                  ),
                  Padding(
                    padding: const EdgeInsets
                        .symmetric(
                        horizontal: 16),
                    child: Slider(
                      value: _ttsSpeed,
                      min: 0.5,
                      max: 2.0,
                      divisions: 6,
                      activeColor:
                          SanggamTheme.primary,
                      inactiveColor: SanggamTheme
                          .surfaceVariant,
                      label: '${_ttsSpeed}x',
                      onChanged: (v) =>
                          setState(() =>
                              _ttsSpeed = v),
                    ),
                  ),
                  ListTile(
                    leading: const Icon(
                        Icons.language,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        'TTS 언어',
                        style: TextStyle(
                            color:
                                Colors.white)),
                    subtitle: Text(
                        _ttsLanguageName(
                            _ttsLanguage),
                        style: const TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () =>
                        _showTtsLanguageDialog(
                            context),
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(
          16, 16, 16, 8),
      child: Text(
        title,
        style: const TextStyle(
          color: SanggamTheme.primary,
          fontSize: 12,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  String _ttsLanguageName(String code) {
    return switch (code) {
      'ko' => '한국어',
      'en' => 'English',
      'ja' => '日本語',
      'zh' => '中文',
      _ => code,
    };
  }

  void _showTtsLanguageDialog(
      BuildContext context) {
    final languages = ['ko', 'en', 'ja', 'zh'];
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text('TTS 언어 선택',
            style: TextStyle(
                color: Colors.white)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: languages.map((lang) {
            final isSelected =
                lang == _ttsLanguage;
            return ListTile(
              title: Text(
                  _ttsLanguageName(lang),
                  style: const TextStyle(
                      color: Colors.white)),
              trailing: isSelected
                  ? const Icon(Icons.check,
                      color: SanggamTheme
                          .jagaeCyan)
                  : null,
              onTap: () {
                setState(() =>
                    _ttsLanguage = lang);
                Navigator.pop(ctx);
              },
            );
          }).toList(),
        ),
      ),
    );
  }
}
