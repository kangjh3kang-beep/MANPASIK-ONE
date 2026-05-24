'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const TELEMEDICINE_TABS = [
  { id: 'doctors', label: '의사 찾기', href: '/domains/telemedicine' },
  { id: 'schedule', label: '예약', href: '/domains/telemedicine/schedule' },
  { id: 'history', label: '진료 이력', href: '/domains/telemedicine/history' },
];

export default function TelemedicineLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={TELEMEDICINE_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
