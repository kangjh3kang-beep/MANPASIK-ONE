'use client';

import React from 'react';
import { Network, Cpu, BarChart3, GitBranch, Zap, ArrowRight, CheckCircle, Clock, AlertTriangle } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { useAgents, useAgentPipeline, AIAgent } from '@mmup/api-client';
import { DomainHeader, KPICard, ErrorState, LoadingSkeleton } from '@mmup/ui';

const STATUS: Record<string, { color: string; label: string }> = {
  running: { color: 'text-emerald-600 bg-emerald-50', label: '가동 중' },
  idle: { color: 'text-amber-600 bg-amber-50', label: '대기' },
  error: { color: 'text-red-600 bg-red-50', label: '오류' },
};

export default function AgentsHubPage() {
  const { data: session } = useSession();
  const { data: agentsData, isLoading: agentsLoading, error: agentsError, refetch: refetchAgents } = useAgents();
  const { data: pipelineData, isLoading: pipelineLoading, error: pipelineError } = useAgentPipeline();

  const user = session?.user as any;

  if (agentsLoading || pipelineLoading) {
    return <LoadingSkeleton columns={4} rows={1} />;
  }

  if (agentsError || pipelineError) {
    return <ErrorState message="에이전트 데이터를 불러오는 중 오류가 발생했습니다." onRetry={() => refetchAgents()} />;
  }

  const agents: AIAgent[] = agentsData?.data || [];
  const pipeline = pipelineData?.data;
  const activeStep = pipeline?.activeStep || 1;
  const totalRequests = pipeline?.totalRequests || 0;
  const pipelineSteps = pipeline?.steps || [];
  const runningCount = agents.filter(a => a.status === 'running').length;
  const avgLatency = agents.filter(a => a.status !== 'error').length > 0
    ? Math.round(agents.filter(a => a.status !== 'error').reduce((sum, a) => sum + a.latency, 0) / agents.filter(a => a.status !== 'error').length)
    : 0;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="AI 에이전트 오케스트레이터">
      <DomainHeader
        title="AI 에이전트 오케스트레이터"
        subtitle="의료 AI 모델 파이프라인 실시간 상태 모니터링 및 성능 관리"
        icon={Network}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'researcher',
          organization: user.organization || 'MMUP AI Center'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <section aria-label="핵심 성과 지표">
          <div className="grid grid-cols-4 gap-4">
            <KPICard title="활성 에이전트" value={`${runningCount}/${agents.length}`} icon={<Cpu className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
            <KPICard title="총 추론 요청" value={totalRequests.toLocaleString()} icon={<BarChart3 className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
            <KPICard title="평균 지연" value={`${avgLatency}ms`} icon={<Clock className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            <KPICard title="활성 단계" value={`${activeStep}/6 단계`} icon={<GitBranch className="h-4 w-4" />} color="text-violet-600 bg-violet-50" />
          </div>
        </section>

        <section aria-label="추론 파이프라인 진행 상태">
          <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <div className="flex items-center justify-between mb-8">
              <h2 className="text-lg font-bold text-slate-900">추론 파이프라인 진행 상태</h2>
              <div className="px-3 py-1 bg-emerald-50 text-emerald-600 rounded-lg text-xs font-bold flex items-center gap-2">
                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" aria-hidden="true" />
                Live Inference Active
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-3" role="status" aria-label={`파이프라인 ${activeStep}단계 진행 중`}>
              {pipelineSteps.map((step, i) => {
                const state = i + 1 < activeStep ? 'done' : i + 1 === activeStep ? 'active' : 'pending';
                return (
                  <React.Fragment key={step.name}>
                    <div className={`flex items-center gap-3 rounded-2xl px-5 py-3 text-sm font-bold transition-all duration-500 ${
                      state === 'done' ? 'bg-emerald-50 text-emerald-700 border border-emerald-100' :
                      state === 'active' ? 'bg-slate-900 text-white shadow-lg scale-105' :
                      'bg-slate-50 text-slate-400 border border-transparent'
                    }`} aria-label={`${step.name}: ${state === 'done' ? '완료' : state === 'active' ? '진행 중' : '대기'}`}>
                      {state === 'done' ? <CheckCircle className="h-4 w-4" aria-hidden="true" /> :
                       state === 'active' ? <Zap className="h-4 w-4 animate-pulse text-yellow-400" aria-hidden="true" /> :
                       <Clock className="h-4 w-4" aria-hidden="true" />}
                      {step.name}
                    </div>
                    {i < pipelineSteps.length - 1 && (
                      <ArrowRight className={`h-4 w-4 flex-shrink-0 transition-colors duration-500 ${i + 1 < activeStep ? 'text-emerald-500' : 'text-slate-200'}`} aria-hidden="true" />
                    )}
                  </React.Fragment>
                );
              })}
            </div>
          </div>
        </section>

        <section aria-label="등록된 에이전트 목록">
          <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <h2 className="text-xl font-bold text-slate-900 mb-6">등록된 지능형 에이전트 인벤토리</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {agents.map(agent => {
                const st = STATUS[agent.status] || STATUS.idle;
                return (
                  <div key={agent.id} className="group rounded-2xl border border-slate-100 bg-slate-50/50 p-5 hover:border-sky-200 hover:bg-white hover:shadow-xl transition-all duration-300" aria-label={`${agent.name} - ${st.label}`}>
                    <div className="flex items-center justify-between mb-4">
                      <div>
                        <h3 className="text-sm font-black text-slate-900 group-hover:text-sky-600 transition-colors">{agent.name}</h3>
                        <p className="text-[10px] font-bold text-slate-400 uppercase tracking-tighter mt-0.5">{agent.model}</p>
                      </div>
                      <span className={`text-[10px] font-black px-2 py-1 rounded-lg shadow-sm ${st.color}`}>{st.label}</span>
                    </div>

                    <div className="grid grid-cols-3 gap-3">
                      {[
                        { l: '정확도', v: `${agent.accuracy}%`, color: 'text-emerald-600' },
                        { l: '지연', v: agent.status === 'error' ? 'N/A' : `${agent.latency}ms`, color: 'text-sky-600' },
                        { l: '요청 수', v: agent.totalRequests.toLocaleString(), color: 'text-slate-600' },
                      ].map(x => (
                        <div key={x.l} className="rounded-xl bg-white border border-slate-100 p-2.5 text-center shadow-sm">
                          <p className="text-[10px] font-bold text-slate-400 mb-1">{x.l}</p>
                          <p className={`text-sm font-black ${x.color} tabular-nums`} role="status">{x.v}</p>
                        </div>
                      ))}
                    </div>

                    {agent.status === 'error' && (
                      <div className="mt-4 flex items-center gap-2 rounded-lg bg-red-50 p-2 text-[10px] font-bold text-red-600" role="alert">
                        <AlertTriangle className="h-3 w-3" aria-hidden="true" />
                        자원 할당 오류: 즉시 복구가 필요합니다.
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
