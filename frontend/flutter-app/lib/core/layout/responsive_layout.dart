import 'package:flutter/material.dart';

/// ManPaSik Adaptive Layout Engine.
/// Switches between Mobile, Tablet, and Desktop UI based on screen width.
class ResponsiveLayout extends StatelessWidget {
  final Widget mobile;
  final Widget? tablet;
  final Widget desktop;

  const ResponsiveLayout({
    super.key,
    required this.mobile,
    this.tablet,
    required this.desktop,
  });

  // Breakpoints standardized for ManPaSik ecosystem.
  static const double mobileLimit = 800; // Mobile < 800

  static bool isMobile(BuildContext context) =>
      MediaQuery.of(context).size.width < mobileLimit;

  static bool isDesktop(BuildContext context) =>
      MediaQuery.of(context).size.width >= mobileLimit;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        if (constraints.maxWidth >= mobileLimit) {
          return desktop;
        } else {
          return mobile;
        }
      },
    );
  }
}
