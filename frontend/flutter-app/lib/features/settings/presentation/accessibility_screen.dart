import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

/// 접근성 설정 화면
class AccessibilityScreen extends ConsumerStatefulWidget {
  const AccessibilityScreen({super.key});

  @override
  ConsumerState<AccessibilityScreen> createState() => _AccessibilityScreenState();
}

class _AccessibilityScreenState extends ConsumerState<AccessibilityScreen> {
  double _fontScale = 1.0;
  bool _highContrast = false;
  bool _screenReader = false;
  double _ttsSpeed = 1.0;
  String _ttsLanguage = 'ko';

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('접근성'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: ListView(
        children: [
          // 글꼴 크기
          _buildSectionHeader(theme, '텍스트'),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              children: [
                const Text('가', style: TextStyle(fontSize: 12)),
                Expanded(
                  child: Slider(
                    value: _fontScale,
                    min: 0.8,
                    max: 1.6,
                    divisions: 8,
                    label: '${(_fontScale * 100).toInt()}%',
                    onChanged: (v) => setState(() => _fontScale = v),
                  ),
                ),
                const Text('가', style: TextStyle(fontSize: 24)),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Text(
                  '미리보기: 혈당 수치가 정상 범위입니다.',
                  style: TextStyle(fontSize: 14 * _fontScale),
                ),
              ),
            ),
          ),
          const Divider(),

          // 시각 보조
          _buildSectionHeader(theme, '시각 보조'),
          SwitchListTile(
            secondary: const Icon(Icons.contrast),
            title: const Text('고대비 모드'),
            subtitle: const Text('화면 대비를 높여 가독성 향상'),
            value: _highContrast,
            onChanged: (v) => setState(() => _highContrast = v),
          ),
          SwitchListTile(
            secondary: const Icon(Icons.record_voice_over),
            title: const Text('화면 읽기'),
            subtitle: const Text('화면의 텍스트를 음성으로 읽어줍니다'),
            value: _screenReader,
            onChanged: (v) => setState(() => _screenReader = v),
          ),
          const Divider(),

          // 음성 설정
          _buildSectionHeader(theme, '음성 설정 (TTS)'),
          ListTile(
            leading: const Icon(Icons.speed),
            title: const Text('읽기 속도'),
            subtitle: Text('${_ttsSpeed}x'),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Slider(
              value: _ttsSpeed,
              min: 0.5,
              max: 2.0,
              divisions: 6,
              label: '${_ttsSpeed}x',
              onChanged: (v) => setState(() => _ttsSpeed = v),
            ),
          ),
          ListTile(
            leading: const Icon(Icons.language),
            title: const Text('TTS 언어'),
            subtitle: Text(_ttsLanguageName(_ttsLanguage)),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showTtsLanguageDialog(context),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(ThemeData theme, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: theme.textTheme.labelLarge?.copyWith(
          color: theme.colorScheme.primary,
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

  void _showTtsLanguageDialog(BuildContext context) {
    final languages = ['ko', 'en', 'ja', 'zh'];
    showDialog(
      context: context,
      builder: (ctx) => SimpleDialog(
        title: const Text('TTS 언어 선택'),
        children: languages.map((lang) {
          final isSelected = lang == _ttsLanguage;
          return ListTile(
            title: Text(_ttsLanguageName(lang)),
            trailing: isSelected ? const Icon(Icons.check, color: Colors.green) : null,
            onTap: () {
              setState(() => _ttsLanguage = lang);
              Navigator.pop(ctx);
            },
          );
        }).toList(),
      ),
    );
  }
}
