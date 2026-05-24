'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { Brain, TrendingUp, TrendingDown, Minus, ArrowRight, Share2, AlertTriangle } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type TrendDirection = 'up' | 'down' | 'stable';

interface Insight {
  id: string;
  title: string;
  description: string;
  trend: TrendDirection;
  trendLabel: string;
  confidence: number;
  action: string;
  severity: 'positive' | 'warning' | 'neutral' | 'alert';
}

const OVERALL_TREND: { direction: TrendDirection; label: string } = {
  direction: 'up',
  label: '개선',
};

const MOCK_INSIGHTS: Insight[] = [
  {
    id: 'i1',
    title: '혈당 변동 패턴 분석',
    description:
      '지난 30일간 공복 혈당의 표준편차가 18mg/dL에서 12mg/dL로 감소하였습니다. ' +
      '식후 2시간 혈당 스파이크가 주 3회에서 1회로 줄어들었으며, 특히 저녁 식사 후 혈당 안정세가 두드러집니다. ' +
      '현재 식이 패턴과 운동 루틴의 효과가 나타나고 있는 것으로 판단됩니다.',
    trend: 'up',
    trendLabel: '변동폭 33% 감소',
    confidence: 92,
    action: '현재 식이 패턴을 유지하고, 저녁 식후 산책을 계속하세요',
    severity: 'positive',
  },
  {
    id: 'i2',
    title: '콜레스테롤 추이 예측',
    description:
      'LDL 콜레스테롤이 3개월 연속 상승 추세를 보이고 있습니다 (128 → 135 → 145mg/dL). ' +
      '현재 추세가 유지될 경우 8주 이내에 160mg/dL을 초과할 수 있으며, 약물 치료 기준(160mg/dL)에 근접합니다. ' +
      'HDL 수치는 52mg/dL로 정상 범위를 유지하고 있습니다.',
    trend: 'down',
    trendLabel: 'LDL 3개월 연속 상승',
    confidence: 87,
    action: '포화지방 섭취를 줄이고, 주 3회 이상 유산소 운동을 권장합니다',
    severity: 'warning',
  },
  {
    id: 'i3',
    title: '측정 빈도 분석',
    description:
      '최근 4주간 혈당 측정은 주 5.2회로 권장 빈도(주 7회)에 미치지 못하나, 이전 대비 개선되었습니다. ' +
      '콜레스테롤 측정은 월 1회로 적정 수준을 유지하고 있습니다. ' +
      '주말 측정 누락이 잦은 패턴이 관찰되며, 리마인더 설정을 권장합니다.',
    trend: 'stable',
    trendLabel: '주 5.2회 (권장 7회)',
    confidence: 95,
    action: '주말 아침 측정 리마인더를 설정하여 규칙적 측정을 유지하세요',
    severity: 'neutral',
  },
  {
    id: 'i4',
    title: '복합 위험도 변화',
    description:
      '혈당, 혈압, 콜레스테롤 수치를 종합한 심혈관 위험도 지수가 지난 분기 대비 8% 감소하였습니다. ' +
      '혈당 관리 개선이 주요 기여 요인이나, LDL 콜레스테롤 상승이 부분적 상쇄 요인으로 작용하고 있습니다. ' +
      '전체적인 건강 관리 방향은 긍정적입니다.',
    trend: 'up',
    trendLabel: '위험도 8% 감소',
    confidence: 84,
    action: '콜레스테롤 관리에 집중하면 복합 위험도를 추가 10% 이상 낮출 수 있습니다',
    severity: 'positive',
  },
];

const trendConfig = {
  up: { icon: TrendingUp, color: 'text-emerald-600' },
  down: { icon: TrendingDown, color: 'text-red-600' },
  stable: { icon: Minus, color: 'text-slate-500' },
};

const severityConfig = {
  positive: { bg: 'bg-emerald-50', border: 'border-emerald-200', accent: 'text-emerald-700' },
  warning: { bg: 'bg-amber-50', border: 'border-amber-200', accent: 'text-amber-700' },
  neutral: { bg: 'bg-sky-50', border: 'border-sky-200', accent: 'text-sky-700' },
  alert: { bg: 'bg-red-50', border: 'border-red-200', accent: 'text-red-700' },
};

const severityLabels = {
  positive: '개선',
  warning: '주의',
  neutral: '유지',
  alert: '경고',
};

export default function InsightsPage() {
  const localePrefix = useLocalePrefix();
  const OverallIcon = trendConfig[OVERALL_TREND.direction].icon;

  return (
    <main className="min-h-screen bg-slate-50" aria-label="AI 트렌드 분석">
      <DomainHeader
        title="AI 트렌드 분석"
        subtitle="측정 데이터 기반 AI 건강 트렌드 인사이트입니다"
        icon={Brain}
      />

      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 전체 건강 추세 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6 text-center">
          <p className="text-xs font-semibold text-slate-500 mb-2">종합 건강 추세</p>
          <div className="flex items-center justify-center gap-3">
            <div className="h-12 w-12 rounded-2xl bg-emerald-50 flex items-center justify-center">
              <OverallIcon className="h-6 w-6 text-emerald-600" />
            </div>
            <span className="text-2xl font-black text-emerald-600">{OVERALL_TREND.label}</span>
          </div>
          <p className="text-xs text-slate-400 mt-2">최근 30일 기준 | 혈당, 혈압, 콜레스테롤, 심박수 종합 분석</p>
        </div>

        {/* 인사이트 카드 */}
        {MOCK_INSIGHTS.map((insight) => {
          const trend = trendConfig[insight.trend];
          const severity = severityConfig[insight.severity];
          const TrendIcon = trend.icon;

          return (
            <div key={insight.id} className="rounded-2xl border border-slate-200 bg-white p-6">
              <div className="flex items-start gap-4">
                <div className={`h-10 w-10 rounded-xl ${severity.bg} flex items-center justify-center flex-shrink-0`}>
                  <TrendIcon className={`h-5 w-5 ${trend.color}`} />
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold ${severity.bg} ${severity.accent}`}>
                      {severityLabels[insight.severity]}
                    </span>
                    <span className="text-[10px] text-slate-400">신뢰도 {insight.confidence}%</span>
                  </div>
                  <h3 className="text-sm font-bold text-slate-900 mb-1">{insight.title}</h3>
                  <p className="text-xs text-slate-500 leading-relaxed mb-3">{insight.description}</p>

                  <div className="bg-slate-50 rounded-xl p-3 mb-3">
                    <p className="text-xs font-bold text-slate-700 mb-1">추천 조치</p>
                    <div className="flex items-start gap-2 text-xs text-slate-600">
                      <ArrowRight className="h-3 w-3 text-sky-500 mt-0.5 flex-shrink-0" />
                      <span>{insight.action}</span>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <span className={`text-xs font-bold ${trend.color}`}>{insight.trendLabel}</span>
                    <button className="text-xs font-semibold text-sky-600 hover:underline inline-flex items-center gap-1">
                      <Share2 className="h-3 w-3" />
                      전문의와 공유
                    </button>
                  </div>
                </div>
              </div>
            </div>
          );
        })}

        {/* 상세 위험도 분석 링크 */}
        <div className="rounded-2xl border border-sky-100 bg-sky-50 p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-bold text-sky-900">더 자세한 위험도 분석이 필요하신가요?</p>
              <p className="text-xs text-sky-600 mt-0.5">예측 분석 도메인에서 상세 위험도 시뮬레이션을 확인하세요</p>
            </div>
            <Link
              href={`${localePrefix}/domains/predictor`}
              className="px-4 py-2 rounded-xl bg-sky-600 text-white text-xs font-semibold hover:bg-sky-700 transition-colors inline-flex items-center gap-1.5"
            >
              위험도 분석
              <ArrowRight className="h-3 w-3" />
            </Link>
          </div>
        </div>

        {/* AI 면책 조항 */}
        <div className="rounded-2xl border border-slate-100 bg-slate-50 p-4">
          <div className="flex items-start gap-2">
            <AlertTriangle className="h-4 w-4 text-slate-400 flex-shrink-0 mt-0.5" />
            <p className="text-xs text-slate-500 leading-relaxed">
              AI 트렌드 분석은 참고 정보이며, 의학적 진단이나 처방을 대체하지 않습니다.
              모든 분석 결과에는 신뢰도가 표시되며, 중요한 건강 결정은 반드시 전문 의료진과 상담하시기 바랍니다.
              분석에 사용된 데이터는 사용자의 측정 기록에 기반합니다.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
