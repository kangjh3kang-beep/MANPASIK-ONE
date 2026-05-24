'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const HARDWARE_TABS = [
  { id: 'devices', label: '장비 현황', href: '/domains/hardware-core' },
  { id: 'ota', label: 'OTA 관리', href: '/domains/hardware-core/ota' },
  { id: 'telemetry', label: '텔레메트리', href: '/domains/hardware-core/telemetry' },
];

export default function HardwareCoreLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={HARDWARE_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
