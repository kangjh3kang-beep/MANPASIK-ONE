'use client';
import React, { useState, useEffect } from 'react';
import { Activity, Menu, ChevronRight, Globe, LogIn } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

/**
 * 현재 pathname에서 locale prefix를 추출
 */
function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export function SiteHeader() {
  const [isScrolled, setIsScrolled] = useState(false);
  const localePrefix = useLocalePrefix();

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <header
      className={`fixed top-0 z-50 w-full transition-all duration-500 ${
        isScrolled
          ? 'bg-white/85 backdrop-blur-xl border-b border-slate-200/50 shadow-sm py-2'
          : 'bg-transparent py-4'
      }`}
    >
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-6 lg:px-8">
        {/* Left Side: Logo */}
        <div className="flex items-center gap-12">
          <Link href={localePrefix} className="flex items-center gap-2 hover:opacity-80 transition-opacity">
            <div className="flex h-[34px] w-[34px] items-center justify-center rounded-[10px] bg-sky-600 text-white shadow-md shadow-sky-600/20">
              <Activity className="h-5 w-5" />
            </div>
            <span className="text-[22px] font-bold tracking-tight text-slate-900">
              MMUP<span className="text-sky-600">.</span>
            </span>
          </Link>

          {/* Desktop Navigation Links */}
          <nav className="hidden lg:flex items-center gap-8">
            <Link href={`${localePrefix}/domains/agents-hub`} className="text-[15px] font-semibold text-slate-600 hover:text-sky-600 transition-colors">AI 분석</Link>
            <Link href={`${localePrefix}/domains/hardware-core`} className="text-[15px] font-semibold text-slate-600 hover:text-sky-600 transition-colors">측정 장비</Link>
            <div className="h-4 w-px bg-slate-300"></div>
            <Link href={`${localePrefix}/domains/partner`} className="text-[15px] font-semibold text-slate-600 hover:text-sky-600 transition-colors">병원 연동</Link>
          </nav>


        </div>

        {/* Right Side: Translation / Login Container */}
        <div className="flex items-center gap-4">
          <button className="hidden sm:flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-slate-900 hover:bg-slate-100 rounded-full py-1.5 px-3 transition-colors">
            <Globe className="h-[18px] w-[18px]" />
            KO
          </button>
          
          <div className="hidden sm:block h-4 w-px bg-slate-300 mx-2"></div>
          
          <Link
            href={`${localePrefix}/login`}
            className="group flex items-center gap-2 rounded-full border border-slate-300 bg-white pl-4 pr-3 py-2 text-sm font-semibold text-slate-700 transition-all hover:border-sky-500 hover:text-sky-600"
          >
            <LogIn className="h-4 w-4" />
            로그인
          </Link>

          <Link
            href={`${localePrefix}/domains/clinical`}
            className="group hidden sm:flex items-center gap-2 rounded-full bg-slate-900 pl-5 pr-4 py-2.5 text-sm font-semibold text-white transition-all hover:bg-sky-600 hover:shadow-lg hover:shadow-sky-600/20"
          >
            임상 콘솔
            <ChevronRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
          </Link>

          {/* Mobile Menu Icon */}
          <button className="lg:hidden ml-2 p-2 rounded-xl text-slate-600 hover:bg-slate-100 transition-colors">
            <Menu className="h-6 w-6" />
          </button>
        </div>
      </div>
    </header>
  );
}
