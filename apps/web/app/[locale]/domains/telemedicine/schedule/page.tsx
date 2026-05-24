'use client';

import React, { useState, useMemo } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Calendar, Clock, CheckCircle, X, ChevronLeft, ChevronRight } from 'lucide-react';

interface Doctor {
  id: string;
  name: string;
  specialty: string;
  photo: string;
  fee: string;
  feeAmount: number;
}

interface TimeSlot {
  time: string;
  available: boolean;
}

interface DaySchedule {
  date: string;
  dayLabel: string;
  slots: TimeSlot[];
}

const DOCTORS: Doctor[] = [
  { id: '1', name: '김내과', specialty: '내과', photo: '👨‍⚕️', fee: '₩25,000', feeAmount: 25000 },
  { id: '2', name: '이피부', specialty: '피부과', photo: '👩‍⚕️', fee: '₩30,000', feeAmount: 30000 },
  { id: '3', name: '박가정', specialty: '가정의학과', photo: '👨‍⚕️', fee: '₩20,000', feeAmount: 20000 },
];

const DAY_LABELS = ['월', '화', '수', '목', '금'];

const TIME_BLOCKS = [
  '09:00', '09:30', '10:00', '10:30', '11:00', '11:30',
  '13:00', '13:30', '14:00', '14:30', '15:00', '15:30',
  '16:00', '16:30', '17:00',
];

function generateMockSchedule(doctorId: string): DaySchedule[] {
  const today = new Date();
  const dayOfWeek = today.getDay();
  const monday = new Date(today);
  monday.setDate(today.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1));

  return Array.from({ length: 5 }, (_, dayIndex) => {
    const date = new Date(monday);
    date.setDate(monday.getDate() + dayIndex);
    const dateStr = `${date.getMonth() + 1}/${date.getDate()}`;

    const slots = TIME_BLOCKS.map((time) => {
      const seed = (parseInt(doctorId) * 7 + dayIndex * 13 + parseInt(time.replace(':', ''))) % 10;
      return { time, available: seed > 3 };
    });

    return { date: dateStr, dayLabel: DAY_LABELS[dayIndex], slots };
  });
}

type BookingState =
  | { step: 'browsing' }
  | { step: 'confirming'; dayIndex: number; slotIndex: number }
  | { step: 'success'; doctorName: string; date: string; time: string; fee: string };

export default function SchedulePage() {
  const [selectedDoctorId, setSelectedDoctorId] = useState<string>(DOCTORS[0].id);
  const [bookingState, setBookingState] = useState<BookingState>({ step: 'browsing' });

  const selectedDoctor = DOCTORS.find((d) => d.id === selectedDoctorId)!;
  const weekSchedule = useMemo(() => generateMockSchedule(selectedDoctorId), [selectedDoctorId]);

  const handleSlotClick = (dayIndex: number, slotIndex: number) => {
    const slot = weekSchedule[dayIndex].slots[slotIndex];
    if (!slot.available) return;
    setBookingState({ step: 'confirming', dayIndex, slotIndex });
  };

  const handleConfirm = () => {
    if (bookingState.step !== 'confirming') return;
    const day = weekSchedule[bookingState.dayIndex];
    const slot = day.slots[bookingState.slotIndex];
    setBookingState({
      step: 'success',
      doctorName: selectedDoctor.name,
      date: `${day.date} (${day.dayLabel})`,
      time: slot.time,
      fee: selectedDoctor.fee,
    });
  };

  const handleReset = () => {
    setBookingState({ step: 'browsing' });
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="진료 예약">
      <DomainHeader
        title="진료 예약"
        subtitle="원하는 의사와 시간을 선택하여 화상진료를 예약하세요"
        icon={Calendar}
      />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8">
        {/* 의사 선택 */}
        <section aria-label="의사 선택" className="mb-8">
          <h2 className="text-sm font-bold text-slate-500 mb-3">담당 의사 선택</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {DOCTORS.map((doctor) => {
              const isSelected = doctor.id === selectedDoctorId;
              return (
                <button
                  key={doctor.id}
                  onClick={() => {
                    setSelectedDoctorId(doctor.id);
                    setBookingState({ step: 'browsing' });
                  }}
                  className={`flex items-center gap-3 p-4 rounded-2xl border-2 text-left transition-all ${
                    isSelected
                      ? 'border-sky-500 bg-sky-50 shadow-sm'
                      : 'border-slate-200 bg-white hover:border-slate-300'
                  }`}
                  aria-pressed={isSelected}
                >
                  <span className="text-3xl">{doctor.photo}</span>
                  <div className="min-w-0">
                    <p className={`text-sm font-bold ${isSelected ? 'text-sky-700' : 'text-slate-800'}`}>
                      {doctor.name} 전문의
                    </p>
                    <p className="text-xs text-slate-500">{doctor.specialty}</p>
                    <p className="text-xs font-semibold text-slate-600 mt-0.5">{doctor.fee}</p>
                  </div>
                </button>
              );
            })}
          </div>
        </section>

        {/* 주간 시간표 */}
        <section aria-label="주간 예약 가능 시간" className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="flex items-center justify-between px-5 py-4 border-b border-slate-100">
            <h2 className="text-base font-bold text-slate-800 flex items-center gap-2">
              <Clock className="h-4 w-4 text-sky-600" />
              이번 주 예약 가능 시간
            </h2>
            <p className="text-xs text-slate-400">30분 단위 · 점심시간 12:00~13:00 제외</p>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full min-w-[600px]">
              <thead>
                <tr className="border-b border-slate-100">
                  <th className="w-20 px-3 py-3 text-xs font-semibold text-slate-400 text-left">시간</th>
                  {weekSchedule.map((day) => (
                    <th key={day.date} className="px-2 py-3 text-center">
                      <span className="block text-xs font-bold text-slate-700">{day.dayLabel}</span>
                      <span className="block text-[10px] text-slate-400">{day.date}</span>
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {TIME_BLOCKS.map((time, slotIndex) => (
                  <tr key={time} className="border-b border-slate-50 last:border-b-0">
                    <td className="px-3 py-1.5 text-xs text-slate-500 font-medium">{time}</td>
                    {weekSchedule.map((day, dayIndex) => {
                      const slot = day.slots[slotIndex];
                      const isConfirming =
                        bookingState.step === 'confirming' &&
                        bookingState.dayIndex === dayIndex &&
                        bookingState.slotIndex === slotIndex;
                      return (
                        <td key={`${dayIndex}-${slotIndex}`} className="px-1.5 py-1">
                          <button
                            onClick={() => handleSlotClick(dayIndex, slotIndex)}
                            disabled={!slot.available || bookingState.step === 'success'}
                            className={`w-full py-1.5 rounded-lg text-[11px] font-medium transition-all ${
                              isConfirming
                                ? 'bg-sky-600 text-white ring-2 ring-sky-300'
                                : slot.available
                                  ? 'bg-emerald-50 text-emerald-700 hover:bg-emerald-100 hover:shadow-sm'
                                  : 'bg-slate-50 text-slate-300 cursor-not-allowed'
                            }`}
                            aria-label={`${day.dayLabel} ${time} ${slot.available ? '예약 가능' : '예약 불가'}`}
                          >
                            {isConfirming ? '선택됨' : slot.available ? '가능' : '—'}
                          </button>
                        </td>
                      );
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="flex items-center gap-4 px-5 py-3 border-t border-slate-100 bg-slate-50">
            <span className="flex items-center gap-1.5 text-[11px] text-slate-500">
              <span className="inline-block h-3 w-6 rounded bg-emerald-50 border border-emerald-200" />
              예약 가능
            </span>
            <span className="flex items-center gap-1.5 text-[11px] text-slate-500">
              <span className="inline-block h-3 w-6 rounded bg-slate-50 border border-slate-200" />
              예약 불가
            </span>
          </div>
        </section>

        {/* 예약 확인 모달 */}
        {bookingState.step === 'confirming' && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
            <div
              className="w-full max-w-md mx-4 rounded-2xl bg-white shadow-xl"
              role="dialog"
              aria-label="예약 확인"
            >
              <div className="flex items-center justify-between px-6 pt-6 pb-2">
                <h3 className="text-lg font-bold text-slate-900">예약 확인</h3>
                <button
                  onClick={handleReset}
                  className="h-8 w-8 rounded-full hover:bg-slate-100 flex items-center justify-center"
                  aria-label="닫기"
                >
                  <X className="h-4 w-4 text-slate-400" />
                </button>
              </div>

              <div className="px-6 py-4 space-y-4">
                <div className="flex items-center gap-3 p-4 rounded-xl bg-slate-50">
                  <span className="text-3xl">{selectedDoctor.photo}</span>
                  <div>
                    <p className="font-bold text-slate-900">{selectedDoctor.name} 전문의</p>
                    <p className="text-sm text-slate-500">{selectedDoctor.specialty}</p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="p-3 rounded-xl bg-sky-50">
                    <p className="text-[10px] font-semibold text-sky-500 mb-0.5">날짜</p>
                    <p className="text-sm font-bold text-sky-900">
                      {weekSchedule[bookingState.dayIndex].date} ({weekSchedule[bookingState.dayIndex].dayLabel})
                    </p>
                  </div>
                  <div className="p-3 rounded-xl bg-sky-50">
                    <p className="text-[10px] font-semibold text-sky-500 mb-0.5">시간</p>
                    <p className="text-sm font-bold text-sky-900">
                      {weekSchedule[bookingState.dayIndex].slots[bookingState.slotIndex].time}
                    </p>
                  </div>
                </div>

                <div className="flex items-center justify-between p-3 rounded-xl bg-amber-50">
                  <span className="text-sm text-amber-800">진료비</span>
                  <span className="text-base font-bold text-amber-900">{selectedDoctor.fee}</span>
                </div>

                <p className="text-xs text-slate-400 text-center">
                  예약 확정 후 진료 시작 10분 전에 알림을 보내드립니다
                </p>
              </div>

              <div className="flex gap-3 px-6 pb-6">
                <button
                  onClick={handleReset}
                  className="flex-1 px-4 py-3 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50"
                >
                  취소
                </button>
                <button
                  onClick={handleConfirm}
                  className="flex-1 px-4 py-3 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
                >
                  예약 확정
                </button>
              </div>
            </div>
          </div>
        )}

        {/* 예약 성공 */}
        {bookingState.step === 'success' && (
          <section className="mt-8 rounded-2xl border-2 border-emerald-200 bg-emerald-50 p-8 text-center">
            <div className="h-14 w-14 rounded-full bg-emerald-100 flex items-center justify-center mx-auto mb-4">
              <CheckCircle className="h-7 w-7 text-emerald-600" />
            </div>
            <h2 className="text-lg font-bold text-emerald-900 mb-2">예약이 확정되었습니다</h2>
            <div className="inline-flex flex-col gap-1 text-sm text-emerald-800 mb-6">
              <p><strong>의사:</strong> {bookingState.doctorName} 전문의</p>
              <p><strong>날짜:</strong> {bookingState.date}</p>
              <p><strong>시간:</strong> {bookingState.time}</p>
              <p><strong>진료비:</strong> {bookingState.fee}</p>
            </div>
            <div className="flex justify-center gap-3">
              <button
                onClick={handleReset}
                className="px-6 py-2.5 rounded-xl border border-emerald-300 text-sm font-semibold text-emerald-700 hover:bg-emerald-100"
              >
                추가 예약
              </button>
            </div>
          </section>
        )}
      </div>
    </main>
  );
}
