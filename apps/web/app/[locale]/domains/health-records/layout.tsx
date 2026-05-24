'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const HEALTH_RECORDS_TABS = [
  { id: 'timeline', label: '타임라인', href: '/domains/health-records' },
  { id: 'report', label: '보고서', href: '/domains/health-records/report' },
  { id: 'sharing', label: '데이터 공유', href: '/domains/health-records/sharing' },
  { id: 'insights', label: '인사이트', href: '/domains/health-records/insights' },
];

export default function HealthRecordsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={HEALTH_RECORDS_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
