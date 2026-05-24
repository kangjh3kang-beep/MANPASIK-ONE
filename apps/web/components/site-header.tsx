'use client';
import React, { useState, useEffect, useCallback } from 'react';
import { Activity, Menu, X, ChevronDown, Globe, LogIn, Moon, Sun, Stethoscope, Network, BarChart3, Coins, Users2, Pill, FileCode2, Cpu, Smartphone, Brain } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

function useCurrentLocale(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? match[1] : 'ko';
}

const DOMAINS = [
  { id: 'measure', label: '건강 측정', icon: Activity, path: '/domains/measure' },
  { id: 'health-records', label: '건강 기록', icon: BarChart3, path: '/domains/health-records' },
  { id: 'clinical', label: '건강 데이터', icon: Stethoscope, path: '/domains/clinical' },
  { id: 'ai-coach', label: 'AI 코치', icon: Brain, path: '/domains/ai-coach' },
  { id: 'agents-hub', label: 'AI 분석', icon: Network, path: '/domains/agents-hub' },
  { id: 'predictor', label: '위험 예측', icon: BarChart3, path: '/domains/predictor' },
  { id: 'telemedicine', label: '화상진료', icon: Smartphone, path: '/domains/telemedicine' },
  { id: 'store', label: '스토어', icon: Coins, path: '/domains/store' },
  { id: 'reward', label: '보상', icon: Coins, path: '/domains/reward' },
  { id: 'partner', label: '병원 연동', icon: Users2, path: '/domains/partner' },
  { id: 'gxp', label: '품질 관리', icon: Pill, path: '/domains/gxp' },
  { id: 'dev-portal', label: '개발자 도구', icon: FileCode2, path: '/domains/dev-portal' },
  { id: 'hardware-core', label: '측정 장비', icon: Cpu, path: '/domains/hardware-core' },
  { id: 'app', label: '건강 앱', icon: Smartphone, path: '/domains/app' },
  { id: 'admin', label: '관리자', icon: Cpu, path: '/domains/admin' },
  { id: 'settings', label: '설정', icon: Pill, path: '/domains/settings' },
];

const LOCALES = [
  { code: 'ko', label: '한국어', flag: '🇰🇷' },
  { code: 'en', label: 'English', flag: '🇺🇸' },
  { code: 'ja', label: '日本語', flag: '🇯🇵' },
  { code: 'zh', label: '中文', flag: '🇨🇳' },
];

export function SiteHeader() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [servicesOpen, setServicesOpen] = useState(false);
  const [langOpen, setLangOpen] = useState(false);
  const [isDark, setIsDark] = useState(false);
  const localePrefix = useLocalePrefix();
  const currentLocale = useCurrentLocale();

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    const saved = localStorage.getItem('mps-theme');
    if (saved === 'dark') {
      setIsDark(true);
      document.documentElement.setAttribute('data-theme', 'dark');
    }
  }, []);

  const toggleTheme = useCallback(() => {
    setIsDark(prev => {
      const next = !prev;
      document.documentElement.setAttribute('data-theme', next ? 'dark' : 'light');
      localStorage.setItem('mps-theme', next ? 'dark' : 'light');
      return next;
    });
  }, []);

  return (
    <>
      <header
        className={`fixed top-0 z-50 w-full transition-all duration-500 ${
          isScrolled
            ? 'backdrop-blur-xl border-b shadow-sm py-1.5'
            : 'py-3'
        }`}
        style={{
          backgroundColor: isScrolled ? 'var(--mps-header-bg)' : 'transparent',
          borderColor: 'var(--mps-border)',
        }}
      >
        <div className="mx-auto flex h-12 max-w-7xl items-center justify-between px-5 lg:px-8">
          {/* Logo */}
          <div className="flex items-center gap-8">
            <Link href={localePrefix} className="flex items-center gap-2.5 hover:opacity-80 transition-opacity">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-sky-600 text-white shadow-md shadow-sky-600/20">
                <Activity className="h-4.5 w-4.5" />
              </div>
              <span className="text-xl font-extrabold tracking-tight" style={{ color: 'var(--mps-text)' }}>
                MPS<span className="text-sky-600">.</span>
              </span>
            </Link>

            {/* Desktop Nav */}
            <nav className="hidden lg:flex items-center gap-1">
              {/* 서비스 드롭다운 */}
              <div className="relative">
                <button
                  onClick={() => setServicesOpen(!servicesOpen)}
                  onBlur={() => setTimeout(() => setServicesOpen(false), 200)}
                  className="flex items-center gap-1 px-3 py-2 rounded-lg text-sm font-semibold transition-colors hover:bg-black/5"
                  style={{ color: 'var(--mps-text-secondary)' }}
                >
                  전체 서비스
                  <ChevronDown className={`h-3.5 w-3.5 transition-transform ${servicesOpen ? 'rotate-180' : ''}`} />
                </button>

                {servicesOpen && (
                  <div
                    className="absolute left-0 top-full mt-2 w-80 rounded-2xl p-3 shadow-xl border z-50"
                    style={{ backgroundColor: 'var(--mps-bg-card)', borderColor: 'var(--mps-border)' }}
                  >
                    <div className="grid grid-cols-1 gap-0.5">
                      {DOMAINS.map(d => (
                        <Link
                          key={d.id}
                          href={`${localePrefix}${d.path}`}
                          className="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors hover:bg-sky-50 dark:hover:bg-sky-950 group"
                          onClick={() => setServicesOpen(false)}
                        >
                          <div className="h-8 w-8 rounded-lg flex items-center justify-center text-sky-600 group-hover:bg-sky-100" style={{ backgroundColor: 'var(--mps-accent-light)' }}>
                            <d.icon className="h-4 w-4" />
                          </div>
                          <span className="text-sm font-medium" style={{ color: 'var(--mps-text)' }}>{d.label}</span>
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* 주요 링크 */}
              <Link href={`${localePrefix}/domains/clinical`} className="px-3 py-2 rounded-lg text-sm font-semibold transition-colors hover:bg-black/5" style={{ color: 'var(--mps-text-secondary)' }}>
                건강 데이터
              </Link>
              <Link href={`${localePrefix}/domains/agents-hub`} className="px-3 py-2 rounded-lg text-sm font-semibold transition-colors hover:bg-black/5" style={{ color: 'var(--mps-text-secondary)' }}>
                AI 분석
              </Link>
              <Link href={`${localePrefix}/domains/partner`} className="px-3 py-2 rounded-lg text-sm font-semibold transition-colors hover:bg-black/5" style={{ color: 'var(--mps-text-secondary)' }}>
                병원 연동
              </Link>
            </nav>
          </div>

          {/* Right Controls */}
          <div className="flex items-center gap-2">
            {/* 테마 토글 */}
            <button
              onClick={toggleTheme}
              className="hidden sm:flex h-9 w-9 items-center justify-center rounded-lg transition-colors hover:bg-black/5"
              aria-label={isDark ? '라이트 모드' : '다크 모드'}
              style={{ color: 'var(--mps-text-secondary)' }}
            >
              {isDark ? <Sun className="h-[18px] w-[18px]" /> : <Moon className="h-[18px] w-[18px]" />}
            </button>

            {/* 언어 선택 */}
            <div className="relative hidden sm:block">
              <button
                onClick={() => setLangOpen(!langOpen)}
                onBlur={() => setTimeout(() => setLangOpen(false), 200)}
                className="flex items-center gap-1.5 px-2.5 py-2 rounded-lg text-sm font-semibold transition-colors hover:bg-black/5"
                style={{ color: 'var(--mps-text-secondary)' }}
              >
                <Globe className="h-4 w-4" />
                {currentLocale.toUpperCase()}
              </button>
              {langOpen && (
                <div
                  className="absolute right-0 top-full mt-2 w-40 rounded-xl p-1.5 shadow-xl border z-50"
                  style={{ backgroundColor: 'var(--mps-bg-card)', borderColor: 'var(--mps-border)' }}
                >
                  {LOCALES.map(l => (
                    <Link
                      key={l.code}
                      href={`/${l.code}`}
                      className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors ${
                        currentLocale === l.code ? 'font-bold' : ''
                      }`}
                      style={{ color: currentLocale === l.code ? 'var(--mps-accent)' : 'var(--mps-text)' }}
                      onClick={() => setLangOpen(false)}
                    >
                      <span>{l.flag}</span>
                      {l.label}
                    </Link>
                  ))}
                </div>
              )}
            </div>

            <div className="hidden sm:block h-5 w-px mx-1" style={{ backgroundColor: 'var(--mps-border)' }} />

            <Link
              href={`${localePrefix}/login`}
              className="hidden sm:flex items-center gap-2 rounded-lg border px-3.5 py-2 text-sm font-semibold transition-all hover:border-sky-500 hover:text-sky-600"
              style={{ borderColor: 'var(--mps-border)', color: 'var(--mps-text)' }}
            >
              <LogIn className="h-4 w-4" />
              로그인
            </Link>

            <Link
              href={`${localePrefix}/domains/clinical`}
              className="hidden md:flex items-center gap-1.5 rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-sky-700 shadow-sm"
            >
              건강 데이터 보기
            </Link>

            {/* 모바일 메뉴 */}
            <button
              onClick={() => setMobileOpen(!mobileOpen)}
              className="lg:hidden p-2 rounded-lg transition-colors hover:bg-black/5"
              style={{ color: 'var(--mps-text)' }}
            >
              {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </header>

      {/* 모바일 메뉴 패널 */}
      {mobileOpen && (
        <div className="fixed inset-0 z-40 lg:hidden" onClick={() => setMobileOpen(false)}>
          <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" />
          <div
            className="absolute right-0 top-0 h-full w-80 max-w-[85vw] p-5 pt-20 overflow-y-auto shadow-2xl"
            style={{ backgroundColor: 'var(--mps-bg-card)' }}
            onClick={e => e.stopPropagation()}
          >
            <p className="text-xs font-bold uppercase tracking-wider mb-3" style={{ color: 'var(--mps-text-muted)' }}>
              전체 서비스
            </p>
            <div className="space-y-1 mb-6">
              {DOMAINS.map(d => (
                <Link
                  key={d.id}
                  href={`${localePrefix}${d.path}`}
                  className="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors"
                  onClick={() => setMobileOpen(false)}
                  style={{ color: 'var(--mps-text)' }}
                >
                  <d.icon className="h-4 w-4 text-sky-600" />
                  <span className="text-sm font-medium">{d.label}</span>
                </Link>
              ))}
            </div>

            <div className="border-t pt-4 space-y-2" style={{ borderColor: 'var(--mps-border)' }}>
              <button onClick={toggleTheme} className="flex items-center gap-3 px-3 py-2.5 w-full rounded-xl text-sm" style={{ color: 'var(--mps-text)' }}>
                {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                {isDark ? '밝은 화면' : '어두운 화면'}
              </button>
              <div className="flex gap-2 px-3">
                {LOCALES.map(l => (
                  <Link key={l.code} href={`/${l.code}`} onClick={() => setMobileOpen(false)}
                    className={`px-3 py-1.5 rounded-lg text-xs font-semibold border ${currentLocale === l.code ? 'border-sky-500 text-sky-600' : ''}`}
                    style={{ borderColor: currentLocale === l.code ? undefined : 'var(--mps-border)', color: currentLocale === l.code ? undefined : 'var(--mps-text-secondary)' }}
                  >
                    {l.flag} {l.code.toUpperCase()}
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
