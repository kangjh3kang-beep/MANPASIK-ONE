'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React from 'react';
import { Pill, ClipboardCheck, Shield, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
  const localePrefix = useLocalePrefix();
import { useAuditLogs, useCompliance, AuditLog } from '@mmup/api-client';
import { DataTable, Column, DomainHeader, KPICard, LoadingSkeleton, ErrorState } from '@mmup/ui';

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function GxPPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const { data: logs, isLoading: logsLoading, error: logsError, refetch } = useAuditLogs();
  const { data: compliances, isLoading: compLoading, error: compError } = useCompliance();
  const user = session?.user as any;

  if (logsLoading || compLoading) {
    return <LoadingSkeleton columns={4} rows={1} />;
  }

  if (logsError || compError) {
    return <ErrorState message="GxP 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

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
    <main className="min-h-screen bg-slate-50" aria-label="의약품 GxP 준수 시스템">
      <DomainHeader
        title="의약품 GxP 준수 시스템"
        subtitle="cGMP 감사 추적, 전자 서명 및 실시간 규제 준수 모니터링"
        icon={Pill}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'pharma',
          organization: user.organization || 'MMUP 제약'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <KPICard title="준수율" value="99.4%" icon={<Shield className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="감사항목" value={`${compData.length}개`} icon={<ClipboardCheck className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
            <KPICard title="미결사항" value="2건" icon={<AlertTriangle className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="전체 로그" value={`${auditData.length}건`} icon={<Clock className="h-4 w-4" />} color="text-slate-600 bg-slate-50" />
          </div>
        </section>

        <section aria-label="감사 추적 로그">
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
              <ClipboardCheck className="w-5 h-5 text-violet-500" aria-hidden="true" />
              감사 추적 로그 (Audit Trail)
            </h2>
            <DataTable
              data={auditData}
              columns={columns}
              itemsPerPage={10}
              searchPlaceholder="로그 검색 (사용자, 수행 기록...)"
            />
          </div>
        </section>

        <section aria-label="규제 준수 매트릭스">
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold text-slate-900 mb-6">규제 준수 매트릭스</h2>
            <div className="space-y-4">
              {compData.map((c: any) => (
                <div key={c.ruleId} className="flex items-center justify-between rounded-2xl border border-slate-100 px-6 py-5 hover:bg-slate-50 transition-colors" aria-label={`${c.ruleId} - ${c.status === 'pass' ? '통과' : '경고'}`}>
                  <div className="flex items-center gap-4">
                    {c.status === 'pass' ? <CheckCircle className="h-6 w-6 text-emerald-500" aria-hidden="true" /> : <AlertTriangle className="h-6 w-6 text-amber-500" aria-hidden="true" />}
                    <div>
                      <span className="text-sm font-bold text-slate-900 block">{c.ruleId}</span>
                      <span className="text-xs text-slate-400">마지막 업데이트: 2026.04.12</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-6">
                    <div className="flex flex-col items-end">
                      <span className="text-sm font-black text-slate-900 tabular-nums" role="status">100%</span>
                      <span className="text-[10px] font-bold text-slate-400 uppercase">Compliance</span>
                    </div>
                    <div className="w-48 h-2.5 rounded-full bg-slate-100 overflow-hidden" role="progressbar" aria-valuenow={100} aria-valuemin={0} aria-valuemax={100}>
                      <div className="h-full rounded-full bg-emerald-500" style={{ width: '100%' }} />
                    </div>
                  </div>
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
          {[{ name: "파트너 연동", path: "/domains/partner", desc: "규제 데이터 교환" },{ name: "임상 콘솔", path: "/domains/clinical", desc: "감사 추적 연동" },{ name: "하드웨어 코어", path: "/domains/hardware-core", desc: "IEC 62304 준수" }].map(d => (
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
