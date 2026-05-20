import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// PrimaryButton — Sanggam Gold 그래디언트 버튼
//
// [Rule 4] theme.colorScheme.primary → SanggamTheme.primary
// [Rule 4] Colors.white 하드코딩 → SanggamTheme.background (금색 위 다크 텍스트)
// [Rule 4] FilledButton 기본 스타일 → SanggamTheme.goldGradient 적용
// [Rule 2] height 56 (8×7), borderRadius 16 (8×2), SizedBox(width:8) — 준수
//
// 비활성 상태: AnimatedOpacity(0.5)로 시각적 피드백.
// ───────────────────────────────────────────────────

class PrimaryButton extends StatelessWidget {
  final String text;
  final VoidCallback? onPressed;
  final bool isLoading;
  final IconData? icon;

  const PrimaryButton({
    super.key,
    required this.text,
    this.onPressed,
    this.isLoading = false,
    this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final enabled = !isLoading && onPressed != null;

    return GestureDetector(
      onTap: enabled ? onPressed : null,
      child: AnimatedOpacity(
        opacity: enabled ? 1.0 : 0.5,
        duration: const Duration(milliseconds: 200),
        child: Container(
          height: 56,
          width: double.infinity,
          decoration: BoxDecoration(
            gradient: SanggamTheme.goldGradient,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(
                color: SanggamTheme.primary.withValues(alpha: 0.3),
                blurRadius: 16,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: Center(
            child: isLoading
                ? const SizedBox(
                    height: 24,
                    width: 24,
                    child: CircularProgressIndicator(
                      strokeWidth: 2.5,
                      color: SanggamTheme.background,
                    ),
                  )
                : Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      if (icon != null) ...[
                        Icon(icon,
                            size: 20, color: SanggamTheme.background),
                        const SizedBox(width: 8),
                      ],
                      Text(
                        text,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: SanggamTheme.background,
                        ),
                      ),
                    ],
                  ),
          ),
        ),
      ),
    );
  }
}
