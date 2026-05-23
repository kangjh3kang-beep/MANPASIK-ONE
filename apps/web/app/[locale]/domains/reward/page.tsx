'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React from 'react';
import { Coins, TrendingUp, Database, Star, Gift, ArrowUpRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { useRewardBalance, useRewardContributions } from '@mmup/api-client';
import { DomainHeader, KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

const REWARDS = [
  { id: 1, title: '건강검진 할인', cost: 500, category: '의료' },
  { id: 2, title: '피트니스 멤버십', cost: 300, category: '운동' },
  { id: 3, title: '건강식품 쿠폰', cost: 200, category: '식품' },
  { id: 4, title: '보험료 할인', cost: 1000, category: '보험' },
];

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function RewardPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const user = session?.user;
  const userId = user?.id || 'user-1';

  const { data: balanceData, isLoading: balanceLoading, error: balanceError, refetch } = useRewardBalance(userId);
  const { data: contribData, isLoading: contribLoading, error: contribError } = useRewardContributions(userId);

  if (balanceLoading || contribLoading) {
    return <LoadingSkeleton columns={3} rows={1} />;
  }

  if (balanceError || contribError) {
    return <ErrorState message="리워드 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

  const balance = balanceData?.data;
  const contributions = contribData?.data || [];

  return (
    <main className="min-h-screen bg-slate-50" aria-label="환자 리워드 풀">
      <DomainHeader
        title="환자 리워드 풀"
        subtitle="개인 건강 데이터 기여에 대한 투명한 보상 및 실제 혜택 제공"
        icon={Coins}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'patient',
          organization: user.organization || 'MMUP 커뮤니티'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="토큰 잔액">
          <div className="rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 p-8 text-white">
            <p className="text-sm font-medium opacity-80">나의 토큰 잔액</p>
            <p className="text-5xl font-bold mt-2 tabular-nums" role="status">
              {balance?.balance?.toLocaleString() || '0'} <span className="text-xl font-medium opacity-70">{balance?.currency || 'MPS'}</span>
            </p>
            <div className="flex items-center gap-2 mt-3">
              <TrendingUp className="h-4 w-4" aria-hidden="true" />
              <span className="text-sm">이번 주 +{balance?.weeklyChange || 0} MPS 적립</span>
            </div>
          </div>
        </section>

        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-3 gap-4">
            <KPICard title="총 기여 횟수" value={`${contributions.length}회`} icon={<Database className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
            <KPICard title="데이터 품질 등급" value="A+" icon={<Star className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="누적 보상" value={`${balance?.totalEarned?.toLocaleString() || '0'} MPS`} icon={<Gift className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
          </div>
        </section>

        <section aria-label="최근 데이터 기여 내역">
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4">최근 데이터 기여 내역</h2>
            <div className="space-y-3">
              {contributions.map((c) => (
                <div key={c.id} className="flex items-center justify-between rounded-xl border border-slate-100 px-4 py-3 hover:bg-slate-50">
                  <div className="flex items-center gap-3">
                    <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-50 text-amber-600" aria-hidden="true">
                      <ArrowUpRight className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-slate-900">{c.type}</p>
                      <p className="text-[11px] text-slate-400">{c.date}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-bold text-amber-600">+{c.points} MPS</p>
                    <p className="text-[10px] text-slate-400">품질 {c.quality}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section aria-label="보상 교환 스토어">
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4">보상 교환 스토어</h2>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              {REWARDS.map(r => (
                <div key={r.id} className="rounded-xl border border-slate-200 p-4 hover:shadow-md transition-shadow cursor-pointer" aria-label={`${r.title} - ${r.cost} MPS`}>
                  <span className="text-[10px] font-bold text-sky-600 bg-sky-50 px-2 py-0.5 rounded-full">{r.category}</span>
                  <h3 className="text-sm font-bold text-slate-900 mt-2">{r.title}</h3>
                  <p className="text-lg font-bold text-amber-600 mt-2">{r.cost} MPS</p>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>

      {/* 관련 도메인 — 모세혈관 교차 연결 */}
      <section aria-label="관련 도메인" className="mt-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[{ name: "임상 콘솔", path: "/domains/clinical", desc: "측정 데이터 기여" },{ name: "건강 앱", path: "/domains/app", desc: "모바일 건강 관리" },{ name: "파트너 연동", path: "/domains/partner", desc: "의료 데이터 공유" }].map(d => (
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
