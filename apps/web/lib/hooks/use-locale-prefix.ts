'use client';

import { usePathname } from 'next/navigation';

/**
 * 현재 pathname에서 locale prefix를 추출합니다.
 * 예: /ko/domains/clinical → '/ko'
 * 감지 실패 시 '/ko' 기본값
 */
export function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}
