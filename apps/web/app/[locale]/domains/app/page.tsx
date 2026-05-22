'use client';
import React from 'react';
import { Activity, Stethoscope, Heart, Droplets, Footprints, Moon, Smartphone, ChevronRight, Users2, Database } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAppMetrics, AppMetric } from '@mmup/api-client';
import { KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const ICON_MAP: Record<string, React.ElementType> = {
  Heart,
  Droplets,
  Footprints,
  Moon,
};

const RECENT = [
  { time: '14:30', type: '혈압', value: '122/78 mmHg' },
  { time: '12:00', type: '혈당', value: '95 mg/dL' },
  { time: '09:15', type: '체온', value: '36.5 °C' },
  { time: '08:00', type: '체중', value: '68.2 kg' },
];

export default function MobileAppPage() {
  const localePrefix = useLocalePrefix();
  const { data: metricsData, isLoading, error, refetch } = useAppMetrics();

  if (isLoading) {
    return <LoadingSkeleton columns={2} rows={1} />;
  }

  if (error) {
    return <ErrorState message="건강 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

  const metrics: AppMetric[] = metricsData?.data || [];

  const ACTIONS = [
    { label: '새 측정 시작', href: `${localePrefix}/domains/clinical`, icon: Stethoscope },
    { label: '의사에게 공유', href: `${localePrefix}/domains/partner`, icon: Users2 },
    { label: '리워드 확인', href: `${localePrefix}/domains/reward`, icon: Database },
  ];

  return (
    <main className="min-h-screen bg-gradient-to-b from-sky-50 to-white" aria-label="만파식 건강 앱">
      <header className="border-b border-slate-200 bg-white/80 backdrop-blur px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-500 text-white" aria-hidden="true">
            <Smartphone className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">만파식 건강 앱 v2.0</h1>
        </div>
        <p className="text-sm text-slate-500">나의 건강 데이터를 한눈에 -- PWA 모바일 뷰</p>
      </header>

      <div className="max-w-lg mx-auto px-6 py-8 space-y-6">
        <section aria-label="건강 점수 요약">
          <div className="rounded-2xl bg-gradient-to-r from-cyan-500 to-sky-500 p-6 text-white shadow-lg shadow-cyan-500/20">
            <p className="text-sm opacity-80">건강 데이터 요약</p>
            <p className="text-xl font-bold mt-1">오늘도 건강한 하루 보내세요!</p>
            <div className="mt-3 flex items-center gap-2 text-sm opacity-90">
              <Activity className="h-4 w-4" aria-hidden="true" />
              <span role="status">건강 점수: <strong>87/100</strong> (안정)</span>
            </div>
          </div>
        </section>

        <section aria-label="건강 지표 카드">
          <div className="grid grid-cols-2 gap-3">
            {metrics.map(card => {
              const IconComp = ICON_MAP[card.iconName];
              return (
                <div key={card.label}>
                  <KPICard
                    title={card.label}
                    value={card.value}
                    unit={card.unit}
                    trend={card.trend}
                    icon={IconComp ? <IconComp className="h-4 w-4" /> : undefined}
                    color={card.color}
                  />
                </div>
              );
            })}
          </div>
        </section>

        <section aria-label="최근 측정 기록">
          <div className="rounded-2xl border border-slate-200 bg-white p-5">
            <h2 className="text-sm font-bold text-slate-900 mb-3">최근 측정 기록</h2>
            <div className="space-y-2">
              {RECENT.map((m, i) => (
                <div key={i} className="flex items-center justify-between py-2 border-b border-slate-50 last:border-0">
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-slate-400 font-mono w-12">{m.time}</span>
                    <span className="text-sm font-medium text-slate-700">{m.type}</span>
                  </div>
                  <span className="text-sm font-bold text-slate-900" role="status">{m.value}</span>
                </div>
              ))}
            </div>
          </div>
        </section>

        <nav aria-label="빠른 작업">
          <div className="space-y-3">
            {ACTIONS.map((action) => (
              <Link
                key={action.label}
                href={action.href}
                className="flex w-full items-center justify-between rounded-2xl border border-slate-200 bg-white px-5 py-4 text-sm font-bold text-slate-700 hover:bg-slate-50 hover:border-sky-500 hover:text-sky-600 transition-all group shadow-sm"
                aria-label={action.label}
              >
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-lg bg-slate-50 flex items-center justify-center group-hover:bg-sky-50 transition-colors" aria-hidden="true">
                    <action.icon className="h-4 w-4" />
                  </div>
                  {action.label}
                </div>
                <ChevronRight className="h-4 w-4 text-slate-300 group-hover:translate-x-1 transition-transform" aria-hidden="true" />
              </Link>
            ))}
          </div>
        </nav>
      </div>
    </main>
  );
}
