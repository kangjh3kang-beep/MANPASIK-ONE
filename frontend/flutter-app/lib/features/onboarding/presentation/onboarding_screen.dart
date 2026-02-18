import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';

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

  // 건강 프로필 데이터
  String _selectedGoal = 'general';
  int _age = 30;

  // 디바이스 페어링
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
    final theme = Theme.of(context);

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // 상단 프로그레스 + 건너뛰기
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
              child: Row(
                children: [
                  Expanded(
                    child: _ProgressIndicator(
                      current: _currentPage,
                      total: _totalPages,
                    ),
                  ),
                  TextButton(
                    onPressed: _completeOnboarding,
                    child: Text(
                      '건너뛰기',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
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
                onPageChanged: (i) => setState(() => _currentPage = i),
                children: [
                  _WelcomePage(onNext: _nextPage),
                  _HealthProfilePage(
                    selectedGoal: _selectedGoal,
                    age: _age,
                    onGoalChanged: (g) => setState(() => _selectedGoal = g),
                    onAgeChanged: (a) => setState(() => _age = a),
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
    );
  }
}

// ── 프로그레스 인디케이터 ──

class _ProgressIndicator extends StatelessWidget {
  final int current;
  final int total;

  const _ProgressIndicator({required this.current, required this.total});

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
                  ? AppTheme.sanggamGold
                  : Theme.of(context).colorScheme.outlineVariant,
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
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [AppTheme.deepSeaBlue, Color(0xFF112240)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(32),
              border: Border.all(
                color: AppTheme.sanggamGold.withOpacity(0.4),
              ),
            ),
            child: const Icon(
              Icons.biotech_rounded,
              size: 64,
              color: AppTheme.sanggamGold,
            ),
          ),
          const SizedBox(height: 40),
          Text(
            '만파식에 오신 것을\n환영합니다',
            textAlign: TextAlign.center,
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          Text(
            '초정밀 차동 계측 기술로\n언제 어디서나 건강을 관리하세요.',
            textAlign: TextAlign.center,
            style: theme.textTheme.bodyLarge?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
              height: 1.6,
            ),
          ),
          const SizedBox(height: 16),
          // 주요 기능 소개
          _FeatureRow(
            icon: Icons.science_rounded,
            text: '15종 이상 바이오마커 분석',
          ),
          _FeatureRow(
            icon: Icons.smart_toy_rounded,
            text: 'AI 건강 코칭 및 트렌드 분석',
          ),
          _FeatureRow(
            icon: Icons.family_restroom_rounded,
            text: '가족 건강 관리 및 공유',
          ),
          const SizedBox(height: 40),
          SizedBox(
            width: double.infinity,
            child: FilledButton(
              onPressed: onNext,
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
              child: const Text('시작하기'),
            ),
          ),
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
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          Icon(icon, size: 20, color: AppTheme.sanggamGold),
          const SizedBox(width: 12),
          Text(text, style: theme.textTheme.bodyMedium),
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
    final theme = Theme.of(context);

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
          Text(
            '건강 목표를\n설정해주세요',
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            '맞춤형 건강 코칭을 위해 필요합니다.',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
          const SizedBox(height: 32),

          // 나이 설정
          Text('나이', style: theme.textTheme.titleSmall),
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
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.sanggamGold,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // 건강 목표 선택
          Text('건강 목표', style: theme.textTheme.titleSmall),
          const SizedBox(height: 12),
          ...goals.entries.map((e) {
            final isSelected = selectedGoal == e.key;
            return Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: InkWell(
                onTap: () => onGoalChanged(e.key),
                borderRadius: BorderRadius.circular(12),
                child: Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: isSelected
                          ? AppTheme.sanggamGold
                          : theme.colorScheme.outlineVariant,
                      width: isSelected ? 2 : 1,
                    ),
                    color: isSelected
                        ? AppTheme.sanggamGold.withOpacity(0.08)
                        : null,
                  ),
                  child: Row(
                    children: [
                      Icon(
                        e.value.$2,
                        color: isSelected
                            ? AppTheme.sanggamGold
                            : theme.colorScheme.onSurfaceVariant,
                      ),
                      const SizedBox(width: 12),
                      Text(
                        e.value.$1,
                        style: theme.textTheme.bodyLarge?.copyWith(
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.normal,
                        ),
                      ),
                      const Spacer(),
                      if (isSelected)
                        const Icon(Icons.check_circle,
                            color: AppTheme.sanggamGold),
                    ],
                  ),
                ),
              ),
            );
          }),
          const SizedBox(height: 32),
          SizedBox(
            width: double.infinity,
            child: FilledButton(
              onPressed: onNext,
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
              child: const Text('다음'),
            ),
          ),
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
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 32),
          Text(
            '디바이스를\n연결해주세요',
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'ManPaSik 측정기와 BLE로 연결합니다.',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
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
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.bluetooth_searching),
              label: Text(isScanning ? '스캔 중...' : '주변 기기 검색'),
              style: OutlinedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
          const SizedBox(height: 24),

          // 디바이스 목록
          if (devices.isNotEmpty) ...[
            Text('발견된 기기', style: theme.textTheme.titleSmall),
            const SizedBox(height: 8),
            ...devices.map((d) {
              final isConnected = connectedDeviceId == d.deviceId;
              return Card(
                margin: const EdgeInsets.only(bottom: 8),
                child: ListTile(
                  leading: Icon(
                    isConnected
                        ? Icons.bluetooth_connected
                        : Icons.bluetooth,
                    color: isConnected ? Colors.green : null,
                  ),
                  title: Text(d.name),
                  subtitle: Text(
                    isConnected
                        ? '연결됨'
                        : 'RSSI: ${d.rssi} dBm',
                  ),
                  trailing: isConnected
                      ? const Icon(Icons.check_circle, color: Colors.green)
                      : TextButton(
                          onPressed: () => onConnect(d.deviceId),
                          child: const Text('연결'),
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
                      color: theme.colorScheme.outline,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      '검색 버튼을 눌러\n주변 기기를 찾아보세요',
                      textAlign: TextAlign.center,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.outline,
                      ),
                    ),
                  ],
                ),
              ),
            ),

          const Spacer(),

          // 다음/나중에 연결
          SizedBox(
            width: double.infinity,
            child: FilledButton(
              onPressed: onNext,
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
              child: Text(connectedDeviceId != null ? '다음' : '나중에 연결하기'),
            ),
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
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 100,
            height: 100,
            decoration: BoxDecoration(
              color: Colors.green.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.check_rounded,
              size: 56,
              color: Colors.green,
            ),
          ),
          const SizedBox(height: 32),
          Text(
            '설정 완료!',
            style: theme.textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            '만파식과 함께 건강한 하루를\n시작해보세요.',
            textAlign: TextAlign.center,
            style: theme.textTheme.bodyLarge?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
              height: 1.6,
            ),
          ),
          const SizedBox(height: 48),
          SizedBox(
            width: double.infinity,
            child: FilledButton(
              onPressed: onComplete,
              style: FilledButton.styleFrom(
                backgroundColor: AppTheme.sanggamGold,
                foregroundColor: AppTheme.deepSeaBlue,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
              child: const Text('홈으로 이동'),
            ),
          ),
        ],
      ),
    );
  }
}
