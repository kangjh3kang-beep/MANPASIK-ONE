'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Users, ChevronDown, ChevronUp, Search } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type UserRole = '환자' | '의사' | '연구원' | '관리자';
type UserStatus = '활성' | '비활성' | '정지';

interface ManagedUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  joinDate: string;
  lastAccess: string;
  measurements: number;
  subscription: string;
  devices: number;
}

const ROLE_STYLES: Record<UserRole, string> = {
  '환자': 'bg-sky-50 text-sky-700 border-sky-200',
  '의사': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '연구원': 'bg-violet-50 text-violet-700 border-violet-200',
  '관리자': 'bg-amber-50 text-amber-700 border-amber-200',
};

const STATUS_STYLES: Record<UserStatus, string> = {
  '활성': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '비활성': 'bg-slate-50 text-slate-500 border-slate-200',
  '정지': 'bg-red-50 text-red-700 border-red-200',
};

const MOCK_USERS: ManagedUser[] = [
  { id: 'u1', name: '김민수', email: 'minsu.kim@email.com', role: '환자', status: '활성', joinDate: '2025-08-15', lastAccess: '2026-05-24', measurements: 234, subscription: '프리미엄', devices: 2 },
  { id: 'u2', name: '이서연', email: 'seoyeon.lee@hospital.kr', role: '의사', status: '활성', joinDate: '2025-06-01', lastAccess: '2026-05-23', measurements: 0, subscription: '의료진', devices: 1 },
  { id: 'u3', name: '박지훈', email: 'jihoon.park@research.ac.kr', role: '연구원', status: '활성', joinDate: '2025-09-20', lastAccess: '2026-05-22', measurements: 1520, subscription: '연구용', devices: 5 },
  { id: 'u4', name: '최유나', email: 'yuna.choi@email.com', role: '환자', status: '비활성', joinDate: '2025-11-03', lastAccess: '2026-03-10', measurements: 45, subscription: '무료', devices: 1 },
  { id: 'u5', name: '정대현', email: 'daehyun.jung@admin.mps.kr', role: '관리자', status: '활성', joinDate: '2025-01-10', lastAccess: '2026-05-24', measurements: 0, subscription: '관리자', devices: 3 },
  { id: 'u6', name: '한소희', email: 'sohee.han@email.com', role: '환자', status: '정지', joinDate: '2025-12-22', lastAccess: '2026-04-15', measurements: 12, subscription: '기본', devices: 1 },
  { id: 'u7', name: '오태준', email: 'taejun.oh@hospital.kr', role: '의사', status: '활성', joinDate: '2026-01-05', lastAccess: '2026-05-24', measurements: 0, subscription: '의료진', devices: 2 },
  { id: 'u8', name: '송예린', email: 'yerin.song@email.com', role: '환자', status: '활성', joinDate: '2026-03-18', lastAccess: '2026-05-21', measurements: 78, subscription: '기본', devices: 1 },
];

const ROLE_FILTERS: (UserRole | '전체')[] = ['전체', '환자', '의사', '연구원', '관리자'];

export default function UsersPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<UserRole | '전체'>('전체');
  const [expandedUser, setExpandedUser] = useState<string | null>(null);

  const filteredUsers = MOCK_USERS.filter(u => {
    const matchesSearch = searchQuery === '' || u.name.includes(searchQuery) || u.email.includes(searchQuery) || u.id.includes(searchQuery);
    const matchesRole = roleFilter === '전체' || u.role === roleFilter;
    return matchesSearch && matchesRole;
  });

  const stats = [
    { label: '전체 사용자', value: MOCK_USERS.length },
    { label: '오늘 신규', value: 2 },
    { label: '활성 사용자', value: MOCK_USERS.filter(u => u.status === '활성').length },
    { label: '정지 계정', value: MOCK_USERS.filter(u => u.status === '정지').length },
  ];

  return (
    <main className="min-h-screen bg-slate-50" aria-label="사용자 관리">
      <DomainHeader
        title="사용자 관리"
        subtitle="플랫폼 사용자를 조회하고 관리합니다"
        icon={Users}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* 통계 카드 */}
        <div className="grid grid-cols-4 gap-4 mb-6">
          {stats.map(s => (
            <div key={s.label} className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
              <p className="text-xs text-slate-500 font-medium mb-1">{s.label}</p>
              <p className="text-2xl font-bold text-slate-900">{s.value}</p>
            </div>
          ))}
        </div>

        {/* 검색 및 필터 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-4 mb-6">
          <div className="flex items-center gap-3 mb-3">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
              <input
                type="text"
                placeholder="이름, 이메일, ID로 검색"
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                className="w-full pl-9 pr-4 py-2 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500"
              />
            </div>
          </div>
          <div className="flex gap-2">
            {ROLE_FILTERS.map(role => (
              <button
                key={role}
                onClick={() => setRoleFilter(role)}
                className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors ${
                  roleFilter === role
                    ? 'bg-sky-600 text-white'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                {role}
              </button>
            ))}
          </div>
        </div>

        {/* 사용자 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="grid grid-cols-[1fr_1.5fr_0.6fr_0.6fr_0.8fr_0.8fr_40px] gap-2 px-5 py-3 bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500">
            <span>이름</span><span>이메일</span><span>역할</span><span>상태</span><span>가입일</span><span>최근접속</span><span />
          </div>
          {filteredUsers.map(u => {
            const isExpanded = expandedUser === u.id;
            return (
              <div key={u.id} className="border-b border-slate-100 last:border-b-0">
                <button
                  onClick={() => setExpandedUser(prev => prev === u.id ? null : u.id)}
                  className="w-full grid grid-cols-[1fr_1.5fr_0.6fr_0.6fr_0.8fr_0.8fr_40px] gap-2 items-center px-5 py-3 text-left hover:bg-slate-50 transition-colors"
                  aria-expanded={isExpanded}
                  aria-label={`${u.name} 상세 보기`}
                >
                  <span className="text-sm font-medium text-slate-900 truncate">{u.name}</span>
                  <span className="text-sm text-slate-500 truncate">{u.email}</span>
                  <span className={`inline-flex justify-center px-2 py-0.5 rounded-full text-[11px] font-semibold border ${ROLE_STYLES[u.role]}`}>{u.role}</span>
                  <span className={`inline-flex justify-center px-2 py-0.5 rounded-full text-[11px] font-semibold border ${STATUS_STYLES[u.status]}`}>{u.status}</span>
                  <span className="text-xs text-slate-500">{u.joinDate}</span>
                  <span className="text-xs text-slate-500">{u.lastAccess}</span>
                  {isExpanded ? <ChevronUp className="h-4 w-4 text-slate-400" /> : <ChevronDown className="h-4 w-4 text-slate-400" />}
                </button>
                {isExpanded && (
                  <div className="px-5 pb-4 pt-1 bg-slate-50/50">
                    <div className="grid grid-cols-3 gap-3 mb-3">
                      <div className="rounded-xl bg-white border border-slate-200 p-3 text-center">
                        <p className="text-xs text-slate-500">측정 횟수</p>
                        <p className="text-lg font-bold text-slate-900">{u.measurements}</p>
                      </div>
                      <div className="rounded-xl bg-white border border-slate-200 p-3 text-center">
                        <p className="text-xs text-slate-500">구독 상태</p>
                        <p className="text-lg font-bold text-slate-900">{u.subscription}</p>
                      </div>
                      <div className="rounded-xl bg-white border border-slate-200 p-3 text-center">
                        <p className="text-xs text-slate-500">디바이스 수</p>
                        <p className="text-lg font-bold text-slate-900">{u.devices}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <select className="px-3 py-1.5 rounded-lg border border-slate-200 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-sky-500" defaultValue={u.role} aria-label="역할 변경">
                        <option value="환자">환자</option>
                        <option value="의사">의사</option>
                        <option value="연구원">연구원</option>
                        <option value="관리자">관리자</option>
                      </select>
                      <button className={`px-4 py-1.5 rounded-lg text-sm font-semibold transition-colors ${
                        u.status === '정지'
                          ? 'bg-emerald-600 text-white hover:bg-emerald-700'
                          : 'bg-red-600 text-white hover:bg-red-700'
                      }`}>
                        {u.status === '정지' ? '활성화' : '계정 정지'}
                      </button>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
