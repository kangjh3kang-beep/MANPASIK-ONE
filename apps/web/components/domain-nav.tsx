'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ArrowLeft, Stethoscope, Network, ActivitySquare, Coins, Users2, Pill, FileCode2, Smartphone, Globe, Brain, Activity, FileText, Video, ShoppingCart, Settings as SettingsIcon } from 'lucide-react';

const DOMAINS = [
  { id: 'measure', label: '측정', path: '/domains/measure', icon: Activity },
  { id: 'health-records', label: '기록', path: '/domains/health-records', icon: FileText },
  { id: 'clinical', label: '건강 데이터', path: '/domains/clinical', icon: Stethoscope },
  { id: 'agents-hub', label: 'AI 분석', path: '/domains/agents-hub', icon: Network },
  { id: 'predictor', label: '위험 예측', path: '/domains/predictor', icon: ActivitySquare },
  { id: 'ai-coach', label: 'AI 코치', path: '/domains/ai-coach', icon: Brain },
  { id: 'reward', label: '보상', path: '/domains/reward', icon: Coins },
  { id: 'partner', label: '병원 연동', path: '/domains/partner', icon: Users2 },
  { id: 'gxp', label: '품질 관리', path: '/domains/gxp', icon: Pill },
  { id: 'dev-portal', label: '개발자', path: '/domains/dev-portal', icon: FileCode2 },
  { id: 'hardware-core', label: '측정 장비', path: '/domains/hardware-core', icon: Globe },
  { id: 'telemedicine', label: '진료', path: '/domains/telemedicine', icon: Video },
  { id: 'store', label: '스토어', path: '/domains/store', icon: ShoppingCart },
  { id: 'app', label: '건강 앱', path: '/domains/app', icon: Smartphone },
  { id: 'settings', label: '설정', path: '/domains/settings', icon: SettingsIcon },
];

interface DomainNavProps {
  currentDomain: string;
}

/**
 * 현재 pathname에서 locale prefix를 추출하여 링크에 적용.
 * 예: /ko/domains/clinical → locale = 'ko'
 */
function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export function DomainNav({ currentDomain }: DomainNavProps) {
  const localePrefix = useLocalePrefix();

  return (
    <div className="w-full bg-slate-900 text-white">
      <div className="max-w-7xl mx-auto px-4 flex items-center gap-1 overflow-x-auto h-10">
        <Link
          href={localePrefix}
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
              href={`${localePrefix}${domain.path}`}
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
