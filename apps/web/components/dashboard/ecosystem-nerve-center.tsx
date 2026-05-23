import React from 'react';
import { Activity, Zap } from 'lucide-react';

export function EcosystemNerveCenter() {
  return (
    <div className="rounded-3xl bg-gradient-to-br from-slate-800 to-slate-900 p-6 border border-white/10">
      {/* 헤더 */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-base font-bold">건강 종합 모니터</h3>
          <p className="text-slate-400 text-[10px] mt-0.5">실시간 생체 균형 지표</p>
        </div>
        <div className="h-8 w-8 rounded-full bg-sky-500/10 flex items-center justify-center ring-1 ring-sky-500/30">
          <Activity className="h-4 w-4 text-sky-400 animate-pulse" />
        </div>
      </div>

      {/* 3개 원형 지표 */}
      <div className="grid grid-cols-3 gap-4 py-4">
        {[
          { label: '신체', value: 92, color: '#10b981', desc: '대사 동기화' },
          { label: '환경', value: 74, color: '#0ea5e9', desc: '미세먼지 지수' },
          { label: '정신', value: 88, color: '#8b5cf6', desc: '뇌파 평균' },
        ].map((item) => (
          <div key={item.label} className="text-center">
            <div className="relative inline-flex items-center justify-center mb-2">
              <svg className="h-20 w-20 -rotate-90">
                <circle className="text-white/5" strokeWidth="5" stroke="currentColor" fill="transparent" r="34" cx="40" cy="40" />
                <circle
                  strokeWidth="5"
                  strokeDasharray={2 * Math.PI * 34}
                  strokeDashoffset={2 * Math.PI * 34 * (1 - item.value / 100)}
                  strokeLinecap="round"
                  stroke={item.color}
                  fill="transparent"
                  r="34" cx="40" cy="40"
                />
              </svg>
              <span className="absolute text-lg font-black">{item.value}%</span>
            </div>
            <p className="text-xs font-bold text-white">{item.label}</p>
            <p className="text-[9px] text-slate-500">{item.desc}</p>
          </div>
        ))}
      </div>

      {/* AI 분석 로그 */}
      <div className="mt-4 pt-4 border-t border-white/5">
        <h4 className="text-[10px] font-bold text-slate-500 mb-3 flex items-center gap-1">
          <Zap className="h-3 w-3 text-amber-500" /> AI 분석 현황
        </h4>
        <div className="space-y-2">
          {[
            { label: '건강 데이터', msg: '데이터 무결성 확인 완료', time: '방금' },
            { label: '위험 예측', msg: 'AI 예측 모델 갱신 완료', time: '2초 전' },
            { label: '보상', msg: '데이터 기여 보상 처리 완료', time: '12초 전' },
            { label: '품질 관리', msg: '전자 서명 검증 완료', time: '1분 전' },
          ].map((log, i) => (
            <div key={i} className="flex items-start gap-2 text-[11px]">
              <div className="h-1.5 w-1.5 rounded-full bg-sky-500 mt-1.5 flex-shrink-0" />
              <div className="flex-1 min-w-0">
                <span className="font-bold text-slate-400">{log.label}</span>
                <span className="text-slate-500 ml-1">{log.msg}</span>
              </div>
              <span className="text-slate-600 flex-shrink-0">{log.time}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
