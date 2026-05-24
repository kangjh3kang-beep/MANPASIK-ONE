'use client';

import React from 'react';

export interface DomainTab {
  id: string;
  label: string;
  href: string;
}

interface DomainTabNavProps {
  tabs: DomainTab[];
  currentPath: string;
  localePrefix: string;
}

export function DomainTabNav({ tabs, currentPath, localePrefix }: DomainTabNavProps) {
  return (
    <nav aria-label="도메인 하위 메뉴" className="border-b border-slate-200 bg-white sticky top-0 z-10">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="flex items-center gap-1 overflow-x-auto scrollbar-hide -mb-px">
          {tabs.map((tab) => {
            const fullHref = `${localePrefix}${tab.href}`;
            const isActive = currentPath === fullHref || currentPath.startsWith(fullHref + '/');
            return (
              <a
                key={tab.id}
                href={fullHref}
                className={`flex-shrink-0 px-4 py-3 text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${
                  isActive
                    ? 'border-sky-600 text-sky-700'
                    : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
                }`}
                aria-current={isActive ? 'page' : undefined}
              >
                {tab.label}
              </a>
            );
          })}
        </div>
      </div>
    </nav>
  );
}
