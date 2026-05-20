'use client';

import { DomainNav } from '@mmup/ui';

import React, { useState, useEffect } from 'react';
import { Pill, ClipboardCheck, FileSignature, Shield, AlertTriangle, CheckCircle, Clock, FileText } from 'lucide-react';

const AUDIT_LOGS = [
  { time: '14:32:07', user: '김품질', action: '배치 검사 승인', batch: 'BT-2026-0412', sig: true },
  { time: '14:28:45', user: '이공정', action: '원자재 입고 검수', batch: 'RM-2026-0089', sig: true },
  { time: '14:15:22', user: '박관리', action: '환경 모니터링 기록', batch: 'ENV-2026-04', sig: false },
  { time: '13:58:11', user: '최약사', action: '제품 출하 승인', batch: 'BT-2026-0411', sig: true },
  { time: '13:42:39', user: '정검사', action: '안정성 시험 결과', batch: 'ST-2026-Q2', sig: true },
];

const COMPLIANCE = [
  { rule: 'cGMP 21 CFR Part 211', status: 'pass', items: 47, passed: 47 },
  { rule: 'ICH Q7 원료의약품', status: 'pass', items: 32, passed: 32 },
  { rule: 'EU GMP Annex 11', status: 'warning', items: 28, passed: 26 },
  { rule: 'KGMP 의약품 제조관리', status: 'pass', items: 41, passed: 41 },
  { rule: 'GDP 유통관리기준', status: 'pass', items: 19, passed: 19 },
];

export default function GxPPage() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;

  return (
    <>
      <DomainNav currentDomain="gxp" />
      <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-600 text-white">
            <Pill className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">의약품 GxP 준수 시스템</h1>
        </div>
        <p className="text-sm text-slate-500">cGMP 감사 추적, 전자 서명, 규제 준수 매트릭스</p>
      </div>

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: '준수율', val: '99.4%', icon: Shield, c: 'text-emerald-600 bg-emerald-50' },
            { label: '감사항목', val: '167개', icon: ClipboardCheck, c: 'text-violet-600 bg-violet-50' },
            { label: '전자서명', val: '284건', icon: FileSignature, c: 'text-sky-600 bg-sky-50' },
            { label: '미결사항', val: '2건', icon: AlertTriangle, c: 'text-amber-600 bg-amber-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-semibold text-slate-500">{s.label}</span>
                <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>

        {/* Compliance Matrix */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">규제 준수 매트릭스</h2>
          <div className="space-y-3">
            {COMPLIANCE.map(c => (
              <div key={c.rule} className="flex items-center justify-between rounded-xl border border-slate-100 px-5 py-4">
                <div className="flex items-center gap-3">
                  {c.status === 'pass' ? <CheckCircle className="h-5 w-5 text-emerald-500" /> : <AlertTriangle className="h-5 w-5 text-amber-500" />}
                  <span className="text-sm font-semibold text-slate-900">{c.rule}</span>
                </div>
                <div className="flex items-center gap-4">
                  <div className="w-32 h-2 rounded-full bg-slate-100 overflow-hidden">
                    <div className={`h-full rounded-full ${c.status === 'pass' ? 'bg-emerald-500' : 'bg-amber-500'}`} style={{ width: `${(c.passed / c.items) * 100}%` }} />
                  </div>
                  <span className="text-sm font-bold text-slate-700 tabular-nums">{c.passed}/{c.items}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Audit Trail */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">감사 추적 로그</h2>
          <div className="space-y-2">
            {AUDIT_LOGS.map((log, i) => (
              <div key={i} className="flex items-center gap-4 rounded-xl border border-slate-100 px-4 py-3 text-sm">
                <span className="text-slate-400 font-mono text-xs tabular-nums w-16">{log.time}</span>
                <span className="font-semibold text-slate-700 w-20">{log.user}</span>
                <span className="flex-1 text-slate-600">{log.action}</span>
                <span className="text-xs bg-slate-50 text-slate-500 px-2 py-0.5 rounded font-mono">{log.batch}</span>
                {log.sig && <FileSignature className="h-4 w-4 text-violet-500" />}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
    </>
  );
}
