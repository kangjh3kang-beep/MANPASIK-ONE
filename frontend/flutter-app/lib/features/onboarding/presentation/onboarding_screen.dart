import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// OnboardingScreen — Sanggam Orbit 온보딩
//
// [Rule 4] AppTheme.* 10건 → SanggamTheme.*
// [Rule 4] theme.colorScheme.* 8건, theme.textTheme.* ~15건 → 직접
// [Rule 4] withOpacity 2건 → withValues(alpha:)
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] FilledButton 4건 → PrimaryButton
// [Rule 4] Card → SanggamContainer
// [Rule 4] CosmicBackground 추가
// [Rule 2] 12→16, 6→8, 14→16, borderRadius 12→16
// ───────────────────────────────────────────────────

/// 온보딩 완료 상태 Provider (메모리 — 추후 SharedPreferences 연동)
final onboardingCompletedProvider = StateProvider<bool>((ref) => false);

/// 온보딩 화면 — 4단계 페이지뷰
///
/// 1. 환영 (앱 소개)
/// 2. 건강 프로필 (나이, 목표)
/// 3. 디바이스 설정 (BLE 페어링)
/// 4. 완료
class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  final _pageController = PageController();
  int _currentPage = 0;
  static const _totalPages = 4;

  String _selectedGoal = 'general';
  int _age = 30;

  bool _isScanning = false;
  List<DeviceInfoDto> _foundDevices = [];
  String? _connectedDeviceId;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _nextPage() {
    if (_currentPage < _totalPages - 1) {
      _pageController.nextPage(
        duration: const Duration(milliseconds: 350),
        curve: Curves.easeInOut,
      );
    }
  }

  void _completeOnboarding() {
    ref.read(onboardingCompletedProvider.notifier).state = true;
    context.go('/home');
  }

  Future<void> _scanDevices() async {
    setState(() => _isScanning = true);
    final devices = await RustBridge.bleScan();
    if (!mounted) return;
    setState(() {
      _foundDevices = devices;
      _isScanning = false;
    });
  }

  Future<void> _connectDevice(String deviceId) async {
    final ok = await RustBridge.bleConnect(deviceId);
    if (!mounted) return;
    if (ok) {
      setState(() => _connectedDeviceId = deviceId);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: CosmicBackground(
        child: SafeArea(
          child: Column(
            children: [
              // 상단 프로그레스 + 건너뛰기
              Padding(
                padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
                child: Row(
                  children: [
                    Expanded(
                      child: _ProgressBar(
                        current: _currentPage,
                        total: _totalPages,
                      ),
                    ),
                    TextButton(
                      onPressed: _completeOnboarding,
                      child: const Text(
                        '건너뛰기',
                        style: TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

              // 페이지 뷰
              Expanded(
                child: PageView(
                  controller: _pageController,
                  physics: const NeverScrollableScrollPhysics(),
                  onPageChanged: (i) =>
                      setState(() => _currentPage = i),
                  children: [
                    _WelcomePage(onNext: _nextPage),
                    _HealthProfilePage(
                      selectedGoal: _selectedGoal,
                      age: _age,
                      onGoalChanged: (g) =>
                          setState(() => _selectedGoal = g),
                      onAgeChanged: (a) =>
                          setState(() => _age = a),
                      onNext: _nextPage,
                    ),
                    _DeviceSetupPage(
                      isScanning: _isScanning,
                      devices: _foundDevices,
                      connectedDeviceId: _connectedDeviceId,
                      onScan: _scanDevices,
                      onConnect: _connectDevice,
                      onNext: _nextPage,
                    ),
                    _CompletePage(onComplete: _completeOnboarding),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ── 프로그레스 바 ──

class _ProgressBar extends StatelessWidget {
  final int current;
  final int total;
  const _ProgressBar({required this.current, required this.total});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: List.generate(total, (i) {
        return Expanded(
          child: Container(
            height: 4,
            margin: const EdgeInsets.symmetric(horizontal: 2),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(2),
              color: i <= current
                  ? SanggamTheme.primary
                  : SanggamTheme.onSurfaceDim.withValues(alpha: 0.3),
            ),
          ),
        );
      }),
    );
  }
}

// ── 1단계: 환영 ──

class _WelcomePage extends StatelessWidget {
  final VoidCallback onNext;
  const _WelcomePage({required this.onNext});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SizedBox(
            width: 160,
            height: 160,
            child: Image.asset(
              'assets/images/premium/premium_onboarding_guide_illustrations_pack_1771750149063.png',
              fit: BoxFit.contain,
              errorBuilder: (_, __, ___) => const Icon(
                  Icons.biotech_rounded,
                  size: 64,
                  color: SanggamTheme.primary),
            ),
          ),
          const SizedBox(height: 40),
          const Text(
            '만파식에 오신 것을\n환영합니다',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.white,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          const Text(
            '초정밀 차동 계측 기술로\n언제 어디서나 건강을 관리하세요.',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 16,
              height: 1.6,
            ),
          ),
          const SizedBox(height: 16),
          const _FeatureRow(
            icon: Icons.science_rounded,
            text: '15종 이상 바이오마커 분석',
          ),
          const _FeatureRow(
            icon: Icons.smart_toy_rounded,
            text: 'AI 건강 코칭 및 트렌드 분석',
          ),
          const _FeatureRow(
            icon: Icons.family_restroom_rounded,
            text: '가족 건강 관리 및 공유',
          ),
          const SizedBox(height: 40),
          PrimaryButton(text: '시작하기', onPressed: onNext),
        ],
      ),
    );
  }
}

class _FeatureRow extends StatelessWidget {
  final IconData icon;
  final String text;
  const _FeatureRow({required this.icon, required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Icon(icon, size: 20, color: SanggamTheme.primary),
          const SizedBox(width: 16),
          Text(
            text,
            style: const TextStyle(color: Colors.white, fontSize: 14),
          ),
        ],
      ),
    );
  }
}

// ── 2단계: 건강 프로필 ──

class _HealthProfilePage extends StatelessWidget {
  final String selectedGoal;
  final int age;
  final ValueChanged<String> onGoalChanged;
  final ValueChanged<int> onAgeChanged;
  final VoidCallback onNext;

  const _HealthProfilePage({
    required this.selectedGoal,
    required this.age,
    required this.onGoalChanged,
    required this.onAgeChanged,
    required this.onNext,
  });

  @override
  Widget build(BuildContext context) {
    final goals = {
      'general': ('전반적 건강 관리', Icons.favorite_rounded),
      'diabetes': ('당뇨 관리', Icons.water_drop_rounded),
      'metabolic': ('대사 증후군 관리', Icons.monitor_heart_rounded),
      'fitness': ('운동/피트니스', Icons.fitness_center_rounded),
    };

    return SingleChildScrollView(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 32),
          const Text(
            '건강 목표를\n설정해주세요',
            style: TextStyle(
              color: Colors.white,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '맞춤형 건강 코칭을 위해 필요합니다.',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 32),

          // 나이 설정
          const Text(
            '나이',
            style: TextStyle(
              color: Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(
                child: Slider(
                  value: age.toDouble(),
                  min: 10,
                  max: 100,
                  divisions: 90,
                  label: '$age세',
                  onChanged: (v) => onAgeChanged(v.round()),
                ),
              ),
              SizedBox(
                width: 56,
                child: Text(
                  '$age세',
                  style: const TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // 건강 목표 선택
          const Text(
            '건강 목표',
            style: TextStyle(
              color: Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          ...goals.entries.map((e) {
            final isSelected = selectedGoal == e.key;
            return Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: InkWell(
                onTap: () => onGoalChanged(e.key),
                borderRadius: BorderRadius.circular(16),
                child: Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: isSelected
                          ? SanggamTheme.primary
                          : SanggamTheme.onSurfaceDim
                              .withValues(alpha: 0.3),
                      width: isSelected ? 2 : 1,
                    ),
                    color: isSelected
                        ? SanggamTheme.primary
                            .withValues(alpha: 0.08)
                        : null,
                  ),
                  child: Row(
                    children: [
                      Icon(
                        e.value.$2,
                        color: isSelected
                            ? SanggamTheme.primary
                            : SanggamTheme.onSurfaceDim,
                      ),
                      const SizedBox(width: 16),
                      Text(
                        e.value.$1,
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight: isSelected
                              ? FontWeight.bold
                              : FontWeight.normal,
                        ),
                      ),
                      const Spacer(),
                      if (isSelected)
                        const Icon(Icons.check_circle,
                            color: SanggamTheme.primary),
                    ],
                  ),
                ),
              ),
            );
          }),
          const SizedBox(height: 32),
          PrimaryButton(text: '다음', onPressed: onNext),
          const SizedBox(height: 32),
        ],
      ),
    );
  }
}

// ── 3단계: 디바이스 설정 ──

class _DeviceSetupPage extends StatelessWidget {
  final bool isScanning;
  final List<DeviceInfoDto> devices;
  final String? connectedDeviceId;
  final VoidCallback onScan;
  final ValueChanged<String> onConnect;
  final VoidCallback onNext;

  const _DeviceSetupPage({
    required this.isScanning,
    required this.devices,
    required this.connectedDeviceId,
    required this.onScan,
    required this.onConnect,
    required this.onNext,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 32),
          const Text(
            '디바이스를\n연결해주세요',
            style: TextStyle(
              color: Colors.white,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'ManPaSik 측정기와 BLE로 연결합니다.',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 32),

          // 스캔 버튼
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: isScanning ? null : onScan,
              icon: isScanning
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child:
                          CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.bluetooth_searching),
              label:
                  Text(isScanning ? '스캔 중...' : '주변 기기 검색'),
              style: OutlinedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                side: BorderSide(
                    color: SanggamTheme.onSurfaceDim
                        .withValues(alpha: 0.3)),
              ),
            ),
          ),
          const SizedBox(height: 24),

          // 디바이스 목록
          if (devices.isNotEmpty) ...[
            const Text(
              '발견된 기기',
              style: TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            ...devices.map((d) {
              final isConnected = connectedDeviceId == d.deviceId;
              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: SanggamContainer(
                  borderRadius: 16,
                  borderWidth: 0.8,
                  borderColor: isConnected
                      ? SanggamTheme.jagaeCyan
                      : SanggamTheme.primary
                          .withValues(alpha: 0.3),
                  blurSigma: 8,
                  jagaeOpacity: 0.03,
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      Icon(
                        isConnected
                            ? Icons.bluetooth_connected
                            : Icons.bluetooth,
                        color: isConnected
                            ? SanggamTheme.jagaeCyan
                            : SanggamTheme.onSurfaceDim,
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment:
                              CrossAxisAlignment.start,
                          children: [
                            Text(
                              d.name,
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 14,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            Text(
                              isConnected
                                  ? '연결됨'
                                  : 'RSSI: ${d.rssi} dBm',
                              style: const TextStyle(
                                color: SanggamTheme.onSurfaceDim,
                                fontSize: 12,
                              ),
                            ),
                          ],
                        ),
                      ),
                      if (isConnected)
                        const Icon(Icons.check_circle,
                            color: SanggamTheme.jagaeCyan)
                      else
                        TextButton(
                          onPressed: () => onConnect(d.deviceId),
                          child: const Text('연결'),
                        ),
                    ],
                  ),
                ),
              );
            }),
          ],

          if (devices.isEmpty && !isScanning)
            Expanded(
              child: Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.bluetooth_disabled_rounded,
                      size: 64,
                      color: SanggamTheme.onSurfaceDim
                          .withValues(alpha: 0.5),
                    ),
                    const SizedBox(height: 16),
                    Text(
                      '검색 버튼을 눌러\n주변 기기를 찾아보세요',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: SanggamTheme.onSurfaceDim
                            .withValues(alpha: 0.5),
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
              ),
            ),

          const Spacer(),

          PrimaryButton(
            text: connectedDeviceId != null ? '다음' : '나중에 연결하기',
            onPressed: onNext,
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }
}

// ── 4단계: 완료 ──

class _CompletePage extends StatelessWidget {
  final VoidCallback onComplete;
  const _CompletePage({required this.onComplete});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 96,
            height: 96,
            decoration: BoxDecoration(
              color: SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.check_rounded,
              size: 56,
              color: SanggamTheme.jagaeCyan,
            ),
          ),
          const SizedBox(height: 32),
          const Text(
            '설정 완료!',
            style: TextStyle(
              color: Colors.white,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          const Text(
            '만파식과 함께 건강한 하루를\n시작해보세요.',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 16,
              height: 1.6,
            ),
          ),
          const SizedBox(height: 48),
          PrimaryButton(text: '홈으로 이동', onPressed: onComplete),
        ],
      ),
    );
  }
}
