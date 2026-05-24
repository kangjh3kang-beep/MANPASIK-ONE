'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Settings, Users, Package, BarChart3, Shield, AlertTriangle, TrendingUp, Activity, MapPin, Search } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const KPI_DATA = [
  { label: '총 사용자', value: '12,847', change: '+3.2%', icon: Users, color: 'text-sky-600 bg-sky-50' },
  { label: '오늘 측정', value: '3,421', change: '+12%', icon: Activity, color: 'text-emerald-600 bg-emerald-50' },
  { label: '활성 리더기', value: '8,923', change: '+1.5%', icon: Package, color: 'text-purple-600 bg-purple-50' },
  { label: '이상 알림', value: '7', change: '-2건', icon: AlertTriangle, color: 'text-amber-600 bg-amber-50' },
];

const RECENT_EVENTS = [
  { time: '14:32', type: '측정', desc: '서울 강남구 — 혈당 측정 이상치 감지', severity: 'warning' },
  { time: '14:28', type: '리더기', desc: 'MPK-A1B2C3 펌웨어 업데이트 완료', severity: 'info' },
  { time: '14:15', type: '사용자', desc: '신규 가입 +23명 (일일 목표 120%)', severity: 'success' },
  { time: '13:50', type: '재고', desc: '혈당 카트리지 재고 부족 경고 (창고 B)', severity: 'warning' },
  { time: '13:30', type: '결제', desc: '정기배송 결제 처리 완료 342건', severity: 'info' },
  { time: '12:45', type: '품질', desc: 'Lot-2026-A 보정 검증 통과', severity: 'success' },
];

const REGIONS = [
  { name: '서울', users: 4521, devices: 3200, measurements: 1230 },
  { name: '경기', users: 3210, devices: 2100, measurements: 890 },
  { name: '부산', users: 1890, devices: 1300, measurements: 560 },
  { name: '대구', users: 1120, devices: 780, measurements: 340 },
  { name: '인천', users: 980, devices: 650, measurements: 290 },
];

export default function AdminPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [nlQuery, setNlQuery] = useState('');
  const [nlResult, setNlResult] = useState<string | null>(null);
  const [drillLevel, setDrillLevel] = useState<'global' | 'region'>('global');

  const handleNlQuery = () => {
    if (!nlQuery.trim()) return;
    setNlResult(`"${nlQuery}"에 대한 분석 결과: 최근 7일간 해당 지표는 전주 대비 +5.3% 상승했으며, 서울 강남구에서 가장 높은 활동량을 보입니다.`);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="관리자 대시보드">
      <DomainHeader
        title="관리자 대시보드"
        subtitle="플랫폼 전체를 한눈에 모니터링하고 관리합니다"
        icon={Settings}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'admin' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-6 py-8 space-y-6">
        {/* AI 자연어 콘솔 */}
        <div className="rounded-2xl bg-slate-900 p-6 text-white">
          <h2 className="text-sm font-bold text-slate-400 mb-3">AI 관리 콘솔</h2>
          <div className="flex gap-3">
            <input
              type="text"
              value={nlQuery}
              onChange={e => setNlQuery(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && handleNlQuery()}
              placeholder="질문하세요: '이번 주 지역별 이상 측정 현황', '재고 부족 지역'"
              aria-label="AI 관리 콘솔 질문 입력"
              className="flex-1 px-4 py-3 rounded-xl bg-slate-800 border border-slate-700 text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-sky-500/30"
            />
            <button onClick={handleNlQuery} aria-label="질문 검색" className="px-6 py-3 rounded-xl bg-sky-600 text-sm font-bold hover:bg-sky-700 transition-colors">
              <Search className="h-4 w-4" />
            </button>
          </div>
          {nlResult && (
            <div className="mt-3 p-4 rounded-xl bg-slate-800 text-sm text-slate-300">
              {nlResult}
            </div>
          )}
        </div>

        {/* KPI 카드 */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {KPI_DATA.map(kpi => (
            <div key={kpi.label} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-semibold text-slate-500">{kpi.label}</span>
                <div className={`h-8 w-8 rounded-lg flex items-center justify-center ${kpi.color}`}>
                  <kpi.icon className="h-4 w-4" />
                </div>
              </div>
              <p className="text-2xl font-bold text-slate-900">{kpi.value}</p>
              <p className={`text-xs font-medium mt-1 ${kpi.change.startsWith('+') ? 'text-emerald-600' : kpi.change.startsWith('-') ? 'text-sky-600' : 'text-slate-400'}`}>
                {kpi.change} 전일 대비
              </p>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 지역별 현황 */}
          <div className="lg:col-span-2 rounded-2xl border border-slate-200 bg-white p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base font-bold text-slate-900 flex items-center gap-2">
                <MapPin className="h-4 w-4 text-sky-600" />
                지역별 현황
              </h2>
              <div className="flex gap-1">
                <button onClick={() => setDrillLevel('global')}
                  className={`px-3 py-1 rounded-lg text-xs font-semibold ${drillLevel === 'global' ? 'bg-sky-600 text-white' : 'text-slate-500'}`}>
                  전체
                </button>
                <button onClick={() => setDrillLevel('region')}
                  className={`px-3 py-1 rounded-lg text-xs font-semibold ${drillLevel === 'region' ? 'bg-sky-600 text-white' : 'text-slate-500'}`}>
                  상세
                </button>
              </div>
            </div>
            <div className="space-y-2">
              {REGIONS.map(region => (
                <div key={region.name} className="flex items-center gap-4 p-3 rounded-xl hover:bg-slate-50 transition-colors">
                  <span className="text-sm font-bold text-slate-700 w-12">{region.name}</span>
                  <div className="flex-1">
                    <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
                      <div className="h-full bg-sky-500 rounded-full" style={{ width: `${(region.users / 5000) * 100}%` }} />
                    </div>
                  </div>
                  <div className="flex gap-4 text-xs text-slate-500">
                    <span>{region.users.toLocaleString()}명</span>
                    <span>{region.devices.toLocaleString()}대</span>
                    <span>{region.measurements.toLocaleString()}건</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* 실시간 이벤트 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
              <Activity className="h-4 w-4 text-emerald-600" />
              실시간 이벤트
            </h2>
            <div className="space-y-3">
              {RECENT_EVENTS.map((event, i) => (
                <div key={i} className="flex items-start gap-3">
                  <div className={`h-2 w-2 rounded-full mt-1.5 flex-shrink-0 ${
                    event.severity === 'warning' ? 'bg-amber-500' :
                    event.severity === 'success' ? 'bg-emerald-500' : 'bg-sky-500'
                  }`} />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-[10px] font-mono text-slate-400">{event.time}</span>
                      <span className="text-[10px] px-1.5 py-0.5 rounded bg-slate-100 font-semibold text-slate-500">{event.type}</span>
                    </div>
                    <p className="text-xs text-slate-600 mt-0.5">{event.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-4">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: '품질 관리', path: '/domains/gxp', desc: '의약품 규정 준수' },
              { name: '건강 데이터', path: '/domains/clinical', desc: '측정 데이터 관리' },
              { name: '병원 연동', path: '/domains/partner', desc: '협력 기관 관리' },
            ].map(d => (
              <Link key={d.path} href={`${localePrefix}${d.path}`}
                className="flex items-center gap-3 p-3 rounded-xl border border-slate-200 bg-white hover:border-sky-300 hover:shadow-sm transition-all group">
                <div className="h-8 w-8 rounded-lg bg-sky-50 flex items-center justify-center text-sky-600 group-hover:bg-sky-100">
                  <span className="text-sm">→</span>
                </div>
                <div>
                  <span className="text-sm font-semibold text-slate-700 group-hover:text-sky-600">{d.name}</span>
                  <p className="text-[11px] text-slate-400">{d.desc}</p>
                </div>
              </Link>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
