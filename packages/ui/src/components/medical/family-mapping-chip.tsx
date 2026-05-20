/**
 * @mmup-axis 1 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard WCAG 2.2 AA
 */
import React from "react";
import type { FamilyId } from "@mmup/types";

export function FamilyMappingChip({ family }: { family: FamilyId | FamilyId[] }) {
  const families = Array.isArray(family) ? family : [family];
  return (
    <div className="flex gap-1" aria-label="Family mapping">
      {families.map(f => (
        <span key={f} className="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">
          Family {f}
        </span>
      ))}
    </div>
  );
}
