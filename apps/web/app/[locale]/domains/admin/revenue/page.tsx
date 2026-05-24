'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { TrendingUp, ArrowUpRight, ArrowDownRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const KPI_CARDS = [
  { label: '월 매출', value: '₩12,450,000', change: '+12.3%', positive: true },
  { label: 'MRR', value: '₩8,200,000', change: '+8.7%', positive: true },
  { label: '구독 전환율', value: '34%', change: '+2.1%p', positive: true },
  { label: '평균 구매액', value: '₩45,000', change: '-1.5%', positive: false },
];

const MONTHLY_REVENUE = [
  { month: '6월', value: 6200 }, { month: '7월', value: 6800 }, { month: '8월', value: 7100 },
  { month: '9월', value: 7500 }, { month: '10월', value: 8200 }, { month: '11월', value: 8800 },
  { month: '12월', value: 9400 }, { month: '1월', value: 9800 }, { month: '2월', value: 10200 },
  { month: '3월', value: 10900 }, { month: '4월', value: 11600 }, { month: '5월', value: 12450 },
];

const REVENUE_BREAKDOWN = [
  { category: '카트리지', percentage: 60, color: 'bg-sky-500' },
  { category: '구독', percentage: 25, color: 'bg-emerald-500' },
  { category: '기기', percentage: 10, color: 'bg-amber-500' },
  { category: '기타', percentage: 5, color: 'bg-slate-400' },
];

const TOP_PRODUCTS = [
  { name: '혈당 카트리지 (10개입)', sales: 2340, revenue: '₩67,860,000', growth: '+15.2%', positive: true },
  { name: '종합검사 카트리지 세트', sales: 1120, revenue: '₩99,680,000', growth: '+22.8%', positive: true },
  { name: 'MPS 리더기 Pro', sales: 340, revenue: '₩98,600,000', growth: '+5.1%', positive: true },
  { name: '프리미엄 구독 (월)', sales: 890, revenue: '₩35,600,000', growth: '+18.4%', positive: true },
  { name: '콜레스테롤 카트리지 (10개입)', sales: 780, revenue: '₩27,300,000', growth: '-3.2%', positive: false },
];

const SUBSCRIPTION_METRICS = [
  { month: '1월', newSub: 120, churn: 45, net: 75 },
  { month: '2월', newSub: 135, churn: 38, net: 97 },
  { month: '3월', newSub: 148, churn: 42, net: 106 },
  { month: '4월', newSub: 162, churn: 35, net: 127 },
  { month: '5월', newSub: 178, churn: 40, net: 138 },
];

export default function RevenuePage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const maxRevenue = Math.max(...MONTHLY_REVENUE.map(m => m.value));

  return (
    <main className="min-h-screen bg-slate-50" aria-label="매출 분석">
      <DomainHeader
        title="매출 분석"
        subtitle="플랫폼 매출 현황과 성장 추이를 분석합니다"
        icon={TrendingUp}
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

        {/* 월별 매출 바 차트 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6 mb-6">
          <h3 className="text-sm font-bold text-slate-900 mb-4">월별 매출 추이 (단위: 만원)</h3>
          <div className="flex items-end gap-2 h-48">
            {MONTHLY_REVENUE.map(m => (
              <div key={m.month} className="flex-1 flex flex-col items-center gap-1">
                <span className="text-[10px] font-semibold text-slate-600">{(m.value / 10).toFixed(0)}</span>
                <div
                  className="w-full bg-sky-500 rounded-t-md transition-all hover:bg-sky-600"
                  style={{ height: `${(m.value / maxRevenue) * 160}px` }}
                />
                <span className="text-[10px] text-slate-500">{m.month}</span>
              </div>
            ))}
          </div>
        </div>

        {/* 매출 구성 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6 mb-6">
          <h3 className="text-sm font-bold text-slate-900 mb-4">매출 구성</h3>
          <div className="flex rounded-xl overflow-hidden h-8 mb-4">
            {REVENUE_BREAKDOWN.map(item => (
              <div
                key={item.category}
                className={`${item.color} transition-all`}
                style={{ width: `${item.percentage}%` }}
                title={`${item.category}: ${item.percentage}%`}
              />
            ))}
          </div>
          <div className="flex gap-6">
            {REVENUE_BREAKDOWN.map(item => (
              <div key={item.category} className="flex items-center gap-2">
                <span className={`h-3 w-3 rounded-full ${item.color}`} />
                <span className="text-xs text-slate-600">{item.category} ({item.percentage}%)</span>
              </div>
            ))}
          </div>
        </div>

        {/* 인기 상품 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden mb-6">
          <div className="px-5 py-3 bg-slate-50 border-b border-slate-200">
            <h3 className="text-sm font-bold text-slate-900">인기 상품 TOP 5</h3>
          </div>
          <div className="grid grid-cols-[2fr_0.8fr_1.2fr_0.8fr] gap-2 px-5 py-2 border-b border-slate-100 text-xs font-semibold text-slate-500">
            <span>상품명</span><span className="text-right">판매량</span><span className="text-right">매출</span><span className="text-right">성장률</span>
          </div>
          {TOP_PRODUCTS.map(p => (
            <div key={p.name} className="grid grid-cols-[2fr_0.8fr_1.2fr_0.8fr] gap-2 items-center px-5 py-3 border-b border-slate-50 last:border-b-0 hover:bg-slate-50 transition-colors">
              <span className="text-sm font-medium text-slate-900 truncate">{p.name}</span>
              <span className="text-sm text-slate-600 text-right">{p.sales.toLocaleString()}</span>
              <span className="text-sm font-semibold text-slate-900 text-right">{p.revenue}</span>
              <span className={`text-xs font-semibold text-right ${p.positive ? 'text-emerald-600' : 'text-red-500'}`}>{p.growth}</span>
            </div>
          ))}
        </div>

        {/* 구독 지표 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h3 className="text-sm font-bold text-slate-900 mb-4">구독 지표 (월별)</h3>
          <div className="grid grid-cols-[0.8fr_1fr_1fr_1fr] gap-2 px-2 py-2 text-xs font-semibold text-slate-500 border-b border-slate-100">
            <span>월</span><span className="text-right">신규</span><span className="text-right">해지</span><span className="text-right">순증</span>
          </div>
          {SUBSCRIPTION_METRICS.map(m => (
            <div key={m.month} className="grid grid-cols-[0.8fr_1fr_1fr_1fr] gap-2 items-center px-2 py-2.5 border-b border-slate-50 last:border-b-0">
              <span className="text-sm text-slate-700 font-medium">{m.month}</span>
              <span className="text-sm text-emerald-600 font-semibold text-right">+{m.newSub}</span>
              <span className="text-sm text-red-500 font-semibold text-right">-{m.churn}</span>
              <span className="text-sm text-sky-600 font-bold text-right">+{m.net}</span>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
