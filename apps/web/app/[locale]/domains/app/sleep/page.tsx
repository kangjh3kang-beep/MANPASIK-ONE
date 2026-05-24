'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Moon, BedDouble, AlarmClock, Lightbulb, Target } from 'lucide-react';

interface SleepStage {
  label: string;
  duration: number; // minutes
  color: string;
  bgClass: string;
}

const LAST_NIGHT_STAGES: SleepStage[] = [
  { label: '깊은수면', duration: 95, color: '#6366f1', bgClass: 'bg-indigo-500' },
  { label: '얕은수면', duration: 180, color: '#0ea5e9', bgClass: 'bg-sky-400' },
  { label: 'REM', duration: 90, color: '#a855f7', bgClass: 'bg-purple-500' },
  { label: '깨어남', duration: 15, color: '#f59e0b', bgClass: 'bg-amber-400' },
];

const WEEKLY_SLEEP = [
  { day: '월', hours: 7.2 },
  { day: '화', hours: 6.5 },
  { day: '수', hours: 8.0 },
  { day: '목', hours: 7.0 },
  { day: '금', hours: 6.8 },
  { day: '토', hours: 8.5 },
  { day: '일', hours: 6.3 },
];

const SLEEP_TIPS = [
  {
    title: '취침 1시간 전 블루라이트 차단',
    description: '최근 야간 스크린 사용 시간이 증가했습니다. 취침 전 블루라이트 차단 모드를 활성화하세요.',
  },
  {
    title: '일정한 기상 시간 유지',
    description: '주말 기상 시간이 평일보다 2시간 이상 늦습니다. 일정한 기상 시간이 수면의 질을 높입니다.',
  },
  {
    title: '카페인 섭취 시간 조절',
    description: '오후 2시 이후 카페인 섭취를 줄이면 깊은 수면 비율이 개선될 수 있습니다.',
  },
];

export default function SleepPage() {
  const [targetHours, setTargetHours] = useState(7.5);
  const [bedtimeReminder, setBedtimeReminder] = useState(true);

  const totalMinutes = LAST_NIGHT_STAGES.reduce((sum, s) => sum + s.duration, 0);
  const totalHours = Math.floor(totalMinutes / 60);
  const totalMins = totalMinutes % 60;
  const sleepScore = 78;

  const maxWeeklyHours = Math.max(...WEEKLY_SLEEP.map((d) => d.hours));

  return (
    <main className="min-h-screen bg-slate-50" aria-label="수면 분석">
      <DomainHeader
        title="수면 분석"
        subtitle="수면 패턴을 분석하고 수면의 질을 개선합니다"
        icon={Moon}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 지난밤 요약 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">지난밤 수면 요약</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="rounded-xl bg-indigo-50 p-4 text-center">
              <p className="text-xs font-semibold text-indigo-500 mb-1">총 수면 시간</p>
              <p className="text-2xl font-bold text-indigo-900">{totalHours}시간 {totalMins}분</p>
            </div>
            <div className="rounded-xl bg-purple-50 p-4 text-center">
              <p className="text-xs font-semibold text-purple-500 mb-1">수면 품질 점수</p>
              <p className="text-2xl font-bold text-purple-900">{sleepScore}<span className="text-sm font-normal text-purple-400">/100</span></p>
            </div>
            <div className="rounded-xl bg-sky-50 p-4 text-center">
              <p className="text-xs font-semibold text-sky-500 mb-1">깊은수면 비율</p>
              <p className="text-2xl font-bold text-sky-900">{Math.round((95 / totalMinutes) * 100)}%</p>
            </div>
          </div>
        </section>

        {/* 수면 단계 바 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">수면 단계</h2>
          <div className="h-8 rounded-full overflow-hidden flex">
            {LAST_NIGHT_STAGES.map((stage) => (
              <div
                key={stage.label}
                className={`${stage.bgClass} transition-all`}
                style={{ width: `${(stage.duration / totalMinutes) * 100}%` }}
                title={`${stage.label}: ${stage.duration}분`}
              />
            ))}
          </div>
          <div className="flex flex-wrap gap-4 mt-3">
            {LAST_NIGHT_STAGES.map((stage) => (
              <span key={stage.label} className="flex items-center gap-1.5 text-xs text-slate-600">
                <span className={`inline-block h-2.5 w-2.5 rounded-full ${stage.bgClass}`} />
                {stage.label} {stage.duration}분 ({Math.round((stage.duration / totalMinutes) * 100)}%)
              </span>
            ))}
          </div>
        </section>

        {/* 취침/기상 시간 */}
        <section className="grid grid-cols-2 gap-3">
          <div className="rounded-2xl border border-slate-200 bg-white p-5 flex items-center gap-4">
            <div className="h-12 w-12 rounded-xl bg-indigo-50 flex items-center justify-center">
              <BedDouble className="h-6 w-6 text-indigo-500" />
            </div>
            <div>
              <p className="text-xs font-semibold text-slate-400">취침 시간</p>
              <p className="text-lg font-bold text-slate-900">23:15</p>
            </div>
          </div>
          <div className="rounded-2xl border border-slate-200 bg-white p-5 flex items-center gap-4">
            <div className="h-12 w-12 rounded-xl bg-amber-50 flex items-center justify-center">
              <AlarmClock className="h-6 w-6 text-amber-500" />
            </div>
            <div>
              <p className="text-xs font-semibold text-slate-400">기상 시간</p>
              <p className="text-lg font-bold text-slate-900">06:35</p>
            </div>
          </div>
        </section>

        {/* 7일 수면 시간 추이 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">7일 수면 시간 추이</h2>
          <div className="flex items-end gap-2 h-36">
            {WEEKLY_SLEEP.map((day) => {
              const heightPercent = (day.hours / maxWeeklyHours) * 100;
              return (
                <div key={day.day} className="flex-1 flex flex-col items-center justify-end h-full">
                  <p className="text-[10px] font-semibold text-slate-600 mb-1">{day.hours}h</p>
                  <div
                    className="w-full rounded-t-lg bg-indigo-400 hover:bg-indigo-500 transition-colors"
                    style={{ height: `${heightPercent}%` }}
                    title={`${day.day}: ${day.hours}시간`}
                  />
                  <p className="text-[10px] text-slate-400 mt-2">{day.day}</p>
                </div>
              );
            })}
          </div>
        </section>

        {/* AI 수면 팁 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Lightbulb className="h-4 w-4 text-amber-500" /> AI 수면 개선 팁
          </h2>
          <div className="space-y-3">
            {SLEEP_TIPS.map((tip, idx) => (
              <div key={idx} className="p-4 rounded-xl bg-amber-50 border border-amber-100">
                <p className="text-sm font-semibold text-amber-900">{tip.title}</p>
                <p className="text-xs text-amber-700 mt-1">{tip.description}</p>
              </div>
            ))}
          </div>
          <p className="text-[10px] text-slate-400 mt-3">
            이 팁은 AI가 수면 패턴을 분석하여 생성하였으며, 의료 조언이 아닙니다. [추정]
          </p>
        </section>

        {/* 수면 목표 설정 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Target className="h-4 w-4 text-sky-600" /> 수면 목표 설정
          </h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">목표 수면 시간</label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="5"
                  max="10"
                  step="0.5"
                  value={targetHours}
                  onChange={(e) => setTargetHours(parseFloat(e.target.value))}
                  className="flex-1 accent-sky-600"
                />
                <span className="text-lg font-bold text-slate-900 w-16 text-right">{targetHours}시간</span>
              </div>
            </div>
            <label className="flex items-center justify-between cursor-pointer">
              <div>
                <p className="text-sm font-medium text-slate-700">취침 시간 알림</p>
                <p className="text-xs text-slate-400">목표 취침 시간 30분 전에 알림을 보냅니다</p>
              </div>
              <div
                role="switch"
                aria-checked={bedtimeReminder}
                onClick={() => setBedtimeReminder(!bedtimeReminder)}
                className={`relative w-11 h-6 rounded-full cursor-pointer transition-colors ${
                  bedtimeReminder ? 'bg-sky-600' : 'bg-slate-300'
                }`}
              >
                <div
                  className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${
                    bedtimeReminder ? 'translate-x-[22px]' : 'translate-x-0.5'
                  }`}
                />
              </div>
            </label>
          </div>
        </section>
      </div>
    </main>
  );
}
