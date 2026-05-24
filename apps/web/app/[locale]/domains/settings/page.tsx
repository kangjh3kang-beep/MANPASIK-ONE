'use client';

import React, { useState, useEffect } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Settings, Moon, Sun, Globe, Eye, Volume2, Bell, Shield, Download, Trash2, ChevronRight } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export default function SettingsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [isDark, setIsDark] = useState(false);
  const [fontSize, setFontSize] = useState<'normal' | 'large' | 'elder'>('normal');
  const [notifications, setNotifications] = useState({ measurement: true, health: true, marketing: false });

  useEffect(() => {
    setIsDark(document.documentElement.getAttribute('data-theme') === 'dark');
  }, []);

  const toggleTheme = () => {
    const next = !isDark;
    setIsDark(next);
    document.documentElement.setAttribute('data-theme', next ? 'dark' : 'light');
    localStorage.setItem('mps-theme', next ? 'dark' : 'light');
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="설정">
      <DomainHeader
        title="설정"
        subtitle="화면, 알림, 개인정보를 관리합니다"
        icon={Settings}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-2xl mx-auto px-6 py-8 space-y-6">
        {/* 화면 설정 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Eye className="h-4 w-4 text-sky-600" /> 화면 설정
          </h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-700">화면 모드</p>
                <p className="text-xs text-slate-400">밝은 화면 또는 어두운 화면을 선택합니다</p>
              </div>
              <button onClick={toggleTheme} className="flex items-center gap-2 px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold hover:bg-slate-50 transition-colors"
                style={{ color: 'var(--mps-text)' }}>
                {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                {isDark ? '밝은 화면' : '어두운 화면'}
              </button>
            </div>

            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-700">글자 크기</p>
                <p className="text-xs text-slate-400">읽기 편한 크기를 선택하세요</p>
              </div>
              <div className="flex gap-1">
                {[
                  { id: 'normal' as const, label: '보통' },
                  { id: 'large' as const, label: '크게' },
                  { id: 'elder' as const, label: '아주 크게' },
                ].map(s => (
                  <button key={s.id} onClick={() => setFontSize(s.id)}
                    className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors ${
                      fontSize === s.id ? 'bg-sky-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}>
                    {s.label}
                  </button>
                ))}
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-700">언어</p>
                <p className="text-xs text-slate-400">인터페이스 언어를 변경합니다</p>
              </div>
              <div className="flex gap-1">
                {[
                  { code: 'ko', label: '한국어' },
                  { code: 'en', label: 'English' },
                  { code: 'ja', label: '日本語' },
                  { code: 'zh', label: '中文' },
                ].map(l => (
                  <Link key={l.code} href={`/${l.code}/domains/settings`}
                    className="px-3 py-1.5 rounded-lg text-xs font-semibold bg-slate-100 text-slate-600 hover:bg-slate-200 transition-colors">
                    {l.label}
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* 알림 설정 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Bell className="h-4 w-4 text-sky-600" /> 알림 설정
          </h2>
          <div className="space-y-4">
            {[
              { key: 'measurement' as const, label: '측정 알림', desc: '측정 시간, 결과 알림을 받습니다' },
              { key: 'health' as const, label: '건강 알림', desc: '이상 수치 감지 시 즉시 알립니다' },
              { key: 'marketing' as const, label: '혜택 알림', desc: '할인, 이벤트 정보를 받습니다' },
            ].map(n => (
              <label key={n.key} className="flex items-center justify-between cursor-pointer">
                <div>
                  <p className="text-sm font-medium text-slate-700">{n.label}</p>
                  <p className="text-xs text-slate-400">{n.desc}</p>
                </div>
                <div className={`relative w-11 h-6 rounded-full transition-colors ${notifications[n.key] ? 'bg-sky-600' : 'bg-slate-300'}`}
                  onClick={() => setNotifications(prev => ({ ...prev, [n.key]: !prev[n.key] }))}>
                  <div className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${notifications[n.key] ? 'translate-x-[22px]' : 'translate-x-0.5'}`} />
                </div>
              </label>
            ))}
          </div>
        </section>

        {/* 개인정보 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Shield className="h-4 w-4 text-sky-600" /> 개인정보 관리
          </h2>
          <div className="space-y-2">
            <button className="w-full flex items-center justify-between p-3 rounded-xl hover:bg-slate-50 transition-colors">
              <div className="flex items-center gap-3">
                <Download className="h-4 w-4 text-slate-400" />
                <span className="text-sm font-medium text-slate-700">내 데이터 내보내기</span>
              </div>
              <ChevronRight className="h-4 w-4 text-slate-300" />
            </button>
            <button className="w-full flex items-center justify-between p-3 rounded-xl hover:bg-slate-50 transition-colors">
              <div className="flex items-center gap-3">
                <Shield className="h-4 w-4 text-slate-400" />
                <span className="text-sm font-medium text-slate-700">데이터 공유 동의 관리</span>
              </div>
              <ChevronRight className="h-4 w-4 text-slate-300" />
            </button>
            <button className="w-full flex items-center justify-between p-3 rounded-xl hover:bg-red-50 transition-colors group">
              <div className="flex items-center gap-3">
                <Trash2 className="h-4 w-4 text-red-400" />
                <span className="text-sm font-medium text-red-600">계정 삭제</span>
              </div>
              <ChevronRight className="h-4 w-4 text-red-300" />
            </button>
          </div>
        </section>

        <p className="text-center text-xs text-slate-400">MPS 만파식 플랫폼 v2.0 · 개인정보처리방침 · 이용약관</p>
      </div>
    </main>
  );
}
