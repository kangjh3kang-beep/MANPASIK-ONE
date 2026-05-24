'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Shield, Download, Trash2, ExternalLink, AlertTriangle } from 'lucide-react';

interface ConsentCell {
  dataType: string;
  purposes: Record<string, boolean>;
}

const PURPOSE_LABELS = ['서비스 제공', 'AI 학습', '연구 활용', '마케팅'];

const INITIAL_CONSENTS: ConsentCell[] = [
  { dataType: '측정 데이터', purposes: { '서비스 제공': true, 'AI 학습': true, '연구 활용': false, '마케팅': false } },
  { dataType: '건강 기록', purposes: { '서비스 제공': true, 'AI 학습': false, '연구 활용': false, '마케팅': false } },
  { dataType: 'AI 분석', purposes: { '서비스 제공': true, 'AI 학습': true, '연구 활용': true, '마케팅': false } },
  { dataType: '활동 데이터', purposes: { '서비스 제공': true, 'AI 학습': false, '연구 활용': false, '마케팅': false } },
  { dataType: '위치 정보', purposes: { '서비스 제공': false, 'AI 학습': false, '연구 활용': false, '마케팅': false } },
];

export default function PrivacyPage() {
  const [consents, setConsents] = useState<ConsentCell[]>(INITIAL_CONSENTS);
  const [exportFormat, setExportFormat] = useState<'JSON' | 'CSV'>('JSON');
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const toggleConsent = (dataIndex: number, purpose: string) => {
    setConsents((prev) =>
      prev.map((cell, i) => {
        if (i !== dataIndex) return cell;
        return {
          ...cell,
          purposes: { ...cell.purposes, [purpose]: !cell.purposes[purpose] },
        };
      }),
    );
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="개인정보 관리">
      <DomainHeader
        title="개인정보 관리"
        subtitle="데이터 동의, 내보내기, 삭제 요청을 관리합니다"
        icon={Shield}
      />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 동의 매트릭스 */}
        <section className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-base font-bold text-slate-900">데이터 활용 동의 관리</h2>
            <p className="text-xs text-slate-400 mt-1">각 데이터 유형별 활용 목적에 대한 동의를 관리합니다</p>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[520px]">
              <thead>
                <tr className="border-b border-slate-100 bg-slate-50">
                  <th className="px-5 py-3 text-left text-xs font-semibold text-slate-500">데이터 유형</th>
                  {PURPOSE_LABELS.map((label) => (
                    <th key={label} className="px-4 py-3 text-center text-xs font-semibold text-slate-500">
                      {label}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {consents.map((row, rowIndex) => (
                  <tr key={row.dataType} className="border-b border-slate-50 last:border-b-0">
                    <td className="px-5 py-3">
                      <span className="text-sm font-medium text-slate-700">{row.dataType}</span>
                    </td>
                    {PURPOSE_LABELS.map((purpose) => (
                      <td key={purpose} className="px-4 py-3 text-center">
                        <div
                          role="switch"
                          aria-checked={row.purposes[purpose]}
                          aria-label={`${row.dataType} - ${purpose}`}
                          onClick={() => toggleConsent(rowIndex, purpose)}
                          className={`relative inline-block w-10 h-5 rounded-full cursor-pointer transition-colors ${
                            row.purposes[purpose] ? 'bg-sky-600' : 'bg-slate-300'
                          }`}
                        >
                          <div
                            className={`absolute top-0.5 h-4 w-4 rounded-full bg-white shadow transition-transform ${
                              row.purposes[purpose] ? 'translate-x-[20px]' : 'translate-x-0.5'
                            }`}
                          />
                        </div>
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="px-6 py-3 border-t border-slate-100 bg-slate-50">
            <p className="text-[11px] text-slate-400">
              마지막 동의 갱신일: 2026년 5월 20일 · 동의 변경은 즉시 적용되며 이전 수집 데이터에는 소급 적용되지 않습니다
            </p>
          </div>
        </section>

        {/* 데이터 내보내기 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Download className="h-4 w-4 text-sky-600" /> 데이터 내보내기
          </h2>
          <p className="text-sm text-slate-600 mb-4">
            만파식 플랫폼에 저장된 모든 개인 데이터를 다운로드할 수 있습니다. GDPR/PIPA 규정에 따라
            데이터 이동권이 보장됩니다.
          </p>
          <div className="flex items-center gap-3">
            <div className="flex gap-1">
              {(['JSON', 'CSV'] as const).map((fmt) => (
                <button
                  key={fmt}
                  onClick={() => setExportFormat(fmt)}
                  className={`px-4 py-2 rounded-lg text-xs font-semibold transition-colors ${
                    exportFormat === fmt
                      ? 'bg-sky-600 text-white'
                      : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                  }`}
                >
                  {fmt}
                </button>
              ))}
            </div>
            <button className="px-5 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors flex items-center gap-1.5">
              <Download className="h-3.5 w-3.5" />
              내 데이터 전체 내보내기
            </button>
          </div>
          <p className="text-[11px] text-slate-400 mt-3">
            데이터 준비에 최대 24시간이 소요될 수 있으며, 완료 시 이메일로 알려드립니다
          </p>
        </section>

        {/* 데이터 삭제 */}
        <section className="rounded-2xl border border-red-200 bg-white p-6">
          <h2 className="text-base font-bold text-red-700 mb-4 flex items-center gap-2">
            <Trash2 className="h-4 w-4 text-red-500" /> 계정 및 데이터 삭제
          </h2>
          <div className="p-4 rounded-xl bg-red-50 border border-red-100 mb-4">
            <div className="flex items-start gap-2">
              <AlertTriangle className="h-4 w-4 text-red-500 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-semibold text-red-800">삭제하면 복구할 수 없습니다</p>
                <p className="text-xs text-red-600 mt-1">
                  모든 측정 데이터, 건강 기록, AI 분석 결과, 구독 정보가 영구적으로 삭제됩니다.
                  진행 중인 구독은 자동 해지되며, 가족 그룹에서 탈퇴됩니다.
                </p>
              </div>
            </div>
          </div>
          {!showDeleteConfirm ? (
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="px-5 py-2.5 rounded-xl border-2 border-red-300 text-sm font-bold text-red-600 hover:bg-red-50 transition-colors"
            >
              계정 및 데이터 삭제 요청
            </button>
          ) : (
            <div className="space-y-3">
              <p className="text-sm font-semibold text-red-800">정말로 삭제하시겠습니까?</p>
              <div className="flex gap-2">
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  className="px-5 py-2.5 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  취소
                </button>
                <button className="px-5 py-2.5 rounded-xl bg-red-600 text-white text-sm font-bold hover:bg-red-700 transition-colors">
                  최종 삭제 확인
                </button>
              </div>
            </div>
          )}
        </section>

        {/* 개인정보처리방침 링크 */}
        <div className="flex justify-center">
          <a
            href="#"
            className="inline-flex items-center gap-1.5 text-sm font-medium text-sky-600 hover:underline"
          >
            <ExternalLink className="h-3.5 w-3.5" />
            전체 개인정보처리방침 보기
          </a>
        </div>
      </div>
    </main>
  );
}
