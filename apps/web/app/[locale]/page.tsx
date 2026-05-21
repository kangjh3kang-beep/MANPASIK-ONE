'use client';
export const runtime = 'edge';
import React from 'react';
import { SiteHeader } from '../../components/site-header';
import { DomainHeader } from '@mmup/ui';
import { DomainBentoGrid } from '../../components/domain-bento-grid';
import { EcosystemNerveCenter } from '../../components/dashboard/ecosystem-nerve-center';
import { LayoutGrid, Database, Activity } from 'lucide-react';
import Link from 'next/link';

export default function UnifiedDashboard() {
  return (
    <div className="min-h-screen bg-slate-50">
      <SiteHeader />
      
      <main className="pt-24 pb-20">
        <div className="max-w-7xl mx-auto px-6 lg:px-8 space-y-12">
          {/* Hero Section / Nerve Center */}
          <div className="relative rounded-3xl overflow-hidden bg-slate-900 text-white p-12 shadow-2xl">
            <div className="absolute inset-0 bg-gradient-to-br from-sky-900/50 via-slate-900 to-slate-900 opacity-90" />
            <div className="relative z-10 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <div>
                <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-sky-500/20 text-sky-400 text-xs font-bold mb-6 border border-sky-500/30">
                  <Activity className="h-3 w-3" />
                  SYSTEM OVERWATCH ACTIVE
                </span>
                <h1 className="text-4xl lg:text-5xl font-black tracking-tight mb-6 leading-tight">
                  만파식 생태계<br />
                  <span className="text-sky-400 font-extrabold">통합 운영 대시보드</span>
                </h1>
                <p className="text-lg text-slate-300 mb-8 max-w-lg leading-relaxed">
                  9대 핵심 도메인과 AI 에이전트 허브를 실시간으로 모니터링하고 제어합니다.
                  정밀 의학 데이터의 전주기 관리를 한눈에 확인하세요.
                </p>
                <div className="flex flex-wrap gap-4">
                  <Link 
                    href="/domains/clinical"
                    className="px-6 py-3.5 bg-sky-600 hover:bg-sky-500 text-white font-bold rounded-2xl transition-all shadow-lg shadow-sky-600/30 flex items-center gap-2"
                  >
                    데모 콘솔 열기
                  </Link>
                  <Link 
                    href="/domains/app"
                    className="px-6 py-3.5 bg-white/10 hover:bg-white/20 text-white font-bold rounded-2xl transition-all border border-white/20 backdrop-blur-sm"
                  >
                    건강 앱 관리
                  </Link>
                </div>

              </div>
              
              <div className="hidden lg:block">
                <EcosystemNerveCenter />
              </div>
            </div>
          </div>

          {/* Infrastructure Nodes */}
          <div>
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-2xl font-bold text-slate-900">도메인 인프라 노드</h2>
                <p className="text-slate-500">생태계의 각 서비스를 독립적으로 관리하고 상태를 확인합니다.</p>
              </div>
              <div className="flex gap-2">
                <div className="flex items-center gap-2 px-3 py-1 bg-white border border-slate-200 rounded-lg text-xs font-semibold text-slate-600">
                  <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  Live Nodes: 9
                </div>
              </div>
            </div>
            
            <DomainBentoGrid />
          </div>
        </div>
      </main>
    </div>
  );
}
