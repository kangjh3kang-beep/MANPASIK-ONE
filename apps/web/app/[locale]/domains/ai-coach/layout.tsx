'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const COACH_TABS = [
  { id: 'chat', label: '상담', href: '/domains/ai-coach' },
  { id: 'goals', label: '목표 관리', href: '/domains/ai-coach/goals' },
  { id: 'insights', label: '건강 인사이트', href: '/domains/ai-coach/insights' },
];

export default function AiCoachLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={COACH_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
