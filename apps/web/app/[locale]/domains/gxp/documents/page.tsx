'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { FileText, ChevronDown, ChevronUp, Plus, PenLine } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type DocStatus = '승인' | '검토중' | '초안';
type DocCategory = 'SOP' | '지침서' | '기록서' | '보고서';
type SignatureStatus = '서명완료' | '대기중';

interface DocVersion {
  version: string;
  date: string;
  author: string;
  changes: string;
}

interface Document {
  id: string;
  title: string;
  docNumber: string;
  version: string;
  status: DocStatus;
  lastModified: string;
  author: string;
  category: DocCategory;
  signatureStatus: SignatureStatus;
  history: DocVersion[];
}

const STATUS_STYLES: Record<DocStatus, string> = {
  '승인': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '검토중': 'bg-amber-50 text-amber-700 border-amber-200',
  '초안': 'bg-slate-100 text-slate-600 border-slate-200',
};

const SIGN_STYLES: Record<SignatureStatus, string> = {
  '서명완료': 'bg-emerald-600 text-white',
  '대기중': 'bg-amber-500 text-white',
};

const MOCK_DOCUMENTS: Document[] = [
  {
    id: 'd1', title: '혈당 측정 표준 운영 절차', docNumber: 'SOP-GLU-001', version: 'v3.2', status: '승인',
    lastModified: '2026-05-10', author: '정대현', category: 'SOP', signatureStatus: '서명완료',
    history: [
      { version: 'v3.2', date: '2026-05-10', author: '정대현', changes: '측정 오차 허용 범위 업데이트' },
      { version: 'v3.1', date: '2026-03-15', author: '정대현', changes: '캘리브레이션 절차 추가' },
      { version: 'v3.0', date: '2026-01-20', author: '김민수', changes: '전면 개정 - ISO 15197 반영' },
    ],
  },
  {
    id: 'd2', title: '카트리지 보관 및 취급 지침', docNumber: 'GUI-CTR-002', version: 'v2.1', status: '승인',
    lastModified: '2026-04-28', author: '박지훈', category: '지침서', signatureStatus: '서명완료',
    history: [
      { version: 'v2.1', date: '2026-04-28', author: '박지훈', changes: '온도 범위 세분화' },
      { version: 'v2.0', date: '2025-12-10', author: '박지훈', changes: '신규 카트리지 추가' },
      { version: 'v1.0', date: '2025-06-01', author: '이서연', changes: '초기 작성' },
    ],
  },
  {
    id: 'd3', title: '임상시험 데이터 기록서', docNumber: 'REC-CLN-003', version: 'v1.3', status: '검토중',
    lastModified: '2026-05-22', author: '이서연', category: '기록서', signatureStatus: '대기중',
    history: [
      { version: 'v1.3', date: '2026-05-22', author: '이서연', changes: '바이오마커 필드 추가' },
      { version: 'v1.2', date: '2026-04-05', author: '이서연', changes: '데이터 검증 규칙 강화' },
      { version: 'v1.1', date: '2026-02-14', author: '김민수', changes: '양식 레이아웃 수정' },
    ],
  },
  {
    id: 'd4', title: '정확도 검증 보고서', docNumber: 'RPT-VAL-004', version: 'v1.0', status: '초안',
    lastModified: '2026-05-20', author: '박지훈', category: '보고서', signatureStatus: '대기중',
    history: [
      { version: 'v1.0', date: '2026-05-20', author: '박지훈', changes: '초기 작성 - 1차 검증 결과' },
      { version: 'v0.2', date: '2026-05-15', author: '박지훈', changes: '데이터 수집 완료' },
      { version: 'v0.1', date: '2026-05-01', author: '박지훈', changes: '보고서 골격 작성' },
    ],
  },
  {
    id: 'd5', title: '품질 관리 시스템 SOP', docNumber: 'SOP-QMS-005', version: 'v2.0', status: '승인',
    lastModified: '2026-03-30', author: '정대현', category: 'SOP', signatureStatus: '서명완료',
    history: [
      { version: 'v2.0', date: '2026-03-30', author: '정대현', changes: 'ISO 13485:2016 반영 전면 개정' },
      { version: 'v1.1', date: '2025-11-20', author: '정대현', changes: '부적합 처리 절차 보완' },
      { version: 'v1.0', date: '2025-06-15', author: '김민수', changes: '초기 작성' },
    ],
  },
  {
    id: 'd6', title: '리더기 교정 절차 지침서', docNumber: 'GUI-CAL-006', version: 'v1.2', status: '검토중',
    lastModified: '2026-05-18', author: '오태준', category: '지침서', signatureStatus: '대기중',
    history: [
      { version: 'v1.2', date: '2026-05-18', author: '오태준', changes: '자동 교정 알고리즘 추가' },
      { version: 'v1.1', date: '2026-03-10', author: '오태준', changes: '교정 주기 변경' },
      { version: 'v1.0', date: '2025-09-01', author: '이서연', changes: '초기 작성' },
    ],
  },
];

const CATEGORY_FILTERS: (DocCategory | '전체')[] = ['전체', 'SOP', '지침서', '기록서', '보고서'];

export default function DocumentsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [expandedDoc, setExpandedDoc] = useState<string | null>(null);
  const [categoryFilter, setCategoryFilter] = useState<DocCategory | '전체'>('전체');

  const filtered = MOCK_DOCUMENTS.filter(d => categoryFilter === '전체' || d.category === categoryFilter);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="문서 관리">
      <DomainHeader
        title="문서 관리"
        subtitle="GxP 규정 문서를 관리하고 전자서명을 진행합니다"
        icon={FileText}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* 필터 및 새 문서 */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex gap-2">
            {CATEGORY_FILTERS.map(cat => (
              <button
                key={cat}
                onClick={() => setCategoryFilter(cat)}
                className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors ${
                  categoryFilter === cat ? 'bg-sky-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                {cat}
              </button>
            ))}
          </div>
          <button className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors">
            <Plus className="h-4 w-4" />
            새 문서 작성
          </button>
        </div>

        {/* 문서 목록 */}
        <div className="space-y-3">
          {filtered.map(doc => {
            const isExpanded = expandedDoc === doc.id;
            return (
              <div key={doc.id} className="rounded-2xl border border-slate-200 bg-white overflow-hidden transition-shadow hover:shadow-sm">
                <button
                  onClick={() => setExpandedDoc(prev => prev === doc.id ? null : doc.id)}
                  className="w-full flex items-center justify-between p-5 text-left"
                  aria-expanded={isExpanded}
                  aria-label={`${doc.title} 상세 보기`}
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-xs font-mono text-slate-400">{doc.docNumber}</span>
                      <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold border ${STATUS_STYLES[doc.status]}`}>{doc.status}</span>
                      <span className="text-xs text-slate-400">{doc.version}</span>
                    </div>
                    <p className="text-sm font-bold text-slate-900 truncate">{doc.title}</p>
                    <p className="text-xs text-slate-500 mt-0.5">작성자: {doc.author} | 최종 수정: {doc.lastModified}</p>
                  </div>
                  <div className="flex items-center gap-3 ml-4">
                    <span className={`px-3 py-1 rounded-lg text-xs font-semibold ${SIGN_STYLES[doc.signatureStatus]}`}>
                      <PenLine className="h-3 w-3 inline mr-1" />
                      {doc.signatureStatus}
                    </span>
                    {isExpanded ? <ChevronUp className="h-4 w-4 text-slate-400" /> : <ChevronDown className="h-4 w-4 text-slate-400" />}
                  </div>
                </button>
                {isExpanded && (
                  <div className="border-t border-slate-100 px-5 pb-5 pt-4">
                    <h4 className="text-xs font-semibold text-slate-500 mb-3">버전 이력</h4>
                    <div className="space-y-2">
                      {doc.history.map((v, idx) => (
                        <div key={v.version} className={`flex items-start gap-3 p-3 rounded-xl ${idx === 0 ? 'bg-sky-50 border border-sky-100' : 'bg-slate-50'}`}>
                          <span className="text-xs font-mono font-bold text-slate-700 mt-0.5 w-10 shrink-0">{v.version}</span>
                          <div className="flex-1 min-w-0">
                            <p className="text-sm text-slate-900">{v.changes}</p>
                            <p className="text-xs text-slate-400 mt-0.5">{v.author} | {v.date}</p>
                          </div>
                        </div>
                      ))}
                    </div>
                    {doc.signatureStatus === '대기중' && (
                      <button className="mt-4 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors">
                        전자서명 진행
                      </button>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
