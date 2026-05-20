import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/providers/grpc_provider.dart';

/// 단일 멤버십 모델 (REST 응답 매핑).
@immutable
class MembershipDto {
  const MembershipDto({
    required this.userId,
    required this.tenantId,
    required this.role,
    required this.active,
    required this.joinedAt,
  });

  final String userId;
  final String tenantId;
  final String role;
  final bool active;
  final int joinedAt;

  /// REST JSON → MembershipDto.
  factory MembershipDto.fromJson(Map<String, dynamic> json) => MembershipDto(
        userId: (json['user_id'] ?? '') as String,
        tenantId: (json['tenant_id'] ?? '') as String,
        role: (json['role'] ?? '') as String,
        active: (json['active'] ?? false) as bool,
        joinedAt: (json['joined_at'] is int)
            ? json['joined_at'] as int
            : int.tryParse(json['joined_at']?.toString() ?? '0') ?? 0,
      );

  @override
  bool operator ==(Object other) =>
      other is MembershipDto &&
      other.userId == userId &&
      other.tenantId == tenantId &&
      other.role == role &&
      other.active == active &&
      other.joinedAt == joinedAt;

  @override
  int get hashCode => Object.hash(userId, tenantId, role, active, joinedAt);
}

/// myMembershipsProvider 는 REST gateway 에서 사용자의 활성 멤버십을 조회.
///
/// 자동 갱신 트리거가 필요하면 `ref.invalidate(myMembershipsProvider)` 호출.
final myMembershipsProvider = FutureProvider<List<MembershipDto>>((ref) async {
  final client = ref.watch(restClientProvider);
  final raw = await client.getMyMemberships();
  return raw
      .map(MembershipDto.fromJson)
      .where((m) => m.active)
      .toList(growable: false);
});

/// tenantMembersProvider 는 특정 조직의 모든 멤버를 조회 (admin 전용).
///
/// 호출 예: `ref.watch(tenantMembersProvider('hospA'))`
final tenantMembersProvider =
    FutureProvider.family<List<MembershipDto>, String>((ref, tenantId) async {
  final client = ref.watch(restClientProvider);
  final raw = await client.listTenantMembers(tenantId);
  return raw.map(MembershipDto.fromJson).toList(growable: false);
});
