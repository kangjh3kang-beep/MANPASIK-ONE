'use client';

import React, { useState } from 'react';
import {
  LayoutDashboard, Users, Activity, Cpu, FileText,
  Settings, ChevronLeft, ChevronRight, Stethoscope,
  Bell, Search, LogOut
} from 'lucide-react';
import Link from 'next/link';

const NAV_ITEMS = [
  { id: 'dashboard', label: '대시보드', icon: LayoutDashboard, href: '#' },
  { id: 'patients', label: '환자 관리', icon: Users, href: '#' },
  { id: 'vitals', label: '생체 신호', icon: Activity, href: '#' },
  { id: 'sensors', label: '센서 디바이스', icon: Cpu, href: '#' },
  { id: 'reports', label: '진단 리포트', icon: FileText, href: '#' },
  { id: 'settings', label: '설정', icon: Settings, href: '#' },
];

type ClinicalShellProps = {
  children?: React.ReactNode;
};

export function ClinicalShell({ children }: ClinicalShellProps) {
  const [collapsed, setCollapsed] = useState(false);
  const [activeId, setActiveId] = useState('dashboard');

  return (
    <div className="flex h-screen bg-slate-50">
      {/* Sidebar */}
      <aside
        className={`flex flex-col border-r border-slate-200 bg-white transition-all duration-300 ${
          collapsed ? 'w-[68px]' : 'w-[240px]'
        }`}
      >
        {/* Logo */}
        <div className="flex h-16 items-center justify-between border-b border-slate-100 px-4">
          {!collapsed && (
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-sky-600 text-white">
                <Stethoscope className="h-4 w-4" />
              </div>
              <span className="text-sm font-bold text-slate-900">Clinical Console</span>
            </div>
          )}
          {collapsed && (
            <div className="mx-auto flex h-8 w-8 items-center justify-center rounded-lg bg-sky-600 text-white">
              <Stethoscope className="h-4 w-4" />
            </div>
          )}
        </div>

        {/* Navigation */}
        <nav className="flex-1 space-y-1 p-3">
          {NAV_ITEMS.map((item) => (
            <button
              key={item.id}
              onClick={() => setActiveId(item.id)}
              className={`flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors ${
                activeId === item.id
                  ? 'bg-sky-50 text-sky-700'
                  : 'text-slate-500 hover:bg-slate-50 hover:text-slate-900'
              }`}
            >
              <item.icon className="h-[18px] w-[18px] flex-shrink-0" />
              {!collapsed && <span>{item.label}</span>}
            </button>
          ))}
        </nav>

        {/* Collapse Toggle */}
        <div className="border-t border-slate-100 p-3">
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="flex w-full items-center justify-center gap-2 rounded-xl px-3 py-2 text-sm text-slate-400 hover:bg-slate-50 hover:text-slate-600"
          >
            {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
            {!collapsed && <span>접기</span>}
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Top Bar */}
        <header className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-6">
          <div className="flex items-center gap-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <input
                type="text"
                placeholder="환자 검색 (이름, ID, 질환코드)..."
                className="w-80 rounded-xl border border-slate-200 bg-slate-50 py-2 pl-10 pr-4 text-sm text-slate-900 placeholder:text-slate-400 focus:border-sky-500 focus:outline-none focus:ring-1 focus:ring-sky-500"
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button className="relative rounded-xl p-2 text-slate-400 hover:bg-slate-50 hover:text-slate-600">
              <Bell className="h-5 w-5" />
              <span className="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-red-500" />
            </button>
            <div className="h-6 w-px bg-slate-200" />
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-sky-100 text-sky-700 text-xs font-bold">
                박
              </div>
              <div className="hidden sm:block">
                <p className="text-sm font-semibold text-slate-900">박의사</p>
                <p className="text-[11px] text-slate-400">KR-MD-2024-00123</p>
              </div>
            </div>
            <Link href="http://localhost:3000/ko" className="rounded-xl p-2 text-slate-400 hover:bg-slate-50 hover:text-slate-600">
              <LogOut className="h-4 w-4" />
            </Link>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
