'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Shield, Smartphone, Monitor, Tablet, LogOut, Eye, EyeOff, Fingerprint } from 'lucide-react';

interface Session {
  id: string;
  device: string;
  icon: React.ElementType;
  ip: string;
  location: string;
  lastActive: string;
  isCurrent: boolean;
}

interface LoginRecord {
  date: string;
  device: string;
  location: string;
  ip: string;
  status: '성공' | '실패';
}

const SESSIONS: Session[] = [
  {
    id: 's1',
    device: 'iPhone 15 Pro',
    icon: Smartphone,
    ip: '192.168.1.105',
    location: '서울특별시, 대한민국',
    lastActive: '현재 활성',
    isCurrent: true,
  },
  {
    id: 's2',
    device: 'MacBook Pro 14"',
    icon: Monitor,
    ip: '203.247.51.12',
    location: '서울특별시, 대한민국',
    lastActive: '2시간 전',
    isCurrent: false,
  },
  {
    id: 's3',
    device: 'iPad Air',
    icon: Tablet,
    ip: '192.168.1.108',
    location: '서울특별시, 대한민국',
    lastActive: '3일 전',
    isCurrent: false,
  },
];

const LOGIN_HISTORY: LoginRecord[] = [
  { date: '2026-05-24 14:30', device: 'iPhone 15 Pro', location: '서울', ip: '192.168.1.105', status: '성공' },
  { date: '2026-05-24 09:15', device: 'MacBook Pro', location: '서울', ip: '203.247.51.12', status: '성공' },
  { date: '2026-05-23 22:40', device: '알 수 없는 기기', location: '부산', ip: '114.200.33.7', status: '실패' },
  { date: '2026-05-22 08:00', device: 'iPhone 15 Pro', location: '서울', ip: '192.168.1.105', status: '성공' },
  { date: '2026-05-21 18:20', device: 'iPad Air', location: '서울', ip: '192.168.1.108', status: '성공' },
];

export default function SecurityPage() {
  const [showBackupCodes, setShowBackupCodes] = useState(false);
  const [showCurrentPw, setShowCurrentPw] = useState(false);
  const [showNewPw, setShowNewPw] = useState(false);
  const [biometricEnabled, setBiometricEnabled] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [newPassword, setNewPassword] = useState('');

  const handleNewPasswordChange = (value: string) => {
    setNewPassword(value);
    let strength = 0;
    if (value.length >= 8) strength++;
    if (/[A-Z]/.test(value)) strength++;
    if (/[0-9]/.test(value)) strength++;
    if (/[^A-Za-z0-9]/.test(value)) strength++;
    setPasswordStrength(strength);
  };

  const strengthLabels = ['', '약함', '보통', '강함', '매우 강함'];
  const strengthColors = ['', 'bg-red-500', 'bg-amber-500', 'bg-emerald-500', 'bg-emerald-600'];

  return (
    <main className="min-h-screen bg-slate-50" aria-label="보안 설정">
      <DomainHeader
        title="보안 설정"
        subtitle="계정 보안과 로그인 활동을 관리합니다"
        icon={Shield}
      />

      <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {/* 2단계 인증 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4 flex items-center gap-2">
            <Shield className="h-4 w-4 text-sky-600" /> 2단계 인증 (2FA)
          </h2>
          <div className="flex items-center justify-between mb-4">
            <div>
              <p className="text-sm font-medium text-slate-700">현재 상태</p>
              <span className="inline-block mt-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-red-50 text-red-600 border border-red-200">
                비활성
              </span>
            </div>
            <button className="px-5 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors">
              활성화
            </button>
          </div>

          <div className="p-4 rounded-xl bg-slate-50 border border-slate-100 mb-4">
            <p className="text-xs font-semibold text-slate-500 mb-2">QR 코드</p>
            <div className="h-32 w-32 mx-auto rounded-lg bg-slate-200 flex items-center justify-center text-slate-400 text-xs">
              QR 코드 영역
            </div>
            <p className="text-[11px] text-slate-400 text-center mt-2">인증 앱(Google Authenticator 등)으로 스캔하세요</p>
          </div>

          <div>
            <div className="flex items-center justify-between mb-2">
              <p className="text-xs font-semibold text-slate-500">백업 코드</p>
              <button
                onClick={() => setShowBackupCodes(!showBackupCodes)}
                className="text-xs text-sky-600 hover:underline"
              >
                {showBackupCodes ? '숨기기' : '보기'}
              </button>
            </div>
            <div className="grid grid-cols-2 gap-1.5">
              {['A3K9-X2M1', 'B7P4-Q8N5', 'C1R6-Y3L0', 'D5S2-Z9H7'].map((code) => (
                <span
                  key={code}
                  className="px-3 py-1.5 rounded-lg bg-slate-100 text-center text-xs font-mono text-slate-600"
                >
                  {showBackupCodes ? code : '****-****'}
                </span>
              ))}
            </div>
          </div>
        </section>

        {/* 활성 세션 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">활성 세션</h2>
          <div className="space-y-3">
            {SESSIONS.map((session) => {
              const IconComp = session.icon;
              return (
                <div
                  key={session.id}
                  className="flex items-center justify-between p-4 rounded-xl border border-slate-100 hover:border-slate-200 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="h-10 w-10 rounded-lg bg-slate-50 flex items-center justify-center">
                      <IconComp className="h-5 w-5 text-slate-500" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="text-sm font-semibold text-slate-800">{session.device}</p>
                        {session.isCurrent && (
                          <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-200">
                            이 기기
                          </span>
                        )}
                      </div>
                      <p className="text-xs text-slate-400">
                        {session.ip} · {session.location} · {session.lastActive}
                      </p>
                    </div>
                  </div>
                  {!session.isCurrent && (
                    <button className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-red-600 hover:bg-red-50 transition-colors">
                      <LogOut className="h-3.5 w-3.5" />
                      로그아웃
                    </button>
                  )}
                </div>
              );
            })}
          </div>
        </section>

        {/* 로그인 기록 */}
        <section className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="text-base font-bold text-slate-900">로그인 기록</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[560px]">
              <thead>
                <tr className="border-b border-slate-100 bg-slate-50">
                  <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">일시</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">기기</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">위치</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">IP</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold text-slate-500">상태</th>
                </tr>
              </thead>
              <tbody>
                {LOGIN_HISTORY.map((record, idx) => (
                  <tr key={idx} className="border-b border-slate-50 last:border-b-0">
                    <td className="px-4 py-3 text-xs text-slate-600 font-mono">{record.date}</td>
                    <td className="px-4 py-3 text-xs text-slate-700 font-medium">{record.device}</td>
                    <td className="px-4 py-3 text-xs text-slate-500">{record.location}</td>
                    <td className="px-4 py-3 text-xs text-slate-400 font-mono">{record.ip}</td>
                    <td className="px-4 py-3">
                      <span
                        className={`px-2 py-0.5 rounded-full text-[11px] font-semibold ${
                          record.status === '성공'
                            ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                            : 'bg-red-50 text-red-700 border border-red-200'
                        }`}
                      >
                        {record.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* 비밀번호 변경 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <h2 className="text-base font-bold text-slate-900 mb-4">비밀번호 변경</h2>
          <div className="space-y-4 max-w-md">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">현재 비밀번호</label>
              <div className="relative">
                <input
                  type={showCurrentPw ? 'text' : 'password'}
                  className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500"
                  placeholder="현재 비밀번호를 입력하세요"
                />
                <button
                  onClick={() => setShowCurrentPw(!showCurrentPw)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                >
                  {showCurrentPw ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">새 비밀번호</label>
              <div className="relative">
                <input
                  type={showNewPw ? 'text' : 'password'}
                  value={newPassword}
                  onChange={(e) => handleNewPasswordChange(e.target.value)}
                  className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500"
                  placeholder="새 비밀번호를 입력하세요"
                />
                <button
                  onClick={() => setShowNewPw(!showNewPw)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                >
                  {showNewPw ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {newPassword && (
                <div className="mt-2">
                  <div className="flex gap-1 mb-1">
                    {[1, 2, 3, 4].map((level) => (
                      <div
                        key={level}
                        className={`h-1.5 flex-1 rounded-full ${
                          level <= passwordStrength ? strengthColors[passwordStrength] : 'bg-slate-200'
                        }`}
                      />
                    ))}
                  </div>
                  <p className="text-[11px] text-slate-500">강도: {strengthLabels[passwordStrength]}</p>
                </div>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">새 비밀번호 확인</label>
              <input
                type="password"
                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500"
                placeholder="새 비밀번호를 다시 입력하세요"
              />
            </div>
            <button className="px-6 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors">
              비밀번호 변경
            </button>
          </div>
        </section>

        {/* 생체 인증 */}
        <section className="rounded-2xl border border-slate-200 bg-white p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 rounded-lg bg-purple-50 flex items-center justify-center">
                <Fingerprint className="h-5 w-5 text-purple-600" />
              </div>
              <div>
                <p className="text-sm font-bold text-slate-900">생체 인증</p>
                <p className="text-xs text-slate-400">지문/Face ID 인증 활성화</p>
              </div>
            </div>
            <div
              role="switch"
              aria-checked={biometricEnabled}
              onClick={() => setBiometricEnabled(!biometricEnabled)}
              className={`relative w-11 h-6 rounded-full cursor-pointer transition-colors ${
                biometricEnabled ? 'bg-sky-600' : 'bg-slate-300'
              }`}
            >
              <div
                className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${
                  biometricEnabled ? 'translate-x-[22px]' : 'translate-x-0.5'
                }`}
              />
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
