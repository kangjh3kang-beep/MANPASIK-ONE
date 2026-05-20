import 'dart:ui_web' as ui_web;
import 'package:flutter/widgets.dart';
import 'package:web/web.dart' as web;

final _registeredTypes = <String>{};
const _holoViewRev = '20260223_rev21';

String _buildViewerUrl(String profileKey, String colorHex) {
  final viewerUri = Uri.base
      .resolve('holo_viewer.html')
      .replace(queryParameters: <String, String>{
    'profile': profileKey,
    'color': colorHex,
    'rev': _holoViewRev,
  });
  return viewerUri.toString();
}

/// Web implementation: embeds holo_viewer.html via IFrame.
/// GLB is the default render path.
Widget buildHoloBodyWeb(String profileKey, String colorHex) {
  final viewType = 'holo-body-$_holoViewRev-$profileKey-$colorHex';
  final viewerUrl = _buildViewerUrl(profileKey, colorHex);

  if (!_registeredTypes.contains(viewType)) {
    ui_web.platformViewRegistry.registerViewFactory(viewType, (int viewId) {
      final container = web.HTMLDivElement();
      container.style
        ..width = '100%'
        ..height = '100%'
        ..overflow = 'visible' // USER FIX: Changed from 'hidden' to prevent potential clipping
        ..position = 'relative'
        ..backgroundColor = 'transparent';

      final iframe = web.HTMLIFrameElement()
        ..src = viewerUrl
        ..allow = 'autoplay; webgl; xr-spatial-tracking'
        ..setAttribute('allowtransparency', 'true')
        ..setAttribute('frameborder', '0')
        ..setAttribute('scrolling', 'no');
      iframe.style
        ..width = '100%'
        ..height = '100%'
        ..border = 'none'
        ..backgroundColor = 'transparent'
        ..position = 'absolute'
        ..top = '0'
        ..left = '0';

      container.append(iframe);
      return container;
    });
    _registeredTypes.add(viewType);
  }

  return HtmlElementView(viewType: viewType);
}
