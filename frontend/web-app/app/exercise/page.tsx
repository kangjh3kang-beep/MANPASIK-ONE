"use client";

import React from 'react';

export default function ExercisePage() {
    return (
        <div className="min-h-screen bg-black overflow-hidden font-sans relative">
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,_var(--wave-cyan-glow),_transparent_70%)] opacity-20"></div>

            <header className="absolute top-8 left-8 z-20">
                <h1 className="text-4xl font-serif shimmer-text">AI Exercise HUD</h1>
                <p className="data-label text-[10px] mt-1">Real-time Skeletal Tracking</p>
            </header>

            <div className="relative h-screen flex items-center justify-center">
                {/* Main AI Perspective */}
                <div className="relative w-[300px] h-[600px] flex items-center justify-center">
                    <div className="absolute inset-0 bg-cyan-500/10 blur-[100px] animate-pulse"></div>
                    <div className="text-[18rem] opacity-30 select-none breathing-animate">🧘‍♂️</div>

                    {/* Tracking Indicators */}
                    <div className="absolute top-[20%] left-1/2 -translate-x-1/2 w-4 h-4 bg-white rounded-full shadow-[0_0_20px_#fff]"></div>
                    <div className="absolute bottom-[30%] left-[20%] w-2 h-2 bg-[var(--wave-cyan)] rounded-full"></div>
                    <div className="absolute bottom-[30%] right-[20%] w-2 h-2 bg-[var(--wave-cyan)] rounded-full shadow-[0_0_10px_var(--wave-cyan)]"></div>
                </div>

                {/* Floating Metrics */}
                <div className="absolute left-10 top-1/2 -translate-y-1/2 space-y-6">
                    <div className="premium-card w-48">
                        <div className="data-label text-[9px]">Accuracy</div>
                        <div className="text-3xl font-serif text-[var(--wave-cyan)]">94%</div>
                    </div>
                    <div className="premium-card w-48">
                        <div className="data-label text-[9px]">Calories</div>
                        <div className="text-3xl font-serif">420</div>
                    </div>
                </div>

                <div className="absolute right-10 top-1/2 -translate-y-1/2 space-y-6">
                    <div className="premium-card w-64 border-[var(--sanggam-gold)]/30">
                        <h3 className="text-sm font-serif mb-4 flex items-center gap-2">
                            <span className="text-[var(--sanggam-gold)]">📈</span> 분석 요약
                        </h3>
                        <div className="space-y-3">
                            <div className="flex justify-between text-[10px] data-label">
                                <span>Alignment</span>
                                <span className="text-green-400">GOOD</span>
                            </div>
                            <div className="h-1 bg-white/5 rounded-full">
                                <div className="h-full bg-[var(--sanggam-gold)] w-4/5"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <footer className="absolute bottom-10 left-10 right-10 flex justify-between items-center z-20">
                <a href="/" className="sanggam-btn sanggam-btn-jagae text-xs px-8">EXIT TRAINING</a>
                <div className="sanggam-panel bg-black/40 px-8 py-4 backdrop-blur-md border-[var(--sanggam-gold)]/20">
                    <p className="text-xs font-serif italic opacity-60">"어깨 근육을 이완하고 호흡을 깊게 유지해 주세요."</p>
                </div>
                <button className="sanggam-btn sanggam-btn-gold text-xs px-8">NEXT FLOW</button>
            </footer>
        </div>
    );
}
