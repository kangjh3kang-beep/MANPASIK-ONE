'use client';

import React, { useState, useEffect } from 'react';

const BIOMARKERS = [
  { id: 'hgb', name: '헤모글로빈', unit: 'g/dL', min: 12, max: 17.5, category: '혈액' },
  { id: 'wbc', name: '백혈구', unit: 'K/μL', min: 4.5, max: 11.0, category: '혈액' },
  { id: 'plt', name: '혈소판', unit: 'K/μL', min: 150, max: 400, category: '혈액' },
  { id: 'rbc', name: '적혈구', unit: 'M/μL', min: 4.5, max: 5.9, category: '혈액' },
  { id: 'glu', name: '혈당', unit: 'mg/dL', min: 70, max: 100, category: '대사' },
  { id: 'hba1c', name: 'HbA1c', unit: '%', min: 4.0, max: 5.7, category: '대사' },
  { id: 'chol', name: '총콜레스테롤', unit: 'mg/dL', min: 0, max: 200, category: '대사' },
  { id: 'ldl', name: 'LDL', unit: 'mg/dL', min: 0, max: 130, category: '대사' },
  { id: 'hdl', name: 'HDL', unit: 'mg/dL', min: 40, max: 100, category: '대사' },
  { id: 'tg', name: '중성지방', unit: 'mg/dL', min: 0, max: 150, category: '대사' },
  { id: 'ast', name: 'AST', unit: 'U/L', min: 0, max: 40, category: '간기능' },
  { id: 'alt', name: 'ALT', unit: 'U/L', min: 0, max: 41, category: '간기능' },
  { id: 'alb', name: '알부민', unit: 'g/dL', min: 3.4, max: 5.4, category: '간기능' },
  { id: 'bun', name: 'BUN', unit: 'mg/dL', min: 6, max: 20, category: '신장' },
  { id: 'cr', name: '크레아티닌', unit: 'mg/dL', min: 0.7, max: 1.3, category: '신장' },
  { id: 'na', name: '나트륨', unit: 'mEq/L', min: 136, max: 145, category: '전해질' },
  { id: 'k', name: '칼륨', unit: 'mEq/L', min: 3.5, max: 5.1, category: '전해질' },
  { id: 'ca', name: '칼슘', unit: 'mg/dL', min: 8.6, max: 10.2, category: '전해질' },
  { id: 'crp', name: 'CRP', unit: 'mg/L', min: 0, max: 3.0, category: '염증' },
  { id: 'esr', name: 'ESR', unit: 'mm/hr', min: 0, max: 20, category: '염증' },
];

function getRiskLevel(value: number, min: number, max: number): 'normal' | 'warning' | 'critical' {
  if (value >= min && value <= max) return 'normal';
  const diff = value < min ? min - value : value - max;
  const range = max - min || 1;
  return diff / range > 0.3 ? 'critical' : 'warning';
}

const RISK_COLORS = {
  normal: 'bg-emerald-100 text-emerald-700 border-emerald-200',
  warning: 'bg-amber-100 text-amber-700 border-amber-200',
  critical: 'bg-red-100 text-red-700 border-red-200 animate-pulse',
};

const RISK_DOT = {
  normal: 'bg-emerald-500',
  warning: 'bg-amber-500',
  critical: 'bg-red-500',
};

export function BiomarkerHeatmap() {
  const [values, setValues] = useState<Record<string, number>>({});
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    const generated: Record<string, number> = {};
    BIOMARKERS.forEach((bm) => {
      const range = bm.max - bm.min;
      // 80% chance normal, 15% warning, 5% critical
      const rand = Math.random();
      if (rand < 0.8) {
        generated[bm.id] = +(bm.min + Math.random() * range).toFixed(1);
      } else if (rand < 0.95) {
        generated[bm.id] = +(bm.max + Math.random() * range * 0.2).toFixed(1);
      } else {
        generated[bm.id] = +(bm.max + Math.random() * range * 0.5).toFixed(1);
      }
    });
    setValues(generated);
  }, []);

  if (!mounted) return <div className="h-64 animate-pulse rounded-2xl bg-slate-100" />;

  const categories = [...new Set(BIOMARKERS.map((b) => b.category))];
  const criticalCount = BIOMARKERS.filter((bm) => getRiskLevel(values[bm.id] ?? 0, bm.min, bm.max) === 'critical').length;
  const warningCount = BIOMARKERS.filter((bm) => getRiskLevel(values[bm.id] ?? 0, bm.min, bm.max) === 'warning').length;

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-bold text-slate-900">바이오마커 히트맵</h2>
        <div className="flex items-center gap-4 text-xs">
          <span className="flex items-center gap-1.5"><span className="h-2.5 w-2.5 rounded-full bg-emerald-500" /> 정상 {BIOMARKERS.length - criticalCount - warningCount}</span>
          <span className="flex items-center gap-1.5"><span className="h-2.5 w-2.5 rounded-full bg-amber-500" /> 주의 {warningCount}</span>
          <span className="flex items-center gap-1.5"><span className="h-2.5 w-2.5 rounded-full bg-red-500" /> 위험 {criticalCount}</span>
        </div>
      </div>

      <div className="space-y-6">
        {categories.map((cat) => (
          <div key={cat}>
            <h3 className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">{cat}</h3>
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2">
              {BIOMARKERS.filter((b) => b.category === cat).map((bm) => {
                const val = values[bm.id] ?? 0;
                const risk = getRiskLevel(val, bm.min, bm.max);
                return (
                  <div
                    key={bm.id}
                    className={`rounded-xl border px-3 py-2.5 transition-all cursor-default hover:shadow-md ${RISK_COLORS[risk]}`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[11px] font-semibold">{bm.name}</span>
                      <span className={`h-2 w-2 rounded-full ${RISK_DOT[risk]}`} />
                    </div>
                    <div className="flex items-end gap-1">
                      <span className="text-lg font-bold tabular-nums">{val}</span>
                      <span className="text-[10px] mb-0.5 opacity-70">{bm.unit}</span>
                    </div>
                    <div className="text-[10px] mt-0.5 opacity-60">
                      기준: {bm.min}–{bm.max}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
