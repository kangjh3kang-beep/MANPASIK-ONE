'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { FileText, ChevronDown, ChevronUp, Download, ExternalLink, Pill, CalendarCheck } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type ConsultationStatus = '진료완료' | '후속예약' | '처방발행';

interface Prescription {
  name: string;
  dosage: string;
  frequency: string;
  duration: string;
}

interface Consultation {
  id: string;
  date: string;
  doctor: string;
  specialty: string;
  diagnosisSummary: string;
  prescriptions: Prescription[];
  status: ConsultationStatus;
  aiNotes: string;
  followUpDate: string | null;
}

const STATUS_STYLES: Record<ConsultationStatus, string> = {
  진료완료: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  후속예약: 'bg-amber-50 text-amber-700 border-amber-200',
  처방발행: 'bg-blue-50 text-blue-700 border-blue-200',
};

const MOCK_CONSULTATIONS: Consultation[] = [
  {
    id: 'c1',
    date: '2026-05-20',
    doctor: '김내과',
    specialty: '내과',
    diagnosisSummary: '경도 고혈압 (1기), 공복혈당 경계치 관찰',
    status: '처방발행',
    prescriptions: [
      { name: '암로디핀 5mg', dosage: '1정', frequency: '1일 1회 아침', duration: '30일' },
      { name: '메트포르민 500mg', dosage: '1정', frequency: '1일 2회 식후', duration: '30일' },
    ],
    aiNotes:
      '환자는 최근 2주간 혈압 상승 추세를 보이며 평균 수축기 혈압 145mmHg로 측정되었습니다. ' +
      '공복혈당 112mg/dL로 경계 범위에 해당합니다. 체중 관리 및 식이요법 병행을 권고하였으며, ' +
      '4주 후 혈압과 혈당 재측정을 통해 약물 용량 조정 여부를 결정할 예정입니다. ' +
      '환자 동의 하에 MPS 측정 데이터가 공유되었으며 진료에 활용되었습니다.',
    followUpDate: '2026-06-17',
  },
  {
    id: 'c2',
    date: '2026-05-12',
    doctor: '이피부',
    specialty: '피부과',
    diagnosisSummary: '계절성 접촉피부염, 아토피 관리 상담',
    status: '후속예약',
    prescriptions: [
      { name: '히드로코르티손 크림 1%', dosage: '적량', frequency: '1일 2회 환부 도포', duration: '14일' },
      { name: '세티리진 10mg', dosage: '1정', frequency: '1일 1회 취침 전', duration: '14일' },
    ],
    aiNotes:
      '양측 전완부에 경미한 발적 및 소양감 호소하였습니다. 계절 변화에 따른 접촉피부염으로 판단되며, ' +
      '스테로이드 외용제 단기 사용 및 항히스타민제 처방하였습니다. ' +
      '보습 관리 강화를 권고하였으며 2주 후 경과 관찰 예정입니다. ' +
      '악화 시 즉시 재진을 권유하였습니다.',
    followUpDate: '2026-05-26',
  },
  {
    id: 'c3',
    date: '2026-04-28',
    doctor: '박가정',
    specialty: '가정의학과',
    diagnosisSummary: '정기 건강검진 결과 상담, 콜레스테롤 경미 상승',
    status: '진료완료',
    prescriptions: [
      { name: '오메가-3 지방산 1000mg', dosage: '1캡슐', frequency: '1일 1회 식후', duration: '90일' },
    ],
    aiNotes:
      '연간 정기 건강검진 결과를 종합 상담하였습니다. 총 콜레스테롤 215mg/dL로 경미 상승하였으나 ' +
      'LDL 135mg/dL, HDL 52mg/dL로 약물 치료 기준에는 미달합니다. ' +
      '식이 조절과 유산소 운동(주 150분 이상)을 권고하였으며, ' +
      '3개월 후 지질 검사 재측정을 계획하였습니다. 기타 항목은 정상 범위입니다.',
    followUpDate: null,
  },
];

function formatDate(dateStr: string): string {
  const d = new Date(dateStr);
  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const weekDay = ['일', '월', '화', '수', '목', '금', '토'][d.getDay()];
  return `${year}년 ${month}월 ${day}일 (${weekDay})`;
}

export default function HistoryPage() {
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const localePrefix = useLocalePrefix();

  const toggleExpand = (id: string) => {
    setExpandedId((prev) => (prev === id ? null : id));
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="진료 이력">
      <DomainHeader
        title="진료 이력"
        subtitle="과거 화상진료 기록과 처방 내역을 확인할 수 있습니다"
        icon={FileText}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8">
        {/* 요약 통계 */}
        <div className="grid grid-cols-3 gap-3 mb-8">
          {[
            { label: '총 진료 횟수', value: `${MOCK_CONSULTATIONS.length}회` },
            { label: '처방 발행', value: `${MOCK_CONSULTATIONS.filter((c) => c.status === '처방발행').length}건` },
            { label: '후속 예약', value: `${MOCK_CONSULTATIONS.filter((c) => c.status === '후속예약').length}건` },
          ].map((stat) => (
            <div key={stat.label} className="rounded-2xl bg-white border border-slate-200 p-4 text-center">
              <p className="text-2xl font-bold text-slate-800">{stat.value}</p>
              <p className="text-xs text-slate-500 mt-1">{stat.label}</p>
            </div>
          ))}
        </div>

        {/* 진료 이력 목록 */}
        <section aria-label="진료 기록 목록" className="space-y-3">
          {MOCK_CONSULTATIONS.map((consultation) => {
            const isExpanded = expandedId === consultation.id;
            return (
              <article
                key={consultation.id}
                className="rounded-2xl border border-slate-200 bg-white overflow-hidden transition-shadow hover:shadow-sm"
              >
                {/* 헤더 (항상 표시) */}
                <button
                  onClick={() => toggleExpand(consultation.id)}
                  className="w-full text-left px-5 py-4 flex items-center gap-4"
                  aria-expanded={isExpanded}
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap mb-1">
                      <span className="text-sm font-bold text-slate-900">
                        {consultation.doctor} 전문의
                      </span>
                      <span className="px-2 py-0.5 rounded-full bg-slate-100 text-[11px] font-medium text-slate-600">
                        {consultation.specialty}
                      </span>
                      <span
                        className={`px-2 py-0.5 rounded-full border text-[11px] font-semibold ${STATUS_STYLES[consultation.status]}`}
                      >
                        {consultation.status}
                      </span>
                    </div>
                    <p className="text-xs text-slate-500 mb-1.5">{formatDate(consultation.date)}</p>
                    <p className="text-sm text-slate-700 truncate">{consultation.diagnosisSummary}</p>
                  </div>

                  <div className="flex items-center gap-2">
                    <span className="text-xs text-slate-400 hidden sm:block">
                      처방 {consultation.prescriptions.length}건
                    </span>
                    {isExpanded ? (
                      <ChevronUp className="h-5 w-5 text-slate-400" />
                    ) : (
                      <ChevronDown className="h-5 w-5 text-slate-400" />
                    )}
                  </div>
                </button>

                {/* 확장 영역 */}
                {isExpanded && (
                  <div className="px-5 pb-5 border-t border-slate-100">
                    {/* AI 진료 요약 */}
                    <div className="mt-4 p-4 rounded-xl bg-slate-50">
                      <h4 className="text-xs font-bold text-slate-500 mb-2 flex items-center gap-1.5">
                        <FileText className="h-3.5 w-3.5" />
                        진료 요약 (AI 생성)
                      </h4>
                      <p className="text-sm text-slate-700 leading-relaxed">{consultation.aiNotes}</p>
                      <p className="text-[10px] text-slate-400 mt-2">
                        이 요약은 AI가 자동 생성하였으며, 의료 판단의 근거로 단독 사용할 수 없습니다.
                      </p>
                    </div>

                    {/* 처방 약물 */}
                    <div className="mt-4">
                      <h4 className="text-xs font-bold text-slate-500 mb-2 flex items-center gap-1.5">
                        <Pill className="h-3.5 w-3.5" />
                        처방 약물
                      </h4>
                      <div className="space-y-2">
                        {consultation.prescriptions.map((rx, idx) => (
                          <div
                            key={idx}
                            className="flex items-center justify-between p-3 rounded-xl bg-blue-50 border border-blue-100"
                          >
                            <div>
                              <p className="text-sm font-semibold text-slate-800">{rx.name}</p>
                              <p className="text-xs text-slate-500">
                                {rx.dosage} · {rx.frequency} · {rx.duration}
                              </p>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* 후속 예약 */}
                    {consultation.followUpDate && (
                      <div className="mt-4 flex items-center gap-2 p-3 rounded-xl bg-amber-50 border border-amber-100">
                        <CalendarCheck className="h-4 w-4 text-amber-600 flex-shrink-0" />
                        <div>
                          <p className="text-xs font-semibold text-amber-800">후속 진료 예정</p>
                          <p className="text-sm font-bold text-amber-900">{formatDate(consultation.followUpDate)}</p>
                        </div>
                      </div>
                    )}

                    {/* 하단 액션 */}
                    <div className="mt-4 flex flex-wrap gap-2">
                      <button
                        className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors"
                        onClick={() => {
                          /* 처방전 다운로드 동작 */
                        }}
                      >
                        <Download className="h-3.5 w-3.5" />
                        처방전 다운로드
                      </button>
                      <Link
                        href={`${localePrefix}/domains/health-records`}
                        className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                      >
                        <ExternalLink className="h-3.5 w-3.5" />
                        건강 기록에서 보기
                      </Link>
                    </div>
                  </div>
                )}
              </article>
            );
          })}
        </section>
      </div>
    </main>
  );
}
