'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Settings, AlertTriangle, CheckCircle, ChevronDown, RefreshCw } from 'lucide-react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

interface CalibrationRecord {
  id: string;
  date: string;
  type: '공장' | '현장';
  score: number;
  technician: string;
}

interface Device {
  id: string;
  name: string;
  lastCalibration: string;
  score: number;
  grade: string;
  nextDue: string;
  expired: boolean;
  history: CalibrationRecord[];
}

const MOCK_DEVICES: Device[] = [
  {
    id: 'dev-1',
    name: 'MPS 리더기 Pro',
    lastCalibration: '2026-05-10',
    score: 94,
    grade: 'A',
    nextDue: '2026-08-10',
    expired: false,
    history: [
      { id: 'c1', date: '2026-05-10', type: '현장', score: 94, technician: '김기술' },
      { id: 'c2', date: '2026-02-12', type: '현장', score: 91, technician: '김기술' },
      { id: 'c3', date: '2025-11-05', type: '공장', score: 98, technician: '제조팀' },
      { id: 'c4', date: '2025-08-20', type: '현장', score: 89, technician: '박보정' },
      { id: 'c5', date: '2025-05-15', type: '공장', score: 97, technician: '제조팀' },
    ],
  },
  {
    id: 'dev-2',
    name: 'MPS 리더기 Mini',
    lastCalibration: '2025-12-01',
    score: 72,
    grade: 'C',
    nextDue: '2026-03-01',
    expired: true,
    history: [
      { id: 'c6', date: '2025-12-01', type: '현장', score: 72, technician: '이보정' },
      { id: 'c7', date: '2025-09-10', type: '현장', score: 80, technician: '이보정' },
      { id: 'c8', date: '2025-06-22', type: '공장', score: 95, technician: '제조팀' },
      { id: 'c9', date: '2025-03-18', type: '현장', score: 85, technician: '이보정' },
      { id: 'c10', date: '2024-12-10', type: '공장', score: 96, technician: '제조팀' },
    ],
  },
];

function gradeColor(grade: string): string {
  if (grade === 'A+' || grade === 'A') return 'text-emerald-600';
  if (grade === 'B') return 'text-sky-600';
  if (grade === 'C') return 'text-amber-600';
  return 'text-red-600';
}

function gaugeColor(score: number): string {
  if (score >= 90) return 'bg-emerald-500';
  if (score >= 75) return 'bg-sky-500';
  if (score >= 60) return 'bg-amber-500';
  return 'bg-red-500';
}

function daysUntil(dateStr: string): number {
  const target = new Date(dateStr);
  const now = new Date('2026-05-24');
  return Math.ceil((target.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
}

export default function CalibrationPage() {
  const localePrefix = useLocalePrefix();
  const [selectedDeviceId, setSelectedDeviceId] = useState(MOCK_DEVICES[0].id);
  const [confirmOpen, setConfirmOpen] = useState(false);

  const device = MOCK_DEVICES.find((d) => d.id === selectedDeviceId) || MOCK_DEVICES[0];
  const countdown = daysUntil(device.nextDue);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="보정 관리">
      <DomainHeader
        title="보정 관리"
        subtitle="기기 보정 상태를 확인하고 재보정을 요청할 수 있습니다"
        icon={Settings}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 기기 선택 */}
        <div className="relative">
          <label htmlFor="device-select" className="block text-xs font-semibold text-slate-500 mb-1">
            기기 선택
          </label>
          <div className="relative">
            <select
              id="device-select"
              value={selectedDeviceId}
              onChange={(e) => setSelectedDeviceId(e.target.value)}
              className="w-full appearance-none rounded-2xl border border-slate-200 bg-white px-4 py-3 pr-10 text-sm font-medium text-slate-900 focus:outline-none focus:ring-2 focus:ring-sky-500"
            >
              {MOCK_DEVICES.map((d) => (
                <option key={d.id} value={d.id}>{d.name}</option>
              ))}
            </select>
            <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400 pointer-events-none" />
          </div>
        </div>

        {/* 만료 경고 배너 */}
        {device.expired && (
          <div className="flex items-center gap-3 rounded-2xl border border-red-200 bg-red-50 px-5 py-4">
            <AlertTriangle className="h-5 w-5 text-red-600 flex-shrink-0" />
            <div>
              <p className="text-sm font-bold text-red-800">보정이 만료되었습니다</p>
              <p className="text-xs text-red-600 mt-0.5">
                마지막 보정일: {device.lastCalibration} | 만료일: {device.nextDue}. 정확한 측정을 위해 즉시 재보정해 주세요.
              </p>
            </div>
          </div>
        )}

        {/* 마지막 보정 정보 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-900 mb-4">마지막 보정 결과</h2>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div>
              <p className="text-xs text-slate-500 mb-1">보정일</p>
              <p className="text-sm font-bold text-slate-900">{device.lastCalibration}</p>
            </div>
            <div>
              <p className="text-xs text-slate-500 mb-1">품질 등급</p>
              <p className={`text-2xl font-black ${gradeColor(device.grade)}`}>{device.grade}</p>
            </div>
            <div className="col-span-2">
              <p className="text-xs text-slate-500 mb-1">품질 점수</p>
              <div className="flex items-center gap-3">
                <div className="flex-1 h-3 rounded-full bg-slate-100 overflow-hidden">
                  <div
                    className={`h-full rounded-full ${gaugeColor(device.score)} transition-all`}
                    style={{ width: `${device.score}%` }}
                  />
                </div>
                <span className="text-sm font-bold text-slate-900 w-10 text-right">{device.score}</span>
              </div>
            </div>
          </div>

          {/* 다음 보정 예정 */}
          <div className="mt-5 pt-4 border-t border-slate-100 flex items-center justify-between">
            <div>
              <p className="text-xs text-slate-500">다음 보정 예정일</p>
              <p className="text-sm font-bold text-slate-900">{device.nextDue}</p>
            </div>
            <div className={`px-3 py-1.5 rounded-xl text-xs font-bold ${
              device.expired
                ? 'bg-red-50 text-red-700'
                : countdown <= 30
                  ? 'bg-amber-50 text-amber-700'
                  : 'bg-emerald-50 text-emerald-700'
            }`}>
              {device.expired ? '만료됨' : `D-${countdown}`}
            </div>
          </div>
        </div>

        {/* 품질 추이 스파크라인 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-900 mb-4">품질 추이</h2>
          <div className="flex items-end gap-2 h-20">
            {[...device.history].reverse().map((record) => (
              <div key={record.id} className="flex-1 flex flex-col items-center gap-1">
                <span className="text-[10px] font-semibold text-slate-500">{record.score}</span>
                <div
                  className={`w-full rounded-t-lg ${gaugeColor(record.score)}`}
                  style={{ height: `${(record.score / 100) * 60}px` }}
                />
                <span className="text-[9px] text-slate-400">{record.date.slice(5)}</span>
              </div>
            ))}
          </div>
        </div>

        {/* 보정 이력 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-sm font-bold text-slate-900">보정 이력</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-slate-50 text-left">
                  <th className="px-6 py-3 text-xs font-semibold text-slate-500">날짜</th>
                  <th className="px-6 py-3 text-xs font-semibold text-slate-500">유형</th>
                  <th className="px-6 py-3 text-xs font-semibold text-slate-500">점수</th>
                  <th className="px-6 py-3 text-xs font-semibold text-slate-500">담당자</th>
                </tr>
              </thead>
              <tbody>
                {device.history.map((record) => (
                  <tr key={record.id} className="border-t border-slate-50 hover:bg-slate-50">
                    <td className="px-6 py-3 text-slate-900 font-medium">{record.date}</td>
                    <td className="px-6 py-3">
                      <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold ${
                        record.type === '공장'
                          ? 'bg-sky-50 text-sky-700'
                          : 'bg-amber-50 text-amber-700'
                      }`}>
                        {record.type}
                      </span>
                    </td>
                    <td className="px-6 py-3 font-bold text-slate-900">{record.score}</td>
                    <td className="px-6 py-3 text-slate-600">{record.technician}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* 재보정 시작 버튼 + 확인 대화 */}
        <div className="flex justify-end">
          {confirmOpen ? (
            <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-lg space-y-3 w-full sm:w-96">
              <p className="text-sm font-bold text-slate-900">재보정을 시작하시겠습니까?</p>
              <p className="text-xs text-slate-500">
                보정 과정은 약 5~10분이 소요되며, 보정 중에는 측정이 불가합니다. 보정 도구 키트가 준비되어 있는지 확인해 주세요.
              </p>
              <div className="flex gap-2 justify-end">
                <button
                  onClick={() => setConfirmOpen(false)}
                  className="px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  취소
                </button>
                <button
                  onClick={() => setConfirmOpen(false)}
                  className="px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors inline-flex items-center gap-1.5"
                >
                  <RefreshCw className="h-3.5 w-3.5" />
                  보정 시작
                </button>
              </div>
            </div>
          ) : (
            <button
              onClick={() => setConfirmOpen(true)}
              className="px-5 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors inline-flex items-center gap-2"
            >
              <RefreshCw className="h-4 w-4" />
              재보정 시작
            </button>
          )}
        </div>
      </div>
    </main>
  );
}
