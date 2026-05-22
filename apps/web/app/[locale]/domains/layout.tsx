'use client';

import React from 'react';
import { DomainNav } from '@/components/domain-nav';
import { usePathname } from 'next/navigation';

export default function DomainsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const currentDomain = pathname.split('/domains/')[1]?.split('/')[0] || '';

  return (
    <>
      <DomainNav currentDomain={currentDomain} />
      {children}
    </>
  );
}
