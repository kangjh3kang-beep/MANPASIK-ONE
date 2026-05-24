'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Activity, Flame, MapPin, Timer, Plus } from 'lucide-react';

interface ActivityRing {
  label: string;
  current: number;
  goal: number;
  unit: string;
  color: string;
  strokeColor: string;
}

const RINGS: ActivityRing[] = [
  { label: '걸음 수', current: 7823, goal: 10000, unit: '걸음', color: 'text-sky-500', strokeColor: '#0ea5e9' },
  { label: '운동 시간', current: 25, goal: 30, unit: '분', color: 'text-emerald-500', strokeColor: '#10b981' },
  { label: '서있기', current: 9, goal: 12, unit: '시간', color: 'text-amber-500', strokeColor: '#f59e0b' },
];

const HOURLY_ACTIVITY = [
  0, 0, 0, 0, 0, 0, 20, 45, 60, 30, 15, 40, 55, 25, 10, 35, 50, 70, 45, 30, 20, 10, 5, 0,
];

const WEEKLY_STEPS = [8200, 6500, 9100, 7800, 10200, 5600, 7823];
const WEEK_LABELS = ['월', '화', '수', '목', '금', '토', '일'];

const EXERCISE_TYPES = ['걷기', '달리기', '자전거', '수영', '기타'];

function RingsSVG({ rings }: { rings: ActivityRing[] }) {
  const size = 180;
  const center = size / 2;
  const radii = [70, 55, 40];
  const strokeWidth = 10;

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} className="mx-auto">
      {rings.map((ring, i) => {
        const r = radii[i];
        const circumference = 2 * Math.PI * r;
        const percent = Math.min(ring.current / ring.goal, 1);
        const dashOffset = circumference * (1 - percent);
        return (
          <g key={ring.label}>
            <circle
              cx={center}
              cy={center}
              r={r}
              fill="none"
              stroke="#e2e8f0"
              strokeWidth={strokeWidth}
            />
            <circle
              cx={center}
              cy={center}
              r={r}
              fill="none"
              stroke={ring.strokeColor}
              strokeWidth={strokeWidth}
              strokeLinecap="round"
              strokeDasharray={circumference}
              strokeDashoffset={dashOffset}
              transform={`rotate(-90 ${center} ${center})`}
              className="transition-all duration-700"
            />
          </g>
        );
      })}
    </svg>
  );
}

export default function ActivityPage() {
  const [showAddForm, setShowAddForm] = useState(false);
  const [selectedExercise, setSelectedExercise] = useState(EXERCISE_TYPES[0]);

  const maxHourly = Math.max(...HOURLY_ACTIVITY);
  const maxWeekly = Math.max(...WEEKLY_STEPS);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="활동 기록">
      <DomainHeader
        title="활동 기록"
        subtitle="일일 활동량과 운동 기록을 확인합니다"
        icon={Activity}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 활동 링 + 요약 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-center">
            <div>
              <RingsSVG rings={RINGS} />
            </div>
            <div className="space-y-3">
              {RINGS.map((ring) => {
                const percent = Math.round((ring.current / ring.goal) * 100);
                return (
                  <div key={ring.label} className="flex items-center justify-between p-3 rounded-xl bg-slate-50">
                    <div className="flex items-center gap-2">
                      <div className={`h-3 w-3 rounded-full`} style={{ backgroundColor: ring.strokeColor }} />
                      <span className="text-sm font-medium text-slate-700">{ring.label}</span>
                    </div>
                    <div className="text-right">
                      <span className="text-sm font-bold text-slate-900">
                        {ring.current.toLocaleString()}/{ring.goal.toLocaleString()} {ring.unit}
                      </span>
                      <span className="ml-2 text-xs text-slate-400">({percent}%)</span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>

        {/* 일일 요약 카드 */}
        <section aria-label="일일 요약" className="grid grid-cols-3 gap-3">
          {[
            { label: '소모 칼로리', value: '342', unit: 'kcal', icon: Flame, color: 'text-orange-500' },
            { label: '이동 거리', value: '5.2', unit: 'km', icon: MapPin, color: 'text-blue-500' },
            { label: '활동 시간', value: '48', unit: '분', icon: Timer, color: 'text-purple-500' },
          ].map((card) => (
            <div key={card.label} className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
              <card.icon className={`h-5 w-5 mx-auto mb-2 ${card.color}`} />
              <p className="text-xl font-bold text-slate-900">
                {card.value}<span className="text-xs font-normal text-slate-400 ml-0.5">{card.unit}</span>
              </p>
              <p className="text-xs text-slate-500 mt-1">{card.label}</p>
            </div>
          ))}
        </section>

        {/* 시간별 활동 바 차트 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">시간별 활동량</h2>
          <div className="flex items-end gap-1 h-32">
            {HOURLY_ACTIVITY.map((value, idx) => (
              <div key={idx} className="flex-1 flex flex-col items-center justify-end h-full">
                <div
                  className="w-full rounded-t bg-sky-400 transition-all hover:bg-sky-500"
                  style={{ height: `${maxHourly > 0 ? (value / maxHourly) * 100 : 0}%`, minHeight: value > 0 ? '2px' : '0' }}
                  title={`${idx}시: ${value}`}
                />
              </div>
            ))}
          </div>
          <div className="flex justify-between mt-2">
            {[0, 6, 12, 18, 23].map((h) => (
              <span key={h} className="text-[10px] text-slate-400">{h}시</span>
            ))}
          </div>
        </section>

        {/* 주간 걸음 수 추이 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">주간 걸음 수 추이</h2>
          <div className="flex items-center justify-between h-8">
            {WEEKLY_STEPS.map((steps, idx) => {
              const size = 8 + (steps / maxWeekly) * 16;
              return (
                <div key={idx} className="flex flex-col items-center gap-2">
                  <div
                    className="rounded-full bg-sky-500"
                    style={{ width: `${size}px`, height: `${size}px` }}
                    title={`${WEEK_LABELS[idx]}: ${steps.toLocaleString()}걸음`}
                  />
                </div>
              );
            })}
          </div>
          <div className="flex justify-between mt-3">
            {WEEK_LABELS.map((label, idx) => (
              <div key={idx} className="text-center">
                <p className="text-[10px] text-slate-400">{label}</p>
                <p className="text-[10px] font-semibold text-slate-600">{(WEEKLY_STEPS[idx] / 1000).toFixed(1)}k</p>
              </div>
            ))}
          </div>
        </section>

        {/* 운동 기록 추가 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-slate-900">운동 기록</h2>
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              운동 기록 추가
            </button>
          </div>

          {showAddForm && (
            <div className="p-4 rounded-xl bg-slate-50 border border-slate-100 space-y-3">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">운동 종류</label>
                <div className="flex flex-wrap gap-1.5">
                  {EXERCISE_TYPES.map((type) => (
                    <button
                      key={type}
                      onClick={() => setSelectedExercise(type)}
                      className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors ${
                        selectedExercise === type
                          ? 'bg-sky-600 text-white'
                          : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-100'
                      }`}
                    >
                      {type}
                    </button>
                  ))}
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">시간 (분)</label>
                  <input
                    type="number"
                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                    placeholder="30"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">거리 (km)</label>
                  <input
                    type="number"
                    step="0.1"
                    className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                    placeholder="3.0"
                  />
                </div>
              </div>
              <div className="flex gap-2 pt-1">
                <button
                  onClick={() => setShowAddForm(false)}
                  className="px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-100"
                >
                  취소
                </button>
                <button className="px-5 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors">
                  저장
                </button>
              </div>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
