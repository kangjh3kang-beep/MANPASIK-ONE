import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// v6.0 상단 HUD 오버레이 — "BODY SCAN ACTIVE" + 진행률 바.
///
/// CustomPaint 내부 8px/60×3px 대비 → 12px/140×4px Flutter 위젯.
class BodyScanHud extends StatelessWidget {
  final double scanProgress;

  const BodyScanHud({super.key, required this.scanProgress});

  @override
  Widget build(BuildContext context) {
    final percent = (scanProgress * 100).toInt();
    return Padding(
      padding: const EdgeInsets.only(top: 10),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // "BODY SCAN ACTIVE" 텍스트
          Text(
            'BODY SCAN ACTIVE',
            style: TextStyle(
              color: AppTheme.sanggamGold.withValues(alpha: 0.7),
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
                  borderRadius: BorderRadius.circular(2),
                  child: Stack(
                    children: [
                      // 배경
                      Container(
                        color: AppTheme.waveCyan.withValues(alpha: 0.15),
                      ),
                      // 진행
                      FractionallySizedBox(
                        widthFactor: scanProgress.clamp(0.0, 1.0),
                        child: Container(
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(2),
                            color: AppTheme.sanggamGold.withValues(alpha: 0.7),
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
                  color: AppTheme.sanggamGold.withValues(alpha: 0.6),
                  fontSize: 10,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
