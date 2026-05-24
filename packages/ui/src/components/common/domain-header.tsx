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
  /** When true, shows a certification shield badge next to the title */
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
  
  // 페르소나별 테마 설정
  const personaThemes: Record<string, { color: string, label: string }> = {
    doctor: { color: 'bg-sky-500', label: '의사' },
    researcher: { color: 'bg-emerald-500', label: '연구원' },
    pharma: { color: 'bg-violet-500', label: '제약 전문가' },
    admin: { color: 'bg-amber-500', label: '관리자' },
    patient: { color: 'bg-rose-500', label: '환자' },
  };

  const theme = user?.persona ? personaThemes[user.persona] : { color: 'bg-slate-500', label: '사용자' };

  return (
    <header className="sticky top-0 z-40 w-full border-b border-slate-200 bg-white/80 backdrop-blur-md">
      <div className="flex h-20 items-center justify-between px-8">
        {/* Left: Branding & Page Title */}
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-3 border-r border-slate-200 pr-6 mr-1 hover:opacity-80 transition-opacity">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-900 text-white">
              <Activity className="h-6 w-6" />
            </div>
            <span className="text-xl font-black tracking-tighter text-slate-900">
              MPS<span className="text-sky-500">.</span>
            </span>
          </Link>
          
          <div className="flex items-center gap-3">
            <div className={`flex h-9 w-9 items-center justify-center rounded-lg ${theme.color} text-white shadow-sm`}>
              <Icon className="h-5 w-5" />
            </div>
            <div>
              <h1 className="text-lg font-bold text-slate-900 leading-none mb-1">
                {title}
                {certified && (
                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-trust-blue/10 text-trust-blue text-[10px] font-semibold ml-2">
                    <svg className="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}><path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" /></svg>
                    인증
                  </span>
                )}
              </h1>
              {subtitle && <p className="text-xs text-slate-400 font-medium tracking-tight">{subtitle}</p>}
            </div>
          </div>
        </div>

        {/* Center: Search (Optional) */}
        <div className="hidden md:flex flex-1 max-w-md px-12">
          <div className="relative w-full">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
            <input 
              type="text" 
              placeholder="환자 번호, 처방 ID 검색..." 
              className="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/10 focus:border-sky-500 transition-all"
            />
          </div>
        </div>

        {/* Right: User Profile & Actions */}
        <div className="flex items-center gap-4">
          <button className="relative p-2 text-slate-400 hover:text-slate-600 transition-colors">
            <Bell className="h-5 w-5" />
            <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-rose-500 border-2 border-white" />
          </button>
          
          <div className="h-8 w-px bg-slate-200 mx-1" />

          {user ? (
            <div className="flex items-center gap-4">
              <div className="flex flex-col items-end">
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-0.5 rounded text-[10px] font-bold text-white uppercase ${theme.color}`}>
                    {theme.label}
                  </span>
                  <span className="text-sm font-bold text-slate-900">{user.name}</span>
                </div>
                {user.organization && <span className="text-[10px] text-slate-400 font-medium">{user.organization}</span>}
              </div>
              
              <div className="group relative cursor-pointer">
                <div className="h-10 w-10 rounded-xl bg-slate-100 border border-slate-200 flex items-center justify-center text-slate-500 hover:bg-slate-200 transition-all overflow-hidden">
                  <User className="h-6 w-6" />
                </div>
                
                {/* Dropdown simulation */}
                <div className="absolute right-0 top-full pt-2 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
                  <div className="w-48 bg-white border border-slate-200 rounded-2xl shadow-xl p-2">
                    <button 
                      onClick={onLogout}
                      className="flex w-full items-center gap-2 px-3 py-2 text-sm text-rose-600 font-bold hover:bg-rose-50 rounded-lg transition-colors"
                    >
                      <LogOut className="h-4 w-4" />
                      로그아웃
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="h-10 w-10 rounded-xl bg-slate-100 animate-pulse" />
          )}
        </div>
      </div>
    </header>
  );
}
