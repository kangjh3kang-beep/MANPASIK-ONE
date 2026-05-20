'use client';

import React from 'react';
import Link from 'next/link';
import { ArrowLeft, Stethoscope, Network, ActivitySquare, Coins, Users2, Pill, FileCode2, Smartphone } from 'lucide-react';

const DOMAINS = [
  { id: 'clinical', label: '임상 콘솔', path: '/domains/clinical', icon: Stethoscope },
  { id: 'agents-hub', label: 'AI 에이전트', path: '/domains/agents-hub', icon: Network },
  { id: 'predictor', label: '생체 예측', path: '/domains/predictor', icon: ActivitySquare },
  { id: 'reward', label: '리워드', path: '/domains/reward', icon: Coins },
  { id: 'partner', label: '파트너', path: '/domains/partner', icon: Users2 },
  { id: 'gxp', label: 'GxP', path: '/domains/gxp', icon: Pill },
  { id: 'dev-portal', label: '개발자', path: '/domains/dev-portal', icon: FileCode2 },
  { id: 'app', label: '건강 앱', path: '/domains/app', icon: Smartphone },
];

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
            <Link
              key={domain.id}
              href={domain.path}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-medium transition-colors flex-shrink-0 ${
                isActive
                  ? 'bg-white/10 text-white'
                  : 'text-slate-400 hover:bg-slate-800 hover:text-white'
              }`}
            >
              <domain.icon className="h-3 w-3" />
              {domain.label}
            </Link>
          );
        })}
      </div>
    </div>
  );
}
