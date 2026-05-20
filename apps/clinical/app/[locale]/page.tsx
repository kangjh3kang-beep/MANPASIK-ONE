/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @sb SB-1
 * @standard IEC 62304 Class B
 *
 * 임상 데이터 콘솔 — 메인 대시보드
 * 의사/연구원 전용 워크스테이션
 */
'use client';

import React from 'react';
import { DomainNav } from '@mmup/ui';
import { ClinicalShell } from '../../components/clinical-shell';
import { PatientTimeline } from '../../components/patient-timeline';
import { BiomarkerHeatmap } from '../../components/biomarker-heatmap';
import { SensorStatusGrid } from '../../components/sensor-status-grid';
import { Activity, Users, AlertTriangle, TrendingUp } from 'lucide-react';

const STAT_CARDS = [
  { label: '활성 환자', value: '1,247', change: '+12', icon: Users, color: 'text-sky-600 bg-sky-50' },
  { label: '금일 측정', value: '3,891', change: '+284', icon: Activity, color: 'text-emerald-600 bg-emerald-50' },
  { label: '이상 징후', value: '23', change: '-4', icon: AlertTriangle, color: 'text-amber-600 bg-amber-50' },
  { label: '진단 정확도', value: '99.7%', change: '+0.2%', icon: TrendingUp, color: 'text-violet-600 bg-violet-50' },
];

export default function ClinicalDashboard() {
  return (
    <>
      <DomainNav currentDomain="clinical" />
      <ClinicalShell>
      {/* Page Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-900">임상 대시보드</h1>
        <p className="text-sm text-slate-500 mt-1">환자 생체 데이터 실시간 모니터링 및 AI 분석 현황</p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {STAT_CARDS.map((stat) => (
          <div key={stat.label} className="rounded-2xl border border-slate-200 bg-white p-5">
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs font-semibold text-slate-500">{stat.label}</span>
              <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${stat.color}`}>
                <stat.icon className="h-4 w-4" />
              </div>
            </div>
            <p className="text-2xl font-bold text-slate-900 tabular-nums">{stat.value}</p>
            <p className="text-xs text-emerald-600 font-medium mt-1">{stat.change} 전일 대비</p>
          </div>
        ))}
      </div>

      {/* Patient Vitals Timeline */}
      <div className="rounded-2xl border border-slate-200 bg-white p-6 mb-6">
        <PatientTimeline />
      </div>

      {/* Biomarker Heatmap */}
      <div className="rounded-2xl border border-slate-200 bg-white p-6 mb-6">
        <BiomarkerHeatmap />
      </div>

      {/* Sensor Status Grid */}
      <div className="rounded-2xl border border-slate-200 bg-white p-6">
        <SensorStatusGrid />
      </div>
    </ClinicalShell>
    </>
  );
}
