import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/glass_morphism_card.dart';
import 'package:manpasik/shared/widgets/medical_stat_card.dart';
import 'package:manpasik/shared/widgets/jarvis_connector.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';

// ───────────────────────────────────────────────────
// OrbitalHud — Sanggam Orbit 궤도 HUD
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] Color(0xFF080C17) → SanggamTheme.background
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// MANPASIK Premium UI - The Sanggam Orbit HUD
///
/// A dynamic, systemic HUD that places health metrics in a polar coordinate system
/// orbiting a central anchor. Incorporates Korean traditional aesthetics and
/// futuristic HCI principles.
class OrbitalHud extends StatefulWidget {
  final Offset center;
  final List<OrbitalNodeData> nodes;
  final Widget centerWidget;
  final double baseRadius;

  const OrbitalHud({
    super.key,
    required this.center,
    required this.nodes,
    required this.centerWidget,
    this.baseRadius = 160.0,
  });

  @override
  State<OrbitalHud> createState() =>
      _OrbitalHudState();
}

class _OrbitalHudState extends State<OrbitalHud>
    with TickerProviderStateMixin {
  late AnimationController _pulseController;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      vsync: this,
      duration:
          const Duration(milliseconds: 1500),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _pulseController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // 1. Central Anchor (HoloBody)
        Positioned(
          left: widget.center.dx - 150,
          top: widget.center.dy - 250,
          width: 300,
          height: 500,
          child: widget.centerWidget,
        ),

        // 2. Orbital Nodes & Connectors
        ..._buildOrbitalLayers(),
      ],
    );
  }

  List<Widget> _buildOrbitalLayers() {
    final List<Widget> layers = [];

    for (int i = 0;
        i < widget.nodes.length;
        i++) {
      final node = widget.nodes[i];
      final angle =
          node.angle * (math.pi / 180);
      
      // [UI 최적화] 중앙 인체 홀로그램(너비 300)과 겹치지 않도록 최소 반경 대폭 확보
      // 어깨 너비와 패널 너비를 고려하여 더 여유 있는 공간 확보
      final baseRadius = widget.baseRadius;
      final xOffset = math.cos(angle).abs();
      double effectiveRadius = baseRadius * node.radiusMultiplier;
      
      // X축 오프셋이 큰 경우(대각선/가로) 반경을 340px 이상으로 확장 (기존 320px)
      if (xOffset > 0.4) {
        effectiveRadius = math.max(effectiveRadius, 340.0 / xOffset);
      } else {
        // 상/하단 노드들도 최소 280px 이상 거리 유지 (기존 220px)
        effectiveRadius = math.max(effectiveRadius, 280.0);
      }

      final nodePos = Offset(
        widget.center.dx +
            effectiveRadius * math.cos(angle) -
            65, // 패널 너비(130)의 절반으로 조정하여 정확한 중앙 정렬
        widget.center.dy +
            effectiveRadius * math.sin(angle) -
            50,
      );

      // Connectors with Pulse Effect
      layers.add(
        JarvisConnector(
          startOffset: widget.center +
              node.bodyAttachmentOffset,
          endOffset:
              nodePos + const Offset(60, 50),
          color: node.isAlert
              ? SanggamTheme.error
                  .withValues(alpha: 0.4)
              : SanggamTheme.primary
                  .withValues(alpha: 0.3),
        ),
      );

      // Orbital Pod (The Glass Card)
      layers.add(
        Positioned(
          left: nodePos.dx,
          top: nodePos.dy,
          width: 130,
          height: 100,
          child: AnimateFadeInUp(
            delay: Duration(
                milliseconds:
                    200 + (i * 150)),
            duration: const Duration(
                milliseconds: 800),
            offset: 40,
            child: _buildOrbitalPod(node),
          ),
        ),
      );
    }

    return layers;
  }

  Widget _buildOrbitalPod(
      OrbitalNodeData node) {
    return GlassmorphismCard(
      useKoreanBorder: true,
      opacity: node.isAlert ? 0.3 : 0.15,
      baseColor: node.isAlert
          ? SanggamTheme.error
              .withValues(alpha: 0.1)
          : SanggamTheme.background,
      padding: const EdgeInsets.all(8),
      child: MedicalStatCard(
        label: node.label,
        value: node.value,
        unit: node.unit,
        icon: node.icon,
        isAlert: node.isAlert,
        color: SanggamTheme.jagaeCyan,
      ),
    );
  }
}

/// Data structure for a node in the Orbital HUD
class OrbitalNodeData {
  final String label;
  final String value;
  final String unit;
  final IconData icon;
  final bool isAlert;
  final double angle;
  final double radiusMultiplier;
  final Offset bodyAttachmentOffset;

  const OrbitalNodeData({
    required this.label,
    required this.value,
    required this.unit,
    required this.icon,
    this.isAlert = false,
    required this.angle,
    this.radiusMultiplier = 1.0,
    this.bodyAttachmentOffset = Offset.zero,
  });
}
