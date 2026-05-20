/**
 * @mmup-axis 3 종합 데이터 분석
 * @mmup-stage 2 진단
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
import React from "react";
import { cn } from "../../lib/utils";

interface ConfidenceBadgeProps {
  score: number; // 0.0 - 1.0
  className?: string;
}

export function ConfidenceBadge({ score, className }: ConfidenceBadgeProps) {
  const isHigh = score >= 0.9;
  return (
    <span
      className={cn(
        "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium",
        isHigh ? "bg-green-100 text-green-800" : "bg-yellow-100 text-yellow-800",
        className
      )}
      role="status"
      aria-label={`Confidence score: ${(score * 100).toFixed(1)}%`}
    >
      {(score * 100).toFixed(1)}% 신뢰도
    </span>
  );
}
