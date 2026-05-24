'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const MEASURE_TABS = [
  { id: 'measure', label: '측정', href: '/domains/measure' },
  { id: 'calibration', label: '보정 관리', href: '/domains/measure/calibration' },
  { id: 'queue', label: '동기화 큐', href: '/domains/measure/queue' },
];

export default function MeasureLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={MEASURE_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
