'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const SETTINGS_TABS = [
  { id: 'general', label: '일반', href: '/domains/settings' },
  { id: 'security', label: '보안', href: '/domains/settings/security' },
  { id: 'privacy', label: '개인정보', href: '/domains/settings/privacy' },
  { id: 'family', label: '가족', href: '/domains/settings/family' },
];

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={SETTINGS_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
