'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React from 'react';
import { Users2, Building2, ArrowLeftRight, CheckCircle2, Clock, Globe } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePartners, PartnerOrg } from '@mmup/api-client';
import { DomainHeader, KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

const STATUS_STYLE: Record<string, { color: string; label: string }> = {
  active: { color: 'text-emerald-700 bg-emerald-50', label: '연결됨' },
  syncing: { color: 'text-sky-700 bg-sky-50', label: '동기화 중' },
  inactive: { color: 'text-slate-500 bg-slate-50', label: '비활성' },
};

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function PartnerPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const { data: partnersData, isLoading, error, refetch } = usePartners();
  const user = session?.user;

  if (isLoading) {
    return <LoadingSkeleton columns={4} rows={1} />;
  }

  if (error) {
    return <ErrorState message="파트너 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

  const partners: PartnerOrg[] = partnersData?.data || [];
  const activeCount = partners.filter(p => p.status === 'active').length;
  const totalExchanges = partners.reduce((sum, p) => sum + p.totalExchanges, 0);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="파트너 통합 연동">
      <DomainHeader
        title="파트너 통합 연동"
        subtitle="HL7/FHIR 국제 표준 의료 데이터 파이프라인"
        icon={Users2}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'admin',
          organization: user.organization || 'MMUP 관리센터'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-4 gap-4">
            <KPICard title="연동 기관" value={`${partners.length}개소`} icon={<Building2 className="h-4 w-4" />} color="text-rose-600 bg-rose-50" />
            <KPICard title="금일 교환" value={`${totalExchanges.toLocaleString()}건`} icon={<ArrowLeftRight className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
            <KPICard title="성공률" value="99.97%" icon={<CheckCircle2 className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="표준 프로토콜" value="FHIR R4" icon={<Globe className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
          </div>
        </section>

        <section aria-label="연동 기관 현황">
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4">연동 기관 현황</h2>
            <div className="space-y-3">
              {partners.map(p => {
                const st = STATUS_STYLE[p.status] || STATUS_STYLE.inactive;
                return (
                  <div key={p.id} className="flex items-center justify-between rounded-xl border border-slate-100 px-5 py-4 hover:shadow-md transition-shadow" aria-label={`${p.name} - ${st.label}`}>
                    <div className="flex items-center gap-4">
                      <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-100" aria-hidden="true">
                        <Building2 className="h-5 w-5 text-slate-600" />
                      </div>
                      <div>
                        <p className="text-sm font-bold text-slate-900">{p.name}</p>
                        <div className="flex items-center gap-2 mt-0.5">
                          <span className="text-[10px] bg-slate-100 text-slate-500 px-1.5 py-0.5 rounded">{p.protocol}</span>
                          <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${st.color}`}>{st.label}</span>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-bold text-slate-700" role="status">{p.totalExchanges.toLocaleString()}건</p>
                      <p className="text-[11px] text-slate-400 flex items-center gap-1 justify-end">
                        <Clock className="h-3 w-3" aria-hidden="true" />{p.lastSync}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>
      </div>

      {/* 관련 도메인 — 모세혈관 교차 연결 */}
      <section aria-label="관련 도메인" className="mt-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[{ name: "임상 콘솔", path: "/domains/clinical", desc: "환자 데이터 연동" },{ name: "GxP 준수", path: "/domains/gxp", desc: "규제 감사 추적" },{ name: "개발자 포털", path: "/domains/dev-portal", desc: "API 연동 가이드" }].map(d => (
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
