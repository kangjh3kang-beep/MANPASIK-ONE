'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const APP_TABS = [
  { id: 'dashboard', label: '대시보드', href: '/domains/app' },
  { id: 'activity', label: '활동', href: '/domains/app/activity' },
  { id: 'sleep', label: '수면', href: '/domains/app/sleep' },
  { id: 'medications', label: '복약', href: '/domains/app/medications' },
];

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={APP_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
