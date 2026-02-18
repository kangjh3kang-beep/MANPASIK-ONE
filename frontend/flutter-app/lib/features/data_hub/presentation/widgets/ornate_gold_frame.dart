import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'dart:math' as math;
import 'dart:ui' as ui;

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

  @override
  void dispose() {
    _shimmerController.dispose();
    _breathingController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return AnimatedBuilder(
      animation: _breathingController,
      builder: (context, _) {
        final scale = 1.0 + (_breathingController.value * 0.005); // Micro-breathing
        
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
                           decoration: BoxDecoration(
                             // Theme-aware Glass Background
                             color: isDark 
                                ? const Color(0xFF051525).withOpacity(0.15) // Dark Mode: Deep Space Glass
                                : const Color(0xFF1A1A1A).withOpacity(0.03), // White Mode: Very Sheer Ink Glass
                             borderRadius: BorderRadius.circular(16),
                             border: Border.all(
                               color: isDark 
                                  ? Colors.white.withOpacity(0.08) 
                                  : Colors.black.withOpacity(0.05), // Subtle rim
                               width: 0.5
                             ), 
                             boxShadow: [
                               BoxShadow(
                                 color: Colors.black.withOpacity(isDark ? 0.3 : 0.05), // Lighter shadow in White Mode
                                 blurRadius: 15,
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
                                          ? [Colors.black.withOpacity(0.5), Colors.transparent]
                                          : [Colors.black.withOpacity(0.05), Colors.transparent], // Lighter Vignette in White Mode
                                       radius: 0.9,
                                     ),
                                   ),
                                 ),
                               ),

                               // 2. Tech Grid (Subtle)
                               Positioned.fill(
                                 child: CustomPaint(
                                   painter: _TechGridPainter(color: AppTheme.waveCyan.withOpacity(0.04)),
                                 ),
                               ),
                               
                               // 3. Cinematic Sheen (Soft & Blurred)
                               Positioned.fill(
                                 child: LayoutBuilder(
                                   builder: (context, constraints) {
                                     final shimmerPos = _shimmerController.value * constraints.maxWidth * 3 - constraints.maxWidth;
                                     return Transform.translate(
                                       offset: Offset(shimmerPos, 0),
                                       child: Transform.rotate(
                                         angle: -math.pi / 5,
                                         child: ImageFiltered( // NEW: Blur effect on sheen
                                            imageFilter: ui.ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                                            child: Container(
                                              width: 30, // Slightly wider to account for blur
                                              decoration: BoxDecoration(
                                                gradient: LinearGradient(
                                                  colors: [
                                                    Colors.transparent,
                                                    Colors.white.withOpacity(0.08), // Much subtler (was 0.4)
                                                    Colors.transparent,
                                                  ],
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

                               // 4. Content
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
    );
  }
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

    // 6. Glow
    final glowPaint = Paint()
      ..color = glowColor.withOpacity(0.4)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 10.0
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
    final yOffset = isTop ? -2.0 : 2.0;
    final path = Path();
    path.moveTo(center.dx - 12, center.dy + yOffset);
    path.quadraticBezierTo(center.dx, center.dy + (isTop ? 12 : -12), center.dx + 12, center.dy + yOffset);
    path.addRect(Rect.fromCenter(center: Offset(center.dx, center.dy + (isTop ? 6 : -6)), width: 4, height: 4));
    canvas.drawPath(path, paint..style = PaintingStyle.fill);
    paint.style = PaintingStyle.stroke; // Reset
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
