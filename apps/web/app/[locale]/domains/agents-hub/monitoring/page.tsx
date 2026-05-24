'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Activity, AlertTriangle, CheckCircle, Clock, TrendingUp, Zap, RefreshCw, Database } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

interface ModelMetrics {
  id: string;
  name: string;
  requestsPerMin: number;
  avgLatency: number;
  p95Latency: number;
  errorRate: number;
  driftScore: number;
  sparkline: number[];
}

interface RecentError {
  id: string;
  model: string;
  inputSummary: string;
  errorType: string;
  timestamp: string;
}

const OVERALL_STATS = [
  { label: '총 추론 횟수', value: '45,230', icon: Activity, color: 'bg-sky-50 text-sky-600' },
  { label: '평균 응답시간', value: '23ms', icon: Clock, color: 'bg-emerald-50 text-emerald-600' },
  { label: '성공률', value: '99.8%', icon: CheckCircle, color: 'bg-emerald-50 text-emerald-600' },
  { label: '활성 모델', value: '5', icon: Database, color: 'bg-violet-50 text-violet-600' },
];

const MODELS: ModelMetrics[] = [
  { id: 'mm-1', name: '혈당 예측', requestsPerMin: 42, avgLatency: 18, p95Latency: 35, errorRate: 0.02, driftScore: 0.08, sparkline: [18, 20, 19, 22, 18, 17, 19, 21] },
  { id: 'mm-2', name: '심전도 분석', requestsPerMin: 28, avgLatency: 45, p95Latency: 82, errorRate: 0.05, driftScore: 0.12, sparkline: [42, 44, 48, 46, 45, 43, 47, 45] },
  { id: 'mm-3', name: '영양 분석', requestsPerMin: 15, avgLatency: 12, p95Latency: 28, errorRate: 0.01, driftScore: 0.22, sparkline: [10, 11, 13, 14, 12, 11, 12, 12] },
  { id: 'mm-4', name: '위험 예측', requestsPerMin: 35, avgLatency: 28, p95Latency: 55, errorRate: 0.03, driftScore: 0.06, sparkline: [25, 27, 30, 28, 26, 29, 28, 28] },
  { id: 'mm-5', name: 'OCR 처리', requestsPerMin: 8, avgLatency: 65, p95Latency: 120, errorRate: 0.12, driftScore: 0.18, sparkline: [60, 62, 68, 70, 65, 63, 66, 65] },
];

const RECENT_ERRORS: RecentError[] = [
  { id: 'err-1', model: 'OCR 처리', inputSummary: '저해상도 처방전 이미지 (320x240)', errorType: '입력 해상도 부족', timestamp: '2026-05-24 09:12:33' },
  { id: 'err-2', model: '심전도 분석', inputSummary: '12리드 ECG 신호 (노이즈 과다)', errorType: '신호 품질 미달', timestamp: '2026-05-24 08:45:10' },
  { id: 'err-3', model: '영양 분석', inputSummary: '미등록 식품 코드 (FD-99999)', errorType: '참조 데이터 누락', timestamp: '2026-05-24 07:22:58' },
];

export default function MonitoringPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [retrainingId, setRetrainingId] = useState<string | null>(null);

  const handleRetrain = (id: string) => {
    setRetrainingId(id);
    setTimeout(() => setRetrainingId(null), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="추론 모니터링">
      <DomainHeader
        title="추론 모니터링"
        subtitle="AI 모델의 실시간 추론 성능과 데이터 드리프트를 모니터링합니다"
        icon={Activity}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'researcher' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8 space-y-6">
        {/* 전체 통계 카드 */}
        <section aria-label="전체 통계">
          <div className="grid grid-cols-4 gap-4">
            {OVERALL_STATS.map(stat => {
              const Icon = stat.icon;
              return (
                <div key={stat.label} className="rounded-2xl border border-slate-200 bg-white p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <div className={`h-8 w-8 rounded-lg flex items-center justify-center ${stat.color}`}>
                      <Icon className="h-4 w-4" />
                    </div>
                  </div>
                  <p className="text-2xl font-black text-slate-900">{stat.value}</p>
                  <p className="text-[11px] font-bold text-slate-400 mt-1">{stat.label}</p>
                </div>
              );
            })}
          </div>
        </section>

        {/* 모델별 모니터링 테이블 */}
        <section aria-label="모델별 추론 현황">
          <h2 className="text-sm font-bold text-slate-900 mb-3">모델별 추론 현황</h2>
          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b border-slate-100 bg-slate-50">
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">모델</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">요청/분</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">평균 지연</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">P95 지연</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">에러율</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">드리프트</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">지연 추세</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">액션</th>
                </tr>
              </thead>
              <tbody>
                {MODELS.map(model => {
                  const hasDrift = model.driftScore > 0.15;
                  const sparkMax = Math.max(...model.sparkline);
                  return (
                    <tr key={model.id} className="border-b border-slate-50 hover:bg-slate-50 transition-colors">
                      <td className="px-5 py-3">
                        <span className="text-xs font-bold text-slate-900">{model.name}</span>
                      </td>
                      <td className="px-5 py-3 text-xs font-mono text-slate-600">{model.requestsPerMin}</td>
                      <td className="px-5 py-3 text-xs font-mono text-slate-600">{model.avgLatency}ms</td>
                      <td className="px-5 py-3 text-xs font-mono text-slate-600">{model.p95Latency}ms</td>
                      <td className="px-5 py-3">
                        <span className={`text-xs font-bold ${model.errorRate > 0.1 ? 'text-red-600' : model.errorRate > 0.05 ? 'text-amber-600' : 'text-emerald-600'}`}>
                          {(model.errorRate * 100).toFixed(1)}%
                        </span>
                      </td>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-1">
                          <span className={`text-xs font-bold ${hasDrift ? 'text-red-600' : 'text-emerald-600'}`}>
                            {model.driftScore.toFixed(2)}
                          </span>
                          {hasDrift && (
                            <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded-full bg-red-50 text-red-600 text-[9px] font-bold border border-red-200">
                              <AlertTriangle className="h-2.5 w-2.5" /> 드리프트
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-5 py-3">
                        {/* 스파크라인 */}
                        <div className="flex items-end gap-px h-5 w-16" aria-label={`${model.name} 지연 추세`}>
                          {model.sparkline.map((val, i) => (
                            <div
                              key={i}
                              className={`flex-1 rounded-t-sm ${hasDrift ? 'bg-red-300' : 'bg-sky-300'}`}
                              style={{ height: `${(val / (sparkMax * 1.2)) * 100}%` }}
                              title={`${val}ms`}
                            />
                          ))}
                        </div>
                      </td>
                      <td className="px-5 py-3">
                        {hasDrift && (
                          <button
                            onClick={() => handleRetrain(model.id)}
                            className="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg bg-sky-600 text-white text-[10px] font-bold hover:bg-sky-700 transition-colors"
                            aria-label={`${model.name} 재학습 요청`}
                          >
                            {retrainingId === model.id ? (
                              <><RefreshCw className="h-3 w-3 animate-spin" /> 요청 중</>
                            ) : (
                              <><RefreshCw className="h-3 w-3" /> 재학습</>
                            )}
                          </button>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </section>

        {/* 드리프트 경고 배너 */}
        {MODELS.some(m => m.driftScore > 0.15) && (
          <section aria-label="드리프트 경고">
            <div className="rounded-2xl border border-amber-200 bg-amber-50 p-4" role="alert">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="h-4 w-4 text-amber-600" />
                <h3 className="text-xs font-bold text-amber-700">데이터 드리프트 감지</h3>
              </div>
              <p className="text-[11px] text-amber-600 ml-6">
                {MODELS.filter(m => m.driftScore > 0.15).map(m => m.name).join(', ')} 모델에서 드리프트 점수가
                임계치(0.15)를 초과했습니다. 모델 재학습을 검토해 주세요. [추정] 드리프트 임계값은 파라미터로 조정 가능합니다.
              </p>
            </div>
          </section>
        )}

        {/* 최근 오류 목록 */}
        <section aria-label="최근 오류">
          <h2 className="text-sm font-bold text-slate-900 mb-3">최근 오류</h2>
          <div className="space-y-2">
            {RECENT_ERRORS.map(err => (
              <div key={err.id} className="rounded-2xl border border-red-100 bg-white p-4">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                    <span className="text-xs font-bold text-slate-900">{err.model}</span>
                    <span className="px-2 py-0.5 rounded-full bg-red-50 text-red-600 text-[10px] font-bold border border-red-200">
                      {err.errorType}
                    </span>
                  </div>
                  <span className="text-[10px] font-mono text-slate-400">{err.timestamp}</span>
                </div>
                <p className="text-[11px] text-slate-500 ml-6">
                  입력 요약: <span className="font-mono text-slate-600">{err.inputSummary}</span>
                </p>
              </div>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
