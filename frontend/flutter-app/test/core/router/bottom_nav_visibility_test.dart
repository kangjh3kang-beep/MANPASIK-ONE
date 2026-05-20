import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/router/bottom_nav_visibility.dart';

void main() {
  group('shouldHideGlobalDock', () {
    test('/home and /measure routes are hidden', () {
      expect(shouldHideGlobalDock('/home'), isTrue);
      expect(shouldHideGlobalDock('/measure'), isTrue);
      expect(shouldHideGlobalDock('/medical/video-call/session-1'), isTrue);
    });

    test('child routes are also hidden', () {
      expect(shouldHideGlobalDock('/home/insight'), isTrue);
      expect(shouldHideGlobalDock('/measure/result'), isTrue);
    });

    test('non-target routes remain visible', () {
      expect(shouldHideGlobalDock('/data/monitoring'), isFalse);
      expect(shouldHideGlobalDock('/market'), isFalse);
      expect(shouldHideGlobalDock('/community/challenge'), isFalse);
      expect(shouldHideGlobalDock('/family'), isFalse);
      expect(shouldHideGlobalDock('/settings'), isFalse);
    });

    test('prefix lookalikes are not hidden', () {
      expect(shouldHideGlobalDock('/homework'), isFalse);
      expect(shouldHideGlobalDock('/measured'), isFalse);
      expect(shouldHideGlobalDock('/measurements'), isFalse);
    });
  });

  group('shouldHideGlobalTopFrame', () {
    test('/home and /measure routes are hidden', () {
      expect(shouldHideGlobalTopFrame('/home'), isTrue);
      expect(shouldHideGlobalTopFrame('/measure'), isTrue);
      expect(shouldHideGlobalTopFrame('/medical/video-call/session-1'), isTrue);
    });

    test('child routes are also hidden', () {
      expect(shouldHideGlobalTopFrame('/home/profile'), isTrue);
      expect(shouldHideGlobalTopFrame('/measure/result'), isTrue);
    });

    test('non-target routes keep global top frame', () {
      expect(shouldHideGlobalTopFrame('/data'), isFalse);
      expect(shouldHideGlobalTopFrame('/market'), isFalse);
      expect(shouldHideGlobalTopFrame('/settings'), isFalse);
    });

    test('prefix lookalikes are not hidden', () {
      expect(shouldHideGlobalTopFrame('/homework'), isFalse);
      expect(shouldHideGlobalTopFrame('/measured'), isFalse);
    });
  });
}
