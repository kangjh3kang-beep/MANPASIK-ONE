'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const CLINICAL_TABS = [
  { id: 'vitals', label: '바이탈 모니터링', href: '/domains/clinical' },
  { id: 'patients', label: '환자 목록', href: '/domains/clinical/patients' },
  { id: 'labs', label: '검사 결과', href: '/domains/clinical/labs' },
];

export default function ClinicalLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={CLINICAL_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
