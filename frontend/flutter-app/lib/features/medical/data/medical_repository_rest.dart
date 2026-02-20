import 'package:dio/dio.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 MedicalRepository 구현체
class MedicalRepositoryRest implements MedicalRepository {
  MedicalRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<TelemedicineReservation> createReservation({
    required String doctorId,
    required DateTime scheduledAt,
  }) async {
    final res = await _client.createReservation(
      userId: userId,
      facilityId: doctorId,
      reason: 'telemedicine',
    );
    return _mapReservation(res);
  }

  @override
  Future<List<TelemedicineReservation>> getReservations() async {
    try {
      final res = await _client.listReservations(userId);
      final list = res['reservations'] as List<dynamic>? ?? [];
      return list
          .map((r) => _mapReservation(r as Map<String, dynamic>))
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<void> cancelReservation(String reservationId) async {
    await _client.cancelReservation(reservationId);
  }

  @override
  Future<List<Prescription>> getPrescriptions() async {
    try {
      // Use health records with type filter for prescriptions
      final res = await _client.listHealthRecords(userId, typeFilter: 2);
      final list = res['records'] as List<dynamic>? ?? [];
      return list
          .map((r) => _mapPrescription(r as Map<String, dynamic>))
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<Prescription> getPrescriptionDetail(String prescriptionId) async {
    final res = await _client.getHealthRecord(prescriptionId);
    return _mapPrescription(res);
  }

  @override
  Future<HealthReport> generateHealthReport() async {
    final res = await _client.generateDailyReport(userId);
    return _mapHealthReport(res);
  }

  @override
  Future<List<HealthReport>> getHealthReports() async {
    try {
      final res = await _client.listHealthRecords(userId, typeFilter: 1);
      final list = res['records'] as List<dynamic>? ?? [];
      return list
          .map((r) => _mapHealthReport(r as Map<String, dynamic>))
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<void> sendEmergencyAlert({
    required String alertType,
    required double abnormalValue,
    required String biomarkerType,
  }) async {
    await _client.sendNotification(
      userId: userId,
      title: '긴급 건강 알림: $biomarkerType',
      body: '$alertType - 측정값: $abnormalValue',
      type: 'emergency',
    );
  }

  @override
  Future<List<DoctorInfo>> getRecommendedDoctors({int limit = 5}) async {
    try {
      final res = await _client.searchFacilities(type: 'doctor', limit: limit);
      final doctors = res['facilities'] as List<dynamic>? ?? [];
      return doctors.map((d) {
        final m = d as Map<String, dynamic>;
        return DoctorInfo(
          doctorId: m['facility_id'] as String? ?? m['id'] as String? ?? '',
          name: m['name'] as String? ?? '',
          specialty: m['specialty'] as String? ?? m['type'] as String? ?? '',
          rating: (m['rating'] as num?)?.toDouble() ?? 4.5,
          reviewCount: m['review_count'] as int? ?? 0,
          avatarUrl: m['avatar_url'] as String?,
          isAvailable: m['is_available'] as bool? ?? true,
          nextSlot: m['next_slot'] as String?,
        );
      }).toList();
    } on DioException {
      // 서버 미연결 시 시뮬레이션 데이터
      return [
        const DoctorInfo(doctorId: 'doc-1', name: '김민수', specialty: '내과', rating: 4.8, reviewCount: 234, isAvailable: true, nextSlot: '오늘 15:00'),
        const DoctorInfo(doctorId: 'doc-2', name: '이서연', specialty: '피부과', rating: 4.9, reviewCount: 189, isAvailable: false, nextSlot: '내일 10:00'),
        const DoctorInfo(doctorId: 'doc-3', name: '박준영', specialty: '정신건강', rating: 4.7, reviewCount: 156, isAvailable: true, nextSlot: '오늘 16:30'),
        const DoctorInfo(doctorId: 'doc-4', name: '최수진', specialty: '내분비', rating: 4.6, reviewCount: 128, isAvailable: true, nextSlot: '오늘 17:00'),
        const DoctorInfo(doctorId: 'doc-5', name: '정하늘', specialty: '가정의학', rating: 4.5, reviewCount: 97, isAvailable: false, nextSlot: '3/1 09:00'),
      ];
    }
  }

  @override
  Future<List<TimeSlot>> getDoctorAvailability(String doctorId) async {
    try {
      final res = await _client.searchFacilities(type: 'doctor');
      // 시뮬레이션: 오늘/내일 시간대
      final now = DateTime.now();
      return List.generate(6, (i) {
        final start = DateTime(now.year, now.month, now.day, 9 + i * 2);
        return TimeSlot(
          startTime: start,
          endTime: start.add(const Duration(hours: 1)),
          isAvailable: i % 3 != 0,
        );
      });
    } on DioException {
      return [];
    }
  }

  @override
  Future<bool> sendPrescriptionToPharmacy(
      String prescriptionId, String pharmacyId) async {
    try {
      await _client.selectPharmacy(prescriptionId, pharmacyId: pharmacyId);
      await _client.sendToPharmacy(prescriptionId);
      return true;
    } catch (_) {
      return false;
    }
  }

  TelemedicineReservation _mapReservation(Map<String, dynamic> m) {
    return TelemedicineReservation(
      id: m['id'] as String? ?? m['reservation_id'] as String? ?? '',
      doctorId: m['doctor_id'] as String? ?? m['facility_id'] as String? ?? '',
      doctorName: m['doctor_name'] as String? ?? m['facility_name'] as String? ?? '',
      specialty: m['specialty'] as String? ?? '',
      scheduledAt: m['scheduled_at'] != null
          ? DateTime.tryParse(m['scheduled_at'] as String) ?? DateTime.now()
          : DateTime.now(),
      durationMinutes: m['duration_minutes'] as int? ?? 30,
      status: _parseReservationStatus(m['status']),
      meetingUrl: m['meeting_url'] as String?,
    );
  }

  ReservationStatus _parseReservationStatus(dynamic v) {
    if (v is String) {
      switch (v.toLowerCase()) {
        case 'confirmed':
          return ReservationStatus.confirmed;
        case 'in_progress':
          return ReservationStatus.inProgress;
        case 'completed':
          return ReservationStatus.completed;
        case 'cancelled':
          return ReservationStatus.cancelled;
      }
    }
    return ReservationStatus.pending;
  }

  Prescription _mapPrescription(Map<String, dynamic> m) {
    return Prescription(
      id: m['id'] as String? ?? m['record_id'] as String? ?? '',
      doctorName: m['doctor_name'] as String? ?? m['provider'] as String? ?? '',
      issuedAt: m['issued_at'] != null
          ? DateTime.tryParse(m['issued_at'] as String) ?? DateTime.now()
          : m['created_at'] != null
              ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
              : DateTime.now(),
      expiresAt: m['expires_at'] != null
          ? DateTime.tryParse(m['expires_at'] as String) ?? DateTime.now().add(const Duration(days: 90))
          : DateTime.now().add(const Duration(days: 90)),
      items: _parsePrescriptionItems(m['items']),
      notes: m['notes'] as String? ?? m['description'] as String?,
    );
  }

  List<PrescriptionItem> _parsePrescriptionItems(dynamic items) {
    if (items is! List) return [];
    return items.map((i) {
      final m = i as Map<String, dynamic>;
      return PrescriptionItem(
        medicineName: m['medicine_name'] as String? ?? '',
        dosage: m['dosage'] as String? ?? '',
        frequency: m['frequency'] as String? ?? '',
        durationDays: m['duration_days'] as int? ?? 0,
      );
    }).toList();
  }

  HealthReport _mapHealthReport(Map<String, dynamic> m) {
    return HealthReport(
      id: m['id'] as String? ?? m['report_id'] as String? ?? '',
      generatedAt: m['generated_at'] != null
          ? DateTime.tryParse(m['generated_at'] as String) ?? DateTime.now()
          : m['created_at'] != null
              ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
              : DateTime.now(),
      periodDescription: m['period_description'] as String? ?? m['title'] as String? ?? '',
      analyses: _parseAnalyses(m['analyses']),
      recommendations: (m['recommendations'] as List<dynamic>?)
              ?.map((r) => r.toString())
              .toList() ??
          [],
      overallStatus: m['overall_status'] as String? ?? 'good',
    );
  }

  List<BiomarkerAnalysis> _parseAnalyses(dynamic analyses) {
    if (analyses is! List) return [];
    return analyses.map((a) {
      final m = a as Map<String, dynamic>;
      return BiomarkerAnalysis(
        biomarkerType: m['biomarker_type'] as String? ?? '',
        displayName: m['display_name'] as String? ?? '',
        latestValue: (m['latest_value'] as num?)?.toDouble() ?? 0.0,
        unit: m['unit'] as String? ?? '',
        trend: m['trend'] as String? ?? 'stable',
        status: m['status'] as String? ?? 'normal',
        advice: m['advice'] as String?,
      );
    }).toList();
  }
}
