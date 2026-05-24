'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Receipt, Database, Trophy, ShoppingBag, Stethoscope, CalendarCheck, ChevronLeft, ChevronRight, TrendingUp } from 'lucide-react';

type TransactionType = '데이터 기여' | '챌린지 달성' | '스토어 구매' | '진료 결제' | '출석 보상';
type FilterType = '전체' | '적립' | '사용' | '교환';

interface Transaction {
  id: string;
  date: string;
  type: TransactionType;
  description: string;
  amount: number;
  balanceAfter: number;
}

const TYPE_CONFIG: Record<TransactionType, { icon: React.ElementType; color: string }> = {
  '데이터 기여': { icon: Database, color: 'bg-sky-50 text-sky-600' },
  '챌린지 달성': { icon: Trophy, color: 'bg-amber-50 text-amber-600' },
  '스토어 구매': { icon: ShoppingBag, color: 'bg-violet-50 text-violet-600' },
  '진료 결제': { icon: Stethoscope, color: 'bg-rose-50 text-rose-600' },
  '출석 보상': { icon: CalendarCheck, color: 'bg-emerald-50 text-emerald-600' },
};

const TRANSACTIONS: Transaction[] = [
  { id: 't-1', date: '2026-05-24', type: '출석 보상', description: '일일 출석 체크', amount: 50, balanceAfter: 2847 },
  { id: 't-2', date: '2026-05-23', type: '데이터 기여', description: '혈당 검사 데이터 제공', amount: 150, balanceAfter: 2797 },
  { id: 't-3', date: '2026-05-22', type: '챌린지 달성', description: '7일 연속 측정 달성', amount: 300, balanceAfter: 2647 },
  { id: 't-4', date: '2026-05-21', type: '스토어 구매', description: '콜레스테롤 카트리지 구매', amount: -350, balanceAfter: 2347 },
  { id: 't-5', date: '2026-05-20', type: '데이터 기여', description: '콜레스테롤 검사 데이터 제공', amount: 150, balanceAfter: 2697 },
  { id: 't-6', date: '2026-05-19', type: '출석 보상', description: '일일 출석 체크', amount: 50, balanceAfter: 2547 },
  { id: 't-7', date: '2026-05-18', type: '진료 결제', description: '화상진료 결제 (이피부 전문의)', amount: -500, balanceAfter: 2497 },
  { id: 't-8', date: '2026-05-17', type: '데이터 기여', description: 'HbA1c 검사 데이터 제공', amount: 200, balanceAfter: 2997 },
  { id: 't-9', date: '2026-05-16', type: '챌린지 달성', description: '걷기 챌린지 주간 목표 달성', amount: 100, balanceAfter: 2797 },
  { id: 't-10', date: '2026-05-15', type: '스토어 구매', description: '건강식품 쿠폰 교환', amount: -200, balanceAfter: 2697 },
];

const FILTERS: FilterType[] = ['전체', '적립', '사용', '교환'];

function filterTransactions(txns: Transaction[], filter: FilterType): Transaction[] {
  switch (filter) {
    case '적립': return txns.filter(t => t.amount > 0 && t.type !== '스토어 구매');
    case '사용': return txns.filter(t => t.amount < 0 && t.type !== '스토어 구매');
    case '교환': return txns.filter(t => t.type === '스토어 구매');
    default: return txns;
  }
}

export default function TransactionsPage() {
  const [filter, setFilter] = useState<FilterType>('전체');
  const [page, setPage] = useState(1);
  const totalPages = 3;

  const filtered = filterTransactions(TRANSACTIONS, filter);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="포인트 거래 내역">
      <DomainHeader title="포인트 거래 내역" subtitle="MPS 포인트 적립 및 사용 내역을 확인하세요" icon={Receipt} />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 잔액 카드 */}
        <div className="rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 p-6 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium opacity-80">현재 잔액</p>
              <p className="text-4xl font-bold mt-1 tabular-nums" role="status">
                2,847 <span className="text-lg font-medium opacity-70">MPS</span>
              </p>
            </div>
            <div className="flex items-center gap-2 bg-white/20 rounded-xl px-3 py-2">
              <TrendingUp className="h-4 w-4" />
              <span className="text-sm font-semibold">이번 주 +650</span>
            </div>
          </div>
        </div>

        {/* 필터 탭 */}
        <div className="flex gap-2" role="tablist" aria-label="거래 유형 필터">
          {FILTERS.map((f) => (
            <button
              key={f}
              role="tab"
              aria-selected={filter === f}
              onClick={() => { setFilter(f); setPage(1); }}
              className={`px-4 py-2 rounded-xl text-sm font-semibold transition-all ${
                filter === f
                  ? 'bg-slate-900 text-white'
                  : 'bg-white text-slate-500 border border-slate-200 hover:border-slate-300'
              }`}
            >
              {f}
            </button>
          ))}
        </div>

        {/* 거래 목록 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="divide-y divide-slate-100">
            {filtered.map((txn) => {
              const config = TYPE_CONFIG[txn.type];
              const Icon = config.icon;
              const isPositive = txn.amount > 0;

              return (
                <div key={txn.id} className="flex items-center gap-4 px-5 py-4 hover:bg-slate-50 transition-colors">
                  <div className={`h-10 w-10 rounded-xl flex items-center justify-center shrink-0 ${config.color}`} aria-hidden="true">
                    <Icon className="h-5 w-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-slate-900 truncate">{txn.description}</p>
                    <div className="flex items-center gap-2 mt-0.5">
                      <span className="text-[11px] text-slate-400">{txn.date}</span>
                      <span className="text-[10px] px-1.5 py-0.5 rounded bg-slate-100 text-slate-500 font-medium">{txn.type}</span>
                    </div>
                  </div>
                  <div className="text-right shrink-0">
                    <p className={`text-sm font-bold tabular-nums ${isPositive ? 'text-emerald-600' : 'text-red-600'}`}>
                      {isPositive ? '+' : ''}{txn.amount.toLocaleString()} MPS
                    </p>
                    <p className="text-[10px] text-slate-400 tabular-nums">잔액 {txn.balanceAfter.toLocaleString()}</p>
                  </div>
                </div>
              );
            })}
          </div>

          {filtered.length === 0 && (
            <div className="py-12 text-center text-sm text-slate-400">
              해당 유형의 거래 내역이 없습니다.
            </div>
          )}
        </div>

        {/* 페이지네이션 */}
        <div className="flex items-center justify-center gap-2">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="h-9 w-9 rounded-lg border border-slate-200 bg-white flex items-center justify-center text-slate-400 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed"
            aria-label="이전 페이지"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
            <button
              key={p}
              onClick={() => setPage(p)}
              className={`h-9 w-9 rounded-lg text-sm font-semibold transition-all ${
                page === p
                  ? 'bg-slate-900 text-white'
                  : 'border border-slate-200 bg-white text-slate-500 hover:bg-slate-50'
              }`}
              aria-label={`${p} 페이지`}
              aria-current={page === p ? 'page' : undefined}
            >
              {p}
            </button>
          ))}
          <button
            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="h-9 w-9 rounded-lg border border-slate-200 bg-white flex items-center justify-center text-slate-400 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed"
            aria-label="다음 페이지"
          >
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      </div>
    </main>
  );
}
