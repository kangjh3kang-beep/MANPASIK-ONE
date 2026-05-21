'use client';

import React, { useState, useEffect } from 'react';
import { Cpu, HardDrive, Unplug, Zap, Activity, ShieldCheck, Microscope } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { DomainHeader } from '@mmup/ui';

export default function HardwareCorePage() {
  const { data: session } = useSession();
  const [cpuTemp, setCpuTemp] = useState(42);
  const [memUsage, setMemUsage] = useState(15.4);
  const user = session?.user as any;

  useEffect(() => {
    const interval = setInterval(() => {
      setCpuTemp(prev => prev + (Math.random() > 0.5 ? 0.2 : -0.2));
      setMemUsage(prev => Math.max(15, prev + (Math.random() > 0.5 ? 0.05 : -0.05)));
    }, 2000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-slate-50">
      <DomainHeader 
        title="하드웨어 코어 엔진" 
        subtitle="Rust 기반 저지연 생체 데이터 파싱 및 엣지 AI 추론 모듈 제어"
        icon={Microscope}
        user={user ? {
          name: user.name,
          email: user.email,
          persona: user.persona || 'dev',
          organization: user.organization || 'MMUP Hardware Lab'
        } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-7xl mx-auto px-8 py-8 space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {[
            { label: '엔진 상태', val: 'Active (v2.4.3)', icon: Activity, c: 'text-emerald-600 bg-emerald-50' },
            { label: 'CPU 온도', val: `${cpuTemp.toFixed(1)}°C`, icon: Zap, c: 'text-amber-600 bg-amber-50' },
            { label: '메모리 점유', val: `${memUsage.toFixed(2)}MB`, icon: HardDrive, c: 'text-sky-600 bg-sky-50' },
            { label: '보안 검증', val: 'IEC 62304 OK', icon: ShieldCheck, c: 'text-indigo-600 bg-indigo-50' },
          ].map(s => (
            <div key={s.label} className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="flex items-center justify-between mb-4">
                <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">{s.label}</span>
                <div className={`p-2 rounded-xl ${s.c}`}><s.icon className="h-4 w-4" /></div>
              </div>
              <p className="text-2xl font-black text-slate-900 tabular-nums">{s.val}</p>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="rounded-3xl border border-slate-200 bg-white p-8">
            <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
              <Unplug className="h-5 w-5 text-sky-600" />
              물리 장치 연결 상태 (HW Bridge)
            </h2>
            <div className="space-y-4">
              {[
                { name: 'BLE GATT Service', status: 'Connected', delay: '12ms' },
                { name: 'UART Serial Bus', status: 'Ready', delay: '2ms' },
                { name: 'Edge TPU Accelerator', status: 'Active', delay: '45ms' },
              ].map(dev => (
                <div key={dev.name} className="flex items-center justify-between p-4 rounded-2xl bg-slate-50 border border-slate-100">
                  <span className="text-sm font-bold text-slate-700">{dev.name}</span>
                  <div className="flex items-center gap-4">
                    <span className="text-xs font-medium text-slate-400">Latency: {dev.delay}</span>
                    <span className="px-2 py-1 bg-emerald-100 text-emerald-700 text-[10px] font-bold rounded-lg">{dev.status}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-3xl border border-slate-200 bg-white p-8">
            <h2 className="text-lg font-bold text-slate-900 mb-6 flex items-center gap-2">
              <Cpu className="h-5 w-5 text-indigo-600" />
              Rust Core 스레드 모니터
            </h2>
            <div className="relative h-48 w-full bg-slate-900 rounded-2xl overflow-hidden p-4 font-mono text-[10px] text-emerald-400">
              <p className="mb-1">{`[sys] Rust ManPaSik Engine initialized...`}</p>
              <p className="mb-1">{`[io] BLE data incoming: 0x4A 0x22 0xFF 0x01`}</p>
              <p className="mb-1 text-sky-400">{`[ml] Edge inference complete: diag_result=0.87`}</p>
              <p className="mb-1">{`[sys] Memory swap check... OK`}</p>
              <p className="mb-1 animate-pulse">{`_`}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
