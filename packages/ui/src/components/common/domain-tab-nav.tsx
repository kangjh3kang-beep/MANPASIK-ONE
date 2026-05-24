'use client';

import React from 'react';

export interface DomainTab {
  id: string;
  label: string;
  href: string;
  count?: number;
}

interface DomainTabNavProps {
  tabs: DomainTab[];
  currentPath: string;
  localePrefix: string;
}

export function DomainTabNav({ tabs, currentPath, localePrefix }: DomainTabNavProps) {
  return (
    <nav
      aria-label="도메인 하위 메뉴"
      className="border-b border-slate-200/80 bg-white/90 backdrop-blur-sm sticky top-0 z-10"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-0.5 overflow-x-auto scrollbar-hide -mb-px py-1">
          {tabs.map((tab) => {
            const fullHref = `${localePrefix}${tab.href}`;
            const isActive = currentPath === fullHref || currentPath.startsWith(fullHref + '/');
            return (
              <a
                key={tab.id}
                href={fullHref}
                className={`
                  relative flex-shrink-0 px-5 py-2.5 text-[13px] font-semibold rounded-lg
                  transition-all duration-200 whitespace-nowrap
                  ${isActive
                    ? 'bg-slate-900 text-white shadow-sm'
                    : 'text-slate-500 hover:text-slate-900 hover:bg-slate-100'
                  }
                `}
                aria-current={isActive ? 'page' : undefined}
              >
                {tab.label}
                {tab.count !== undefined && (
                  <span className={`ml-1.5 px-1.5 py-0.5 rounded-full text-[10px] font-bold ${
                    isActive ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-500'
                  }`}>
                    {tab.count}
                  </span>
                )}
              </a>
            );
          })}
        </div>
      </div>
    </nav>
  );
}
