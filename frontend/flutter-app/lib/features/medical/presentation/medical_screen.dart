import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/utils/navigation_utils.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

class MedicalScreen extends StatefulWidget {
  const MedicalScreen({super.key});

  @override
  State<MedicalScreen> createState() => _MedicalScreenState();
}

class _MedicalScreenState extends State<MedicalScreen> {
  bool _shareVitals = true;
  bool _shareLabs = true;
  bool _shareRecentHistory = true;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(8, 8, 8, 0),
              child: Row(
                children: [
                  IconButton(
                    tooltip: '뒤로',
                    onPressed: () => context.popOrHome(),
                    icon: const Icon(Icons.arrow_back_ios_new_rounded,
                        color: Colors.white),
                  ),
                  const SizedBox(width: 6),
                  const Expanded(
                    child: Text(
                      '의료 서비스 허브',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
                  TextButton.icon(
                    onPressed: () => context.go('/home'),
                    icon: const Icon(Icons.home_rounded, size: 16),
                    label: const Text('홈'),
                    style: TextButton.styleFrom(
                      foregroundColor: SanggamTheme.primary,
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  _buildJourneyOverview(),
                  const SizedBox(height: 16),
                  _buildActionGrid(context),
                  const SizedBox(height: 16),
                  _buildLiveEntryCard(context),
                  const SizedBox(height: 16),
                  _buildRecentConsultation(context),
                  const SizedBox(height: 16),
                  _buildSupportMenu(context),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildJourneyOverview() {
    return SanggamContainer(
      borderRadius: 20,
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.health_and_safety_rounded,
                  size: 20,
                  color: SanggamTheme.primary.withValues(alpha: 0.95)),
              const SizedBox(width: 8),
              const Text(
                '진료 스토리라인',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
          const SizedBox(height: 10),
          const Text(
            '의료기관 선택 → 예약 → 화상진료 → 처방전/약국 전송 → 결과 확인',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 13,
              height: 1.45,
            ),
          ),
          const SizedBox(height: 12),
          const Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _JourneyChip(icon: Icons.search_rounded, label: '기관 선택'),
              _JourneyChip(icon: Icons.event_available_rounded, label: '예약'),
              _JourneyChip(icon: Icons.videocam_rounded, label: '화상진료'),
              _JourneyChip(icon: Icons.local_pharmacy_rounded, label: '약국 전송'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildActionGrid(BuildContext context) {
    return GridView.count(
      crossAxisCount: 2,
      crossAxisSpacing: 10,
      mainAxisSpacing: 10,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      childAspectRatio: 1.28,
      children: [
        _ActionCard(
          icon: Icons.apartment_rounded,
          title: '의료기관 선택',
          description: '전문의/병원 검색',
          onTap: () => context.push('/medical/facility-search'),
        ),
        _ActionCard(
          icon: Icons.video_camera_front_rounded,
          title: '화상진료 예약',
          description: '진료과·의사 선택',
          onTap: () => context.push('/medical/telemedicine'),
        ),
        _ActionCard(
          icon: Icons.local_pharmacy_rounded,
          title: '약국 선택',
          description: '처방전 전송·수령',
          onTap: () => context.push('/medical/pharmacy'),
        ),
        _ActionCard(
          icon: Icons.share_rounded,
          title: '바이오데이터 공유',
          description: '진료용 데이터 제공',
          onTap: _openShareSheet,
        ),
      ],
    );
  }

  Widget _buildLiveEntryCard(BuildContext context) {
    return SanggamContainer(
      borderRadius: 18,
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Container(
            width: 46,
            height: 46,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: SanggamTheme.primary.withValues(alpha: 0.16),
            ),
            child:
                const Icon(Icons.videocam_rounded, color: SanggamTheme.primary),
          ),
          const SizedBox(width: 12),
          const Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '진료실 바로 입장',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w700,
                    fontSize: 15,
                  ),
                ),
                SizedBox(height: 4),
                Text(
                  '예약 세션이 있으면 바로 화상진료로 이동합니다.',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          FilledButton(
            onPressed: () => context.push('/medical/video-call/demo-session'),
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.primary,
              foregroundColor: SanggamTheme.background,
            ),
            child: const Text('입장'),
          ),
        ],
      ),
    );
  }

  Widget _buildRecentConsultation(BuildContext context) {
    return SanggamContainer(
      borderRadius: 18,
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '최근 진료/처방',
            style: TextStyle(
              color: Colors.white,
              fontWeight: FontWeight.w700,
              fontSize: 15,
            ),
          ),
          const SizedBox(height: 8),
          ListTile(
            contentPadding: EdgeInsets.zero,
            dense: true,
            leading: const Icon(Icons.description_outlined,
                color: SanggamTheme.jagaeCyan),
            title: const Text(
              '진료 결과 보기',
              style: TextStyle(color: Colors.white, fontSize: 14),
            ),
            subtitle: const Text(
              '2026-02-23 / 내과 / 김건강 전문의',
              style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 11),
            ),
            trailing:
                const Icon(Icons.chevron_right_rounded, color: Colors.white54),
            onTap: () => context.push('/medical/consultation/demo/result'),
          ),
          ListTile(
            contentPadding: EdgeInsets.zero,
            dense: true,
            leading: const Icon(Icons.receipt_long_rounded,
                color: SanggamTheme.primary),
            title: const Text(
              '처방전 상세',
              style: TextStyle(color: Colors.white, fontSize: 14),
            ),
            subtitle: const Text(
              '약국 전송 상태 확인',
              style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 11),
            ),
            trailing:
                const Icon(Icons.chevron_right_rounded, color: Colors.white54),
            onTap: () => context.push('/medical/prescription/demo'),
          ),
        ],
      ),
    );
  }

  Widget _buildSupportMenu(BuildContext context) {
    return SanggamContainer(
      borderRadius: 18,
      padding: const EdgeInsets.all(12),
      child: Column(
        children: [
          _SupportMenuTile(
            icon: Icons.place_rounded,
            title: '내 주변 의료기관',
            subtitle: '거리/운영시간/접수 상태 확인',
            onTap: () => context.push('/medical/facility-search'),
          ),
          const Divider(height: 1, color: Color(0x22D4AF37)),
          _SupportMenuTile(
            icon: Icons.local_shipping_rounded,
            title: '약 배송/수령 설정',
            subtitle: '수령 방식과 알림 관리',
            onTap: () => context.push('/medical/pharmacy'),
          ),
          const Divider(height: 1, color: Color(0x22D4AF37)),
          _SupportMenuTile(
            icon: Icons.monitor_heart_rounded,
            title: '진료 전 바이오데이터 점검',
            subtitle: '공유 항목과 최근 측정값 확인',
            onTap: _openShareSheet,
          ),
        ],
      ),
    );
  }

  Future<void> _openShareSheet() async {
    await showModalBottomSheet<void>(
      context: context,
      backgroundColor: const Color(0xFF0A1428),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (sheetContext) {
        return StatefulBuilder(
          builder: (context, setSheetState) {
            return SafeArea(
              top: false,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '진료용 바이오데이터 공유',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    const SizedBox(height: 10),
                    _shareSwitchTile(
                      title: '실시간 활력징후',
                      value: _shareVitals,
                      onChanged: (value) =>
                          setSheetState(() => _shareVitals = value),
                    ),
                    _shareSwitchTile(
                      title: '최근 검사 결과',
                      value: _shareLabs,
                      onChanged: (value) =>
                          setSheetState(() => _shareLabs = value),
                    ),
                    _shareSwitchTile(
                      title: '최근 진료 이력',
                      value: _shareRecentHistory,
                      onChanged: (value) =>
                          setSheetState(() => _shareRecentHistory = value),
                    ),
                    const SizedBox(height: 12),
                    SizedBox(
                      width: double.infinity,
                      child: FilledButton(
                        onPressed: () {
                          Navigator.of(sheetContext).pop();
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('선택한 바이오데이터 공유 설정이 적용되었습니다.'),
                            ),
                          );
                        },
                        style: FilledButton.styleFrom(
                          backgroundColor: SanggamTheme.primary,
                          foregroundColor: SanggamTheme.background,
                        ),
                        child: const Text('공유 설정 저장'),
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
    setState(() {});
  }

  Widget _shareSwitchTile({
    required String title,
    required bool value,
    required ValueChanged<bool> onChanged,
  }) {
    return SwitchListTile(
      value: value,
      onChanged: onChanged,
      contentPadding: EdgeInsets.zero,
      activeThumbColor: SanggamTheme.primary,
      title: Text(
        title,
        style: const TextStyle(color: Colors.white, fontSize: 14),
      ),
      subtitle: const Text(
        '진료 세션 중 의료진에게만 제공됩니다.',
        style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 11),
      ),
    );
  }
}

class _JourneyChip extends StatelessWidget {
  const _JourneyChip({
    required this.icon,
    required this.label,
  });

  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        color: SanggamTheme.primary.withValues(alpha: 0.12),
        border: Border.all(
          color: SanggamTheme.primary.withValues(alpha: 0.35),
          width: 0.7,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: SanggamTheme.primary),
          const SizedBox(width: 5),
          Text(
            label,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 11,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
}

class _ActionCard extends StatelessWidget {
  const _ActionCard({
    required this.icon,
    required this.title,
    required this.description,
    required this.onTap,
  });

  final IconData icon;
  final String title;
  final String description;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      onTap: onTap,
      borderRadius: 16,
      padding: const EdgeInsets.all(14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: SanggamTheme.primary, size: 22),
          const Spacer(),
          Text(
            title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            description,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 11,
            ),
          ),
        ],
      ),
    );
  }
}

class _SupportMenuTile extends StatelessWidget {
  const _SupportMenuTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: onTap,
      dense: true,
      leading: Icon(icon, color: SanggamTheme.primary),
      title: Text(
        title,
        style: const TextStyle(
          color: Colors.white,
          fontWeight: FontWeight.w600,
          fontSize: 14,
        ),
      ),
      subtitle: Text(
        subtitle,
        style: const TextStyle(
          color: SanggamTheme.onSurfaceDim,
          fontSize: 11,
        ),
      ),
      trailing: const Icon(Icons.chevron_right_rounded, color: Colors.white54),
    );
  }
}
