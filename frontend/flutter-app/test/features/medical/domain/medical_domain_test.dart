import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';

void main() {
  group('ReservationStatus enum', () {
    test('값 목록 확인 (5개)', () {
      expect(ReservationStatus.values, hasLength(5));
      expect(ReservationStatus.values, contains(ReservationStatus.pending));
      expect(ReservationStatus.values, contains(ReservationStatus.confirmed));
      expect(ReservationStatus.values, contains(ReservationStatus.inProgress));
      expect(ReservationStatus.values, contains(ReservationStatus.completed));
      expect(ReservationStatus.values, contains(ReservationStatus.cancelled));
    });
  });

  group('TelemedicineReservation', () {
    test('기본 생성 확인', () {
      final reservation = TelemedicineReservation(
        id: 'res-1',
        doctorId: 'doc-1',
        doctorName: '김의사',
        specialty: '내과',
        scheduledAt: DateTime(2026, 4, 20, 14, 0),
        durationMinutes: 30,
        status: ReservationStatus.confirmed,
      );
      expect(reservation.id, 'res-1');
      expect(reservation.doctorName, '김의사');
      expect(reservation.specialty, '내과');
      expect(reservation.durationMinutes, 30);
      expect(reservation.status, ReservationStatus.confirmed);
      expect(reservation.meetingUrl, isNull);
    });

    test('meetingUrl 포함 생성', () {
      final reservation = TelemedicineReservation(
        id: 'res-2',
        doctorId: 'doc-2',
        doctorName: '이의사',
        specialty: '가정의학과',
        scheduledAt: DateTime(2026, 4, 21),
        durationMinutes: 20,
        status: ReservationStatus.inProgress,
        meetingUrl: 'https://meet.example.com/room-123',
      );
      expect(reservation.meetingUrl, isNotNull);
      expect(reservation.status, ReservationStatus.inProgress);
    });
  });

  group('Prescription', () {
    test('기본 생성 확인', () {
      final rx = Prescription(
        id: 'rx-1',
        doctorName: '김의사',
        issuedAt: DateTime(2026, 4, 19),
        expiresAt: DateTime(2026, 5, 19),
        items: const [
          PrescriptionItem(
            medicineName: '메트포르민',
            dosage: '500mg',
            frequency: '하루 2회',
            durationDays: 30,
          ),
        ],
      );
      expect(rx.id, 'rx-1');
      expect(rx.items, hasLength(1));
      expect(rx.notes, isNull);
    });

    test('notes 포함 생성', () {
      final rx = Prescription(
        id: 'rx-2',
        doctorName: '이의사',
        issuedAt: DateTime(2026, 4, 19),
        expiresAt: DateTime(2026, 5, 19),
        items: const [],
        notes: '식후 30분 복용',
      );
      expect(rx.notes, '식후 30분 복용');
    });
  });

  group('PrescriptionItem', () {
    test('생성 확인', () {
      const item = PrescriptionItem(
        medicineName: '아스피린',
        dosage: '100mg',
        frequency: '하루 1회',
        durationDays: 90,
      );
      expect(item.medicineName, '아스피린');
      expect(item.dosage, '100mg');
      expect(item.frequency, '하루 1회');
      expect(item.durationDays, 90);
    });
  });

  group('HealthReport', () {
    test('기본 생성 확인', () {
      final report = HealthReport(
        id: 'rpt-1',
        generatedAt: DateTime(2026, 3, 31),
        periodDescription: '2026년 3월 건강 리포트',
        analyses: const [],
        recommendations: const ['규칙적인 운동을 권장합니다.'],
        overallStatus: 'good',
      );
      expect(report.id, 'rpt-1');
      expect(report.recommendations, hasLength(1));
      expect(report.overallStatus, 'good');
    });

    test('분석 항목 포함 생성', () {
      final report = HealthReport(
        id: 'rpt-2',
        generatedAt: DateTime(2026, 4, 1),
        periodDescription: '2026년 3월',
        analyses: const [
          BiomarkerAnalysis(
            biomarkerType: 'blood_glucose',
            displayName: '혈당',
            latestValue: 95.0,
            unit: 'mg/dL',
            trend: 'stable',
            status: 'normal',
          ),
          BiomarkerAnalysis(
            biomarkerType: 'cholesterol',
            displayName: '콜레스테롤',
            latestValue: 210.0,
            unit: 'mg/dL',
            trend: 'increasing',
            status: 'borderline',
            advice: '식이요법 필요',
          ),
        ],
        recommendations: [],
        overallStatus: 'caution',
      );
      expect(report.analyses, hasLength(2));
      expect(report.analyses[1].advice, '식이요법 필요');
    });
  });

  group('BiomarkerAnalysis', () {
    test('기본 생성 확인', () {
      const analysis = BiomarkerAnalysis(
        biomarkerType: 'blood_pressure',
        displayName: '혈압',
        latestValue: 120.0,
        unit: 'mmHg',
        trend: 'stable',
        status: 'normal',
      );
      expect(analysis.biomarkerType, 'blood_pressure');
      expect(analysis.latestValue, 120.0);
      expect(analysis.advice, isNull);
    });
  });

  group('DoctorInfo', () {
    test('기본 생성 확인', () {
      const doctor = DoctorInfo(
        doctorId: 'doc-1',
        name: '김의사',
        specialty: '내과',
        rating: 4.8,
        reviewCount: 150,
        isAvailable: true,
      );
      expect(doctor.doctorId, 'doc-1');
      expect(doctor.rating, 4.8);
      expect(doctor.reviewCount, 150);
      expect(doctor.avatarUrl, isNull);
      expect(doctor.isAvailable, isTrue);
      expect(doctor.nextSlot, isNull);
    });

    test('nextSlot 포함 생성', () {
      const doctor = DoctorInfo(
        doctorId: 'doc-2',
        name: '이의사',
        specialty: '가정의학과',
        rating: 4.5,
        reviewCount: 80,
        isAvailable: true,
        nextSlot: '오늘 15:00',
      );
      expect(doctor.nextSlot, '오늘 15:00');
    });
  });

  group('TimeSlot', () {
    test('생성 확인', () {
      final slot = TimeSlot(
        startTime: DateTime(2026, 4, 20, 14, 0),
        endTime: DateTime(2026, 4, 20, 14, 30),
        isAvailable: true,
      );
      expect(slot.startTime.hour, 14);
      expect(slot.endTime.hour, 14);
      expect(slot.endTime.minute, 30);
      expect(slot.isAvailable, isTrue);
    });

    test('이용 불가 슬롯', () {
      final slot = TimeSlot(
        startTime: DateTime(2026, 4, 20, 10, 0),
        endTime: DateTime(2026, 4, 20, 10, 30),
        isAvailable: false,
      );
      expect(slot.isAvailable, isFalse);
    });
  });
}
