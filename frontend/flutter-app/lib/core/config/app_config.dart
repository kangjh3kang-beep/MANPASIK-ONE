/// ManPaSik 앱 환경 설정
///
/// 환경별(dev/staging/prod) 설정을 중앙 관리합니다.
/// BaseUrl, SSL 핀, 기능 플래그 등을 포함합니다.
class AppConfig {
  AppConfig._();

  static AppEnvironment _env = AppEnvironment.dev;
  static AppEnvironment get environment => _env;

  /// 앱 시작 시 환경 설정 (main.dart에서 호출)
  static void initialize({AppEnvironment env = AppEnvironment.dev}) {
    _env = env;
  }

  /// REST API Gateway Base URL
  static String get baseUrl => switch (_env) {
        AppEnvironment.dev => 'http://localhost:8080/api/v1',
        AppEnvironment.staging => 'https://staging-api.manpasik.com/api/v1',
        AppEnvironment.prod => 'https://api.manpasik.com/api/v1',
      };

  /// gRPC 게이트웨이 호스트
  static String get grpcHost => switch (_env) {
        AppEnvironment.dev => 'localhost',
        AppEnvironment.staging => 'staging-grpc.manpasik.com',
        AppEnvironment.prod => 'grpc.manpasik.com',
      };

  /// gRPC 포트
  static int get grpcPort => switch (_env) {
        AppEnvironment.dev => 50051,
        AppEnvironment.staging => 443,
        AppEnvironment.prod => 443,
      };

  /// SSL Pinning 허용 호스트
  static List<String> get allowedHosts => switch (_env) {
        AppEnvironment.dev => ['localhost', '10.0.2.2', '127.0.0.1'],
        AppEnvironment.staging => [
            'staging-api.manpasik.com',
            'staging-grpc.manpasik.com',
          ],
        AppEnvironment.prod => [
            'api.manpasik.com',
            'gateway.manpasik.com',
            'auth.manpasik.com',
            'grpc.manpasik.com',
          ],
      };

  /// SSL 인증서 SHA-256 핀 (프로덕션 전용)
  /// 인증서 갱신 시 반드시 백업 핀을 포함해야 합니다.
  static List<String> get certificatePins => switch (_env) {
        AppEnvironment.dev => [], // 개발 환경에서는 핀 검증 비활성화
        AppEnvironment.staging => [
            // 스테이징 인증서 핀 (Let's Encrypt)
            'sha256/jQJTbIh0grw0/1TkHSumWb+Fs0Ggogr621gT3PvPKG0=',
          ],
        AppEnvironment.prod => [
            // 프로덕션 기본 인증서 핀
            'sha256/YLh1dUR9y6Kja30RrAn7JKnbQG/uEtLMkBgFF2Fuihg=',
            // 프로덕션 백업 인증서 핀 (갱신 대비)
            'sha256/Vjs8r4z+80wjNcr1YKepWQboSIRi63WsWXhIMN+eWys=',
            // Let's Encrypt ISRG Root X1 (중간 인증서)
            'sha256/C5+lpZ7tcVwmwQIMcRtPbsQtWLABXhQzejna0wHFr8M=',
          ],
      };

  /// SSL Pinning 활성화 여부
  static bool get sslPinningEnabled => _env == AppEnvironment.prod;

  /// 디버그 모드 여부
  static bool get isDebug => _env == AppEnvironment.dev;

  /// WebSocket URL (실시간 알림용)
  static String get wsUrl => switch (_env) {
        AppEnvironment.dev => 'ws://localhost:8080/ws',
        AppEnvironment.staging => 'wss://staging-api.manpasik.com/ws',
        AppEnvironment.prod => 'wss://api.manpasik.com/ws',
      };

  /// 119 자동 신고 전화번호
  static String get emergencyNumber => '119';

  /// 기능 플래그
  static bool get enableRustFfi => _env == AppEnvironment.prod;
  static bool get enableWebRtc => _env != AppEnvironment.dev;
  static bool get enableHealthKit => _env != AppEnvironment.dev;
}

/// 앱 실행 환경
enum AppEnvironment {
  dev,
  staging,
  prod,
}
