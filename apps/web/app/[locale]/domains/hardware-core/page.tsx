'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

import React from 'react';
import { Cpu, HardDrive, Unplug, Zap, Activity, ShieldCheck, Microscope } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
  const localePrefix = useLocalePrefix();
import { useHardwareDevices, useHardwareDiagnostics, HardwareDevice } from '@mmup/api-client';
import { DomainHeader, KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function HardwareCorePage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const { data: devicesData, isLoading: devicesLoading, error: devicesError, refetch } = useHardwareDevices();
  const { data: diagData, isLoading: diagLoading, error: diagError } = useHardwareDiagnostics();
  const user = session?.user as any;

  if (devicesLoading || diagLoading) {
    return <LoadingSkeleton columns={4} rows={1} />;
  }

  if (devicesError || diagError) {
    return <ErrorState message="하드웨어 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

  const devices: HardwareDevice[] = devicesData?.data || [];
  const diag = diagData?.data;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="하드웨어 코어 엔진">
      <DomainHeader
        title="하드웨어 코어 엔진"
        subtitle="Rust 기반 저지연 생체 데이터 파싱 및 엣지 AI 추론 모듈 제어"
        icon={Microscope}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'dev',
          organization: user.organization || 'MMUP Hardware Lab'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <KPICard title="엔진 상태" value={`Active (${diag?.engineVersion || 'N/A'})`} icon={<span className="text-sm">→</span>} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="CPU 온도" value={`${diag?.cpuTemp?.toFixed(1) || '0'}°C`} icon={<Zap className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="메모리 점유" value={`${diag?.memoryMB?.toFixed(2) || '0'}MB`} icon={<HardDrive className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
            <KPICard title="보안 검증" value={diag?.securityStatus || 'N/A'} icon={<ShieldCheck className="h-4 w-4" />} color="text-indigo-600 bg-indigo-50" />
          </div>
        </section>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <section aria-label="물리 장치 연결 상태">
            <div className="rounded-3xl border border-slate-200 bg-white p-8">
              <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
                <Unplug className="h-5 w-5 text-sky-600" aria-hidden="true" />
                물리 장치 연결 상태 (HW Bridge)
              </h2>
              <div className="space-y-4">
                {devices.map(dev => (
                  <div key={dev.name} className="flex items-center justify-between p-4 rounded-2xl bg-slate-50 border border-slate-100" aria-label={`${dev.name} - ${dev.status}`}>
                    <span className="text-sm font-bold text-slate-700">{dev.name}</span>
                    <div className="flex items-center gap-4">
                      <span className="text-xs font-medium text-slate-400">Latency: {dev.latency}</span>
                      <span className="px-2 py-1 bg-emerald-100 text-emerald-700 text-[10px] font-bold rounded-lg" role="status">{dev.status}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>

          <section aria-label="Rust Core 스레드 모니터">
            <div className="rounded-3xl border border-slate-200 bg-white p-8">
              <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
                <Cpu className="h-5 w-5 text-indigo-600" aria-hidden="true" />
                Rust Core 스레드 모니터
              </h2>
              <div className="relative h-48 w-full bg-slate-900 rounded-2xl overflow-hidden p-4 font-mono text-[10px] text-emerald-400" role="log" aria-label="시스템 로그">
                <p className="mb-1">{`[sys] Rust ManPaSik Engine initialized...`}</p>
                <p className="mb-1">{`[io] BLE data incoming: 0x4A 0x22 0xFF 0x01`}</p>
                <p className="mb-1 text-sky-400">{`[ml] Edge inference complete: diag_result=0.87`}</p>
                <p className="mb-1">{`[sys] Memory swap check... OK`}</p>
                <p className="mb-1 animate-pulse">{`_`}</p>
              </div>
            </div>
          </section>
        </div>
      </div>

      {/* 관련 도메인 — 모세혈관 교차 연결 */}
      <section aria-label="관련 도메인" className="mt-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[{ name: "개발자 포털", path: "/domains/dev-portal", desc: "SDK 도구" },{ name: "AI 에이전트", path: "/domains/agents-hub", desc: "AI 추론 엔진" },{ name: "임상 콘솔", path: "/domains/clinical", desc: "센서 데이터 전달" }].map(d => (
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
