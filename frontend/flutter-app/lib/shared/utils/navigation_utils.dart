import 'package:flutter/widgets.dart';
import 'package:go_router/go_router.dart';

extension NavigationContextX on BuildContext {
  void popOrHome({String fallback = '/home'}) {
    final router = GoRouter.of(this);
    if (router.canPop()) {
      router.pop();
      return;
    }
    router.go(fallback);
  }
}
