import 'package:flutter/material.dart';

/// Sanggam Orbit Global Frame.
/// Handles the top slim header image globally to avoid redundancy in cards.
class SanggamOrbitFrame extends StatelessWidget {
  final Widget child;

  const SanggamOrbitFrame({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Content
        Positioned.fill(child: child),

        // Global Top Header Frame (v5)
        Positioned(
          top: 0,
          left: 0,
          right: 0,
          height: 84,
          child: IgnorePointer(
            child: Stack(
              children: [
                const Positioned.fill(
                  child: DecoratedBox(
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        colors: [
                          Color(0xFF101A2F),
                          Color(0xAA0B1021),
                          Colors.transparent,
                        ],
                      ),
                    ),
                  ),
                ),
                Positioned.fill(
                  child: Opacity(
                    opacity: 0.32,
                    child: ColorFiltered(
                      colorFilter: const ColorFilter.mode(
                        Color(0xFF1A2A44),
                        BlendMode.multiply,
                      ),
                      child: Image.asset(
                        'assets/images/header_3d_frame_slim.png',
                        fit: BoxFit.cover,
                        alignment: Alignment.center,
                        filterQuality: FilterQuality.high,
                        errorBuilder: (context, error, stackTrace) =>
                            const SizedBox.shrink(),
                      ),
                    ),
                  ),
                ),
                Positioned(
                  left: 0,
                  right: 0,
                  bottom: 0,
                  height: 1,
                  child: Container(
                    decoration: const BoxDecoration(
                      gradient: LinearGradient(
                        colors: [
                          Colors.transparent,
                          Color(0x66D4AF37),
                          Color(0x99D4AF37),
                          Color(0x66D4AF37),
                          Colors.transparent,
                        ],
                        stops: [0.0, 0.22, 0.5, 0.78, 1.0],
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
