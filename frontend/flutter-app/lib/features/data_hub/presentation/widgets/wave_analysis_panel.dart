import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';

class WaveAnalysisPanel extends StatelessWidget {
  final String title;
  final Widget child;
  final Widget? footer;
  final double? width;
  final double? height;
  final bool isActive;

  const WaveAnalysisPanel({
    super.key,
    required this.title,
    required this.child,
    this.footer,
    this.width,
    this.height,
    this.isActive = false,
  });

  @override
  Widget build(BuildContext context) {
    return OrnateGoldFrame(
      width: width,
      height: height,
      isActive: isActive,
      padding: const EdgeInsets.fromLTRB(20, 24, 20, 20), // Adjust for frame
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header Center Aligned per design reference
          Center(
            child: Column(
              children: [
                Text(
                  title.toUpperCase(),
                  style: Theme.of(context).textTheme.labelSmall?.copyWith(
                    color: AppTheme.sanggamGold, 
                    fontSize: 10,
                    letterSpacing: 1.5,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 4),
                Container(
                  width: 30,
                  height: 1,
                  color: AppTheme.sanggamGold.withOpacity(0.5),
                )
              ],
            ),
          ),
          const SizedBox(height: 16),
          
          // Content
          Expanded(child: child),
          
          // Footer
          if (footer != null) ...[
            const SizedBox(height: 12),
            // Decorative Footer Line (Dashed or styled)
            Container(height: 1, color: Colors.white10), 
            const SizedBox(height: 8),
            footer!,
          ],
        ],
      ),
    );
  }
}
