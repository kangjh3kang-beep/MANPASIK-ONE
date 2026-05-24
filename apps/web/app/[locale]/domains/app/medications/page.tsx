'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Pill, Clock, CheckCircle, XCircle, AlertTriangle, Plus, Calendar, ExternalLink } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type MedStatus = '복용완료' | '대기중' | '건너뜀';

interface Medication {
  id: string;
  name: string;
  dosage: string;
  time: string;
  timing: string;
  status: MedStatus;
}

const STATUS_CONFIG: Record<MedStatus, { icon: React.ElementType; style: string; label: string }> = {
  '복용완료': { icon: CheckCircle, style: 'bg-emerald-50 text-emerald-700 border-emerald-200', label: '복용완료' },
  '대기중': { icon: Clock, style: 'bg-sky-50 text-sky-700 border-sky-200', label: '대기중' },
  '건너뜀': { icon: XCircle, style: 'bg-red-50 text-red-700 border-red-200', label: '건너뜀' },
};

const INITIAL_MEDICATIONS: Medication[] = [
  { id: 'med1', name: '메트포르민', dosage: '500mg', time: '08:00', timing: '아침 식후', status: '복용완료' },
  { id: 'med2', name: '아토르바스타틴', dosage: '10mg', time: '18:00', timing: '저녁 식후', status: '대기중' },
  { id: 'med3', name: '오메가3', dosage: '1000mg', time: '12:00', timing: '점심 식후', status: '복용완료' },
  { id: 'med4', name: '비타민D', dosage: '1000IU', time: '08:00', timing: '아침', status: '건너뜀' },
];

const INTERACTIONS_WARNING = {
  drugs: ['메트포르민', '아토르바스타틴'],
  message: '메트포르민과 아토르바스타틴을 함께 복용 시 드물게 근육통 증상이 나타날 수 있습니다. 이상 증상 발생 시 담당 의사와 상담하세요.',
};

const FREQUENCY_OPTIONS = ['매일', '주 2회', '주 3회', '필요 시'];

export default function MedicationsPage() {
  const localePrefix = useLocalePrefix();
  const [medications, setMedications] = useState<Medication[]>(INITIAL_MEDICATIONS);
  const [showAddForm, setShowAddForm] = useState(false);

  const handleStatusChange = (medId: string, newStatus: MedStatus) => {
    setMedications((prev) =>
      prev.map((m) => (m.id === medId ? { ...m, status: newStatus } : m)),
    );
  };

  const completedCount = medications.filter((m) => m.status === '복용완료').length;
  const adherenceRate = 87;
  const consecutiveDays = 12;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="복약 관리">
      <DomainHeader
        title="복약 관리"
        subtitle="오늘의 복약 스케줄과 순응 기록을 관리합니다"
        icon={Pill}
      />

      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 순응 통계 */}
        <section aria-label="복약 순응 통계" className="grid grid-cols-3 gap-3">
          <div className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
            <p className="text-2xl font-bold text-slate-900">{completedCount}/{medications.length}</p>
            <p className="text-xs text-slate-500 mt-1">오늘 복용</p>
          </div>
          <div className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
            <p className="text-2xl font-bold text-emerald-600">{adherenceRate}%</p>
            <p className="text-xs text-slate-500 mt-1">이번 주 순응률</p>
          </div>
          <div className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
            <p className="text-2xl font-bold text-sky-600">{consecutiveDays}일</p>
            <p className="text-xs text-slate-500 mt-1">연속 복약</p>
          </div>
        </section>

        {/* 오늘 복약 스케줄 */}
        <section className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-base font-bold text-slate-900 flex items-center gap-2">
              <Calendar className="h-4 w-4 text-sky-600" /> 오늘의 복약 스케줄
            </h2>
          </div>
          <div className="divide-y divide-slate-50">
            {medications.map((med) => {
              const statusConfig = STATUS_CONFIG[med.status];
              const StatusIcon = statusConfig.icon;
              return (
                <div key={med.id} className="px-6 py-4 flex items-center gap-4">
                  <div className="flex-shrink-0">
                    <StatusIcon
                      className={`h-5 w-5 ${
                        med.status === '복용완료'
                          ? 'text-emerald-500'
                          : med.status === '건너뜀'
                            ? 'text-red-400'
                            : 'text-sky-400'
                      }`}
                    />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-0.5">
                      <p className="text-sm font-bold text-slate-900">{med.name}</p>
                      <span className="text-xs text-slate-500">{med.dosage}</span>
                    </div>
                    <p className="text-xs text-slate-400">
                      {med.timing} ({med.time})
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <span
                      className={`px-2 py-0.5 rounded-full text-[10px] font-semibold border ${statusConfig.style}`}
                    >
                      {statusConfig.label}
                    </span>
                    {med.status === '대기중' && (
                      <div className="flex gap-1">
                        <button
                          onClick={() => handleStatusChange(med.id, '복용완료')}
                          className="px-3 py-1.5 rounded-lg bg-emerald-600 text-white text-xs font-semibold hover:bg-emerald-700 transition-colors"
                        >
                          복용
                        </button>
                        <button
                          onClick={() => handleStatusChange(med.id, '건너뜀')}
                          className="px-3 py-1.5 rounded-lg border border-slate-200 text-xs font-semibold text-slate-500 hover:bg-slate-50 transition-colors"
                        >
                          건너뜀
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </section>

        {/* 약물 상호작용 경고 */}
        <section className="rounded-2xl border border-amber-200 bg-amber-50 p-5">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-5 w-5 text-amber-600 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-bold text-amber-900">약물 상호작용 참고</p>
              <p className="text-xs text-amber-700 mt-1">{INTERACTIONS_WARNING.message}</p>
              <p className="text-[10px] text-amber-500 mt-2">
                관련 약물: {INTERACTIONS_WARNING.drugs.join(', ')} · 이 정보는 AI가 생성하였으며 의료 판단의 근거가 아닙니다. [추정]
              </p>
            </div>
          </div>
        </section>

        {/* 약 추가 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-slate-900">약 관리</h2>
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              약 추가
            </button>
          </div>

          {showAddForm && (
            <div className="p-4 rounded-xl bg-slate-50 border border-slate-100 space-y-3">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">약 이름</label>
                <input
                  type="text"
                  className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                  placeholder="약 이름을 입력하세요"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">용량</label>
                  <input
                    type="text"
                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                    placeholder="500mg"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">복용 시간</label>
                  <input
                    type="time"
                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">복용 빈도</label>
                <div className="flex flex-wrap gap-1.5">
                  {FREQUENCY_OPTIONS.map((freq) => (
                    <button
                      key={freq}
                      className="px-3 py-1.5 rounded-lg text-xs font-semibold bg-white text-slate-600 border border-slate-200 hover:bg-slate-100 transition-colors"
                    >
                      {freq}
                    </button>
                  ))}
                </div>
              </div>
              <div className="flex gap-2 pt-1">
                <button
                  onClick={() => setShowAddForm(false)}
                  className="px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-100"
                >
                  취소
                </button>
                <button className="px-5 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors">
                  저장
                </button>
              </div>
            </div>
          )}
        </section>

        {/* 처방 갱신 링크 */}
        <div className="flex justify-center">
          <Link
            href={`${localePrefix}/domains/telemedicine`}
            className="inline-flex items-center gap-1.5 text-sm font-medium text-sky-600 hover:underline"
          >
            <ExternalLink className="h-3.5 w-3.5" />
            화상진료로 처방 갱신하기
          </Link>
        </div>
      </div>
    </main>
  );
}
