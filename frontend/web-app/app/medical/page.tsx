"use client";

import React from 'react';

export default function MedicalPage() {
    return (
        <div className="min-h-screen bg-[var(--deep-sea)] p-6 text-[var(--hanji-white)] selection:bg-[var(--sanggam-gold)]/30">
            <header className="max-w-7xl mx-auto flex justify-between items-center mb-10">
                <div>
                    <h1 className="text-3xl font-serif shimmer-text">Telemedicine HUD</h1>
                    <p className="data-label text-[10px]">Secure Medical Connection</p>
                </div>
                <a href="/" className="sanggam-btn sanggam-btn-jagae text-xs px-6">DASHBOARD</a>
            </header>

            <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Video Column */}
                <div className="lg:col-span-2 space-y-6">
                    <div className="relative aspect-video sanggam-panel overflow-hidden border-2 border-[var(--sanggam-gold)]/40">
                        <div className="absolute inset-0 bg-neutral-900 flex items-center justify-center text-9xl opacity-20">👩‍⚕️</div>
                        <div className="absolute top-6 left-6 flex items-center gap-3">
                            <div className="w-3 h-3 rounded-full bg-red-500 animate-pulse"></div>
                            <span className="text-[10px] font-mono tracking-widest text-white">REC ENCRYPTED</span>
                        </div>
                        <div className="absolute bottom-8 left-1/2 -translate-x-1/2">
                            <div className="flex gap-6 p-4 sanggam-panel bg-black/60 rounded-full border border-white/10 backdrop-blur-xl">
                                <button className="w-10 h-10 rounded-full border border-white/20 flex items-center justify-center">🎤</button>
                                <button className="w-10 h-10 rounded-full border border-white/20 flex items-center justify-center">📷</button>
                                <button className="px-8 rounded-full bg-red-600 font-serif text-sm">EXIT</button>
                            </div>
                        </div>
                    </div>
                    <div className="grid grid-cols-3 gap-6">
                        <div className="premium-card text-center py-4">
                            <div className="data-label text-[9px] mb-1">HR</div>
                            <div className="text-2xl font-serif text-white">78</div>
                        </div>
                        <div className="premium-card text-center py-4">
                            <div className="data-label text-[9px] mb-1">Temp</div>
                            <div className="text-2xl font-serif text-white">36.5</div>
                        </div>
                        <div className="premium-card text-center py-4">
                            <div className="data-label text-[9px] mb-1">Stress</div>
                            <div className="text-2xl font-serif text-[var(--wave-cyan)]">LOW</div>
                        </div>
                    </div>
                </div>

                {/* Records Column */}
                <div className="premium-card flex flex-col h-full border-[var(--sanggam-gold)]/30">
                    <h3 className="text-xl font-serif mb-6">진료 기록</h3>
                    <div className="flex-1 space-y-4">
                        {[1, 2].map((_, i) => (
                            <div key={i} className="p-4 rounded-lg bg-white/5 border border-white/5">
                                <div className="data-label text-[9px] mb-2 opacity-50">2026.02.22</div>
                                <p className="text-xs text-[var(--hanji-white)]/80">정기 검진 및 건강 상태 확인 완료.</p>
                            </div>
                        ))}
                    </div>
                    <button className="w-full sanggam-btn sanggam-btn-gold mt-8 py-4">처방전 발행하기</button>
                </div>
            </div>
        </div>
    );
}
