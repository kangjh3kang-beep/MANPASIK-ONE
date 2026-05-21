'use client';

import React, { useState, useEffect } from 'react';
import { Network, Cpu, BarChart3, GitBranch, Zap, ArrowRight, CheckCircle, Clock, AlertTriangle } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { DomainHeader } from '@mmup/ui';

const AGENTS_DATA = [
  { id: 'diag-nlp', name: '진단 NLP', model: 'BioBERT-KR v3.2', status: 'running', accuracy: 98.7, latency: 42, requests: 12847 },
  { id: 'ecg-cls', name: 'ECG 분류', model: 'ResNet-ECG v2.1', status: 'running', accuracy: 99.2, latency: 18, requests: 8923 },
  { id: 'xray-det', name: 'X-Ray 탐지', model: 'YOLO-Med v4.0', status: 'running', accuracy: 97.8, latency: 85, requests: 3421 },
  { id: 'drug-int', name: '약물 상호작용', model: 'GNN-Drug v1.8', status: 'idle', accuracy: 96.1, latency: 120, requests: 1876 },
  { id: 'risk-pred', name: '위험도 예측', model: 'XGBoost-Risk v5.0', status: 'running', accuracy: 94.5, latency: 35, requests: 6234 },
  { id: 'ocr-rx', name: '처방전 OCR', model: 'TrOCR-Med v2.4', status: 'error', accuracy: 91.3, latency: 0, requests: 425 },
];

const STATUS: Record<string, { color: string; label: string }> = {
  running: { color: 'text-emerald-600 bg-emerald-50', label: '가동 중' },
  idle: { color: 'text-amber-600 bg-amber-50', label: '대기' },
  error: { color: 'text-red-600 bg-red-50', label: '오류' },
};

export default function AgentsHubPage() {
  const { data: session } = useSession();
  const [activeStep, setActiveStep] = useState(3); // '모델 추론' 단계 시작
  const [requests, setRequests] = useState(33726);
  
  const user = session?.user as any;

  // 실시간 시뮬레이션: 파이프라인 단계 및 요청 수 가변
  useEffect(() => {
    const interval = setInterval(() => {
      setRequests(prev => prev + Math.floor(Math.random() * 5));
      setActiveStep(prev => (prev % 6) + 1);
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  const PIPELINE_STEPS = [
    { id: 1, name: '데이터 수집' },
    { id: 2, name: '전처리' },
    { id: 3, name: '특성 추출' },
    { id: 4, name: '모델 추론' },
    { id: 5, name: '후처리' },
    { id: 6, name: '결과 저장' },
  ];

  const runningCount = AGENTS_DATA.filter(a => a.status === 'running').length;

  return (
    <div className="min-h-screen bg-slate-50">
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
        {/* KPI 섹션 */}
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: '활성 에이전트', val: `${runningCount}/${AGENTS_DATA.length}`, icon: Cpu, c: 'text-emerald-600 bg-emerald-50' },
            { label: '총 추론 요청', val: requests.toLocaleString(), icon: BarChart3, c: 'text-sky-600 bg-sky-50' },
            { label: '평균 지연', val: '52ms', icon: Clock, c: 'text-amber-600 bg-amber-50' },
            { label: '활성 단계', val: `${activeStep}/6 단계`, icon: GitBranch, c: 'text-violet-600 bg-violet-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-semibold text-slate-500">{s.label}</span>
                <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-bold text-slate-900 tabular-nums transition-all">{s.val}</p>
            </div>
          ))}
        </div>

        {/* 실시간 파이프라인 시각화 */}
        <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-lg font-bold text-slate-900">추론 파이프라인 진행 상태</h2>
            <div className="px-3 py-1 bg-emerald-50 text-emerald-600 rounded-lg text-xs font-bold flex items-center gap-2">
              <div className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
              Live Inference Active
            </div>
          </div>
          
          <div className="flex flex-wrap items-center gap-3">
            {PIPELINE_STEPS.map((step, i) => {
              const state = i + 1 < activeStep ? 'done' : i + 1 === activeStep ? 'active' : 'pending';
              return (
                <React.Fragment key={step.name}>
                  <div className={`flex items-center gap-3 rounded-2xl px-5 py-3 text-sm font-bold transition-all duration-500 ${
                    state === 'done' ? 'bg-emerald-50 text-emerald-700 border border-emerald-100' :
                    state === 'active' ? 'bg-slate-900 text-white shadow-lg scale-105' :
                    'bg-slate-50 text-slate-400 border border-transparent'
                  }`}>
                    {state === 'done' ? <CheckCircle className="h-4 w-4" /> :
                     state === 'active' ? <Zap className="h-4 w-4 animate-pulse text-yellow-400" /> :
                     <Clock className="h-4 w-4" />}
                    {step.name}
                  </div>
                  {i < PIPELINE_STEPS.length - 1 && (
                    <ArrowRight className={`h-4 w-4 flex-shrink-0 transition-colors duration-500 ${i + 1 < activeStep ? 'text-emerald-500' : 'text-slate-200'}`} />
                  )}
                </React.Fragment>
              );
            })}
          </div>
        </div>

        {/* 에이전트 상세 목록 */}
        <div className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
          <h2 className="text-xl font-bold text-slate-900 mb-6">등록된 지능형 에이전트 인벤토리</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {AGENTS_DATA.map(agent => {
              const st = STATUS[agent.status];
              return (
                <div key={agent.id} className="group rounded-2xl border border-slate-100 bg-slate-50/50 p-5 hover:border-sky-200 hover:bg-white hover:shadow-xl transition-all duration-300">
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
                      { l: '요청 수', v: agent.requests.toLocaleString(), color: 'text-slate-600' },
                    ].map(x => (
                      <div key={x.l} className="rounded-xl bg-white border border-slate-100 p-2.5 text-center shadow-sm">
                        <p className="text-[10px] font-bold text-slate-400 mb-1">{x.l}</p>
                        <p className={`text-sm font-black ${x.color} tabular-nums`}>{x.v}</p>
                      </div>
                    ))}
                  </div>

                  {agent.status === 'error' && (
                    <div className="mt-4 flex items-center gap-2 rounded-lg bg-red-50 p-2 text-[10px] font-bold text-red-600">
                      <AlertTriangle className="h-3 w-3" />
                      자원 할단 오류: 즉시 복구가 필요합니다.
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

