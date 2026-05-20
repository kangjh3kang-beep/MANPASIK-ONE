"use client";

import React from 'react';

export default function MarketPage() {
    const categories = ['자양강장', '디지털 건강', '한방 보조제', '스마트 기기'];
    const products = [
        { name: 'Sanggam Bio-Ring', price: '₩289,000', tag: 'BEST', img: '💍' },
        { name: 'Jagae Smart Patch', price: '₩45,000', tag: 'NEW', img: '🩹' },
        { name: 'Ginseng Nano-Extract', price: '₩120,000', tag: 'HOT', img: '🧪' },
        { name: 'Holo-Body Scanner', price: '₩560,000', tag: 'PREMIUM', img: '📡' },
    ];

    return (
        <div className="min-h-screen bg-[var(--deep-sea)] p-8 font-sans">
            <header className="flex justify-between items-center mb-12 border-b border-[var(--sanggam-gold)]/20 pb-6">
                <div>
                    <h1 className="text-4xl font-serif shimmer-text">Manpasik Market</h1>
                    <p className="data-label text-xs">Premium Health Curations</p>
                </div>
                <div className="flex gap-4">
                    <a href="/" className="sanggam-btn sanggam-btn-jagae text-xs px-6">DASHBOARD</a>
                    <button className="w-10 h-10 rounded-full border border-[var(--sanggam-gold)]/30 flex items-center justify-center text-[var(--sanggam-gold)] text-xl bg-[var(--sanggam-gold)]/5 hover:bg-[var(--sanggam-gold)]/20 transition-all">🛒</button>
                </div>
            </header>

            {/* Category List */}
            <div className="flex gap-4 mb-12 overflow-x-auto pb-4 scrollbar-hide">
                {categories.map((cat, i) => (
                    <button key={i} className="px-6 py-2 rounded-full border border-[var(--sanggam-gold)]/20 text-[var(--sanggam-gold)]/60 hover:border-[var(--sanggam-gold)] hover:text-[var(--sanggam-gold)] transition-all whitespace-nowrap text-sm font-serif">
                        {cat}
                    </button>
                ))}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                {products.map((p, i) => (
                    <div key={i} className="premium-card group cursor-pointer">
                        <div className="absolute top-4 right-4 bg-[var(--sanggam-gold)]/10 text-[var(--sanggam-gold)] text-[10px] px-2 py-1 rounded border border-[var(--sanggam-gold)]/30 font-mono">
                            {p.tag}
                        </div>
                        <div className="h-48 flex items-center justify-center text-7xl mb-6 bg-white/5 rounded-xl border border-white/5 group-hover:bg-[var(--sanggam-gold)]/5 transition-colors">
                            {p.img}
                        </div>
                        <h3 className="text-xl font-serif mb-2 group-hover:text-white transition-colors">{p.name}</h3>
                        <div className="flex justify-between items-end">
                            <span className="text-xl text-[var(--sanggam-gold)] font-mono">{p.price}</span>
                            <button className="w-8 h-8 rounded-full bg-[var(--sanggam-gold)]/10 flex items-center justify-center text-[var(--sanggam-gold)] border border-[var(--sanggam-gold)]/20 hover:bg-[var(--sanggam-gold)] hover:text-[var(--ink-black)] transition-all">
                                +
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            {/* Recommended Section */}
            <section className="mt-20">
                <h2 className="text-2xl font-serif mb-8 flex items-center gap-4">
                    <span className="w-8 h-px bg-[var(--sanggam-gold)]"></span>
                    맞춤형 추천
                </h2>
                <div className="sanggam-panel p-8 flex flex-col md:flex-row items-center gap-8">
                    <div className="flex-1">
                        <h3 className="text-3xl font-serif mb-4 shimmer-text">AI 기반 건강 분석팩</h3>
                        <p className="text-[var(--hanji-white)] opacity-60 text-sm leading-relaxed mb-6">
                            사용자님의 최신 바이오 데이터를 분석한 결과입니다. 정밀한 영양 관리가 필요합니다.
                        </p>
                        <button className="sanggam-btn sanggam-btn-gold px-12">분석 결과 보기</button>
                    </div>
                    <div className="w-64 h-64 bg-[var(--sanggam-gold)]/5 rounded-full border border-[var(--sanggam-gold)]/10 flex items-center justify-center text-9xl breathing-animate">
                        📦
                    </div>
                </div>
            </section>
        </div>
    );
}
