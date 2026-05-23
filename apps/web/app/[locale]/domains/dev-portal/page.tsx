'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React, { useState } from 'react';
import { FileCode2, Terminal, Key, Webhook, Copy, CheckCircle, BookOpen } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { useDevPortalKeys, useDevPortalEndpoints, DevPortalKey, DevPortalEndpoint } from '@mmup/api-client';
import { DomainHeader, KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

const METHOD_COLORS: Record<string, string> = { GET: 'text-emerald-700 bg-emerald-50', POST: 'text-sky-700 bg-sky-50' };

function useLocalePrefix(): string { const pathname = usePathname(); const match = pathname.match(/^\/(ko|en|ja|zh)/); return match ? `/${match[1]}` : '/ko'; }

export default function DevPortalPage() {
  const { data: session } = useSession();
  const localePrefix = useLocalePrefix();
  const [copied, setCopied] = useState<string | null>(null);
  const user = session?.user;

  const { data: keysData, isLoading: keysLoading, error: keysError, refetch } = useDevPortalKeys();
  const { data: endpointsData, isLoading: endpointsLoading, error: endpointsError } = useDevPortalEndpoints();

  if (keysLoading || endpointsLoading) {
    return <LoadingSkeleton columns={4} rows={1} />;
  }

  if (keysError || endpointsError) {
    return <ErrorState message="개발자 도구 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetch()} />;
  }

  const keys: DevPortalKey[] = keysData?.data || [];
  const endpoints: DevPortalEndpoint[] = endpointsData?.data || [];
  const activeKeys = keys.filter(k => k.status === 'active');

  return (
    <main className="min-h-screen bg-slate-50" aria-label="개발자 도구">
      <DomainHeader
        title="개발자 도구"
        subtitle="외부 개발자가 만파식 기능을 활용할 수 있는 연동 도구입니다"
        icon={FileCode2}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'dev',
          organization: user.organization || 'MPS 개발자 센터'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-4 gap-4">
            <KPICard title="API 엔드포인트" value={`${endpoints.length}개`} icon={<Terminal className="h-4 w-4" />} color="text-slate-600 bg-slate-100" />
            <KPICard title="활성 API 키" value={`${activeKeys.length}개`} icon={<Key className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="웹훅" value="7개" icon={<Webhook className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
            <KPICard title="API 문서" value="OpenAPI 3.1" icon={<BookOpen className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
          </div>
        </section>

        <section aria-label="API 키 관리">
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-3">나의 API 키</h2>
            <div className="space-y-2">
              {activeKeys.map(k => (
                <div key={k.id} className="flex items-center gap-2 rounded-xl bg-slate-900 px-4 py-3">
                  <div className="flex-1">
                    <p className="text-[10px] text-slate-400 font-bold mb-0.5">{k.name}</p>
                    <code className="text-sm text-emerald-400 font-mono">{k.key}</code>
                  </div>
                  <button
                    onClick={() => { navigator.clipboard.writeText(k.key); setCopied(k.id); setTimeout(() => setCopied(null), 2000); }}
                    className="text-slate-400 hover:text-white transition-colors"
                    aria-label={`${k.name} API 키 복사`}
                  >
                    {copied === k.id ? <CheckCircle className="h-4 w-4 text-emerald-400" /> : <Copy className="h-4 w-4" />}
                  </button>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section aria-label="API 엔드포인트 탐색기">
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4">API 엔드포인트 탐색기</h2>
            <div className="space-y-2">
              {endpoints.map((ep, i) => (
                <div key={i} className="flex items-center gap-4 rounded-xl border border-slate-100 px-4 py-3 hover:bg-slate-50 transition-colors cursor-pointer" aria-label={`${ep.method} ${ep.path} - ${ep.description}`}>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${METHOD_COLORS[ep.method] || 'text-slate-600 bg-slate-50'}`}>{ep.method}</span>
                  <code className="text-sm font-mono text-slate-700 flex-1">{ep.path}</code>
                  <span className="text-xs text-slate-400">{ep.description}</span>
                  <span className="text-[10px] bg-slate-100 text-slate-500 px-1.5 py-0.5 rounded">{ep.auth}</span>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${ep.status === 'stable' ? 'text-emerald-700 bg-emerald-50' : 'text-amber-700 bg-amber-50'}`}>{ep.status}</span>
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
          {[{ name: "하드웨어 코어", path: "/domains/hardware-core", desc: "SDK 연동" },{ name: "AI 에이전트", path: "/domains/agents-hub", desc: "AI API 활용" },{ name: "파트너 연동", path: "/domains/partner", desc: "FHIR 연동 API" }].map(d => (
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
