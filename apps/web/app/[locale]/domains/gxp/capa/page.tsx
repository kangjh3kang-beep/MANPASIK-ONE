'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { AlertTriangle, ChevronDown, ChevronUp, Plus, CheckCircle, Circle } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type CapaType = '시정' | '예방';
type Severity = '높음' | '중간' | '낮음';
type CapaStatus = '진행중' | '완료' | '기한초과';
type RootCause = '인적오류' | '설비' | '절차' | '원자재' | '환경';

const WORKFLOW_STEPS = ['식별', '조사', '시정', '검증', '종료'] as const;

interface CapaItem {
  id: string;
  capaNumber: string;
  title: string;
  type: CapaType;
  severity: Severity;
  status: CapaStatus;
  dueDate: string;
  currentStep: number;
  rootCause: RootCause;
  assignee: string;
  description: string;
}

const SEVERITY_STYLES: Record<Severity, string> = {
  '높음': 'bg-red-50 text-red-700 border-red-200',
  '중간': 'bg-amber-50 text-amber-700 border-amber-200',
  '낮음': 'bg-slate-100 text-slate-600 border-slate-200',
};

const STATUS_STYLES: Record<CapaStatus, string> = {
  '진행중': 'bg-sky-50 text-sky-700 border-sky-200',
  '완료': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '기한초과': 'bg-red-50 text-red-700 border-red-200',
};

const TYPE_STYLES: Record<CapaType, string> = {
  '시정': 'bg-violet-50 text-violet-700 border-violet-200',
  '예방': 'bg-sky-50 text-sky-700 border-sky-200',
};

const MOCK_CAPAS: CapaItem[] = [
  {
    id: 'c1', capaNumber: 'CAPA-2026-001', title: '혈당 측정 편차 초과 건',
    type: '시정', severity: '높음', status: '진행중', dueDate: '2026-06-15',
    currentStep: 2, rootCause: '설비', assignee: '박지훈',
    description: '특정 로트 카트리지에서 혈당 측정값이 기준 대비 ±15% 초과 편차 발생',
  },
  {
    id: 'c2', capaNumber: 'CAPA-2026-002', title: '리더기 BLE 연결 간헐적 끊김',
    type: '시정', severity: '중간', status: '진행중', dueDate: '2026-06-30',
    currentStep: 1, rootCause: '설비', assignee: '오태준',
    description: 'nRF52840 펌웨어 v2.3에서 간헐적 BLE 연결 해제 현상 보고',
  },
  {
    id: 'c3', capaNumber: 'CAPA-2026-003', title: 'SOP 미준수 — 샘플 보관 온도 이탈',
    type: '예방', severity: '낮음', status: '진행중', dueDate: '2026-07-01',
    currentStep: 3, rootCause: '인적오류', assignee: '정대현',
    description: '보관 온도 로그 확인 결과 2건의 일시적 온도 이탈 확인, 예방 조치 수립',
  },
  {
    id: 'c4', capaNumber: 'CAPA-2025-018', title: '원자재 입고 검사 절차 미비',
    type: '예방', severity: '중간', status: '기한초과', dueDate: '2026-05-01',
    currentStep: 2, rootCause: '절차', assignee: '이서연',
    description: '원자재 수입검사 항목 누락으로 인한 부적합 원자재 투입 방지 절차 개선',
  },
];

export default function CapaPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [expandedCapa, setExpandedCapa] = useState<string | null>(MOCK_CAPAS[0].id);

  const stats = [
    { label: '진행중', value: MOCK_CAPAS.filter(c => c.status === '진행중').length, color: 'text-sky-600' },
    { label: '완료', value: 12, color: 'text-emerald-600' },
    { label: '기한초과', value: MOCK_CAPAS.filter(c => c.status === '기한초과').length, color: 'text-red-600' },
  ];

  return (
    <main className="min-h-screen bg-slate-50" aria-label="시정 및 예방 조치">
      <DomainHeader
        title="시정 및 예방 조치 (CAPA)"
        subtitle="CAPA 현황을 추적하고 워크플로를 관리합니다"
        icon={AlertTriangle}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* 통계 카드 */}
        <div className="grid grid-cols-3 gap-4 mb-6">
          {stats.map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
              <p className="text-xs text-slate-500 font-medium mb-1">{s.label}</p>
              <p className={`text-2xl font-bold ${s.color}`}>{s.value}</p>
            </div>
          ))}
        </div>

        {/* 새 CAPA 등록 */}
        <div className="flex justify-end mb-4">
          <button className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors">
            <Plus className="h-4 w-4" />
            새 CAPA 등록
          </button>
        </div>

        {/* CAPA 목록 */}
        <div className="space-y-3">
          {MOCK_CAPAS.map(capa => {
            const isExpanded = expandedCapa === capa.id;
            return (
              <div key={capa.id} className="rounded-2xl border border-slate-200 bg-white overflow-hidden transition-shadow hover:shadow-sm">
                <button
                  onClick={() => setExpandedCapa(prev => prev === capa.id ? null : capa.id)}
                  className="w-full flex items-center justify-between p-5 text-left"
                  aria-expanded={isExpanded}
                  aria-label={`${capa.capaNumber} 상세 보기`}
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-xs font-mono text-slate-400">{capa.capaNumber}</span>
                      <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold border ${TYPE_STYLES[capa.type]}`}>{capa.type}</span>
                      <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold border ${SEVERITY_STYLES[capa.severity]}`}>{capa.severity}</span>
                      <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold border ${STATUS_STYLES[capa.status]}`}>{capa.status}</span>
                    </div>
                    <p className="text-sm font-bold text-slate-900 truncate">{capa.title}</p>
                    <p className="text-xs text-slate-500 mt-0.5">담당: {capa.assignee} | 기한: {capa.dueDate}</p>
                  </div>
                  {isExpanded ? <ChevronUp className="h-4 w-4 text-slate-400" /> : <ChevronDown className="h-4 w-4 text-slate-400" />}
                </button>
                {isExpanded && (
                  <div className="border-t border-slate-100 px-5 pb-5 pt-4">
                    {/* 설명 */}
                    <p className="text-sm text-slate-700 mb-4">{capa.description}</p>

                    {/* 5단계 워크플로 */}
                    <div className="mb-4">
                      <h4 className="text-xs font-semibold text-slate-500 mb-3">워크플로 진행 상태</h4>
                      <div className="flex items-center justify-between relative">
                        <div className="absolute top-4 left-0 right-0 h-0.5 bg-slate-200" />
                        <div
                          className="absolute top-4 left-0 h-0.5 bg-sky-600 transition-all"
                          style={{ width: `${(capa.currentStep / (WORKFLOW_STEPS.length - 1)) * 100}%` }}
                        />
                        {WORKFLOW_STEPS.map((step, idx) => {
                          const isComplete = idx <= capa.currentStep;
                          const isCurrent = idx === capa.currentStep;
                          return (
                            <div key={step} className="relative flex flex-col items-center z-10">
                              <div className={`h-8 w-8 rounded-full flex items-center justify-center transition-colors ${
                                isCurrent
                                  ? 'bg-sky-600 text-white ring-4 ring-sky-100'
                                  : isComplete
                                    ? 'bg-sky-600 text-white'
                                    : 'bg-slate-100 text-slate-400 border border-slate-200'
                              }`}>
                                {isComplete ? <CheckCircle className="h-4 w-4" /> : <Circle className="h-4 w-4" />}
                              </div>
                              <span className={`mt-2 text-[11px] font-medium whitespace-nowrap ${isComplete ? 'text-slate-900' : 'text-slate-400'}`}>{step}</span>
                            </div>
                          );
                        })}
                      </div>
                    </div>

                    {/* 근본 원인 */}
                    <div className="flex items-center gap-4 p-3 rounded-xl bg-slate-50">
                      <div>
                        <p className="text-xs text-slate-500">근본 원인 분류</p>
                        <p className="text-sm font-semibold text-slate-900">{capa.rootCause}</p>
                      </div>
                      <div className="border-l border-slate-200 pl-4">
                        <p className="text-xs text-slate-500">담당자</p>
                        <p className="text-sm font-semibold text-slate-900">{capa.assignee}</p>
                      </div>
                      <div className="border-l border-slate-200 pl-4">
                        <p className="text-xs text-slate-500">완료 기한</p>
                        <p className={`text-sm font-semibold ${capa.status === '기한초과' ? 'text-red-600' : 'text-slate-900'}`}>{capa.dueDate}</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
