'use client';
export const runtime = 'edge';

import React, { useState, useEffect } from 'react';
import { Users2, Building2, ArrowLeftRight, CheckCircle2, Clock, Globe } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { DomainHeader } from '@mmup/ui';

const PARTNERS = [
  { name: '서울대학교병원', type: 'FHIR R4', status: 'active', exchanges: 12847, lastSync: '2분 전' },
  { name: '연세 세브란스', type: 'HL7 v2.5', status: 'active', exchanges: 8923, lastSync: '5분 전' },
  { name: '삼성서울병원', type: 'FHIR R4', status: 'active', exchanges: 7456, lastSync: '8분 전' },
  { name: '아산병원', type: 'FHIR R4', status: 'syncing', exchanges: 3421, lastSync: '동기화 중...' },
  { name: '국립암센터', type: 'HL7 v2.5', status: 'inactive', exchanges: 1205, lastSync: '3시간 전' },
];

const STATUS_STYLE: Record<string, { color: string; label: string }> = {
  active: { color: 'text-emerald-700 bg-emerald-50', label: '연결됨' },
  syncing: { color: 'text-sky-700 bg-sky-50', label: '동기화 중' },
  inactive: { color: 'text-slate-500 bg-slate-50', label: '비활성' },
};

export default function PartnerPage() {
  const { data: session } = useSession();
  const user = session?.user as any;

  return (
    <div className="min-h-screen bg-slate-50">
      <DomainHeader 
        title="파트너 통합 연동" 
        subtitle="HL7/FHIR 국제 표준 의료 데이터 파이프라인 및 글로벌 네트워크 현황"
        icon={Users2}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'partner',
          organization: user.organization || '글로벌 의료 협력센터'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {[
            { label: '연동 기관', val: '5개소', icon: Building2, c: 'text-rose-600 bg-rose-50' },
            { label: '금일 교환', val: '33,852건', icon: ArrowLeftRight, c: 'text-sky-600 bg-sky-50' },
            { label: '성공률', val: '99.97%', icon: CheckCircle2, c: 'text-emerald-600 bg-emerald-50' },
            { label: '표준 프로토콜', val: 'FHIR R4', icon: Globe, c: 'text-violet-600 bg-violet-50' },
          ].map(s => (
            <div key={s.label} className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-tight">{s.label}</span>
                <div className={`p-2 rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-black text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>

        <div className="rounded-3xl border border-slate-200 bg-white overflow-hidden shadow-xl">
          <div className="px-8 py-6 border-b border-slate-100 flex items-center justify-between bg-white">
            <h2 className="text-lg font-bold text-slate-900">실시간 연동 기관 현황</h2>
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-xs font-bold text-slate-500 uppercase tracking-widest">HL7 Bridge Active</span>
            </div>
          </div>
          <div className="divide-y divide-slate-50">
            {PARTNERS.map(p => {
              const st = STATUS_STYLE[p.status];
              return (
                <div key={p.name} className="flex items-center justify-between px-8 py-5 hover:bg-slate-50/50 transition-colors">
                  <div className="flex items-center gap-4">
                    <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-100/50 border border-slate-100">
                      <Building2 className="h-6 w-6 text-slate-400" />
                    </div>
                    <div>
                      <p className="text-base font-bold text-slate-900">{p.name}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <span className="text-[10px] bg-slate-100 text-slate-500 px-2 py-0.5 rounded font-bold uppercase tracking-wider">{p.type}</span>
                        <span className={`text-[10px] font-bold px-2 py-0.5 rounded-lg ${st.color}`}>{st.label}</span>
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-lg font-black text-slate-900 tabular-nums">{p.exchanges.toLocaleString()}</p>
                    <p className="text-[10px] font-bold text-slate-400 flex items-center gap-1 justify-end uppercase tracking-tighter">
                      <Clock className="h-3 w-3" />
                      Last Sync: {p.lastSync}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
