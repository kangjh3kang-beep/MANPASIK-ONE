import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// BodyScanHud — Sanggam Orbit 바디스캔 HUD
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// v6.0 상단 HUD 오버레이 — "바디스캔 활성" + 진행률 바.
class BodyScanHud extends StatelessWidget {
  final double scanProgress;

  const BodyScanHud(
      {super.key,
      required this.scanProgress});

  @override
  Widget build(BuildContext context) {
    final percent =
        (scanProgress * 100).toInt();
    return Padding(
      padding:
          const EdgeInsets.only(top: 10),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // "바디스캔 활성" 텍스트
          Text(
            '바디스캔 활성',
            style: TextStyle(
              color: SanggamTheme.primary
                  .withValues(alpha: 0.7),
              fontSize: 12,
              fontWeight: FontWeight.w600,
              letterSpacing: 3.0,
            ),
          ),
          const SizedBox(height: 6),
          // 진행률 바 + 퍼센트
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              // 바 배경 + 진행
              SizedBox(
                width: 140,
                height: 4,
                child: ClipRRect(
                  borderRadius:
                      BorderRadius.circular(
                          2),
                  child: Stack(
                    children: [
                      // 배경
                      Container(
                        color: SanggamTheme
                            .jagaeCyan
                            .withValues(
                                alpha:
                                    0.15),
                      ),
                      // 진행
                      FractionallySizedBox(
                        widthFactor:
                            scanProgress
                                .clamp(
                                    0.0,
                                    1.0),
                        child: Container(
                          decoration:
                              BoxDecoration(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        2),
                            color: SanggamTheme
                                .primary
                                .withValues(
                                    alpha:
                                        0.7),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(width: 8),
              // 퍼센트 표시
              Text(
                '$percent%',
                style: TextStyle(
                  color: SanggamTheme
                      .primary
                      .withValues(
                          alpha: 0.6),
                  fontSize: 10,
                  fontWeight:
                      FontWeight.w600,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
