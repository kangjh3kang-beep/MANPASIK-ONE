'use client';

import React from 'react';
import { Activity, AlertTriangle, CheckCircle, Info } from 'lucide-react';

export interface RiskProgressProps {
  diseaseName: string;
  riskScore: number; // 0.0 to 1.0
  factors?: string[];
  trend?: 'up' | 'down' | 'stable';
}

export function RiskProgress({ diseaseName, riskScore, factors = [], trend = 'stable' }: RiskProgressProps) {
  // 위험도에 따른 색상 및 상태 결정
  const getRiskLevel = (score: number) => {
    if (score >= 0.75) return { color: 'bg-rose-500', text: 'text-rose-500', label: '고위험', icon: AlertTriangle };
    if (score >= 0.5) return { color: 'bg-amber-500', text: 'text-amber-500', label: '주의', icon: Activity };
    if (score >= 0.25) return { color: 'bg-emerald-500', text: 'text-emerald-500', label: '정상', icon: CheckCircle };
    return { color: 'bg-blue-500', text: 'text-blue-500', label: '우수', icon: Info };
  };

  const level = getRiskLevel(riskScore);
  const percentage = Math.round(riskScore * 100);

  return (
    <div className="w-full bg-white rounded-xl p-5 border border-slate-100 shadow-sm" data-testid="risk-progress">
      <div className="flex justify-between items-end mb-3">
        <div className="flex items-center gap-2">
          <h3 className="text-base font-bold text-slate-800">{diseaseName}</h3>
          {trend === 'up' && <span className="text-rose-500 text-xs font-bold w-4">↗</span>}
          {trend === 'down' && <span className="text-emerald-500 text-xs font-bold w-4">↘</span>}
          {trend === 'stable' && <span className="text-slate-400 text-xs font-bold w-4">→</span>}
        </div>
        <div className={`flex items-center gap-1.5 text-xs font-bold px-2 py-1 rounded-full bg-slate-50 ${level.text}`}>
          <level.icon className="w-3 h-3" />
          {level.label}
        </div>
      </div>

      <div className="relative h-2.5 w-full bg-slate-100 rounded-full overflow-hidden mb-4">
        <div 
          className={`absolute top-0 left-0 h-full rounded-full transition-all duration-1000 ease-out ${level.color}`}
          style={{ width: `${percentage}%` }}
        />
      </div>

      <div className="flex justify-between items-center">
        <div className="flex flex-wrap gap-2">
          {factors.map((factor, idx) => (
            <span key={idx} className="text-[10px] font-medium text-slate-500 bg-slate-100 px-2 py-0.5 rounded">
              {factor}
            </span>
          ))}
          {factors.length === 0 && <span className="text-[10px] text-slate-400">분석된 주요 지표 없음</span>}
        </div>
        <div className={`text-sm font-black ${level.text}`}>
          {percentage}%
        </div>
      </div>
    </div>
  );
}
