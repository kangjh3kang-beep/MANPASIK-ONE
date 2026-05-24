'use client';

import React, { useState, useMemo } from 'react';
import { DomainHeader } from '@mmup/ui';
import { SlidersHorizontal, RotateCcw, AlertTriangle, ArrowUp, ArrowDown, Minus } from 'lucide-react';

interface SliderConfig {
  id: string;
  label: string;
  unit: string;
  min: number;
  max: number;
  step: number;
  defaultValue: number;
}

interface DiseaseRisk {
  name: string;
  baseRisk: number;
  weights: Record<string, number>;
}

const SLIDERS: SliderConfig[] = [
  { id: 'bmi', label: '체질량지수(BMI)', unit: 'kg/m²', min: 15, max: 40, step: 0.1, defaultValue: 24.5 },
  { id: 'glucose', label: '공복혈당', unit: 'mg/dL', min: 70, max: 200, step: 1, defaultValue: 108 },
  { id: 'cholesterol', label: '총콜레스테롤', unit: 'mg/dL', min: 100, max: 300, step: 1, defaultValue: 195 },
  { id: 'exercise', label: '주간운동시간', unit: '시간', min: 0, max: 14, step: 0.5, defaultValue: 3 },
  { id: 'sleep', label: '수면시간', unit: '시간', min: 4, max: 10, step: 0.5, defaultValue: 6.5 },
];

const DISEASES: DiseaseRisk[] = [
  {
    name: '당뇨',
    baseRisk: 62,
    weights: { bmi: 1.8, glucose: 0.35, cholesterol: 0.05, exercise: -2.5, sleep: -1.2 },
  },
  {
    name: '심혈관',
    baseRisk: 38,
    weights: { bmi: 1.2, glucose: 0.15, cholesterol: 0.18, exercise: -3.0, sleep: -0.8 },
  },
  {
    name: '고혈압',
    baseRisk: 45,
    weights: { bmi: 1.5, glucose: 0.1, cholesterol: 0.12, exercise: -2.0, sleep: -1.5 },
  },
];

function computeRisk(disease: DiseaseRisk, values: Record<string, number>, defaults: Record<string, number>): number {
  let risk = disease.baseRisk;
  for (const [key, weight] of Object.entries(disease.weights)) {
    const delta = values[key] - defaults[key];
    risk += delta * weight;
  }
  return Math.max(0, Math.min(100, Math.round(risk)));
}

function ChangeIndicator({ before, after }: { before: number; after: number }) {
  const diff = after - before;
  if (diff === 0) return <span className="flex items-center gap-1 text-xs text-slate-400"><Minus className="h-3 w-3" />변동 없음</span>;
  const isUp = diff > 0;
  return (
    <span className={`flex items-center gap-1 text-xs font-bold ${isUp ? 'text-red-600' : 'text-emerald-600'}`}>
      {isUp ? <ArrowUp className="h-3 w-3" /> : <ArrowDown className="h-3 w-3" />}
      {isUp ? '+' : ''}{diff}%p
    </span>
  );
}

export default function SimulatorPage() {
  const defaultValues: Record<string, number> = {};
  SLIDERS.forEach(s => { defaultValues[s.id] = s.defaultValue; });

  const [values, setValues] = useState<Record<string, number>>({ ...defaultValues });

  const handleSliderChange = (id: string, value: number) => {
    setValues(prev => ({ ...prev, [id]: value }));
  };

  const handleReset = () => {
    setValues({ ...defaultValues });
  };

  const results = useMemo(() => {
    return DISEASES.map(disease => ({
      name: disease.name,
      before: disease.baseRisk,
      after: computeRisk(disease, values, defaultValues),
    }));
  }, [values]);

  const isModified = Object.keys(values).some(k => values[k] !== defaultValues[k]);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="건강 시뮬레이터">
      <DomainHeader
        title="건강 시뮬레이터"
        subtitle="생활습관 변화에 따른 위험도 변화를 시뮬레이션합니다"
        icon={SlidersHorizontal}
      />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 안내 문구 */}
        <div className="rounded-2xl border border-sky-200 bg-sky-50 p-4">
          <p className="text-sm text-sky-800 font-medium">
            생활습관을 조절하면 위험도가 어떻게 변하는지 확인해보세요
          </p>
          <p className="text-xs text-sky-600 mt-1">
            슬라이더를 움직여 각 수치를 변경하면 실시간으로 예측 위험도가 업데이트됩니다.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
          {/* 입력 슬라이더 */}
          <div className="lg:col-span-3 space-y-4">
            <div className="rounded-2xl border border-slate-200 bg-white p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-sm font-bold text-slate-500">생활습관 파라미터</h2>
                <button
                  onClick={handleReset}
                  disabled={!isModified}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                    isModified
                      ? 'bg-slate-100 text-slate-700 hover:bg-slate-200'
                      : 'bg-slate-50 text-slate-300 cursor-not-allowed'
                  }`}
                >
                  <RotateCcw className="h-3 w-3" />
                  초기화
                </button>
              </div>

              <div className="space-y-6">
                {SLIDERS.map((slider) => {
                  const current = values[slider.id];
                  const isChanged = current !== slider.defaultValue;
                  return (
                    <div key={slider.id}>
                      <div className="flex items-center justify-between mb-2">
                        <label htmlFor={slider.id} className="text-sm font-medium text-slate-700">
                          {slider.label}
                        </label>
                        <div className="flex items-center gap-2">
                          {isChanged && (
                            <span className="text-[10px] text-slate-400 line-through">
                              {slider.defaultValue}{slider.unit}
                            </span>
                          )}
                          <span className={`text-sm font-bold tabular-nums ${isChanged ? 'text-sky-600' : 'text-slate-900'}`}>
                            {current}{slider.unit}
                          </span>
                        </div>
                      </div>
                      <input
                        id={slider.id}
                        type="range"
                        min={slider.min}
                        max={slider.max}
                        step={slider.step}
                        value={current}
                        onChange={(e) => handleSliderChange(slider.id, parseFloat(e.target.value))}
                        className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-sky-600"
                        aria-valuemin={slider.min}
                        aria-valuemax={slider.max}
                        aria-valuenow={current}
                        aria-label={`${slider.label}: ${current} ${slider.unit}`}
                      />
                      <div className="flex justify-between mt-1">
                        <span className="text-[10px] text-slate-400">{slider.min}</span>
                        <span className="text-[10px] text-slate-400">{slider.max}</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>

          {/* 결과 패널 */}
          <div className="lg:col-span-2 space-y-4">
            <div className="rounded-2xl border border-slate-200 bg-white p-6 sticky top-4">
              <h2 className="text-sm font-bold text-slate-500 mb-4">시뮬레이션 결과</h2>

              <div className="space-y-4">
                {results.map((result) => {
                  const riskColor =
                    result.after >= 75 ? 'text-red-600' :
                    result.after >= 50 ? 'text-amber-600' :
                    result.after >= 25 ? 'text-sky-600' : 'text-emerald-600';
                  const barColor =
                    result.after >= 75 ? 'bg-red-400' :
                    result.after >= 50 ? 'bg-amber-400' :
                    result.after >= 25 ? 'bg-sky-400' : 'bg-emerald-400';

                  return (
                    <div key={result.name} className="rounded-xl border border-slate-100 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-bold text-slate-800">{result.name}</span>
                        <ChangeIndicator before={result.before} after={result.after} />
                      </div>

                      <div className="flex items-center gap-3 mb-2">
                        <div className="flex-1">
                          <div className="flex items-baseline gap-1">
                            <span className="text-xs text-slate-400">현재</span>
                            <span className="text-sm font-bold text-slate-500">{result.before}%</span>
                          </div>
                        </div>
                        <span className="text-slate-300">→</span>
                        <div className="flex-1 text-right">
                          <div className="flex items-baseline gap-1 justify-end">
                            <span className="text-xs text-slate-400">변경 후</span>
                            <span className={`text-lg font-black tabular-nums ${riskColor}`}>{result.after}%</span>
                          </div>
                        </div>
                      </div>

                      {/* 진행 바 */}
                      <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full transition-all duration-500 ${barColor}`}
                          style={{ width: `${result.after}%` }}
                        />
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* 종합 변화 */}
              {isModified && (
                <div className="mt-4 p-3 rounded-xl bg-slate-50">
                  <p className="text-xs text-slate-500 text-center">
                    {results.every(r => r.after <= r.before)
                      ? '모든 질환의 위험도가 감소하였습니다'
                      : results.every(r => r.after >= r.before)
                        ? '모든 질환의 위험도가 증가하였습니다'
                        : '일부 질환의 위험도가 변동되었습니다'}
                  </p>
                </div>
              )}
            </div>

            {/* 면책 조항 */}
            <div className="rounded-xl border border-amber-200 bg-amber-50 p-4">
              <div className="flex items-start gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-600 shrink-0 mt-0.5" />
                <div>
                  <p className="text-[11px] text-amber-800 font-semibold mb-1">시뮬레이션 안내</p>
                  <p className="text-[10px] text-amber-700 leading-relaxed">
                    본 시뮬레이션은 통계 모델 기반 [추정] 결과이며, 실제 의학적 진단을 대체하지 않습니다.
                    정확한 건강 평가는 반드시 의료 전문가와 상담하세요.
                    모델 신뢰도: 92~96% [검증됨, 임상 미검증].
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
