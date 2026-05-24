'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Play, Code, Clock, CheckCircle, Loader2, Search } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

interface Endpoint {
  method: HttpMethod;
  path: string;
  label: string;
  sampleBody?: string;
}

const METHOD_STYLES: Record<HttpMethod, string> = {
  GET: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  POST: 'bg-sky-50 text-sky-700 border-sky-200',
  PUT: 'bg-amber-50 text-amber-700 border-amber-200',
  DELETE: 'bg-red-50 text-red-700 border-red-200',
};

const ENDPOINTS: Endpoint[] = [
  { method: 'GET', path: '/api/v1/patients', label: '환자 목록 조회' },
  { method: 'GET', path: '/api/v1/measurements/history', label: '측정 이력 조회' },
  { method: 'POST', path: '/api/v1/predictions/run', label: '예측 실행', sampleBody: '{\n  "patient_id": "PAT-001",\n  "model": "glucose_predictor",\n  "features": [5.2, 3.1, 98]\n}' },
  { method: 'GET', path: '/api/v1/devices', label: '장비 목록 조회' },
  { method: 'POST', path: '/api/v1/measurements', label: '측정 결과 등록', sampleBody: '{\n  "patient_id": "PAT-001",\n  "cartridge_id": "CRT-GLU-044",\n  "value": 112,\n  "unit": "mg/dL"\n}' },
  { method: 'PUT', path: '/api/v1/patients/{id}', label: '환자 정보 수정', sampleBody: '{\n  "name": "홍길동",\n  "phone": "010-1234-5678"\n}' },
  { method: 'DELETE', path: '/api/v1/webhooks/{id}', label: '웹훅 삭제' },
  { method: 'GET', path: '/api/v1/cartridges/inventory', label: '카트리지 재고 조회' },
];

interface HistoryEntry {
  method: HttpMethod;
  path: string;
  status: number;
  time: string;
}

const SAMPLE_RESPONSES: Record<number, object> = {
  200: { success: true, data: { id: 'PAT-001', name: '김철수', measurements: 24, lastVisit: '2026-05-23' }, meta: { requestId: 'req-abc123', latency: '23ms' } },
  201: { success: true, data: { id: 'MSR-0045', value: 112, unit: 'mg/dL', confidence: 0.96 }, meta: { requestId: 'req-def456', latency: '45ms' } },
};

export default function PlaygroundPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [selectedIdx, setSelectedIdx] = useState(0);
  const [requestBody, setRequestBody] = useState(ENDPOINTS[0].sampleBody || '');
  const [isLoading, setIsLoading] = useState(false);
  const [response, setResponse] = useState<{ status: number; body: object } | null>(null);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [copied, setCopied] = useState(false);

  const endpoint = ENDPOINTS[selectedIdx];

  const handleSelectEndpoint = (idx: number) => {
    setSelectedIdx(idx);
    setRequestBody(ENDPOINTS[idx].sampleBody || '');
    setResponse(null);
  };

  const handleExecute = () => {
    setIsLoading(true);
    setTimeout(() => {
      const status = endpoint.method === 'POST' ? 201 : endpoint.method === 'DELETE' ? 204 : 200;
      const body = SAMPLE_RESPONSES[status] || { success: true, message: '처리 완료' };
      setResponse({ status, body });
      setHistory(prev => [
        { method: endpoint.method, path: endpoint.path, status, time: new Date().toLocaleTimeString('ko-KR') },
        ...prev,
      ].slice(0, 3));
      setIsLoading(false);
    }, 800);
  };

  const curlCommand = `curl -X ${endpoint.method} \\\n  "https://api.manpasik.com${endpoint.path}" \\\n  -H "Authorization: Bearer mps_sk_xxxx" \\\n  -H "Content-Type: application/json"${endpoint.sampleBody ? ` \\\n  -d '${endpoint.sampleBody.replace(/\n/g, '')}'` : ''}`;

  const statusStyle = (code: number) =>
    code < 300 ? 'bg-emerald-50 text-emerald-700' :
    code < 500 ? 'bg-amber-50 text-amber-700' :
    'bg-red-50 text-red-700';

  return (
    <main className="min-h-screen bg-slate-50" aria-label="API 플레이그라운드">
      <DomainHeader
        title="API 플레이그라운드"
        subtitle="API 엔드포인트를 선택하고 직접 호출해 보세요"
        icon={Play}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'dev' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-6xl mx-auto px-6 py-8">
        <div className="grid grid-cols-12 gap-6">
          {/* 좌측: 요청 패널 */}
          <div className="col-span-8 space-y-4">
            {/* 엔드포인트 선택 */}
            <div className="rounded-2xl border border-slate-200 bg-white p-5">
              <label className="block text-xs font-bold text-slate-500 mb-2">엔드포인트 선택</label>
              <select
                value={selectedIdx}
                onChange={e => handleSelectEndpoint(Number(e.target.value))}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm font-mono text-slate-700 focus:outline-none focus:ring-2 focus:ring-sky-500"
                aria-label="엔드포인트 선택"
              >
                {ENDPOINTS.map((ep, i) => (
                  <option key={i} value={i}>{ep.method} {ep.path} — {ep.label}</option>
                ))}
              </select>

              <div className="flex items-center gap-3 mt-4">
                <span className={`px-3 py-1 rounded-full text-xs font-bold border ${METHOD_STYLES[endpoint.method]}`}>
                  {endpoint.method}
                </span>
                <code className="flex-1 text-sm font-mono text-slate-700 bg-slate-50 rounded-lg px-3 py-2">
                  https://api.manpasik.com{endpoint.path}
                </code>
              </div>
            </div>

            {/* 헤더 & 바디 */}
            <div className="rounded-2xl border border-slate-200 bg-white p-5">
              <h3 className="text-xs font-bold text-slate-500 mb-3">요청 헤더</h3>
              <div className="bg-slate-900 rounded-xl p-4 mb-4">
                <code className="text-xs text-emerald-400 font-mono block">Authorization: Bearer mps_sk_xxxx</code>
                <code className="text-xs text-emerald-400 font-mono block mt-1">Content-Type: application/json</code>
              </div>

              {(endpoint.method === 'POST' || endpoint.method === 'PUT') && (
                <>
                  <h3 className="text-xs font-bold text-slate-500 mb-2">요청 본문</h3>
                  <textarea
                    value={requestBody}
                    onChange={e => setRequestBody(e.target.value)}
                    rows={6}
                    className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm font-mono text-slate-700 focus:outline-none focus:ring-2 focus:ring-sky-500 resize-none"
                    aria-label="요청 본문"
                  />
                </>
              )}

              <button
                onClick={handleExecute}
                disabled={isLoading}
                className="mt-4 w-full flex items-center justify-center gap-2 rounded-xl bg-sky-600 px-6 py-3 text-sm font-bold text-white hover:bg-sky-700 transition-colors disabled:opacity-50"
                aria-label="API 실행"
              >
                {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />}
                {isLoading ? '실행 중...' : '실행'}
              </button>
            </div>

            {/* 응답 패널 */}
            {response && (
              <div className="rounded-2xl border border-slate-200 bg-white p-5">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-xs font-bold text-slate-500">응답</h3>
                  <span className={`px-3 py-1 rounded-full text-xs font-bold ${statusStyle(response.status)}`}>
                    {response.status}
                  </span>
                </div>
                <pre className="bg-slate-900 rounded-xl p-4 text-xs text-emerald-400 font-mono overflow-x-auto whitespace-pre-wrap">
                  {JSON.stringify(response.body, null, 2)}
                </pre>
              </div>
            )}

            {/* cURL 명령어 */}
            <div className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-xs font-bold text-slate-500 flex items-center gap-2">
                  <Code className="h-3.5 w-3.5" /> cURL 명령어
                </h3>
                <button
                  onClick={() => { navigator.clipboard.writeText(curlCommand); setCopied(true); setTimeout(() => setCopied(false), 2000); }}
                  className="text-xs text-sky-600 hover:text-sky-800 font-semibold flex items-center gap-1"
                  aria-label="cURL 명령어 복사"
                >
                  {copied ? <><CheckCircle className="h-3.5 w-3.5" /> 복사됨</> : '복사'}
                </button>
              </div>
              <pre className="bg-slate-900 rounded-xl p-4 text-xs text-slate-300 font-mono overflow-x-auto whitespace-pre-wrap">
                {curlCommand}
              </pre>
            </div>
          </div>

          {/* 우측: 요청 기록 */}
          <div className="col-span-4">
            <div className="rounded-2xl border border-slate-200 bg-white p-5 sticky top-8">
              <h3 className="text-xs font-bold text-slate-500 flex items-center gap-2 mb-4">
                <Clock className="h-3.5 w-3.5" /> 최근 요청 기록
              </h3>
              {history.length === 0 ? (
                <p className="text-xs text-slate-400 text-center py-6">아직 요청 기록이 없습니다</p>
              ) : (
                <div className="space-y-2">
                  {history.map((h, i) => (
                    <div key={i} className="rounded-xl bg-slate-50 border border-slate-100 p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${METHOD_STYLES[h.method]}`}>
                          {h.method}
                        </span>
                        <span className={`text-[10px] font-bold ${statusStyle(h.status)} px-2 py-0.5 rounded`}>
                          {h.status}
                        </span>
                      </div>
                      <code className="text-[11px] font-mono text-slate-600 block truncate">{h.path}</code>
                      <p className="text-[10px] text-slate-400 mt-1">{h.time}</p>
                    </div>
                  ))}
                </div>
              )}

              {/* 검색 팁 */}
              <div className="mt-6 rounded-xl bg-sky-50 border border-sky-100 p-3">
                <h4 className="text-[11px] font-bold text-sky-700 flex items-center gap-1 mb-1">
                  <Search className="h-3 w-3" /> 도움말
                </h4>
                <p className="text-[10px] text-sky-600 leading-relaxed">
                  엔드포인트를 선택하고 실행 버튼을 누르면 모의 응답을 확인할 수 있습니다.
                  실제 API 키는 개발자 도구 페이지에서 발급받으세요.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
