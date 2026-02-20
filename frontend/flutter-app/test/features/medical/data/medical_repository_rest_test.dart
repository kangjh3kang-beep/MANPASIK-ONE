import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/medical/data/medical_repository_rest.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('MedicalRepositoryRest', () {
    test('MedicalRepositoryRest는 MedicalRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<MedicalRepository>());
    });

    test('getReservations는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final reservations = await repo.getReservations();
      expect(reservations, isEmpty);
    });

    test('getPrescriptions는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final prescriptions = await repo.getPrescriptions();
      expect(prescriptions, isEmpty);
    });

    test('getHealthReports는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final reports = await repo.getHealthReports();
      expect(reports, isEmpty);
    });

    test('getRecommendedDoctors는 DioException 시 시뮬레이션 데이터 5명을 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final doctors = await repo.getRecommendedDoctors();
      expect(doctors, hasLength(5));
      expect(doctors.first.name, isNotEmpty);
      expect(doctors.first.specialty, isNotEmpty);
      expect(doctors.first.rating, greaterThan(0));
    });

    test('getDoctorAvailability는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final slots = await repo.getDoctorAvailability('doc-1');
      expect(slots, isEmpty);
    });

    test('sendPrescriptionToPharmacy는 실패 시 false를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MedicalRepositoryRest(client, userId: 'user-1');
      final result = await repo.sendPrescriptionToPharmacy('prsc-1', 'pharm-1');
      expect(result, isFalse);
    });
  });

  group('Medical 도메인 모델', () {
    test('TelemedicineReservation 생성 확인', () {
      final reservation = TelemedicineReservation(
        id: 'res-1',
        doctorId: 'doc-1',
        doctorName: '김의사',
        specialty: '내과',
        scheduledAt: DateTime(2026, 2, 20, 15, 0),
        durationMinutes: 30,
        status: ReservationStatus.confirmed,
      );
      expect(reservation.id, 'res-1');
      expect(reservation.status, ReservationStatus.confirmed);
      expect(reservation.meetingUrl, isNull);
      expect(reservation.durationMinutes, 30);
    });

    test('ReservationStatus enum 값 확인', () {
      expect(ReservationStatus.values, hasLength(5));
      expect(ReservationStatus.values, contains(ReservationStatus.inProgress));
    });

    test('Prescription 생성 확인', () {
      final prescription = Prescription(
        id: 'prsc-1',
        doctorName: '김의사',
        issuedAt: DateTime(2026, 2, 19),
        expiresAt: DateTime(2026, 5, 19),
        items: const [
          PrescriptionItem(
            medicineName: '아스피린',
            dosage: '100mg',
            frequency: '1일 1회',
            durationDays: 30,
          ),
        ],
        notes: '식후 복용',
      );
      expect(prescription.items, hasLength(1));
      expect(prescription.items.first.medicineName, '아스피린');
      expect(prescription.notes, '식후 복용');
    });

    test('HealthReport 생성 확인', () {
      final report = HealthReport(
        id: 'rep-1',
        generatedAt: DateTime(2026, 2, 19),
        periodDescription: '2026년 2월 건강 리포트',
        analyses: const [],
        recommendations: const ['운동 시간 늘리기', '수분 섭취 증가'],
        overallStatus: 'good',
      );
      expect(report.recommendations, hasLength(2));
      expect(report.overallStatus, 'good');
    });

    test('BiomarkerAnalysis 생성 확인', () {
      const analysis = BiomarkerAnalysis(
        biomarkerType: 'glucose',
        displayName: '혈당',
        latestValue: 95.0,
        unit: 'mg/dL',
        trend: 'stable',
        status: 'normal',
        advice: '정상 범위입니다',
      );
      expect(analysis.latestValue, 95.0);
      expect(analysis.status, 'normal');
    });

    test('DoctorInfo 생성 확인', () {
      const doctor = DoctorInfo(
        doctorId: 'doc-1',
        name: '김민수',
        specialty: '내과',
        rating: 4.8,
        reviewCount: 234,
        isAvailable: true,
        nextSlot: '오늘 15:00',
      );
      expect(doctor.rating, 4.8);
      expect(doctor.isAvailable, isTrue);
      expect(doctor.avatarUrl, isNull);
    });

    test('TimeSlot 생성 확인', () {
      final slot = TimeSlot(
        startTime: DateTime(2026, 2, 20, 10, 0),
        endTime: DateTime(2026, 2, 20, 11, 0),
        isAvailable: true,
      );
      expect(slot.isAvailable, isTrue);
      expect(slot.endTime.difference(slot.startTime).inHours, 1);
    });
  });
}
