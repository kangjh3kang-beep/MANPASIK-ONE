'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const DEV_PORTAL_TABS = [
  { id: 'api', label: 'API 관리', href: '/domains/dev-portal' },
  { id: 'playground', label: '플레이그라운드', href: '/domains/dev-portal/playground' },
  { id: 'webhooks', label: '웹훅', href: '/domains/dev-portal/webhooks' },
];

export default function DevPortalLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={DEV_PORTAL_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
