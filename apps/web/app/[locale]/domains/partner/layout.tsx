'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const PARTNER_TABS = [
  { id: 'list', label: '파트너 목록', href: '/domains/partner' },
  { id: 'connections', label: '연결 상태', href: '/domains/partner/connections' },
  { id: 'analytics', label: '데이터 교환', href: '/domains/partner/analytics' },
];

export default function PartnerLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={PARTNER_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
