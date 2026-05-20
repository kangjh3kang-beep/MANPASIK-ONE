import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:webview_cef/webview_cef.dart' as cef;
import 'package:webview_flutter/webview_flutter.dart' as wv;

const _linuxViewerRev = '20260221_glb5_v6';

bool _supportsNativeHoloWebView() {
  switch (defaultTargetPlatform) {
    case TargetPlatform.android:
    case TargetPlatform.iOS:
    case TargetPlatform.macOS:
    case TargetPlatform.linux:
      return true;
    case TargetPlatform.windows:
    case TargetPlatform.fuchsia:
      return false;
  }
}

Widget? buildHoloBodyNative(String profileKey, String colorHex) {
  if (!_supportsNativeHoloWebView()) {
    return null;
  }

  if (defaultTargetPlatform == TargetPlatform.linux) {
    return _LinuxInlineHoloView(
      profileKey: profileKey,
      colorHex: colorHex,
    );
  }

  return _NativeHoloBodyWebView(
    profileKey: profileKey,
    colorHex: colorHex,
  );
}

class _NativeHoloBodyWebView extends StatefulWidget {
  const _NativeHoloBodyWebView({
    required this.profileKey,
    required this.colorHex,
  });

  final String profileKey;
  final String colorHex;

  @override
  State<_NativeHoloBodyWebView> createState() => _NativeHoloBodyWebViewState();
}

class _NativeHoloBodyWebViewState extends State<_NativeHoloBodyWebView> {
  late final wv.WebViewController _controller;
  bool _isPageReady = false;

  @override
  void initState() {
    super.initState();
    _controller = wv.WebViewController()
      ..setJavaScriptMode(wv.JavaScriptMode.unrestricted)
      ..setBackgroundColor(const Color(0x00000000))
      ..setNavigationDelegate(
        wv.NavigationDelegate(
          onPageFinished: (_) {
            _isPageReady = true;
            _syncViewerState();
          },
        ),
      )
      ..loadFlutterAsset('web/holo_viewer.html');
  }

  @override
  void didUpdateWidget(covariant _NativeHoloBodyWebView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.profileKey != widget.profileKey ||
        oldWidget.colorHex != widget.colorHex) {
      _syncViewerState();
    }
  }

  Future<void> _syncViewerState() async {
    if (!_isPageReady) {
      return;
    }

    final profile = jsonEncode(widget.profileKey);
    final color = jsonEncode(widget.colorHex);

    try {
      await _controller.runJavaScript(
        'window.postMessage(JSON.stringify({type:"setProfile", value:$profile}), "*");',
      );
      await _controller.runJavaScript(
        'window.postMessage(JSON.stringify({type:"setColor", value:$color}), "*");',
      );
    } catch (_) {
      // Ignore runtime bridge errors and keep the screen stable.
    }
  }

  @override
  Widget build(BuildContext context) {
    return wv.WebViewWidget(controller: _controller);
  }
}

class _LinuxInlineHoloView extends StatefulWidget {
  const _LinuxInlineHoloView({
    required this.profileKey,
    required this.colorHex,
  });

  final String profileKey;
  final String colorHex;

  @override
  State<_LinuxInlineHoloView> createState() => _LinuxInlineHoloViewState();
}

class _LinuxInlineHoloViewState extends State<_LinuxInlineHoloView> {
  static Future<void>? _cefInitFuture;

  late final cef.WebViewController _controller;
  bool _ready = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _controller = cef.WebviewManager().createWebView(
      loading: const Center(
        child: Text(
          'Loading hologram engine...',
          textAlign: TextAlign.center,
          style: TextStyle(fontSize: 12, color: Color(0xFF7FE8FF)),
        ),
      ),
    );
    unawaited(_bootstrap());
  }

  @override
  void didUpdateWidget(covariant _LinuxInlineHoloView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.profileKey != widget.profileKey ||
        oldWidget.colorHex != widget.colorHex) {
      unawaited(_syncUrl());
    }
  }

  @override
  void dispose() {
    unawaited(_controller.dispose());
    super.dispose();
  }

  static Future<void> _ensureCefInitialized() {
    return _cefInitFuture ??= cef.WebviewManager().initialize();
  }

  Future<String> _viewerUrl() async {
    final server = await _HoloAssetServer.ensureStarted();
    return Uri(
      scheme: 'http',
      host: '127.0.0.1',
      port: server.port,
      path: '/holo_viewer.html',
      queryParameters: <String, String>{
        'profile': widget.profileKey,
        'color': widget.colorHex,
        'rev': _linuxViewerRev,
      },
    ).toString();
  }

  Future<void> _bootstrap() async {
    try {
      await _ensureCefInitialized();
      final url = await _viewerUrl();
      await _controller.initialize(url);
      if (!mounted) {
        return;
      }
      setState(() {
        _ready = true;
        _error = null;
      });
    } catch (e) {
      if (!mounted) {
        return;
      }
      setState(() {
        _ready = false;
        _error = e.toString();
      });
    }
  }

  Future<void> _syncUrl() async {
    if (!_ready || !_controller.value) {
      return;
    }

    try {
      await _controller.loadUrl(await _viewerUrl());
    } catch (e) {
      if (!mounted) {
        return;
      }
      setState(() => _error = e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) {
      return Center(
        child: Text(
          _error!,
          textAlign: TextAlign.center,
          style: const TextStyle(
            fontSize: 12,
            color: Color(0xFFFF6E6E),
          ),
        ),
      );
    }

    return ValueListenableBuilder<bool>(
      valueListenable: _controller,
      builder: (context, ready, _) {
        if (!ready || !_ready) {
          return _controller.loadingWidget;
        }
        return _controller.webviewWidget;
      },
    );
  }
}

class _HoloAssetServer {
  static HttpServer? _server;
  static Future<HttpServer>? _serverFuture;

  static Future<HttpServer> ensureStarted() {
    final existing = _server;
    if (existing != null) {
      return Future<HttpServer>.value(existing);
    }

    _serverFuture ??= _createServer();
    return _serverFuture!;
  }

  static Future<HttpServer> _createServer() async {
    try {
      final root = _resolveWebAssetRoot();
      final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);
      server.listen((HttpRequest request) {
        unawaited(_serveAssetFile(root, request));
      });
      _server = server;
      return server;
    } catch (e) {
      _serverFuture = null;
      rethrow;
    }
  }

  static Directory _resolveWebAssetRoot() {
    final executableDir = File(Platform.resolvedExecutable).parent.path;
    final candidates = <String>[
      '$executableDir/data/flutter_assets/web',
      '$executableDir/flutter_assets/web',
      '${Directory.current.path}/data/flutter_assets/web',
      '${Directory.current.path}/flutter_assets/web',
      '${Directory.current.path}/build/linux/x64/debug/bundle/data/flutter_assets/web',
      '${Directory.current.path}/build/linux/x64/release/bundle/data/flutter_assets/web',
      '${Directory.current.path}/web',
    ];

    for (final candidate in candidates) {
      final dir = Directory(candidate);
      if (dir.existsSync()) {
        return dir;
      }
    }

    throw StateError('Unable to locate web assets directory for Linux GLB viewer.');
  }

  static Future<void> _serveAssetFile(Directory root, HttpRequest request) async {
    final response = request.response;
    try {
      var path = request.uri.path;
      if (path.isEmpty || path == '/') {
        path = '/holo_viewer.html';
      }

      if (path.contains('..')) {
        response.statusCode = HttpStatus.forbidden;
        await response.close();
        return;
      }

      final relative = path.startsWith('/') ? path.substring(1) : path;
      final file = File('${root.path}/$relative');
      if (!file.existsSync()) {
        response.statusCode = HttpStatus.notFound;
        await response.close();
        return;
      }

      response.headers.set(HttpHeaders.contentTypeHeader, _contentTypeFor(file.path));
      await response.addStream(file.openRead());
      await response.close();
    } catch (e) {
      debugPrint('Linux inline hologram asset server error: $e');
      try {
        response.statusCode = HttpStatus.internalServerError;
      } catch (_) {}
      await response.close();
    }
  }

  static String _contentTypeFor(String path) {
    if (path.endsWith('.html')) return 'text/html; charset=utf-8';
    if (path.endsWith('.js')) return 'application/javascript; charset=utf-8';
    if (path.endsWith('.json')) return 'application/json; charset=utf-8';
    if (path.endsWith('.css')) return 'text/css; charset=utf-8';
    if (path.endsWith('.svg')) return 'image/svg+xml';
    if (path.endsWith('.png')) return 'image/png';
    if (path.endsWith('.jpg') || path.endsWith('.jpeg')) return 'image/jpeg';
    if (path.endsWith('.ico')) return 'image/x-icon';
    if (path.endsWith('.glb')) return 'model/gltf-binary';
    return 'application/octet-stream';
  }
}
