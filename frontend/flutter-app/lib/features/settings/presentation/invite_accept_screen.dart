import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// InviteAcceptScreen — 조직 초대 수락 화면.
///
/// 사용자가 받은 토큰을 입력하면 가입 처리되며, 성공 시 myMembershipsProvider
/// 를 invalidate 하여 새로 가입한 조직이 목록에 즉시 반영됨.
class InviteAcceptScreen extends ConsumerStatefulWidget {
  const InviteAcceptScreen({super.key, this.initialToken});

  /// 딥링크/푸시 알림에서 토큰을 미리 채워줄 때 사용.
  final String? initialToken;

  @override
  ConsumerState<InviteAcceptScreen> createState() =>
      _InviteAcceptScreenState();
}

class _InviteAcceptScreenState extends ConsumerState<InviteAcceptScreen> {
  late final TextEditingController _tokenController;
  bool _submitting = false;
  String? _resultMessage;
  bool _success = false;

  @override
  void initState() {
    super.initState();
    _tokenController = TextEditingController(text: widget.initialToken ?? '');
  }

  @override
  void dispose() {
    _tokenController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final token = _tokenController.text.trim();
    if (token.isEmpty) {
      setState(() {
        _resultMessage = '토큰을 입력하세요';
        _success = false;
      });
      return;
    }
    setState(() {
      _submitting = true;
      _resultMessage = null;
    });
    try {
      final client = ref.read(restClientProvider);
      final resp = await client.acceptTenantInvitation(token);
      ref.invalidate(myMembershipsProvider);
      setState(() {
        _success = true;
        final tid = (resp['tenant_id'] ?? '') as String;
        final role = (resp['role'] ?? '') as String;
        _resultMessage = '가입 완료: $tid ($role)';
      });
    } catch (e) {
      setState(() {
        _success = false;
        _resultMessage = e.toString();
      });
    } finally {
      setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            _Header(),
            Expanded(
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  SanggamContainer(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            '받은 초대 토큰을 입력하세요',
                            style: TextStyle(
                                color: SanggamTheme.primary,
                                fontSize: 14,
                                fontWeight: FontWeight.w600),
                          ),
                          const SizedBox(height: 12),
                          TextField(
                            key: const Key('invite-token-field'),
                            controller: _tokenController,
                            style: const TextStyle(
                              color: Colors.white,
                              fontFamily: 'monospace',
                            ),
                            decoration: const InputDecoration(
                              hintText: '32자 헥스 토큰',
                              hintStyle: TextStyle(
                                  color: SanggamTheme.onSurfaceDim),
                              filled: true,
                              fillColor: SanggamTheme.surfaceVariant,
                              border: OutlineInputBorder(
                                  borderSide: BorderSide.none),
                            ),
                          ),
                          const SizedBox(height: 16),
                          SizedBox(
                            width: double.infinity,
                            child: ElevatedButton(
                              key: const Key('submit-accept'),
                              onPressed: _submitting ? null : _submit,
                              style: ElevatedButton.styleFrom(
                                backgroundColor: SanggamTheme.primary,
                                foregroundColor: SanggamTheme.onPrimary,
                                padding: const EdgeInsets.symmetric(
                                    vertical: 14),
                              ),
                              child: _submitting
                                  ? const SizedBox(
                                      height: 20,
                                      width: 20,
                                      child: CircularProgressIndicator(
                                          strokeWidth: 2),
                                    )
                                  : const Text('가입 수락'),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  if (_resultMessage != null) ...[
                    const SizedBox(height: 16),
                    SanggamContainer(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          children: [
                            Icon(
                              _success ? Icons.check_circle : Icons.error_outline,
                              color: _success
                                  ? SanggamTheme.primary
                                  : SanggamTheme.error,
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: Text(
                                _resultMessage!,
                                style: TextStyle(
                                  color: _success
                                      ? Colors.white
                                      : SanggamTheme.error,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
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
            '조직 초대 수락',
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
