'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { Award, Check, Lock, TrendingUp, Database } from 'lucide-react';

type TierName = 'Bronze' | 'Silver' | 'Gold' | 'Platinum';

interface Tier {
  name: TierName;
  nameKo: string;
  threshold: number;
  color: string;
  bgColor: string;
  borderColor: string;
  iconBg: string;
  benefits: string[];
  achieved: boolean;
  current: boolean;
}

interface BenefitComparison {
  label: string;
  values: Record<TierName, string>;
}

const TIERS: Tier[] = [
  {
    name: 'Bronze',
    nameKo: '브론즈',
    threshold: 0,
    color: 'text-orange-700',
    bgColor: 'bg-orange-50',
    borderColor: 'border-orange-200',
    iconBg: 'bg-orange-100',
    benefits: ['기본 적립률 1%', '월간 건강 리포트', '커뮤니티 접근'],
    achieved: true,
    current: false,
  },
  {
    name: 'Silver',
    nameKo: '실버',
    threshold: 1000,
    color: 'text-slate-700',
    bgColor: 'bg-slate-50',
    borderColor: 'border-slate-300',
    iconBg: 'bg-slate-200',
    benefits: ['적립률 1.5%', '카트리지 5% 할인', '주간 건강 리포트', '우선 고객 지원'],
    achieved: true,
    current: true,
  },
  {
    name: 'Gold',
    nameKo: '골드',
    threshold: 5000,
    color: 'text-amber-700',
    bgColor: 'bg-amber-50',
    borderColor: 'border-amber-300',
    iconBg: 'bg-amber-100',
    benefits: ['적립률 2%', '카트리지 10% 할인', '우선 예약 진료', '무료 배송', '프리미엄 AI 코칭'],
    achieved: false,
    current: false,
  },
  {
    name: 'Platinum',
    nameKo: '플래티넘',
    threshold: 15000,
    color: 'text-violet-700',
    bgColor: 'bg-violet-50',
    borderColor: 'border-violet-300',
    iconBg: 'bg-violet-100',
    benefits: ['적립률 3%', '카트리지 15% 할인', '전용 의료 상담', '무료 배송', 'VIP 챌린지 초대', '연구 참여 보상 2배'],
    achieved: false,
    current: false,
  },
];

const BENEFIT_COMPARISON: BenefitComparison[] = [
  { label: '적립률', values: { Bronze: '1%', Silver: '1.5%', Gold: '2%', Platinum: '3%' } },
  { label: '할인율', values: { Bronze: '-', Silver: '5%', Gold: '10%', Platinum: '15%' } },
  { label: '우선 예약', values: { Bronze: '-', Silver: '-', Gold: 'O', Platinum: 'O' } },
  { label: '무료 배송', values: { Bronze: '-', Silver: '-', Gold: 'O', Platinum: 'O' } },
  { label: 'AI 코칭', values: { Bronze: '기본', Silver: '기본', Gold: '프리미엄', Platinum: '프리미엄' } },
  { label: '전용 상담', values: { Bronze: '-', Silver: '-', Gold: '-', Platinum: 'O' } },
];

const CURRENT_POINTS = 2847;
const NEXT_TIER_THRESHOLD = 5000;
const PROGRESS_TO_NEXT = Math.round((CURRENT_POINTS / NEXT_TIER_THRESHOLD) * 100);

const MONTHLY_CONTRIBUTIONS = [
  { label: '건강 측정', count: 23, points: 460 },
  { label: '데이터 기여', count: 12, points: 1800 },
  { label: '챌린지 달성', count: 3, points: 450 },
  { label: '출석 보상', count: 18, points: 137 },
];

export default function TiersPage() {
  return (
    <main className="min-h-screen bg-slate-50" aria-label="등급 혜택">
      <DomainHeader title="등급 혜택" subtitle="등급별 혜택을 확인하고 다음 등급에 도전하세요" icon={Award} />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-8">
        {/* 현재 등급 카드 */}
        <div className="rounded-2xl border-2 border-slate-300 bg-gradient-to-r from-slate-50 to-white p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="h-12 w-12 rounded-2xl bg-slate-200 flex items-center justify-center">
                <Award className="h-6 w-6 text-slate-600" />
              </div>
              <div>
                <p className="text-xs font-semibold text-slate-500">현재 등급</p>
                <p className="text-xl font-black text-slate-900">Silver (실버)</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-xs text-slate-500">보유 포인트</p>
              <p className="text-2xl font-bold text-slate-900 tabular-nums">{CURRENT_POINTS.toLocaleString()} MPS</p>
            </div>
          </div>

          <div className="mb-2">
            <div className="flex items-center justify-between mb-1.5">
              <span className="text-xs text-slate-500">Gold 등급까지</span>
              <span className="text-xs font-bold text-slate-700">{PROGRESS_TO_NEXT}% ({(NEXT_TIER_THRESHOLD - CURRENT_POINTS).toLocaleString()} MPS 남음)</span>
            </div>
            <div className="h-3 bg-slate-200 rounded-full overflow-hidden">
              <div
                className="h-full rounded-full bg-gradient-to-r from-slate-400 to-amber-400 transition-all duration-700"
                style={{ width: `${PROGRESS_TO_NEXT}%` }}
              />
            </div>
          </div>
          <div className="flex justify-between text-[10px] text-slate-400">
            <span>Silver (1,000)</span>
            <span>Gold (5,000)</span>
          </div>
        </div>

        {/* 등급 카드 그리드 */}
        <section aria-label="등급 목록">
          <h2 className="text-lg font-bold text-slate-900 mb-4">전체 등급</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {TIERS.map((tier) => (
              <div
                key={tier.name}
                className={`rounded-2xl border-2 p-5 transition-all ${
                  tier.current
                    ? `${tier.borderColor} ${tier.bgColor} ring-2 ring-offset-2 ring-slate-300 shadow-md`
                    : tier.achieved
                      ? `${tier.borderColor} ${tier.bgColor}`
                      : 'border-slate-200 bg-white opacity-75'
                }`}
              >
                <div className="flex items-center justify-between mb-3">
                  <div className={`h-10 w-10 rounded-xl ${tier.iconBg} flex items-center justify-center`} aria-hidden="true">
                    {tier.achieved ? <Award className={`h-5 w-5 ${tier.color}`} /> : <Lock className="h-5 w-5 text-slate-400" />}
                  </div>
                  {tier.current && (
                    <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-sky-100 text-sky-700 border border-sky-200">
                      현재 등급
                    </span>
                  )}
                  {tier.achieved && !tier.current && (
                    <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-700 border border-emerald-200">
                      달성 완료
                    </span>
                  )}
                </div>

                <h3 className={`text-base font-black ${tier.achieved ? tier.color : 'text-slate-400'}`}>
                  {tier.nameKo}
                </h3>
                <p className="text-[11px] text-slate-500 mt-0.5 mb-3">
                  {tier.threshold === 0 ? '시작 등급' : `${tier.threshold.toLocaleString()} MPS 이상`}
                </p>

                <ul className="space-y-1.5">
                  {tier.benefits.map((benefit, idx) => (
                    <li key={idx} className="flex items-center gap-1.5">
                      <Check className={`h-3 w-3 shrink-0 ${tier.achieved ? 'text-emerald-500' : 'text-slate-300'}`} />
                      <span className={`text-xs ${tier.achieved ? 'text-slate-700' : 'text-slate-400'}`}>{benefit}</span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </section>

        {/* 혜택 비교표 */}
        <section aria-label="등급별 혜택 비교">
          <h2 className="text-lg font-bold text-slate-900 mb-4">등급별 혜택 비교</h2>
          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full min-w-[500px]">
                <thead>
                  <tr className="border-b border-slate-100">
                    <th className="px-5 py-3 text-left text-xs font-bold text-slate-500">혜택</th>
                    {TIERS.map((tier) => (
                      <th key={tier.name} className={`px-4 py-3 text-center text-xs font-bold ${tier.current ? 'bg-slate-50' : ''} ${tier.color}`}>
                        {tier.nameKo}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {BENEFIT_COMPARISON.map((row, idx) => (
                    <tr key={row.label} className={`border-b border-slate-50 ${idx % 2 === 0 ? '' : 'bg-slate-50/50'}`}>
                      <td className="px-5 py-3 text-sm font-medium text-slate-700">{row.label}</td>
                      {TIERS.map((tier) => (
                        <td key={tier.name} className={`px-4 py-3 text-center text-sm ${tier.current ? 'bg-slate-50/50 font-bold text-slate-900' : 'text-slate-600'}`}>
                          {row.values[tier.name] === 'O' ? (
                            <Check className="h-4 w-4 text-emerald-500 mx-auto" />
                          ) : row.values[tier.name] === '-' ? (
                            <span className="text-slate-300">-</span>
                          ) : (
                            row.values[tier.name]
                          )}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* 월간 데이터 기여 요약 */}
        <section aria-label="월간 데이터 기여 요약">
          <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Database className="h-5 w-5 text-sky-600" />
            이번 달 데이터 기여 요약
          </h2>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            {MONTHLY_CONTRIBUTIONS.map((item) => (
              <div key={item.label} className="rounded-xl border border-slate-200 bg-white p-4 text-center">
                <p className="text-xs text-slate-500 font-medium">{item.label}</p>
                <p className="text-2xl font-bold text-slate-900 mt-1 tabular-nums">{item.count}회</p>
                <p className="text-xs font-semibold text-emerald-600 mt-0.5 flex items-center justify-center gap-1">
                  <TrendingUp className="h-3 w-3" />
                  +{item.points} MPS
                </p>
              </div>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
