import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'package:qr_flutter/qr_flutter.dart';

// ───────────────────────────────────────────────────
// FamilyCreateScreen — Sanggam Orbit 가족 그룹 생성
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 4x → SanggamTheme.primary
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.* 3x → 직접 TextStyle
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] h:12→16
// ───────────────────────────────────────────────────

/// 가족 그룹 생성 화면
class FamilyCreateScreen extends ConsumerStatefulWidget {
  const FamilyCreateScreen({super.key, this.mode});

  final String? mode;

  @override
  ConsumerState<FamilyCreateScreen> createState() =>
      _FamilyCreateScreenState();
}

class _FamilyCreateScreenState
    extends ConsumerState<FamilyCreateScreen>
    with SingleTickerProviderStateMixin {
  final _nameController = TextEditingController();
  final _inviteController = TextEditingController();
  String _inviteMethod = 'link';
  bool _isSubmitting = false;
  late final AnimationController _qrAnimController;
  late final Animation<double> _qrScale;

  bool get _isInviteMode => widget.mode == 'invite';

  @override
  void initState() {
    super.initState();
    _qrAnimController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 500),
    );
    _qrScale = CurvedAnimation(
        parent: _qrAnimController, curve: Curves.elasticOut);
  }

  @override
  void dispose() {
    _qrAnimController.dispose();
    _nameController.dispose();
    _inviteController.dispose();
    super.dispose();
  }

  Future<void> _createGroup() async {
    if (_nameController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('그룹 이름을 입력해주세요.')),
      );
      return;
    }
    setState(() => _isSubmitting = true);
    try {
      final client = ref.read(restClientProvider);
      await client.createFamilyGroup(
          userId: 'current-user',
          name: _nameController.text.trim());
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
              content: Text('가족 그룹이 생성되었습니다.')),
        );
        context.pop();
      }
    } catch (_) {
      if (mounted) {
        setState(() => _isSubmitting = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
              content: Text('그룹 생성에 실패했습니다.')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      _isInviteMode ? '가족 초대' : '가족 그룹 만들기',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment.stretch,
                  children: [
                    if (!_isInviteMode) ...[
                      const Icon(Icons.family_restroom,
                          size: 64,
                          color: SanggamTheme.primary),
                      const SizedBox(height: 16),
                      const Text(
                        '가족 그룹을 만들어\n건강 데이터를 함께 관리하세요.',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                        ),
                      ),
                      const SizedBox(height: 32),
                      TextFormField(
                        controller: _nameController,
                        decoration: const InputDecoration(
                          labelText: '그룹 이름',
                          hintText: '예: 우리 가족',
                          prefixIcon: Icon(Icons.group),
                        ),
                      ),
                      const SizedBox(height: 24),
                    ],
                    const Text(
                      '초대 방법',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    SegmentedButton<String>(
                      segments: const [
                        ButtonSegment(
                            value: 'link',
                            label: Text('초대 링크'),
                            icon: Icon(Icons.link)),
                        ButtonSegment(
                            value: 'email',
                            label: Text('이메일'),
                            icon: Icon(Icons.email)),
                        ButtonSegment(
                            value: 'qr',
                            label: Text('QR 코드'),
                            icon: Icon(Icons.qr_code)),
                      ],
                      selected: {_inviteMethod},
                      onSelectionChanged: (s) {
                        setState(
                            () => _inviteMethod = s.first);
                        if (s.first == 'qr') {
                          _qrAnimController.forward(from: 0);
                        }
                      },
                    ),
                    const SizedBox(height: 16),
                    if (_inviteMethod == 'email')
                      TextFormField(
                        controller: _inviteController,
                        keyboardType:
                            TextInputType.emailAddress,
                        decoration: const InputDecoration(
                          labelText: '초대할 이메일',
                          hintText: 'family@email.com',
                          prefixIcon: Icon(Icons.email),
                        ),
                      ),
                    if (_inviteMethod == 'link')
                      SanggamContainer(
                        borderRadius: 16,
                        padding: EdgeInsets.zero,
                        child: ListTile(
                          leading: const Icon(Icons.link,
                              color: SanggamTheme.primary),
                          title: const Text('초대 링크 생성'),
                          subtitle: const Text(
                              '링크를 복사하여 가족에게 공유하세요.'),
                          trailing: IconButton(
                            icon: const Icon(Icons.copy),
                            tooltip: '링크 복사',
                            onPressed: () {
                              ScaffoldMessenger.of(context)
                                  .showSnackBar(
                                const SnackBar(
                                    content: Text(
                                        '초대 링크가 복사되었습니다.')),
                              );
                            },
                          ),
                        ),
                      ),
                    if (_inviteMethod == 'qr')
                      ScaleTransition(
                        scale: _qrScale,
                        child: SanggamContainer(
                          borderRadius: 16,
                          padding: const EdgeInsets.all(32),
                          child: Column(
                            children: [
                              QrImageView(
                                data:
                                    'https://manpasik.com/invite/${_nameController.text.trim().isNotEmpty ? _nameController.text.trim().hashCode.abs() : "new-group"}',
                                version: QrVersions.auto,
                                size: 160,
                                dataModuleStyle:
                                    const QrDataModuleStyle(
                                  dataModuleShape:
                                      QrDataModuleShape
                                          .square,
                                  color: Colors.white,
                                ),
                                eyeStyle: const QrEyeStyle(
                                  eyeShape:
                                      QrEyeShape.square,
                                  color:
                                      SanggamTheme.primary,
                                ),
                                errorCorrectionLevel:
                                    QrErrorCorrectLevel.M,
                              ),
                              const SizedBox(height: 16),
                              const Text(
                                'QR 코드를 스캔하여 참여하세요',
                                style: TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim,
                                  fontSize: 12,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    const SizedBox(height: 32),
                    if (!_isInviteMode)
                      FilledButton(
                        onPressed:
                            _isSubmitting ? null : _createGroup,
                        style: FilledButton.styleFrom(
                          minimumSize:
                              const Size.fromHeight(48),
                          backgroundColor:
                              SanggamTheme.primary,
                          foregroundColor:
                              SanggamTheme.background,
                          shape: RoundedRectangleBorder(
                              borderRadius:
                                  BorderRadius.circular(16)),
                        ),
                        child: _isSubmitting
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child:
                                    CircularProgressIndicator(
                                        strokeWidth: 2),
                              )
                            : const Text('그룹 만들기'),
                      ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
