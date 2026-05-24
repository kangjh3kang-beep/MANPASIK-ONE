'use client';

import React from 'react';
import { DomainHeader } from '@mmup/ui';
import { Flag, Users, Trophy, CheckCircle, Star } from 'lucide-react';

interface Challenge {
  id: string;
  title: string;
  description: string;
  progressCurrent: number;
  progressTotal: number;
  unit: string;
  reward: number;
  participants: number;
  isJoined: boolean;
  isCompleted?: boolean;
}

interface LeaderboardEntry {
  rank: number;
  name: string;
  avatar: string;
  points: number;
  tier: string;
  tierColor: string;
}

const ACTIVE_CHALLENGES: Challenge[] = [
  {
    id: 'c-1',
    title: '30일 걷기 챌린지',
    description: '매일 6,000보 이상 걷기를 달성하세요',
    progressCurrent: 15,
    progressTotal: 30,
    unit: '일',
    reward: 500,
    participants: 234,
    isJoined: true,
  },
  {
    id: 'c-2',
    title: '매일 측정 챌린지',
    description: '14일 연속으로 건강 데이터를 측정하세요',
    progressCurrent: 7,
    progressTotal: 14,
    unit: '일',
    reward: 300,
    participants: 89,
    isJoined: true,
  },
  {
    id: 'c-3',
    title: '수면 개선 챌린지',
    description: '21일 동안 7시간 이상 수면을 달성하세요',
    progressCurrent: 3,
    progressTotal: 21,
    unit: '일',
    reward: 400,
    participants: 156,
    isJoined: false,
  },
];

const COMPLETED_CHALLENGES: Challenge[] = [
  {
    id: 'c-4',
    title: '첫 번째 측정 챌린지',
    description: '첫 건강 데이터를 측정하세요',
    progressCurrent: 1,
    progressTotal: 1,
    unit: '회',
    reward: 100,
    participants: 1247,
    isJoined: true,
    isCompleted: true,
  },
  {
    id: 'c-5',
    title: '데이터 기여 챌린지',
    description: '5회 이상 데이터를 기여하세요',
    progressCurrent: 5,
    progressTotal: 5,
    unit: '회',
    reward: 250,
    participants: 412,
    isJoined: true,
    isCompleted: true,
  },
];

const LEADERBOARD: LeaderboardEntry[] = [
  { rank: 1, name: '김건강', avatar: '🏆', points: 12450, tier: 'Platinum', tierColor: 'bg-violet-50 text-violet-700' },
  { rank: 2, name: '이활력', avatar: '🥈', points: 9820, tier: 'Gold', tierColor: 'bg-amber-50 text-amber-700' },
  { rank: 3, name: '박꾸준', avatar: '🥉', points: 8340, tier: 'Gold', tierColor: 'bg-amber-50 text-amber-700' },
  { rank: 4, name: '최노력', avatar: '👤', points: 6200, tier: 'Silver', tierColor: 'bg-slate-100 text-slate-600' },
  { rank: 5, name: '정성실', avatar: '👤', points: 5180, tier: 'Silver', tierColor: 'bg-slate-100 text-slate-600' },
];

function ProgressBar({ current, total, color = 'bg-sky-500' }: { current: number; total: number; color?: string }) {
  const pct = Math.min(100, (current / total) * 100);
  return (
    <div className="h-2.5 bg-slate-100 rounded-full overflow-hidden">
      <div className={`h-full rounded-full transition-all duration-500 ${color}`} style={{ width: `${pct}%` }} />
    </div>
  );
}

export default function ChallengesPage() {
  return (
    <main className="min-h-screen bg-slate-50" aria-label="건강 챌린지">
      <DomainHeader title="건강 챌린지" subtitle="건강 목표를 달성하고 MPS 포인트를 획득하세요" icon={Flag} />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-8">
        {/* 진행 중인 챌린지 */}
        <section aria-label="진행 중인 챌린지">
          <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Flag className="h-5 w-5 text-sky-600" />
            진행 중인 챌린지
          </h2>
          <div className="space-y-4">
            {ACTIVE_CHALLENGES.map((challenge) => {
              const pct = Math.round((challenge.progressCurrent / challenge.progressTotal) * 100);
              return (
                <div key={challenge.id} className="rounded-2xl border border-slate-200 bg-white p-5 hover:shadow-sm transition-shadow">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex-1">
                      <h3 className="text-base font-bold text-slate-900">{challenge.title}</h3>
                      <p className="text-sm text-slate-500 mt-0.5">{challenge.description}</p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0 ml-4">
                      <span className="text-sm font-bold text-amber-600 flex items-center gap-1">
                        <Trophy className="h-4 w-4" />
                        {challenge.reward} MPS
                      </span>
                    </div>
                  </div>

                  <div className="mb-3">
                    <div className="flex items-center justify-between mb-1.5">
                      <span className="text-xs text-slate-500">
                        {challenge.progressCurrent}/{challenge.progressTotal}{challenge.unit} 완료
                      </span>
                      <span className="text-xs font-bold text-slate-700">{pct}%</span>
                    </div>
                    <ProgressBar current={challenge.progressCurrent} total={challenge.progressTotal} />
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-1 text-xs text-slate-400">
                      <Users className="h-3.5 w-3.5" />
                      <span>{challenge.participants}명 참여</span>
                    </div>
                    <button
                      className={`px-4 py-2 rounded-xl text-sm font-semibold transition-all ${
                        challenge.isJoined
                          ? 'bg-sky-50 text-sky-600 border border-sky-200 cursor-default'
                          : 'bg-sky-600 text-white hover:bg-sky-700'
                      }`}
                      disabled={challenge.isJoined}
                    >
                      {challenge.isJoined ? '참여중' : '참여하기'}
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        </section>

        {/* 리더보드 */}
        <section aria-label="리더보드">
          <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Star className="h-5 w-5 text-amber-500" />
            리더보드 (상위 5명)
          </h2>
          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
            <div className="divide-y divide-slate-100">
              {LEADERBOARD.map((entry) => (
                <div key={entry.rank} className={`flex items-center gap-4 px-5 py-4 ${entry.rank <= 3 ? 'bg-amber-50/30' : ''}`}>
                  <span className="text-lg font-black text-slate-300 w-6 text-center tabular-nums">{entry.rank}</span>
                  <span className="text-2xl w-8 text-center" aria-hidden="true">{entry.avatar}</span>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-bold text-slate-900">{entry.name}</p>
                    <span className={`inline-block mt-0.5 px-2 py-0.5 rounded text-[10px] font-semibold ${entry.tierColor}`}>
                      {entry.tier}
                    </span>
                  </div>
                  <span className="text-sm font-bold text-amber-600 tabular-nums">{entry.points.toLocaleString()} MPS</span>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* 완료된 챌린지 */}
        <section aria-label="완료된 챌린지">
          <h2 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
            <CheckCircle className="h-5 w-5 text-emerald-500" />
            완료된 챌린지
          </h2>
          <div className="space-y-3">
            {COMPLETED_CHALLENGES.map((challenge) => (
              <div key={challenge.id} className="rounded-2xl border border-emerald-100 bg-emerald-50/50 p-5">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="h-9 w-9 rounded-xl bg-emerald-100 flex items-center justify-center" aria-hidden="true">
                      <CheckCircle className="h-5 w-5 text-emerald-600" />
                    </div>
                    <div>
                      <h3 className="text-sm font-bold text-slate-800">{challenge.title}</h3>
                      <p className="text-xs text-slate-500">{challenge.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-sm font-bold text-emerald-600">+{challenge.reward} MPS</span>
                    <span className="px-3 py-1 rounded-full text-[10px] font-semibold bg-emerald-100 text-emerald-700 border border-emerald-200">
                      보상 수령 완료
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
