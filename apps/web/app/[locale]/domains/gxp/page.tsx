'use client';
import React from 'react';
import { Pill, ClipboardCheck, Shield, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { useSession } from 'next-auth/react';
import { useAuditLogs, useCompliance, AuditLog } from '@mmup/api-client';
import { DataTable, Column } from '@mmup/ui';

export default function GxPPage() {
  const { data: session } = useSession();
  const { data: logs, isLoading: logsLoading } = useAuditLogs();
  const { data: compliances, isLoading: compLoading } = useCompliance();
  const user = session?.user as any;


  if (logsLoading || compLoading) return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center animate-pulse">
      <ClipboardCheck className="w-12 h-12 text-violet-500" />
    </div>
  );

  const columns: Column<AuditLog>[] = [
    { key: 'timestamp', label: '시간' },
    { key: 'user', label: '사용자', render: (val) => <span className="font-bold text-slate-900">{val}</span> },
    { key: 'action', label: '수행 기록' },
    { key: 'batchId', label: '배치 ID' },
    { 
      key: 'hasSignature', 
      label: '전자서명',
      render: (val) => (
        <span className={`px-2 py-1 rounded-lg text-[10px] font-black uppercase ${
          val ? 'bg-emerald-50 text-emerald-600' : 'bg-rose-50 text-rose-600'
        }`}>
          {val ? 'Signed' : 'Missing'}
        </span>
      )
    },
  ];

  const auditData = logs?.data || [];
  const compData = compliances?.data || [];

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-600 text-white"><Pill className="h-5 w-5" /></div>
          <h1 className="text-2xl font-bold text-slate-900">의약품 GxP 준수 시스템</h1>
        </div>
        <p className="text-sm text-slate-500">cGMP 감사 추적, 전자 서명 및 실시간 규제 준수 모니터링</p>
      </div>
      
      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {[
            { label: '준수율', val: '99.4%', icon: Shield, c: 'text-emerald-600 bg-emerald-50' },
            { label: '감사항목', val: `${compData.length}개`, icon: ClipboardCheck, c: 'text-violet-600 bg-violet-50' },
            { label: '미결사항', val: '2건', icon: AlertTriangle, c: 'text-amber-600 bg-amber-50' },
            { label: '전체 로그', val: `${auditData.length}건`, icon: Clock, c: 'text-slate-600 bg-slate-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between mb-3"><span className="text-xs font-semibold text-slate-500">{s.label}</span><div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div></div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>

        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
            <ClipboardCheck className="w-5 h-5 text-violet-500" />
            감사 추적 로그 (Audit Trail)
          </h2>
          <DataTable 
            data={auditData} 
            columns={columns} 
            itemsPerPage={10} 
            searchPlaceholder="로그 검색 (사용자, 수행 기록...)"
          />
        </div>

        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-bold text-slate-900 mb-6">규제 준수 매트릭스</h2>
          <div className="space-y-4">
            {compData.map((c: any) => (
              <div key={c.ruleId} className="flex items-center justify-between rounded-2xl border border-slate-100 px-6 py-5 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-4">
                  {c.status === 'pass' ? <CheckCircle className="h-6 w-6 text-emerald-500" /> : <AlertTriangle className="h-6 w-6 text-amber-500" />}
                  <div>
                    <span className="text-sm font-bold text-slate-900 block">{c.ruleId}</span>
                    <span className="text-xs text-slate-400">마지막 업데이트: 2026.04.12</span>
                  </div>
                </div>
                <div className="flex items-center gap-6">
                  <div className="flex flex-col items-end">
                    <span className="text-sm font-black text-slate-900 tabular-nums">100%</span>
                    <span className="text-[10px] font-bold text-slate-400 uppercase">Compliance</span>
                  </div>
                  <div className="w-48 h-2.5 rounded-full bg-slate-100 overflow-hidden">
                    <div className="h-full rounded-full bg-emerald-500" style={{ width: '100%' }} />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
