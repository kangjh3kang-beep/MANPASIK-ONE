"use client";

import React from 'react';

export default function PremiumMedical() {
    return (
        <div className="min-h-screen bg-deep-sea p-6 text-hanji selection:bg-sanggam/30">
            <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column: Video & Controls */}
                <div className="lg:col-span-2 space-y-6">
                    <div className="relative aspect-video sanggam-panel overflow-hidden border-2 border-sanggam/40">
                        {/* Simulated Video Feed */}
                        <div className="absolute inset-0 bg-neutral-900 flex items-center justify-center">
                            <span className="text-sanggam/20 text-9xl font-serif">👩‍⚕️</span>
                            {/* Scan Overlay */}
                            <div className="absolute inset-4 border border-wave-cyan/20 pointer-events-none">
                                <div className="absolute top-0 left-0 w-8 h-8 border-t-2 border-l-2 border-wave-cyan"></div>
                                <div className="absolute top-0 right-0 w-8 h-8 border-t-2 border-r-2 border-wave-cyan"></div>
                                <div className="absolute bottom-0 left-0 w-8 h-8 border-b-2 border-l-2 border-wave-cyan"></div>
                                <div className="absolute bottom-0 right-0 w-8 h-8 border-b-2 border-r-2 border-wave-cyan"></div>
                            </div>
                        </div>

                        <div className="absolute top-6 left-6 flex items-center gap-3">
                            <div className="w-3 h-3 rounded-full bg-red-500 animate-pulse"></div>
                            <span className="text-xs font-mono tracking-widest text-white shadow-sm">REC LIVE 1080P</span>
                        </div>

                        <div className="absolute bottom-8 left-1/2 -translate-x-1/2">
                            <div className="flex gap-6 p-4 sanggam-panel bg-black/60 rounded-full border border-white/10 backdrop-blur-xl">
                                <button className="w-12 h-12 rounded-full border border-white/20 flex items-center justify-center hover:bg-white/10 transition-colors">🎤</button>
                                <button className="w-12 h-12 rounded-full border border-white/20 flex items-center justify-center hover:bg-white/10 transition-colors">📷</button>
                                <button className="w-12 h-12 rounded-full bg-red-600 flex items-center justify-center hover:bg-red-700 transition-all shadow-lg shadow-red-900/40 px-6 font-serif">진료 종료</button>
                                <button className="w-12 h-12 rounded-full border border-white/20 flex items-center justify-center hover:bg-white/10 transition-colors">⚙️</button>
                            </div>
                        </div>
                    </div>

                    <div className="grid grid-cols-3 gap-6">
                        <div className="premium-card text-center py-4">
                            <div className="data-label text-[10px] mb-1">Heart Rate</div>
                            <div className="text-2xl font-serif text-white">78 BPM</div>
                        </div>
                        <div className="premium-card text-center py-4">
                            <div className="data-label text-[10px] mb-1">BP</div>
                            <div className="text-2xl font-serif text-white">120/80</div>
                        </div>
                        <div className="premium-card text-center py-4 border-dancheong/30">
                            <div className="data-label text-[10px] mb-1 text-dancheong">Stress</div>
                            <div className="text-2xl font-serif text-dancheong animate-pulse">HIGH</div>
                        </div>
                    </div>
                </div>

                {/* Right Column: Records & Prescription */}
                <div className="space-y-6">
                    <div className="premium-card flex flex-col h-full border-sanggam/30">
                        <h3 className="text-xl font-serif mb-6 flex items-center gap-2">
                            <span className="text-sanggam">📜</span> 진료 기록부
                        </h3>
                        <div className="flex-1 space-y-4 overflow-y-auto max-h-[400px] pr-2 scrollbar-thin scrollbar-thumb-sanggam/20">
                            {[1, 2, 3].map((_, i) => (
                                <div key={i} className="p-4 rounded-lg bg-white/5 border border-white/5 hover:border-sanggam/20 transition-all">
                                    <div className="flex justify-between text-[10px] data-label mb-2">
                                        <span>2026.02.22</span>
                                        <span className="text-sanggam">Dr. Kim Ji-Won</span>
                                    </div>
                                    <p className="text-sm text-hanji/80 font-serif">계절성 알레르기 증상 완화를 위한 항히스타민제 처방 건...</p>
                                </div>
                            ))}
                        </div>
                        <div className="mt-8 pt-6 border-t border-white/10">
                            <div className="flex items-center gap-4 mb-6">
                                <div className="w-16 h-16 rounded-xl bg-sanggam/10 border border-sanggam/30 flex items-center justify-center text-3xl">💊</div>
                                <div>
                                    <h4 className="text-lg font-serif">처방전 발행 완료</h4>
                                    <p className="text-[10px] data-label text-wave-cyan">Digital RX: #MPSK-402</p>
                                </div>
                            </div>
                            <button className="w-full sanggam-btn py-4 flex items-center justify-center gap-2 group">
                                처방전 전송하기
                                <span className="inline-block transition-transform group-hover:translate-x-1">→</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
