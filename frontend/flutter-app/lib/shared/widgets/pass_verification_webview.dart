import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

/// PASS 본인인증 WebView 위젯
///
/// PASS 본인인증 페이지를 WebView로 호출하고
/// 인증 성공/실패 시 콜백을 반환합니다.
class PassVerificationWebView extends StatefulWidget {
  const PassVerificationWebView({
    super.key,
    required this.merchantId,
    this.callbackUrl = 'manpasik://identity/callback',
    this.verificationUrl,
  });

  final String merchantId;
  final String callbackUrl;

  /// 서버에서 생성한 본인인증 요청 URL (없으면 기본 PASS URL 사용)
  final String? verificationUrl;

  @override
  State<PassVerificationWebView> createState() =>
      _PassVerificationWebViewState();
}

class _PassVerificationWebViewState extends State<PassVerificationWebView> {
  late final WebViewController _controller;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onPageStarted: (_) => setState(() => _isLoading = true),
          onPageFinished: (_) => setState(() => _isLoading = false),
          onNavigationRequest: _handleNavigation,
        ),
      )
      ..addJavaScriptChannel(
        'PassResult',
        onMessageReceived: _handleJsMessage,
      );

    if (widget.verificationUrl != null) {
      _controller.loadRequest(Uri.parse(widget.verificationUrl!));
    } else {
      _controller.loadHtmlString(_buildVerificationHtml());
    }
  }

  void _handleJsMessage(JavaScriptMessage message) {
    // JavaScript에서 인증 결과를 전달받는 채널
    final data = message.message;
    if (data.startsWith('success:')) {
      final parts = data.substring(8).split('|');
      Navigator.pop(
        context,
        PassVerificationResult(
          success: true,
          name: parts.isNotEmpty ? parts[0] : null,
          phone: parts.length > 1 ? parts[1] : null,
          birthDate: parts.length > 2 ? parts[2] : null,
          ci: parts.length > 3 ? parts[3] : null,
        ),
      );
    } else if (data.startsWith('fail:')) {
      Navigator.pop(
        context,
        PassVerificationResult(
          success: false,
          errorMessage: data.substring(5),
        ),
      );
    }
  }

  NavigationDecision _handleNavigation(NavigationRequest request) {
    final uri = Uri.tryParse(request.url);
    if (uri == null) return NavigationDecision.navigate;

    // 인증 성공 콜백
    if (request.url.startsWith(widget.callbackUrl) ||
        request.url.contains('identity/callback')) {
      final name = uri.queryParameters['name'];
      final phone = uri.queryParameters['phone'];
      final birthDate = uri.queryParameters['birth_date'];
      final ci = uri.queryParameters['ci'];
      final success = uri.queryParameters['success'] != 'false';

      Navigator.pop(
        context,
        PassVerificationResult(
          success: success,
          name: name,
          phone: phone,
          birthDate: birthDate,
          ci: ci,
          errorMessage: success ? null : uri.queryParameters['message'],
        ),
      );
      return NavigationDecision.prevent;
    }

    return NavigationDecision.navigate;
  }

  String _buildVerificationHtml() {
    return '''
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <style>
    body { display: flex; justify-content: center; align-items: center;
           min-height: 100vh; margin: 0; background: #f8f9fa;
           font-family: -apple-system, sans-serif; }
    .container { text-align: center; padding: 24px; }
    .title { font-size: 20px; font-weight: bold; color: #333; }
    .desc { margin-top: 12px; color: #666; font-size: 14px; }
    .btn { margin-top: 24px; padding: 14px 32px; border-radius: 8px;
           background: #3182f6; color: white; border: none;
           font-size: 16px; cursor: pointer; }
    .btn:hover { background: #1b6fef; }
  </style>
</head>
<body>
  <div class="container">
    <div class="title">PASS 본인인증</div>
    <div class="desc">휴대폰 본인인증을 진행합니다.<br/>
      통신사 인증 앱이 실행됩니다.</div>
    <button class="btn" onclick="startVerification()">본인인증 시작</button>
  </div>
  <script>
    function startVerification() {
      // PASS 인증 모듈 초기화
      // 실제 환경에서는 서버에서 발급받은 인증 URL로 리다이렉트
      // 개발 환경에서는 시뮬레이션 결과 반환
      try {
        if (window.PassResult) {
          // 서버 연동 시 실제 PASS 인증 결과가 콜백됨
          // 시뮬레이션 모드: 성공 결과 반환
          PassResult.postMessage('success:홍길동|010-1234-5678|1990-01-15|CI_HASH_VALUE');
        } else {
          window.location.href = '${widget.callbackUrl}?success=true&name=홍길동&phone=010-1234-5678&birth_date=1990-01-15';
        }
      } catch (e) {
        if (window.PassResult) {
          PassResult.postMessage('fail:인증 모듈 오류: ' + e.message);
        }
      }
    }
  </script>
</body>
</html>
''';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('본인인증'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.pop(context, null),
        ),
      ),
      body: Stack(
        children: [
          WebViewWidget(controller: _controller),
          if (_isLoading) const Center(child: CircularProgressIndicator()),
        ],
      ),
    );
  }
}

/// PASS 본인인증 결과
class PassVerificationResult {
  final bool success;
  final String? name;
  final String? phone;
  final String? birthDate;
  final String? ci;
  final String? errorMessage;

  const PassVerificationResult({
    required this.success,
    this.name,
    this.phone,
    this.birthDate,
    this.ci,
    this.errorMessage,
  });
}
