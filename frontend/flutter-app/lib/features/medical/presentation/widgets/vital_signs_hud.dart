import 'dart:async';
import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 화상 진료 중 생체 신호 HUD 오버레이
///
/// 심박수, SpO2, 혈압, 체온, 스트레스 지수를 실시간 표시합니다.
/// 현재는 시뮬레이션 데이터를 사용하며, BLE 연동 시 실시간 스트림으로 교체합니다.
class VitalSignsHud extends StatefulWidget {
  const VitalSignsHud({super.key, this.isExpanded = false, this.onToggle});

  final bool isExpanded;
  final VoidCallback? onToggle;

  @override
  State<VitalSignsHud> createState() => _VitalSignsHudState();
}

class _VitalSignsHudState extends State<VitalSignsHud> with SingleTickerProviderStateMixin {
  late AnimationController _pulseController;
  Timer? _simulationTimer;

  // 시뮬레이션 데이터
  int _heartRate = 72;
  int _spo2 = 98;
  String _bloodPressure = '120/80';
  double _temperature = 36.5;
  int _stressLevel = 35;

  final _random = math.Random();

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1000),
    )..repeat(reverse: true);

    // 3초마다 시뮬레이션 데이터 갱신
    _simulationTimer = Timer.periodic(const Duration(seconds: 3), (_) {
      if (!mounted) return;
      setState(() {
        _heartRate = 68 + _random.nextInt(12); // 68-79
        _spo2 = 96 + _random.nextInt(4); // 96-99
        final sys = 115 + _random.nextInt(15); // 115-129
        final dia = 75 + _random.nextInt(10); // 75-84
        _bloodPressure = '$sys/$dia';
        _temperature = 36.3 + _random.nextDouble() * 0.6; // 36.3-36.9
        _stressLevel = 25 + _random.nextInt(25); // 25-49
      });
    });
  }

  @override
  void dispose() {
    _pulseController.dispose();
    _simulationTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!widget.isExpanded) {
      return _buildCompactHud();
    }
    return _buildExpandedHud();
  }

  /// 축소 상태: 좌하단에 심박수만 표시
  Widget _buildCompactHud() {
    return GestureDetector(
      onTap: widget.onToggle,
      child: AnimatedBuilder(
        animation: _pulseController,
        builder: (context, child) {
          final pulseScale = 1.0 + _pulseController.value * 0.05;
          return Transform.scale(
            scale: pulseScale,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: Colors.black.withOpacity(0.6),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(
                  color: _heartRateColor(_heartRate).withOpacity(0.5),
                  width: 1,
                ),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.favorite, size: 16, color: _heartRateColor(_heartRate)),
                  const SizedBox(width: 4),
                  Text(
                    '$_heartRate',
                    style: TextStyle(
                      color: _heartRateColor(_heartRate),
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      fontFeatures: const [FontFeature.tabularFigures()],
                    ),
                  ),
                  const SizedBox(width: 2),
                  Text(
                    'bpm',
                    style: TextStyle(color: Colors.white.withOpacity(0.5), fontSize: 10),
                  ),
                  const SizedBox(width: 8),
                  Icon(Icons.expand_less, size: 14, color: Colors.white.withOpacity(0.4)),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  /// 확장 상태: 전체 바이탈 사인 패널
  Widget _buildExpandedHud() {
    return GestureDetector(
      onTap: widget.onToggle,
      child: Container(
        width: double.infinity,
        margin: const EdgeInsets.symmetric(horizontal: 12),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: Colors.black.withOpacity(0.75),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.3), width: 0.5),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // 헤더
            Row(
              children: [
                Container(
                  width: 6,
                  height: 6,
                  decoration: const BoxDecoration(
                    color: Color(0xFF00E676),
                    shape: BoxShape.circle,
                  ),
                ),
                const SizedBox(width: 6),
                const Text(
                  'VITAL SIGNS',
                  style: TextStyle(
                    color: Colors.white54,
                    fontSize: 10,
                    fontWeight: FontWeight.w600,
                    letterSpacing: 1.5,
                  ),
                ),
                const Spacer(),
                Icon(Icons.expand_more, size: 14, color: Colors.white.withOpacity(0.4)),
              ],
            ),
            const SizedBox(height: 10),
            // 바이탈 그리드
            Row(
              children: [
                Expanded(child: _buildVitalTile(
                  icon: Icons.favorite,
                  label: '심박수',
                  value: '$_heartRate',
                  unit: 'bpm',
                  color: _heartRateColor(_heartRate),
                  isPulsing: true,
                )),
                const SizedBox(width: 8),
                Expanded(child: _buildVitalTile(
                  icon: Icons.water_drop,
                  label: 'SpO₂',
                  value: '$_spo2',
                  unit: '%',
                  color: _spo2Color(_spo2),
                )),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(child: _buildVitalTile(
                  icon: Icons.speed,
                  label: '혈압',
                  value: _bloodPressure,
                  unit: 'mmHg',
                  color: const Color(0xFF29B6F6),
                )),
                const SizedBox(width: 8),
                Expanded(child: _buildVitalTile(
                  icon: Icons.thermostat,
                  label: '체온',
                  value: _temperature.toStringAsFixed(1),
                  unit: '°C',
                  color: _temperatureColor(_temperature),
                )),
              ],
            ),
            const SizedBox(height: 8),
            // 스트레스 바
            _buildStressBar(),
          ],
        ),
      ),
    );
  }

  Widget _buildVitalTile({
    required IconData icon,
    required String label,
    required String value,
    required String unit,
    required Color color,
    bool isPulsing = false,
  }) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: color.withOpacity(0.08),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: color.withOpacity(0.15)),
      ),
      child: Row(
        children: [
          if (isPulsing)
            AnimatedBuilder(
              animation: _pulseController,
              builder: (_, __) => Icon(
                icon,
                size: 16,
                color: color.withOpacity(0.6 + _pulseController.value * 0.4),
              ),
            )
          else
            Icon(icon, size: 16, color: color),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: TextStyle(color: Colors.white.withOpacity(0.4), fontSize: 9),
                ),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.baseline,
                  textBaseline: TextBaseline.alphabetic,
                  children: [
                    Text(
                      value,
                      style: TextStyle(
                        color: color,
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        fontFeatures: const [FontFeature.tabularFigures()],
                      ),
                    ),
                    const SizedBox(width: 2),
                    Text(
                      unit,
                      style: TextStyle(color: Colors.white.withOpacity(0.3), fontSize: 9),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStressBar() {
    final stressColor = _stressColor(_stressLevel);
    final stressLabel = _stressLevel < 30
        ? '안정'
        : _stressLevel < 50
            ? '보통'
            : _stressLevel < 70
                ? '높음'
                : '매우 높음';

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: stressColor.withOpacity(0.08),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: stressColor.withOpacity(0.15)),
      ),
      child: Row(
        children: [
          Icon(Icons.psychology, size: 16, color: stressColor),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('스트레스', style: TextStyle(color: Colors.white.withOpacity(0.4), fontSize: 9)),
              Text(
                '$_stressLevel ($stressLabel)',
                style: TextStyle(color: stressColor, fontSize: 14, fontWeight: FontWeight.w700),
              ),
            ],
          ),
          const Spacer(),
          SizedBox(
            width: 100,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: _stressLevel / 100,
                backgroundColor: Colors.white.withOpacity(0.1),
                valueColor: AlwaysStoppedAnimation(stressColor),
                minHeight: 6,
              ),
            ),
          ),
        ],
      ),
    );
  }

  // ── 색상 헬퍼 ──

  Color _heartRateColor(int bpm) {
    if (bpm < 60) return const Color(0xFF29B6F6); // 서맥 — 파랑
    if (bpm <= 100) return const Color(0xFF00E676); // 정상 — 녹색
    return Colors.redAccent; // 빈맥 — 빨강
  }

  Color _spo2Color(int spo2) {
    if (spo2 >= 95) return const Color(0xFF00E676);
    if (spo2 >= 90) return Colors.orange;
    return Colors.redAccent;
  }

  Color _temperatureColor(double temp) {
    if (temp >= 36.1 && temp <= 37.2) return const Color(0xFF00E676);
    if (temp > 37.2 && temp <= 38.0) return Colors.orange;
    if (temp > 38.0) return Colors.redAccent;
    return const Color(0xFF29B6F6); // 저체온
  }

  Color _stressColor(int level) {
    if (level < 30) return const Color(0xFF00E676);
    if (level < 50) return const Color(0xFFFFCA28);
    if (level < 70) return Colors.orange;
    return Colors.redAccent;
  }
}
