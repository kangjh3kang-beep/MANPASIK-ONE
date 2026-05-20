'use client';

import React from 'react';
import Link from 'next/link';
import { Activity, ArrowLeft, Stethoscope, Network, ActivitySquare, Coins, Users2, Pill, FileCode2, Smartphone } from 'lucide-react';

const DOMAINS = [
  { id: 'clinical', label: '임상 콘솔', path: '/domains/clinical', port: 3001, icon: Stethoscope },
  { id: 'agents-hub', label: 'AI 에이전트', path: '/domains/agents-hub', port: 3002, icon: Network },
  { id: 'predictor', label: '생체 예측', path: '/domains/predictor', port: 3003, icon: ActivitySquare },
  { id: 'reward', label: '리워드', path: '/domains/reward', port: 3004, icon: Coins },
  { id: 'partner', label: '파트너', path: '/domains/partner', port: 3005, icon: Users2 },
  { id: 'gxp', label: 'GxP', path: '/domains/gxp', port: 3006, icon: Pill },
  { id: 'dev-portal', label: '개발자', path: '/domains/dev-portal', port: 3007, icon: FileCode2 },
  { id: 'app', label: '건강 앱', path: '/domains/app', port: 3008, icon: Smartphone },
];

function getDomainUrl(domain: typeof DOMAINS[0]) {
  if (typeof window !== 'undefined' && window.location.hostname === 'localhost') {
    // MSA 환경에서는 각 도메인 앱이 루트(/)에서 서빙됩니다.
    return `http://localhost:${domain.port}/`;
  }
  return domain.path;
}



interface DomainNavProps {
  currentDomain: string;
}

export function DomainNav({ currentDomain }: DomainNavProps) {
  return (
    <div className="w-full bg-slate-900 text-white">
      <div className="max-w-7xl mx-auto px-4 flex items-center gap-1 overflow-x-auto h-10">
        <Link
          href="/ko"
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-medium text-slate-300 hover:bg-slate-800 hover:text-white transition-colors flex-shrink-0"
        >
          <ArrowLeft className="h-3 w-3" />
          메인
        </Link>
        <div className="h-4 w-px bg-slate-700 mx-1 flex-shrink-0" />
        
        {DOMAINS.map(domain => {
          const isActive = domain.id === currentDomain;
          return (
            <a
              key={domain.id}
              href={getDomainUrl(domain)}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-medium transition-colors flex-shrink-0 ${
                isActive
                  ? 'bg-white/10 text-white'
                  : 'text-slate-400 hover:bg-slate-800 hover:text-white'
              }`}
            >
              <domain.icon className="h-3 w-3" />
              {domain.label}
            </a>
          );

        })}
      </div>
    </div>
  );
}
