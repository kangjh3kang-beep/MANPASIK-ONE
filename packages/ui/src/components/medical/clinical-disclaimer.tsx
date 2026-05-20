/**
 * @mmup-axis 4 종합병원
 * @mmup-stage 4 처방
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
import React from "react";

export function ClinicalDisclaimer() {
  return (
    <div 
      className="text-xs text-gray-500 mt-4 p-3 bg-gray-50 border-t border-gray-200"
      role="alert"
    >
      <p>본 정보는 의료 자문 보조이며 진료 행위를 대체하지 않습니다. 전문가와 상담하십시오.</p>
    </div>
  );
}
