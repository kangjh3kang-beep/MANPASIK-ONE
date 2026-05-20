import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// TenantInterceptor는 모든 REST 요청에 활성 조직 ID(`X-Tenant-ID`)와
/// 사용자 ID(`X-User-ID`)를 자동으로 부착합니다.
///
/// 백엔드 gateway 의 `TenantPropagation` 미들웨어가 이 헤더를 읽어
/// gRPC outgoing metadata 로 전파하고, telemedicine/health-record/family/auth
/// 등 인터셉터가 활성화된 서비스에서 멀티테넌트 격리 검증에 사용합니다.
///
/// 헤더 미설정 시 (예: 로그인 전 또는 단일테넌트 모드) 부착하지 않음.
class TenantInterceptor extends Interceptor {
  TenantInterceptor({this.preferencesProvider});

  /// 테스트용 — SharedPreferences 인스턴스 주입.
  /// nil 이면 SharedPreferences.getInstance() 사용.
  final Future<SharedPreferences> Function()? preferencesProvider;

  static const _activeTenantKey = 'active_tenant_id';
  static const _activeUserKey = 'active_user_id';

  /// SharedPreferences 에 활성 조직 ID 저장.
  ///
  /// 사용자가 조직 전환 시 호출 (예: 가족 그룹 전환, 병원 선택).
  /// `tenantId == null` 또는 빈 문자열이면 저장된 값 제거 (개인 모드).
  static Future<void> setActiveTenant(String? tenantId) async {
    final prefs = await SharedPreferences.getInstance();
    if (tenantId == null || tenantId.isEmpty) {
      await prefs.remove(_activeTenantKey);
    } else {
      await prefs.setString(_activeTenantKey, tenantId);
    }
  }

  /// 현재 저장된 활성 조직 ID 반환. 없으면 null.
  static Future<String?> getActiveTenant() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_activeTenantKey);
  }

  /// 활성 사용자 ID 저장 (보통 로그인 시 호출).
  static Future<void> setActiveUser(String? userId) async {
    final prefs = await SharedPreferences.getInstance();
    if (userId == null || userId.isEmpty) {
      await prefs.remove(_activeUserKey);
    } else {
      await prefs.setString(_activeUserKey, userId);
    }
  }

  /// 현재 저장된 사용자 ID 반환.
  static Future<String?> getActiveUser() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_activeUserKey);
  }

  /// 모든 테넌시 데이터 제거 (로그아웃 시 호출).
  static Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_activeTenantKey);
    await prefs.remove(_activeUserKey);
  }

  Future<SharedPreferences> _prefs() async {
    if (preferencesProvider != null) {
      return preferencesProvider!();
    }
    return SharedPreferences.getInstance();
  }

  @override
  void onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    late final SharedPreferences prefs;
    try {
      prefs = await _prefs();
    } catch (_) {
      super.onRequest(options, handler);
      return;
    }
    final tenantId = prefs.getString(_activeTenantKey);
    if (tenantId != null && tenantId.isNotEmpty) {
      options.headers['X-Tenant-ID'] = tenantId;
    }
    final userId = prefs.getString(_activeUserKey);
    if (userId != null && userId.isNotEmpty) {
      options.headers['X-User-ID'] = userId;
    }
    super.onRequest(options, handler);
  }
}
