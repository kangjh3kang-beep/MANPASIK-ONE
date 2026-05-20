'use client';
import React, { useState, useEffect } from 'react';
import { Coins, TrendingUp, Database, Star, Gift, ArrowUpRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { DomainHeader } from '@mmup/ui';

const CONTRIBUTIONS = [
  { date: '2026-04-26', type: '심박 데이터', points: 150, quality: 'A+' },
  { date: '2026-04-25', type: '혈압 데이터', points: 120, quality: 'A' },
  { date: '2026-04-25', type: '수면 패턴', points: 200, quality: 'A+' },
  { date: '2026-04-24', type: '활동량 데이터', points: 80, quality: 'B+' },
  { date: '2026-04-23', type: '식이 기록', points: 100, quality: 'A' },
];

const REWARDS = [
  { id: 1, title: '건강검진 할인', cost: 500, category: '의료' },
  { id: 2, title: '피트니스 멤버십', cost: 300, category: '운동' },
  { id: 3, title: '건강식품 쿠폰', cost: 200, category: '식품' },
  { id: 4, title: '보험료 할인', cost: 1000, category: '보험' },
];

export default function RewardPage() {
  const { data: session } = useSession();
  const user = session?.user as any;

  return (
    <div className="min-h-screen bg-slate-50">
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
        <div className="rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 p-8 text-white">
          <p className="text-sm font-medium opacity-80">나의 토큰 잔액</p>
          <p className="text-5xl font-bold mt-2 tabular-nums">2,847 <span className="text-xl font-medium opacity-70">MPS</span></p>
          <div className="flex items-center gap-2 mt-3"><TrendingUp className="h-4 w-4" /><span className="text-sm">이번 주 +650 MPS 적립</span></div>
        </div>
        <div className="grid grid-cols-3 gap-4">
          {[
            { label: '총 기여 횟수', val: '1,247회', icon: Database, c: 'text-sky-600 bg-sky-50' },
            { label: '데이터 품질 등급', val: 'A+', icon: Star, c: 'text-amber-600 bg-amber-50' },
            { label: '누적 보상', val: '12,480 MPS', icon: Gift, c: 'text-emerald-600 bg-emerald-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3"><span className="text-xs font-semibold text-slate-500">{s.label}</span><div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div></div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">최근 데이터 기여 내역</h2>
          <div className="space-y-3">
            {CONTRIBUTIONS.map((c, i) => (
              <div key={i} className="flex items-center justify-between rounded-xl border border-slate-100 px-4 py-3 hover:bg-slate-50">
                <div className="flex items-center gap-3"><div className="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-50 text-amber-600"><ArrowUpRight className="h-4 w-4" /></div><div><p className="text-sm font-semibold text-slate-900">{c.type}</p><p className="text-[11px] text-slate-400">{c.date}</p></div></div>
                <div className="text-right"><p className="text-sm font-bold text-amber-600">+{c.points} MPS</p><p className="text-[10px] text-slate-400">품질 {c.quality}</p></div>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">보상 교환 스토어</h2>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            {REWARDS.map(r => (
              <div key={r.id} className="rounded-xl border border-slate-200 p-4 hover:shadow-md transition-shadow cursor-pointer">
                <span className="text-[10px] font-bold text-sky-600 bg-sky-50 px-2 py-0.5 rounded-full">{r.category}</span>
                <h3 className="text-sm font-bold text-slate-900 mt-2">{r.title}</h3>
                <p className="text-lg font-bold text-amber-600 mt-2">{r.cost} MPS</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
