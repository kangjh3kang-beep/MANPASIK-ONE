'use client';

import React from 'react';

interface LoadingSkeletonProps {
  columns?: number;
  rows?: number;
}

export function LoadingSkeleton({ columns = 4, rows = 1 }: LoadingSkeletonProps) {
  return (
    <div className="min-h-screen bg-slate-50" role="status" aria-label="페이지 로딩 중">
      <div className="border-b border-slate-200 bg-white px-8 py-6 animate-pulse">
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-10 bg-slate-200 rounded-xl" />
          <div className="h-6 w-48 bg-slate-200 rounded" />
        </div>
        <div className="h-4 w-72 bg-slate-100 rounded mt-2" />
      </div>
      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        {Array.from({ length: rows }).map((_, rowIdx) => (
          <div key={rowIdx} className={`grid gap-4`} style={{ gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))` }}>
            {Array.from({ length: columns }).map((_, colIdx) => (
              <div key={colIdx} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm animate-pulse">
                <div className="flex items-center justify-between mb-3">
                  <div className="h-3 w-20 bg-slate-200 rounded" />
                  <div className="h-8 w-8 bg-slate-100 rounded-xl" />
                </div>
                <div className="h-7 w-24 bg-slate-200 rounded" />
              </div>
            ))}
          </div>
        ))}
        <div className="rounded-3xl border border-slate-200 bg-white p-8 animate-pulse">
          <div className="h-5 w-48 bg-slate-200 rounded mb-6" />
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="h-16 bg-slate-100 rounded-xl" />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
