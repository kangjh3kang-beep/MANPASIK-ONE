import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

import 'package:manpasik/core/services/rest_client.dart';

/// 약국 API 서비스 인터페이스 (B8)
///
/// 공공데이터포털 약국 API 또는 건강보험심사평가원 API 연동 시
/// 이 인터페이스를 구현합니다. API 키 미설정 시 SimulatedPharmacyService 사용.
abstract class PharmacyService {
  Future<List<PharmacyInfo>> searchNearby({
    required double latitude,
    required double longitude,
    double radiusKm = 3.0,
    int limit = 20,
  });

  Future<PharmacyInfo?> getById(String pharmacyId);

  Future<bool> sendPrescription({
    required String pharmacyId,
    required String prescriptionId,
  });

  /// 구성 상태에 따라 적절한 구현체를 반환하는 팩토리
  static PharmacyService create({ManPaSikRestClient? restClient}) {
    if (RestPharmacyService.isConfigured && restClient != null) {
      return RestPharmacyService(restClient: restClient);
    }
    return SimulatedPharmacyService();
  }
}

/// REST API 기반 약국 서비스
///
/// 공공데이터포털 약국 정보 API와 ManPaSik 백엔드를 조합하여
/// 주변 약국 검색 및 처방전 전송 기능을 제공합니다.
/// PHARMACY_API_KEY 환경변수로 공공데이터 API 키를 설정합니다.
class RestPharmacyService implements PharmacyService {
  RestPharmacyService({required this.restClient});

  final ManPaSikRestClient restClient;

  static const _apiKey = String.fromEnvironment('PHARMACY_API_KEY');
  static bool get isConfigured => _apiKey.isNotEmpty;

  // 공공데이터포털 약국 API
  static const _publicApiBaseUrl =
      'http://apis.data.go.kr/B552657/ErmctInsttInfoInqireService';

  final _publicDio = Dio(BaseOptions(
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
  ));

  @override
  Future<List<PharmacyInfo>> searchNearby({
    required double latitude,
    required double longitude,
    double radiusKm = 3.0,
    int limit = 20,
  }) async {
    if (!isConfigured) {
      debugPrint('[Pharmacy] API 키 미설정 → 시뮬레이션');
      return _simulateNearbySearch(latitude, longitude);
    }

    try {
      final resp = await _publicDio.get(
        '$_publicApiBaseUrl/getParmacyListInfoInqire',
        queryParameters: {
          'serviceKey': _apiKey,
          'WGS84_LON': longitude.toString(),
          'WGS84_LAT': latitude.toString(),
          'pageNo': '1',
          'numOfRows': limit.toString(),
          '_type': 'json',
        },
      );

      final body = resp.data as Map<String, dynamic>? ?? {};
      final response = body['response'] as Map<String, dynamic>? ?? {};
      final bodyData = response['body'] as Map<String, dynamic>? ?? {};
      final items = bodyData['items'] as Map<String, dynamic>? ?? {};
      final itemList = items['item'] as List<dynamic>? ?? [];

      return itemList.map((item) {
        final m = item as Map<String, dynamic>;
        final lat = double.tryParse(m['wgs84Lat']?.toString() ?? '') ?? latitude;
        final lon = double.tryParse(m['wgs84Lon']?.toString() ?? '') ?? longitude;
        return PharmacyInfo(
          id: m['hpid']?.toString() ?? '',
          name: m['dutyName']?.toString() ?? '',
          address: m['dutyAddr']?.toString() ?? '',
          phone: m['dutyTel1']?.toString() ?? '',
          latitude: lat,
          longitude: lon,
          distanceKm: _calcDistance(latitude, longitude, lat, lon),
          isOpen: _isCurrentlyOpen(m),
          operatingHours: _parseOperatingHours(m),
        );
      }).toList()
        ..sort((a, b) => a.distanceKm.compareTo(b.distanceKm));
    } on DioException catch (e) {
      debugPrint('[Pharmacy] 공공데이터 API 오류: $e');
      return _simulateNearbySearch(latitude, longitude);
    }
  }

  @override
  Future<PharmacyInfo?> getById(String pharmacyId) async {
    if (!isConfigured) {
      return _simulatePharmacyDetail(pharmacyId);
    }

    try {
      final resp = await _publicDio.get(
        '$_publicApiBaseUrl/getParmacyBassInfoInqire',
        queryParameters: {
          'serviceKey': _apiKey,
          'HPID': pharmacyId,
          '_type': 'json',
        },
      );

      final body = resp.data as Map<String, dynamic>? ?? {};
      final response = body['response'] as Map<String, dynamic>? ?? {};
      final bodyData = response['body'] as Map<String, dynamic>? ?? {};
      final items = bodyData['items'] as Map<String, dynamic>? ?? {};
      final item = items['item'];
      if (item == null) return null;

      final m = item is List ? item.first as Map<String, dynamic> : item as Map<String, dynamic>;
      final lat = double.tryParse(m['wgs84Lat']?.toString() ?? '') ?? 0;
      final lon = double.tryParse(m['wgs84Lon']?.toString() ?? '') ?? 0;
      return PharmacyInfo(
        id: m['hpid']?.toString() ?? pharmacyId,
        name: m['dutyName']?.toString() ?? '',
        address: m['dutyAddr']?.toString() ?? '',
        phone: m['dutyTel1']?.toString() ?? '',
        latitude: lat,
        longitude: lon,
        distanceKm: 0,
        isOpen: _isCurrentlyOpen(m),
        operatingHours: _parseOperatingHours(m),
      );
    } on DioException {
      return _simulatePharmacyDetail(pharmacyId);
    }
  }

  @override
  Future<bool> sendPrescription({
    required String pharmacyId,
    required String prescriptionId,
  }) async {
    try {
      await restClient.selectPharmacy(
        prescriptionId,
        pharmacyId: pharmacyId,
      );
      await restClient.sendToPharmacy(prescriptionId);
      return true;
    } catch (e) {
      debugPrint('[Pharmacy] 처방전 전송 실패: $e');
      return false;
    }
  }

  double _calcDistance(double lat1, double lon1, double lat2, double lon2) {
    // 간단한 Haversine 근사 (km)
    const earthRadius = 6371.0;
    final dLat = (lat2 - lat1) * 3.14159265 / 180;
    final dLon = (lon2 - lon1) * 3.14159265 / 180;
    final a = _sin2(dLat / 2) +
        _cos(lat1 * 3.14159265 / 180) *
            _cos(lat2 * 3.14159265 / 180) *
            _sin2(dLon / 2);
    return earthRadius * 2 * _asin(_sqrt(a));
  }

  // 간단한 삼각함수 (dart:math 불필요)
  double _sin2(double x) {
    final s = x - x * x * x / 6 + x * x * x * x * x / 120;
    return s * s;
  }
  double _cos(double x) => 1 - x * x / 2 + x * x * x * x / 24;
  double _asin(double x) => x + x * x * x / 6;
  double _sqrt(double x) {
    if (x <= 0) return 0;
    double guess = x / 2;
    for (int i = 0; i < 10; i++) {
      guess = (guess + x / guess) / 2;
    }
    return guess;
  }

  bool _isCurrentlyOpen(Map<String, dynamic> m) {
    final now = DateTime.now();
    final day = now.weekday; // 1=Mon, 7=Sun
    final timeKey = 'dutyTime${day}s';
    final closeKey = 'dutyTime${day}c';
    final openStr = m[timeKey]?.toString();
    final closeStr = m[closeKey]?.toString();
    if (openStr == null || closeStr == null) return false;
    final nowMin = now.hour * 100 + now.minute;
    final openMin = int.tryParse(openStr) ?? 0;
    final closeMin = int.tryParse(closeStr) ?? 2400;
    return nowMin >= openMin && nowMin <= closeMin;
  }

  String _parseOperatingHours(Map<String, dynamic> m) {
    final open = m['dutyTime1s']?.toString() ?? '0900';
    final close = m['dutyTime1c']?.toString() ?? '1800';
    return '${open.substring(0, 2)}:${open.substring(2)} ~ ${close.substring(0, 2)}:${close.substring(2)}';
  }

  List<PharmacyInfo> _simulateNearbySearch(double lat, double lon) {
    return [
      PharmacyInfo(id: 'pharm_001', name: '건강 약국', address: '서울시 강남구 테헤란로 123', phone: '02-1234-5678', latitude: lat + 0.002, longitude: lon + 0.001, distanceKm: 0.3, isOpen: true, operatingHours: '09:00 ~ 21:00'),
      PharmacyInfo(id: 'pharm_002', name: '사랑 약국', address: '서울시 강남구 역삼로 456', phone: '02-9876-5432', latitude: lat - 0.001, longitude: lon + 0.003, distanceKm: 0.8, isOpen: true, operatingHours: '08:30 ~ 22:00'),
      PharmacyInfo(id: 'pharm_003', name: '온누리 약국', address: '서울시 강남구 선릉로 789', phone: '02-5555-1234', latitude: lat + 0.004, longitude: lon - 0.002, distanceKm: 1.2, isOpen: false, operatingHours: '09:00 ~ 18:00'),
    ];
  }

  PharmacyInfo _simulatePharmacyDetail(String pharmacyId) {
    return PharmacyInfo(
      id: pharmacyId,
      name: '건강 약국',
      address: '서울시 강남구 테헤란로 123',
      phone: '02-1234-5678',
      latitude: 37.5665,
      longitude: 126.9780,
      distanceKm: 0.3,
      isOpen: true,
      operatingHours: '09:00 ~ 21:00',
    );
  }
}

/// 시뮬레이션 약국 서비스
class SimulatedPharmacyService implements PharmacyService {
  @override
  Future<List<PharmacyInfo>> searchNearby({
    required double latitude,
    required double longitude,
    double radiusKm = 3.0,
    int limit = 20,
  }) async {
    debugPrint('[SimulatedPharmacy] 주변 약국 검색: ($latitude, $longitude)');
    await Future.delayed(const Duration(milliseconds: 500));
    return [
      PharmacyInfo(id: 'pharm_001', name: '건강 약국', address: '서울시 강남구 테헤란로 123', phone: '02-1234-5678', latitude: latitude + 0.002, longitude: longitude + 0.001, distanceKm: 0.3, isOpen: true, operatingHours: '09:00 ~ 21:00'),
      PharmacyInfo(id: 'pharm_002', name: '사랑 약국', address: '서울시 강남구 역삼로 456', phone: '02-9876-5432', latitude: latitude - 0.001, longitude: longitude + 0.003, distanceKm: 0.8, isOpen: true, operatingHours: '08:30 ~ 22:00'),
      PharmacyInfo(id: 'pharm_003', name: '온누리 약국', address: '서울시 강남구 선릉로 789', phone: '02-5555-1234', latitude: latitude + 0.004, longitude: longitude - 0.002, distanceKm: 1.2, isOpen: false, operatingHours: '09:00 ~ 18:00'),
    ];
  }

  @override
  Future<PharmacyInfo?> getById(String pharmacyId) async {
    debugPrint('[SimulatedPharmacy] 약국 조회: $pharmacyId');
    await Future.delayed(const Duration(milliseconds: 300));
    return PharmacyInfo(id: pharmacyId, name: '건강 약국', address: '서울시 강남구 테헤란로 123', phone: '02-1234-5678', latitude: 37.5665, longitude: 126.9780, distanceKm: 0.3, isOpen: true, operatingHours: '09:00 ~ 21:00');
  }

  @override
  Future<bool> sendPrescription({
    required String pharmacyId,
    required String prescriptionId,
  }) async {
    debugPrint('[SimulatedPharmacy] 처방전 전송: $prescriptionId → $pharmacyId');
    await Future.delayed(const Duration(seconds: 1));
    return true;
  }
}

class PharmacyInfo {
  final String id;
  final String name;
  final String address;
  final String phone;
  final double latitude;
  final double longitude;
  final double distanceKm;
  final bool isOpen;
  final String operatingHours;

  const PharmacyInfo({
    required this.id,
    required this.name,
    required this.address,
    required this.phone,
    required this.latitude,
    required this.longitude,
    required this.distanceKm,
    required this.isOpen,
    required this.operatingHours,
  });
}
