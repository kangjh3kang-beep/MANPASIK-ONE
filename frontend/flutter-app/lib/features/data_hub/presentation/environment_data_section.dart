import 'package:flutter/material.dart';

import 'package:manpasik/core/services/public_data_service.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// EnvironmentDataSection — Sanggam Orbit 환경 데이터 섹션
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.dancheongRed 2x → SanggamTheme.error
// [Rule 4] Theme.of(context).textTheme → 직접 TextStyle
// [Rule 4] theme.colorScheme.onSurfaceVariant → SanggamTheme.onSurfaceDim
// [Rule 4] Card → 다크 테마 Container
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// 환경 데이터 섹션 위젯 (C4)
///
/// DataHub 화면에서 날씨, 대기질, 감염병 경보를 표시합니다.
class EnvironmentDataSection extends StatefulWidget {
  const EnvironmentDataSection({super.key});

  @override
  State<EnvironmentDataSection> createState() =>
      _EnvironmentDataSectionState();
}

class _EnvironmentDataSectionState extends State<EnvironmentDataSection> {
  final _service = SimulatedPublicDataService();
  WeatherData? _weather;
  AirQualityData? _airQuality;
  List<DiseaseAlert> _alerts = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    try {
      final results = await Future.wait([
        _service.getWeather(latitude: 37.5665, longitude: 126.9780),
        _service.getAirQuality(stationName: '강남구'),
        _service.getDiseaseAlerts(),
      ]);
      if (mounted) {
        setState(() {
          _weather = results[0] as WeatherData;
          _airQuality = results[1] as AirQualityData;
          _alerts = results[2] as List<DiseaseAlert>;
          _loading = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return Container(
        decoration: BoxDecoration(
          color: Colors.white.withValues(alpha: 0.05),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: SanggamTheme.surfaceVariant),
        ),
        padding: const EdgeInsets.all(24),
        child: const Center(
          child: CircularProgressIndicator(color: SanggamTheme.primary),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '환경 데이터',
          style: TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),

        // 날씨 + 대기질
        Row(
          children: [
            Expanded(child: _buildWeatherCard()),
            const SizedBox(width: 12),
            Expanded(child: _buildAirQualityCard()),
          ],
        ),

        // 질병 경보
        if (_alerts.isNotEmpty) ...[
          const SizedBox(height: 12),
          ..._alerts.map((a) => _buildAlertCard(a)),
        ],
      ],
    );
  }

  Widget _buildWeatherCard() {
    if (_weather == null) return const SizedBox.shrink();
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: SanggamTheme.surfaceVariant),
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Icon(Icons.wb_sunny_rounded, size: 18,
                  color: SanggamTheme.primary),
              SizedBox(width: 4),
              Text('날씨', style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              )),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '${_weather!.temperature.toStringAsFixed(1)}°C',
            style: const TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            '${_weather!.condition} | 체감 ${_weather!.feelsLike.toStringAsFixed(1)}°C',
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
          Text(
            '습도 ${_weather!.humidity}% | UV ${_weather!.uvIndex}',
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 11,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAirQualityCard() {
    if (_airQuality == null) return const SizedBox.shrink();
    final gradeColor = switch (_airQuality!.grade) {
      AirQualityGrade.good => SanggamTheme.jagaeCyan,
      AirQualityGrade.moderate => SanggamTheme.primary,
      _ => SanggamTheme.error,
    };
    final gradeLabel = switch (_airQuality!.grade) {
      AirQualityGrade.good => '좋음',
      AirQualityGrade.moderate => '보통',
      AirQualityGrade.unhealthySensitive => '민감군 나쁨',
      AirQualityGrade.unhealthy => '나쁨',
      AirQualityGrade.veryUnhealthy => '매우 나쁨',
      AirQualityGrade.hazardous => '위험',
    };

    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: SanggamTheme.surfaceVariant),
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Icon(Icons.air_rounded, size: 18,
                  color: SanggamTheme.jagaeCyan),
              SizedBox(width: 4),
              Text('대기질', style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              )),
            ],
          ),
          const SizedBox(height: 8),
          Container(
            padding:
                const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: gradeColor.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              gradeLabel,
              style: TextStyle(
                color: gradeColor,
                fontSize: 14,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(height: 4),
          Text(
            'PM10: ${_airQuality!.pm10}㎍ | PM2.5: ${_airQuality!.pm25}㎍',
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 11,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAlertCard(DiseaseAlert alert) {
    final color = switch (alert.severity) {
      AlertSeverity.info => SanggamTheme.jagaeCyan,
      AlertSeverity.caution => SanggamTheme.primary,
      AlertSeverity.warning => SanggamTheme.error,
      AlertSeverity.critical => const Color(0xFF8B0000),
    };

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: color.withValues(alpha: 0.2),
        ),
      ),
      child: ListTile(
        leading: Icon(Icons.warning_amber_rounded, color: color),
        title: Text(alert.title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.bold,
            )),
        subtitle: Text(alert.description,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
            maxLines: 2),
        dense: true,
      ),
    );
  }
}
