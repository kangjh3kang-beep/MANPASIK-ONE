'use client';

import { DomainNav } from '@mmup/ui';

import React, { useState, useEffect } from 'react';
import { FileCode2, Terminal, Key, Webhook, Copy, CheckCircle, ExternalLink, BookOpen } from 'lucide-react';

const API_ENDPOINTS = [
  { method: 'GET', path: '/api/v1/patients', desc: '환자 목록 조회', auth: 'Bearer', status: 'stable' },
  { method: 'POST', path: '/api/v1/measurements', desc: '측정 데이터 전송', auth: 'Bearer', status: 'stable' },
  { method: 'GET', path: '/api/v1/biomarkers/:id', desc: '바이오마커 조회', auth: 'Bearer', status: 'stable' },
  { method: 'POST', path: '/api/v1/predictions/run', desc: 'AI 예측 실행', auth: 'API Key', status: 'beta' },
  { method: 'GET', path: '/api/v1/rewards/balance', desc: '토큰 잔액 조회', auth: 'Bearer', status: 'stable' },
  { method: 'POST', path: '/api/v1/fhir/bundle', desc: 'FHIR 번들 전송', auth: 'mTLS', status: 'beta' },
];

const METHOD_COLORS: Record<string, string> = {
  GET: 'text-emerald-700 bg-emerald-50',
  POST: 'text-sky-700 bg-sky-50',
  PUT: 'text-amber-700 bg-amber-50',
  DELETE: 'text-red-700 bg-red-50',
};

export default function DevPortalPage() {
  const [mounted, setMounted] = useState(false);
  const [copied, setCopied] = useState(false);

  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;

  const apiKey = 'mmup_live_sk_x7k9m2n4p6q8r0t2v4w6y8a0b2c4d6e8';

  return (
    <>
      <DomainNav currentDomain="dev-portal" />
      <div className="min-h-screen bg-slate-50">
      <div className="border-b border-slate-200 bg-white px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-800 text-white">
            <FileCode2 className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">개발자 포털</h1>
        </div>
        <p className="text-sm text-slate-500">MMUP Open API 명세서, SDK 다운로드 및 실시간 샌드박스</p>
      </div>

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: 'API 엔드포인트', val: '24개', icon: Terminal, c: 'text-slate-600 bg-slate-100' },
            { label: '활성 API 키', val: '3개', icon: Key, c: 'text-amber-600 bg-amber-50' },
            { label: '웹훅', val: '7개', icon: Webhook, c: 'text-violet-600 bg-violet-50' },
            { label: 'API 문서', val: 'OpenAPI 3.1', icon: BookOpen, c: 'text-sky-600 bg-sky-50' },
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

        {/* API Key */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-3">나의 API 키</h2>
          <div className="flex items-center gap-2 rounded-xl bg-slate-900 px-4 py-3">
            <code className="flex-1 text-sm text-emerald-400 font-mono">{apiKey}</code>
            <button
              onClick={() => { navigator.clipboard.writeText(apiKey); setCopied(true); setTimeout(() => setCopied(false), 2000); }}
              className="text-slate-400 hover:text-white transition-colors"
            >
              {copied ? <CheckCircle className="h-4 w-4 text-emerald-400" /> : <Copy className="h-4 w-4" />}
            </button>
          </div>
        </div>

        {/* API Explorer */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">API 엔드포인트 탐색기</h2>
          <div className="space-y-2">
            {API_ENDPOINTS.map((ep, i) => (
              <div key={i} className="flex items-center gap-4 rounded-xl border border-slate-100 px-4 py-3 hover:bg-slate-50 transition-colors cursor-pointer">
                <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${METHOD_COLORS[ep.method]}`}>{ep.method}</span>
                <code className="text-sm font-mono text-slate-700 flex-1">{ep.path}</code>
                <span className="text-xs text-slate-400">{ep.desc}</span>
                <span className="text-[10px] bg-slate-100 text-slate-500 px-1.5 py-0.5 rounded">{ep.auth}</span>
                <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${ep.status === 'stable' ? 'text-emerald-700 bg-emerald-50' : 'text-amber-700 bg-amber-50'}`}>{ep.status}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
    </>
  );
}
