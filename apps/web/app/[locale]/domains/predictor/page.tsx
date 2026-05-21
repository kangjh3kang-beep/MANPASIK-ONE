'use client';

import React from 'react';
import { Brain, AlertTriangle, ActivitySquare, Target } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePredictions, DiseasePrediction } from '@mmup/api-client';
import { RiskProgress, DomainHeader } from '@mmup/ui';

export default function PredictorPage() {
  const { data: session } = useSession();
  const { data: predictions, isLoading, error } = usePredictions();
  const user = session?.user as any;

  if (isLoading) return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center">
      <div className="flex flex-col items-center gap-4">
        <Brain className="w-12 h-12 text-emerald-500 animate-bounce" />
        <p className="text-slate-500 font-medium">GNN 예측 모델 분석 중...</p>
      </div>
    </div>
  );

  if (error) return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center">
      <div className="bg-white p-8 rounded-3xl border border-slate-200 shadow-xl max-w-md text-center">
        <AlertTriangle className="w-16 h-16 text-rose-500 mx-auto mb-4" />
        <h2 className="text-xl font-bold text-slate-900 mb-2">예측 데이터 로드 실패</h2>
        <button onClick={() => window.location.reload()} className="px-6 py-2 bg-emerald-600 text-white rounded-xl font-bold">재시도</button>
      </div>
    </div>
  );

  const data = predictions?.data || [];
  const highRiskCount = data.filter(p => p.riskScore >= 0.75).length;

  return (
    <div className="min-h-screen bg-slate-50">
      <DomainHeader 
        title="생체 지표 예측 엔진" 
        subtitle="GNN 기반 다중 질환 시계열 예측 및 데이터 기반 예방 의학"
        icon={ActivitySquare}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'researcher',
          organization: user.organization || 'MMUP AI Lab'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            { label: '분석 대상', val: `${data.length}개 질환`, icon: Target, c: 'text-emerald-600 bg-emerald-50' },
            { label: '고위험 감지', val: `${highRiskCount}건`, icon: AlertTriangle, c: highRiskCount > 0 ? 'text-red-600 bg-red-50' : 'text-emerald-600 bg-emerald-50' },
            { label: '모델 정확도', val: '96.3%', icon: Brain, c: 'text-violet-600 bg-violet-50' },
          ].map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-semibold text-slate-500">{s.label}</span>
                <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-bold text-slate-900">{s.val}</p>
            </div>
          ))}
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {data.map((pred: DiseasePrediction) => (
            <div key={pred.disease} className="rounded-2xl border border-slate-200 bg-white p-2 shadow-sm">
              <RiskProgress 
                diseaseName={pred.disease}
                riskScore={pred.riskScore}
                factors={pred.contributingFactors}
                trend={pred.trend as any}
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
