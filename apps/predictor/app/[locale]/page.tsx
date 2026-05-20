'use client';

import { DomainNav } from '@mmup/ui';

import React, { useState, useEffect } from 'react';
import { ActivitySquare, TrendingUp, TrendingDown, AlertTriangle, Target, Brain, BarChart3 } from 'lucide-react';

const PREDICTIONS = [
  { disease: '제2형 당뇨', risk: 0.82, trend: 'up', factors: ['HbA1c 6.8%', '공복혈당 126', 'BMI 28.3'] },
  { disease: '고혈압', risk: 0.67, trend: 'stable', factors: ['수축기 138', '가족력', 'Na 섭취'] },
  { disease: '관상동맥 질환', risk: 0.34, trend: 'down', factors: ['LDL 95', '운동 빈도 ↑', 'HDL 52'] },
  { disease: '만성 신장 질환', risk: 0.21, trend: 'stable', factors: ['GFR 88', 'Cr 1.0', '정상범위'] },
  { disease: '간경변', risk: 0.15, trend: 'down', factors: ['AST 22', 'ALT 18', '알부민 4.2'] },
];

function RiskBar({ risk }: { risk: number }) {
  const pct = Math.round(risk * 100);
  const color = risk > 0.7 ? 'bg-red-500' : risk > 0.4 ? 'bg-amber-500' : 'bg-emerald-500';
  return (
    <div className="flex items-center gap-3">
      <div className="flex-1 h-2.5 rounded-full bg-slate-100 overflow-hidden">
        <div className={`h-full rounded-full transition-all duration-1000 ${color}`} style={{ width: `${pct}%` }} />
      </div>
      <span className={`text-sm font-bold tabular-nums ${risk > 0.7 ? 'text-red-600' : risk > 0.4 ? 'text-amber-600' : 'text-emerald-600'}`}>{pct}%</span>
    </div>
  );
}


export default function PredictorPage() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;

  return (
    <>
      <DomainNav currentDomain="predictor" />
      <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-600 text-white">
            <ActivitySquare className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">생체 지표 예측 엔진</h1>
        </div>
        <p className="text-sm text-slate-500">GNN 기반 다중 질환 시계열 예측 및 위험도 분석</p>
      </div>

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-3 gap-4">
          {[
            { label: '분석 대상', val: '5개 질환', icon: Target, c: 'text-emerald-600 bg-emerald-50' },
            { label: '고위험 감지', val: '2건', icon: AlertTriangle, c: 'text-red-600 bg-red-50' },
            { label: '모델 정확도', val: '96.3%', icon: Brain, c: 'text-violet-600 bg-violet-50' },
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

        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-6">질환별 위험도 예측 결과</h2>
          <div className="space-y-6">
            {PREDICTIONS.map(pred => (
              <div key={pred.disease} className="rounded-xl border border-slate-100 p-5 hover:shadow-md transition-shadow">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <h3 className="text-sm font-bold text-slate-900">{pred.disease}</h3>
                    {pred.trend === 'up' ? <TrendingUp className="h-4 w-4 text-red-500" /> :
                     pred.trend === 'down' ? <TrendingDown className="h-4 w-4 text-emerald-500" /> :
                     <BarChart3 className="h-4 w-4 text-slate-400" />}
                  </div>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                    pred.risk > 0.7 ? 'text-red-700 bg-red-50' : pred.risk > 0.4 ? 'text-amber-700 bg-amber-50' : 'text-emerald-700 bg-emerald-50'
                  }`}>
                    {pred.risk > 0.7 ? '고위험' : pred.risk > 0.4 ? '주의' : '정상'}
                  </span>
                </div>
                <RiskBar risk={pred.risk} />
                <div className="flex gap-2 mt-3">
                  {pred.factors.map(f => (
                    <span key={f} className="text-[11px] bg-slate-50 text-slate-500 px-2 py-1 rounded-lg">{f}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
    </>
  );
}
