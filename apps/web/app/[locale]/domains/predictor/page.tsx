'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

import React from 'react';
import { Brain, AlertTriangle, ActivitySquare, Target } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
  const localePrefix = useLocalePrefix();
import { usePredictions, DiseasePrediction } from '@mmup/api-client';
import { RiskProgress, DomainHeader, KPICard, ErrorState } from '@mmup/ui';

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function PredictorPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const { data: predictions, isLoading, error, refetch } = usePredictions();
  const user = session?.user as any;

  if (isLoading) return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center" role="status" aria-label="예측 모델 분석 중">
      <div className="flex flex-col items-center gap-4">
        <Brain className="w-12 h-12 text-emerald-500 animate-bounce" aria-hidden="true" />
        <p className="text-slate-500 font-medium">GNN 예측 모델 분석 중...</p>
      </div>
    </div>
  );

  if (error) return (
    <ErrorState message="예측 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />
  );

  const data = predictions?.data || [];
  const highRiskCount = data.filter(p => p.riskScore >= 0.75).length;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="생체 지표 예측 엔진">
      <DomainHeader
        title="생체 지표 예측 엔진"
        subtitle="GNN 기반 다중 질환 시계열 예측 및 데이터 기반 예방 의학"
        icon={ActivitySquare}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'researcher',
          organization: user.organization || 'MMUP AI Lab'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <KPICard title="분석 대상" value={`${data.length}개 질환`} icon={<Target className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="고위험 감지" value={`${highRiskCount}건`} icon={<AlertTriangle className="h-4 w-4" />} color={highRiskCount > 0 ? 'text-red-600 bg-red-50' : 'text-emerald-600 bg-emerald-50'} />
            <KPICard title="모델 정확도" value="96.3%" icon={<Brain className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
          </div>
        </section>

        <section aria-label="질환별 위험도 분석">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {data.map((pred: DiseasePrediction) => (
              <div key={pred.disease} className="rounded-2xl border border-slate-200 bg-white p-2 shadow-sm" aria-label={`${pred.disease} 위험도 ${(pred.riskScore * 100).toFixed(0)}%`}>
                <RiskProgress
                  diseaseName={pred.disease}
                  riskScore={pred.riskScore}
                  factors={pred.contributingFactors}
                  trend={pred.trend as any}
                />
              </div>
            ))}
          </div>
        </section>
      </div>

      {/* 관련 도메인 — 모세혈관 교차 연결 */}
      <section aria-label="관련 도메인" className="mt-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[{ name: "임상 콘솔", path: "/domains/clinical", desc: "실시간 바이탈 데이터" },{ name: "AI 에이전트", path: "/domains/agents-hub", desc: "ML 파이프라인" },{ name: "리워드", path: "/domains/reward", desc: "예측 데이터 기여" }].map(d => (
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
