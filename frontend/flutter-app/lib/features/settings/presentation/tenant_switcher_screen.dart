import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// TenantSwitcherScreen — 활성 조직 전환 화면.
///
/// 사용자가 속한 조직 목록을 보여주고, 선택 시 [activeTenantProvider] 갱신.
/// REST 호출은 자동으로 X-Tenant-ID 헤더가 부착되며 백엔드 멀티테넌트 격리에
/// 활용됨.
///
/// 데이터 소스 우선순위:
///   1. [tenants] 주입 (테스트/직접 제공) — null 이 아니면 이를 사용
///   2. [myMembershipsProvider] 자동 로드 (REST gateway 호출)
class TenantSwitcherScreen extends ConsumerWidget {
  const TenantSwitcherScreen({super.key, this.tenants});

  /// 사용자가 속한 조직 목록 (조직 ID + 표시 이름 + 역할).
  ///
  /// null 이면 [myMembershipsProvider] 로 자동 조회.
  final List<TenantOption>? tenants;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final active = ref.watch(activeTenantProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            _Header(activeTenantId: active.tenantId),
            Expanded(child: _buildBody(context, ref, active)),
          ],
        ),
      ),
    );
  }

  Widget _buildBody(BuildContext context, WidgetRef ref, ActiveTenant active) {
    // 외부 주입이 우선
    if (tenants != null) {
      return _TenantList(
        tenants: tenants!,
        activeTenantId: active.tenantId,
        onSelect: (option) => _switch(ref, option),
        onClearToPersonal: () => _switchToPersonal(ref),
      );
    }

    // 자동 로드 (REST 호출)
    final asyncMems = ref.watch(myMembershipsProvider);
    return asyncMems.when(
      data: (list) {
        final options = list
            .map((m) => TenantOption(
                  tenantId: m.tenantId,
                  displayName: m.tenantId,
                  role: m.role,
                ))
            .toList();
        return _TenantList(
          tenants: options,
          activeTenantId: active.tenantId,
          onSelect: (option) => _switch(ref, option),
          onClearToPersonal: () => _switchToPersonal(ref),
        );
      },
      loading: () => const Center(
        child: CircularProgressIndicator(color: SanggamTheme.primary),
      ),
      error: (err, _) => _ErrorState(
        message: '$err',
        onRetry: () => ref.invalidate(myMembershipsProvider),
      ),
    );
  }

  void _switch(WidgetRef ref, TenantOption opt) {
    ref.read(activeTenantProvider.notifier).switchTo(
          tenantId: opt.tenantId,
          displayName: opt.displayName,
        );
  }

  void _switchToPersonal(WidgetRef ref) {
    ref.read(activeTenantProvider.notifier).clear();
  }
}

/// TenantOption — 조직 표시 옵션.
class TenantOption {
  const TenantOption({
    required this.tenantId,
    required this.displayName,
    required this.role,
  });
  final String tenantId;
  final String displayName;
  final String role;
}

class _Header extends StatelessWidget {
  const _Header({required this.activeTenantId});

  final String? activeTenantId;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
      child: Row(
        children: [
          IconButton(
            icon: const Icon(Icons.arrow_back, color: SanggamTheme.primary),
            onPressed: () => Navigator.of(context).maybePop(),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '조직 전환',
              style: const TextStyle(
                color: SanggamTheme.primary,
                fontSize: 20,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _TenantList extends StatelessWidget {
  const _TenantList({
    required this.tenants,
    required this.activeTenantId,
    required this.onSelect,
    required this.onClearToPersonal,
  });

  final List<TenantOption> tenants;
  final String? activeTenantId;
  final ValueChanged<TenantOption> onSelect;
  final VoidCallback onClearToPersonal;

  @override
  Widget build(BuildContext context) {
    final isPersonal = activeTenantId == null || activeTenantId!.isEmpty;
    return ListView(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      children: [
        // 개인 모드 (조직 미선택)
        SanggamContainer(
          margin: const EdgeInsets.only(bottom: 8),
          child: _PersonalTile(
            isActive: isPersonal,
            onTap: isPersonal ? null : onClearToPersonal,
          ),
        ),
        const SizedBox(height: 16),
        if (tenants.isEmpty)
          const _EmptyState()
        else
          ...tenants.map(
            (t) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: SanggamContainer(
                child: _TenantTile(
                  option: t,
                  isActive: t.tenantId == activeTenantId,
                  onTap: () => onSelect(t),
                ),
              ),
            ),
          ),
      ],
    );
  }
}

class _PersonalTile extends StatelessWidget {
  const _PersonalTile({required this.isActive, this.onTap});
  final bool isActive;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            const Icon(Icons.person, color: SanggamTheme.primary, size: 32),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '개인',
                    style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '조직 헤더 없이 사용 (개인 데이터)',
                    style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                  ),
                ],
              ),
            ),
            if (isActive)
              const Icon(Icons.check_circle, color: SanggamTheme.primary),
          ],
        ),
      ),
    );
  }
}

class _TenantTile extends StatelessWidget {
  const _TenantTile({
    required this.option,
    required this.isActive,
    required this.onTap,
  });

  final TenantOption option;
  final bool isActive;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            const Icon(
              Icons.business,
              color: SanggamTheme.primary,
              size: 32,
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    option.displayName,
                    style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '역할: ${option.role}',
                    style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                  ),
                ],
              ),
            ),
            if (isActive)
              const Icon(Icons.check_circle, color: SanggamTheme.primary),
          ],
        ),
      ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  const _EmptyState();

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 48),
      child: Center(
        child: Text(
          '소속된 조직이 없습니다',
          style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
        ),
      ),
    );
  }
}

class _ErrorState extends StatelessWidget {
  const _ErrorState({required this.message, required this.onRetry});

  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline,
                color: SanggamTheme.error, size: 48),
            const SizedBox(height: 16),
            Text(
              '멤버십 조회 실패',
              style: const TextStyle(
                color: SanggamTheme.error,
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              message,
              textAlign: TextAlign.center,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              ),
            ),
            const SizedBox(height: 16),
            TextButton(
              onPressed: onRetry,
              child: const Text(
                '다시 시도',
                style: TextStyle(color: SanggamTheme.primary),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
