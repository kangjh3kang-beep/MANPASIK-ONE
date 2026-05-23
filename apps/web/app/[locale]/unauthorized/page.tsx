'use client';

import React from 'react';
import { Activity, ShieldAlert, ArrowLeft } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export default function UnauthorizedPage() {
  const localePrefix = useLocalePrefix();
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-white to-red-50 px-4">
      <div className="text-center max-w-md">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-red-100 text-red-600">
          <ShieldAlert className="h-8 w-8" />
        </div>

        <h1 className="text-2xl font-bold text-slate-900 mb-3">
          접근 권한이 없습니다
        </h1>
        <p className="text-sm text-slate-500 mb-8 leading-relaxed">
          현재 로그인된 계정의 페르소나(역할)로는 이 도메인에 접근할 수 없습니다.
          다른 계정으로 로그인하거나, 관리자에게 권한을 요청하세요.
        </p>

        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          <Link
            href={`${localePrefix}/login`}
            className="inline-flex items-center justify-center gap-2 rounded-xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white hover:bg-sky-600 transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            다른 계정으로 로그인
          </Link>
          <Link
            href={localePrefix}
            className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-300 bg-white px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition-colors"
          >
            메인으로 돌아가기
          </Link>
        </div>

        <p className="mt-10 text-xs text-slate-400">
          IEC 62304 Class C 준수 — 역할 기반 접근 제어(RBAC) 정책에 의해 차단되었습니다.
        </p>
      </div>
    </div>
  );
}
