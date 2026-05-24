'use client';

import React from 'react';
import Link from 'next/link';
import { LogOut, User, Bell, Search, Activity } from 'lucide-react';

interface DomainHeaderProps {
  title: string;
  subtitle?: string;
  user?: {
    name: string;
    email: string;
    persona: string;
    organization?: string;
  };
  onLogout?: () => void;
  icon?: React.ElementType;
  certified?: boolean;
}

export function DomainHeader({
  title,
  subtitle,
  user,
  onLogout,
  icon: Icon = Activity,
  certified = false
}: DomainHeaderProps) {

  const personaThemes: Record<string, { color: string, bg: string, label: string }> = {
    doctor: { color: 'text-sky-700', bg: 'bg-sky-50', label: '의사' },
    researcher: { color: 'text-emerald-700', bg: 'bg-emerald-50', label: '연구원' },
    pharma: { color: 'text-violet-700', bg: 'bg-violet-50', label: '제약' },
    admin: { color: 'text-amber-700', bg: 'bg-amber-50', label: '관리자' },
    patient: { color: 'text-rose-700', bg: 'bg-rose-50', label: '환자' },
  };

  const theme = user?.persona ? personaThemes[user.persona] : { color: 'text-slate-700', bg: 'bg-slate-50', label: '사용자' };

  return (
    <header className="sticky top-0 z-40 w-full bg-white/95 backdrop-blur-xl border-b border-slate-200/60">
      <div className="max-w-7xl mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Left: Logo + Title */}
        <div className="flex items-center gap-4 min-w-0">
          <Link href="/" className="flex items-center gap-2.5 hover:opacity-80 transition-opacity flex-shrink-0">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-900 text-white">
              <Activity className="h-4 w-4" />
            </div>
            <span className="text-lg font-black tracking-tight text-slate-900 hidden sm:block">
              MPS<span className="text-[#0891B2]">.</span>
            </span>
          </Link>

          <div className="h-6 w-px bg-slate-200 flex-shrink-0" />

          <div className="flex items-center gap-2.5 min-w-0">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-900/5 text-slate-700 flex-shrink-0">
              <Icon className="h-4 w-4" />
            </div>
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <h1 className="text-sm font-bold text-slate-900 truncate">{title}</h1>
                {certified && (
                  <span className="flex-shrink-0 inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded bg-emerald-50 text-emerald-700 text-[10px] font-bold">
                    <svg className="h-2.5 w-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
                    </svg>
                    인증
                  </span>
                )}
              </div>
              {subtitle && <p className="text-[11px] text-slate-400 truncate">{subtitle}</p>}
            </div>
          </div>
        </div>

        {/* Center: Search */}
        <div className="hidden lg:flex flex-1 max-w-sm mx-8">
          <div className="relative w-full">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-slate-300" />
            <input
              type="text"
              placeholder="검색..."
              aria-label="검색"
              className="w-full pl-9 pr-4 py-2 bg-slate-50 border border-slate-200/60 rounded-lg text-xs focus:outline-none focus:ring-2 focus:ring-[#0891B2]/20 focus:border-[#0891B2]/40 transition-all placeholder:text-slate-300"
            />
          </div>
        </div>

        {/* Right: Actions */}
        <div className="flex items-center gap-2">
          <button className="relative p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-50 rounded-lg transition-all" aria-label="알림">
            <Bell className="h-4 w-4" />
            <span className="absolute top-1.5 right-1.5 h-1.5 w-1.5 rounded-full bg-rose-500" />
          </button>

          {user ? (
            <div className="group relative">
              <button className="flex items-center gap-2 px-2.5 py-1.5 rounded-lg hover:bg-slate-50 transition-all cursor-pointer">
                <div className="h-7 w-7 rounded-lg bg-slate-100 flex items-center justify-center text-slate-500">
                  <User className="h-3.5 w-3.5" />
                </div>
                <div className="hidden sm:block text-left">
                  <p className="text-[11px] font-semibold text-slate-700 leading-none">{user.name}</p>
                  <p className={`text-[10px] font-medium ${theme.color} leading-none mt-0.5`}>{theme.label}</p>
                </div>
              </button>

              <div className="absolute right-0 top-full pt-1.5 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                <div className="w-40 bg-white border border-slate-200 rounded-xl shadow-lg p-1.5">
                  <button
                    onClick={onLogout}
                    className="flex w-full items-center gap-2 px-3 py-2 text-xs text-rose-600 font-semibold hover:bg-rose-50 rounded-lg transition-colors"
                  >
                    <LogOut className="h-3.5 w-3.5" />
                    로그아웃
                  </button>
                </div>
              </div>
            </div>
          ) : (
            <div className="h-7 w-7 rounded-lg bg-slate-100 animate-pulse" />
          )}
        </div>
      </div>
    </header>
  );
}
