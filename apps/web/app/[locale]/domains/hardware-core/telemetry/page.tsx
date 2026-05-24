'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Monitor, CheckCircle, AlertTriangle, Activity, Heart, Zap, Wifi, TrendingUp, Shield } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type HealthStatus = '정상' | '주의' | '오류';
type AlertSeverity = '경고' | '주의' | '위험';

interface SensorInfo {
  name: string;
  status: HealthStatus;
  signalQuality: number;
  batteryLevel: number;
  lastReading: string;
}

interface AlertItem {
  id: string;
  title: string;
  severity: AlertSeverity;
  timestamp: string;
  description: string;
}

const STATUS_DOT: Record<HealthStatus, string> = {
  '정상': 'bg-emerald-500',
  '주의': 'bg-amber-500',
  '오류': 'bg-red-500',
};

const STATUS_CARD: Record<HealthStatus, { bg: string; text: string; border: string }> = {
  '정상': { bg: 'bg-emerald-50', text: 'text-emerald-700', border: 'border-emerald-200' },
  '주의': { bg: 'bg-amber-50', text: 'text-amber-700', border: 'border-amber-200' },
  '오류': { bg: 'bg-red-50', text: 'text-red-700', border: 'border-red-200' },
};

const SEVERITY_STYLES: Record<AlertSeverity, string> = {
  '주의': 'bg-amber-50 text-amber-700 border-amber-200',
  '경고': 'bg-sky-50 text-sky-700 border-sky-200',
  '위험': 'bg-red-50 text-red-700 border-red-200',
};

const DEVICE_OVERVIEW: { status: HealthStatus; count: number }[] = [
  { status: '정상', count: 45 },
  { status: '주의', count: 3 },
  { status: '오류', count: 1 },
];

const SENSORS: SensorInfo[] = [
  { name: '광학 센서 (LED)', status: '정상', signalQuality: 98, batteryLevel: 87, lastReading: '2026-05-24 09:30' },
  { name: '전기화학 센서', status: '정상', signalQuality: 95, batteryLevel: 92, lastReading: '2026-05-24 09:28' },
  { name: '온도 센서', status: '주의', signalQuality: 78, batteryLevel: 65, lastReading: '2026-05-24 09:25' },
  { name: 'BLE 통신 모듈', status: '정상', signalQuality: 91, batteryLevel: 88, lastReading: '2026-05-24 09:30' },
  { name: 'NFC 리더', status: '정상', signalQuality: 94, batteryLevel: 90, lastReading: '2026-05-24 09:29' },
  { name: '가속도계', status: '오류', signalQuality: 32, batteryLevel: 45, lastReading: '2026-05-24 08:15' },
];

const CPU_READINGS = [42, 45, 48, 52, 47, 44];
const CPU_LABELS = ['09:00', '09:05', '09:10', '09:15', '09:20', '09:25'];

const ALERTS: AlertItem[] = [
  { id: 'alt-1', title: '센서 과열 감지', severity: '경고', timestamp: '2026-05-24 09:18', description: '온도 센서 판독값이 임계치(55°C)에 근접합니다.' },
  { id: 'alt-2', title: '배터리 부족 경고', severity: '주의', timestamp: '2026-05-24 08:45', description: '가속도계 모듈 배터리가 45% 이하입니다. 충전이 필요합니다.' },
  { id: 'alt-3', title: '신호 약함 감지', severity: '위험', timestamp: '2026-05-24 08:15', description: '가속도계 신호 품질이 32%로 정상 범위 이하입니다. 센서 점검이 필요합니다.' },
];

export default function TelemetryPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [diagnosing, setDiagnosing] = useState<string | null>(null);

  const handleDiagnose = (sensorName: string) => {
    setDiagnosing(sensorName);
    setTimeout(() => setDiagnosing(null), 2000);
  };

  const memoryUsage = 68;
  const storageUsage = 42;
  const cpuMax = Math.max(...CPU_READINGS);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="텔레메트리 대시보드">
      <DomainHeader
        title="텔레메트리 대시보드"
        subtitle="실시간 장비 상태와 센서 데이터를 모니터링합니다"
        icon={Monitor}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'admin' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8 space-y-6">
        {/* 장비 건강 요약 */}
        <section aria-label="장비 건강 요약">
          <h2 className="text-sm font-bold text-slate-900 mb-3">장비 건강 현황</h2>
          <div className="grid grid-cols-3 gap-4">
            {DEVICE_OVERVIEW.map(item => {
              const style = STATUS_CARD[item.status];
              return (
                <div key={item.status} className={`rounded-2xl border ${style.border} ${style.bg} p-5 text-center`}>
                  <p className={`text-3xl font-black ${style.text}`}>{item.count}</p>
                  <p className={`text-xs font-bold mt-1 ${style.text}`}>{item.status}</p>
                </div>
              );
            })}
          </div>
        </section>

        {/* 센서 건강 그리드 */}
        <section aria-label="센서 상태">
          <h2 className="text-sm font-bold text-slate-900 mb-3">센서 상태</h2>
          <div className="grid grid-cols-2 gap-3">
            {SENSORS.map(sensor => (
              <div key={sensor.name} className="rounded-2xl border border-slate-200 bg-white p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <span className={`h-2.5 w-2.5 rounded-full ${STATUS_DOT[sensor.status]}`} aria-label={sensor.status} />
                    <h3 className="text-xs font-bold text-slate-900">{sensor.name}</h3>
                  </div>
                  <button
                    onClick={() => handleDiagnose(sensor.name)}
                    className="px-2 py-1 rounded-lg border border-slate-200 text-[10px] font-semibold text-slate-500 hover:bg-slate-50 transition-colors"
                    aria-label={`${sensor.name} 원격 진단`}
                  >
                    {diagnosing === sensor.name ? '진단 중...' : '원격 진단'}
                  </button>
                </div>

                <div className="grid grid-cols-3 gap-2">
                  <div className="rounded-xl bg-slate-50 p-2 text-center">
                    <p className="text-[10px] text-slate-400 font-bold mb-0.5">신호 품질</p>
                    <p className={`text-sm font-black ${sensor.signalQuality >= 80 ? 'text-emerald-600' : sensor.signalQuality >= 50 ? 'text-amber-600' : 'text-red-600'}`}>
                      {sensor.signalQuality}%
                    </p>
                  </div>
                  <div className="rounded-xl bg-slate-50 p-2 text-center">
                    <p className="text-[10px] text-slate-400 font-bold mb-0.5">배터리</p>
                    <p className={`text-sm font-black ${sensor.batteryLevel >= 60 ? 'text-emerald-600' : sensor.batteryLevel >= 30 ? 'text-amber-600' : 'text-red-600'}`}>
                      {sensor.batteryLevel}%
                    </p>
                  </div>
                  <div className="rounded-xl bg-slate-50 p-2 text-center">
                    <p className="text-[10px] text-slate-400 font-bold mb-0.5">마지막</p>
                    <p className="text-[10px] font-semibold text-slate-600">{sensor.lastReading.split(' ')[1]}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* 시스템 메트릭 */}
        <section aria-label="시스템 메트릭">
          <h2 className="text-sm font-bold text-slate-900 mb-3">시스템 메트릭</h2>
          <div className="rounded-2xl border border-slate-200 bg-white p-5 space-y-5">
            {/* CPU 온도 추세 */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-bold text-slate-600 flex items-center gap-1">
                  <TrendingUp className="h-3.5 w-3.5 text-sky-500" /> CPU 온도 추세
                </span>
                <span className="text-xs font-mono text-slate-500">{CPU_READINGS[CPU_READINGS.length - 1]}°C</span>
              </div>
              <div className="flex items-end gap-1 h-16">
                {CPU_READINGS.map((val, i) => (
                  <div key={i} className="flex-1 flex flex-col items-center gap-1">
                    <div
                      className={`w-full rounded-t-md ${val > 50 ? 'bg-amber-400' : 'bg-sky-400'}`}
                      style={{ height: `${(val / (cpuMax + 10)) * 100}%` }}
                      title={`${val}°C`}
                    />
                    <span className="text-[8px] text-slate-400">{CPU_LABELS[i]}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* 메모리 사용률 */}
            <div>
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs font-bold text-slate-600">메모리 사용률</span>
                <span className="text-xs font-bold text-slate-700">{memoryUsage}%</span>
              </div>
              <div className="h-2.5 bg-slate-100 rounded-full overflow-hidden">
                <div className="h-full bg-sky-500 rounded-full" style={{ width: `${memoryUsage}%` }} role="progressbar" aria-valuenow={memoryUsage} aria-valuemin={0} aria-valuemax={100} />
              </div>
            </div>

            {/* 저장공간 사용률 */}
            <div>
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs font-bold text-slate-600">저장공간 사용률</span>
                <span className="text-xs font-bold text-slate-700">{storageUsage}%</span>
              </div>
              <div className="h-2.5 bg-slate-100 rounded-full overflow-hidden">
                <div className="h-full bg-emerald-500 rounded-full" style={{ width: `${storageUsage}%` }} role="progressbar" aria-valuenow={storageUsage} aria-valuemin={0} aria-valuemax={100} />
              </div>
            </div>
          </div>
        </section>

        {/* 알림 */}
        <section aria-label="최근 알림">
          <h2 className="text-sm font-bold text-slate-900 mb-3">최근 알림</h2>
          <div className="space-y-2">
            {ALERTS.map(alert => (
              <div key={alert.id} className={`rounded-2xl border p-4 ${SEVERITY_STYLES[alert.severity]}`} role="alert">
                <div className="flex items-center justify-between mb-1">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-4 w-4" />
                    <h3 className="text-xs font-bold">{alert.title}</h3>
                  </div>
                  <span className="text-[10px] font-mono">{alert.timestamp}</span>
                </div>
                <p className="text-[11px] ml-6">{alert.description}</p>
              </div>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
