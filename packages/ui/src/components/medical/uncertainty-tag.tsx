/**
 * @mmup-axis 3 데이터 분석
 * @mmup-stage 2 진단
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */
import React from "react";

export function UncertaintyTag({ value }: { value: string }) {
  return (
    <div 
      className="bg-red-50 text-red-700 p-2 rounded text-sm font-semibold border border-red-200"
      role="note"
      aria-label={`Uncertainty factor: ${value}`}
    >
      불확실도 (135-P5): {value}
    </div>
  );
}
