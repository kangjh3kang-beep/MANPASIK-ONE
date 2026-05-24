'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Target, Plus, CheckCircle, TrendingUp, Flame, Award } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type GoalStatus = 'active' | 'completed' | 'paused';

const MOCK_GOALS = [
  { id: '1', title: '공복 혈당 100 이하 유지', category: '혈당', target: 100, current: 108, unit: 'mg/dL', status: 'active' as GoalStatus, streak: 12, progress: 72, milestones: ['1주 달성', '2주 달성', '1개월 달성'], completedMilestones: 2 },
  { id: '2', title: '매일 30분 이상 걷기', category: '운동', target: 30, current: 25, unit: '분', status: 'active' as GoalStatus, streak: 5, progress: 83, milestones: ['7일 연속', '14일 연속', '30일 연속'], completedMilestones: 0 },
  { id: '3', title: '체중 75kg 도달', category: '체중', target: 75, current: 78.2, unit: 'kg', status: 'active' as GoalStatus, streak: 0, progress: 56, milestones: ['79kg 달성', '77kg 달성', '75kg 달성'], completedMilestones: 1 },
  { id: '4', title: '수면 7시간 이상', category: '수면', target: 7, current: 7.2, unit: '시간', status: 'completed' as GoalStatus, streak: 30, progress: 100, milestones: ['1주 달성', '2주 달성', '1개월 달성'], completedMilestones: 3 },
];

export default function GoalsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [goals, setGoals] = useState(MOCK_GOALS);
  const [showAddForm, setShowAddForm] = useState(false);

  const activeGoals = goals.filter(g => g.status === 'active');
  const completedGoals = goals.filter(g => g.status === 'completed');
  const totalStreak = Math.max(...goals.map(g => g.streak));

  return (
    <main className="min-h-screen bg-slate-50" aria-label="목표 관리">
      <DomainHeader
        title="목표 관리"
        subtitle="건강 목표를 설정하고 달성 현황을 추적합니다"
        icon={Target}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-4xl mx-auto px-6 py-8 space-y-6">
        {/* 요약 카드 */}
        <div className="grid grid-cols-3 gap-4">
          <div className="rounded-2xl border border-slate-200 bg-white p-5 text-center">
            <div className="h-10 w-10 rounded-xl bg-sky-50 flex items-center justify-center mx-auto mb-2">
              <Target className="h-5 w-5 text-sky-600" />
            </div>
            <p className="text-2xl font-bold text-slate-900">{activeGoals.length}</p>
            <p className="text-xs text-slate-500">진행 중인 목표</p>
          </div>
          <div className="rounded-2xl border border-slate-200 bg-white p-5 text-center">
            <div className="h-10 w-10 rounded-xl bg-emerald-50 flex items-center justify-center mx-auto mb-2">
              <CheckCircle className="h-5 w-5 text-emerald-600" />
            </div>
            <p className="text-2xl font-bold text-slate-900">{completedGoals.length}</p>
            <p className="text-xs text-slate-500">달성 완료</p>
          </div>
          <div className="rounded-2xl border border-slate-200 bg-white p-5 text-center">
            <div className="h-10 w-10 rounded-xl bg-orange-50 flex items-center justify-center mx-auto mb-2">
              <Flame className="h-5 w-5 text-orange-600" />
            </div>
            <p className="text-2xl font-bold text-slate-900">{totalStreak}일</p>
            <p className="text-xs text-slate-500">최장 연속 달성</p>
          </div>
        </div>

        {/* 목표 추가 버튼 */}
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-2xl border-2 border-dashed border-slate-300 text-slate-500 hover:border-sky-400 hover:text-sky-600 transition-colors"
        >
          <Plus className="h-4 w-4" />
          새 목표 추가
        </button>

        {/* 새 목표 입력 폼 */}
        {showAddForm && (
          <div className="rounded-2xl border border-sky-200 bg-sky-50/50 p-6 space-y-4">
            <h3 className="text-sm font-bold text-slate-900">새 건강 목표 설정</h3>
            <input type="text" placeholder="목표를 입력하세요 (예: 공복 혈당 100 이하 유지)" className="w-full px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/30" />
            <div className="grid grid-cols-3 gap-3">
              <input type="number" placeholder="목표값" className="px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/30" />
              <input type="text" placeholder="단위" className="px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/30" />
              <select className="px-4 py-3 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/30">
                <option>혈당</option><option>체중</option><option>운동</option><option>수면</option><option>식이</option>
              </select>
            </div>
            <div className="flex gap-2 justify-end">
              <button onClick={() => setShowAddForm(false)} className="px-4 py-2 rounded-xl text-sm font-medium text-slate-500 hover:bg-slate-100">취소</button>
              <button className="px-6 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700">저장</button>
            </div>
          </div>
        )}

        {/* 진행 중인 목표 */}
        <section>
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-sky-600" /> 진행 중
          </h2>
          <div className="space-y-3">
            {activeGoals.map(goal => (
              <div key={goal.id} className="rounded-2xl border border-slate-200 bg-white p-5">
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h3 className="text-sm font-bold text-slate-900">{goal.title}</h3>
                    <span className="inline-block mt-1 px-2 py-0.5 rounded-full bg-sky-50 text-sky-600 text-[10px] font-semibold">{goal.category}</span>
                  </div>
                  {goal.streak > 0 && (
                    <div className="flex items-center gap-1 px-2 py-1 rounded-lg bg-orange-50">
                      <Flame className="h-3 w-3 text-orange-500" />
                      <span className="text-xs font-bold text-orange-600">{goal.streak}일 연속</span>
                    </div>
                  )}
                </div>

                <div className="flex items-center gap-3 mb-3">
                  <div className="flex-1">
                    <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
                      <div className="h-full bg-sky-500 rounded-full transition-all" style={{ width: `${goal.progress}%` }} />
                    </div>
                  </div>
                  <span className="text-sm font-bold text-slate-700">{goal.progress}%</span>
                </div>

                <div className="flex items-center justify-between text-xs text-slate-500">
                  <span>현재: <strong className="text-slate-700">{goal.current}{goal.unit}</strong></span>
                  <span>목표: <strong className="text-sky-600">{goal.target}{goal.unit}</strong></span>
                </div>

                {/* 마일스톤 */}
                <div className="mt-3 flex items-center gap-2">
                  {goal.milestones.map((m, i) => (
                    <div key={m} className={`flex items-center gap-1 px-2 py-1 rounded-lg text-[10px] font-medium ${
                      i < goal.completedMilestones ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-50 text-slate-400'
                    }`}>
                      {i < goal.completedMilestones && <CheckCircle className="h-3 w-3" />}
                      {m}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* 달성 완료 */}
        {completedGoals.length > 0 && (
          <section>
            <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
              <Award className="h-4 w-4 text-emerald-600" /> 달성 완료
            </h2>
            <div className="space-y-3">
              {completedGoals.map(goal => (
                <div key={goal.id} className="rounded-2xl border border-emerald-200 bg-emerald-50/50 p-5 opacity-80">
                  <div className="flex items-center gap-3">
                    <CheckCircle className="h-5 w-5 text-emerald-500 flex-shrink-0" />
                    <div>
                      <h3 className="text-sm font-bold text-slate-900">{goal.title}</h3>
                      <p className="text-xs text-slate-500 mt-0.5">{goal.streak}일 연속 달성 완료</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </main>
  );
}
