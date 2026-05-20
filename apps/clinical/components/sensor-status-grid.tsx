'use client';

import React, { useState, useEffect } from 'react';
import { Wifi, WifiOff, Battery, BatteryWarning, Cpu, Clock, Signal } from 'lucide-react';

interface SensorDevice {
  id: string;
  name: string;
  type: string;
  status: 'online' | 'offline' | 'warning';
  battery: number;
  signal: number;
  lastSeen: string;
  firmware: string;
}

function generateSensors(): SensorDevice[] {
  const types = ['NIR 분광계', 'PPG 센서', '가속도계', '온도 센서', 'ECG 리드', '압전 센서'];
  const names = ['SEN-A01', 'SEN-A02', 'SEN-B01', 'SEN-B02', 'SEN-C01', 'SEN-C02', 'SEN-D01', 'SEN-D02'];

  return names.map((name, i) => {
    const rand = Math.random();
    const status: SensorDevice['status'] = rand > 0.85 ? 'offline' : rand > 0.7 ? 'warning' : 'online';
    return {
      id: `dev_${name.toLowerCase().replace('-', '')}`,
      name,
      type: types[i % types.length],
      status,
      battery: status === 'offline' ? 0 : Math.round(20 + Math.random() * 80),
      signal: status === 'offline' ? 0 : Math.round(40 + Math.random() * 60),
      lastSeen: status === 'offline' ? '24분 전' : status === 'warning' ? '3분 전' : '방금',
      firmware: 'v2.4.' + (3 + Math.floor(Math.random() * 5)),
    };
  });
}

const STATUS_STYLES = {
  online: { dot: 'bg-emerald-500', bg: 'border-emerald-100', label: '정상', labelColor: 'text-emerald-700 bg-emerald-50' },
  warning: { dot: 'bg-amber-500', bg: 'border-amber-100', label: '주의', labelColor: 'text-amber-700 bg-amber-50' },
  offline: { dot: 'bg-red-500', bg: 'border-red-200', label: '오프라인', labelColor: 'text-red-700 bg-red-50' },
};

export function SensorStatusGrid() {
  const [sensors, setSensors] = useState<SensorDevice[]>([]);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    setSensors(generateSensors());
  }, []);

  if (!mounted) return <div className="h-48 animate-pulse rounded-2xl bg-slate-100" />;

  const online = sensors.filter((s) => s.status === 'online').length;
  const total = sensors.length;

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-bold text-slate-900">센서 디바이스 상태</h2>
        <div className="flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1.5">
          <span className={`h-2 w-2 rounded-full ${online === total ? 'bg-emerald-500' : 'bg-amber-500'}`} />
          <span className="text-xs font-semibold text-slate-600">
            {online}/{total} 온라인
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {sensors.map((sensor) => {
          const style = STATUS_STYLES[sensor.status];
          return (
            <div
              key={sensor.id}
              className={`rounded-2xl border bg-white p-4 transition-all hover:shadow-md ${style.bg}`}
            >
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-100">
                    <Cpu className="h-4 w-4 text-slate-600" />
                  </div>
                  <div>
                    <p className="text-sm font-bold text-slate-900">{sensor.name}</p>
                    <p className="text-[11px] text-slate-400">{sensor.type}</p>
                  </div>
                </div>
                <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${style.labelColor}`}>
                  {style.label}
                </span>
              </div>

              <div className="grid grid-cols-3 gap-2 text-center">
                <div className="rounded-lg bg-slate-50 p-2">
                  <div className="flex items-center justify-center gap-1 mb-0.5">
                    {sensor.battery > 20 ? (
                      <Battery className="h-3 w-3 text-emerald-500" />
                    ) : (
                      <BatteryWarning className="h-3 w-3 text-red-500" />
                    )}
                  </div>
                  <p className="text-xs font-bold text-slate-700">{sensor.battery}%</p>
                </div>
                <div className="rounded-lg bg-slate-50 p-2">
                  <div className="flex items-center justify-center gap-1 mb-0.5">
                    <Signal className="h-3 w-3 text-sky-500" />
                  </div>
                  <p className="text-xs font-bold text-slate-700">{sensor.signal}%</p>
                </div>
                <div className="rounded-lg bg-slate-50 p-2">
                  <div className="flex items-center justify-center gap-1 mb-0.5">
                    <Clock className="h-3 w-3 text-slate-400" />
                  </div>
                  <p className="text-xs font-bold text-slate-700">{sensor.lastSeen}</p>
                </div>
              </div>

              <p className="text-[10px] text-slate-400 mt-2 text-right">FW {sensor.firmware}</p>
            </div>
          );
        })}
      </div>
    </div>
  );
}
