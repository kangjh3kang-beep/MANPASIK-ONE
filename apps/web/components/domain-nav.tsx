'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ArrowLeft, Stethoscope, Network, ActivitySquare, Coins, Users2, Pill, FileCode2, Smartphone, Globe, Brain, Activity, FileText, Video, ShoppingCart, Settings as SettingsIcon, UserCog, ChevronDown } from 'lucide-react';

const NAV_GROUPS = [
  {
    label: '건강 관리',
    items: [
      { id: 'measure', label: '측정', path: '/domains/measure', icon: Activity },
      { id: 'health-records', label: '기록', path: '/domains/health-records', icon: FileText },
      { id: 'clinical', label: '건강 데이터', path: '/domains/clinical', icon: Stethoscope },
    ],
  },
  {
    label: 'AI · 예측',
    items: [
      { id: 'agents-hub', label: 'AI 분석', path: '/domains/agents-hub', icon: Network },
      { id: 'predictor', label: '위험 예측', path: '/domains/predictor', icon: ActivitySquare },
      { id: 'ai-coach', label: 'AI 코치', path: '/domains/ai-coach', icon: Brain },
    ],
  },
  {
    label: '서비스',
    items: [
      { id: 'telemedicine', label: '화상 진료', path: '/domains/telemedicine', icon: Video },
      { id: 'store', label: '스토어', path: '/domains/store', icon: ShoppingCart },
      { id: 'reward', label: '보상', path: '/domains/reward', icon: Coins },
      { id: 'app', label: '건강 앱', path: '/domains/app', icon: Smartphone },
    ],
  },
  {
    label: '운영',
    items: [
      { id: 'partner', label: '병원 연동', path: '/domains/partner', icon: Users2 },
      { id: 'gxp', label: '품질 관리', path: '/domains/gxp', icon: Pill },
      { id: 'dev-portal', label: '개발자', path: '/domains/dev-portal', icon: FileCode2 },
      { id: 'hardware-core', label: '장비', path: '/domains/hardware-core', icon: Globe },
      { id: 'admin', label: '관리자', path: '/domains/admin', icon: UserCog },
      { id: 'settings', label: '설정', path: '/domains/settings', icon: SettingsIcon },
    ],
  },
];

interface DomainNavProps {
  currentDomain: string;
}

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export function DomainNav({ currentDomain }: DomainNavProps) {
  const localePrefix = useLocalePrefix();

  const currentGroup = NAV_GROUPS.find(g => g.items.some(i => i.id === currentDomain));

  return (
    <nav className="w-full bg-slate-900" aria-label="도메인 네비게이션">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center h-11 gap-0.5 overflow-x-auto scrollbar-hide">
          {/* 홈 */}
          <Link
            href={localePrefix}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium text-slate-400 hover:text-white hover:bg-white/10 transition-all flex-shrink-0"
          >
            <ArrowLeft className="h-3.5 w-3.5" />
            <span className="hidden sm:inline">메인</span>
          </Link>

          <div className="h-5 w-px bg-slate-700 mx-1 flex-shrink-0" />

          {/* 카테고리 그룹 */}
          {NAV_GROUPS.map((group, gi) => (
            <React.Fragment key={group.label}>
              {gi > 0 && <div className="h-5 w-px bg-slate-700/50 mx-0.5 flex-shrink-0" />}

              {/* 그룹 라벨 */}
              <span className="text-[9px] font-bold text-slate-600 uppercase tracking-wider px-2 flex-shrink-0 hidden lg:block">
                {group.label}
              </span>

              {/* 그룹 아이템 */}
              {group.items.map(item => {
                const isActive = item.id === currentDomain;
                return (
                  <Link
                    key={item.id}
                    href={`${localePrefix}${item.path}`}
                    className={`
                      flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs font-medium
                      transition-all flex-shrink-0
                      ${isActive
                        ? 'bg-white text-slate-900 shadow-sm'
                        : 'text-slate-400 hover:text-white hover:bg-white/10'
                      }
                    `}
                    aria-current={isActive ? 'page' : undefined}
                  >
                    <item.icon className="h-3.5 w-3.5" />
                    {item.label}
                  </Link>
                );
              })}
            </React.Fragment>
          ))}
        </div>
      </div>
    </nav>
  );
}
