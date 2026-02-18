import 'package:flutter/material.dart';

import 'package:manpasik/core/services/public_data_service.dart';
import 'package:manpasik/core/theme/app_theme.dart';

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
    final theme = Theme.of(context);

    if (_loading) {
      return const Card(
        child: Padding(
          padding: EdgeInsets.all(24),
          child: Center(child: CircularProgressIndicator()),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '환경 데이터',
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),

        // 날씨 + 대기질
        Row(
          children: [
            Expanded(child: _buildWeatherCard(theme)),
            const SizedBox(width: 12),
            Expanded(child: _buildAirQualityCard(theme)),
          ],
        ),

        // 질병 경보
        if (_alerts.isNotEmpty) ...[
          const SizedBox(height: 12),
          ..._alerts.map((a) => _buildAlertCard(theme, a)),
        ],
      ],
    );
  }

  Widget _buildWeatherCard(ThemeData theme) {
    if (_weather == null) return const SizedBox.shrink();
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.wb_sunny_rounded, size: 18),
                const SizedBox(width: 4),
                Text('날씨', style: theme.textTheme.bodySmall),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              '${_weather!.temperature.toStringAsFixed(1)}°C',
              style: theme.textTheme.headlineSmall
                  ?.copyWith(fontWeight: FontWeight.bold),
            ),
            Text(
              '${_weather!.condition} | 체감 ${_weather!.feelsLike.toStringAsFixed(1)}°C',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            Text(
              '습도 ${_weather!.humidity}% | UV ${_weather!.uvIndex}',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
                fontSize: 11,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAirQualityCard(ThemeData theme) {
    if (_airQuality == null) return const SizedBox.shrink();
    final gradeColor = switch (_airQuality!.grade) {
      AirQualityGrade.good => Colors.green,
      AirQualityGrade.moderate => Colors.orange,
      _ => AppTheme.dancheongRed,
    };
    final gradeLabel = switch (_airQuality!.grade) {
      AirQualityGrade.good => '좋음',
      AirQualityGrade.moderate => '보통',
      AirQualityGrade.unhealthySensitive => '민감군 나쁨',
      AirQualityGrade.unhealthy => '나쁨',
      AirQualityGrade.veryUnhealthy => '매우 나쁨',
      AirQualityGrade.hazardous => '위험',
    };

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.air_rounded, size: 18),
                const SizedBox(width: 4),
                Text('대기질', style: theme.textTheme.bodySmall),
              ],
            ),
            const SizedBox(height: 8),
            Container(
              padding:
                  const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: gradeColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(
                gradeLabel,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: gradeColor,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const SizedBox(height: 4),
            Text(
              'PM10: ${_airQuality!.pm10}㎍ | PM2.5: ${_airQuality!.pm25}㎍',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
                fontSize: 11,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAlertCard(ThemeData theme, DiseaseAlert alert) {
    final color = switch (alert.severity) {
      AlertSeverity.info => Colors.blue,
      AlertSeverity.caution => Colors.orange,
      AlertSeverity.warning => AppTheme.dancheongRed,
      AlertSeverity.critical => const Color(0xFF8B0000),
    };

    return Card(
      color: color.withOpacity(0.05),
      child: ListTile(
        leading: Icon(Icons.warning_amber_rounded, color: color),
        title: Text(alert.title,
            style: theme.textTheme.bodyMedium
                ?.copyWith(fontWeight: FontWeight.bold)),
        subtitle: Text(alert.description,
            style: theme.textTheme.bodySmall, maxLines: 2),
        dense: true,
      ),
    );
  }
}
