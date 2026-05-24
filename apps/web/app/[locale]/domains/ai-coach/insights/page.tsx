'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { Lightbulb, TrendingUp, TrendingDown, AlertTriangle, CheckCircle, ArrowRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const WEEKLY_INSIGHTS = [
  {
    id: '1',
    type: 'positive' as const,
    title: '혈당 관리가 개선되고 있습니다',
    description: '지난 7일간 공복 혈당 평균이 115mg/dL에서 108mg/dL로 6% 감소했습니다. 식후 혈당 변동폭도 안정세를 보이고 있습니다.',
    metric: '공복혈당 -6%',
    confidence: 0.91,
    actions: ['식후 30분 가벼운 산책을 계속 유지하세요', '현재 식단 패턴을 유지하세요'],
  },
  {
    id: '2',
    type: 'warning' as const,
    title: '수면 패턴에 변화가 감지되었습니다',
    description: '이번 주 평균 수면 시간이 6.2시간으로 지난 주 대비 45분 감소했습니다. 특히 목~금요일 취침 시간이 불규칙합니다.',
    metric: '수면시간 -45분',
    confidence: 0.87,
    actions: ['취침 1시간 전 스마트폰 사용을 줄여보세요', '주중 취침 시간을 일정하게 유지하세요'],
  },
  {
    id: '3',
    type: 'neutral' as const,
    title: '활동량이 목표에 근접합니다',
    description: '일평균 걸음 수 7,823보로 목표(8,000보) 대비 97.8% 달성 중입니다. 주말 활동량이 평일 대비 40% 높습니다.',
    metric: '목표 달성률 97.8%',
    confidence: 0.94,
    actions: ['평일 점심시간 10분 산책을 추가해보세요'],
  },
  {
    id: '4',
    type: 'alert' as const,
    title: '콜레스테롤 수치 모니터링 필요',
    description: '최근 측정된 LDL 콜레스테롤이 145mg/dL로 이전 대비 12% 상승했습니다. 3개월 연속 상승 추세를 보이고 있어 주의가 필요합니다.',
    metric: 'LDL +12%',
    confidence: 0.89,
    actions: ['화상 진료를 통해 전문의와 상담하세요', '포화지방 섭취를 줄이세요'],
  },
];

const typeConfig = {
  positive: { color: 'emerald', icon: TrendingUp, label: '개선' },
  warning: { color: 'amber', icon: TrendingDown, label: '주의' },
  neutral: { color: 'sky', icon: CheckCircle, label: '유지' },
  alert: { color: 'red', icon: AlertTriangle, label: '경고' },
};

export default function InsightsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  return (
    <main className="min-h-screen bg-slate-50" aria-label="건강 인사이트">
      <DomainHeader
        title="주간 건강 인사이트"
        subtitle="AI가 분석한 이번 주 건강 변화와 맞춤 조언입니다"
        icon={Lightbulb}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-3xl mx-auto px-6 py-8 space-y-4">
        <p className="text-center text-xs text-slate-400 mb-6">
          2026년 5월 18일 ~ 5월 24일 기준 · 신뢰도는 AI 분석 확신 정도를 나타냅니다
        </p>

        {WEEKLY_INSIGHTS.map(insight => {
          const config = typeConfig[insight.type];
          const Icon = config.icon;
          return (
            <div key={insight.id} className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="flex items-start gap-4">
                <div className={`h-10 w-10 rounded-xl bg-${config.color}-50 flex items-center justify-center flex-shrink-0`}>
                  <Icon className={`h-5 w-5 text-${config.color}-600`} />
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold bg-${config.color}-50 text-${config.color}-700`}>
                      {config.label}
                    </span>
                    <span className="text-[10px] text-slate-400">신뢰도 {(insight.confidence * 100).toFixed(0)}%</span>
                  </div>
                  <h3 className="text-sm font-bold text-slate-900 mb-1">{insight.title}</h3>
                  <p className="text-xs text-slate-500 leading-relaxed mb-3">{insight.description}</p>

                  <div className="bg-slate-50 rounded-xl p-3 mb-3">
                    <p className="text-xs font-bold text-slate-700 mb-2">추천 조치</p>
                    <ul className="space-y-1">
                      {insight.actions.map((action, i) => (
                        <li key={i} className="flex items-start gap-2 text-xs text-slate-600">
                          <ArrowRight className="h-3 w-3 text-sky-500 mt-0.5 flex-shrink-0" />
                          {action}
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="flex items-center justify-between">
                    <span className={`text-xs font-bold text-${config.color}-600`}>{insight.metric}</span>
                    {insight.type === 'alert' && (
                      <Link href={`${localePrefix}/domains/telemedicine`} className="text-xs font-semibold text-sky-600 hover:underline">
                        전문의 상담 →
                      </Link>
                    )}
                  </div>
                </div>
              </div>
            </div>
          );
        })}

        <p className="text-center text-xs text-slate-400 mt-6">
          AI 인사이트는 참고 정보이며, 의학적 진단이나 처방을 대체하지 않습니다.
        </p>
      </div>
    </main>
  );
}
