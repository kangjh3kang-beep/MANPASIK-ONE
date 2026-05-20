"use client";

import React from 'react';

export default function PremiumDashboard() {
    const metrics = [
        { label: '심박수', value: '72', unit: 'BPM', pos: 'top-10 left-10' },
        { label: '산소포화도', value: '98', unit: '%', pos: 'top-10 right-10' },
        { label: '혈당', value: '105', unit: 'mg/dL', pos: 'bottom-10 left-10' },
        { label: '뇌활동', value: 'Alpha', unit: 'Hz', pos: 'bottom-10 right-10' },
    ];

    return (
        <div className="relative min-h-screen overflow-hidden p-8 flex items-center justify-center">
            {/* Background Patterns */}
            <div className="absolute inset-0 opacity-10 pointer-events-none">
                <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
                    <defs>
                        <pattern id="lattice" width="40" height="40" patternUnits="userSpaceOnUse">
                            <path d="M 40 0 L 0 0 0 40" fill="none" stroke="var(--sanggam-gold)" strokeWidth="0.5" />
                        </pattern>
                    </defs>
                    <rect width="100%" height="100%" fill="url(#lattice)" />
                </svg>
            </div>

            <div className="z-10 w-full max-w-5xl relative h-[600px]">
                {/* Central Anchor: Golden Globe */}
                <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
                    <div className="globe-container">
                        <div className="golden-globe">
                            <div className="globe-wire wire-h"></div>
                            <div className="globe-wire wire-v"></div>
                        </div>
                        {/* Aura */}
                        <div className="absolute inset-0 rounded-full bg-sanggam-gold/5 blur-3xl animate-pulse"></div>
                    </div>
                    <div className="mt-8 text-center">
                        <h2 className="text-3xl font-serif shimmer-text">Sanggam Orbit System</h2>
                        <p className="data-label mt-2">Global Health Integration Active</p>
                    </div>
                </div>

                {/* Orbital Nodes */}
                {metrics.map((m, i) => (
                    <div key={i} className={`absolute ${m.pos} w-48 transition-all duration-700 hover:scale-105`}>
                        <div className="premium-card group">
                            <div className="flex justify-between items-start mb-4">
                                <span className="data-label text-xs">{m.label}</span>
                                <div className="w-2 h-2 rounded-full bg-wave-cyan animate-ping"></div>
                            </div>
                            <div className="flex items-baseline gap-1">
                                <span className="text-4xl font-serif text-white">{m.value}</span>
                                <span className="text-xs text-wave-cyan/60">{m.unit}</span>
                            </div>
                            <div className="mt-4 h-1 w-full bg-white/5 rounded-full overflow-hidden">
                                <div className="h-full bg-sanggam-gold w-3/4 shadow-[0_0_10px_var(--sanggam-gold-glow)]"></div>
                            </div>

                            {/* Dynamic Tether Link (Pseudo) */}
                            <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-px h-16 bg-gradient-to-t from-sanggam-gold/20 to-transparent"></div>
                        </div>
                    </div>
                ))}
            </div>

            {/* Footer System Info */}
            <div className="absolute bottom-10 left-10 right-10 flex justify-between items-end border-t border-sanggam-gold/20 pt-4">
                <div className="flex gap-8">
                    <div>
                        <div className="data-label text-[10px] opacity-40">System Core</div>
                        <div className="text-sm font-serif text-sanggam-gold">V6.2.0-ALBION</div>
                    </div>
                    <div>
                        <div className="data-label text-[10px] opacity-40">Security Status</div>
                        <div className="text-sm font-serif text-wave-cyan">ENCRYPTED</div>
                    </div>
                </div>
                <div className="text-right">
                    <button className="sanggam-btn text-sm">ACCESS PROTOCOL</button>
                </div>
            </div>
        </div>
    );
}
