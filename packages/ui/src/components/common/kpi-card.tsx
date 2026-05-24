'use client';

import React from 'react';

export interface KPICardProps {
  title: string;
  value: string | number;
  unit?: string;
  trend?: string;
  icon?: React.ReactNode;
  color?: string;
  loading?: boolean;
  /** Optional status indicator shown as a colored top border stripe */
  status?: 'normal' | 'caution' | 'danger';
}

const statusBorderStyles: Record<string, string> = {
  normal: '2px solid rgb(16 185 129)', // emerald-500
  caution: '2px solid rgb(245 158 11)', // amber-500
  danger: '2px solid rgb(239 68 68)', // red-500
};

export function KPICard({ title, value, unit, trend, icon, color = 'text-sky-600 bg-sky-50', loading = false, status }: KPICardProps) {
  const borderTopStyle = status ? statusBorderStyles[status] : undefined;
  if (loading) {
    return (
      <div
        className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm animate-pulse"
        style={borderTopStyle ? { borderTop: borderTopStyle } : undefined}
        aria-label={`${title} 로딩 중`}
        role="status"
      >
        <div className="flex items-center justify-between mb-3">
          <div className="h-3 w-20 bg-slate-200 rounded" />
          <div className="h-8 w-8 bg-slate-100 rounded-xl" />
        </div>
        <div className="h-7 w-24 bg-slate-200 rounded mt-1" />
        {trend !== undefined && <div className="h-3 w-16 bg-slate-100 rounded mt-2" />}
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
      style={borderTopStyle ? { borderTop: borderTopStyle } : undefined}
      aria-label={`${title}: ${value}${unit ? ` ${unit}` : ''}`}
    >
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs font-semibold text-slate-500">{title}</span>
        {icon && (
          <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${color}`}>
            {icon}
          </div>
        )}
      </div>
      <p className="text-2xl font-bold text-slate-900 tabular-nums" role="status">
        {value}
        {unit && <span className="text-sm font-normal text-slate-400 ml-1">{unit}</span>}
      </p>
      {trend && (
        <p className="text-xs text-slate-400 mt-1">{trend}</p>
      )}
    </div>
  );
}
