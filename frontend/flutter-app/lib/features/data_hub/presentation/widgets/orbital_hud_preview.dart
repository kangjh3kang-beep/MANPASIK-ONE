import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/holo_body.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/orbital_hud.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';

// ───────────────────────────────────────────────────
// OrbitalHudPreview — Sanggam Orbit 디자인 시안
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// MANPASIK Design Preview — Sanggam Orbit System
class OrbitalHudPreview
    extends StatelessWidget {
  const OrbitalHudPreview({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: RoyalCloudBackground(
        child: Stack(
          children: [
            // Background Lattice
            const Positioned.fill(
              child: JagaeContainer(
                showLattice: true,
                opacity: 0.05,
                child: SizedBox.expand(),
              ),
            ),

            // Main Orbital HUD
            LayoutBuilder(
              builder:
                  (context, constraints) {
                final screenCenter =
                    Offset(
                        constraints
                                .maxWidth /
                            2,
                        constraints
                                .maxHeight *
                            0.45);

                return OrbitalHud(
                  center: screenCenter,
                  baseRadius:
                      constraints.maxWidth *
                          0.38,
                  nodes: const [
                    OrbitalNodeData(
                      label: '심박수',
                      value: '72',
                      unit: 'bpm',
                      icon: Icons.favorite,
                      angle: -135,
                      radiusMultiplier:
                          1.1,
                      bodyAttachmentOffset:
                          Offset(0, -100),
                    ),
                    OrbitalNodeData(
                      label: '산소포화도',
                      value: '98',
                      unit: '%',
                      icon: Icons.air,
                      angle: -45,
                      radiusMultiplier:
                          1.1,
                      bodyAttachmentOffset:
                          Offset(0, -80),
                    ),
                    OrbitalNodeData(
                      label: '혈당',
                      value: '105',
                      unit: 'mg/dL',
                      icon:
                          Icons.bloodtype,
                      angle: 135,
                      radiusMultiplier:
                          1.2,
                      bodyAttachmentOffset:
                          Offset(0, 50),
                      isAlert: true,
                    ),
                    OrbitalNodeData(
                      label: '뇌활동',
                      value: '42',
                      unit: '%',
                      icon:
                          Icons.psychology,
                      angle: 225,
                      radiusMultiplier:
                          1.0,
                      bodyAttachmentOffset:
                          Offset(0, -180),
                    ),
                  ],
                  centerWidget:
                      const HoloBody(
                    showEcg: true,
                    showHud: true,
                    color: SanggamTheme
                        .jagaeCyan,
                  ),
                );
              },
            ),

            // Header Typography
            Positioned(
              top: 60,
              left: 20,
              right: 20,
              child: Column(
                children: [
                  Text(
                    'MANPASIK V6.0',
                    style: TextStyle(
                      fontFamily:
                          'Orbitron',
                      color: SanggamTheme
                          .primary
                          .withValues(
                              alpha: 0.8),
                      fontSize: 12,
                      letterSpacing: 4,
                      fontWeight:
                          FontWeight.bold,
                    ),
                  ),
                  const SizedBox(
                      height: 8),
                  const Text(
                    '상감 오르빗 시스템 활성',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight:
                          FontWeight.w900,
                      letterSpacing: 1.2,
                    ),
                  ),
                  Container(
                    margin:
                        const EdgeInsets
                            .only(top: 12),
                    width: 40,
                    height: 1,
                    color: SanggamTheme
                        .primary,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
