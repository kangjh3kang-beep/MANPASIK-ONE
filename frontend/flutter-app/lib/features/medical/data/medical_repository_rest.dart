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
    // No direct cancel endpoint; placeholder
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
    // No direct emergency alert REST endpoint; placeholder
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
