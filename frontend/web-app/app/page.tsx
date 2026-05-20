"use client";

import React from 'react';

export default function Home() {
  const metrics = [
    { label: '심박수', value: '74', unit: 'BPM', pos: 'top-10 left-10' },
    { label: '산소포화도', value: '99', unit: '%', pos: 'top-10 right-10' },
    { label: '혈당', value: '98', unit: 'mg/dL', pos: 'bottom-10 left-10' },
    { label: '뇌활동', value: 'Alpha', unit: 'Hz', pos: 'bottom-10 right-10' },
  ];

  return (
    <div className="relative min-h-screen overflow-hidden p-8 flex flex-col items-center justify-center">
      {/* ═══ Traditional Lattice Background ═══ */}
      <div className="absolute inset-0 opacity-10 pointer-events-none">
        <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <pattern id="lattice" width="60" height="60" patternUnits="userSpaceOnUse">
              <path d="M 60 0 L 0 0 0 60" fill="none" stroke="var(--sanggam-gold)" strokeWidth="0.5" />
              <circle cx="30" cy="30" r="1" fill="var(--sanggam-gold)" />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#lattice)" />
        </svg>
      </div>

      <header className="absolute top-12 text-center z-20">
        <h1 className="text-5xl font-serif shimmer-text mb-2">만파식 전역 통합 관제</h1>
        <p className="data-label text-sm tracking-[0.3em]">GLOBAL DATA HUB V6.2.0</p>
      </header>

      <div className="z-10 w-full max-w-6xl relative h-[650px] flex items-center justify-center">
        {/* Central Anchor: Golden Globe (CSS 3D) */}
        <div className="relative">
          <div className="globe-container">
            <div className="golden-globe">
              <div className="globe-wire wire-h"></div>
              <div className="globe-wire wire-v"></div>
            </div>
            {/* Energy Aura */}
            <div className="absolute inset-0 rounded-full bg-[var(--sanggam-gold-glow)] blur-[80px] animate-pulse"></div>
          </div>
          <div className="mt-12 text-center">
            <h2 className="text-2xl font-serif text-white/90">Sanggam Orbit System</h2>
            <div className="sanggam-glow-line w-48 mx-auto mt-2" />
          </div>
        </div>

        {/* Orbital Nodes Mapping */}
        {metrics.map((m, i) => (
          <div key={i} className={`absolute ${m.pos} w-52 transition-all duration-1000 hover:scale-105`}>
            <div className="premium-card group">
              <div className="flex justify-between items-start mb-4">
                <span className="data-label text-[10px] opacity-70">{m.label}</span>
                <div className="w-2 h-2 rounded-full bg-[var(--wave-cyan)] animate-ping"></div>
              </div>
              <div className="flex items-baseline gap-1">
                <span className="text-4xl font-serif text-white">{m.value}</span>
                <span className="text-xs text-[var(--wave-cyan)] opacity-60">{m.unit}</span>
              </div>
              <div className="mt-4 h-1 w-full bg-white/5 rounded-full overflow-hidden">
                <div
                  className="h-full bg-[var(--sanggam-gold)] shadow-[0_0_15px_var(--sanggam-gold-glow)]"
                  style={{ width: '75%' }}
                ></div>
              </div>

              {/* Virtual Connection Line */}
              <div className="absolute -bottom-4 left-1/2 -translate-x-1/2 w-px h-24 bg-gradient-to-t from-[var(--sanggam-gold)]/20 to-transparent"></div>
            </div>
          </div>
        ))}
      </div>

      {/* Footer Navigation Overlay */}
      <footer className="absolute bottom-12 w-full max-w-7xl px-8 flex justify-between items-end border-t border-white/5 pt-6 z-20">
        <div className="flex gap-10">
          <div>
            <div className="data-label text-[10px] opacity-40">System Core</div>
            <div className="text-sm font-serif text-[var(--sanggam-gold)]">ACTIVE-ALBION</div>
          </div>
          <div>
            <div className="data-label text-[10px] opacity-40">Security Status</div>
            <div className="text-sm font-serif text-[var(--wave-cyan)]">LEVEL 9 ENCRYPTED</div>
          </div>
        </div>

        <nav className="flex gap-6">
          <a href="/market" className="sanggam-btn sanggam-btn-gold text-xs px-8">MARKET</a>
          <a href="/medical" className="sanggam-btn sanggam-btn-jagae text-xs px-8">MEDICAL</a>
          <a href="/exercise" className="sanggam-btn sanggam-btn-jagae text-xs px-8">EXERCISE</a>
        </nav>
      </footer>
    </div>
  );
}

