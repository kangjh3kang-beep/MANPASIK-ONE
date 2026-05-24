'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const STORE_TABS = [
  { id: 'products', label: '상품', href: '/domains/store' },
  { id: 'orders', label: '주문 내역', href: '/domains/store/orders' },
  { id: 'subscriptions', label: '정기배송', href: '/domains/store/subscriptions' },
];

export default function StoreLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={STORE_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
