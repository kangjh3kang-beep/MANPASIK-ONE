/// ManPaSik 앱 상수 정의
class AppConstants {
  AppConstants._();

  // ── 앱 정보 ──
  static const String appName = '만파식';
  static const String appNameEn = 'ManPaSik';
  static const String appVersion = '1.0.0';

  // ── API 엔드포인트 (S5에서 환경별 분리) ──
  static const String baseUrl = 'http://localhost:8000';
  static const String grpcHost = 'localhost';
  static const int grpcAuthPort = 50051;
  static const int grpcUserPort = 50052;
  static const int grpcDevicePort = 50053;
  static const int grpcMeasurementPort = 50054;
  static const int grpcAdminPort = 50055;
  static const int grpcAiInferencePort = 50058;

  // ── 토큰 ──
  static const Duration accessTokenExpiry = Duration(minutes: 15);
  static const Duration refreshTokenExpiry = Duration(days: 7);

  // ── 구독 티어 ──
  static const String tierFree = 'FREE';
  static const String tierBasic = 'BASIC';
  static const String tierPro = 'PRO';
  static const String tierClinical = 'CLINICAL';

  // ── 디바이스 제한 (티어별) ──
  static const Map<String, int> maxDevicesPerTier = {
    tierFree: 1,
    tierBasic: 3,
    tierPro: 5,
    tierClinical: 10,
  };

  // ── UI ──
  static const double borderRadiusSmall = 8.0;
  static const double borderRadiusMedium = 16.0;
  static const double borderRadiusLarge = 24.0;
  static const double paddingSmall = 8.0;
  static const double paddingMedium = 16.0;
  static const double paddingLarge = 24.0;
}
