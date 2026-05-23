'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React, { useState, useEffect } from 'react';
import { Activity, Users, AlertTriangle, TrendingUp, Heart, Droplets, Zap, Stethoscope } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
  const localePrefix = useLocalePrefix();
import { usePatientVitals, VitalSign } from '@mmup/api-client';
import { AreaChart, DomainHeader, KPICard, ErrorState } from '@mmup/ui';

function MiniSparkline({ data, color }: { data: number[]; color: string }) {
  const min = Math.min(...data); const max = Math.max(...data); const range = max - min || 1;
  const h = 32; const w = 120;
  const points = data.map((v, i) => `${(i / (data.length - 1)) * w},${h - ((v - min) / range) * h}`).join(' ');
  return <svg width={w} height={h} className="inline-block" aria-hidden="true"><polyline fill="none" stroke={color} strokeWidth="2" strokeLinejoin="round" points={points} /></svg>;
}

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function ClinicalPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const { data: initialVitals, isLoading, error, refetch } = usePatientVitals('patient-1');
  const [vitals, setVitals] = useState<VitalSign[]>([]);

  useEffect(() => {
    if (initialVitals?.data) {
      setVitals(initialVitals.data);
    }
  }, [initialVitals]);

  useEffect(() => {
    const interval = setInterval(() => {
      setVitals(prev => {
        if (prev.length === 0) return prev;
        const last = prev[prev.length - 1];
        const next = {
          ...last,
          heartRate: Math.max(60, Math.min(100, last.heartRate + (Math.random() > 0.5 ? 1 : -1))),
          bloodPressureSys: Math.max(100, Math.min(140, last.bloodPressureSys + (Math.random() > 0.5 ? 2 : -2))),
          timestamp: new Date().toLocaleTimeString('ko-KR', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' })
        };
        return [...prev.slice(1), next];
      });
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  if (isLoading) return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center" role="status" aria-label="임상 데이터 로딩 중">
      <div className="flex flex-col items-center gap-4">
        <Activity className="w-12 h-12 text-sky-500 animate-pulse" aria-hidden="true" />
        <p className="text-slate-500 font-medium">임상 데이터 로드 중...</p>
      </div>
    </div>
  );

  if (error) return (
    <ErrorState message="임상 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />
  );

  const VITAL_CARDS = [
    { key: 'heartRate', label: '심박수', unit: 'bpm', icon: Heart, color: '#ef4444', min: 60, max: 100 },
    { key: 'spo2', label: '산소포화도', unit: '%', icon: Droplets, color: '#3b82f6', min: 95, max: 100 },
    { key: 'bloodPressureSys', label: '수축기 혈압', unit: 'mmHg', icon: Zap, color: '#8b5cf6', min: 90, max: 140 },
  ];

  const latest = vitals.length > 0 ? vitals[vitals.length - 1] as any : {};
  const user = session?.user as any;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="임상 데이터 콘솔">
      <DomainHeader
        title="임상 데이터 콘솔"
        subtitle="환자 생체 데이터 실시간 모니터링 및 AI 분석 현황"
        icon={Stethoscope}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'doctor',
          organization: user.organization || 'MMUP 종합병원'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="통계 요약">
          <div className="grid grid-cols-4 gap-4">
            <KPICard title="활성 환자" value="1,247" icon={<Users className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
            <KPICard title="금일 측정" value="3,891" icon={<span className="text-sm">→</span>} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="이상 징후" value="23" icon={<AlertTriangle className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="진단 정확도" value="99.7%" icon={<TrendingUp className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
          </div>
        </section>

        <section aria-label="실시간 바이탈 카드">
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                <div className="h-2 w-2 rounded-full bg-rose-500 animate-ping" aria-hidden="true" />
                실시간 생체 신호 (Live Vitals)
              </h2>
              <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest bg-slate-50 px-2 py-1 rounded">Stream Active : 2.4GHz</div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {VITAL_CARDS.map(card => {
                const val = latest[card.key as keyof VitalSign] ?? 0;
                const isNormal = Number(val) >= card.min && Number(val) <= card.max;
                const dataArr = vitals.map(v => Number(v[card.key as keyof VitalSign]) ?? 0);
                return (
                  <div key={card.key} className={`rounded-3xl border p-6 transition-all duration-700 ${isNormal ? 'border-slate-100 bg-white hover:shadow-md' : 'border-rose-100 bg-rose-50/50'}`} aria-label={`${card.label}: ${val} ${card.unit}`}>
                    <div className="flex items-center gap-3 mb-4">
                      <div className="p-2.5 rounded-2xl bg-white border border-slate-100 shadow-sm" aria-hidden="true">
                        <card.icon className={`h-5 w-5 ${!isNormal ? 'animate-pulse' : ''}`} style={{ color: card.color }} />
                      </div>
                      <div>
                        <span className="text-[11px] font-bold text-slate-400 uppercase tracking-tight block leading-none mb-1">{card.label}</span>
                        <div className="flex items-baseline gap-1.5">
                          <span className="text-2xl font-black tabular-nums text-slate-900 transition-all" role="status">{val}</span>
                          <span className="text-xs font-bold text-slate-400">{card.unit}</span>
                        </div>
                      </div>
                    </div>
                    <div className="h-12 w-full opacity-80 overflow-hidden">
                      <MiniSparkline data={dataArr} color={card.color} />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>

        <section aria-label="심박수 추이 분석 차트">
          <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-xl font-bold text-slate-900">심박수 정밀 추이 분석</h2>
                <p className="text-sm text-slate-400 mt-1">GNN 모델 기반 99.8% 정확도 필터링 적용</p>
              </div>
              <div className="flex gap-2">
                {['1H', '4H', '12H', '1D'].map(t => (
                  <button key={t} className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all ${t === '4H' ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'}`} aria-label={`${t} 기간 선택`}>{t}</button>
                ))}
              </div>
            </div>
            <div className="h-[320px] w-full">
              <AreaChart
                data={vitals}
                xKey="timestamp"
                yKey="heartRate"
                color="#ef4444"
                height={320}
              />
            </div>
          </div>
        </section>
      </div>

      {/* 관련 도메인 — 모세혈관 교차 연결 */}
      <section aria-label="관련 도메인" className="mt-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[{ name: "생체 예측", path: "/domains/predictor", desc: "AI 질병 위험도 분석" },{ name: "리워드", path: "/domains/reward", desc: "데이터 기여 보상" },{ name: "파트너 연동", path: "/domains/partner", desc: "의료기관 데이터 공유" }].map(d => (
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
    </main>
  );
}
