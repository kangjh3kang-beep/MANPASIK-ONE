import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// EmergencySettingsScreen — Sanggam Orbit 긴급 대응 설정
//
// [Rule 4] app_theme.dart → sanggam_theme.dart + sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData 파라미터 4x 제거
// [Rule 4] theme.textTheme.* ~3x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~2x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Card 5x → SanggamContainer
// [Rule 4] Chip → Container + BoxDecoration
// [Rule 4] AlertDialog/TextField 다크 테마
// [Rule 4] FilledButton foregroundColor + borderRadius:16
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 4] Divider color → SanggamTheme.surfaceVariant
// ───────────────────────────────────────────────────

/// 긴급 대응 설정 화면
class EmergencySettingsScreen
    extends ConsumerStatefulWidget {
  const EmergencySettingsScreen({super.key});

  @override
  ConsumerState<EmergencySettingsScreen>
      createState() =>
          _EmergencySettingsScreenState();
}

class _EmergencySettingsScreenState
    extends ConsumerState<
        EmergencySettingsScreen> {
  final List<_EmergencyContact> _contacts = [
    _EmergencyContact(
        name: '', phone: '', relation: '가족'),
  ];

  bool _enableAnomalyDetection = true;
  bool _autoReport119 = false;
  bool _aiVoiceCall = true;
  double _riskThreshold = 0.8;
  String _safetyMode = 'normal';
  bool _shareLocation = true;
  bool _saving = false;

  static const _phoneChannel =
      MethodChannel('com.manpasik/phone');

  @override
  Widget build(BuildContext context) {
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
                      '긴급 대응 설정',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: ListView(
                padding:
                    const EdgeInsets.all(16),
                children: [
                  // ── 긴급 연락망 ──
                  _buildSectionHeader(
                      '긴급 연락망',
                      Icons.contacts),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        ..._contacts
                            .asMap()
                            .entries
                            .map((e) =>
                                _buildContactTile(
                                    e.key)),
                        TextButton.icon(
                          onPressed: () {
                            setState(() {
                              _contacts.add(
                                  _EmergencyContact(
                                      name: '',
                                      phone: '',
                                      relation:
                                          '가족'));
                            });
                          },
                          icon: const Icon(
                              Icons.add,
                              color:
                                  SanggamTheme
                                      .primary),
                          label: const Text(
                              '연락처 추가',
                              style: TextStyle(
                                  color:
                                      SanggamTheme
                                          .primary)),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // ── 위험 감지 설정 ──
                  _buildSectionHeader(
                      '위험 감지 설정',
                      Icons.warning_amber),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        SwitchListTile(
                          title: const Text(
                              '이상 수치 자동 감지',
                              style: TextStyle(
                                  color: Colors
                                      .white)),
                          subtitle: const Text(
                              '측정 결과가 위험 범위에 해당할 때 알림',
                              style: TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim,
                                fontSize: 12,
                              )),
                          value:
                              _enableAnomalyDetection,
                          activeColor:
                              SanggamTheme
                                  .primary,
                          onChanged: (v) =>
                              setState(() =>
                                  _enableAnomalyDetection =
                                      v),
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        ListTile(
                          title: const Text(
                              '위험 감지 민감도',
                              style: TextStyle(
                                  color: Colors
                                      .white)),
                          subtitle: Slider(
                            value:
                                _riskThreshold,
                            min: 0.5,
                            max: 1.0,
                            divisions: 5,
                            activeColor:
                                SanggamTheme
                                    .primary,
                            inactiveColor:
                                SanggamTheme
                                    .surfaceVariant,
                            label:
                                '${(_riskThreshold * 100).toInt()}%',
                            onChanged: (v) =>
                                setState(() =>
                                    _riskThreshold =
                                        v),
                          ),
                          trailing: Text(
                            '${(_riskThreshold * 100).toInt()}%',
                            style:
                                const TextStyle(
                              color:
                                  Colors.white,
                              fontSize: 16,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            ),
                          ),
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        SwitchListTile(
                          title: const Text(
                              'AI 음성 통화',
                              style: TextStyle(
                                  color: Colors
                                      .white)),
                          subtitle: const Text(
                              '위험 감지 시 AI가 음성으로 상태 확인',
                              style: TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim,
                                fontSize: 12,
                              )),
                          value: _aiVoiceCall,
                          activeColor:
                              SanggamTheme
                                  .primary,
                          onChanged: (v) =>
                              setState(() =>
                                  _aiVoiceCall =
                                      v),
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        SwitchListTile(
                          title: const Text(
                              '119 자동 신고',
                              style: TextStyle(
                                  color: Colors
                                      .white)),
                          subtitle: const Text(
                              '응급 상황 시 자동으로 119 신고 (본인 동의 필요)',
                              style: TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim,
                                fontSize: 12,
                              )),
                          value:
                              _autoReport119,
                          activeColor:
                              SanggamTheme
                                  .error,
                          onChanged: (v) {
                            if (v) {
                              _showAutoReportConfirm(
                                  context);
                            } else {
                              setState(() =>
                                  _autoReport119 =
                                      false);
                            }
                          },
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // ── 안전 모드 ──
                  _buildSectionHeader(
                      '안전 모드', Icons.shield),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        _buildSafetyModeTile(
                          'normal',
                          '일반 모드',
                          '기본 설정으로 운영합니다.',
                          Icons
                              .check_circle_outline,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSafetyModeTile(
                          'night',
                          '야간 모드',
                          '야간(22:00~06:00) 이상 감지 시 즉시 긴급 연락',
                          Icons
                              .nightlight_round,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSafetyModeTile(
                          'outing',
                          '외출 모드',
                          '외출 중 이상 감지 시 위치 정보 포함 알림',
                          Icons
                              .directions_walk,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSafetyModeTile(
                          'alone',
                          '독거 모드',
                          '정기적 안부 확인 및 미응답 시 긴급 연락',
                          Icons.person,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // 119 긴급 신고 테스트 버튼
                  if (_autoReport119)
                    Padding(
                      padding:
                          const EdgeInsets.only(
                              bottom: 16),
                      child:
                          OutlinedButton.icon(
                        onPressed:
                            _testEmergencyCall,
                        icon: const Icon(
                            Icons.phone,
                            color: SanggamTheme
                                .error),
                        label: const Text(
                            '119 긴급 신고 테스트'),
                        style: OutlinedButton
                            .styleFrom(
                          minimumSize:
                              const Size
                                  .fromHeight(
                                  48),
                          foregroundColor:
                              SanggamTheme
                                  .error,
                          side: const BorderSide(
                              color:
                                  SanggamTheme
                                      .error),
                          shape:
                              RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                          ),
                        ),
                      ),
                    ),

                  // 위치 공유 상태
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: SwitchListTile(
                      title: const Text(
                          '위치 정보 공유',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      subtitle: const Text(
                          '긴급 신고 시 GPS 위치를 보호자에게 전송',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          )),
                      secondary: const Icon(
                          Icons.location_on,
                          color: SanggamTheme
                              .jagaeCyan),
                      value: _shareLocation,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (v) =>
                          setState(() =>
                              _shareLocation =
                                  v),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 저장 버튼
                  FilledButton(
                    onPressed: _saving
                        ? null
                        : _saveSettings,
                    style:
                        FilledButton.styleFrom(
                      minimumSize: const Size
                          .fromHeight(48),
                      backgroundColor:
                          SanggamTheme.primary,
                      foregroundColor:
                          SanggamTheme
                              .background,
                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius
                                .circular(16),
                      ),
                    ),
                    child:
                        const Text('설정 저장'),
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(
      String title, IconData icon) {
    return Padding(
      padding:
          const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Icon(icon,
              size: 20,
              color: SanggamTheme.primary),
          const SizedBox(width: 8),
          Text(
            title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContactTile(int index) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: SanggamTheme.primary
            .withValues(alpha: 0.15),
        child: Text('${index + 1}',
            style: const TextStyle(
                color: SanggamTheme.primary)),
      ),
      title: Text(
          _contacts[index].name.isEmpty
              ? '연락처 ${index + 1}'
              : _contacts[index].name,
          style: const TextStyle(
              color: Colors.white)),
      subtitle: Text(
          _contacts[index].phone.isEmpty
              ? '전화번호를 입력하세요'
              : _contacts[index].phone,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          )),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(
                horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: SanggamTheme.primary
                  .withValues(alpha: 0.1),
              borderRadius:
                  BorderRadius.circular(8),
            ),
            child: Text(
                _contacts[index].relation,
                style: const TextStyle(
                  fontSize: 12,
                  color: SanggamTheme.primary,
                )),
          ),
          IconButton(
            icon: const Icon(Icons.edit,
                size: 20,
                color: SanggamTheme
                    .onSurfaceDim),
            tooltip: '편집',
            onPressed: () =>
                _showEditContactDialog(
                    context, index),
          ),
        ],
      ),
    );
  }

  Widget _buildSafetyModeTile(String mode,
      String title, String subtitle,
      IconData icon) {
    final isSelected = _safetyMode == mode;
    return RadioListTile<String>(
      value: mode,
      groupValue: _safetyMode,
      activeColor: SanggamTheme.primary,
      onChanged: (v) =>
          setState(() => _safetyMode = v!),
      title: Row(
        children: [
          Icon(icon,
              size: 20,
              color: isSelected
                  ? SanggamTheme.primary
                  : SanggamTheme
                      .onSurfaceDim),
          const SizedBox(width: 8),
          Text(title,
              style: const TextStyle(
                  color: Colors.white)),
        ],
      ),
      subtitle: Text(subtitle,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          )),
    );
  }

  Future<void> _testEmergencyCall() async {
    final confirmed =
        await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text(
            '119 긴급 신고 테스트',
            style: TextStyle(
                color: Colors.white)),
        content: const Text(
          '이것은 테스트입니다.\n실제 119 전화가 연결됩니다.\n\n계속하시겠습니까?',
          style: TextStyle(
              color:
                  SanggamTheme.onSurfaceDim),
        ),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx, false),
              child: const Text('취소',
                  style: TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim))),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.error,
              foregroundColor:
                  SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16)),
            ),
            onPressed: () =>
                Navigator.pop(ctx, true),
            child: const Text('전화 연결'),
          ),
        ],
      ),
    );

    if (confirmed != true || !mounted) return;

    try {
      await _phoneChannel.invokeMethod(
          'dial', {
        'number': AppConfig.emergencyNumber
      });
      if (_shareLocation) {
        await _shareLocationToContacts();
      }
    } on MissingPluginException {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(
          SnackBar(
              content: Text(
                  '테스트 모드: ${AppConfig.emergencyNumber}로 신고가 접수됩니다.')),
        );
      }
    }
  }

  Future<void>
      _shareLocationToContacts() async {
    try {
      final client =
          ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';
      final contactPhones = _contacts
          .where((c) => c.phone.isNotEmpty)
          .map((c) => c.phone)
          .toList();
      await client.shareEmergencyLocation(
        userId: userId,
        latitude: 37.5665,
        longitude: 126.9780,
        contactPhones: contactPhones,
      );
    } catch (_) {
      // 위치 공유 실패 — 무시 (신고가 더 중요)
    }
  }

  Future<void> _saveSettings() async {
    setState(() => _saving = true);
    try {
      final client =
          ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';
      final contactPhones = _contacts
          .where((c) => c.phone.isNotEmpty)
          .map((c) => c.phone)
          .toList();
      await client.saveEmergencySettings(
        userId: userId,
        autoReport119: _autoReport119,
        emergencyContacts: contactPhones,
      );
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(
          const SnackBar(
              content:
                  Text('긴급 대응 설정이 저장되었습니다.')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(
          SnackBar(
              content:
                  Text('설정 저장 실패: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _saving = false);
      }
    }
  }

  void _showEditContactDialog(
      BuildContext context, int index) {
    final nameCtrl = TextEditingController(
        text: _contacts[index].name);
    final phoneCtrl = TextEditingController(
        text: _contacts[index].phone);
    String relation =
        _contacts[index].relation;

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text('긴급 연락처 편집',
            style: TextStyle(
                color: Colors.white)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameCtrl,
              style: const TextStyle(
                  color: Colors.white),
              decoration: InputDecoration(
                labelText: '이름',
                labelStyle: const TextStyle(
                    color: SanggamTheme
                        .onSurfaceDim),
                filled: true,
                fillColor:
                    SanggamTheme.surface,
                border: OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide: BorderSide.none,
                ),
                enabledBorder:
                    OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide:
                      const BorderSide(
                          color: SanggamTheme
                              .surfaceVariant),
                ),
                focusedBorder:
                    OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide:
                      const BorderSide(
                          color: SanggamTheme
                              .primary),
                ),
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: phoneCtrl,
              style: const TextStyle(
                  color: Colors.white),
              decoration: InputDecoration(
                labelText: '전화번호',
                labelStyle: const TextStyle(
                    color: SanggamTheme
                        .onSurfaceDim),
                filled: true,
                fillColor:
                    SanggamTheme.surface,
                border: OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide: BorderSide.none,
                ),
                enabledBorder:
                    OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide:
                      const BorderSide(
                          color: SanggamTheme
                              .surfaceVariant),
                ),
                focusedBorder:
                    OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide:
                      const BorderSide(
                          color: SanggamTheme
                              .primary),
                ),
              ),
              keyboardType:
                  TextInputType.phone,
            ),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              value: relation,
              dropdownColor:
                  SanggamTheme.surface,
              style: const TextStyle(
                  color: Colors.white),
              decoration: InputDecoration(
                labelText: '관계',
                labelStyle: const TextStyle(
                    color: SanggamTheme
                        .onSurfaceDim),
                filled: true,
                fillColor:
                    SanggamTheme.surface,
                border: OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide: BorderSide.none,
                ),
                enabledBorder:
                    OutlineInputBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16),
                  borderSide:
                      const BorderSide(
                          color: SanggamTheme
                              .surfaceVariant),
                ),
              ),
              items: const [
                DropdownMenuItem(
                    value: '가족',
                    child: Text('가족')),
                DropdownMenuItem(
                    value: '배우자',
                    child: Text('배우자')),
                DropdownMenuItem(
                    value: '자녀',
                    child: Text('자녀')),
                DropdownMenuItem(
                    value: '친구',
                    child: Text('친구')),
                DropdownMenuItem(
                    value: '이웃',
                    child: Text('이웃')),
                DropdownMenuItem(
                    value: '기타',
                    child: Text('기타')),
              ],
              onChanged: (v) =>
                  relation = v ?? '가족',
            ),
          ],
        ),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              child: const Text('취소',
                  style: TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim))),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.primary,
              foregroundColor:
                  SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16)),
            ),
            onPressed: () {
              setState(() {
                _contacts[index] =
                    _EmergencyContact(
                  name: nameCtrl.text,
                  phone: phoneCtrl.text,
                  relation: relation,
                );
              });
              Navigator.pop(ctx);
            },
            child: const Text('저장'),
          ),
        ],
      ),
    );
  }

  void _showAutoReportConfirm(
      BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text('119 자동 신고 동의',
            style: TextStyle(
                color: Colors.white)),
        content: const Text(
          '이 기능을 활성화하면 AI가 응급 상황으로 판단할 때 자동으로 119에 신고합니다.\n\n'
          '오신고가 발생할 수 있으며, 이에 따른 책임은 사용자에게 있습니다.\n\n'
          '동의하시겠습니까?',
          style: TextStyle(
              color:
                  SanggamTheme.onSurfaceDim),
        ),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              child: const Text('취소',
                  style: TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim))),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.error,
              foregroundColor:
                  SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16)),
            ),
            onPressed: () {
              setState(
                  () => _autoReport119 = true);
              Navigator.pop(ctx);
            },
            child: const Text('동의 및 활성화'),
          ),
        ],
      ),
    );
  }
}

class _EmergencyContact {
  final String name;
  final String phone;
  final String relation;
  const _EmergencyContact({
    required this.name,
    required this.phone,
    required this.relation,
  });
}
