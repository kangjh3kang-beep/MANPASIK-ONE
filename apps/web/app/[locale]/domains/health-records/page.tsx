'use client';

import React, { useState } from 'react';
import { DomainHeader, MeasurementResultCard, MiniSparkline, NextActionGroup } from '@mmup/ui';
import { FileText, Download, Filter, Calendar, TrendingUp, TrendingDown, Minus, Brain, ShoppingCart, BarChart3 } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type RecordStatus = 'normal' | 'caution' | 'danger';

const MOCK_RECORDS: { id: string; date: string; time: string; biomarker: string; value: number; unit: string; status: RecordStatus; confidence: number; trend: string }[] = [
  { id: '1', date: '2026-05-24', time: '14:32', biomarker: '혈당', value: 95, unit: 'mg/dL', status: 'normal', confidence: 0.96, trend: 'stable' },
  { id: '2', date: '2026-05-23', time: '08:15', biomarker: '혈당', value: 112, unit: 'mg/dL', status: 'normal', confidence: 0.94, trend: 'up' },
  { id: '3', date: '2026-05-22', time: '14:00', biomarker: '당화혈색소', value: 6.8, unit: '%', status: 'caution', confidence: 0.93, trend: 'stable' },
  { id: '4', date: '2026-05-21', time: '09:30', biomarker: '총 콜레스테롤', value: 195, unit: 'mg/dL', status: 'normal', confidence: 0.91, trend: 'down' },
  { id: '5', date: '2026-05-20', time: '08:00', biomarker: '혈당', value: 88, unit: 'mg/dL', status: 'normal', confidence: 0.95, trend: 'down' },
  { id: '6', date: '2026-05-19', time: '14:45', biomarker: '요산', value: 7.2, unit: 'mg/dL', status: 'caution', confidence: 0.92, trend: 'up' },
];

const RANGES: Record<string, { low: number; high: number }> = {
  '혈당': { low: 70, high: 140 },
  '당화혈색소': { low: 4.0, high: 6.5 },
  '총 콜레스테롤': { low: 0, high: 200 },
  '요산': { low: 2.0, high: 7.0 },
};

const TrendIcon = ({ trend }: { trend: string }) => {
  if (trend === 'up') return <TrendingUp className="h-3.5 w-3.5 text-amber-500" />;
  if (trend === 'down') return <TrendingDown className="h-3.5 w-3.5 text-emerald-500" />;
  return <Minus className="h-3.5 w-3.5 text-slate-400" />;
};

export default function HealthRecordsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [selectedRecord, setSelectedRecord] = useState<string | null>(null);
  const [filterBiomarker, setFilterBiomarker] = useState<string>('전체');

  const biomarkers = ['전체', ...new Set(MOCK_RECORDS.map(r => r.biomarker))];
  const filteredRecords = filterBiomarker === '전체'
    ? MOCK_RECORDS
    : MOCK_RECORDS.filter(r => r.biomarker === filterBiomarker);

  const selected = MOCK_RECORDS.find(r => r.id === selectedRecord);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="건강 기록">
      <DomainHeader
        title="건강 기록 타임라인"
        subtitle="모든 측정 결과를 시간순으로 확인하고 추세를 파악합니다"
        icon={FileText}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* 상단: 요약 + 필터 + 내보내기 */}
        <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-2 rounded-lg border border-slate-200 bg-white text-sm">
              <Calendar className="h-4 w-4 text-slate-400" />
              <span className="text-slate-600 font-medium">최근 7일</span>
            </div>
            <div className="flex items-center gap-1">
              <Filter className="h-4 w-4 text-slate-400" />
              {biomarkers.map(b => (
                <button
                  key={b}
                  onClick={() => setFilterBiomarker(b)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors ${
                    filterBiomarker === b
                      ? 'bg-sky-600 text-white'
                      : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-50'
                  }`}
                >
                  {b}
                </button>
              ))}
            </div>
          </div>
          <button className="flex items-center gap-2 px-4 py-2 rounded-lg border border-slate-200 bg-white text-sm font-medium text-slate-600 hover:bg-slate-50 transition-colors">
            <Download className="h-4 w-4" />
            내보내기
          </button>
        </div>

        {/* 범위 요약 바 (FreeStyle Libre AGP 참고) */}
        {filteredRecords.length > 0 && (
          <div className="rounded-2xl border border-slate-200 bg-white p-5 mb-6">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-bold text-slate-700">측정 범위 요약</h3>
              <div className="flex items-center gap-1">
                <MiniSparkline data={filteredRecords.map(r => r.value)} color="#0ea5e9" width={80} height={24} />
              </div>
            </div>
            <div className="flex h-4 rounded-full overflow-hidden">
              {(() => {
                const total = filteredRecords.length;
                const normal = filteredRecords.filter(r => r.status === 'normal').length;
                const caution = filteredRecords.filter(r => r.status === 'caution').length;
                const danger = filteredRecords.filter(r => r.status === 'danger').length;
                return (
                  <>
                    <div className="bg-emerald-400 transition-all" style={{ width: `${(normal/total)*100}%` }} />
                    <div className="bg-amber-400 transition-all" style={{ width: `${(caution/total)*100}%` }} />
                    <div className="bg-red-400 transition-all" style={{ width: `${(danger/total)*100}%` }} />
                  </>
                );
              })()}
            </div>
            <div className="flex gap-4 mt-2 text-xs text-slate-500">
              <span className="flex items-center gap-1"><div className="h-2 w-2 rounded-full bg-emerald-400" />정상 {filteredRecords.filter(r=>r.status==='normal').length}건</span>
              <span className="flex items-center gap-1"><div className="h-2 w-2 rounded-full bg-amber-400" />주의 {filteredRecords.filter(r=>r.status==='caution').length}건</span>
              <span className="flex items-center gap-1"><div className="h-2 w-2 rounded-full bg-red-400" />위험 {filteredRecords.filter(r=>r.status==='danger').length}건</span>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 왼쪽: 타임라인 목록 */}
          <div className="lg:col-span-2 space-y-3">
            {filteredRecords.map(record => {
              const isSelected = selectedRecord === record.id;
              return (
                <button
                  key={record.id}
                  onClick={() => setSelectedRecord(isSelected ? null : record.id)}
                  className={`w-full text-left rounded-2xl border p-4 transition-all ${
                    isSelected
                      ? 'border-sky-500 bg-sky-50 shadow-sm'
                      : 'border-slate-200 bg-white hover:border-sky-300 hover:shadow-sm'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="text-center">
                        <p className="text-xs text-slate-400">{record.date.slice(5)}</p>
                        <p className="text-sm font-bold text-slate-700">{record.time}</p>
                      </div>
                      <div className="h-8 w-px bg-slate-200" />
                      <div>
                        <p className="text-sm font-semibold text-slate-900">{record.biomarker}</p>
                        <div className="flex items-center gap-1.5 mt-0.5">
                          <span className="text-lg font-bold text-slate-900">{record.value}</span>
                          <span className="text-xs text-slate-400">{record.unit}</span>
                          <TrendIcon trend={record.trend} />
                        </div>
                      </div>
                    </div>
                    <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                      record.status === 'normal' ? 'bg-emerald-50 text-emerald-700' :
                      record.status === 'caution' ? 'bg-amber-50 text-amber-700' :
                      'bg-red-50 text-red-700'
                    }`}>
                      {record.status === 'normal' ? '○ 정상' : record.status === 'caution' ? '△ 주의' : '⚠ 위험'}
                    </span>
                  </div>
                </button>
              );
            })}

            {filteredRecords.length === 0 && (
              <div className="text-center py-12">
                <FileText className="h-12 w-12 text-slate-300 mx-auto mb-3" />
                <p className="text-sm text-slate-500">측정 기록이 없습니다</p>
                <Link href={`${localePrefix}/domains/measure`} className="text-sm text-sky-600 font-semibold mt-2 inline-block hover:underline">
                  첫 측정 시작하기 →
                </Link>
              </div>
            )}
          </div>

          {/* 오른쪽: 선택된 기록 상세 */}
          <div className="lg:col-span-1">
            {selected ? (
              <div className="sticky top-24">
                <h3 className="text-sm font-bold text-slate-500 mb-3">상세 결과</h3>
                <MeasurementResultCard
                  biomarker={selected.biomarker}
                  value={selected.value}
                  unit={selected.unit}
                  confidence={selected.confidence}
                  uncertainty={{ value: selected.value * 0.03, ci_lower: Math.round((selected.value * 0.97) * 10) / 10, ci_upper: Math.round((selected.value * 1.03) * 10) / 10 }}
                  referenceRange={RANGES[selected.biomarker] || { low: 0, high: 100 }}
                  status={selected.status}
                  timestamp={`${selected.date} ${selected.time}`}
                />
                <div className="mt-4 flex gap-2">
                  <Link
                    href={`${localePrefix}/domains/ai-coach`}
                    className="flex-1 px-3 py-2 text-center text-sm font-semibold rounded-xl bg-sky-600 text-white hover:bg-sky-700 transition-colors"
                  >
                    AI에게 물어보기
                  </Link>
                  <button className="px-3 py-2 text-sm font-semibold rounded-xl border border-slate-200 text-slate-600 hover:bg-slate-50 transition-colors">
                    공유
                  </button>
                </div>
              </div>
            ) : (
              <div className="sticky top-24 rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-8 text-center">
                <FileText className="h-8 w-8 text-slate-300 mx-auto mb-2" />
                <p className="text-sm text-slate-400">기록을 선택하면<br />상세 결과를 볼 수 있습니다</p>
              </div>
            )}
          </div>
        </div>

        {/* 다음 단계 추천 */}
        <div className="mt-8">
          <NextActionGroup
            title="다음 단계"
            actions={[
              { id: 'coach', title: 'AI 코치 상담', description: '기록을 바탕으로 건강 상담을 받으세요', icon: <Brain className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/ai-coach`, priority: 'high' },
              { id: 'store', title: '카트리지 재구매', description: '소모품을 재주문합니다', icon: <ShoppingCart className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/store` },
              { id: 'predictor', title: '위험 예측 확인', description: '건강 위험도를 미리 분석합니다', icon: <BarChart3 className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/predictor` },
            ]}
          />
        </div>

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-8">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: '건강 측정', path: '/domains/measure', desc: '새 측정 시작' },
              { name: 'AI 코치', path: '/domains/ai-coach', desc: '결과 해석 상담' },
              { name: '건강 데이터', path: '/domains/clinical', desc: '상세 분석' },
              { name: '위험 예측', path: '/domains/predictor', desc: '건강 위험도 분석' },
            ].map(d => (
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
      </div>
    </main>
  );
}
