import 'package:flutter/foundation.dart';

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
      PharmacyInfo(
        id: 'pharm_001',
        name: '건강 약국',
        address: '서울시 강남구 테헤란로 123',
        phone: '02-1234-5678',
        latitude: latitude + 0.002,
        longitude: longitude + 0.001,
        distanceKm: 0.3,
        isOpen: true,
        operatingHours: '09:00 ~ 21:00',
      ),
      PharmacyInfo(
        id: 'pharm_002',
        name: '사랑 약국',
        address: '서울시 강남구 역삼로 456',
        phone: '02-9876-5432',
        latitude: latitude - 0.001,
        longitude: longitude + 0.003,
        distanceKm: 0.8,
        isOpen: true,
        operatingHours: '08:30 ~ 22:00',
      ),
      PharmacyInfo(
        id: 'pharm_003',
        name: '온누리 약국',
        address: '서울시 강남구 선릉로 789',
        phone: '02-5555-1234',
        latitude: latitude + 0.004,
        longitude: longitude - 0.002,
        distanceKm: 1.2,
        isOpen: false,
        operatingHours: '09:00 ~ 18:00',
      ),
    ];
  }

  @override
  Future<PharmacyInfo?> getById(String pharmacyId) async {
    debugPrint('[SimulatedPharmacy] 약국 조회: $pharmacyId');
    await Future.delayed(const Duration(milliseconds: 300));
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
