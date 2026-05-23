'use client';

import React from 'react';

type MeasurementStatus = 'normal' | 'caution' | 'danger';

export interface MeasurementResultCardProps {
  biomarker: string;
  value: number;
  unit: string;
  confidence: number;
  uncertainty: { value: number; ci_lower?: number; ci_upper?: number };
  referenceRange: { low: number; high: number };
  status: MeasurementStatus;
  timestamp: string;
  loading?: boolean;
  error?: string;
  offline?: boolean;
  onRetry?: () => void;
}

const STATUS_CONFIG: Record<MeasurementStatus, { icon: string; label: string; color: string; bg: string }> = {
  normal: { icon: '○', label: '정상', color: 'text-emerald-700', bg: 'bg-emerald-50 border-emerald-200' },
  caution: { icon: '△', label: '주의', color: 'text-amber-700', bg: 'bg-amber-50 border-amber-200' },
  danger: { icon: '⚠', label: '위험', color: 'text-red-700', bg: 'bg-red-50 border-red-200' },
};

export function MeasurementResultCard({
  biomarker, value, unit, confidence, uncertainty, referenceRange,
  status, timestamp, loading, error, offline, onRetry,
}: MeasurementResultCardProps) {
  // Loading skeleton
  if (loading) {
    return (
      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm animate-pulse" role="status" aria-label="측정 결과 로딩 중">
        <div className="h-4 w-32 bg-slate-200 rounded mb-4" />
        <div className="h-10 w-24 bg-slate-200 rounded mb-4" />
        <div className="h-3 w-full bg-slate-100 rounded mb-2" />
        <div className="h-3 w-48 bg-slate-100 rounded" />
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="rounded-2xl border border-red-200 bg-red-50 p-6 shadow-sm" role="alert" aria-live="polite">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-red-500 text-lg">⚠</span>
          <span className="text-sm font-semibold text-red-700">측정 오류</span>
        </div>
        <p className="text-sm text-red-600 mb-3">{error}</p>
        {onRetry && (
          <button onClick={onRetry} className="text-sm font-medium text-red-700 underline" aria-label="재측정">
            재측정
          </button>
        )}
      </div>
    );
  }

  const statusCfg = STATUS_CONFIG[status];
  const rangeWidth = referenceRange.high - referenceRange.low;
  const valuePosition = rangeWidth > 0
    ? Math.max(0, Math.min(100, ((value - referenceRange.low) / rangeWidth) * 100))
    : 50;

  return (
    <div
      className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
      aria-label={`${biomarker} 측정 결과: ${value} ${unit}, ${statusCfg.label}`}
    >
      {/* Header: biomarker + time + offline badge */}
      <div className="flex items-center justify-between mb-4">
        <span className="text-sm font-semibold text-slate-700">{biomarker}</span>
        <div className="flex items-center gap-2">
          {offline && (
            <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-100 text-slate-500 border border-slate-200" aria-label="오프라인 캐시">
              오프라인
            </span>
          )}
          <span className="text-xs text-slate-400">{timestamp}</span>
        </div>
      </div>

      {/* Center: large value + unit + status chip (3-triple: color + icon + text) */}
      <div className="flex items-end justify-between mb-4">
        <div>
          <span className="text-4xl font-bold text-slate-900 tabular-nums" role="status" aria-live="polite">
            {value}
          </span>
          <span className="text-lg text-slate-400 ml-1">{unit}</span>
        </div>
        <span
          className={`inline-flex items-center gap-1 px-3 py-1 rounded-full border text-sm font-semibold ${statusCfg.bg} ${statusCfg.color}`}
          role="status"
          aria-label={`상태: ${statusCfg.label}`}
        >
          <span aria-hidden="true">{statusCfg.icon}</span>
          {statusCfg.label}
        </span>
      </div>

      {/* Reference range bar */}
      <div className="mb-4" aria-label={`기준 범위: ${referenceRange.low}~${referenceRange.high} ${unit}`}>
        <div className="flex justify-between text-[10px] text-slate-400 mb-1">
          <span>{referenceRange.low}</span>
          <span className="text-slate-500 font-medium">기준 범위</span>
          <span>{referenceRange.high}</span>
        </div>
        <div className="relative h-2 bg-slate-100 rounded-full overflow-hidden">
          <div className="absolute inset-y-0 bg-emerald-100 rounded-full" style={{ left: '10%', right: '10%' }} />
          <div
            className={`absolute top-0 w-3 h-3 -mt-0.5 rounded-full border-2 border-white shadow ${
              status === 'normal' ? 'bg-emerald-500' : status === 'caution' ? 'bg-amber-500' : 'bg-red-500'
            }`}
            style={{ left: `calc(${valuePosition}% - 6px)` }}
            aria-hidden="true"
          />
        </div>
      </div>

      {/* Footer: confidence + uncertainty */}
      <div className="flex items-center justify-between pt-3 border-t border-slate-100">
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1" aria-label={`신뢰도 ${Math.round(confidence * 100)}%`}>
            <div className="w-16 h-1.5 bg-slate-100 rounded-full overflow-hidden">
              <div className="h-full bg-sky-500 rounded-full" style={{ width: `${confidence * 100}%` }} />
            </div>
            <span className="text-xs text-slate-500">{Math.round(confidence * 100)}%</span>
          </div>
        </div>
        <span className="text-xs text-slate-400" aria-label={`불확실성 ±${uncertainty.value}`}>
          ±{uncertainty.value}
          {uncertainty.ci_lower != null && uncertainty.ci_upper != null && (
            <span className="ml-1">({uncertainty.ci_lower}~{uncertainty.ci_upper})</span>
          )}
        </span>
      </div>
    </div>
  );
}
