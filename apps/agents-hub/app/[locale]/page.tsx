'use client';

import { DomainNav } from '@mmup/ui';

import React, { useState, useEffect } from 'react';
import { Network, Play, Pause, CheckCircle, XCircle, Clock, Cpu, BarChart3, GitBranch, Zap, ArrowRight, Activity } from 'lucide-react';

const AGENTS = [
  { id: 'diag-nlp', name: '진단 NLP', model: 'BioBERT-KR v3.2', status: 'running', accuracy: 98.7, latency: 42, requests: 12847 },
  { id: 'ecg-cls', name: 'ECG 분류', model: 'ResNet-ECG v2.1', status: 'running', accuracy: 99.2, latency: 18, requests: 8923 },
  { id: 'xray-det', name: 'X-Ray 탐지', model: 'YOLO-Med v4.0', status: 'running', accuracy: 97.8, latency: 85, requests: 3421 },
  { id: 'drug-int', name: '약물 상호작용', model: 'GNN-Drug v1.8', status: 'idle', accuracy: 96.1, latency: 120, requests: 1876 },
  { id: 'risk-pred', name: '위험도 예측', model: 'XGBoost-Risk v5.0', status: 'running', accuracy: 94.5, latency: 35, requests: 6234 },
  { id: 'ocr-rx', name: '처방전 OCR', model: 'TrOCR-Med v2.4', status: 'error', accuracy: 91.3, latency: 0, requests: 425 },
];

const STATUS = {
  running: { icon: Play, color: 'text-emerald-600 bg-emerald-50', label: '가동 중' },
  idle: { icon: Pause, color: 'text-amber-600 bg-amber-50', label: '대기' },
  error: { icon: XCircle, color: 'text-red-600 bg-red-50', label: '오류' },
};

const PIPELINE_STEPS = [
  { name: '데이터 수집', status: 'done' },
  { name: '전처리', status: 'done' },
  { name: '특성 추출', status: 'done' },
  { name: '모델 추론', status: 'active' },
  { name: '후처리', status: 'pending' },
  { name: '결과 저장', status: 'pending' },
];

export default function AgentsHubPage() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;

  const running = AGENTS.filter(a => a.status === 'running').length;

  return (
    <>
      <DomainNav currentDomain="agents-hub" />
      <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-600 text-white">
            <Network className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">AI 에이전트 오케스트레이터</h1>
        </div>
        <p className="text-sm text-slate-500">의료 AI 모델 파이프라인 실시간 상태 모니터링 및 관리</p>
      </div>

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        {/* KPI */}
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: '활성 에이전트', val: `${running}/${AGENTS.length}`, icon: Cpu, c: 'text-emerald-600 bg-emerald-50' },
            { label: '총 추론 요청', val: '33,726', icon: BarChart3, c: 'text-sky-600 bg-sky-50' },
            { label: '평균 지연', val: '52ms', icon: Clock, c: 'text-amber-600 bg-amber-50' },
            { label: '파이프라인', val: '4/6 단계', icon: GitBranch, c: 'text-violet-600 bg-violet-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-semibold text-slate-500">{s.label}</span>
                <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>

        {/* Pipeline */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">추론 파이프라인 진행 상태</h2>
          <div className="flex items-center gap-2">
            {PIPELINE_STEPS.map((step, i) => (
              <React.Fragment key={step.name}>
                <div className={`flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-medium ${
                  step.status === 'done' ? 'bg-emerald-50 text-emerald-700' :
                  step.status === 'active' ? 'bg-sky-50 text-sky-700 ring-2 ring-sky-500 animate-pulse' :
                  'bg-slate-50 text-slate-400'
                }`}>
                  {step.status === 'done' ? <CheckCircle className="h-4 w-4" /> :
                   step.status === 'active' ? <Zap className="h-4 w-4" /> :
                   <Clock className="h-4 w-4" />}
                  {step.name}
                </div>
                {i < PIPELINE_STEPS.length - 1 && <ArrowRight className="h-4 w-4 text-slate-300 flex-shrink-0" />}
              </React.Fragment>
            ))}
          </div>
        </div>

        {/* Agent Cards */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">등록된 AI 에이전트</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {AGENTS.map(agent => {
              const st = STATUS[agent.status as keyof typeof STATUS];
              return (
                <div key={agent.id} className="rounded-xl border border-slate-200 p-4 hover:shadow-md transition-shadow">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-sm font-bold text-slate-900">{agent.name}</h3>
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${st.color}`}>{st.label}</span>
                  </div>
                  <p className="text-xs text-slate-400 mb-3">{agent.model}</p>
                  <div className="grid grid-cols-3 gap-2 text-center">
                    <div className="rounded-lg bg-slate-50 p-2">
                      <p className="text-xs text-slate-400">정확도</p>
                      <p className="text-sm font-bold text-slate-700">{agent.accuracy}%</p>
                    </div>
                    <div className="rounded-lg bg-slate-50 p-2">
                      <p className="text-xs text-slate-400">지연</p>
                      <p className="text-sm font-bold text-slate-700">{agent.latency}ms</p>
                    </div>
                    <div className="rounded-lg bg-slate-50 p-2">
                      <p className="text-xs text-slate-400">요청</p>
                      <p className="text-sm font-bold text-slate-700">{agent.requests.toLocaleString()}</p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
    </>
  );
}
