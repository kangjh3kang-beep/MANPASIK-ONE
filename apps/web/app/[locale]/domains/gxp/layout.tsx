'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const GXP_TABS = [
  { id: 'compliance', label: '규정 준수', href: '/domains/gxp' },
  { id: 'documents', label: '문서 관리', href: '/domains/gxp/documents' },
  { id: 'capa', label: '시정조치', href: '/domains/gxp/capa' },
];

export default function GxpLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={GXP_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
