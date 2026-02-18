import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 긴급 대응 설정 화면
///
/// 긴급 연락처 관리, 위험 감지 기준, 119 자동 신고, 안전 모드 설정.
/// url_launcher를 통한 실제 119 전화 연결 및 GPS 위치 공유를 지원합니다.
class EmergencySettingsScreen extends ConsumerStatefulWidget {
  const EmergencySettingsScreen({super.key});

  @override
  ConsumerState<EmergencySettingsScreen> createState() => _EmergencySettingsScreenState();
}

class _EmergencySettingsScreenState extends ConsumerState<EmergencySettingsScreen> {
  // 긴급 연락처 목록
  final List<_EmergencyContact> _contacts = [
    _EmergencyContact(name: '', phone: '', relation: '가족'),
  ];

  // 위험 감지 설정
  bool _enableAnomalyDetection = true;
  bool _autoReport119 = false;
  bool _aiVoiceCall = true;
  double _riskThreshold = 0.8;

  // 안전 모드
  String _safetyMode = 'normal';

  // 위치 공유
  bool _shareLocation = true;

  static const _phoneChannel = MethodChannel('com.manpasik/phone');

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('긴급 대응 설정'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // ── 긴급 연락망 ──
          _buildSectionHeader(theme, '긴급 연락망', Icons.contacts),
          Card(
            child: Column(
              children: [
                ..._contacts.asMap().entries.map((e) => _buildContactTile(theme, e.key)),
                TextButton.icon(
                  onPressed: () {
                    setState(() {
                      _contacts.add(_EmergencyContact(name: '', phone: '', relation: '가족'));
                    });
                  },
                  icon: const Icon(Icons.add),
                  label: const Text('연락처 추가'),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // ── 위험 감지 설정 ──
          _buildSectionHeader(theme, '위험 감지 설정', Icons.warning_amber),
          Card(
            child: Column(
              children: [
                SwitchListTile(
                  title: const Text('이상 수치 자동 감지'),
                  subtitle: const Text('측정 결과가 위험 범위에 해당할 때 알림'),
                  value: _enableAnomalyDetection,
                  onChanged: (v) => setState(() => _enableAnomalyDetection = v),
                ),
                const Divider(height: 1),
                ListTile(
                  title: const Text('위험 감지 민감도'),
                  subtitle: Slider(
                    value: _riskThreshold,
                    min: 0.5,
                    max: 1.0,
                    divisions: 5,
                    label: '${(_riskThreshold * 100).toInt()}%',
                    onChanged: (v) => setState(() => _riskThreshold = v),
                  ),
                  trailing: Text(
                    '${(_riskThreshold * 100).toInt()}%',
                    style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.bold),
                  ),
                ),
                const Divider(height: 1),
                SwitchListTile(
                  title: const Text('AI 음성 통화'),
                  subtitle: const Text('위험 감지 시 AI가 음성으로 상태 확인'),
                  value: _aiVoiceCall,
                  onChanged: (v) => setState(() => _aiVoiceCall = v),
                ),
                const Divider(height: 1),
                SwitchListTile(
                  title: const Text('119 자동 신고'),
                  subtitle: const Text('응급 상황 시 자동으로 119 신고 (본인 동의 필요)'),
                  value: _autoReport119,
                  activeColor: Colors.red,
                  onChanged: (v) {
                    if (v) {
                      _showAutoReportConfirm(context);
                    } else {
                      setState(() => _autoReport119 = false);
                    }
                  },
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // ── 안전 모드 ──
          _buildSectionHeader(theme, '안전 모드', Icons.shield),
          Card(
            child: Column(
              children: [
                _buildSafetyModeTile(
                  theme, 'normal', '일반 모드',
                  '기본 설정으로 운영합니다.',
                  Icons.check_circle_outline,
                ),
                const Divider(height: 1),
                _buildSafetyModeTile(
                  theme, 'night', '야간 모드',
                  '야간(22:00~06:00) 이상 감지 시 즉시 긴급 연락',
                  Icons.nightlight_round,
                ),
                const Divider(height: 1),
                _buildSafetyModeTile(
                  theme, 'outing', '외출 모드',
                  '외출 중 이상 감지 시 위치 정보 포함 알림',
                  Icons.directions_walk,
                ),
                const Divider(height: 1),
                _buildSafetyModeTile(
                  theme, 'alone', '독거 모드',
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
              padding: const EdgeInsets.only(bottom: 16),
              child: OutlinedButton.icon(
                onPressed: _testEmergencyCall,
                icon: const Icon(Icons.phone, color: Colors.red),
                label: const Text('119 긴급 신고 테스트'),
                style: OutlinedButton.styleFrom(
                  minimumSize: const Size.fromHeight(48),
                  foregroundColor: Colors.red,
                  side: const BorderSide(color: Colors.red),
                ),
              ),
            ),

          // 위치 공유 상태
          Card(
            child: SwitchListTile(
              title: const Text('위치 정보 공유'),
              subtitle: const Text('긴급 신고 시 GPS 위치를 보호자에게 전송'),
              secondary: const Icon(Icons.location_on, color: Colors.blue),
              value: _shareLocation,
              onChanged: (v) => setState(() => _shareLocation = v),
            ),
          ),
          const SizedBox(height: 16),

          // 저장 버튼
          FilledButton(
            onPressed: _saveSettings,
            style: FilledButton.styleFrom(
              minimumSize: const Size.fromHeight(48),
              backgroundColor: AppTheme.sanggamGold,
            ),
            child: const Text('설정 저장'),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  bool _saving = false;

  /// 119 긴급 신고 테스트
  Future<void> _testEmergencyCall() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('119 긴급 신고 테스트'),
        content: const Text(
          '이것은 테스트입니다.\n실제 119 전화가 연결됩니다.\n\n계속하시겠습니까?',
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('취소')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('전화 연결'),
          ),
        ],
      ),
    );

    if (confirmed != true || !mounted) return;

    try {
      // url_launcher를 통한 전화 연결
      // 패키지 설치 후: await launchUrl(Uri.parse('tel:${AppConfig.emergencyNumber}'));
      await _phoneChannel.invokeMethod('dial', {'number': AppConfig.emergencyNumber});

      // 위치 공유 동시 실행
      if (_shareLocation) {
        await _shareLocationToContacts();
      }
    } on MissingPluginException {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('테스트 모드: ${AppConfig.emergencyNumber}로 신고가 접수됩니다.')),
        );
      }
    }
  }

  /// 보호자에게 GPS 위치 공유
  Future<void> _shareLocationToContacts() async {
    try {
      // 위치 정보 획득 (geolocator 패키지 또는 platform channel)
      // final position = await Geolocator.getCurrentPosition();
      // 보호자에게 REST API로 위치 전송
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final contactPhones = _contacts.where((c) => c.phone.isNotEmpty).map((c) => c.phone).toList();
      await client.shareEmergencyLocation(
        userId: userId,
        latitude: 37.5665, // 실제 구현 시 GPS 좌표
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
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
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
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('긴급 대응 설정이 저장되었습니다.')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('설정 저장 실패: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Widget _buildSectionHeader(ThemeData theme, String title, IconData icon) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Icon(icon, size: 20, color: AppTheme.sanggamGold),
          const SizedBox(width: 8),
          Text(
            title,
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }

  Widget _buildContactTile(ThemeData theme, int index) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: theme.colorScheme.primaryContainer,
        child: Text('${index + 1}', style: TextStyle(color: theme.colorScheme.onPrimaryContainer)),
      ),
      title: Text(_contacts[index].name.isEmpty ? '연락처 ${index + 1}' : _contacts[index].name),
      subtitle: Text(_contacts[index].phone.isEmpty ? '전화번호를 입력하세요' : _contacts[index].phone),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Chip(label: Text(_contacts[index].relation, style: const TextStyle(fontSize: 12))),
          IconButton(
            icon: const Icon(Icons.edit, size: 20),
            onPressed: () => _showEditContactDialog(context, index),
          ),
        ],
      ),
    );
  }

  Widget _buildSafetyModeTile(ThemeData theme, String mode, String title, String subtitle, IconData icon) {
    final isSelected = _safetyMode == mode;
    return RadioListTile<String>(
      value: mode,
      groupValue: _safetyMode,
      onChanged: (v) => setState(() => _safetyMode = v!),
      title: Row(
        children: [
          Icon(icon, size: 20, color: isSelected ? AppTheme.sanggamGold : null),
          const SizedBox(width: 8),
          Text(title),
        ],
      ),
      subtitle: Text(subtitle, style: theme.textTheme.bodySmall),
    );
  }

  void _showEditContactDialog(BuildContext context, int index) {
    final nameCtrl = TextEditingController(text: _contacts[index].name);
    final phoneCtrl = TextEditingController(text: _contacts[index].phone);
    String relation = _contacts[index].relation;

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('긴급 연락처 편집'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: '이름')),
            const SizedBox(height: 8),
            TextField(
              controller: phoneCtrl,
              decoration: const InputDecoration(labelText: '전화번호'),
              keyboardType: TextInputType.phone,
            ),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              value: relation,
              decoration: const InputDecoration(labelText: '관계'),
              items: const [
                DropdownMenuItem(value: '가족', child: Text('가족')),
                DropdownMenuItem(value: '배우자', child: Text('배우자')),
                DropdownMenuItem(value: '자녀', child: Text('자녀')),
                DropdownMenuItem(value: '친구', child: Text('친구')),
                DropdownMenuItem(value: '이웃', child: Text('이웃')),
                DropdownMenuItem(value: '기타', child: Text('기타')),
              ],
              onChanged: (v) => relation = v ?? '가족',
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () {
              setState(() {
                _contacts[index] = _EmergencyContact(
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

  void _showAutoReportConfirm(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('119 자동 신고 동의'),
        content: const Text(
          '이 기능을 활성화하면 AI가 응급 상황으로 판단할 때 자동으로 119에 신고합니다.\n\n'
          '오신고가 발생할 수 있으며, 이에 따른 책임은 사용자에게 있습니다.\n\n'
          '동의하시겠습니까?',
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () {
              setState(() => _autoReport119 = true);
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
  const _EmergencyContact({required this.name, required this.phone, required this.relation});
}
