'use client';

import React from 'react';
import { AlertTriangle } from 'lucide-react';

interface ErrorStateProps {
  message?: string;
  onRetry?: () => void;
}

export function ErrorState({ message = '데이터를 불러오는 중 오류가 발생했습니다.', onRetry }: ErrorStateProps) {
  return (
    <div
      className="min-h-screen bg-slate-50 flex items-center justify-center"
      aria-live="polite"
      role="alert"
    >
      <div className="bg-white p-8 rounded-3xl border border-slate-200 shadow-xl max-w-md text-center">
        <AlertTriangle className="w-16 h-16 text-rose-500 mx-auto mb-4" aria-hidden="true" />
        <h2 className="text-xl font-bold text-slate-900 mb-2">데이터 통신 에러</h2>
        <p className="text-sm text-slate-500 mb-6">{message}</p>
        {onRetry && (
          <button
            onClick={onRetry}
            className="px-6 py-2 bg-sky-600 text-white rounded-xl font-bold hover:bg-sky-700 transition-colors"
            aria-label="데이터 다시 불러오기"
          >
            재시도
          </button>
        )}
      </div>
    </div>
  );
}
