'use client';

import React, { useState } from 'react';

export interface ThreeTierExplanationProps {
  /** 일반인 설명 (항상 표시) */
  laymanText: string;
  /** 중급 설명 데이터 */
  intermediateData?: { label: string; value: string }[];
  /** 전문가 설명 데이터 */
  expertData?: { label: string; value: string }[];
  /** AI 면책 고지 표시 여부 */
  showDisclaimer?: boolean;
}

/**
 * 3단계 설명 계층 — 일반인 → 중급 → 전문가 (점진적 공개)
 *
 * 벤치마크: K Health (점진적 공개), Abbott FreeStyle (AGP 상세)
 * 규제: IEC 62366 use error 방지 — 모든 문해력 수준에서 이해 가능
 */
export function ThreeTierExplanation({
  laymanText,
  intermediateData,
  expertData,
  showDisclaimer = true,
}: ThreeTierExplanationProps) {
  const [tier, setTier] = useState<0 | 1 | 2>(0);

  return (
    <div className="pt-3 border-t border-slate-100">
      {/* 일반인 (기본, 항상 표시) */}
      <p className="text-sm text-slate-600 mb-2">{laymanText}</p>

      {/* 자세히 버튼 */}
      {tier === 0 && intermediateData && (
        <button onClick={() => setTier(1)} className="text-xs font-semibold text-sky-600 hover:text-sky-700 transition-colors">
          자세히 보기 ▾
        </button>
      )}

      {/* 중급 */}
      {tier >= 1 && intermediateData && (
        <div className="mt-2 p-3 rounded-lg bg-slate-50 text-xs text-slate-600 space-y-1">
          {intermediateData.map((item, i) => (
            <p key={i}>{item.label}: <strong>{item.value}</strong></p>
          ))}
          {tier === 1 && expertData && (
            <button onClick={() => setTier(2)} className="text-sky-600 font-semibold mt-1 hover:text-sky-700">
              전문가 정보 ▾
            </button>
          )}
        </div>
      )}

      {/* 전문가 */}
      {tier === 2 && expertData && (
        <div className="mt-2 p-3 rounded-lg bg-sky-50 text-xs text-slate-600 space-y-1">
          {expertData.map((item, i) => (
            <p key={i}>{item.label}: <strong>{item.value}</strong></p>
          ))}
          {showDisclaimer && (
            <p className="text-slate-400 mt-1">이 결과는 AI 참고 정보이며, 의료 진단을 대체하지 않습니다.</p>
          )}
          <button onClick={() => setTier(0)} className="text-sky-600 font-semibold mt-1 hover:text-sky-700">
            접기 ▴
          </button>
        </div>
      )}
    </div>
  );
}
