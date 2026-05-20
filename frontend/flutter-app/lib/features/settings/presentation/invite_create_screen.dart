import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// InviteCreateScreen — 조직 초대 발급 화면.
///
/// admin/owner 권한이 있는 조직에 대해 새 멤버 초대 token 을 발급하고,
/// 발급된 token 을 표시하여 외부 채널(이메일/카톡 등)로 공유 가능.
class InviteCreateScreen extends ConsumerStatefulWidget {
  const InviteCreateScreen({super.key});

  @override
  ConsumerState<InviteCreateScreen> createState() => _InviteCreateScreenState();
}

class _InviteCreateScreenState extends ConsumerState<InviteCreateScreen> {
  String? _selectedTenantId;
  String _selectedRole = 'member';
  final _hintController = TextEditingController();
  bool _submitting = false;
  String? _generatedToken;
  String? _error;

  static const _roles = <String>['member', 'medical_staff', 'admin', 'viewer'];

  @override
  void dispose() {
    _hintController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (_selectedTenantId == null || _selectedTenantId!.isEmpty) {
      setState(() => _error = '조직을 선택하세요');
      return;
    }
    setState(() {
      _submitting = true;
      _error = null;
      _generatedToken = null;
    });
    try {
      final client = ref.read(restClientProvider);
      final resp = await client.createTenantInvitation(
        tenantId: _selectedTenantId!,
        role: _selectedRole,
        inviteeHint: _hintController.text.trim().isEmpty
            ? null
            : _hintController.text.trim(),
      );
      setState(() {
        _generatedToken = resp['token'] as String?;
      });
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final asyncMems = ref.watch(myMembershipsProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _Header(),
            Expanded(
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  asyncMems.when(
                    data: (list) {
                      // admin/owner 인 멤버십만 초대 가능
                      final adminTenants = list.where((m) =>
                          m.role == 'admin' ||
                          m.role == 'owner').toList();
                      if (adminTenants.isEmpty) {
                        return const _EmptyAdminState();
                      }
                      return _InviteForm(
                        adminTenants: adminTenants,
                        selectedTenantId: _selectedTenantId,
                        onTenantChanged: (v) =>
                            setState(() => _selectedTenantId = v),
                        roles: _roles,
                        selectedRole: _selectedRole,
                        onRoleChanged: (v) =>
                            setState(() => _selectedRole = v ?? 'member'),
                        hintController: _hintController,
                        submitting: _submitting,
                        onSubmit: _submit,
                      );
                    },
                    loading: () => const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary),
                    ),
                    error: (e, _) => Text(
                      '멤버십 조회 실패: $e',
                      style: const TextStyle(color: SanggamTheme.error),
                    ),
                  ),
                  if (_error != null) ...[
                    const SizedBox(height: 12),
                    SanggamContainer(
                      child: Padding(
                        padding: const EdgeInsets.all(12),
                        child: Text(_error!,
                            style: const TextStyle(color: SanggamTheme.error)),
                      ),
                    ),
                  ],
                  if (_generatedToken != null) ...[
                    const SizedBox(height: 16),
                    _TokenCard(token: _generatedToken!),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Header extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Row(
        children: [
          IconButton(
            icon: const Icon(Icons.arrow_back, color: SanggamTheme.primary),
            onPressed: () => Navigator.of(context).maybePop(),
          ),
          const SizedBox(width: 8),
          const Text(
            '조직 초대 발급',
            style: TextStyle(
              color: SanggamTheme.primary,
              fontSize: 20,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
}

class _EmptyAdminState extends StatelessWidget {
  const _EmptyAdminState();

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(32),
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: const [
            Icon(Icons.lock_outline,
                color: SanggamTheme.onSurfaceDim, size: 48),
            SizedBox(height: 12),
            Text(
              '관리자 권한이 있는 조직이 없습니다',
              style: TextStyle(color: SanggamTheme.onSurfaceDim),
            ),
          ],
        ),
      ),
    );
  }
}

class _InviteForm extends StatelessWidget {
  const _InviteForm({
    required this.adminTenants,
    required this.selectedTenantId,
    required this.onTenantChanged,
    required this.roles,
    required this.selectedRole,
    required this.onRoleChanged,
    required this.hintController,
    required this.submitting,
    required this.onSubmit,
  });

  final List<MembershipDto> adminTenants;
  final String? selectedTenantId;
  final ValueChanged<String?> onTenantChanged;
  final List<String> roles;
  final String selectedRole;
  final ValueChanged<String?> onRoleChanged;
  final TextEditingController hintController;
  final bool submitting;
  final VoidCallback onSubmit;

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const _Label('조직 선택'),
            DropdownButtonFormField<String>(
              key: const Key('tenant-dropdown'),
              dropdownColor: SanggamTheme.surface,
              style: const TextStyle(color: Colors.white),
              value: selectedTenantId,
              items: adminTenants
                  .map((m) => DropdownMenuItem(
                        value: m.tenantId,
                        child: Text('${m.tenantId} (${m.role})'),
                      ))
                  .toList(),
              onChanged: onTenantChanged,
              decoration: const InputDecoration(
                filled: true,
                fillColor: SanggamTheme.surfaceVariant,
                border: OutlineInputBorder(borderSide: BorderSide.none),
              ),
            ),
            const SizedBox(height: 16),
            const _Label('초대 역할'),
            DropdownButtonFormField<String>(
              key: const Key('role-dropdown'),
              dropdownColor: SanggamTheme.surface,
              style: const TextStyle(color: Colors.white),
              value: selectedRole,
              items: roles
                  .map((r) =>
                      DropdownMenuItem(value: r, child: Text(r)))
                  .toList(),
              onChanged: onRoleChanged,
              decoration: const InputDecoration(
                filled: true,
                fillColor: SanggamTheme.surfaceVariant,
                border: OutlineInputBorder(borderSide: BorderSide.none),
              ),
            ),
            const SizedBox(height: 16),
            const _Label('수락자 식별 (이메일/전화, 선택)'),
            TextField(
              key: const Key('hint-field'),
              controller: hintController,
              style: const TextStyle(color: Colors.white),
              decoration: const InputDecoration(
                hintText: 'user@example.com',
                hintStyle: TextStyle(color: SanggamTheme.onSurfaceDim),
                filled: true,
                fillColor: SanggamTheme.surfaceVariant,
                border: OutlineInputBorder(borderSide: BorderSide.none),
              ),
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                key: const Key('submit-invite'),
                onPressed: submitting ? null : onSubmit,
                style: ElevatedButton.styleFrom(
                  backgroundColor: SanggamTheme.primary,
                  foregroundColor: SanggamTheme.onPrimary,
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child: submitting
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('초대 토큰 발급'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Label extends StatelessWidget {
  const _Label(this.text);
  final String text;
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Text(text,
          style: const TextStyle(
              color: SanggamTheme.onSurfaceDim, fontSize: 12)),
    );
  }
}

class _TokenCard extends StatelessWidget {
  const _TokenCard({required this.token});
  final String token;

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '초대 토큰 발급 완료',
              style: TextStyle(
                  color: SanggamTheme.primary,
                  fontSize: 14,
                  fontWeight: FontWeight.w600),
            ),
            const SizedBox(height: 8),
            const Text(
              '아래 토큰을 수락자에게 전달하세요. 7일 후 만료됩니다.',
              style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 11),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: SanggamTheme.surfaceVariant,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  Expanded(
                    child: SelectableText(
                      token,
                      style: const TextStyle(
                          color: Colors.white,
                          fontFamily: 'monospace',
                          fontSize: 12),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.copy, color: SanggamTheme.primary),
                    onPressed: () {
                      Clipboard.setData(ClipboardData(text: token));
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('클립보드에 복사됨')),
                      );
                    },
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
