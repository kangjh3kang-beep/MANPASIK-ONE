'use client';
import React, { useState, useEffect } from 'react';
import { Users2, Building2, ArrowLeftRight, CheckCircle2, Clock, Globe } from 'lucide-react';
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
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;
  return (
    <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2"><div className="flex h-10 w-10 items-center justify-center rounded-xl bg-rose-600 text-white"><Users2 className="h-5 w-5" /></div><h1 className="text-2xl font-bold text-slate-900">파트너 통합 연동</h1></div>
        <p className="text-sm text-slate-500">HL7/FHIR 국제 표준 의료 데이터 파이프라인</p>
      </div>
      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: '연동 기관', val: '5개소', icon: Building2, c: 'text-rose-600 bg-rose-50' },
            { label: '금일 교환', val: '33,852건', icon: ArrowLeftRight, c: 'text-sky-600 bg-sky-50' },
            { label: '성공률', val: '99.97%', icon: CheckCircle2, c: 'text-emerald-600 bg-emerald-50' },
            { label: '표준 프로토콜', val: 'FHIR R4', icon: Globe, c: 'text-violet-600 bg-violet-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3"><span className="text-xs font-semibold text-slate-500">{s.label}</span><div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div></div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">연동 기관 현황</h2>
          <div className="space-y-3">
            {PARTNERS.map(p => {
              const st = STATUS_STYLE[p.status];
              return (
                <div key={p.name} className="flex items-center justify-between rounded-xl border border-slate-100 px-5 py-4 hover:shadow-md transition-shadow">
                  <div className="flex items-center gap-4"><div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-100"><Building2 className="h-5 w-5 text-slate-600" /></div><div><p className="text-sm font-bold text-slate-900">{p.name}</p><div className="flex items-center gap-2 mt-0.5"><span className="text-[10px] bg-slate-100 text-slate-500 px-1.5 py-0.5 rounded">{p.type}</span><span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${st.color}`}>{st.label}</span></div></div></div>
                  <div className="text-right"><p className="text-sm font-bold text-slate-700">{p.exchanges.toLocaleString()}건</p><p className="text-[11px] text-slate-400 flex items-center gap-1 justify-end"><Clock className="h-3 w-3" />{p.lastSync}</p></div>
                </div>);
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
