'use client';

import React, { useState, useCallback } from 'react';
import { DomainHeader, MeasurementResultCard, NextActionGroup } from '@mmup/ui';
import { Activity, CheckCircle, Clock, Loader2, AlertTriangle, Smartphone, ClipboardList, Brain, BarChart3 } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type MeasureStep = 'ready' | 'cartridge' | 'sample' | 'measuring' | 'result';

const STEPS = [
  { id: 'cartridge', label: '카트리지 인식', icon: Smartphone },
  { id: 'sample', label: '시료 준비', icon: Clock },
  { id: 'measuring', label: '측정 진행', icon: Loader2 },
  { id: 'result', label: '결과 확인', icon: CheckCircle },
];

// 데모 측정 시뮬레이션
function simulateMeasurement(): Promise<{
  biomarker: string; value: number; unit: string; confidence: number;
  uncertainty: { value: number; ci_lower: number; ci_upper: number };
  referenceRange: { low: number; high: number };
  status: 'normal' | 'caution' | 'danger';
}> {
  return new Promise(resolve => {
    setTimeout(() => {
      const value = 85 + Math.random() * 40;
      const status = value < 140 ? 'normal' : value < 180 ? 'caution' : 'danger';
      resolve({
        biomarker: '혈당',
        value: Math.round(value),
        unit: 'mg/dL',
        confidence: 0.93 + Math.random() * 0.05,
        uncertainty: { value: 2.8, ci_lower: Math.round(value - 2.8), ci_upper: Math.round(value + 2.8) },
        referenceRange: { low: 70, high: 140 },
        status,
      });
    }, 3000);
  });
}

export default function MeasurePage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [step, setStep] = useState<MeasureStep>('ready');
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState<any>(null);
  const [cartridgeInfo, setCartridgeInfo] = useState<{ type: string; remaining: number } | null>(null);

  const startMeasurement = useCallback(async () => {
    // Step 1: 카트리지 인식
    setStep('cartridge');
    setProgress(0);
    await new Promise(r => setTimeout(r, 1500));
    setCartridgeInfo({ type: '혈당 카트리지', remaining: 8 });
    setProgress(25);

    // Step 2: 시료 준비
    setStep('sample');
    await new Promise(r => setTimeout(r, 2000));
    setProgress(50);

    // Step 3: 측정 진행
    setStep('measuring');
    const interval = setInterval(() => {
      setProgress(prev => Math.min(prev + 2, 90));
    }, 100);

    const measureResult = await simulateMeasurement();
    clearInterval(interval);
    setProgress(100);

    // Step 4: 결과
    setResult(measureResult);
    setStep('result');
  }, []);

  const resetMeasurement = () => {
    setStep('ready');
    setProgress(0);
    setResult(null);
    setCartridgeInfo(null);
  };

  const currentStepIndex = STEPS.findIndex(s => s.id === step);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="측정 진행">
      <DomainHeader
        title="건강 측정"
        subtitle="카트리지를 인식하고 측정을 시작합니다"
        icon={Activity}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-2xl mx-auto px-6 py-8">
        {/* 단계 표시 (스테퍼) */}
        {step !== 'ready' && (
          <div className="flex items-center justify-between mb-8">
            {STEPS.map((s, i) => {
              const isActive = i === currentStepIndex;
              const isDone = i < currentStepIndex;
              return (
                <React.Fragment key={s.id}>
                  <div className="flex flex-col items-center gap-1">
                    <div className={`h-10 w-10 rounded-full flex items-center justify-center text-sm font-bold transition-all ${
                      isDone ? 'bg-emerald-500 text-white' :
                      isActive ? 'bg-sky-600 text-white ring-4 ring-sky-100' :
                      'bg-slate-200 text-slate-400'
                    }`}>
                      {isDone ? <CheckCircle className="h-5 w-5" /> : <s.icon className={`h-5 w-5 ${isActive && s.id === 'measuring' ? 'animate-spin' : ''}`} />}
                    </div>
                    <span className={`text-[11px] font-medium ${isActive ? 'text-sky-600' : isDone ? 'text-emerald-600' : 'text-slate-400'}`}>
                      {s.label}
                    </span>
                  </div>
                  {i < STEPS.length - 1 && (
                    <div className={`flex-1 h-0.5 mx-2 rounded ${isDone ? 'bg-emerald-300' : 'bg-slate-200'}`} />
                  )}
                </React.Fragment>
              );
            })}
          </div>
        )}

        {/* 진행률 바 */}
        {step !== 'ready' && step !== 'result' && (
          <div className="mb-8">
            <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
              <div className="h-full bg-sky-500 rounded-full transition-all duration-300" style={{ width: `${progress}%` }} />
            </div>
            <p className="text-xs text-slate-400 mt-1 text-right">{progress}%</p>
          </div>
        )}

        {/* 단계별 콘텐츠 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">

          {/* 준비 상태 — 측정 시작 버튼 */}
          {step === 'ready' && (
            <div className="text-center py-8">
              <div className="h-20 w-20 rounded-full bg-sky-50 flex items-center justify-center mx-auto mb-6">
                <Activity className="h-10 w-10 text-sky-600" />
              </div>
              <h2 className="text-xl font-bold text-slate-900 mb-2">측정을 시작하세요</h2>
              <p className="text-sm text-slate-500 mb-8">
                카트리지를 리더기에 삽입한 후 아래 버튼을 누르세요.<br />
                오프라인에서도 측정이 가능합니다.
              </p>
              <button
                onClick={startMeasurement}
                className="px-8 py-4 bg-sky-600 hover:bg-sky-700 text-white font-bold text-lg rounded-2xl transition-all shadow-lg shadow-sky-600/20"
              >
                측정 시작
              </button>
            </div>
          )}

          {/* 카트리지 인식 */}
          {step === 'cartridge' && (
            <div className="text-center py-8">
              <Loader2 className="h-12 w-12 text-sky-500 animate-spin mx-auto mb-4" />
              <h2 className="text-lg font-bold text-slate-900 mb-2">카트리지를 인식하고 있습니다</h2>
              <p className="text-sm text-slate-500">NFC로 카트리지 정보를 읽는 중...</p>
            </div>
          )}

          {/* 시료 준비 가이드 */}
          {step === 'sample' && cartridgeInfo && (
            <div className="py-4">
              <div className="flex items-center gap-3 mb-6 p-4 rounded-xl bg-emerald-50 border border-emerald-200">
                <CheckCircle className="h-5 w-5 text-emerald-600 flex-shrink-0" />
                <div>
                  <p className="text-sm font-bold text-emerald-700">{cartridgeInfo.type} 인식 완료</p>
                  <p className="text-xs text-emerald-600">남은 사용 횟수: {cartridgeInfo.remaining}회</p>
                </div>
              </div>

              <h2 className="text-lg font-bold text-slate-900 mb-4">시료를 준비해 주세요</h2>
              <ol className="space-y-4 mb-6">
                <li className="flex items-start gap-3">
                  <span className="h-7 w-7 rounded-full bg-sky-100 text-sky-700 flex items-center justify-center text-sm font-bold flex-shrink-0">1</span>
                  <div>
                    <p className="text-sm font-medium text-slate-700">알코올 솜으로 손가락을 소독하세요</p>
                    <p className="text-xs text-slate-400">깨끗한 표면에서 정확한 측정이 가능합니다</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="h-7 w-7 rounded-full bg-sky-100 text-sky-700 flex items-center justify-center text-sm font-bold flex-shrink-0">2</span>
                  <div>
                    <p className="text-sm font-medium text-slate-700">채혈침으로 손끝을 살짝 찌르세요</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <span className="h-7 w-7 rounded-full bg-sky-100 text-sky-700 flex items-center justify-center text-sm font-bold flex-shrink-0">3</span>
                  <div>
                    <p className="text-sm font-medium text-slate-700">카트리지 투입구에 혈액 한 방울을 떨어뜨리세요</p>
                  </div>
                </li>
              </ol>
              <p className="text-xs text-slate-400 text-center">시료 준비가 자동으로 감지되면 측정이 시작됩니다</p>
            </div>
          )}

          {/* 측정 진행 */}
          {step === 'measuring' && (
            <div className="text-center py-8">
              <div className="relative inline-flex items-center justify-center mb-6">
                <svg className="h-24 w-24 -rotate-90">
                  <circle className="text-slate-100" strokeWidth="6" stroke="currentColor" fill="transparent" r="44" cx="48" cy="48" />
                  <circle
                    className="text-sky-500 transition-all duration-300"
                    strokeWidth="6"
                    strokeDasharray={2 * Math.PI * 44}
                    strokeDashoffset={2 * Math.PI * 44 * (1 - progress / 100)}
                    strokeLinecap="round"
                    stroke="currentColor"
                    fill="transparent"
                    r="44" cx="48" cy="48"
                  />
                </svg>
                <span className="absolute text-lg font-bold text-slate-700">{progress}%</span>
              </div>
              <h2 className="text-lg font-bold text-slate-900 mb-2">측정 중입니다</h2>
              <p className="text-sm text-slate-500">잠시만 기다려 주세요. 약 30초 소요됩니다.</p>
              <p className="text-xs text-slate-400 mt-4">차동 보정 (Sdiff = S - α×R) 적용 중...</p>
            </div>
          )}

          {/* 결과 카드 */}
          {step === 'result' && result && (
            <div>
              <div className="flex items-center gap-2 mb-6">
                <CheckCircle className="h-5 w-5 text-emerald-500" />
                <h2 className="text-lg font-bold text-slate-900">측정 완료</h2>
              </div>

              <MeasurementResultCard
                biomarker={result.biomarker}
                value={result.value}
                unit={result.unit}
                confidence={result.confidence}
                uncertainty={result.uncertainty}
                referenceRange={result.referenceRange}
                status={result.status}
                timestamp={new Date().toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' })}
              />

              <div className="flex gap-3 mt-6">
                <button
                  onClick={resetMeasurement}
                  className="flex-1 px-4 py-3 bg-sky-600 hover:bg-sky-700 text-white font-bold rounded-xl transition-all"
                >
                  다시 측정하기
                </button>
                <Link
                  href={`${localePrefix}/domains/clinical`}
                  className="flex-1 px-4 py-3 border border-slate-200 text-slate-700 font-bold rounded-xl text-center hover:bg-slate-50 transition-all"
                >
                  기록 보기
                </Link>
              </div>

              {/* 다음 단계 추천 */}
              <div className="mt-6">
                <NextActionGroup
                  title="다음 단계"
                  actions={[
                    { id: 'records', title: '기록에서 추이 보기', description: '과거 측정값과 비교하여 변화를 확인합니다', icon: <ClipboardList className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/health-records`, priority: 'high' },
                    { id: 'coach', title: 'AI 코치에게 해석 요청', description: '측정 결과를 AI가 쉽게 설명해드립니다', icon: <Brain className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/ai-coach` },
                    { id: 'predictor', title: '위험 예측 확인', description: '건강 위험도를 미리 분석합니다', icon: <BarChart3 className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/predictor` },
                  ]}
                />
              </div>
            </div>
          )}
        </div>

        {/* 안내 문구 */}
        <div className="mt-6 flex items-start gap-2 text-xs text-slate-400">
          <AlertTriangle className="h-4 w-4 flex-shrink-0 mt-0.5" />
          <p>이 측정 결과는 AI 참고 정보이며, 의학적 진단을 대체하지 않습니다. 정확한 진단은 전문의와 상담하세요.</p>
        </div>

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-8">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: '건강 기록', path: '/domains/health-records', desc: '측정 이력 확인' },
              { name: 'AI 코치', path: '/domains/ai-coach', desc: '결과 해석 상담' },
              { name: '건강 데이터', path: '/domains/clinical', desc: '상세 분석' },
            ].map(d => (
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
      </div>
    </main>
  );
}
