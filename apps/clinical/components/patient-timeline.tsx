'use client';

import React, { useState, useEffect } from 'react';
import { TrendingUp, TrendingDown, Minus, Heart, Thermometer, Droplets, Wind, Zap } from 'lucide-react';

interface VitalRecord {
  time: string;
  hr: number;
  spo2: number;
  temp: number;
  bp_sys: number;
  bp_dia: number;
  rr: number;
}

function generateMockVitals(): VitalRecord[] {
  const now = new Date();
  return Array.from({ length: 24 }, (_, i) => {
    const t = new Date(now.getTime() - (23 - i) * 3600000);
    return {
      time: `${t.getHours().toString().padStart(2, '0')}:00`,
      hr: 68 + Math.round(Math.random() * 20 - 10),
      spo2: 96 + Math.round(Math.random() * 4),
      temp: +(36.2 + Math.random() * 1.2).toFixed(1),
      bp_sys: 115 + Math.round(Math.random() * 20 - 10),
      bp_dia: 72 + Math.round(Math.random() * 14 - 7),
      rr: 14 + Math.round(Math.random() * 6 - 3),
    };
  });
}

function MiniSparkline({ data, color }: { data: number[]; color: string }) {
  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;
  const h = 32;
  const w = 120;
  const points = data.map((v, i) => `${(i / (data.length - 1)) * w},${h - ((v - min) / range) * h}`).join(' ');

  return (
    <svg width={w} height={h} className="inline-block">
      <polyline fill="none" stroke={color} strokeWidth="1.5" strokeLinejoin="round" points={points} />
    </svg>
  );
}

const VITAL_CARDS = [
  { key: 'hr', label: '심박수', unit: 'bpm', icon: Heart, color: '#ef4444', normalMin: 60, normalMax: 100 },
  { key: 'spo2', label: '산소포화도', unit: '%', icon: Droplets, color: '#3b82f6', normalMin: 95, normalMax: 100 },
  { key: 'temp', label: '체온', unit: '°C', icon: Thermometer, color: '#f59e0b', normalMin: 36.1, normalMax: 37.2 },
  { key: 'bp_sys', label: '수축기 혈압', unit: 'mmHg', icon: Zap, color: '#8b5cf6', normalMin: 90, normalMax: 140 },
  { key: 'rr', label: '호흡수', unit: '/min', icon: Wind, color: '#10b981', normalMin: 12, normalMax: 20 },
] as const;

export function PatientTimeline() {
  const [vitals, setVitals] = useState<VitalRecord[]>([]);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    setVitals(generateMockVitals());
    const interval = setInterval(() => {
      setVitals((prev) => {
        const next = [...prev.slice(1)];
        const now = new Date();
        next.push({
          time: `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`,
          hr: 68 + Math.round(Math.random() * 20 - 10),
          spo2: 96 + Math.round(Math.random() * 4),
          temp: +(36.2 + Math.random() * 1.2).toFixed(1),
          bp_sys: 115 + Math.round(Math.random() * 20 - 10),
          bp_dia: 72 + Math.round(Math.random() * 14 - 7),
          rr: 14 + Math.round(Math.random() * 6 - 3),
        });
        return next;
      });
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  if (!mounted) return <div className="h-48 animate-pulse rounded-2xl bg-slate-100" />;

  const latest = vitals[vitals.length - 1];
  const prev = vitals[vitals.length - 2];

  return (
    <div>
      <h2 className="text-lg font-bold text-slate-900 mb-4">실시간 생체 신호 모니터링</h2>
      <div className="grid grid-cols-2 lg:grid-cols-5 gap-4">
        {VITAL_CARDS.map((card) => {
          if (!latest) return null;
          const val = latest[card.key as keyof VitalRecord] as number;
          const prevVal = prev ? (prev[card.key as keyof VitalRecord] as number) : val;
          const diff = val - prevVal;
          const isNormal = val >= card.normalMin && val <= card.normalMax;
          const dataArr = vitals.map((v) => v[card.key as keyof VitalRecord] as number);

          return (
            <div
              key={card.key}
              className={`rounded-2xl border p-4 transition-all ${
                isNormal
                  ? 'border-slate-200 bg-white'
                  : 'border-red-200 bg-red-50 animate-pulse'
              }`}
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <card.icon className="h-4 w-4" style={{ color: card.color }} />
                  <span className="text-xs font-semibold text-slate-500">{card.label}</span>
                </div>
                {!isNormal && (
                  <span className="text-[10px] font-bold text-red-600 bg-red-100 px-1.5 py-0.5 rounded-full">이상</span>
                )}
              </div>

              <div className="flex items-end gap-1.5">
                <span className="text-2xl font-bold tabular-nums" style={{ color: isNormal ? '#0f172a' : '#dc2626' }}>
                  {val}
                </span>
                <span className="text-xs text-slate-400 mb-1">{card.unit}</span>
              </div>

              <div className="flex items-center gap-1 mt-1">
                {diff > 0 ? (
                  <TrendingUp className="h-3 w-3 text-red-500" />
                ) : diff < 0 ? (
                  <TrendingDown className="h-3 w-3 text-emerald-500" />
                ) : (
                  <Minus className="h-3 w-3 text-slate-400" />
                )}
                <span className={`text-[11px] font-medium ${diff > 0 ? 'text-red-500' : diff < 0 ? 'text-emerald-500' : 'text-slate-400'}`}>
                  {diff > 0 ? '+' : ''}{diff.toFixed(card.key === 'temp' ? 1 : 0)}
                </span>
              </div>

              <div className="mt-3">
                <MiniSparkline data={dataArr} color={card.color} />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
