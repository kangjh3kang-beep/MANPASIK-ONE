"use client";

import React from 'react';

export default function PremiumExercise() {
    return (
        <div className="min-h-screen bg-black overflow-hidden font-sans">
            {/* Dynamic Background */}
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,_rgba(0,242,255,0.05),_transparent_70%)]"></div>

            <div className="relative h-screen grid grid-cols-4 gap-6 p-8">
                {/* Left Stats Bar */}
                <div className="col-span-1 space-y-4 flex flex-col justify-center">
                    <div className="premium-card border-wave-cyan/30">
                        <div className="data-label mb-2">Skeletal Accuracy</div>
                        <div className="text-4xl font-serif text-wave-cyan">94%</div>
                        <div className="mt-4 flex gap-1">
                            {[1, 1, 1, 1, 0, 0].map((v, i) => (
                                <div key={i} className={`h-1 flex-1 rounded-full ${v ? 'bg-wave-cyan' : 'bg-white/10'}`}></div>
                            ))}
                        </div>
                    </div>
                    <div className="premium-card">
                        <div className="data-label mb-2">Calories Burned</div>
                        <div className="text-4xl font-serif text-white">420 <span className="text-xs data-label">kcal</span></div>
                    </div>
                    <div className="premium-card">
                        <div className="data-label mb-2">Workout Time</div>
                        <div className="text-4xl font-serif text-white">24:12</div>
                    </div>
                </div>

                {/* Center Main AI Guide */}
                <div className="col-span-2 relative flex items-center justify-center">
                    {/* Holograph Trainer Mockup */}
                    <div className="w-full h-full max-h-[80vh] flex items-center justify-center">
                        <div className="relative w-80 h-full">
                            <div className="absolute inset-0 bg-gradient-to-t from-wave-cyan/20 to-transparent blur-2xl animate-pulse"></div>
                            <div className="w-full h-full flex flex-col items-center justify-center text-[20rem] opacity-40 select-none">
                                🧘‍♂️
                            </div>
                            {/* Joint Circles Overlay */}
                            <div className="absolute top-[20%] left-1/2 -translate-x-1/2 w-4 h-4 bg-white rounded-full shadow-[0_0_15px_#fff] animate-ping"></div>
                            <div className="absolute top-[40%] left-[30%] w-3 h-3 bg-wave-cyan rounded-full"></div>
                            <div className="absolute top-[40%] right-[30%] w-3 h-3 bg-wave-cyan rounded-full shadow-[0_0_10px_var(--wave-cyan)]"></div>
                        </div>
                    </div>

                    {/* Bottom Control HUD */}
                    <div className="absolute bottom-4 left-0 right-0">
                        <div className="sanggam-panel p-6 flex items-center justify-between border-sanggam/30 bg-black/80">
                            <div className="flex items-center gap-6">
                                <div className="h-12 w-1 bg-sanggam animate-pulse"></div>
                                <div>
                                    <h4 className="text-xs data-label">Current Pose</h4>
                                    <h2 className="text-xl font-serif shimmer-text">Downward Facing Dog</h2>
                                </div>
                            </div>
                            <div className="flex gap-4">
                                <button className="px-6 py-2 bg-white/5 border border-white/10 rounded font-serif text-sm hover:bg-white/10 transition-all">PAUSE</button>
                                <button className="px-8 py-2 bg-sanggam border border-sanggam rounded font-serif text-sm text-black font-bold hover:shadow-[0_0_20px_var(--sanggam-gold-glow)] transition-all">NEXT FLOW</button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Right Activity Bar */}
                <div className="col-span-1 flex flex-col justify-center space-y-6">
                    <h2 className="text-xl font-serif border-b border-sanggam/20 pb-4 flex items-center gap-2">
                        <span className="text-sanggam">📈</span> 실시간 분석
                    </h2>
                    <div className="space-y-4">
                        {[1, 2, 3].map((_, i) => (
                            <div key={i} className="premium-card p-4 group">
                                <div className="flex justify-between mb-3">
                                    <span className="text-[10px] data-label opacity-60">Shoulder Alignment</span>
                                    <span className="text-[10px] data-label text-sanggam">OPTIMAL</span>
                                </div>
                                <div className="h-12 flex items-end gap-1">
                                    {[4, 7, 5, 8, 6, 9, 7].map((h, j) => (
                                        <div key={j} className="flex-1 bg-sanggam/20 rounded-t group-hover:bg-sanggam/40 transition-all" style={{ height: `${h * 10}%` }}></div>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>

                    <div className="sanggam-panel p-4 border-wave-cyan/20">
                        <div className="flex items-center gap-4">
                            <div className="w-10 h-10 rounded-full bg-wave-cyan/10 flex items-center justify-center text-wave-cyan">💬</div>
                            <div className="text-[10px] font-serif leading-tight">
                                "어깨에 힘을 조금 더 빼고 <br />호흡을 길게 유지하세요."
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
