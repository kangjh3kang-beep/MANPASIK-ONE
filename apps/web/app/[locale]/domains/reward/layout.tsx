'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const REWARD_TABS = [
  { id: 'overview', label: '보상 개요', href: '/domains/reward' },
  { id: 'transactions', label: '거래 내역', href: '/domains/reward/transactions' },
  { id: 'challenges', label: '챌린지', href: '/domains/reward/challenges' },
  { id: 'tiers', label: '등급', href: '/domains/reward/tiers' },
];

export default function RewardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={REWARD_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
