'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Users, Phone, Mail, Plus, ShieldAlert, Baby } from 'lucide-react';

interface FamilyMember {
  id: string;
  name: string;
  emoji: string;
  role: '관리자' | '보호자' | '구성원';
  relationship: string;
  isMinor: boolean;
  permissions: {
    viewMeasurements: boolean;
    receiveAlerts: boolean;
    shareReports: boolean;
  };
}

const ROLE_STYLES: Record<string, string> = {
  '관리자': 'bg-sky-50 text-sky-700 border-sky-200',
  '보호자': 'bg-amber-50 text-amber-700 border-amber-200',
  '구성원': 'bg-slate-50 text-slate-600 border-slate-200',
};

const INITIAL_MEMBERS: FamilyMember[] = [
  {
    id: 'm1',
    name: '김건강',
    emoji: '👨',
    role: '관리자',
    relationship: '본인',
    isMinor: false,
    permissions: { viewMeasurements: true, receiveAlerts: true, shareReports: true },
  },
  {
    id: 'm2',
    name: '이돌봄',
    emoji: '👩',
    role: '보호자',
    relationship: '배우자',
    isMinor: false,
    permissions: { viewMeasurements: true, receiveAlerts: true, shareReports: false },
  },
  {
    id: 'm3',
    name: '김자녀',
    emoji: '👦',
    role: '구성원',
    relationship: '자녀',
    isMinor: true,
    permissions: { viewMeasurements: false, receiveAlerts: false, shareReports: false },
  },
];

export default function FamilyPage() {
  const [members, setMembers] = useState<FamilyMember[]>(INITIAL_MEMBERS);
  const [showInviteForm, setShowInviteForm] = useState(false);
  const [emergencyAlertEnabled, setEmergencyAlertEnabled] = useState(true);

  const togglePermission = (memberId: string, key: keyof FamilyMember['permissions']) => {
    setMembers((prev) =>
      prev.map((m) => {
        if (m.id !== memberId) return m;
        return { ...m, permissions: { ...m.permissions, [key]: !m.permissions[key] } };
      }),
    );
  };

  const permissionLabels: { key: keyof FamilyMember['permissions']; label: string }[] = [
    { key: 'viewMeasurements', label: '측정 데이터 조회' },
    { key: 'receiveAlerts', label: '위험 알림 수신' },
    { key: 'shareReports', label: '보고서 공유' },
  ];

  return (
    <main className="min-h-screen bg-slate-50" aria-label="가족 계정 관리">
      <DomainHeader
        title="가족 계정 관리"
        subtitle="가족 구성원의 건강 데이터 접근 권한을 관리합니다"
        icon={Users}
      />

      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 가족 그룹 카드 */}
        <section className="rounded-2xl border border-slate-200 bg-gradient-to-r from-sky-50 to-cyan-50 p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-bold text-slate-900">김씨 가족</h2>
              <p className="text-sm text-slate-500 mt-0.5">구성원 {members.length}명</p>
            </div>
            <div className="flex -space-x-2">
              {members.map((m) => (
                <div
                  key={m.id}
                  className="h-10 w-10 rounded-full bg-white border-2 border-white shadow flex items-center justify-center text-lg"
                  title={m.name}
                >
                  {m.emoji}
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* 구성원 카드 */}
        <section aria-label="가족 구성원" className="space-y-4">
          {members.map((member) => (
            <div key={member.id} className="rounded-2xl border border-slate-200 bg-white p-5">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-12 w-12 rounded-full bg-slate-50 flex items-center justify-center text-2xl">
                  {member.emoji}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-bold text-slate-900">{member.name}</p>
                    <span
                      className={`px-2 py-0.5 rounded-full text-[10px] font-semibold border ${ROLE_STYLES[member.role]}`}
                    >
                      {member.role}
                    </span>
                    {member.isMinor && (
                      <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-purple-50 text-purple-600 border border-purple-200">
                        미성년
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-slate-400">{member.relationship}</p>
                </div>
              </div>

              {/* 권한 토글 */}
              <div className="space-y-2.5 pl-15">
                {permissionLabels.map(({ key, label }) => (
                  <label key={key} className="flex items-center justify-between cursor-pointer">
                    <span className="text-sm text-slate-600">{label}</span>
                    <div
                      role="switch"
                      aria-checked={member.permissions[key]}
                      aria-label={`${member.name} ${label}`}
                      onClick={() => togglePermission(member.id, key)}
                      className={`relative w-10 h-5 rounded-full cursor-pointer transition-colors ${
                        member.permissions[key] ? 'bg-sky-600' : 'bg-slate-300'
                      }`}
                    >
                      <div
                        className={`absolute top-0.5 h-4 w-4 rounded-full bg-white shadow transition-transform ${
                          member.permissions[key] ? 'translate-x-[20px]' : 'translate-x-0.5'
                        }`}
                      />
                    </div>
                  </label>
                ))}
              </div>

              {/* 미성년 보호 설명 */}
              {member.isMinor && (
                <div className="mt-4 p-3 rounded-xl bg-purple-50 border border-purple-100">
                  <div className="flex items-start gap-2">
                    <Baby className="h-4 w-4 text-purple-500 flex-shrink-0 mt-0.5" />
                    <div>
                      <p className="text-xs font-semibold text-purple-800">미성년 보호 모드</p>
                      <p className="text-[11px] text-purple-600 mt-0.5">
                        미성년 구성원의 데이터 접근은 보호자 승인이 필요하며, 민감한 건강 정보(정신건강, 성관련)는
                        보호자에게도 제한적으로 공유됩니다.
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </section>

        {/* 가족 구성원 초대 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-slate-900">가족 구성원 초대</h2>
            <button
              onClick={() => setShowInviteForm(!showInviteForm)}
              className="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              초대하기
            </button>
          </div>

          {showInviteForm && (
            <div className="space-y-3 p-4 rounded-xl bg-slate-50 border border-slate-100">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">이름</label>
                <input
                  type="text"
                  className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                  placeholder="구성원 이름"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">관계</label>
                <select className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 bg-white">
                  <option>부모</option>
                  <option>배우자</option>
                  <option>자녀</option>
                  <option>형제/자매</option>
                  <option>기타</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">이메일</label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                  <input
                    type="email"
                    className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
                    placeholder="초대할 이메일 주소"
                  />
                </div>
              </div>
              <div className="flex gap-2 pt-1">
                <button
                  onClick={() => setShowInviteForm(false)}
                  className="px-4 py-2 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-100 transition-colors"
                >
                  취소
                </button>
                <button className="px-5 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors">
                  초대 보내기
                </button>
              </div>
            </div>
          )}
        </section>

        {/* 긴급 연락처 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <ShieldAlert className="h-4 w-4 text-red-500" /> 긴급 연락처
          </h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 rounded-xl bg-red-50 border border-red-100">
              <div className="flex items-center gap-3">
                <Phone className="h-4 w-4 text-red-500" />
                <div>
                  <p className="text-sm font-semibold text-slate-800">이돌봄 (배우자)</p>
                  <p className="text-xs text-slate-500">010-1234-5678</p>
                </div>
              </div>
              <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-red-100 text-red-600">
                주 연락처
              </span>
            </div>
            <label className="flex items-center justify-between cursor-pointer p-3 rounded-xl hover:bg-slate-50 transition-colors">
              <div>
                <p className="text-sm font-medium text-slate-700">긴급 알림 자동발송</p>
                <p className="text-xs text-slate-400">위험 수치 감지 시 긴급 연락처에 자동 알림</p>
              </div>
              <div
                role="switch"
                aria-checked={emergencyAlertEnabled}
                onClick={() => setEmergencyAlertEnabled(!emergencyAlertEnabled)}
                className={`relative w-11 h-6 rounded-full cursor-pointer transition-colors ${
                  emergencyAlertEnabled ? 'bg-red-500' : 'bg-slate-300'
                }`}
              >
                <div
                  className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${
                    emergencyAlertEnabled ? 'translate-x-[22px]' : 'translate-x-0.5'
                  }`}
                />
              </div>
            </label>
          </div>
        </section>
      </div>
    </main>
  );
}
