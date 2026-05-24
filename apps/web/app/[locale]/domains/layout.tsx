'use client';

import React from 'react';
import { DomainNav } from '@/components/domain-nav';
import { usePathname } from 'next/navigation';

export default function DomainsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const currentDomain = pathname.split('/domains/')[1]?.split('/')[0] || '';

  return (
    <>
      <a href="#main-content" className="sr-only focus:not-sr-only focus:absolute focus:z-[100] focus:top-2 focus:left-2 focus:px-4 focus:py-2 focus:bg-sky-600 focus:text-white focus:rounded-lg focus:text-sm focus:font-bold">
        본문으로 건너뛰기
      </a>
      <DomainNav currentDomain={currentDomain} />
      <div id="main-content">
        {children}
      </div>
    </>
  );
}
