'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Search, AlertTriangle, TrendingUp, BookOpen } from 'lucide-react';

type Disease = '당뇨' | '심혈관' | '고혈압';
type RiskLevel = '낮음' | '보통' | '높음' | '매우높음';
type TimePeriod = '7일' | '30일' | '90일';

interface FeatureContribution {
  name: string;
  value: number;
}

interface Recommendation {
  text: string;
  source: string;
  sourceType: 'guideline' | 'study' | 'expert';
}

interface DiseaseDetail {
  riskScore: number;
  riskLevel: RiskLevel;
  confidence: number;
  uncertainty: string;
  features: FeatureContribution[];
  projections: Record<TimePeriod, number[]>;
  recommendations: Recommendation[];
}

const RISK_LEVEL_STYLES: Record<RiskLevel, string> = {
  '낮음': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '보통': 'bg-amber-50 text-amber-700 border-amber-200',
  '높음': 'bg-orange-50 text-orange-700 border-orange-200',
  '매우높음': 'bg-red-50 text-red-700 border-red-200',
};

const DISEASE_DATA: Record<Disease, DiseaseDetail> = {
  '당뇨': {
    riskScore: 62,
    riskLevel: '높음',
    confidence: 94.2,
    uncertainty: '95% CI: 57~67%',
    features: [
      { name: 'BMI', value: 12 },
      { name: '공복혈당', value: 8 },
      { name: '나이', value: 5 },
      { name: '운동', value: -7 },
      { name: '식이', value: -3 },
    ],
    projections: {
      '7일': [62, 63, 61, 64, 62, 63, 65],
      '30일': [62, 61, 63, 60, 58, 57, 55, 54, 56, 55],
      '90일': [62, 60, 58, 55, 52, 50, 48, 46, 45, 44],
    },
    recommendations: [
      { text: '공복혈당 관리를 위해 식이섬유 섭취를 늘리세요', source: '대한당뇨학회 가이드라인', sourceType: 'guideline' },
      { text: '주 150분 이상의 중등도 유산소 운동을 권장합니다', source: 'ADA 2026 Standards', sourceType: 'guideline' },
      { text: 'HbA1c 수치를 3개월 주기로 모니터링하세요', source: '주치의 권고', sourceType: 'expert' },
    ],
  },
  '심혈관': {
    riskScore: 38,
    riskLevel: '보통',
    confidence: 91.8,
    uncertainty: '95% CI: 33~43%',
    features: [
      { name: '총콜레스테롤', value: 10 },
      { name: '혈압', value: 7 },
      { name: '흡연이력', value: 4 },
      { name: '운동', value: -9 },
      { name: 'HDL', value: -6 },
    ],
    projections: {
      '7일': [38, 37, 38, 39, 38, 37, 36],
      '30일': [38, 37, 36, 35, 34, 34, 33, 32, 33, 32],
      '90일': [38, 36, 34, 32, 30, 28, 27, 26, 25, 24],
    },
    recommendations: [
      { text: 'LDL 콜레스테롤을 100mg/dL 이하로 유지하세요', source: 'ESC 심혈관 예방 가이드라인', sourceType: 'guideline' },
      { text: '오메가-3 지방산이 풍부한 식단을 권장합니다', source: 'NEJM 2025 연구', sourceType: 'study' },
      { text: '혈압을 주 2회 이상 자가 모니터링하세요', source: '주치의 권고', sourceType: 'expert' },
    ],
  },
  '고혈압': {
    riskScore: 45,
    riskLevel: '보통',
    confidence: 93.1,
    uncertainty: '95% CI: 40~50%',
    features: [
      { name: '나트륨 섭취', value: 14 },
      { name: '체중', value: 9 },
      { name: '스트레스', value: 6 },
      { name: '칼륨 섭취', value: -5 },
      { name: '수면', value: -4 },
    ],
    projections: {
      '7일': [45, 44, 46, 45, 43, 44, 43],
      '30일': [45, 44, 43, 42, 41, 40, 39, 38, 38, 37],
      '90일': [45, 43, 41, 39, 37, 35, 33, 32, 31, 30],
    },
    recommendations: [
      { text: '일일 나트륨 섭취량을 2,000mg 이하로 제한하세요', source: '대한고혈압학회 가이드라인', sourceType: 'guideline' },
      { text: 'DASH 식단은 혈압 감소에 효과적입니다', source: 'Lancet 2024 메타분석', sourceType: 'study' },
      { text: '수면 시간을 7~8시간으로 확보하세요', source: '주치의 권고', sourceType: 'expert' },
    ],
  },
};

const DISEASES: Disease[] = ['당뇨', '심혈관', '고혈압'];
const TIME_PERIODS: TimePeriod[] = ['7일', '30일', '90일'];

const SOURCE_STYLES: Record<string, string> = {
  guideline: 'bg-sky-50 text-sky-700',
  study: 'bg-violet-50 text-violet-700',
  expert: 'bg-amber-50 text-amber-700',
};

const SOURCE_LABELS: Record<string, string> = {
  guideline: '가이드라인',
  study: '연구',
  expert: '전문가',
};

function CircularGauge({ score, size = 120 }: { score: number; size?: number }) {
  const radius = (size - 12) / 2;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;
  const color = score >= 75 ? '#ef4444' : score >= 50 ? '#f59e0b' : score >= 25 ? '#3b82f6' : '#10b981';

  return (
    <svg width={size} height={size} className="transform -rotate-90" aria-label={`위험도 ${score}점`}>
      <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke="#e2e8f0" strokeWidth="10" />
      <circle
        cx={size / 2} cy={size / 2} r={radius}
        fill="none" stroke={color} strokeWidth="10"
        strokeLinecap="round"
        strokeDasharray={circumference}
        strokeDashoffset={offset}
        className="transition-all duration-700"
      />
      <text
        x={size / 2} y={size / 2}
        textAnchor="middle" dominantBaseline="central"
        className="fill-slate-900 text-2xl font-black"
        transform={`rotate(90, ${size / 2}, ${size / 2})`}
      >
        {score}
      </text>
    </svg>
  );
}

function MiniProjectionChart({ data, color }: { data: number[]; color: string }) {
  const min = Math.min(...data) - 5;
  const max = Math.max(...data) + 5;
  const range = max - min || 1;
  const w = 280;
  const h = 80;
  const points = data.map((v, i) => `${(i / (data.length - 1)) * w},${h - ((v - min) / range) * h}`).join(' ');

  return (
    <svg width={w} height={h} className="w-full" viewBox={`0 0 ${w} ${h}`} preserveAspectRatio="none" aria-hidden="true">
      <polyline fill="none" stroke={color} strokeWidth="2" strokeLinejoin="round" points={points} />
      {data.map((v, i) => (
        <circle key={i} cx={(i / (data.length - 1)) * w} cy={h - ((v - min) / range) * h} r="3" fill={color} />
      ))}
    </svg>
  );
}

export default function DetailAnalysisPage() {
  const [selectedDisease, setSelectedDisease] = useState<Disease>('당뇨');
  const [timePeriod, setTimePeriod] = useState<TimePeriod>('30일');

  const detail = DISEASE_DATA[selectedDisease];
  const maxAbsValue = Math.max(...detail.features.map(f => Math.abs(f.value)));
  const gaugeColor = detail.riskScore >= 75 ? '#ef4444' : detail.riskScore >= 50 ? '#f59e0b' : '#3b82f6';

  return (
    <main className="min-h-screen bg-slate-50" aria-label="질환별 상세 분석">
      <DomainHeader title="질환별 상세 분석" subtitle="각 질환의 위험 요인을 심층 분석합니다" icon={Search} />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 질환 선택 탭 */}
        <div className="flex gap-2" role="tablist" aria-label="질환 선택">
          {DISEASES.map((disease) => (
            <button
              key={disease}
              role="tab"
              aria-selected={selectedDisease === disease}
              onClick={() => setSelectedDisease(disease)}
              className={`px-5 py-2.5 rounded-xl text-sm font-bold transition-all ${
                selectedDisease === disease
                  ? 'bg-slate-900 text-white shadow-sm'
                  : 'bg-white text-slate-600 border border-slate-200 hover:border-slate-300'
              }`}
            >
              {disease}
            </button>
          ))}
        </div>

        {/* 위험도 & 기여도 */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* 위험도 게이지 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-sm font-bold text-slate-500 mb-4">위험도 점수</h2>
            <div className="flex items-center gap-6">
              <CircularGauge score={detail.riskScore} size={120} />
              <div className="space-y-2">
                <span className={`inline-block px-3 py-1 rounded-full text-xs font-semibold border ${RISK_LEVEL_STYLES[detail.riskLevel]}`}>
                  {detail.riskLevel}
                </span>
                <p className="text-xs text-slate-500">모델 신뢰도: {detail.confidence}%</p>
                <p className="text-xs text-slate-400">{detail.uncertainty}</p>
              </div>
            </div>
          </div>

          {/* SHAP 기여도 차트 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-6">
            <h2 className="text-sm font-bold text-slate-500 mb-4">주요 기여 요인 (SHAP)</h2>
            <div className="space-y-3">
              {detail.features.map((f) => {
                const widthPct = (Math.abs(f.value) / maxAbsValue) * 100;
                const isPositive = f.value > 0;
                return (
                  <div key={f.name} className="flex items-center gap-3">
                    <span className="w-16 text-xs font-medium text-slate-700 text-right shrink-0">{f.name}</span>
                    <div className="flex-1 flex items-center h-6 relative">
                      <div className="absolute inset-0 flex items-center">
                        <div className="w-full h-px bg-slate-200" />
                      </div>
                      <div
                        className={`relative h-5 rounded transition-all ${isPositive ? 'bg-red-400' : 'bg-emerald-400'}`}
                        style={{
                          width: `${widthPct * 0.5}%`,
                          marginLeft: isPositive ? '50%' : `${50 - widthPct * 0.5}%`,
                        }}
                      />
                    </div>
                    <span className={`w-10 text-xs font-bold text-right ${isPositive ? 'text-red-600' : 'text-emerald-600'}`}>
                      {isPositive ? '+' : ''}{f.value}
                    </span>
                  </div>
                );
              })}
            </div>
            <p className="text-[10px] text-slate-400 mt-3">빨간색: 위험 증가 / 녹색: 위험 감소</p>
          </div>
        </div>

        {/* 시간 투영 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-bold text-slate-500">위험도 추이 예측</h2>
            <div className="flex gap-1">
              {TIME_PERIODS.map((tp) => (
                <button
                  key={tp}
                  onClick={() => setTimePeriod(tp)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-bold transition-all ${
                    timePeriod === tp ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'
                  }`}
                >
                  {tp}
                </button>
              ))}
            </div>
          </div>
          <div className="h-24">
            <MiniProjectionChart data={detail.projections[timePeriod]} color={gaugeColor} />
          </div>
          <p className="text-[10px] text-slate-400 mt-2">
            [추정] 생활습관 유지 시 예측 추이. 실제 결과는 다를 수 있습니다.
          </p>
        </div>

        {/* 근거 기반 권고 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-500 mb-4 flex items-center gap-2">
            <BookOpen className="h-4 w-4" />
            근거 기반 권고사항
          </h2>
          <div className="space-y-3">
            {detail.recommendations.map((rec, idx) => (
              <div key={idx} className="flex items-start gap-3 p-4 rounded-xl bg-slate-50">
                <div className="h-6 w-6 rounded-full bg-sky-100 text-sky-600 flex items-center justify-center text-xs font-bold shrink-0 mt-0.5">
                  {idx + 1}
                </div>
                <div className="flex-1">
                  <p className="text-sm text-slate-800 font-medium">{rec.text}</p>
                  <span className={`inline-block mt-1.5 px-2 py-0.5 rounded text-[10px] font-semibold ${SOURCE_STYLES[rec.sourceType]}`}>
                    {SOURCE_LABELS[rec.sourceType]}: {rec.source}
                  </span>
                </div>
              </div>
            ))}
          </div>
          <p className="text-[10px] text-slate-400 mt-4">
            <AlertTriangle className="inline h-3 w-3 mr-1" />
            본 권고는 AI 분석 결과이며, 의료 행위를 대체하지 않습니다. 반드시 주치의와 상담하세요.
          </p>
        </div>
      </div>
    </main>
  );
}
