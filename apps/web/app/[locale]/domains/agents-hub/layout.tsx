'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const AGENTS_TABS = [
  { id: 'agents', label: 'AI 에이전트', href: '/domains/agents-hub' },
  { id: 'registry', label: '모델 레지스트리', href: '/domains/agents-hub/registry' },
  { id: 'monitoring', label: '추론 모니터링', href: '/domains/agents-hub/monitoring' },
];

export default function AgentsHubLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={AGENTS_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
