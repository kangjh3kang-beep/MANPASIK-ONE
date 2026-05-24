'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { BarChart3, ArrowUpRight, ArrowDownRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type ErrorStatus = '해결' | '미해결';

const KPI_CARDS = [
  { label: '총 교환건수', value: '12,450', change: '+18.5%', positive: true },
  { label: '성공률', value: '99.2%', change: '+0.3%p', positive: true },
  { label: '평균 응답시간', value: '89ms', change: '-12ms', positive: true },
  { label: '오류 건수', value: '98', change: '-15.3%', positive: true },
];

interface PartnerVolume {
  name: string;
  volume: number;
  color: string;
}

const PARTNER_VOLUMES: PartnerVolume[] = [
  { name: '서울대병원', volume: 4250, color: 'bg-sky-500' },
  { name: '연세의료원', volume: 3800, color: 'bg-emerald-500' },
  { name: '삼성서울병원', volume: 2900, color: 'bg-violet-500' },
  { name: '아산병원', volume: 1500, color: 'bg-amber-500' },
];

const DATA_TYPES = [
  { type: 'FHIR', percentage: 45, color: 'bg-sky-500' },
  { type: 'HL7', percentage: 30, color: 'bg-emerald-500' },
  { type: 'CSV', percentage: 15, color: 'bg-amber-500' },
  { type: 'API', percentage: 10, color: 'bg-violet-500' },
];

interface ErrorLog {
  id: string;
  timestamp: string;
  partner: string;
  errorType: string;
  status: ErrorStatus;
}

const ERROR_STATUS_STYLES: Record<ErrorStatus, string> = {
  '해결': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '미해결': 'bg-red-50 text-red-700 border-red-200',
};

const MOCK_ERRORS: ErrorLog[] = [
  { id: 'e1', timestamp: '2026-05-24 13:45:22', partner: '아산병원', errorType: '연결 시간 초과', status: '미해결' },
  { id: 'e2', timestamp: '2026-05-24 11:20:05', partner: '삼성서울병원', errorType: 'FHIR 스키마 불일치', status: '해결' },
  { id: 'e3', timestamp: '2026-05-23 18:32:41', partner: '아산병원', errorType: 'TLS 인증서 만료', status: '미해결' },
  { id: 'e4', timestamp: '2026-05-23 09:15:33', partner: '연세의료원', errorType: '데이터 변환 실패', status: '해결' },
  { id: 'e5', timestamp: '2026-05-22 16:48:19', partner: '서울대병원', errorType: '페이로드 크기 초과', status: '해결' },
];

const MONTHLY_TREND = [8200, 8800, 9100, 9500, 10200, 10800, 11200, 11500, 11800, 12000, 12200, 12450];

function TrendSparkline({ data }: { data: number[] }) {
  const max = Math.max(...data);
  const min = Math.min(...data);
  const range = max - min || 1;
  const height = 40;
  const width = 200;
  const step = width / (data.length - 1);

  const points = data.map((v, i) => `${i * step},${height - ((v - min) / range) * height}`).join(' ');

  return (
    <svg width={width} height={height} className="shrink-0" aria-label="월별 추이">
      <polyline fill="none" stroke="#0284c7" strokeWidth="2" points={points} />
    </svg>
  );
}

export default function AnalyticsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const maxVolume = Math.max(...PARTNER_VOLUMES.map(p => p.volume));

  return (
    <main className="min-h-screen bg-slate-50" aria-label="데이터 교환 분석">
      <DomainHeader
        title="데이터 교환 분석"
        subtitle="파트너 기관 간 데이터 교환 현황을 분석합니다"
        icon={BarChart3}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* KPI 카드 */}
        <div className="grid grid-cols-4 gap-4 mb-6">
          {KPI_CARDS.map(kpi => (
            <div key={kpi.label} className="rounded-2xl border border-slate-200 bg-white p-4">
              <p className="text-xs text-slate-500 font-medium mb-1">{kpi.label}</p>
              <p className="text-xl font-bold text-slate-900">{kpi.value}</p>
              <div className={`flex items-center gap-1 mt-1 text-xs font-semibold ${kpi.positive ? 'text-emerald-600' : 'text-red-500'}`}>
                {kpi.positive ? <ArrowUpRight className="h-3 w-3" /> : <ArrowDownRight className="h-3 w-3" />}
                {kpi.change}
              </div>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-2 gap-6 mb-6">
          {/* 파트너별 교환량 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h3 className="text-sm font-bold text-slate-900 mb-4">파트너별 교환량</h3>
            <div className="space-y-3">
              {PARTNER_VOLUMES.map(p => (
                <div key={p.name}>
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs font-medium text-slate-700">{p.name}</span>
                    <span className="text-xs font-semibold text-slate-900">{p.volume.toLocaleString()}건</span>
                  </div>
                  <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
                    <div
                      className={`h-full ${p.color} rounded-full transition-all`}
                      style={{ width: `${(p.volume / maxVolume) * 100}%` }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* 데이터 유형 분포 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h3 className="text-sm font-bold text-slate-900 mb-4">데이터 유형 분포</h3>
            <div className="flex rounded-xl overflow-hidden h-6 mb-4">
              {DATA_TYPES.map(dt => (
                <div
                  key={dt.type}
                  className={`${dt.color} transition-all`}
                  style={{ width: `${dt.percentage}%` }}
                  title={`${dt.type}: ${dt.percentage}%`}
                />
              ))}
            </div>
            <div className="grid grid-cols-2 gap-2">
              {DATA_TYPES.map(dt => (
                <div key={dt.type} className="flex items-center gap-2">
                  <span className={`h-3 w-3 rounded-full ${dt.color}`} />
                  <span className="text-xs text-slate-600">{dt.type} ({dt.percentage}%)</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* 월별 추이 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6 mb-6">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-slate-900">월별 교환량 추이</h3>
            <TrendSparkline data={MONTHLY_TREND} />
          </div>
          <p className="text-xs text-slate-500 mt-2">최근 12개월 | 총 {MONTHLY_TREND.reduce((a, b) => a + b, 0).toLocaleString()}건</p>
        </div>

        {/* 오류 로그 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-5 py-3 bg-slate-50 border-b border-slate-200">
            <h3 className="text-sm font-bold text-slate-900">최근 오류 로그</h3>
          </div>
          <div className="grid grid-cols-[1.2fr_0.8fr_1.2fr_0.6fr] gap-2 px-5 py-2 border-b border-slate-100 text-xs font-semibold text-slate-500">
            <span>발생 시각</span><span>파트너</span><span>오류 유형</span><span className="text-center">상태</span>
          </div>
          {MOCK_ERRORS.map(err => (
            <div key={err.id} className="grid grid-cols-[1.2fr_0.8fr_1.2fr_0.6fr] gap-2 items-center px-5 py-3 border-b border-slate-50 last:border-b-0 hover:bg-slate-50 transition-colors">
              <span className="text-xs font-mono text-slate-600">{err.timestamp}</span>
              <span className="text-sm text-slate-700">{err.partner}</span>
              <span className="text-sm text-slate-900">{err.errorType}</span>
              <span className={`inline-flex justify-center px-2 py-0.5 rounded-full text-[11px] font-semibold border ${ERROR_STATUS_STYLES[err.status]}`}>{err.status}</span>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
