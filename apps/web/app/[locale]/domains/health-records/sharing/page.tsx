'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Shield, ChevronDown, Clock, Eye, AlertTriangle } from 'lucide-react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type AccessLevel = '읽기전용' | '읽기+쓰기';

interface DataPermission {
  type: string;
  enabled: boolean;
}

interface Institution {
  id: string;
  name: string;
  type: string;
  accessLevel: AccessLevel;
  expiryDate: string;
  permissions: DataPermission[];
}

interface ActivityEntry {
  id: string;
  institution: string;
  action: string;
  dataType: string;
  timestamp: string;
}

const MOCK_INSTITUTIONS: Institution[] = [
  {
    id: 'inst-1',
    name: '서울대병원',
    type: '상급종합병원',
    accessLevel: '읽기전용',
    expiryDate: '2027-05-24',
    permissions: [
      { type: '측정 데이터', enabled: true },
      { type: '처방전', enabled: true },
      { type: 'AI 분석 결과', enabled: false },
    ],
  },
  {
    id: 'inst-2',
    name: '연세의료원',
    type: '상급종합병원',
    accessLevel: '읽기+쓰기',
    expiryDate: '2026-12-31',
    permissions: [
      { type: '측정 데이터', enabled: true },
      { type: '처방전', enabled: true },
      { type: 'AI 분석 결과', enabled: true },
    ],
  },
  {
    id: 'inst-3',
    name: '삼성서울병원',
    type: '상급종합병원',
    accessLevel: '읽기전용',
    expiryDate: '2026-11-15',
    permissions: [
      { type: '측정 데이터', enabled: true },
      { type: '처방전', enabled: false },
      { type: 'AI 분석 결과', enabled: false },
    ],
  },
  {
    id: 'inst-4',
    name: '가정의학과 김의사',
    type: '의원',
    accessLevel: '읽기+쓰기',
    expiryDate: '2027-01-31',
    permissions: [
      { type: '측정 데이터', enabled: true },
      { type: '처방전', enabled: true },
      { type: 'AI 분석 결과', enabled: true },
    ],
  },
];

const ACTIVITY_LOG: ActivityEntry[] = [
  { id: 'a1', institution: '연세의료원', action: '데이터 조회', dataType: '혈당 측정 기록', timestamp: '2026-05-24 10:15' },
  { id: 'a2', institution: '가정의학과 김의사', action: '처방전 등록', dataType: '처방전', timestamp: '2026-05-23 14:32' },
  { id: 'a3', institution: '서울대병원', action: '데이터 조회', dataType: '콜레스테롤 추이', timestamp: '2026-05-22 09:48' },
];

export default function SharingPage() {
  const localePrefix = useLocalePrefix();
  const [institutions, setInstitutions] = useState(MOCK_INSTITUTIONS);

  const togglePermission = (instId: string, permType: string) => {
    setInstitutions((prev) =>
      prev.map((inst) =>
        inst.id === instId
          ? {
              ...inst,
              permissions: inst.permissions.map((p) =>
                p.type === permType ? { ...p, enabled: !p.enabled } : p,
              ),
            }
          : inst,
      ),
    );
  };

  const updateAccessLevel = (instId: string, level: AccessLevel) => {
    setInstitutions((prev) =>
      prev.map((inst) =>
        inst.id === instId ? { ...inst, accessLevel: level } : inst,
      ),
    );
  };

  const revokeConsent = (instId: string) => {
    setInstitutions((prev) => prev.filter((inst) => inst.id !== instId));
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="데이터 공유 관리">
      <DomainHeader
        title="데이터 공유 관리"
        subtitle="의료기관별 데이터 접근 권한을 관리합니다"
        icon={Shield}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 기관 목록 */}
        <div className="space-y-4">
          {institutions.map((inst) => (
            <div key={inst.id} className="rounded-2xl border border-slate-200 bg-white p-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="text-sm font-bold text-slate-900">{inst.name}</h3>
                    <span className="px-2 py-0.5 rounded-full bg-slate-100 text-[11px] font-medium text-slate-600">
                      {inst.type}
                    </span>
                  </div>
                  <p className="text-xs text-slate-400">
                    동의 만료일: <span className="font-semibold text-slate-600">{inst.expiryDate}</span>
                  </p>
                </div>
                <button
                  onClick={() => revokeConsent(inst.id)}
                  className="px-3 py-1.5 rounded-xl border border-red-200 text-xs font-semibold text-red-600 hover:bg-red-50 transition-colors"
                >
                  동의 철회
                </button>
              </div>

              {/* 데이터 유형별 토글 */}
              <div className="space-y-2 mb-4">
                {inst.permissions.map((perm) => (
                  <div
                    key={perm.type}
                    className="flex items-center justify-between p-3 rounded-xl bg-slate-50"
                  >
                    <span className="text-sm text-slate-700">{perm.type}</span>
                    <button
                      role="switch"
                      aria-checked={perm.enabled}
                      aria-label={`${inst.name} ${perm.type} 공유 ${perm.enabled ? '비활성화' : '활성화'}`}
                      onClick={() => togglePermission(inst.id, perm.type)}
                      className={`relative w-10 h-5 rounded-full transition-colors ${
                        perm.enabled ? 'bg-sky-600' : 'bg-slate-300'
                      }`}
                    >
                      <span
                        className={`absolute top-0.5 left-0.5 h-4 w-4 rounded-full bg-white shadow transition-transform ${
                          perm.enabled ? 'translate-x-5' : 'translate-x-0'
                        }`}
                      />
                    </button>
                  </div>
                ))}
              </div>

              {/* 접근 수준 선택 */}
              <div className="flex items-center gap-3">
                <span className="text-xs font-semibold text-slate-500">접근 수준</span>
                <div className="relative">
                  <select
                    value={inst.accessLevel}
                    onChange={(e) => updateAccessLevel(inst.id, e.target.value as AccessLevel)}
                    className="appearance-none rounded-xl border border-slate-200 bg-white px-3 py-1.5 pr-8 text-xs font-medium text-slate-700 focus:outline-none focus:ring-2 focus:ring-sky-500"
                  >
                    <option value="읽기전용">읽기전용</option>
                    <option value="읽기+쓰기">읽기+쓰기</option>
                  </select>
                  <ChevronDown className="absolute right-2 top-1/2 -translate-y-1/2 h-3 w-3 text-slate-400 pointer-events-none" />
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* 접근 활동 로그 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-sm font-bold text-slate-900">최근 접근 활동</h2>
          </div>
          <div className="divide-y divide-slate-50">
            {ACTIVITY_LOG.map((entry) => (
              <div key={entry.id} className="flex items-center gap-3 px-6 py-4">
                <div className="h-8 w-8 rounded-xl bg-sky-50 flex items-center justify-center flex-shrink-0">
                  <Eye className="h-4 w-4 text-sky-600" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-slate-900">
                    <span className="font-semibold">{entry.institution}</span>
                    <span className="text-slate-400"> | </span>
                    {entry.action}
                  </p>
                  <p className="text-xs text-slate-400">
                    {entry.dataType} | {entry.timestamp}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 개인정보 안내 */}
        <div className="rounded-2xl border border-slate-100 bg-slate-50 p-5">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-4 w-4 text-slate-400 flex-shrink-0 mt-0.5" />
            <div className="text-xs text-slate-500 leading-relaxed space-y-1">
              <p className="font-semibold text-slate-600">개인정보 보호 안내</p>
              <p>
                데이터 공유는 사용자의 명시적 동의에 기반하며, 언제든지 동의를 철회할 수 있습니다.
                공유된 데이터는 HIPAA, GDPR, PIPA 규정에 따라 암호화되어 전송되며,
                접근 기록은 모두 감사 로그에 저장됩니다.
                동의 철회 시 해당 기관의 데이터 접근 권한이 즉시 회수됩니다.
              </p>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
