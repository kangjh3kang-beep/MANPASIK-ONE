'use client';
import React from 'react';
import { DomainTabNav } from '@mmup/ui';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const PREDICTOR_TABS = [
  { id: 'risk', label: '위험 예측', href: '/domains/predictor' },
  { id: 'detail', label: '상세 분석', href: '/domains/predictor/detail' },
  { id: 'simulator', label: '시뮬레이터', href: '/domains/predictor/simulator' },
];

export default function PredictorLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const localePrefix = useLocalePrefix();
  return (
    <>
      <DomainTabNav tabs={PREDICTOR_TABS} currentPath={pathname} localePrefix={localePrefix} />
      {children}
    </>
  );
}
