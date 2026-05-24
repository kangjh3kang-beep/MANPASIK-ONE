'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const ADMIN_TABS = [
  { id: 'dashboard', label: '대시보드', href: '/domains/admin' },
  { id: 'users', label: '사용자 관리', href: '/domains/admin/users' },
  { id: 'revenue', label: '매출 분석', href: '/domains/admin/revenue' },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={ADMIN_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
