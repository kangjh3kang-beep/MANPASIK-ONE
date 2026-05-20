import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/network/tenant_interceptor.dart';

/// 활성 조직 정보 (조직 ID + 표시 이름).
@immutable
class ActiveTenant {
  const ActiveTenant({this.tenantId, this.displayName});

  final String? tenantId;
  final String? displayName;

  bool get isPersonal => tenantId == null || tenantId!.isEmpty;

  ActiveTenant copyWith({String? tenantId, String? displayName}) =>
      ActiveTenant(
        tenantId: tenantId ?? this.tenantId,
        displayName: displayName ?? this.displayName,
      );

  @override
  bool operator ==(Object other) =>
      other is ActiveTenant &&
      other.tenantId == tenantId &&
      other.displayName == displayName;

  @override
  int get hashCode => Object.hash(tenantId, displayName);
}

/// ActiveTenantNotifier 는 SharedPreferences 의 활성 조직 ID 와
/// Riverpod 상태를 동기화.
///
/// 조직 전환 시:
///   1. SharedPreferences 업데이트 (TenantInterceptor 가 즉시 헤더에 반영)
///   2. Riverpod 상태 갱신 (UI 리빌드)
class ActiveTenantNotifier extends StateNotifier<ActiveTenant> {
  ActiveTenantNotifier() : super(const ActiveTenant()) {
    _loadFromPreferences();
  }

  Future<void> _loadFromPreferences() async {
    final tid = await TenantInterceptor.getActiveTenant();
    if (tid != null && tid.isNotEmpty) {
      state = ActiveTenant(tenantId: tid);
    }
  }

  /// 활성 조직 변경. tenantId=null 이면 개인 모드(헤더 미부착).
  Future<void> switchTo({String? tenantId, String? displayName}) async {
    await TenantInterceptor.setActiveTenant(tenantId);
    state = ActiveTenant(tenantId: tenantId, displayName: displayName);
  }

  /// 개인 모드로 전환 (조직 헤더 제거).
  Future<void> clear() async {
    await TenantInterceptor.setActiveTenant(null);
    state = const ActiveTenant();
  }
}

/// activeTenantProvider 는 앱 전역에서 현재 활성 조직 상태를 노출.
///
/// 조직 전환 UI 에서 `ref.read(activeTenantProvider.notifier).switchTo(...)` 호출.
/// 다른 화면에서 `ref.watch(activeTenantProvider)` 로 변경 시 리빌드.
final activeTenantProvider =
    StateNotifierProvider<ActiveTenantNotifier, ActiveTenant>(
  (_) => ActiveTenantNotifier(),
);
