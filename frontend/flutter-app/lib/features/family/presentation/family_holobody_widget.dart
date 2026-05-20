import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/holo_body.dart';

/// HoloBodyBuilder 는 HoloBody 위젯 생성 함수.
///
/// 기본값은 실제 HoloBody — 테스트 시 단순 Container 로 교체 가능
/// (HoloBody 가 webview_flutter 의존이라 테스트 환경에서 빌드 불가).
typedef HoloBodyBuilder = Widget Function({
  required double width,
  required double height,
  required Map<String, dynamic> bioData,
  required Color color,
  Key? key,
});

/// defaultHoloBodyBuilder 는 실제 HoloBody 위젯 사용.
Widget defaultHoloBodyBuilder({
  required double width,
  required double height,
  required Map<String, dynamic> bioData,
  required Color color,
  Key? key,
}) =>
    HoloBody(
      key: key,
      width: width,
      height: height,
      bioData: bioData,
      color: color,
    );

/// FamilyHoloBodyWidget — 가족 그룹 활성 시 멤버 선택 토글이 포함된 HoloBody.
///
/// activeTenantProvider 가 가족 그룹 (isPersonal=false) 일 때 멤버 chip rail 을
/// 표시하고, 선택된 멤버의 bioData 를 HoloBody 에 전달.
///
/// 색상 차별:
///   - 본인 (selectedUserId=null): SanggamTheme.primary (gold)
///   - 가족 멤버: SanggamTheme.jagaeCyan
class FamilyHoloBodyWidget extends ConsumerStatefulWidget {
  const FamilyHoloBodyWidget({
    super.key,
    this.width = 300,
    this.height = 500,
    this.selfBioData = const {},
    this.fetchMemberBioData,
    this.holoBodyBuilder = defaultHoloBodyBuilder,
  });

  final double width;
  final double height;

  /// 본인 측정 데이터 (멤버 미선택 시 사용).
  final Map<String, dynamic> selfBioData;

  /// 멤버 ID → bioData 변환 함수 (선택, 미주입 시 빈 데이터).
  final Future<Map<String, dynamic>> Function(String userId)? fetchMemberBioData;

  /// HoloBody 빌더 — 테스트에서 모킹 가능.
  final HoloBodyBuilder holoBodyBuilder;

  @override
  ConsumerState<FamilyHoloBodyWidget> createState() =>
      _FamilyHoloBodyWidgetState();
}

class _FamilyHoloBodyWidgetState extends ConsumerState<FamilyHoloBodyWidget> {
  /// 선택된 멤버 ID. null = 본인.
  String? _selectedUserId;
  Map<String, dynamic> _selectedBioData = const {};
  bool _loading = false;

  Future<void> _selectMember(String? userId) async {
    setState(() {
      _selectedUserId = userId;
      _loading = userId != null;
      if (userId == null) {
        _selectedBioData = const {};
      }
    });
    if (userId == null) return;

    if (widget.fetchMemberBioData == null) {
      setState(() {
        _selectedBioData = const {};
        _loading = false;
      });
      return;
    }
    try {
      final data = await widget.fetchMemberBioData!(userId);
      if (mounted) {
        setState(() {
          _selectedBioData = data;
          _loading = false;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() {
          _selectedBioData = const {};
          _loading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final activeTenant = ref.watch(activeTenantProvider);
    final isPersonal = activeTenant.isPersonal;

    // 본인 모드면 멤버 선택 UI 숨김
    if (isPersonal) {
      return widget.holoBodyBuilder(
        width: widget.width,
        height: widget.height,
        bioData: widget.selfBioData,
        color: SanggamTheme.primary,
      );
    }

    // 가족 그룹 활성: 멤버 목록 + 토글
    final asyncMembers = ref.watch(tenantMembersProvider(activeTenant.tenantId!));

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        SizedBox(
          height: 56,
          child: asyncMembers.when(
            data: (list) => _MemberRail(
              key: const Key('member-rail'),
              members: list,
              selectedUserId: _selectedUserId,
              onSelect: _selectMember,
            ),
            loading: () => const Center(
              child: SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  color: SanggamTheme.primary,
                ),
              ),
            ),
            error: (_, __) => const SizedBox(),
          ),
        ),
        const SizedBox(height: 8),
        if (_loading)
          SizedBox(
            width: widget.width,
            height: widget.height,
            child: const Center(
              child: CircularProgressIndicator(color: SanggamTheme.primary),
            ),
          )
        else
          widget.holoBodyBuilder(
            key: ValueKey('holo-${_selectedUserId ?? "self"}'),
            width: widget.width,
            height: widget.height,
            bioData: _selectedUserId == null
                ? widget.selfBioData
                : _selectedBioData,
            color: _selectedUserId == null
                ? SanggamTheme.primary
                : SanggamTheme.jagaeCyan,
          ),
      ],
    );
  }
}

class _MemberRail extends StatelessWidget {
  const _MemberRail({
    super.key,
    required this.members,
    required this.selectedUserId,
    required this.onSelect,
  });

  final List<MembershipDto> members;
  final String? selectedUserId;
  final ValueChanged<String?> onSelect;

  @override
  Widget build(BuildContext context) {
    return ListView(
      scrollDirection: Axis.horizontal,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      children: [
        _Chip(
          key: const Key('chip-self'),
          label: '본인',
          icon: Icons.person,
          isSelected: selectedUserId == null,
          color: SanggamTheme.primary,
          onTap: () => onSelect(null),
        ),
        ...members.map((m) => _Chip(
              key: Key('chip-${m.userId}'),
              label: m.userId,
              icon: Icons.group,
              isSelected: selectedUserId == m.userId,
              color: SanggamTheme.jagaeCyan,
              onTap: () => onSelect(m.userId),
            )),
      ],
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip({
    super.key,
    required this.label,
    required this.icon,
    required this.isSelected,
    required this.color,
    required this.onTap,
  });

  final String label;
  final IconData icon;
  final bool isSelected;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
      child: ActionChip(
        avatar: Icon(icon, size: 16, color: isSelected ? color : SanggamTheme.onSurfaceDim),
        label: Text(
          label,
          style: TextStyle(
            color: isSelected ? color : Colors.white,
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
          ),
        ),
        backgroundColor: isSelected
            ? color.withValues(alpha: 0.15)
            : SanggamTheme.surfaceVariant,
        side: BorderSide(
          color: isSelected ? color : SanggamTheme.surfaceVariant,
          width: isSelected ? 2 : 1,
        ),
        onPressed: onTap,
      ),
    );
  }
}
