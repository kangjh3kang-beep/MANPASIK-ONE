'use client';

import React, { useState } from 'react';
import { Activity, ArrowRight, Eye, EyeOff, ShieldCheck, Stethoscope, FlaskConical, Pill, Users2, FileCode2, Smartphone, Settings, Home } from 'lucide-react';
import { signIn } from 'next-auth/react';
import Link from 'next/link';

const PERSONA_CARDS = [
  { email: 'patient@mmup.io', label: '환자', desc: '건강 데이터 조회 및 리워드', icon: Smartphone, color: 'bg-cyan-500' },
  { email: 'doctor@mmup.io', label: '의사', desc: '임상 데이터 콘솔 접근', icon: Stethoscope, color: 'bg-sky-600' },
  { email: 'researcher@mmup.io', label: '연구원', desc: 'AI 에이전트 허브 & 예측', icon: FlaskConical, color: 'bg-emerald-600' },
  { email: 'pharma@mmup.io', label: '제약사', desc: 'GxP 준수 시스템', icon: Pill, color: 'bg-violet-600' },
  { email: 'partner@mmup.io', label: '파트너', desc: 'HL7/FHIR 파이프라인', icon: Users2, color: 'bg-rose-600' },
  { email: 'dev@mmup.io', label: '개발자', desc: 'API 포털 & 샌드박스', icon: FileCode2, color: 'bg-slate-600' },
  { email: 'admin@mmup.io', label: '관리자', desc: '전체 플랫폼 관리', icon: Settings, color: 'bg-amber-600' },
];

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {

    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const result = await signIn('credentials', {
        email,
        password,
        redirect: false,
      });

      if (result?.error) {
        setError('이메일 또는 비밀번호가 올바르지 않습니다.');
      } else {
        window.location.href = '/ko';
      }
    } catch {
      setError('로그인 처리 중 오류가 발생했습니다.');
    } finally {
      setIsLoading(false);
    }
  };


  const selectPersona = (personaEmail: string) => {
    setEmail(personaEmail);
    setPassword('demo1234');
    setError('');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-sky-50 flex">
      {/* Left Panel - Login Form */}
      <div className="flex flex-1 flex-col justify-center px-8 py-12 lg:px-16">
        <div className="mx-auto w-full max-w-md">
          {/* Logo */}
          <div className="flex items-center justify-between mb-10">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-sky-600 text-white shadow-lg shadow-sky-600/20">
                <Activity className="h-6 w-6" />
              </div>
              <span className="text-2xl font-bold tracking-tight text-slate-900">
                MMUP<span className="text-sky-600">.</span>
              </span>
            </div>
            
            <Link 
              href="/"
              className="flex items-center gap-2 rounded-lg bg-slate-100 px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition-colors"
            >
              <Home className="h-4 w-4" />
              홈으로 가기
            </Link>
          </div>

          {/* Title Section */}
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            플랫폼에 로그인
          </h1>
          <p className="mt-2 text-sm text-slate-500">
            MMUP 통합 의료 생태계에 접속합니다. 로그인 후 역할(페르소나)에 따라 적합한 도메인 허브로 자동 이동합니다.
          </p>

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="mt-8 space-y-5">
            <div>
              <label htmlFor="email" className="block text-sm font-semibold text-slate-700 mb-1.5">
                이메일 주소
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="doctor@mmup.io"
                required
                className="block w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition-all focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20 focus:outline-none"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-semibold text-slate-700 mb-1.5">
                비밀번호
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  className="block w-full rounded-xl border border-slate-300 bg-white px-4 py-3 pr-12 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition-all focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20 focus:outline-none"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
            </div>

            {error && (
              <div className="rounded-xl bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="flex w-full items-center justify-center gap-2 rounded-xl bg-slate-900 px-4 py-3.5 text-sm font-semibold text-white shadow-lg transition-all hover:bg-sky-600 hover:shadow-sky-600/20 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent" />
              ) : (
                <>
                  로그인
                  <ArrowRight className="h-4 w-4" />
                </>
              )}
            </button>
          </form>

          {/* Security Badge */}
          <div className="mt-8 flex items-center gap-2 rounded-xl bg-emerald-50 border border-emerald-100 px-4 py-3">
            <ShieldCheck className="h-5 w-5 text-emerald-600 flex-shrink-0" />
            <p className="text-xs text-emerald-700">
              모든 통신은 TLS 1.3 암호화로 보호됩니다. IEC 62304 Class C 준수.
            </p>
          </div>
        </div>
      </div>

      {/* Right Panel - Persona Quick Select */}
      <div className="hidden lg:flex flex-1 flex-col justify-center bg-slate-900 px-12 py-12 rounded-l-[3rem]">
        <div className="mx-auto w-full max-w-lg">
          <h2 className="text-xl font-bold text-white mb-2">
            🧪 개발 데모 — 페르소나 바로 선택
          </h2>
          <p className="text-sm text-slate-400 mb-8">
            아래 카드를 클릭하면 해당 역할의 이메일과 비밀번호가 자동으로 입력됩니다. 비밀번호는 모두 <code className="bg-slate-800 px-1.5 py-0.5 rounded text-sky-400">demo1234</code> 입니다.
          </p>

          <div className="grid grid-cols-2 gap-3">
            {PERSONA_CARDS.map((persona) => (
              <button
                key={persona.email}
                onClick={() => selectPersona(persona.email)}
                className={`group flex items-start gap-3 rounded-2xl border border-slate-700 bg-slate-800/50 px-4 py-4 text-left transition-all hover:border-sky-500 hover:bg-slate-800 ${
                  email === persona.email ? 'ring-2 ring-sky-500 border-sky-500 bg-slate-800' : ''
                }`}
              >
                <div className={`flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-xl ${persona.color} text-white`}>
                  <persona.icon className="h-5 w-5" />
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-bold text-white">{persona.label}</p>
                  <p className="text-xs text-slate-400 truncate">{persona.desc}</p>
                </div>
              </button>
            ))}
          </div>

          <div className="mt-8 border-t border-slate-700 pt-6">
            <p className="text-xs text-slate-500">
              © 2026 ManPaSik Healthcare Platform. 페르소나 기반 SSO(Single Sign-On) 시스템.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
