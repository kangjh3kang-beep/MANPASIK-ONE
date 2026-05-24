'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { FileText, Download, Loader2, Info, CheckCircle } from 'lucide-react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type ReportFormat = 'pdf' | 'fhir' | 'csv';

interface RecentReport {
  id: string;
  name: string;
  format: ReportFormat;
  date: string;
  size: string;
}

const FORMAT_OPTIONS: { value: ReportFormat; label: string; description: string }[] = [
  { value: 'pdf', label: 'PDF', description: '인쇄 및 공유에 적합한 시각적 보고서' },
  { value: 'fhir', label: 'FHIR R4 JSON', description: '의료기관 호환 국제 표준 데이터 형식' },
  { value: 'csv', label: 'CSV', description: '스프레드시트 분석용 원시 데이터' },
];

const BIOMARKER_OPTIONS = [
  { id: 'glucose', label: '혈당', checked: true },
  { id: 'cholesterol', label: '콜레스테롤', checked: true },
  { id: 'hba1c', label: '당화혈색소', checked: false },
  { id: 'bp', label: '혈압', checked: true },
  { id: 'hr', label: '심박수', checked: false },
];

const RECENT_REPORTS: RecentReport[] = [
  { id: 'r1', name: '5월 건강 종합 보고서', format: 'pdf', date: '2026-05-20', size: '1.2 MB' },
  { id: 'r2', name: '4월 FHIR 데이터 내보내기', format: 'fhir', date: '2026-04-30', size: '340 KB' },
  { id: 'r3', name: '1분기 혈당 추이 보고서', format: 'csv', date: '2026-04-01', size: '128 KB' },
];

const FORMAT_LABELS: Record<ReportFormat, string> = { pdf: 'PDF', fhir: 'FHIR R4', csv: 'CSV' };
const FORMAT_BADGE_COLORS: Record<ReportFormat, string> = {
  pdf: 'bg-red-50 text-red-700 border-red-200',
  fhir: 'bg-sky-50 text-sky-700 border-sky-200',
  csv: 'bg-emerald-50 text-emerald-700 border-emerald-200',
};

export default function ReportPage() {
  const localePrefix = useLocalePrefix();
  const [format, setFormat] = useState<ReportFormat>('pdf');
  const [startDate, setStartDate] = useState('2026-04-24');
  const [endDate, setEndDate] = useState('2026-05-24');
  const [biomarkers, setBiomarkers] = useState(BIOMARKER_OPTIONS);
  const [generating, setGenerating] = useState(false);
  const [generated, setGenerated] = useState(false);

  const selectedCount = biomarkers.filter((b) => b.checked).length;

  const toggleBiomarker = (id: string) => {
    setBiomarkers((prev) =>
      prev.map((b) => (b.id === id ? { ...b, checked: !b.checked } : b)),
    );
  };

  const handleGenerate = () => {
    setGenerating(true);
    setGenerated(false);
    setTimeout(() => {
      setGenerating(false);
      setGenerated(true);
    }, 2000);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="건강 보고서">
      <DomainHeader
        title="건강 보고서"
        subtitle="측정 데이터를 기반으로 건강 보고서를 생성하고 내보낼 수 있습니다"
        icon={FileText}
      />

      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 형식 선택 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-900 mb-4">보고서 형식</h2>
          <div className="space-y-2">
            {FORMAT_OPTIONS.map((opt) => (
              <label
                key={opt.value}
                className={`flex items-center gap-3 p-3 rounded-xl border cursor-pointer transition-colors ${
                  format === opt.value
                    ? 'border-sky-300 bg-sky-50'
                    : 'border-slate-100 hover:bg-slate-50'
                }`}
              >
                <input
                  type="radio"
                  name="format"
                  value={opt.value}
                  checked={format === opt.value}
                  onChange={() => setFormat(opt.value)}
                  className="accent-sky-600"
                />
                <div>
                  <p className="text-sm font-semibold text-slate-900">{opt.label}</p>
                  <p className="text-xs text-slate-500">{opt.description}</p>
                </div>
              </label>
            ))}
          </div>
          {format === 'fhir' && (
            <div className="mt-3 flex items-start gap-2 p-3 rounded-xl bg-sky-50 border border-sky-100">
              <Info className="h-4 w-4 text-sky-600 flex-shrink-0 mt-0.5" />
              <p className="text-xs text-sky-700">FHIR R4 형식은 의료기관과 호환되는 국제 표준입니다</p>
            </div>
          )}
        </div>

        {/* 기간 및 바이오마커 선택 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-900 mb-4">기간 및 항목 선택</h2>

          <div className="grid grid-cols-2 gap-4 mb-5">
            <div>
              <label htmlFor="start-date" className="block text-xs font-semibold text-slate-500 mb-1">시작일</label>
              <input
                id="start-date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-sky-500"
              />
            </div>
            <div>
              <label htmlFor="end-date" className="block text-xs font-semibold text-slate-500 mb-1">종료일</label>
              <input
                id="end-date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="w-full rounded-xl border border-slate-200 px-3 py-2 text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-sky-500"
              />
            </div>
          </div>

          <p className="text-xs font-semibold text-slate-500 mb-2">바이오마커 필터</p>
          <div className="flex flex-wrap gap-2">
            {biomarkers.map((b) => (
              <label
                key={b.id}
                className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border cursor-pointer text-xs font-semibold transition-colors ${
                  b.checked
                    ? 'border-sky-300 bg-sky-50 text-sky-700'
                    : 'border-slate-200 bg-white text-slate-500 hover:bg-slate-50'
                }`}
              >
                <input
                  type="checkbox"
                  checked={b.checked}
                  onChange={() => toggleBiomarker(b.id)}
                  className="sr-only"
                />
                {b.checked && <CheckCircle className="h-3 w-3" />}
                {b.label}
              </label>
            ))}
          </div>
        </div>

        {/* 미리보기 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-sm font-bold text-slate-900 mb-3">포함 내용 미리보기</h2>
          <div className="grid grid-cols-3 gap-3 mb-4">
            <div className="rounded-xl bg-slate-50 p-3 text-center">
              <p className="text-lg font-bold text-slate-900">{selectedCount}</p>
              <p className="text-[11px] text-slate-500">선택 항목</p>
            </div>
            <div className="rounded-xl bg-slate-50 p-3 text-center">
              <p className="text-lg font-bold text-slate-900">31</p>
              <p className="text-[11px] text-slate-500">측정 일수</p>
            </div>
            <div className="rounded-xl bg-slate-50 p-3 text-center">
              <p className="text-lg font-bold text-slate-900">127</p>
              <p className="text-[11px] text-slate-500">총 측정 건수</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={handleGenerate}
              disabled={generating || selectedCount === 0}
              className="px-5 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors inline-flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {generating ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  생성 중...
                </>
              ) : (
                <>
                  <FileText className="h-4 w-4" />
                  보고서 생성
                </>
              )}
            </button>
            {generated && (
              <a
                href="#"
                className="px-4 py-2 rounded-xl border border-emerald-200 bg-emerald-50 text-emerald-700 text-sm font-semibold inline-flex items-center gap-1.5 hover:bg-emerald-100 transition-colors"
              >
                <Download className="h-3.5 w-3.5" />
                다운로드
              </a>
            )}
          </div>
        </div>

        {/* 최근 보고서 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-sm font-bold text-slate-900">최근 생성 보고서</h2>
          </div>
          <div className="divide-y divide-slate-50">
            {RECENT_REPORTS.map((report) => (
              <div key={report.id} className="flex items-center justify-between px-6 py-4 hover:bg-slate-50 transition-colors">
                <div className="flex items-center gap-3">
                  <FileText className="h-4 w-4 text-slate-400" />
                  <div>
                    <p className="text-sm font-medium text-slate-900">{report.name}</p>
                    <p className="text-xs text-slate-400">{report.date} | {report.size}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-semibold border ${FORMAT_BADGE_COLORS[report.format]}`}>
                    {FORMAT_LABELS[report.format]}
                  </span>
                  <button className="p-1.5 rounded-lg hover:bg-slate-100 transition-colors" aria-label={`${report.name} 다운로드`}>
                    <Download className="h-4 w-4 text-slate-500" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </main>
  );
}
