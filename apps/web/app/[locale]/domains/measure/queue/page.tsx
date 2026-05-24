'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { RefreshCw, WifiOff, Clock, CheckCircle, AlertCircle, Loader2, Upload } from 'lucide-react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type SyncStatus = '대기중' | '동기화중' | '완료' | '실패';

interface QueueItem {
  id: string;
  measurementType: string;
  timestamp: string;
  size: string;
  status: SyncStatus;
  retryCount: number;
  errorMessage?: string;
}

const STATUS_CONFIG: Record<SyncStatus, { color: string; bgColor: string; borderColor: string; icon: React.ElementType }> = {
  '대기중': { color: 'text-amber-700', bgColor: 'bg-amber-50', borderColor: 'border-amber-200', icon: Clock },
  '동기화중': { color: 'text-sky-700', bgColor: 'bg-sky-50', borderColor: 'border-sky-200', icon: Loader2 },
  '완료': { color: 'text-emerald-700', bgColor: 'bg-emerald-50', borderColor: 'border-emerald-200', icon: CheckCircle },
  '실패': { color: 'text-red-700', bgColor: 'bg-red-50', borderColor: 'border-red-200', icon: AlertCircle },
};

const MOCK_QUEUE: QueueItem[] = [
  {
    id: 'q1',
    measurementType: '혈당 측정',
    timestamp: '2026-05-24 09:32:15',
    size: '2.4 KB',
    status: '대기중',
    retryCount: 0,
  },
  {
    id: 'q2',
    measurementType: '콜레스테롤 측정',
    timestamp: '2026-05-24 09:28:41',
    size: '3.1 KB',
    status: '동기화중',
    retryCount: 0,
  },
  {
    id: 'q3',
    measurementType: '당화혈색소 측정',
    timestamp: '2026-05-24 08:15:03',
    size: '2.8 KB',
    status: '대기중',
    retryCount: 0,
  },
  {
    id: 'q4',
    measurementType: '혈압 측정',
    timestamp: '2026-05-23 22:10:55',
    size: '1.9 KB',
    status: '실패',
    retryCount: 3,
    errorMessage: '서버 응답 시간 초과 (504)',
  },
];

const STATS = {
  pending: 3,
  completed: 47,
  failed: 1,
};

export default function SyncQueuePage() {
  const localePrefix = useLocalePrefix();
  const [isOffline] = useState(true);
  const [queue, setQueue] = useState(MOCK_QUEUE);
  const lastSyncTime = '2026-05-24 08:00:12';

  const handleRetry = (itemId: string) => {
    setQueue((prev) =>
      prev.map((item) =>
        item.id === itemId ? { ...item, status: '대기중' as SyncStatus, retryCount: item.retryCount + 1 } : item,
      ),
    );
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="동기화 큐">
      <DomainHeader
        title="동기화 큐"
        subtitle="오프라인 측정 데이터의 동기화 상태를 관리합니다"
        icon={Upload}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 오프라인 배너 */}
        {isOffline && (
          <div className="flex items-center gap-3 rounded-2xl border border-amber-200 bg-amber-50 px-5 py-4">
            <WifiOff className="h-5 w-5 text-amber-600 flex-shrink-0" />
            <div>
              <p className="text-sm font-bold text-amber-800">현재 오프라인 상태입니다</p>
              <p className="text-xs text-amber-600 mt-0.5">
                연결 복구 시 자동 동기화됩니다. 측정 데이터는 안전하게 로컬에 저장되어 있습니다.
              </p>
            </div>
          </div>
        )}

        {/* 통계 카드 */}
        <div className="grid grid-cols-3 gap-4">
          {[
            { label: '대기 중', value: STATS.pending, color: 'text-amber-600', bg: 'bg-amber-50', icon: Clock },
            { label: '동기화 완료', value: STATS.completed, color: 'text-emerald-600', bg: 'bg-emerald-50', icon: CheckCircle },
            { label: '실패', value: STATS.failed, color: 'text-red-600', bg: 'bg-red-50', icon: AlertCircle },
          ].map((stat) => (
            <div key={stat.label} className="rounded-2xl border border-slate-200 bg-white p-5 text-center">
              <div className={`h-10 w-10 rounded-xl ${stat.bg} flex items-center justify-center mx-auto mb-2`}>
                <stat.icon className={`h-5 w-5 ${stat.color}`} />
              </div>
              <p className="text-2xl font-bold text-slate-900">{stat.value}</p>
              <p className="text-xs text-slate-500 mt-1">{stat.label}</p>
            </div>
          ))}
        </div>

        {/* 마지막 동기화 시간 + 전체 동기화 버튼 */}
        <div className="flex items-center justify-between">
          <p className="text-xs text-slate-400">
            마지막 동기화: <span className="font-semibold text-slate-600">{lastSyncTime}</span>
          </p>
          <button
            className="px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors inline-flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={isOffline}
            title={isOffline ? '오프라인 상태에서는 동기화할 수 없습니다' : '대기 중인 항목을 모두 동기화합니다'}
          >
            <RefreshCw className="h-3.5 w-3.5" />
            전체 동기화
          </button>
        </div>

        {/* 큐 목록 */}
        <div className="space-y-3">
          {queue.map((item) => {
            const config = STATUS_CONFIG[item.status];
            const StatusIcon = config.icon;

            return (
              <div
                key={item.id}
                className="rounded-2xl border border-slate-200 bg-white p-5 transition-shadow hover:shadow-sm"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex items-start gap-3 flex-1 min-w-0">
                    <div className={`h-10 w-10 rounded-xl ${config.bgColor} flex items-center justify-center flex-shrink-0`}>
                      <StatusIcon className={`h-5 w-5 ${config.color} ${item.status === '동기화중' ? 'animate-spin' : ''}`} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-sm font-bold text-slate-900">{item.measurementType}</h3>
                        <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold border ${config.bgColor} ${config.color} ${config.borderColor}`}>
                          {item.status}
                        </span>
                      </div>
                      <p className="text-xs text-slate-500">
                        {item.timestamp} | {item.size}
                        {item.retryCount > 0 && (
                          <span className="ml-2 text-amber-600">재시도 {item.retryCount}회</span>
                        )}
                      </p>
                      {item.errorMessage && (
                        <p className="text-xs text-red-500 mt-1">{item.errorMessage}</p>
                      )}
                    </div>
                  </div>

                  {item.status === '실패' && (
                    <button
                      onClick={() => handleRetry(item.id)}
                      className="px-3 py-1.5 rounded-xl border border-slate-200 text-xs font-semibold text-slate-600 hover:bg-slate-50 transition-colors inline-flex items-center gap-1"
                    >
                      <RefreshCw className="h-3 w-3" />
                      재시도
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        {/* 동기화 정보 안내 */}
        <div className="rounded-2xl border border-slate-100 bg-slate-50 p-4">
          <p className="text-xs text-slate-500 leading-relaxed">
            측정 데이터는 오프라인 상태에서도 기기에 안전하게 저장됩니다. 네트워크 연결이 복구되면
            CRDT 기반 충돌 해결 알고리즘을 통해 자동으로 동기화됩니다.
            동기화 실패 항목은 최대 5회까지 자동 재시도되며, 이후 수동 재시도가 필요합니다.
          </p>
        </div>
      </div>
    </main>
  );
}
