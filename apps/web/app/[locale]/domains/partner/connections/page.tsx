'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { Wifi, CheckCircle, AlertTriangle, XCircle, RefreshCw } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type ConnectionStatus = '정상' | '지연' | '오류';

interface PartnerConnection {
  id: string;
  name: string;
  status: ConnectionStatus;
  uptime: string;
  avgLatency: string;
  lastCheck: string;
  sparkline: number[];
}

const STATUS_CONFIG: Record<ConnectionStatus, { dot: string; bg: string; text: string; icon: React.ElementType }> = {
  '정상': { dot: 'bg-emerald-500', bg: 'bg-emerald-50 border-emerald-200', text: 'text-emerald-700', icon: CheckCircle },
  '지연': { dot: 'bg-amber-500', bg: 'bg-amber-50 border-amber-200', text: 'text-amber-700', icon: AlertTriangle },
  '오류': { dot: 'bg-red-500', bg: 'bg-red-50 border-red-200', text: 'text-red-700', icon: XCircle },
};

const MOCK_CONNECTIONS: PartnerConnection[] = [
  {
    id: 'p1', name: '서울대병원', status: '정상', uptime: '99.9%', avgLatency: '45ms',
    lastCheck: '2026-05-24 14:30:00',
    sparkline: [42, 38, 45, 41, 39, 44, 46, 43, 40, 47, 44, 42, 38, 45, 43, 41, 39, 46, 44, 42, 45, 43, 41, 44],
  },
  {
    id: 'p2', name: '연세의료원', status: '정상', uptime: '99.2%', avgLatency: '120ms',
    lastCheck: '2026-05-24 14:28:00',
    sparkline: [110, 115, 125, 118, 122, 130, 128, 115, 112, 120, 118, 125, 130, 122, 115, 118, 120, 125, 128, 115, 118, 122, 120, 118],
  },
  {
    id: 'p3', name: '삼성서울병원', status: '지연', uptime: '98.5%', avgLatency: '350ms',
    lastCheck: '2026-05-24 14:25:00',
    sparkline: [200, 250, 300, 280, 350, 400, 380, 320, 350, 370, 340, 310, 350, 380, 360, 340, 350, 370, 380, 350, 340, 360, 350, 345],
  },
  {
    id: 'p4', name: '아산병원', status: '오류', uptime: '95.1%', avgLatency: '-',
    lastCheck: '2026-05-24 13:45:00',
    sparkline: [80, 85, 90, 120, 200, 500, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
  },
];

function MiniSparkline({ data, status }: { data: number[]; status: ConnectionStatus }) {
  const max = Math.max(...data, 1);
  const height = 32;
  const width = 120;
  const step = width / (data.length - 1);
  const color = status === '정상' ? '#10b981' : status === '지연' ? '#f59e0b' : '#ef4444';

  const points = data.map((v, i) => `${i * step},${height - (v / max) * height}`).join(' ');

  return (
    <svg width={width} height={height} className="shrink-0" aria-label="지연시간 추이">
      <polyline fill="none" stroke={color} strokeWidth="1.5" points={points} />
    </svg>
  );
}

export default function ConnectionsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const healthyCount = MOCK_CONNECTIONS.filter(c => c.status === '정상').length;
  const totalCount = MOCK_CONNECTIONS.length;
  const overallHealthy = healthyCount === totalCount;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="연결 상태 모니터링">
      <DomainHeader
        title="연결 상태 모니터링"
        subtitle="파트너 기관 API 연결 상태를 실시간으로 확인합니다"
        icon={Wifi}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* 전체 상태 배너 */}
        <div className={`rounded-2xl border p-4 mb-6 flex items-center justify-between ${
          overallHealthy
            ? 'border-emerald-200 bg-emerald-50'
            : 'border-amber-200 bg-amber-50'
        }`}>
          <div className="flex items-center gap-3">
            {overallHealthy ? (
              <CheckCircle className="h-5 w-5 text-emerald-600" />
            ) : (
              <AlertTriangle className="h-5 w-5 text-amber-600" />
            )}
            <div>
              <p className={`text-sm font-bold ${overallHealthy ? 'text-emerald-800' : 'text-amber-800'}`}>
                {overallHealthy ? '모든 연결 정상' : '일부 연결에 문제가 감지되었습니다'}
              </p>
              <p className={`text-xs ${overallHealthy ? 'text-emerald-600' : 'text-amber-600'}`}>
                정상: {healthyCount}/{totalCount} | 마지막 확인: 2026-05-24 14:30
              </p>
            </div>
          </div>
          <button className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white border border-slate-200 text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors">
            <RefreshCw className="h-3.5 w-3.5" />
            전체 점검
          </button>
        </div>

        {/* 파트너 연결 카드 */}
        <div className="grid grid-cols-2 gap-4">
          {MOCK_CONNECTIONS.map(conn => {
            const config = STATUS_CONFIG[conn.status];
            const StatusIcon = config.icon;
            return (
              <div key={conn.id} className="rounded-2xl border border-slate-200 bg-white p-5 transition-shadow hover:shadow-sm">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-2">
                    <span className={`h-2.5 w-2.5 rounded-full ${config.dot} animate-pulse`} />
                    <h3 className="text-sm font-bold text-slate-900">{conn.name}</h3>
                  </div>
                  <span className={`flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold border ${config.bg} ${config.text}`}>
                    <StatusIcon className="h-3 w-3" />
                    {conn.status}
                  </span>
                </div>

                <div className="grid grid-cols-2 gap-3 mb-4">
                  <div className="rounded-xl bg-slate-50 p-3">
                    <p className="text-xs text-slate-500">Uptime SLA</p>
                    <p className="text-lg font-bold text-slate-900">{conn.uptime}</p>
                  </div>
                  <div className="rounded-xl bg-slate-50 p-3">
                    <p className="text-xs text-slate-500">평균 지연시간</p>
                    <p className="text-lg font-bold text-slate-900">{conn.avgLatency}</p>
                  </div>
                </div>

                <div className="flex items-center justify-between mb-3">
                  <span className="text-xs text-slate-400">최근 24시간 지연시간</span>
                  <MiniSparkline data={conn.sparkline} status={conn.status} />
                </div>

                <div className="flex items-center justify-between">
                  <span className="text-[11px] text-slate-400">마지막 점검: {conn.lastCheck}</span>
                  <button className="px-3 py-1.5 rounded-lg text-xs font-semibold bg-sky-600 text-white hover:bg-sky-700 transition-colors">
                    연결 테스트
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
