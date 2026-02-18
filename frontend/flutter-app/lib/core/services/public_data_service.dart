import 'package:flutter/foundation.dart';

/// 공공데이터 서비스 인터페이스 (C11)
///
/// 기상청 날씨, 미세먼지(에어코리아), 질병관리청 감염병 경보 등
/// 공공 데이터 API 연동 시 이 인터페이스를 구현합니다.
/// API 키 미설정 시 SimulatedPublicDataService 사용.
abstract class PublicDataService {
  Future<WeatherData> getWeather({
    required double latitude,
    required double longitude,
  });

  Future<AirQualityData> getAirQuality({
    required String stationName,
  });

  Future<List<DiseaseAlert>> getDiseaseAlerts();
}

/// 시뮬레이션 공공데이터 서비스
class SimulatedPublicDataService implements PublicDataService {
  @override
  Future<WeatherData> getWeather({
    required double latitude,
    required double longitude,
  }) async {
    debugPrint('[SimulatedPublicData] 날씨 조회: ($latitude, $longitude)');
    await Future.delayed(const Duration(milliseconds: 300));
    return WeatherData(
      temperature: 5.2,
      humidity: 45,
      condition: '맑음',
      pm10: 32,
      pm25: 18,
      uvIndex: 3,
      feelsLike: 2.1,
      updatedAt: DateTime.now(),
    );
  }

  @override
  Future<AirQualityData> getAirQuality({
    required String stationName,
  }) async {
    debugPrint('[SimulatedPublicData] 대기질 조회: $stationName');
    await Future.delayed(const Duration(milliseconds: 300));
    return AirQualityData(
      stationName: stationName,
      pm10: 32,
      pm25: 18,
      o3: 0.035,
      no2: 0.028,
      co: 0.5,
      so2: 0.003,
      grade: AirQualityGrade.good,
      updatedAt: DateTime.now(),
    );
  }

  @override
  Future<List<DiseaseAlert>> getDiseaseAlerts() async {
    debugPrint('[SimulatedPublicData] 감염병 경보 조회');
    await Future.delayed(const Duration(milliseconds: 300));
    return [
      DiseaseAlert(
        id: 'alert_001',
        title: '인플루엔자 주의보',
        description: '전국적으로 인플루엔자 환자 증가 추세입니다. 예방접종을 권장합니다.',
        severity: AlertSeverity.caution,
        region: '전국',
        publishedAt: DateTime.now().subtract(const Duration(hours: 6)),
      ),
    ];
  }
}

class WeatherData {
  final double temperature;
  final int humidity;
  final String condition;
  final int pm10;
  final int pm25;
  final int uvIndex;
  final double feelsLike;
  final DateTime updatedAt;

  const WeatherData({
    required this.temperature,
    required this.humidity,
    required this.condition,
    required this.pm10,
    required this.pm25,
    required this.uvIndex,
    required this.feelsLike,
    required this.updatedAt,
  });
}

class AirQualityData {
  final String stationName;
  final int pm10;
  final int pm25;
  final double o3;
  final double no2;
  final double co;
  final double so2;
  final AirQualityGrade grade;
  final DateTime updatedAt;

  const AirQualityData({
    required this.stationName,
    required this.pm10,
    required this.pm25,
    required this.o3,
    required this.no2,
    required this.co,
    required this.so2,
    required this.grade,
    required this.updatedAt,
  });
}

enum AirQualityGrade { good, moderate, unhealthySensitive, unhealthy, veryUnhealthy, hazardous }

enum AlertSeverity { info, caution, warning, critical }

class DiseaseAlert {
  final String id;
  final String title;
  final String description;
  final AlertSeverity severity;
  final String region;
  final DateTime publishedAt;

  const DiseaseAlert({
    required this.id,
    required this.title,
    required this.description,
    required this.severity,
    required this.region,
    required this.publishedAt,
  });
}
