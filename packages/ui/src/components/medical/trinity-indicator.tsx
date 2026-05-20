/**
 * @mmup-axis 1 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard WCAG 2.2 AA
 */
import React from "react";
import { IconIP1, IconIP2, IconIP3 } from "../../icons/trinity";

export function TrinityIndicator({ active = [] }: { active?: ("IP1" | "IP2" | "IP3")[] }) {
  return (
    <div className="flex gap-2 items-center" aria-label="Trinity IP indicator">
      <div className={active.includes("IP1") ? "text-blue-600" : "text-gray-300"}><IconIP1 /></div>
      <div className={active.includes("IP2") ? "text-blue-600" : "text-gray-300"}><IconIP2 /></div>
      <div className={active.includes("IP3") ? "text-blue-600" : "text-gray-300"}><IconIP3 /></div>
    </div>
  );
}
