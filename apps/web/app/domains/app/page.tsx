'use client';
import React, { useState, useEffect } from 'react';
import { Activity, Stethoscope, Heart, Droplets, Footprints, Moon, Smartphone, ChevronRight, Users2, Database } from 'lucide-react';
import Link from 'next/link';

const HEALTH_CARDS = [
  { label: '심박수', val: '72', unit: 'bpm', icon: Heart, color: 'text-red-500 bg-red-50', trend: '정상' },
  { label: '산소포화도', val: '98', unit: '%', icon: Droplets, color: 'text-sky-500 bg-sky-50', trend: '정상' },
  { label: '걸음 수', val: '8,247', unit: '보', icon: Footprints, color: 'text-emerald-500 bg-emerald-50', trend: '목표 82%' },
  { label: '수면', val: '7h 23m', unit: '', icon: Moon, color: 'text-violet-500 bg-violet-50', trend: '양호' },
];

const RECENT = [
  { time: '14:30', type: '혈압', value: '122/78 mmHg' },
  { time: '12:00', type: '혈당', value: '95 mg/dL' },
  { time: '09:15', type: '체온', value: '36.5 °C' },
  { time: '08:00', type: '체중', value: '68.2 kg' },
];

export default function MobileAppPage() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="min-h-screen bg-slate-50 animate-pulse" />;

  const lastUpdate = new Date().toLocaleTimeString();

  return (
    <div className="min-h-screen bg-gradient-to-b from-sky-50 to-white">
      <div className="border-b border-slate-200 bg-white/80 backdrop-blur px-8 py-6">
        <div className="flex items-center gap-3 mb-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-500 text-white">
            <Smartphone className="h-5 w-5" />
          </div>
          <h1 className="text-2xl font-bold text-slate-900">만파식 건강 앱 v2.0</h1>
        </div>
        <p className="text-sm text-slate-500">나의 건강 데이터를 한눈에 — PWA 모바일 뷰 (Last Sync: {lastUpdate})</p>
      </div>

      <div className="max-w-lg mx-auto px-6 py-8 space-y-6">
        <div className="rounded-2xl bg-gradient-to-r from-cyan-500 to-sky-500 p-6 text-white shadow-lg shadow-cyan-500/20">
          <p className="text-sm opacity-80">안녕하세요, 김환자님 👋</p>
          <p className="text-xl font-bold mt-1">오늘도 건강한 하루 보내세요!</p>
          <div className="mt-3 flex items-center gap-2 text-sm opacity-90">
            <Activity className="h-4 w-4" />
            <span>건강 점수: <strong>87/100</strong> (안정)</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-3">
          {HEALTH_CARDS.map(card => (
            <div key={card.label} className="rounded-2xl border border-slate-200 bg-white p-4">
              <div className="flex items-center gap-2 mb-2">
                <div className={`flex h-8 w-8 items-center justify-center rounded-xl ${card.color}`}>
                  <card.icon className="h-4 w-4" />
                </div>
                <span className="text-xs font-semibold text-slate-500">{card.label}</span>
              </div>
              <p className="text-xl font-bold text-slate-900">{card.val} <span className="text-xs font-normal text-slate-400">{card.unit}</span></p>
              <p className="text-[11px] text-emerald-600 mt-0.5">{card.trend}</p>
            </div>
          ))}
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <h2 className="text-sm font-bold text-slate-900 mb-3">최근 측정 기록</h2>
          <div className="space-y-2">
            {RECENT.map((m, i) => (
              <div key={i} className="flex items-center justify-between py-2 border-b border-slate-50 last:border-0">
                <div className="flex items-center gap-3">
                  <span className="text-xs text-slate-400 font-mono w-12">{m.time}</span>
                  <span className="text-sm font-medium text-slate-700">{m.type}</span>
                </div>
                <span className="text-sm font-bold text-slate-900">{m.value}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="space-y-3">
          {[
            { label: '새 측정 시작', href: 'http://localhost:3001/', icon: Stethoscope },
            { label: '의사에게 공유', href: 'http://localhost:3005/', icon: Users2 },
            { label: '리워드 확인', href: 'http://localhost:3004/', icon: Database },
          ].map((action) => (
            <a 
              key={action.label} 
              href={action.href}
              className="flex w-full items-center justify-between rounded-2xl border border-slate-200 bg-white px-5 py-4 text-sm font-bold text-slate-700 hover:bg-slate-50 hover:border-sky-500 hover:text-sky-600 transition-all group shadow-sm"
            >
              <div className="flex items-center gap-3">
                <div className="h-8 w-8 rounded-lg bg-slate-50 flex items-center justify-center group-hover:bg-sky-50 transition-colors">
                  <action.icon className="h-4 w-4" />
                </div>
                {action.label}
              </div>
              <ChevronRight className="h-4 w-4 text-slate-300 group-hover:translate-x-1 transition-transform" />
            </a>
          ))}
        </div>

      </div>
    </div>
  );
}
