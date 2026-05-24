'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { FlaskConical, ChevronDown, TrendingUp, Eye } from 'lucide-react';

type LabStatus = '정상' | '주의' | '이상';

interface LabResult {
  id: string;
  testName: string;
  currentValue: number;
  unit: string;
  referenceMin: number;
  referenceMax: number;
  previousValue: number;
  changeDirection: 'up' | 'down' | 'stable';
  date: string;
  status: LabStatus;
  history: number[];
}

interface PatientOption {
  id: string;
  name: string;
  patientId: string;
}

const PATIENTS: PatientOption[] = [
  { id: 'p-1', name: '김영수', patientId: 'MPS-P001' },
  { id: 'p-2', name: '이미경', patientId: 'MPS-P002' },
  { id: 'p-3', name: '박준호', patientId: 'MPS-P003' },
];

const LAB_RESULTS: LabResult[] = [
  { id: 'l-1', testName: '혈당', currentValue: 126, unit: 'mg/dL', referenceMin: 70, referenceMax: 100, previousValue: 118, changeDirection: 'up', date: '2026-05-22', status: '주의', history: [105, 112, 118, 122, 126] },
  { id: 'l-2', testName: 'HbA1c', currentValue: 6.8, unit: '%', referenceMin: 4.0, referenceMax: 5.6, previousValue: 6.5, changeDirection: 'up', date: '2026-05-22', status: '이상', history: [6.0, 6.2, 6.5, 6.6, 6.8] },
  { id: 'l-3', testName: 'LDL', currentValue: 128, unit: 'mg/dL', referenceMin: 0, referenceMax: 130, previousValue: 135, changeDirection: 'down', date: '2026-05-22', status: '주의', history: [142, 138, 135, 132, 128] },
  { id: 'l-4', testName: 'HDL', currentValue: 52, unit: 'mg/dL', referenceMin: 40, referenceMax: 60, previousValue: 50, changeDirection: 'up', date: '2026-05-22', status: '정상', history: [45, 47, 50, 51, 52] },
  { id: 'l-5', testName: '중성지방', currentValue: 168, unit: 'mg/dL', referenceMin: 0, referenceMax: 150, previousValue: 175, changeDirection: 'down', date: '2026-05-22', status: '주의', history: [190, 182, 175, 170, 168] },
  { id: 'l-6', testName: '크레아티닌', currentValue: 1.1, unit: 'mg/dL', referenceMin: 0.7, referenceMax: 1.3, previousValue: 1.0, changeDirection: 'up', date: '2026-05-22', status: '정상', history: [0.9, 0.9, 1.0, 1.0, 1.1] },
  { id: 'l-7', testName: 'AST', currentValue: 35, unit: 'U/L', referenceMin: 0, referenceMax: 40, previousValue: 32, changeDirection: 'up', date: '2026-05-22', status: '정상', history: [28, 30, 32, 33, 35] },
  { id: 'l-8', testName: 'ALT', currentValue: 48, unit: 'U/L', referenceMin: 0, referenceMax: 41, previousValue: 42, changeDirection: 'up', date: '2026-05-22', status: '이상', history: [35, 38, 42, 45, 48] },
];

const STATUS_STYLES: Record<LabStatus, string> = {
  '정상': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '주의': 'bg-amber-50 text-amber-700 border-amber-200',
  '이상': 'bg-red-50 text-red-700 border-red-200',
};

const DIRECTION_SYMBOL: Record<string, { symbol: string; color: string }> = {
  up: { symbol: '\u2191', color: 'text-red-500' },
  down: { symbol: '\u2193', color: 'text-emerald-500' },
  stable: { symbol: '\u2192', color: 'text-slate-400' },
};

function ReferenceBar({ min, max, value, rangeStart, rangeEnd }: { min: number; max: number; value: number; rangeStart: number; rangeEnd: number }) {
  const totalRange = max - min;
  const normalStart = ((rangeStart - min) / totalRange) * 100;
  const normalWidth = ((rangeEnd - rangeStart) / totalRange) * 100;
  const valuePos = Math.max(0, Math.min(100, ((value - min) / totalRange) * 100));
  const isInRange = value >= rangeStart && value <= rangeEnd;

  return (
    <div className="relative h-4 w-full" aria-label={`기준 범위 ${rangeStart}-${rangeEnd}, 현재 값 ${value}`}>
      <div className="absolute inset-0 h-2 top-1 bg-slate-100 rounded-full" />
      <div
        className="absolute h-2 top-1 bg-emerald-200 rounded-full"
        style={{ left: `${normalStart}%`, width: `${normalWidth}%` }}
      />
      <div
        className={`absolute top-0 h-4 w-2.5 rounded-full -translate-x-1/2 transition-all ${isInRange ? 'bg-emerald-500' : 'bg-red-500'}`}
        style={{ left: `${valuePos}%` }}
      />
    </div>
  );
}

function MiniTrendChart({ data, status }: { data: number[]; status: LabStatus }) {
  const min = Math.min(...data) * 0.95;
  const max = Math.max(...data) * 1.05;
  const range = max - min || 1;
  const w = 120;
  const h = 40;
  const color = status === '정상' ? '#10b981' : status === '주의' ? '#f59e0b' : '#ef4444';
  const points = data.map((v, i) => `${(i / (data.length - 1)) * w},${h - ((v - min) / range) * h}`).join(' ');

  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} className="inline-block" aria-hidden="true">
      <polyline fill="none" stroke={color} strokeWidth="2" strokeLinejoin="round" points={points} />
      {data.map((v, i) => (
        <circle key={i} cx={(i / (data.length - 1)) * w} cy={h - ((v - min) / range) * h} r="2.5" fill={color} />
      ))}
    </svg>
  );
}

export default function LabsPage() {
  const [selectedPatient, setSelectedPatient] = useState(PATIENTS[0].id);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const currentPatient = PATIENTS.find(p => p.id === selectedPatient);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="검사 결과">
      <DomainHeader title="검사 결과" subtitle="환자별 검사 결과와 추이를 확인합니다" icon={FlaskConical} />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 환자 선택 */}
        <div className="flex items-center gap-3">
          <label htmlFor="patient-select" className="text-sm font-bold text-slate-700">환자 선택</label>
          <div className="relative">
            <select
              id="patient-select"
              value={selectedPatient}
              onChange={(e) => setSelectedPatient(e.target.value)}
              className="appearance-none pl-4 pr-10 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-medium text-slate-900 focus:outline-none focus:ring-2 focus:ring-sky-500"
            >
              {PATIENTS.map(p => (
                <option key={p.id} value={p.id}>{p.name} ({p.patientId})</option>
              ))}
            </select>
            <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400 pointer-events-none" />
          </div>
          {currentPatient && (
            <span className="text-xs text-slate-400">ID: {currentPatient.patientId}</span>
          )}
        </div>

        {/* 상태 요약 */}
        <div className="grid grid-cols-3 gap-3">
          {(['정상', '주의', '이상'] as LabStatus[]).map(status => {
            const count = LAB_RESULTS.filter(r => r.status === status).length;
            return (
              <div key={status} className={`rounded-xl border p-4 text-center ${STATUS_STYLES[status]}`}>
                <p className="text-xs font-semibold opacity-80">{status}</p>
                <p className="text-2xl font-bold mt-1">{count}항목</p>
              </div>
            );
          })}
        </div>

        {/* 검사 결과 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          {/* 헤더 */}
          <div className="hidden sm:grid grid-cols-12 gap-2 px-5 py-3 border-b border-slate-100 bg-slate-50">
            <span className="col-span-2 text-xs font-bold text-slate-500">검사명</span>
            <span className="col-span-2 text-xs font-bold text-slate-500">현재값</span>
            <span className="col-span-3 text-xs font-bold text-slate-500">기준 범위</span>
            <span className="col-span-2 text-xs font-bold text-slate-500">이전값</span>
            <span className="col-span-1 text-xs font-bold text-slate-500">상태</span>
            <span className="col-span-2 text-xs font-bold text-slate-500 text-right">추이</span>
          </div>

          <div className="divide-y divide-slate-100">
            {LAB_RESULTS.map((result) => {
              const direction = DIRECTION_SYMBOL[result.changeDirection];
              const isExpanded = expandedId === result.id;
              const barMin = Math.min(result.referenceMin * 0.7, result.currentValue * 0.8);
              const barMax = Math.max(result.referenceMax * 1.3, result.currentValue * 1.2);

              return (
                <div key={result.id}>
                  <div className="grid grid-cols-12 gap-2 items-center px-5 py-4 hover:bg-slate-50 transition-colors">
                    <div className="col-span-2">
                      <p className="text-sm font-bold text-slate-900">{result.testName}</p>
                      <p className="text-[10px] text-slate-400">{result.date}</p>
                    </div>
                    <div className="col-span-2">
                      <span className="text-sm font-bold text-slate-900 tabular-nums">{result.currentValue}</span>
                      <span className="text-xs text-slate-400 ml-1">{result.unit}</span>
                    </div>
                    <div className="col-span-3 px-1">
                      <ReferenceBar
                        min={barMin}
                        max={barMax}
                        value={result.currentValue}
                        rangeStart={result.referenceMin}
                        rangeEnd={result.referenceMax}
                      />
                      <p className="text-[10px] text-slate-400 mt-1">{result.referenceMin}~{result.referenceMax} {result.unit}</p>
                    </div>
                    <div className="col-span-2 flex items-center gap-1">
                      <span className="text-sm text-slate-500 tabular-nums">{result.previousValue}</span>
                      <span className={`text-sm font-bold ${direction.color}`}>{direction.symbol}</span>
                    </div>
                    <div className="col-span-1">
                      <span className={`inline-block px-2 py-0.5 rounded-full text-[10px] font-semibold border ${STATUS_STYLES[result.status]}`}>
                        {result.status}
                      </span>
                    </div>
                    <div className="col-span-2 flex justify-end">
                      <button
                        onClick={() => setExpandedId(isExpanded ? null : result.id)}
                        className="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-semibold text-sky-600 bg-sky-50 hover:bg-sky-100 transition-colors"
                        aria-label={`${result.testName} 추이 보기`}
                      >
                        <Eye className="h-3 w-3" />
                        추이 보기
                      </button>
                    </div>
                  </div>

                  {/* 추이 차트 */}
                  {isExpanded && (
                    <div className="px-5 pb-5 border-t border-slate-100 bg-slate-50/50">
                      <div className="pt-4">
                        <h4 className="text-xs font-bold text-slate-500 mb-3 flex items-center gap-1.5">
                          <TrendingUp className="h-3 w-3" />
                          최근 5회 검사 추이 ({result.testName})
                        </h4>
                        <div className="flex items-center gap-6">
                          <div className="flex-1">
                            <MiniTrendChart data={result.history} status={result.status} />
                          </div>
                          <div className="flex gap-3">
                            {result.history.map((val, idx) => (
                              <div key={idx} className="text-center">
                                <p className="text-xs font-bold text-slate-700 tabular-nums">{val}</p>
                                <p className="text-[9px] text-slate-400">#{idx + 1}</p>
                              </div>
                            ))}
                          </div>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* 면책 */}
        <div className="rounded-xl border border-slate-200 bg-white p-4">
          <p className="text-[10px] text-slate-400 leading-relaxed">
            검사 결과는 참고용이며, 정확한 진단은 담당 의료진과 상담하세요.
            기준 범위는 검사 기관 및 방법에 따라 다를 수 있습니다.
            모든 결과에는 측정 불확실성이 포함되어 있습니다.
          </p>
        </div>
      </div>
    </main>
  );
}
