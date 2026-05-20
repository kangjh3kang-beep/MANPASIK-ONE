/**
 * @mmup-axis 8 이행 검증
 * @mmup-stage 5 이행
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-2
 * @standard WCAG 2.2 AA
 */
import React from "react";

export function NextActionCard({ title, description, onClick }: { title: string, description: string, onClick?: () => void }) {
  return (
    <button 
      onClick={onClick}
      className="block w-full text-left p-4 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
      aria-label={`Next action: ${title}`}
    >
      <h3 className="font-bold text-gray-900">{title}</h3>
      <p className="text-gray-600 text-sm mt-1">{description}</p>
    </button>
  );
}
