/// 의료 서비스 도메인 모델 및 리포지토리
///
/// 비대면 진료, 처방전, 건강 리포트, 응급 알림

/// 비대면 진료 예약
class TelemedicineReservation {
  final String id;
  final String doctorId;
  final String doctorName;
  final String specialty;
  final DateTime scheduledAt;
  final int durationMinutes;
  final ReservationStatus status;
  final String? meetingUrl;

  const TelemedicineReservation({
    required this.id,
    required this.doctorId,
    required this.doctorName,
    required this.specialty,
    required this.scheduledAt,
    required this.durationMinutes,
    required this.status,
    this.meetingUrl,
  });
}

/// 예약 상태
enum ReservationStatus { pending, confirmed, inProgress, completed, cancelled }

/// 처방전
class Prescription {
  final String id;
  final String doctorName;
  final DateTime issuedAt;
  final DateTime expiresAt;
  final List<PrescriptionItem> items;
  final String? notes;

  const Prescription({
    required this.id,
    required this.doctorName,
    required this.issuedAt,
    required this.expiresAt,
    required this.items,
    this.notes,
  });
}

/// 처방 항목
class PrescriptionItem {
  final String medicineName;
  final String dosage;
  final String frequency;
  final int durationDays;

  const PrescriptionItem({
    required this.medicineName,
    required this.dosage,
    required this.frequency,
    required this.durationDays,
  });
}

/// 건강 리포트
class HealthReport {
  final String id;
  final DateTime generatedAt;
  final String periodDescription;
  final List<BiomarkerAnalysis> analyses;
  final List<String> recommendations;
  final String overallStatus; // 'excellent', 'good', 'caution', 'alert'

  const HealthReport({
    required this.id,
    required this.generatedAt,
    required this.periodDescription,
    required this.analyses,
    required this.recommendations,
    required this.overallStatus,
  });
}

/// 바이오마커 분석 결과
class BiomarkerAnalysis {
  final String biomarkerType;
  final String displayName;
  final double latestValue;
  final String unit;
  final String trend;
  final String status; // 'normal', 'borderline', 'abnormal'
  final String? advice;

  const BiomarkerAnalysis({
    required this.biomarkerType,
    required this.displayName,
    required this.latestValue,
    required this.unit,
    required this.trend,
    required this.status,
    this.advice,
  });
}

/// 의사 정보 (추천 목록용)
class DoctorInfo {
  final String doctorId;
  final String name;
  final String specialty;
  final double rating;
  final int reviewCount;
  final String? avatarUrl;
  final bool isAvailable;
  final String? nextSlot; // "오늘 15:00" 등

  const DoctorInfo({
    required this.doctorId,
    required this.name,
    required this.specialty,
    required this.rating,
    required this.reviewCount,
    this.avatarUrl,
    required this.isAvailable,
    this.nextSlot,
  });
}

/// 진료 시간 슬롯
class TimeSlot {
  final DateTime startTime;
  final DateTime endTime;
  final bool isAvailable;

  const TimeSlot({
    required this.startTime,
    required this.endTime,
    required this.isAvailable,
  });
}

/// 의료 서비스 리포지토리 인터페이스
abstract class MedicalRepository {
  /// 비대면 진료 예약 생성
  Future<TelemedicineReservation> createReservation({
    required String doctorId,
    required DateTime scheduledAt,
  });

  /// 예약 목록 조회
  Future<List<TelemedicineReservation>> getReservations();

  /// 예약 취소
  Future<void> cancelReservation(String reservationId);

  /// 처방전 목록
  Future<List<Prescription>> getPrescriptions();

  /// 처방전 상세
  Future<Prescription> getPrescriptionDetail(String prescriptionId);

  /// 건강 리포트 생성 요청
  Future<HealthReport> generateHealthReport();

  /// 건강 리포트 목록
  Future<List<HealthReport>> getHealthReports();

  /// 응급 알림 전송
  Future<void> sendEmergencyAlert({
    required String alertType,
    required double abnormalValue,
    required String biomarkerType,
  });

  /// 추천 의사 목록
  Future<List<DoctorInfo>> getRecommendedDoctors({int limit = 5});

  /// 의사 가용성 확인
  Future<List<TimeSlot>> getDoctorAvailability(String doctorId);

  /// 처방전 약국 전송
  Future<bool> sendPrescriptionToPharmacy(
      String prescriptionId, String pharmacyId);
}
