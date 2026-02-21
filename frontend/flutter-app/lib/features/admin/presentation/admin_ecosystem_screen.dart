import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

// ─────────────────────────────────────────────
// 데이터 모델
// ─────────────────────────────────────────────

/// 생태계 설정 항목
class EcoConfigItem {
  final String key;
  final String label;
  final String value;
  final String type; // string, number, boolean, secret, select, url
  final String description;
  final List<String> allowedValues;
  final bool restartRequired;

  const EcoConfigItem({
    required this.key,
    required this.label,
    required this.value,
    this.type = 'string',
    this.description = '',
    this.allowedValues = const [],
    this.restartRequired = false,
  });

  EcoConfigItem copyWith({String? value}) =>
      EcoConfigItem(
        key: key,
        label: label,
        value: value ?? this.value,
        type: type,
        description: description,
        allowedValues: allowedValues,
        restartRequired: restartRequired,
      );
}

/// 마이크로서비스 정보
class ServiceInfo {
  final String name;
  final String displayName;
  final int grpcPort;
  final String status; // up, down, degraded
  final String category; // infra, user, medical, ai, commerce, utility

  const ServiceInfo({
    required this.name,
    required this.displayName,
    required this.grpcPort,
    this.status = 'up',
    required this.category,
  });
}

// ─────────────────────────────────────────────
// Riverpod 상태 관리
// ─────────────────────────────────────────────

/// 생태계 설정 상태
class EcosystemState {
  final Map<String, List<EcoConfigItem>> sections;
  final bool isLoading;
  final String? successMessage;

  const EcosystemState({
    this.sections = const {},
    this.isLoading = false,
    this.successMessage,
  });

  EcosystemState copyWith({
    Map<String, List<EcoConfigItem>>? sections,
    bool? isLoading,
    String? successMessage,
  }) {
    return EcosystemState(
      sections: sections ?? this.sections,
      isLoading: isLoading ?? this.isLoading,
      successMessage: successMessage,
    );
  }
}

/// 생태계 설정 Notifier
class EcosystemNotifier extends StateNotifier<EcosystemState> {
  EcosystemNotifier() : super(const EcosystemState(isLoading: true)) {
    _loadDefaults();
  }

  void _loadDefaults() {
    state = EcosystemState(sections: _buildAllSections());
  }

  /// 설정값 업데이트
  void updateConfig(String section, String key, String newValue) {
    final updated = Map<String, List<EcoConfigItem>>.from(state.sections);
    final items = updated[section];
    if (items == null) return;

    updated[section] = items.map((item) {
      if (item.key == key) return item.copyWith(value: newValue);
      return item;
    }).toList();

    state = state.copyWith(
      sections: updated,
      successMessage: '설정이 저장되었습니다.',
    );
  }

  void clearMessage() {
    state = state.copyWith(successMessage: null);
  }

  /// 전체 섹션 데이터 구성
  static Map<String, List<EcoConfigItem>> _buildAllSections() {
    return {
      'external_keys': _externalKeys,
      'feature_flags': _featureFlags,
      'infra': _infraConfigs,
      'security': _securityConfigs,
      'notification': _notificationConfigs,
      'ai_measurement': _aiMeasurementConfigs,
      'payment': _paymentConfigs,
      'environment': _environmentConfigs,
    };
  }

  // ── 외부 연동 키 ──
  static final _externalKeys = [
    const EcoConfigItem(
      key: 'kakao.native_app_key',
      label: '카카오 네이티브 앱 키',
      value: 'YOUR_NATIVE_APP_KEY',
      type: 'secret',
      description: '카카오 개발자 콘솔에서 발급받은 네이티브 앱 키',
    ),
    const EcoConfigItem(
      key: 'kakao.rest_api_key',
      label: '카카오 REST API 키',
      value: '',
      type: 'secret',
      description: '카카오 REST API 호출에 사용되는 키',
    ),
    const EcoConfigItem(
      key: 'kakao.javascript_key',
      label: '카카오 JavaScript 키',
      value: '',
      type: 'secret',
      description: '카카오 웹 SDK 초기화에 사용되는 키',
    ),
    const EcoConfigItem(
      key: 'toss.secret_key',
      label: 'Toss Payments 시크릿 키',
      value: '',
      type: 'secret',
      description: 'Toss Payments PG 연동 시크릿 키',
    ),
    const EcoConfigItem(
      key: 'toss.api_url',
      label: 'Toss Payments API URL',
      value: 'https://api.tosspayments.com',
      type: 'url',
      description: 'Toss Payments REST API 엔드포인트',
    ),
    const EcoConfigItem(
      key: 'fcm.server_key',
      label: 'Firebase 서버 키',
      value: '',
      type: 'secret',
      description: 'Firebase Cloud Messaging 서버 키',
    ),
    const EcoConfigItem(
      key: 'fcm.project_id',
      label: 'Firebase 프로젝트 ID',
      value: '',
      type: 'string',
      description: 'Firebase 프로젝트 식별자',
    ),
    const EcoConfigItem(
      key: 'keycloak.url',
      label: 'Keycloak URL',
      value: 'http://keycloak:9090',
      type: 'url',
      description: 'Keycloak OIDC 인증 서버 주소',
    ),
    const EcoConfigItem(
      key: 'keycloak.realm',
      label: 'Keycloak Realm',
      value: 'manpasik',
      type: 'string',
      description: 'Keycloak Realm 이름',
    ),
    const EcoConfigItem(
      key: 'keycloak.client_id',
      label: 'Keycloak Client ID',
      value: 'manpasik-api',
      type: 'string',
      description: 'Keycloak 클라이언트 식별자',
    ),
    const EcoConfigItem(
      key: 'keycloak.client_secret',
      label: 'Keycloak Client Secret',
      value: '',
      type: 'secret',
      description: 'Keycloak 클라이언트 시크릿',
    ),
    const EcoConfigItem(
      key: 'google.oauth_client_id',
      label: 'Google OAuth Client ID',
      value: '',
      type: 'secret',
      description: 'Google 소셜 로그인 클라이언트 ID',
    ),
    const EcoConfigItem(
      key: 'apple.service_id',
      label: 'Apple Service ID',
      value: '',
      type: 'secret',
      description: 'Apple 로그인 서비스 ID',
    ),
  ];

  // ── 기능 플래그 ──
  static final _featureFlags = [
    const EcoConfigItem(key: 'flag.kakao_login', label: '카카오 로그인', value: 'true', type: 'boolean', description: '카카오 소셜 로그인 활성화'),
    const EcoConfigItem(key: 'flag.google_login', label: 'Google 로그인', value: 'true', type: 'boolean', description: 'Google 소셜 로그인 활성화'),
    const EcoConfigItem(key: 'flag.apple_login', label: 'Apple 로그인', value: 'true', type: 'boolean', description: 'Apple 소셜 로그인 활성화'),
    const EcoConfigItem(key: 'flag.rust_ffi', label: 'Rust FFI 엔진', value: 'false', type: 'boolean', description: 'Rust 네이티브 코어 엔진 활성화'),
    const EcoConfigItem(key: 'flag.vision_analyzer', label: '비전 분석기', value: 'true', type: 'boolean', description: '카메라 카트리지 비전 분석 활성화'),
    const EcoConfigItem(key: 'flag.telemedicine', label: '원격 진료', value: 'true', type: 'boolean', description: '원격 의료 상담 기능 활성화'),
    const EcoConfigItem(key: 'flag.marketplace', label: '마켓플레이스', value: 'true', type: 'boolean', description: '건강 식품/보조제 마켓 활성화'),
    const EcoConfigItem(key: 'flag.family_sharing', label: '가족 공유', value: 'true', type: 'boolean', description: '가족 그룹 데이터 공유 활성화'),
    const EcoConfigItem(key: 'flag.ai_coaching', label: 'AI 코칭', value: 'true', type: 'boolean', description: 'AI 맞춤 건강 코칭 활성화'),
    const EcoConfigItem(key: 'flag.push_notification', label: '푸시 알림', value: 'true', type: 'boolean', description: 'FCM 푸시 알림 발송 활성화'),
    const EcoConfigItem(key: 'flag.demo_mode', label: '데모 모드', value: 'true', type: 'boolean', description: '비회원 가상 데이터 체험 모드'),
    const EcoConfigItem(key: 'flag.dark_mode_default', label: '다크 모드 기본값', value: 'true', type: 'boolean', description: '앱 기본 테마를 다크 모드로 설정'),
  ];

  // ── 인프라 설정 ──
  static final _infraConfigs = [
    const EcoConfigItem(key: 'db.host', label: 'PostgreSQL 호스트', value: 'postgres', type: 'string', description: '메인 데이터베이스 호스트'),
    const EcoConfigItem(key: 'db.port', label: 'PostgreSQL 포트', value: '5432', type: 'number', description: '메인 데이터베이스 포트'),
    const EcoConfigItem(key: 'db.name', label: 'DB 이름', value: 'manpasik', type: 'string', description: '데이터베이스 이름'),
    const EcoConfigItem(key: 'db.user', label: 'DB 사용자', value: 'manpasik', type: 'string', description: '데이터베이스 접속 사용자'),
    const EcoConfigItem(key: 'db.password', label: 'DB 비밀번호', value: '', type: 'secret', description: '데이터베이스 접속 비밀번호', restartRequired: true),
    const EcoConfigItem(key: 'db.max_conns', label: 'DB 최대 연결 수', value: '20', type: 'number', description: '연결 풀 최대 크기'),
    const EcoConfigItem(key: 'db.ssl_mode', label: 'DB SSL 모드', value: 'disable', type: 'select', description: 'PostgreSQL SSL 모드', allowedValues: ['disable', 'require', 'verify-ca', 'verify-full']),
    const EcoConfigItem(key: 'redis.host', label: 'Redis 호스트', value: 'redis', type: 'string', description: 'Redis 캐시 서버 호스트'),
    const EcoConfigItem(key: 'redis.port', label: 'Redis 포트', value: '6379', type: 'number', description: 'Redis 포트'),
    const EcoConfigItem(key: 'redis.password', label: 'Redis 비밀번호', value: '', type: 'secret', description: 'Redis 접속 비밀번호'),
    const EcoConfigItem(key: 'kafka.brokers', label: 'Kafka 브로커', value: 'redpanda:19092', type: 'string', description: 'Kafka/Redpanda 브로커 주소 (쉼표 구분)'),
    const EcoConfigItem(key: 's3.endpoint', label: 'S3 엔드포인트', value: 'minio:9000', type: 'string', description: 'S3/MinIO 스토리지 엔드포인트'),
    const EcoConfigItem(key: 's3.access_key', label: 'S3 Access Key', value: '', type: 'secret', description: 'S3 액세스 키'),
    const EcoConfigItem(key: 's3.secret_key', label: 'S3 Secret Key', value: '', type: 'secret', description: 'S3 시크릿 키'),
    const EcoConfigItem(key: 's3.bucket', label: 'S3 버킷', value: 'manpasik', type: 'string', description: '기본 S3 버킷 이름'),
    const EcoConfigItem(key: 'milvus.host', label: 'Milvus 호스트', value: 'milvus', type: 'string', description: '벡터 DB 호스트'),
    const EcoConfigItem(key: 'milvus.port', label: 'Milvus 포트', value: '19530', type: 'number', description: '벡터 DB 포트'),
    const EcoConfigItem(key: 'es.url', label: 'Elasticsearch URL', value: 'http://elasticsearch:9200', type: 'url', description: '검색 엔진 URL'),
  ];

  // ── 보안 설정 ──
  static final _securityConfigs = [
    const EcoConfigItem(key: 'jwt.secret', label: 'JWT 시크릿', value: '', type: 'secret', description: 'JWT 서명 키 (HS256)', restartRequired: true),
    const EcoConfigItem(key: 'jwt.access_ttl_minutes', label: 'Access Token TTL (분)', value: '15', type: 'number', description: 'Access Token 유효 시간'),
    const EcoConfigItem(key: 'jwt.refresh_ttl_days', label: 'Refresh Token TTL (일)', value: '7', type: 'number', description: 'Refresh Token 유효 기간'),
    const EcoConfigItem(key: 'jwt.issuer', label: 'JWT Issuer', value: 'manpasik-auth', type: 'string', description: 'JWT 발급자 식별자'),
    const EcoConfigItem(key: 'cors.allowed_origins', label: 'CORS 허용 도메인', value: '*', type: 'string', description: '허용할 Origin 목록 (쉼표 구분)'),
    const EcoConfigItem(key: 'rate_limit.requests_per_minute', label: '분당 요청 제한', value: '100', type: 'number', description: 'IP당 분당 최대 요청 수'),
    const EcoConfigItem(key: 'rate_limit.burst', label: '버스트 허용량', value: '20', type: 'number', description: '순간 최대 허용 요청 수'),
    const EcoConfigItem(key: 'security.max_body_size_mb', label: '최대 요청 바디 (MB)', value: '10', type: 'number', description: 'HTTP 요청 바디 크기 제한'),
    const EcoConfigItem(key: 'security.config_encryption_key', label: '설정 암호화 키', value: '', type: 'secret', description: 'AES-256-GCM 설정값 암호화 키', restartRequired: true),
    const EcoConfigItem(key: 'security.bcrypt_cost', label: 'bcrypt 비용', value: '12', type: 'number', description: '비밀번호 해싱 비용 (10~14)'),
  ];

  // ── 알림 채널 ──
  static final _notificationConfigs = [
    const EcoConfigItem(key: 'noti.fcm_enabled', label: 'FCM 푸시 활성화', value: 'true', type: 'boolean', description: 'Firebase Cloud Messaging 활성화'),
    const EcoConfigItem(key: 'noti.email_enabled', label: '이메일 알림 활성화', value: 'false', type: 'boolean', description: 'SMTP 이메일 알림 활성화'),
    const EcoConfigItem(key: 'noti.sms_enabled', label: 'SMS 알림 활성화', value: 'false', type: 'boolean', description: 'SMS 알림 활성화'),
    const EcoConfigItem(key: 'noti.smtp_host', label: 'SMTP 호스트', value: '', type: 'string', description: '이메일 발송 SMTP 서버'),
    const EcoConfigItem(key: 'noti.smtp_port', label: 'SMTP 포트', value: '587', type: 'number', description: 'SMTP 포트 (587: TLS, 465: SSL)'),
    const EcoConfigItem(key: 'noti.smtp_user', label: 'SMTP 사용자', value: '', type: 'string', description: 'SMTP 인증 사용자'),
    const EcoConfigItem(key: 'noti.smtp_password', label: 'SMTP 비밀번호', value: '', type: 'secret', description: 'SMTP 인증 비밀번호'),
    const EcoConfigItem(key: 'noti.sender_email', label: '발신 이메일', value: 'noreply@manpasik.com', type: 'string', description: '알림 발신자 이메일 주소'),
    const EcoConfigItem(key: 'noti.sms_api_key', label: 'SMS API 키', value: '', type: 'secret', description: 'SMS 발송 서비스 API 키'),
  ];

  // ── 측정/AI ──
  static final _aiMeasurementConfigs = [
    const EcoConfigItem(key: 'ai.model_version', label: 'AI 모델 버전', value: 'v2.1.0', type: 'string', description: '현재 배포된 AI 추론 모델 버전'),
    const EcoConfigItem(key: 'ai.inference_timeout_ms', label: 'AI 추론 타임아웃 (ms)', value: '5000', type: 'number', description: 'AI 추론 최대 대기 시간'),
    const EcoConfigItem(key: 'ai.confidence_threshold', label: '신뢰도 임계값', value: '0.7', type: 'number', description: '분석 결과 최소 신뢰도 (0.0~1.0)'),
    const EcoConfigItem(key: 'ai.max_batch_size', label: '최대 배치 크기', value: '32', type: 'number', description: 'AI 배치 추론 최대 크기'),
    const EcoConfigItem(key: 'measurement.glucose_normal_min', label: '혈당 정상 하한', value: '70', type: 'number', description: '공복 혈당 정상 범위 하한 (mg/dL)'),
    const EcoConfigItem(key: 'measurement.glucose_normal_max', label: '혈당 정상 상한', value: '100', type: 'number', description: '공복 혈당 정상 범위 상한 (mg/dL)'),
    const EcoConfigItem(key: 'measurement.hba1c_warning', label: 'HbA1c 주의 기준', value: '6.5', type: 'number', description: '당화혈색소 주의 기준값 (%)'),
    const EcoConfigItem(key: 'measurement.calibration_interval_days', label: '캘리브레이션 주기 (일)', value: '30', type: 'number', description: '카트리지 교정 주기'),
    const EcoConfigItem(key: 'vision.fps_throttle', label: '비전 FPS 제한', value: '5', type: 'number', description: '카메라 프레임 처리 초당 횟수'),
    const EcoConfigItem(key: 'vision.roi_percent', label: 'ROI 비율 (%)', value: '40', type: 'number', description: '카메라 관심 영역 중앙 비율'),
  ];

  // ── 결제 ──
  static final _paymentConfigs = [
    const EcoConfigItem(key: 'pay.pg_provider', label: 'PG 사업자', value: 'toss', type: 'select', description: '결제 게이트웨이 사업자', allowedValues: ['toss', 'iamport', 'nice']),
    const EcoConfigItem(key: 'pay.currency', label: '기본 통화', value: 'KRW', type: 'select', description: '결제 기본 통화', allowedValues: ['KRW', 'USD', 'JPY']),
    const EcoConfigItem(key: 'pay.basic_monthly_price', label: 'Basic 월 요금', value: '9900', type: 'number', description: 'Basic 구독 월 요금 (원)'),
    const EcoConfigItem(key: 'pay.pro_monthly_price', label: 'Pro 월 요금', value: '19900', type: 'number', description: 'Pro 구독 월 요금 (원)'),
    const EcoConfigItem(key: 'pay.clinical_monthly_price', label: 'Clinical 월 요금', value: '49900', type: 'number', description: 'Clinical 구독 월 요금 (원)'),
    const EcoConfigItem(key: 'pay.free_trial_days', label: '무료 체험 기간 (일)', value: '14', type: 'number', description: '신규 가입 무료 체험 기간'),
    const EcoConfigItem(key: 'pay.tax_rate_percent', label: '부가세율 (%)', value: '10', type: 'number', description: '부가가치세 세율'),
    const EcoConfigItem(key: 'pay.refund_policy_days', label: '환불 가능 기간 (일)', value: '7', type: 'number', description: '결제 후 환불 가능 기간'),
  ];

  // ── 환경/운영 ──
  static final _environmentConfigs = [
    const EcoConfigItem(key: 'env.name', label: '환경 이름', value: 'development', type: 'select', description: '현재 배포 환경', allowedValues: ['development', 'staging', 'production']),
    const EcoConfigItem(key: 'env.maintenance_mode', label: '유지보수 모드', value: 'false', type: 'boolean', description: '유지보수 모드 활성화 (모든 API 차단)'),
    const EcoConfigItem(key: 'env.log_level', label: '로그 레벨', value: 'info', type: 'select', description: '전역 로그 출력 레벨', allowedValues: ['debug', 'info', 'warn', 'error']),
    const EcoConfigItem(key: 'env.version', label: '시스템 버전', value: '1.0.0', type: 'string', description: '만파식 생태계 현재 버전'),
    const EcoConfigItem(key: 'env.gateway_port', label: 'Gateway HTTP 포트', value: '8080', type: 'number', description: 'REST API Gateway 리스닝 포트', restartRequired: true),
    const EcoConfigItem(key: 'env.shutdown_timeout_sec', label: '종료 대기 (초)', value: '5', type: 'number', description: 'Graceful Shutdown 대기 시간'),
    const EcoConfigItem(key: 'env.app_name', label: '앱 이름', value: 'MANPASIK', type: 'string', description: '앱 표시 이름 (모든 로케일)'),
    const EcoConfigItem(key: 'env.support_email', label: '지원 이메일', value: 'support@manpasik.com', type: 'string', description: '고객 지원 이메일 주소'),
  ];
}

final ecosystemProvider =
    StateNotifierProvider<EcosystemNotifier, EcosystemState>(
  (ref) => EcosystemNotifier(),
);

// ─────────────────────────────────────────────
// 서비스 목록 데이터
// ─────────────────────────────────────────────

const _allServices = <ServiceInfo>[
  ServiceInfo(name: 'gateway', displayName: 'API Gateway', grpcPort: 8080, category: 'infra'),
  ServiceInfo(name: 'auth-service', displayName: '인증 서비스', grpcPort: 50051, category: 'infra'),
  ServiceInfo(name: 'admin-service', displayName: '관리자 서비스', grpcPort: 50068, category: 'infra'),
  ServiceInfo(name: 'user-service', displayName: '사용자 서비스', grpcPort: 50052, category: 'user'),
  ServiceInfo(name: 'family-service', displayName: '가족 서비스', grpcPort: 50063, category: 'user'),
  ServiceInfo(name: 'community-service', displayName: '커뮤니티 서비스', grpcPort: 50067, category: 'user'),
  ServiceInfo(name: 'device-service', displayName: '디바이스 서비스', grpcPort: 50053, category: 'medical'),
  ServiceInfo(name: 'measurement-service', displayName: '측정 서비스', grpcPort: 50054, category: 'medical'),
  ServiceInfo(name: 'health-record-service', displayName: '건강기록 서비스', grpcPort: 50064, category: 'medical'),
  ServiceInfo(name: 'calibration-service', displayName: '캘리브레이션 서비스', grpcPort: 50060, category: 'medical'),
  ServiceInfo(name: 'ai-inference-service', displayName: 'AI 추론 서비스', grpcPort: 50058, category: 'ai'),
  ServiceInfo(name: 'vision-service', displayName: '비전 AI 서비스', grpcPort: 50071, category: 'ai'),
  ServiceInfo(name: 'cartridge-service', displayName: '카트리지 서비스', grpcPort: 50059, category: 'ai'),
  ServiceInfo(name: 'shop-service', displayName: '상점 서비스', grpcPort: 50056, category: 'commerce'),
  ServiceInfo(name: 'payment-service', displayName: '결제 서비스', grpcPort: 50057, category: 'commerce'),
  ServiceInfo(name: 'subscription-service', displayName: '구독 서비스', grpcPort: 50055, category: 'commerce'),
  ServiceInfo(name: 'marketplace-service', displayName: '마켓플레이스', grpcPort: 50072, category: 'commerce'),
  ServiceInfo(name: 'telemedicine-service', displayName: '원격진료 서비스', grpcPort: 50065, category: 'medical'),
  ServiceInfo(name: 'reservation-service', displayName: '예약 서비스', grpcPort: 50066, category: 'medical'),
  ServiceInfo(name: 'coaching-service', displayName: '코칭 서비스', grpcPort: 50061, category: 'medical'),
  ServiceInfo(name: 'prescription-service', displayName: '처방전 서비스', grpcPort: 50069, category: 'medical'),
  ServiceInfo(name: 'notification-service', displayName: '알림 서비스', grpcPort: 50062, category: 'utility'),
  ServiceInfo(name: 'translation-service', displayName: '번역 서비스', grpcPort: 50070, category: 'utility'),
  ServiceInfo(name: 'video-service', displayName: '비디오 서비스', grpcPort: 50071, category: 'utility'),
  ServiceInfo(name: 'emergency-service', displayName: '응급 서비스', grpcPort: 50073, category: 'utility'),
  ServiceInfo(name: 'analytics-service', displayName: '분석 서비스', grpcPort: 50074, category: 'utility'),
  ServiceInfo(name: 'iot-gateway-service', displayName: 'IoT 게이트웨이', grpcPort: 50075, category: 'utility'),
  ServiceInfo(name: 'nlp-service', displayName: 'NLP 서비스', grpcPort: 50076, category: 'ai'),
];

// ─────────────────────────────────────────────
// 메인 화면
// ─────────────────────────────────────────────

/// 만파식 생태계 총괄 관리 화면
///
/// 28개 마이크로서비스, 외부 연동 키, 기능 플래그,
/// 인프라/보안/알림/AI/결제/환경 설정을 단일 화면에서 관리합니다.
class AdminEcosystemScreen extends ConsumerStatefulWidget {
  const AdminEcosystemScreen({super.key});

  @override
  ConsumerState<AdminEcosystemScreen> createState() =>
      _AdminEcosystemScreenState();
}

class _AdminEcosystemScreenState extends ConsumerState<AdminEcosystemScreen> {
  final _searchController = TextEditingController();
  String _searchQuery = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final ecoState = ref.watch(ecosystemProvider);

    // 저장 성공 메시지
    ref.listen<EcosystemState>(ecosystemProvider, (prev, next) {
      if (next.successMessage != null &&
          next.successMessage != prev?.successMessage) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(next.successMessage!),
            behavior: SnackBarBehavior.floating,
            backgroundColor: Colors.green.shade700,
          ),
        );
        ref.read(ecosystemProvider.notifier).clearMessage();
      }
    });

    return Scaffold(
      appBar: AppBar(
        title: const Text('생태계 총괄 관리'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            tooltip: '카테고리별 설정',
            onPressed: () => context.push('/admin/settings'),
          ),
        ],
      ),
      body: Column(
        children: [
          // 검색 바
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 4),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '설정 검색 (키, 이름, 설명)...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchQuery.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() => _searchQuery = '');
                        },
                      )
                    : null,
                contentPadding: const EdgeInsets.symmetric(vertical: 12),
              ),
              onChanged: (v) => setState(() => _searchQuery = v.toLowerCase()),
            ),
          ),

          // 콘텐츠
          Expanded(
            child: ListView(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 32),
              children: [
                // 1) 생태계 개요
                _buildOverviewSection(theme),
                const SizedBox(height: 16),

                // 2) 서비스 관리
                _buildServiceGrid(theme),
                const SizedBox(height: 16),

                // 3~10) 설정 섹션들
                ..._buildConfigSections(theme, ecoState),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // ── 1. 생태계 개요 ──
  Widget _buildOverviewSection(ThemeData theme) {
    final upCount =
        _allServices.where((s) => s.status == 'up').length;

    return Card(
      color: AppTheme.sanggamGold.withValues(alpha: 0.08),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.hub_outlined, color: AppTheme.sanggamGold, size: 28),
                const SizedBox(width: 10),
                Text('만파식 생태계', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                _overviewChip(theme, '$upCount/${_allServices.length}', '서비스 가동', Icons.cloud_done, Colors.green),
                const SizedBox(width: 12),
                _overviewChip(theme, '177', 'REST 엔드포인트', Icons.api, Colors.blue),
                const SizedBox(width: 12),
                _overviewChip(theme, '25', 'DB 스키마', Icons.storage, Colors.orange),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                _overviewChip(theme, '6', '외부 연동', Icons.extension, Colors.purple),
                const SizedBox(width: 12),
                _overviewChip(theme, '12', '기능 플래그', Icons.flag, Colors.teal),
                const SizedBox(width: 12),
                _overviewChip(theme, '80+', '설정 항목', Icons.tune, Colors.red),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _overviewChip(ThemeData theme, String value, String label, IconData icon, Color color) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 8),
        decoration: BoxDecoration(
          color: color.withValues(alpha: 0.1),
          borderRadius: BorderRadius.circular(10),
        ),
        child: Column(
          children: [
            Icon(icon, size: 18, color: color),
            const SizedBox(height: 4),
            Text(value, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
            Text(label, style: theme.textTheme.bodySmall?.copyWith(fontSize: 10, color: color)),
          ],
        ),
      ),
    );
  }

  // ── 2. 서비스 관리 그리드 ──
  Widget _buildServiceGrid(ThemeData theme) {
    final categories = {
      'infra': ('인프라', Icons.dns, Colors.blue),
      'user': ('사용자', Icons.people, Colors.green),
      'medical': ('의료/측정', Icons.medical_services, Colors.red),
      'ai': ('AI/비전', Icons.psychology, Colors.purple),
      'commerce': ('상거래', Icons.shopping_bag, Colors.orange),
      'utility': ('유틸리티', Icons.build, Colors.teal),
    };

    return _EcoSection(
      title: '마이크로서비스 현황',
      icon: Icons.cloud_outlined,
      subtitle: '${_allServices.length}개 서비스',
      children: categories.entries.map((entry) {
        final (catLabel, catIcon, catColor) = entry.value;
        final services =
            _allServices.where((s) => s.category == entry.key).toList();
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(catIcon, size: 16, color: catColor),
                  const SizedBox(width: 6),
                  Text(catLabel,
                      style: theme.textTheme.labelLarge
                          ?.copyWith(fontWeight: FontWeight.bold, color: catColor)),
                  const Spacer(),
                  Text('${services.length}개',
                      style: theme.textTheme.bodySmall?.copyWith(color: catColor)),
                ],
              ),
              const SizedBox(height: 4),
              Wrap(
                spacing: 6,
                runSpacing: 4,
                children: services.map((svc) {
                  return Chip(
                    avatar: Icon(
                      svc.status == 'up' ? Icons.check_circle : Icons.error,
                      size: 14,
                      color: svc.status == 'up' ? Colors.green : Colors.red,
                    ),
                    label: Text(svc.displayName,
                        style: theme.textTheme.bodySmall),
                    visualDensity: VisualDensity.compact,
                    materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                    padding: EdgeInsets.zero,
                  );
                }).toList(),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  // ── 3~10. 설정 섹션들 ──
  List<Widget> _buildConfigSections(ThemeData theme, EcosystemState ecoState) {
    final sectionMeta = {
      'external_keys': ('외부 연동 키 관리', Icons.vpn_key, '카카오, Toss, Firebase, Keycloak, Google, Apple'),
      'feature_flags': ('기능 플래그', Icons.flag_outlined, '기능별 ON/OFF 토글'),
      'infra': ('인프라 설정', Icons.dns_outlined, 'PostgreSQL, Redis, Kafka, S3, Milvus, ES'),
      'security': ('보안 설정', Icons.shield_outlined, 'JWT, CORS, Rate Limit, 암호화'),
      'notification': ('알림 채널 설정', Icons.notifications_outlined, 'FCM, Email, SMS'),
      'ai_measurement': ('측정 / AI 설정', Icons.science_outlined, 'AI 모델, 임계값, 캘리브레이션, 비전'),
      'payment': ('결제 설정', Icons.payment_outlined, 'PG, 구독 요금, 환불 정책'),
      'environment': ('환경 / 운영', Icons.settings_applications_outlined, '환경, 유지보수, 로그, 버전'),
    };

    final widgets = <Widget>[];
    for (final entry in sectionMeta.entries) {
      final sectionKey = entry.key;
      final (title, icon, subtitle) = entry.value;
      final items = ecoState.sections[sectionKey] ?? [];

      // 검색 필터링
      final filtered = _searchQuery.isEmpty
          ? items
          : items.where((item) =>
              item.key.toLowerCase().contains(_searchQuery) ||
              item.label.toLowerCase().contains(_searchQuery) ||
              item.description.toLowerCase().contains(_searchQuery)).toList();

      if (_searchQuery.isNotEmpty && filtered.isEmpty) continue;

      widgets.add(
        Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: _EcoSection(
            title: title,
            icon: icon,
            subtitle: '${filtered.length}개 항목 · $subtitle',
            children: filtered
                .map((item) => _ConfigTile(
                      item: item,
                      sectionKey: sectionKey,
                    ))
                .toList(),
          ),
        ),
      );
    }
    return widgets;
  }
}

// ─────────────────────────────────────────────
// 재사용 위젯
// ─────────────────────────────────────────────

/// 접을 수 있는 설정 섹션
class _EcoSection extends StatelessWidget {
  const _EcoSection({
    required this.title,
    required this.icon,
    this.subtitle,
    required this.children,
  });

  final String title;
  final IconData icon;
  final String? subtitle;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      clipBehavior: Clip.antiAlias,
      child: ExpansionTile(
        leading: Icon(icon, color: AppTheme.sanggamGold),
        title: Text(title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
        subtitle: subtitle != null
            ? Text(subtitle!, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant))
            : null,
        childrenPadding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
        expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
        children: children,
      ),
    );
  }
}

/// 개별 설정 항목 타일
class _ConfigTile extends ConsumerWidget {
  const _ConfigTile({
    required this.item,
    required this.sectionKey,
  });

  final EcoConfigItem item;
  final String sectionKey;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final isSecret = item.type == 'secret';
    final displayValue = isSecret && item.value.isNotEmpty
        ? '••••••••'
        : item.value.isEmpty
            ? '(미설정)'
            : item.value;

    // boolean 타입은 인라인 스위치
    if (item.type == 'boolean') {
      final isOn = item.value.toLowerCase() == 'true';
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 2),
        child: SwitchListTile(
          contentPadding: const EdgeInsets.symmetric(horizontal: 4),
          title: Text(item.label, style: theme.textTheme.bodyMedium),
          subtitle: item.description.isNotEmpty
              ? Text(item.description,
                  style: theme.textTheme.bodySmall
                      ?.copyWith(color: theme.colorScheme.onSurfaceVariant))
              : null,
          value: isOn,
          dense: true,
          onChanged: (v) {
            ref.read(ecosystemProvider.notifier).updateConfig(
                  sectionKey,
                  item.key,
                  v.toString(),
                );
          },
        ),
      );
    }

    // 기타 타입은 탭하여 편집
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: ListTile(
        contentPadding: const EdgeInsets.symmetric(horizontal: 4),
        dense: true,
        title: Row(
          children: [
            Expanded(
              child: Text(item.label, style: theme.textTheme.bodyMedium),
            ),
            if (item.restartRequired)
              Tooltip(
                message: '변경 시 재시작 필요',
                child: Icon(Icons.restart_alt, size: 14,
                    color: theme.colorScheme.error.withValues(alpha: 0.7)),
              ),
            const SizedBox(width: 4),
            _typeBadge(theme, item.type),
          ],
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 2),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: theme.colorScheme.surfaceContainerHighest.withValues(alpha: 0.5),
                borderRadius: BorderRadius.circular(6),
              ),
              child: Text(
                displayValue,
                style: theme.textTheme.bodySmall?.copyWith(
                  fontFamily: 'monospace',
                  color: item.value.isEmpty
                      ? theme.colorScheme.onSurfaceVariant
                      : null,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (item.description.isNotEmpty) ...[
              const SizedBox(height: 2),
              Text(item.description,
                  style: theme.textTheme.bodySmall?.copyWith(
                    fontSize: 11,
                    color: theme.colorScheme.onSurfaceVariant,
                  )),
            ],
          ],
        ),
        trailing: const Icon(Icons.edit_outlined, size: 18),
        onTap: () => _showEditDialog(context, ref),
      ),
    );
  }

  Widget _typeBadge(ThemeData theme, String type) {
    final (label, color) = switch (type) {
      'secret' => ('비밀', Colors.red),
      'number' => ('숫자', Colors.orange),
      'url' => ('URL', Colors.blue),
      'select' => ('선택', Colors.purple),
      _ => ('문자열', Colors.grey),
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 1),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(label,
          style: TextStyle(fontSize: 9, color: color, fontWeight: FontWeight.w600)),
    );
  }

  Future<void> _showEditDialog(BuildContext context, WidgetRef ref) async {
    final controller = TextEditingController(text: item.value);
    String selectedValue = item.value;

    final result = await showDialog<String>(
      context: context,
      builder: (ctx) {
        return StatefulBuilder(builder: (ctx2, setDialogState) {
          final theme = Theme.of(ctx2);
          return AlertDialog(
            title: Text(item.label),
            content: SizedBox(
              width: MediaQuery.of(ctx2).size.width * 0.85,
              child: SingleChildScrollView(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (item.description.isNotEmpty) ...[
                      Text(item.description,
                          style: theme.textTheme.bodySmall
                              ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                      const SizedBox(height: 12),
                    ],
                    if (item.type == 'select' && item.allowedValues.isNotEmpty)
                      DropdownButtonFormField<String>(
                        value: item.allowedValues.contains(selectedValue)
                            ? selectedValue
                            : null,
                        decoration: const InputDecoration(labelText: '값 선택'),
                        items: item.allowedValues
                            .map((v) => DropdownMenuItem(value: v, child: Text(v)))
                            .toList(),
                        onChanged: (v) {
                          if (v != null) {
                            setDialogState(() => selectedValue = v);
                          }
                        },
                      )
                    else if (item.type == 'number')
                      TextField(
                        controller: controller,
                        keyboardType: const TextInputType.numberWithOptions(decimal: true),
                        inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[\d.\-]'))],
                        decoration: const InputDecoration(labelText: '값'),
                      )
                    else
                      TextField(
                        controller: controller,
                        obscureText: item.type == 'secret',
                        decoration: InputDecoration(
                          labelText: '값',
                          hintText: item.type == 'secret' ? '비밀 값 입력' : '값 입력',
                        ),
                      ),
                    const SizedBox(height: 12),
                    // 메타 정보
                    Container(
                      padding: const EdgeInsets.all(10),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainerHighest
                            .withValues(alpha: 0.5),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('키: ${item.key}',
                              style: theme.textTheme.bodySmall
                                  ?.copyWith(fontFamily: 'monospace', fontSize: 11)),
                          if (item.restartRequired)
                            Padding(
                              padding: const EdgeInsets.only(top: 4),
                              child: Row(
                                children: [
                                  Icon(Icons.warning_amber, size: 13,
                                      color: theme.colorScheme.error),
                                  const SizedBox(width: 4),
                                  Text('변경 후 서비스 재시작이 필요합니다.',
                                      style: theme.textTheme.bodySmall
                                          ?.copyWith(color: theme.colorScheme.error, fontSize: 11)),
                                ],
                              ),
                            ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(ctx2).pop(),
                child: const Text('취소'),
              ),
              FilledButton(
                onPressed: () {
                  final val = item.type == 'select'
                      ? selectedValue
                      : controller.text.trim();
                  Navigator.of(ctx2).pop(val);
                },
                child: const Text('저장'),
              ),
            ],
          );
        });
      },
    );

    controller.dispose();

    if (result != null) {
      ref.read(ecosystemProvider.notifier).updateConfig(
            sectionKey,
            item.key,
            result,
          );
    }
  }
}
