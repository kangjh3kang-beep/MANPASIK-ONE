import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';

// ───────────────────────────────────────────────────
// WaveAnalysisPanel — Sanggam Orbit 파형 분석 패널
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Theme.of(context).textTheme → 직접 TextStyle
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Colors.white10 → SanggamTheme.surfaceVariant
// ───────────────────────────────────────────────────

class WaveAnalysisPanel
    extends StatelessWidget {
  final String title;
  final Widget child;
  final Widget? footer;
  final double? width;
  final double? height;
  final bool isActive;
  final EdgeInsets? contentPadding;

  const WaveAnalysisPanel({
    super.key,
    required this.title,
    required this.child,
    this.footer,
    this.width,
    this.height,
    this.isActive = false,
    this.contentPadding,
  });

  @override
  Widget build(BuildContext context) {
    return OrnateGoldFrame(
      width: width,
      height: height,
      isActive: isActive,
      padding: contentPadding ??
          const EdgeInsets.fromLTRB(
              20, 24, 20, 20),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          // Header Center Aligned
          Center(
            child: Column(
              children: [
                Text(
                  title.toUpperCase(),
                  style: const TextStyle(
                    color:
                        SanggamTheme.primary,
                    fontSize: 10,
                    letterSpacing: 1.5,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 4),
                Container(
                  width: 30,
                  height: 1,
                  color: SanggamTheme.primary
                      .withValues(alpha: 0.5),
                ),
              ],
            ),
          ),
          const SizedBox(height: 16),

          // Content
          Expanded(child: child),

          // Footer
          if (footer != null) ...[
            const SizedBox(height: 12),
            Container(
              height: 1,
              color: SanggamTheme
                  .surfaceVariant,
            ),
            const SizedBox(height: 8),
            footer!,
          ],
        ],
      ),
    );
  }
}
