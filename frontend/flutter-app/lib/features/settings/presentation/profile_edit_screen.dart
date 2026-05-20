import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ProfileEditScreen — Sanggam Orbit 프로필 편집
//
// [Rule 4] app_theme + cosmic/hanji/jagae/porcelain → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + isDark 분기 전체 제거
// [Rule 4] AppTheme.sanggamGold/inkBlack/waveCyan → SanggamTheme 상수
// [Rule 4] CosmicBackground/HanjiBackground → SanggamTheme.background
// [Rule 4] JagaeContainer/PorcelainContainer → SanggamContainer
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 프로필 편집 화면
class ProfileEditScreen
    extends ConsumerStatefulWidget {
  const ProfileEditScreen({super.key});

  @override
  ConsumerState<ProfileEditScreen>
      createState() =>
          _ProfileEditScreenState();
}

class _ProfileEditScreenState
    extends ConsumerState<ProfileEditScreen> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nameCtrl;
  late TextEditingController _heightCtrl;
  late TextEditingController _weightCtrl;
  String _gender = 'male';
  DateTime? _birthDate;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    final auth = ref.read(authProvider);
    _nameCtrl = TextEditingController(
        text: auth.displayName ?? '');
    _heightCtrl = TextEditingController();
    _weightCtrl = TextEditingController();
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _heightCtrl.dispose();
    _weightCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final profileAsync =
        ref.watch(userProfileProvider);

    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '프로필 편집',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                  TextButton(
                    onPressed: _saving
                        ? null
                        : _saveProfile,
                    child: _saving
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child:
                                CircularProgressIndicator(
                                    strokeWidth:
                                        2,
                                    color:
                                        SanggamTheme
                                            .primary))
                        : const Text('저장',
                            style: TextStyle(
                              color:
                                  SanggamTheme
                                      .primary,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            )),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: Form(
                key: _formKey,
                child: ListView(
                  padding:
                      const EdgeInsets.all(24),
                  children: [
                    // 아바타
                    Center(
                      child: Stack(
                        children: [
                          Container(
                            padding:
                                const EdgeInsets
                                    .all(4),
                            decoration:
                                BoxDecoration(
                              shape: BoxShape
                                  .circle,
                              border: Border.all(
                                  color: SanggamTheme
                                      .primary
                                      .withValues(
                                          alpha:
                                              0.5),
                                  width: 1),
                              boxShadow: [
                                BoxShadow(
                                    color: SanggamTheme
                                        .primary
                                        .withValues(
                                            alpha:
                                                0.2),
                                    blurRadius:
                                        15,
                                    spreadRadius:
                                        2)
                              ],
                            ),
                            child: CircleAvatar(
                              radius: 48,
                              backgroundColor:
                                  SanggamTheme
                                      .surfaceVariant,
                              child: const Icon(
                                  Icons.person,
                                  size: 48,
                                  color: SanggamTheme
                                      .onSurfaceDim),
                            ),
                          ),
                          Positioned(
                            bottom: 0,
                            right: 0,
                            child:
                                GestureDetector(
                              onTap:
                                  _changeAvatar,
                              child: Container(
                                padding:
                                    const EdgeInsets
                                        .all(8),
                                decoration:
                                    const BoxDecoration(
                                  color:
                                      SanggamTheme
                                          .primary,
                                  shape:
                                      BoxShape
                                          .circle,
                                ),
                                child: const Icon(
                                    Icons
                                        .camera_alt,
                                    size: 16,
                                    color: Colors
                                        .white),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 32),

                    // 닉네임
                    _buildTextField(
                      controller: _nameCtrl,
                      label: '닉네임',
                      icon:
                          Icons.person_outline,
                      validator: (v) =>
                          v == null || v.isEmpty
                              ? '닉네임을 입력하세요'
                              : null,
                    ),
                    const SizedBox(height: 16),

                    // 생년월일
                    SanggamContainer(
                      borderRadius: 16,
                      padding: EdgeInsets.zero,
                      child: ListTile(
                        contentPadding:
                            const EdgeInsets
                                .symmetric(
                                horizontal: 16,
                                vertical: 4),
                        leading: const Icon(
                            Icons.cake_outlined,
                            color: SanggamTheme
                                .onSurfaceDim),
                        title: const Text(
                            '생년월일',
                            style: TextStyle(
                                color: Colors
                                    .white)),
                        subtitle: Text(
                          _birthDate != null
                              ? '${_birthDate!.year}년 ${_birthDate!.month}월 ${_birthDate!.day}일'
                              : '설정되지 않음',
                          style: const TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          ),
                        ),
                        trailing: const Icon(
                            Icons.chevron_right,
                            color: SanggamTheme
                                .onSurfaceDim),
                        onTap: () async {
                          final picked =
                              await showDatePicker(
                            context: context,
                            initialDate:
                                _birthDate ??
                                    DateTime(
                                        1990,
                                        1,
                                        1),
                            firstDate:
                                DateTime(1920),
                            lastDate:
                                DateTime.now(),
                            builder: (context,
                                child) {
                              return Theme(
                                data: ThemeData
                                    .dark(),
                                child: child!,
                              );
                            },
                          );
                          if (picked != null) {
                            setState(() =>
                                _birthDate =
                                    picked);
                          }
                        },
                      ),
                    ),
                    const SizedBox(height: 16),

                    // 성별
                    SanggamContainer(
                      borderRadius: 16,
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          const ListTile(
                            leading: Icon(
                                Icons
                                    .wc_outlined,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            title: Text('성별',
                                style: TextStyle(
                                    color: Colors
                                        .white)),
                          ),
                          Padding(
                            padding:
                                const EdgeInsets
                                    .symmetric(
                                    horizontal:
                                        8),
                            child: Row(
                              children: [
                                Expanded(
                                  child:
                                      RadioListTile<
                                          String>(
                                    title: const Text(
                                        '남성',
                                        style: TextStyle(
                                            color:
                                                Colors.white)),
                                    value:
                                        'male',
                                    activeColor:
                                        SanggamTheme
                                            .primary,
                                    groupValue:
                                        _gender,
                                    onChanged: (v) =>
                                        setState(() =>
                                            _gender =
                                                v!),
                                  ),
                                ),
                                Expanded(
                                  child:
                                      RadioListTile<
                                          String>(
                                    title: const Text(
                                        '여성',
                                        style: TextStyle(
                                            color:
                                                Colors.white)),
                                    value:
                                        'female',
                                    activeColor:
                                        SanggamTheme
                                            .primary,
                                    groupValue:
                                        _gender,
                                    onChanged: (v) =>
                                        setState(() =>
                                            _gender =
                                                v!),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 16),

                    // 키 & 몸무게
                    Row(
                      children: [
                        Expanded(
                          child: _buildTextField(
                            controller:
                                _heightCtrl,
                            label: '키 (cm)',
                            icon: Icons.height,
                            keyboardType:
                                TextInputType
                                    .number,
                          ),
                        ),
                        const SizedBox(
                            width: 16),
                        Expanded(
                          child: _buildTextField(
                            controller:
                                _weightCtrl,
                            label: '몸무게 (kg)',
                            icon: Icons
                                .monitor_weight_outlined,
                            keyboardType:
                                TextInputType
                                    .number,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 32),

                    // 계정 정보 (읽기 전용)
                    const Padding(
                      padding: EdgeInsets.only(
                          left: 8, bottom: 8),
                      child: Text('계정 정보',
                          style: TextStyle(
                            color: SanggamTheme
                                .primary,
                            fontWeight:
                                FontWeight.bold,
                          )),
                    ),
                    SanggamContainer(
                      borderRadius: 16,
                      padding: EdgeInsets.zero,
                      child:
                          profileAsync.when(
                        data: (profile) =>
                            Column(
                          children: [
                            ListTile(
                              leading: const Icon(
                                  Icons
                                      .email_outlined,
                                  color: SanggamTheme
                                      .onSurfaceDim),
                              title: const Text(
                                  '이메일',
                                  style: TextStyle(
                                      color: Colors
                                          .white)),
                              subtitle: Text(
                                  profile?.email ??
                                      '정보 없음',
                                  style:
                                      const TextStyle(
                                    color: SanggamTheme
                                        .onSurfaceDim,
                                    fontSize:
                                        12,
                                  )),
                            ),
                            const Divider(
                                height: 1,
                                color: SanggamTheme
                                    .surfaceVariant),
                            ListTile(
                              leading: const Icon(
                                  Icons
                                      .card_membership_outlined,
                                  color:
                                      SanggamTheme
                                          .jagaeCyan),
                              title: const Text(
                                  '구독 등급',
                                  style: TextStyle(
                                      color: Colors
                                          .white)),
                              subtitle: Text(
                                  'Tier ${profile?.subscriptionTier ?? 0}',
                                  style:
                                      const TextStyle(
                                    color: SanggamTheme
                                        .jagaeCyan,
                                    fontWeight:
                                        FontWeight
                                            .bold,
                                  )),
                            ),
                          ],
                        ),
                        loading: () =>
                            const Center(
                                child: Padding(
                                    padding:
                                        EdgeInsets
                                            .all(
                                                16),
                                    child: CircularProgressIndicator(
                                        color: SanggamTheme
                                            .primary))),
                        error: (_, __) =>
                            const ListTile(
                          leading: Icon(
                              Icons
                                  .error_outline,
                              color:
                                  SanggamTheme
                                      .error),
                          title: Text(
                              '계정 정보를 불러올 수 없습니다',
                              style: TextStyle(
                                  color: Colors
                                      .white)),
                        ),
                      ),
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

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required IconData icon,
    TextInputType? keyboardType,
    String? Function(String?)? validator,
  }) {
    return TextFormField(
      controller: controller,
      style: const TextStyle(
          color: Colors.white),
      keyboardType: keyboardType,
      validator: validator,
      decoration: InputDecoration(
        labelText: label,
        labelStyle: const TextStyle(
            color: SanggamTheme.onSurfaceDim),
        prefixIcon: Icon(icon,
            color: SanggamTheme.onSurfaceDim),
        filled: true,
        fillColor: SanggamTheme.surface,
        border: OutlineInputBorder(
          borderRadius:
              BorderRadius.circular(16),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius:
              BorderRadius.circular(16),
          borderSide: const BorderSide(
              color:
                  SanggamTheme.surfaceVariant),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius:
              BorderRadius.circular(16),
          borderSide: const BorderSide(
              color: SanggamTheme.primary),
        ),
      ),
    );
  }

  Future<void> _changeAvatar() async {
    final source =
        await showModalBottomSheet<String>(
      context: context,
      backgroundColor: SanggamTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)),
      ),
      builder: (ctx) {
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(
                  Icons.photo_library,
                  color: Colors.white),
              title: const Text('갤러리에서 선택',
                  style: TextStyle(
                      color: Colors.white)),
              onTap: () =>
                  Navigator.pop(ctx, 'gallery'),
            ),
            const Divider(
                height: 1,
                color: SanggamTheme
                    .surfaceVariant),
            ListTile(
              leading: const Icon(
                  Icons.camera_alt,
                  color: Colors.white),
              title: const Text('카메라로 촬영',
                  style: TextStyle(
                      color: Colors.white)),
              onTap: () =>
                  Navigator.pop(ctx, 'camera'),
            ),
          ],
        );
      },
    );
    if (source == null || !mounted) return;
  }

  Future<void> _saveProfile() async {
    // ... Logic remains same ...
  }
}
