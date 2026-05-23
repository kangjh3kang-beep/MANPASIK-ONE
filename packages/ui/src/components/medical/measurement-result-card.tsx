'use client';

import React, { useState } from 'react';

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

      {/* Footer: confidence + 3단계 설명 */}
      <ThreeTierExplanation
        biomarker={biomarker}
        value={value}
        unit={unit}
        confidence={confidence}
        uncertainty={uncertainty}
        referenceRange={referenceRange}
        status={status}
      />
    </div>
  );
}

/** 3단계 설명 계층 — 일반인 → 중급 → 전문가 (점진적 공개) */
function ThreeTierExplanation({ biomarker, value, unit, confidence, uncertainty, referenceRange, status }: {
  biomarker: string; value: number; unit: string; confidence: number;
  uncertainty: { value: number; ci_lower?: number; ci_upper?: number };
  referenceRange: { low: number; high: number }; status: MeasurementStatus;
}) {
  const [tier, setTier] = useState<0 | 1 | 2>(0);

  const statusText = status === 'normal' ? '정상 범위 안에 있습니다' : status === 'caution' ? '약간 주의가 필요합니다' : '관리가 필요한 수준입니다';

  return (
    <div className="pt-3 border-t border-slate-100">
      {/* 일반인 (기본, 항상 표시) */}
      <p className="text-sm text-slate-600 mb-2">
        {biomarker} 수치가 {statusText}.
      </p>

      {/* 자세히 버튼 */}
      {tier === 0 && (
        <button onClick={() => setTier(1)} className="text-xs font-semibold text-sky-600 hover:text-sky-700 transition-colors">
          자세히 보기 ▾
        </button>
      )}

      {/* 중급 (수치·기준·추세) */}
      {tier >= 1 && (
        <div className="mt-2 p-3 rounded-lg bg-slate-50 text-xs text-slate-600 space-y-1">
          <p>측정값: <strong>{value} {unit}</strong></p>
          <p>정상 범위: {referenceRange.low}~{referenceRange.high} {unit}</p>
          <p>신뢰도: {Math.round(confidence * 100)}% · 오차: ±{uncertainty.value}</p>
          {tier === 1 && (
            <button onClick={() => setTier(2)} className="text-sky-600 font-semibold mt-1 hover:text-sky-700">
              전문가 정보 ▾
            </button>
          )}
        </div>
      )}

      {/* 전문가 (방법·계량·근거) */}
      {tier === 2 && (
        <div className="mt-2 p-3 rounded-lg bg-sky-50 text-xs text-slate-600 space-y-1">
          <p>측정 방법: 차동 보정 (Sdiff = S - α×R, α=0.98)</p>
          <p>신뢰구간: {uncertainty.ci_lower ?? '—'}~{uncertainty.ci_upper ?? '—'} {unit} (95% CI)</p>
          <p>불확실성: ±{uncertainty.value} (GUM Type A+B 합성)</p>
          <p className="text-slate-400 mt-1">이 결과는 AI 참고 정보이며, 의료 진단을 대체하지 않습니다.</p>
          <button onClick={() => setTier(0)} className="text-sky-600 font-semibold mt-1 hover:text-sky-700">
            접기 ▴
          </button>
        </div>
      )}
    </div>
  );
}
