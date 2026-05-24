'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Video, Search, Star, Clock, MapPin, Phone, Calendar, Shield, ChevronRight, CheckCircle } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const DOCTORS = [
  { id: '1', name: '김내과', specialty: '내과', rating: 4.8, reviews: 234, available: true, nextSlot: '오늘 15:00', photo: '👨‍⚕️', fee: '₩25,000' },
  { id: '2', name: '이피부', specialty: '피부과', rating: 4.9, reviews: 189, available: true, nextSlot: '오늘 16:30', photo: '👩‍⚕️', fee: '₩30,000' },
  { id: '3', name: '박가정', specialty: '가정의학과', rating: 4.7, reviews: 312, available: false, nextSlot: '내일 09:00', photo: '👨‍⚕️', fee: '₩20,000' },
  { id: '4', name: '최소아', specialty: '소아과', rating: 4.9, reviews: 156, available: true, nextSlot: '오늘 17:00', photo: '👩‍⚕️', fee: '₩25,000' },
];

type BookingStep = 'search' | 'select' | 'consent' | 'waiting' | 'call';

export default function TelemedicinePage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [step, setStep] = useState<BookingStep>('search');
  const [selectedDoctor, setSelectedDoctor] = useState<typeof DOCTORS[0] | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [consentGiven, setConsentGiven] = useState({ dataShare: false, recording: false });

  const filteredDoctors = DOCTORS.filter(d =>
    d.name.includes(searchQuery) || d.specialty.includes(searchQuery) || searchQuery === ''
  );

  const handleSelectDoctor = (doctor: typeof DOCTORS[0]) => {
    setSelectedDoctor(doctor);
    setStep('consent');
  };

  const handleStartCall = () => {
    setStep('waiting');
    setTimeout(() => setStep('call'), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="화상진료">
      <DomainHeader
        title="화상진료"
        subtitle="전문 의료진과 화상으로 진료받을 수 있습니다"
        icon={Video}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-4xl mx-auto px-6 py-8">

        {/* 의사 검색 */}
        {step === 'search' && (
          <>
            <div className="relative mb-6">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder="진료과목 또는 의사 이름으로 검색"
                className="w-full pl-12 pr-4 py-3.5 rounded-2xl border border-slate-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-400"
              />
            </div>

            <div className="space-y-3">
              {filteredDoctors.map(doctor => (
                <button
                  key={doctor.id}
                  onClick={() => handleSelectDoctor(doctor)}
                  className="w-full text-left rounded-2xl border border-slate-200 bg-white p-5 hover:border-sky-300 hover:shadow-sm transition-all group"
                >
                  <div className="flex items-center gap-4">
                    <span className="text-4xl">{doctor.photo}</span>
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-base font-bold text-slate-900">{doctor.name} 전문의</h3>
                        <span className="px-2 py-0.5 rounded-full bg-slate-100 text-xs font-medium text-slate-600">{doctor.specialty}</span>
                      </div>
                      <div className="flex items-center gap-3 text-xs text-slate-500">
                        <span className="flex items-center gap-1"><Star className="h-3 w-3 text-amber-500" />{doctor.rating} ({doctor.reviews})</span>
                        <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{doctor.nextSlot}</span>
                        <span className="font-semibold text-slate-700">{doctor.fee}</span>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {doctor.available ? (
                        <span className="px-3 py-1.5 rounded-full bg-emerald-50 text-emerald-700 text-xs font-semibold">진료 가능</span>
                      ) : (
                        <span className="px-3 py-1.5 rounded-full bg-slate-100 text-slate-500 text-xs font-semibold">예약만 가능</span>
                      )}
                      <ChevronRight className="h-4 w-4 text-slate-300 group-hover:text-sky-500 transition-colors" />
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </>
        )}

        {/* 데이터 공유 동의 */}
        {step === 'consent' && selectedDoctor && (
          <div className="rounded-2xl border border-slate-200 bg-white p-8">
            <div className="flex items-center gap-3 mb-6">
              <Shield className="h-6 w-6 text-sky-600" />
              <h2 className="text-lg font-bold text-slate-900">진료 전 동의</h2>
            </div>

            <div className="flex items-center gap-4 p-4 rounded-xl bg-slate-50 mb-6">
              <span className="text-3xl">{selectedDoctor.photo}</span>
              <div>
                <p className="font-bold text-slate-900">{selectedDoctor.name} 전문의 ({selectedDoctor.specialty})</p>
                <p className="text-sm text-slate-500">{selectedDoctor.nextSlot} · {selectedDoctor.fee}</p>
              </div>
            </div>

            <div className="space-y-4 mb-8">
              <label className="flex items-start gap-3 cursor-pointer">
                <input type="checkbox" checked={consentGiven.dataShare}
                  onChange={e => setConsentGiven(prev => ({ ...prev, dataShare: e.target.checked }))}
                  className="mt-1 h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
                <div>
                  <p className="text-sm font-medium text-slate-700">건강 측정 데이터 공유에 동의합니다</p>
                  <p className="text-xs text-slate-400">최근 90일 측정 기록과 AI 분석 요약이 의사에게 공유됩니다. 언제든 철회할 수 있습니다.</p>
                </div>
              </label>
              <label className="flex items-start gap-3 cursor-pointer">
                <input type="checkbox" checked={consentGiven.recording}
                  onChange={e => setConsentGiven(prev => ({ ...prev, recording: e.target.checked }))}
                  className="mt-1 h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500" />
                <div>
                  <p className="text-sm font-medium text-slate-700">진료 내용 기록에 동의합니다</p>
                  <p className="text-xs text-slate-400">진료 요약이 자동 생성되어 건강 기록에 저장됩니다.</p>
                </div>
              </label>
            </div>

            <div className="flex gap-3">
              <button onClick={() => setStep('search')} className="px-6 py-3 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50">
                뒤로
              </button>
              <button
                onClick={handleStartCall}
                disabled={!consentGiven.dataShare}
                className="flex-1 px-6 py-3 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
              >
                진료 시작
              </button>
            </div>
          </div>
        )}

        {/* 연결 대기 */}
        {step === 'waiting' && (
          <div className="rounded-2xl border border-slate-200 bg-white p-12 text-center">
            <div className="h-16 w-16 rounded-full bg-sky-50 flex items-center justify-center mx-auto mb-4 animate-pulse">
              <Phone className="h-8 w-8 text-sky-600" />
            </div>
            <h2 className="text-lg font-bold text-slate-900 mb-2">의사를 연결하고 있습니다</h2>
            <p className="text-sm text-slate-500">잠시만 기다려 주세요...</p>
          </div>
        )}

        {/* 화상 통화 (시뮬레이션) */}
        {step === 'call' && selectedDoctor && (
          <div className="rounded-2xl border border-slate-200 bg-slate-900 overflow-hidden">
            {/* 영상 영역 */}
            <div className="relative aspect-video bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center">
              <div className="text-center">
                <span className="text-6xl mb-4 block">{selectedDoctor.photo}</span>
                <p className="text-white font-bold">{selectedDoctor.name} 전문의</p>
                <p className="text-slate-400 text-sm">화상 진료 중</p>
                <div className="flex items-center gap-1 mt-2 justify-center">
                  <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
                  <span className="text-xs text-emerald-400">연결됨</span>
                </div>
              </div>
              {/* 내 화면 (작은 창) */}
              <div className="absolute bottom-4 right-4 w-32 h-24 rounded-xl bg-slate-700 border-2 border-slate-600 flex items-center justify-center">
                <span className="text-sm text-slate-400">내 화면</span>
              </div>
            </div>

            {/* 하단 컨트롤 */}
            <div className="flex items-center justify-center gap-4 p-4 bg-slate-800">
              <button aria-label="카메라 켜기/끄기" className="h-12 w-12 rounded-full bg-slate-700 text-white flex items-center justify-center hover:bg-slate-600">
                <Video className="h-5 w-5" aria-hidden="true" />
              </button>
              <button aria-label="마이크 켜기/끄기" className="h-12 w-12 rounded-full bg-slate-700 text-white flex items-center justify-center hover:bg-slate-600">
                <Phone className="h-5 w-5" aria-hidden="true" />
              </button>
              <button
                onClick={() => { setStep('search'); setSelectedDoctor(null); }}
                className="h-12 w-12 rounded-full bg-red-600 text-white flex items-center justify-center hover:bg-red-700"
              >
                <Phone className="h-5 w-5 rotate-[135deg]" />
              </button>
            </div>

            {/* 측정 데이터 사이드바 */}
            <div className="p-4 bg-slate-800 border-t border-slate-700">
              <p className="text-xs text-slate-400 mb-2 flex items-center gap-1">
                <CheckCircle className="h-3 w-3 text-emerald-500" /> 공유된 측정 데이터
              </p>
              <div className="grid grid-cols-3 gap-2">
                {[
                  { label: '혈당', value: '95 mg/dL', status: '정상' },
                  { label: '당화혈색소', value: '6.8%', status: '주의' },
                  { label: '콜레스테롤', value: '195 mg/dL', status: '정상' },
                ].map(d => (
                  <div key={d.label} className="rounded-lg bg-slate-700 p-2 text-center">
                    <p className="text-[10px] text-slate-400">{d.label}</p>
                    <p className="text-sm font-bold text-white">{d.value}</p>
                    <p className={`text-[10px] ${d.status === '정상' ? 'text-emerald-400' : 'text-amber-400'}`}>{d.status}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-8">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: '건강 기록', path: '/domains/health-records', desc: '진료 기록 확인' },
              { name: '건강 측정', path: '/domains/measure', desc: '진료 전 측정' },
              { name: '건강 데이터', path: '/domains/clinical', desc: '데이터 관리' },
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
