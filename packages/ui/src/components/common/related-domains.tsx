'use client';

import React from 'react';
import Link from 'next/link';

interface RelatedDomain {
  name: string;
  path: string;
  desc: string;
}

interface RelatedDomainsProps {
  domains: RelatedDomain[];
  localePrefix: string;
}

export function RelatedDomains({ domains, localePrefix }: RelatedDomainsProps) {
  return (
    <section aria-label="관련 도메인" className="mt-8">
      <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        {domains.map(d => (
          <Link key={d.path} href={`${localePrefix}${d.path}`}
            className="flex items-center gap-3 p-3 rounded-xl border border-slate-200 bg-white hover:border-sky-300 hover:shadow-sm transition-all group">
            <div className="h-8 w-8 rounded-lg bg-sky-50 flex items-center justify-center text-sky-600 group-hover:bg-sky-100">
              <span className="text-sm">→</span>
            </div>
            <div>
              <span className="text-sm font-semibold text-slate-700 group-hover:text-sky-600">{d.name}</span>
              <p className="text-[11px] text-slate-400">{d.desc}</p>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
