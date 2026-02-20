import 'package:flutter/material.dart';
import 'dart:ui' as ui;
import 'dart:math' as math;
import 'package:manpasik/core/theme/app_theme.dart';

class OrnateGoldFrame extends StatefulWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsets padding;
  final bool isActive;

  const OrnateGoldFrame({
    super.key,
    required this.child,
     this.width,
     this.height,
    this.padding = const EdgeInsets.all(16),
    this.isActive = false,
  });

  @override
  State<OrnateGoldFrame> createState() => _OrnateGoldFrameState();
}

class _OrnateGoldFrameState extends State<OrnateGoldFrame> with TickerProviderStateMixin {
  late AnimationController _shimmerController;
  late AnimationController _breathingController;
  
  @override
  void initState() {
    super.initState();
    _shimmerController = AnimationController(
       vsync: this, 
       duration: const Duration(seconds: 12), // Slower, more elegant
    )..repeat(reverse: false);

    _breathingController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 5),
    )..repeat(reverse: true);
  }

  bool _isHovered = false;

  @override
  void dispose() {
    _shimmerController.dispose();
    _breathingController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedBuilder(
        animation: _breathingController,
        builder: (context, _) {
          // Hover Scale: 1.02x when hovered, plus breathing
          final baseScale = _isHovered ? 1.02 : 1.0;
          final scale = baseScale + (_breathingController.value * 0.005); 
          
          return Transform.scale(
            scale: scale,
            child: Container(
              width: widget.width,
              height: widget.height,
            child: AnimatedBuilder(
              animation: _shimmerController,
              builder: (context, _) {
                return CustomPaint(
                  painter: _OrnateFramePainter(
                    color: AppTheme.sanggamGold,
                    glowColor: widget.isActive ? AppTheme.waveCyan : AppTheme.sanggamGold,
                    shimmerValue: _shimmerController.value,
                  ),
                  child: Container(
                    padding: const EdgeInsets.all(4), // Frame spacing
                    child: ClipRRect( 
                      borderRadius: BorderRadius.circular(16),
                      child: BackdropFilter(
                        filter: ui.ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                        child: Container(
                           // Renovated: Ultra-Glass (Cyan Tint + Max Transparency)
                           // BOOSTED VISIBILITY: Increased opacity slightly and added stronger borders/glows
                           decoration: BoxDecoration(
                             color: isDark 
                                    ? const Color(0xFF001020).withOpacity(0.01) // Ultra Transparent (Deep Void)
                                    : const Color(0xFFFFFFFF).withOpacity(0.05),
                             borderRadius: BorderRadius.circular(16),
                             border: Border.all(
                               color: _isHovered 
                                  ? AppTheme.sanggamGold // Bright Gold on Hover
                                  : (isDark 
                                      ? AppTheme.sanggamGold.withOpacity(0.8) 
                                      : Colors.black.withOpacity(0.2)),
                               width: _isHovered ? 2.5 : 1.5 // Thicker on Hover
                             ), 
                             gradient: LinearGradient(
                                begin: Alignment.topLeft,
                                end: Alignment.bottomRight,
                                colors: [
                                  AppTheme.waveCyan.withOpacity(0.15), // Stronger Glint
                                  Colors.transparent,
                                  Colors.transparent,
                                  AppTheme.sanggamGold.withOpacity(0.1), // Stronger Gold Tint
                                ],
                                stops: const [0.0, 0.3, 0.7, 1.0],
                             ),
                             boxShadow: [
                               BoxShadow(
                                 color: AppTheme.waveCyan.withOpacity(0.1), // Stronger Glow
                                 blurRadius: 20,
                                 spreadRadius: -2,
                               ),
                             ]
                           ),
                           child: Stack(
                             children: [
                               // 1. Vignette (Depth)
                               Positioned.fill(
                                 child: Container(
                                   decoration: BoxDecoration(
                                     gradient: RadialGradient(
                                       colors: isDark
                                          ? [Colors.black.withOpacity(0.05), Colors.transparent] // Almost invisible vignette
                                          : [Colors.black.withOpacity(0.02), Colors.transparent],
                                       radius: 1.2,
                                     ),
                                   ),
                                 ),
                               ),

                               // 2. Korean Patterned Glass Etching
                               Positioned.fill(
                                 child: Opacity(
                                   opacity: 0.6, // Keep 60% Pattern
                                   child: CustomPaint(
                                     painter: _KoreanGlassPatternPainter(
                                       color: AppTheme.sanggamGold, 
                                       patternType: PatternType.cloud
                                     ),
                                   ),
                                 ),
                               ),

                               // 3. Tech Grid (REMOVED: User requested removal of meaningless lines)
                               // Positioned.fill(
                               //   child: CustomPaint(
                               //     painter: _TechGridPainter(color: AppTheme.waveCyan.withOpacity(0.04)),
                               //   ),
                               // ),
                               
                               // 4. Cinematic Sheen (Subtle & Soft)
                               Positioned.fill(
                                 child: LayoutBuilder(
                                   builder: (context, constraints) {
                                     final shimmerPos = _shimmerController.value * constraints.maxWidth * 3 - constraints.maxWidth;
                                     return Transform.translate(
                                       offset: Offset(shimmerPos, 0),
                                       child: Transform.rotate(
                                         angle: -math.pi / 5,
                                         child: ImageFiltered( 
                                            imageFilter: ui.ImageFilter.blur(sigmaX: 8, sigmaY: 8), // Softer blur
                                            child: Container(
                                              width: 60, 
                                              decoration: BoxDecoration(
                                                gradient: LinearGradient(
                                                  colors: [
                                                    Colors.transparent,
                                                    Colors.white.withOpacity(0.05),
                                                    Colors.white.withOpacity(0.15), // Soft Highlight
                                                    Colors.white.withOpacity(0.05),
                                                    Colors.transparent,
                                                  ],
                                                  stops: const [0.0, 0.4, 0.5, 0.6, 1.0],
                                                  begin: Alignment.centerLeft,
                                                  end: Alignment.centerRight,
                                                ),
                                              ),
                                            ),
                                         ),
                                       ),
                                     );
                                   },
                                 ),
                               ),

                               // 5. Content
                               Padding(
                                 padding: widget.padding,
                                 child: widget.child,
                               ),
                             ],
                           ),
                        ),
                      ),
                    ),
                  ),
                );
              }
            ),
          ),
        );
      }
    ),
    );
  }
}

enum PatternType { cloud, geometric }

class _KoreanGlassPatternPainter extends CustomPainter {
  final Color color;
  final PatternType patternType;

  _KoreanGlassPatternPainter({required this.color, required this.patternType});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.0;

    if (patternType == PatternType.cloud) {
      _drawCloudPattern(canvas, size, paint);
    } 
  }

  void _drawCloudPattern(Canvas canvas, Size size, Paint paint) {
    // Detailed "Un-Mun" (Cloud) Pattern
    // Replaces the simple 3-line curve with a proper spiral motif
    double step = 80.0;
    
    // Golden sub-pattern paint
    final goldPaint = Paint()
      ..color = AppTheme.sanggamGold.withOpacity(0.15)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5;

    for (double y = 0; y < size.height; y += step) {
      for (double x = 0; x < size.width; x += step) {
        if ((x + y) % (step * 2) == 0) continue; // Checkerboard placement
        
        final cx = x + step/2;
        final cy = y + step/2;
        
        // 1. Cloud Spiral (Main)
        final path = Path();
        // A recognizable traditional cloud tail and head
        path.moveTo(cx - 15, cy + 5);
        path.cubicTo(cx - 10, cy - 10, cx + 5, cy - 15, cx + 15, cy - 5); // Upper arch
        path.cubicTo(cx + 25, cy + 5, cx + 15, cy + 15, cx + 5, cy + 10); // Lower loop
        path.cubicTo(cx, cy + 5, cx - 10, cy + 15, cx - 20, cy + 10); // Tail
        
        canvas.drawPath(path, paint..color = color.withOpacity(0.1));

        // 2. Subtle Gold Accent (Detail)
        canvas.drawCircle(Offset(cx + 15, cy - 5), 1.5, goldPaint);
        canvas.drawCircle(Offset(cx - 5, cy + 10), 1.0, goldPaint);
      }
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

class _OrnateFramePainter extends CustomPainter {
  final Color color;
  final Color glowColor;
  final double shimmerValue;

  _OrnateFramePainter({required this.color, required this.glowColor, required this.shimmerValue});

  @override
  void paint(Canvas canvas, Size size) {
    // 1. Brushed Metal Texture
    final rect = Rect.fromLTWH(0, 0, size.width, size.height);
    final frameGradient = LinearGradient(
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
      colors: [
         const Color(0xFF5D4037), // Dark Bronze
         const Color(0xFFFFD54F), // Bright Gold
         const Color(0xFFFFF8E1), // Pale Gold
         const Color(0xFFFFB300), // Amber Gold
         const Color(0xFF4E342E), // Darker Brown
      ],
      stops: const [0.0, 0.3, 0.5, 0.7, 1.0],
      tileMode: TileMode.mirror,
    );

    // 2. Path
    final cs = 25.0; 
    final w = size.width;
    final h = size.height;
    final path = Path();
    path.moveTo(0, cs); path.lineTo(0, 10); path.quadraticBezierTo(0, 0, 10, 0); path.lineTo(cs, 0); // TL
    path.moveTo(w - cs, 0); path.lineTo(w - 10, 0); path.quadraticBezierTo(w, 0, w, 10); path.lineTo(w, cs); // TR
    path.moveTo(w, h - cs); path.lineTo(w, h - 10); path.quadraticBezierTo(w, h, w - 10, h); path.lineTo(w - cs, h); // BR
    path.moveTo(cs, h); path.lineTo(10, h); path.quadraticBezierTo(0, h, 0, h - 10); path.lineTo(0, h - cs); // BL

    // 3. Bevel Lighting
    canvas.drawPath(path.shift(const Offset(3.0, 3.0)), Paint()..color = Colors.black.withOpacity(0.9)..style = PaintingStyle.stroke..strokeWidth = 6.0);
    canvas.drawPath(path.shift(const Offset(-1.5, -1.5)), Paint()..color = Colors.white.withOpacity(0.7)..style = PaintingStyle.stroke..strokeWidth = 4.0..maskFilter = const MaskFilter.blur(BlurStyle.solid, 2));

    // 4. Main Metal Body
    final paint = Paint()
      ..shader = frameGradient.createShader(rect)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 6.0;
    canvas.drawPath(path, paint);

    // 5. Korean "Sanggam" (Inlay) Pattern - Dragon Scales/Lattice
    final patternPaint = Paint()
      ..color = const Color(0xFF3E2723).withOpacity(0.5) // Dark inlay marks
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.0;

    _drawDragonScalePattern(canvas, path, patternPaint, w, h, cs);

    // 6. Glow (Pulsating)
    final pulse = math.sin(shimmerValue * 2 * math.pi) * 0.5 + 0.5; // 0.0 to 1.0
    final glowPaint = Paint()
      ..color = glowColor.withOpacity(0.3 + (pulse * 0.2)) // Pulse Opacity
      ..style = PaintingStyle.stroke
      ..strokeWidth = 8.0 + (pulse * 4.0) // Pulse Width
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 12);
    canvas.drawPath(path, glowPaint);

    // 7. Nodes
    final nodePaint = Paint()..color = glowColor..style = PaintingStyle.fill;
    final nodeGlow = Paint()..color = glowColor.withOpacity(0.6)..maskFilter = const MaskFilter.blur(BlurStyle.normal, 6);
    _drawNode(canvas, Offset(cs, 0), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(w-cs, 0), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(w, cs), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(w, h-cs), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(w-cs, h), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(cs, h), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(0, h-cs), nodePaint, nodeGlow);
    _drawNode(canvas, Offset(0, cs), nodePaint, nodeGlow);

    // Inner Decor Lines (Delicate Sanggam Inlay)
    final innerPath = Path();
    final inset = 6.0;
    // TL Inner
    innerPath.moveTo(inset, cs); innerPath.lineTo(inset, inset + 10); innerPath.quadraticBezierTo(inset, inset, inset + 10, inset); innerPath.lineTo(cs, inset);
    // TR Inner
    innerPath.moveTo(w - cs, inset); innerPath.lineTo(w - inset - 10, inset); innerPath.quadraticBezierTo(w - inset, inset, w - inset, inset + 10); innerPath.lineTo(w - inset, cs);
    // BR Inner
    innerPath.moveTo(w - inset, h - cs); innerPath.lineTo(w - inset, h - inset - 10); innerPath.quadraticBezierTo(w - inset, h - inset, w - inset - 10, h - inset); innerPath.lineTo(w - cs, h - inset);
    // BL Inner
    innerPath.moveTo(cs, h - inset); innerPath.lineTo(inset + 10, h - inset); innerPath.quadraticBezierTo(inset, h - inset, inset, h - inset - 10); innerPath.lineTo(inset, h - cs);
    
    canvas.drawPath(innerPath, Paint()..color = color.withOpacity(0.4)..style = PaintingStyle.stroke..strokeWidth = 1.0);
    
    // Tech Accents (Cybernetic Corners)
    final techPaint = Paint()..color = glowColor.withOpacity(0.6)..style = PaintingStyle.fill;
    final double gap = 4.0;
    // TL
    canvas.drawCircle(Offset(cs + gap, gap + 2), 1.5, techPaint);
    canvas.drawCircle(Offset(gap + 2, cs + gap), 1.5, techPaint);
    // TR
    canvas.drawCircle(Offset(w - cs - gap, gap + 2), 1.5, techPaint);
    canvas.drawCircle(Offset(w - gap - 2, cs + gap), 1.5, techPaint);
    // BR
    canvas.drawCircle(Offset(w - cs - gap, h - gap - 2), 1.5, techPaint);
    canvas.drawCircle(Offset(w - gap - 2, h - cs - gap), 1.5, techPaint);
    // BL
    canvas.drawCircle(Offset(cs + gap, h - gap - 2), 1.5, techPaint);
    canvas.drawCircle(Offset(gap + 2, h - cs - gap), 1.5, techPaint);
    
    // Decor Center
    _drawLotusDecoration(canvas, Offset(w/2, 0), paint, isTop: true);
    _drawLotusDecoration(canvas, Offset(w/2, h), paint, isTop: false);
  }

  void _drawDragonScalePattern(Canvas canvas, Path framePath, Paint paint, double w, double h, double cs) {
    // Draw cross-hatch pattern only near straight edges to simulate grip/inlay
    // Top Edge Scales
    canvas.save();
    canvas.clipPath(framePath); // Clip to the frame shape
    
    double step = 6.0;
    // We manually draw lines along the approximate frame area since clipPath works
    // Diagonal lines /
    for(double i=0; i<w+h; i+=step) {
      canvas.drawLine(Offset(i, 0), Offset(0, i), paint);
    }
    // Diagonal lines \
    for(double i=-h; i<w; i+=step) {
      canvas.drawLine(Offset(i, 0), Offset(i+h, h), paint);
    }
    canvas.restore();
  }

  void _drawNode(Canvas canvas, Offset center, Paint paint, Paint glow) {
    canvas.drawCircle(center, 8.0, glow); // Larger Glow
    canvas.drawCircle(center, 3.0, paint); // Explicit Node Core
  }
  
  void _drawLotusDecoration(Canvas canvas, Offset center, Paint paint, {required bool isTop}) {
    // Renovated: "Cyber-Lotus" (Geometric + Tech)
    final yFlip = isTop ? 1.0 : -1.0;
    
    // 1. Central Diamond (Cyber-Core)
    final diamondPath = Path();
    diamondPath.moveTo(center.dx, center.dy + (8 * yFlip)); // Top/Bottom tip
    diamondPath.lineTo(center.dx + 6, center.dy + (14 * yFlip)); // Right
    diamondPath.lineTo(center.dx, center.dy + (20 * yFlip)); // Bottom Tip
    diamondPath.lineTo(center.dx - 6, center.dy + (14 * yFlip)); // Left
    diamondPath.close();
    
    canvas.drawPath(diamondPath, Paint()..color = AppTheme.sanggamGold..style = PaintingStyle.fill);
    
    // 2. Digital Wings (Circuit Lines)
    final wingPath = Path();
    wingPath.moveTo(center.dx - 6, center.dy + (14 * yFlip));
    wingPath.lineTo(center.dx - 15, center.dy + (14 * yFlip));
    wingPath.lineTo(center.dx - 20, center.dy + (8 * yFlip));
    
    wingPath.moveTo(center.dx + 6, center.dy + (14 * yFlip));
    wingPath.lineTo(center.dx + 15, center.dy + (14 * yFlip));
    wingPath.lineTo(center.dx + 20, center.dy + (8 * yFlip));
    
    canvas.drawPath(wingPath, Paint()..color = AppTheme.sanggamGold.withOpacity(0.8)..style = PaintingStyle.stroke..strokeWidth = 1.5);
    
    // 3. Glowing Dot
    canvas.drawCircle(Offset(center.dx, center.dy + (14 * yFlip)), 1.5, Paint()..color = Colors.white);
  }

  @override
  bool shouldRepaint(covariant _OrnateFramePainter oldDelegate) => oldDelegate.shimmerValue != shimmerValue;
}

class _TechGridPainter extends CustomPainter {
  final Color color;
  _TechGridPainter({required this.color});
  @override
  void paint(Canvas canvas, Size size) {
    // Very faint, thin lines for subtle texture
    final paint = Paint()..color = color..strokeWidth = 0.3;
    double step = 25.0; // Larger cells
    
    // Draw Grid
    for(double x=0; x<size.width; x+=step) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for(double y=0; y<size.height; y+=step) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
    
    // Crosshairs
    final crossPaint = Paint()..color = color.withOpacity(0.5)..strokeWidth = 1.0;
    double cx = size.width / 2;
    double cy = size.height / 2;
    canvas.drawLine(Offset(cx - 5, cy), Offset(cx + 5, cy), crossPaint);
    canvas.drawLine(Offset(cx, cy - 5), Offset(cx, cy + 5), crossPaint);
  }
  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
