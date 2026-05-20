/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
'use client';

import React, { Suspense } from 'react';
import dynamic from 'next/dynamic';
import { ArrowRight, PlayCircle, ShieldCheck, BadgeCheck, Stethoscope } from 'lucide-react';

const MedicalHero3DWrapper = dynamic(
  () => import('./medical-hero-3d').then((mod) => mod.MedicalHero3D),
  { ssr: false }
);

interface MedicalHeroProps {
  onPrimary: () => void;
  onSecondary: () => void;
  locale?: string;
}



export function MedicalHero({ onPrimary, onSecondary, locale = 'ko' }: MedicalHeroProps) {
  const isKo = locale === 'ko';

  return (
    <section className="relative overflow-hidden bg-gradient-to-b from-slate-50 via-white to-white">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-[0.35]"
        style={{
          backgroundImage:
            'radial-gradient(circle at 20% 0%, rgba(14,165,233,0.10), transparent 45%), radial-gradient(circle at 80% 30%, rgba(59,130,246,0.08), transparent 50%)',
        }}
      />
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          backgroundImage:
            'linear-gradient(to right, rgba(15,23,42,0.04) 1px, transparent 1px), linear-gradient(to bottom, rgba(15,23,42,0.04) 1px, transparent 1px)',
          backgroundSize: '56px 56px',
          maskImage: 'radial-gradient(ellipse at center, black 30%, transparent 75%)',
          WebkitMaskImage: 'radial-gradient(ellipse at center, black 30%, transparent 75%)',
        }}
      />

      <div className="relative mx-auto grid max-w-7xl grid-cols-1 gap-12 px-4 pb-20 pt-14 sm:px-6 lg:grid-cols-2 lg:gap-8 lg:px-8 lg:pt-20">
        <div className="flex flex-col justify-center">
          <span className="inline-flex w-fit items-center gap-2 rounded-full border border-sky-200 bg-sky-50 px-3 py-1 text-xs font-semibold tracking-wide text-sky-700">
            <BadgeCheck className="h-3.5 w-3.5" />
            {isKo ? 'KFDA · CE · ISO 13485 인증' : 'KFDA · CE · ISO 13485 Certified'}
          </span>

          <h1 className="mt-5 text-[2.25rem] font-semibold leading-[1.12] tracking-tight text-slate-900 sm:text-5xl lg:text-[3.4rem]">
            {isKo ? (
              <>
                정밀 의학을 일상으로,
                <br />
                <span className="bg-gradient-to-r from-sky-700 to-blue-700 bg-clip-text text-transparent">
                  당신의 손 안에서.
                </span>
              </>
            ) : (
              <>
                Precision medicine,
                <br />
                <span className="bg-gradient-to-r from-sky-700 to-blue-700 bg-clip-text text-transparent">
                  in the palm of your hand.
                </span>
              </>
            )}
          </h1>

          <p className="mt-5 max-w-xl text-[15px] leading-[1.7] text-slate-600 sm:text-base">
            {isKo
              ? '연구실 수준의 분광 진단을 가정에서 수행합니다. 50가지 이상의 생체 지표를 비침습적으로 측정하고, 임상 가이드라인에 기반한 인사이트를 제공합니다.'
              : 'Lab-grade spectral diagnostics, at home. Track 50+ biomarkers non-invasively and receive clinically-grounded insights tailored to you.'}
          </p>

          <div className="mt-8 flex flex-wrap gap-3">
            <button
              onClick={onPrimary}
              className="inline-flex items-center gap-2 rounded-md bg-gradient-to-b from-sky-600 to-blue-700 px-5 py-3 text-sm font-semibold text-white shadow-[0_1px_0_rgba(255,255,255,0.2)_inset,0_8px_20px_-8px_rgba(2,132,199,0.6)] transition-transform hover:-translate-y-px hover:shadow-[0_1px_0_rgba(255,255,255,0.2)_inset,0_12px_24px_-8px_rgba(2,132,199,0.7)]"
            >
              {isKo ? '솔루션 살펴보기' : 'Explore Solutions'}
              <ArrowRight className="h-4 w-4" />
            </button>
            <button
              onClick={onSecondary}
              className="inline-flex items-center gap-2 rounded-md border border-slate-300 bg-white px-5 py-3 text-sm font-semibold text-slate-800 transition-colors hover:border-slate-400 hover:bg-slate-50"
            >
              <PlayCircle className="h-4 w-4 text-sky-700" />
              {isKo ? '데모 보기' : 'Watch Demo'}
            </button>
          </div>

          <dl className="mt-10 grid grid-cols-3 gap-6 border-t border-slate-200 pt-6">
            {[
              { v: '99.8%', l: isKo ? '진단 정확도' : 'Diagnostic accuracy' },
              { v: '< 90s', l: isKo ? '결과 도출 시간' : 'Time to result' },
              { v: '50+', l: isKo ? '생체 지표' : 'Biomarkers' },
            ].map((s) => (
              <div key={s.l}>
                <dt className="text-2xl font-semibold tracking-tight text-slate-900 tabular-nums">{s.v}</dt>
                <dd className="mt-1 text-xs font-medium text-slate-500">{s.l}</dd>
              </div>
            ))}
          </dl>
        </div>

        <div className="relative">
          <div className="relative mx-auto aspect-square w-full max-w-[520px]">
            <div className="absolute inset-0 rounded-[28px] bg-gradient-to-br from-sky-50 via-white to-blue-50 ring-1 ring-slate-200 shadow-[0_30px_60px_-30px_rgba(15,23,42,0.18)]" />
            <div className="absolute inset-3 overflow-hidden rounded-[22px]">
              <MedicalHero3DWrapper />
            </div>

            <div className="absolute -left-3 top-10 hidden rounded-xl border border-slate-200 bg-white/90 px-4 py-3 shadow-lg backdrop-blur sm:block">
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-sky-50 text-sky-700">
                  <Stethoscope className="h-4 w-4" />
                </div>
                <div>
                  <p className="text-[11px] font-medium uppercase tracking-wider text-slate-500">
                    {isKo ? '실시간 분석' : 'Live analysis'}
                  </p>
                  <p className="text-sm font-semibold text-slate-900">
                    {isKo ? 'Hb 14.2 g/dL' : 'Hb 14.2 g/dL'}
                  </p>
                </div>
              </div>
            </div>
            <div className="absolute -right-3 bottom-12 hidden rounded-xl border border-slate-200 bg-white/90 px-4 py-3 shadow-lg backdrop-blur sm:block">
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700">
                  <ShieldCheck className="h-4 w-4" />
                </div>
                <div>
                  <p className="text-[11px] font-medium uppercase tracking-wider text-slate-500">
                    {isKo ? '임상 검증' : 'Clinically validated'}
                  </p>
                  <p className="text-sm font-semibold text-slate-900">
                    {isKo ? '12개 종합병원' : '12 hospitals'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
