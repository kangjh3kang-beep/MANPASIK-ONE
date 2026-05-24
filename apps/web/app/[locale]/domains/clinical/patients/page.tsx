'use client';

import React, { useState, useMemo } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Users, Search, ChevronDown, ChevronUp, Calendar, Pill, Clock } from 'lucide-react';

type RiskLevel = '고위험' | '중간' | '저위험';
type FilterType = '전체' | '고위험' | '중간' | '저위험';

interface Patient {
  id: string;
  name: string;
  patientId: string;
  age: number;
  gender: '남' | '여';
  riskLevel: RiskLevel;
  lastVisit: string;
  primaryCondition: string;
  recentVitals: {
    heartRate: number;
    bloodPressure: string;
    spo2: number;
  };
  activeMedications: number;
  nextAppointment: string;
}

const RISK_STYLES: Record<RiskLevel, string> = {
  '고위험': 'bg-red-50 text-red-700 border-red-200',
  '중간': 'bg-amber-50 text-amber-700 border-amber-200',
  '저위험': 'bg-emerald-50 text-emerald-700 border-emerald-200',
};

const PATIENTS: Patient[] = [
  {
    id: 'p-1', name: '김영수', patientId: 'MPS-P001', age: 67, gender: '남',
    riskLevel: '고위험', lastVisit: '2026-05-22', primaryCondition: '제2형 당뇨',
    recentVitals: { heartRate: 88, bloodPressure: '145/92', spo2: 96 },
    activeMedications: 4, nextAppointment: '2026-05-28',
  },
  {
    id: 'p-2', name: '이미경', patientId: 'MPS-P002', age: 54, gender: '여',
    riskLevel: '중간', lastVisit: '2026-05-20', primaryCondition: '고혈압',
    recentVitals: { heartRate: 72, bloodPressure: '132/84', spo2: 98 },
    activeMedications: 2, nextAppointment: '2026-06-03',
  },
  {
    id: 'p-3', name: '박준호', patientId: 'MPS-P003', age: 71, gender: '남',
    riskLevel: '고위험', lastVisit: '2026-05-23', primaryCondition: '만성 신장 질환',
    recentVitals: { heartRate: 78, bloodPressure: '152/96', spo2: 95 },
    activeMedications: 6, nextAppointment: '2026-05-26',
  },
  {
    id: 'p-4', name: '최수진', patientId: 'MPS-P004', age: 42, gender: '여',
    riskLevel: '저위험', lastVisit: '2026-05-18', primaryCondition: '이상지질혈증',
    recentVitals: { heartRate: 68, bloodPressure: '118/76', spo2: 99 },
    activeMedications: 1, nextAppointment: '2026-06-15',
  },
  {
    id: 'p-5', name: '정민우', patientId: 'MPS-P005', age: 58, gender: '남',
    riskLevel: '중간', lastVisit: '2026-05-21', primaryCondition: '관상동맥 질환',
    recentVitals: { heartRate: 82, bloodPressure: '138/88', spo2: 97 },
    activeMedications: 3, nextAppointment: '2026-06-01',
  },
  {
    id: 'p-6', name: '한지은', patientId: 'MPS-P006', age: 35, gender: '여',
    riskLevel: '저위험', lastVisit: '2026-05-19', primaryCondition: '갑상선 기능 저하',
    recentVitals: { heartRate: 64, bloodPressure: '112/72', spo2: 99 },
    activeMedications: 1, nextAppointment: '2026-06-20',
  },
];

const FILTERS: FilterType[] = ['전체', '고위험', '중간', '저위험'];

export default function PatientsPage() {
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<FilterType>('전체');
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const filtered = useMemo(() => {
    return PATIENTS.filter(p => {
      const matchesSearch = search === '' || p.name.includes(search) || p.patientId.toLowerCase().includes(search.toLowerCase());
      const matchesFilter = filter === '전체' || p.riskLevel === filter;
      return matchesSearch && matchesFilter;
    });
  }, [search, filter]);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="환자 관리">
      <DomainHeader title="환자 관리" subtitle="등록된 환자의 위험도와 상태를 관리합니다" icon={Users} />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 검색 및 필터 */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
            <input
              type="text"
              placeholder="환자 이름 또는 ID로 검색"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 bg-white text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
              aria-label="환자 검색"
            />
          </div>
          <div className="flex gap-2" role="group" aria-label="위험도 필터">
            {FILTERS.map((f) => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-4 py-2.5 rounded-xl text-sm font-semibold transition-all whitespace-nowrap ${
                  filter === f
                    ? 'bg-slate-900 text-white'
                    : 'bg-white text-slate-500 border border-slate-200 hover:border-slate-300'
                }`}
              >
                {f}
              </button>
            ))}
          </div>
        </div>

        {/* 환자 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          {/* 헤더 */}
          <div className="hidden sm:grid grid-cols-12 gap-2 px-5 py-3 border-b border-slate-100 bg-slate-50">
            <span className="col-span-3 text-xs font-bold text-slate-500">환자 정보</span>
            <span className="col-span-2 text-xs font-bold text-slate-500">위험도</span>
            <span className="col-span-2 text-xs font-bold text-slate-500">최근 방문</span>
            <span className="col-span-3 text-xs font-bold text-slate-500">주 진단</span>
            <span className="col-span-2 text-xs font-bold text-slate-500 text-right">상세</span>
          </div>

          <div className="divide-y divide-slate-100">
            {filtered.map((patient) => {
              const isExpanded = expandedId === patient.id;
              return (
                <div key={patient.id}>
                  <button
                    onClick={() => setExpandedId(isExpanded ? null : patient.id)}
                    className="w-full grid grid-cols-12 gap-2 items-center px-5 py-4 text-left hover:bg-slate-50 transition-colors"
                    aria-expanded={isExpanded}
                    aria-label={`${patient.name} 환자 상세 정보`}
                  >
                    <div className="col-span-3">
                      <p className="text-sm font-bold text-slate-900">{patient.name}</p>
                      <p className="text-[11px] text-slate-400">{patient.patientId} / {patient.age}세 / {patient.gender}</p>
                    </div>
                    <div className="col-span-2">
                      <span className={`inline-block px-2.5 py-1 rounded-full text-xs font-semibold border ${RISK_STYLES[patient.riskLevel]}`}>
                        {patient.riskLevel}
                      </span>
                    </div>
                    <div className="col-span-2">
                      <p className="text-sm text-slate-600">{patient.lastVisit}</p>
                    </div>
                    <div className="col-span-3">
                      <p className="text-sm text-slate-700 truncate">{patient.primaryCondition}</p>
                    </div>
                    <div className="col-span-2 flex justify-end items-center gap-2">
                      <span className="text-xs text-sky-600 font-semibold">상세 보기</span>
                      {isExpanded ? <ChevronUp className="h-4 w-4 text-slate-400" /> : <ChevronDown className="h-4 w-4 text-slate-400" />}
                    </div>
                  </button>

                  {/* 확장 상세 */}
                  {isExpanded && (
                    <div className="px-5 pb-5 border-t border-slate-100 bg-slate-50/50">
                      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-4">
                        {/* 최근 바이탈 */}
                        <div className="rounded-xl border border-slate-200 bg-white p-4">
                          <h4 className="text-xs font-bold text-slate-500 mb-3 flex items-center gap-1.5">
                            <Clock className="h-3 w-3" />
                            최근 바이탈
                          </h4>
                          <div className="space-y-2">
                            <div className="flex justify-between text-sm">
                              <span className="text-slate-500">심박수</span>
                              <span className="font-semibold text-slate-900">{patient.recentVitals.heartRate} bpm</span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-slate-500">혈압</span>
                              <span className="font-semibold text-slate-900">{patient.recentVitals.bloodPressure} mmHg</span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-slate-500">산소포화도</span>
                              <span className="font-semibold text-slate-900">{patient.recentVitals.spo2}%</span>
                            </div>
                          </div>
                        </div>

                        {/* 활성 약물 */}
                        <div className="rounded-xl border border-slate-200 bg-white p-4">
                          <h4 className="text-xs font-bold text-slate-500 mb-3 flex items-center gap-1.5">
                            <Pill className="h-3 w-3" />
                            투약 현황
                          </h4>
                          <p className="text-2xl font-bold text-slate-900">{patient.activeMedications}종</p>
                          <p className="text-xs text-slate-400 mt-1">현재 활성 처방</p>
                        </div>

                        {/* 다음 예약 */}
                        <div className="rounded-xl border border-slate-200 bg-white p-4">
                          <h4 className="text-xs font-bold text-slate-500 mb-3 flex items-center gap-1.5">
                            <Calendar className="h-3 w-3" />
                            다음 예약
                          </h4>
                          <p className="text-lg font-bold text-slate-900">{patient.nextAppointment}</p>
                          <p className="text-xs text-slate-400 mt-1">정기 검진 예정</p>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}

            {filtered.length === 0 && (
              <div className="py-12 text-center text-sm text-slate-400">
                검색 조건에 일치하는 환자가 없습니다.
              </div>
            )}
          </div>
        </div>

        {/* 요약 */}
        <div className="grid grid-cols-3 gap-4">
          {(['고위험', '중간', '저위험'] as RiskLevel[]).map((level) => {
            const count = PATIENTS.filter(p => p.riskLevel === level).length;
            return (
              <div key={level} className={`rounded-xl border p-4 text-center ${RISK_STYLES[level]}`}>
                <p className="text-xs font-semibold opacity-80">{level}</p>
                <p className="text-2xl font-bold mt-1">{count}명</p>
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
