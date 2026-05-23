'use client';

import React from 'react';
import { Network, Database, ShieldCheck, ActivitySquare, Pill, FileCode2, Users2, Stethoscope, ChevronRight, Globe, Brain, Activity } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

/**
 * 현재 pathname에서 locale prefix를 추출
 */
function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const domains = [
  {
    id: 'measure',
    title: '건강 측정',
    desc: '카트리지를 인식하고 측정을 시작합니다. 오프라인에서도 사용 가능합니다',
    icon: Activity,
    colorClasses: 'bg-rose-50 text-rose-600 ring-rose-500/20 group-hover:bg-rose-500 group-hover:text-white',
    url: '/domains/measure'
  },
  {
    id: 'clinical',
    title: '건강 데이터 관리',
    desc: '환자의 생체 신호를 실시간으로 모니터링하고 AI가 분석 결과를 제공합니다',
    icon: Stethoscope,
    colorClasses: 'bg-sky-50 text-sky-600 ring-sky-500/20 group-hover:bg-sky-500 group-hover:text-white',
    url: '/domains/clinical'
  },
  {
    id: 'agents-hub',
    title: 'AI 분석 센터',
    desc: '인공지능 모델의 작동 상태를 확인하고 분석 성능을 관리합니다',
    icon: Network,
    colorClasses: 'bg-indigo-50 text-indigo-600 ring-indigo-500/20 group-hover:bg-indigo-500 group-hover:text-white',
    url: '/domains/agents-hub'
  },
  {
    id: 'predictor',
    title: '건강 위험도 예측',
    desc: '과거 검사 데이터를 분석하여 질병 발생 가능성을 미리 알려드립니다',
    icon: ActivitySquare,
    colorClasses: 'bg-emerald-50 text-emerald-600 ring-emerald-500/20 group-hover:bg-emerald-500 group-hover:text-white',
    url: '/domains/predictor'
  },
  {
    id: 'reward',
    title: '건강 데이터 보상',
    desc: '건강 검사 데이터를 제공하면 포인트로 보상받을 수 있습니다',
    icon: Database,
    colorClasses: 'bg-amber-50 text-amber-600 ring-amber-500/20 group-hover:bg-amber-500 group-hover:text-white',
    url: '/domains/reward'
  },
  {
    id: 'partner',
    title: '병원 연동 관리',
    desc: '협력 병원과 건강 데이터를 국제 표준 형식으로 안전하게 주고받습니다',
    icon: Users2,
    colorClasses: 'bg-rose-50 text-rose-600 ring-rose-500/20 group-hover:bg-rose-500 group-hover:text-white',
    url: '/domains/partner'
  },
  {
    id: 'gxp',
    title: '의약품 품질 관리',
    desc: '의약품 제조·유통 과정의 모든 기록을 추적하고 규정 준수를 확인합니다',
    icon: Pill,
    colorClasses: 'bg-violet-50 text-violet-600 ring-violet-500/20 group-hover:bg-violet-500 group-hover:text-white',
    url: '/domains/gxp'
  },
  {
    id: 'dev-portal',
    title: '개발자 도구',
    desc: '외부 개발자가 만파식 기능을 활용할 수 있도록 연동 도구를 제공합니다',
    icon: FileCode2,
    colorClasses: 'bg-slate-50 text-slate-600 ring-slate-500/20 group-hover:bg-slate-700 group-hover:text-white',
    url: '/domains/dev-portal'
  },
  {
    id: 'hardware-core',
    title: '측정 장비 관리',
    desc: '리더기 센서의 작동 상태를 확인하고 측정 정밀도를 관리합니다',
    icon: Globe,
    colorClasses: 'bg-orange-50 text-orange-600 ring-orange-500/20 group-hover:bg-orange-500 group-hover:text-white',
    url: '/domains/hardware-core'
  },
  {
    id: 'ai-coach',
    title: 'AI 건강 코치',
    desc: '생활비서와 AI주치의가 건강 궁금증에 맞춤 답변을 제공합니다',
    icon: Brain,
    colorClasses: 'bg-purple-50 text-purple-600 ring-purple-500/20 group-hover:bg-purple-500 group-hover:text-white',
    url: '/domains/ai-coach'
  },
  {
    id: 'app',
    title: '건강 관리 앱',
    desc: '스마트폰에서 건강 데이터를 확인하고 관리할 수 있는 모바일 앱입니다',
    icon: ShieldCheck,
    colorClasses: 'bg-cyan-50 text-cyan-600 ring-cyan-500/20 group-hover:bg-cyan-500 group-hover:text-white',
    url: '/domains/app'
  }
];

export function DomainBentoGrid() {
  const localePrefix = useLocalePrefix();

  return (
    <section id="domains" className="py-24 bg-slate-50">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-sm font-bold tracking-widest text-sky-600 uppercase">서비스 구성</h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl">
            하나로 연결된 건강 생태계, <span className="text-sky-600">9개 핵심 서비스</span>
          </p>
          <p className="mt-4 max-w-2xl mx-auto text-lg text-slate-600">
            각 서비스가 독립적으로 운영되면서도 하나의 플랫폼 위에서 유기적으로 연결됩니다.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {domains.map((domain) => (
            <Link 
              key={domain.id} 
              href={`${localePrefix}${domain.url}`}
              className="group relative flex flex-col justify-between overflow-hidden rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-200 transition-all duration-300 hover:-translate-y-2 hover:shadow-xl hover:ring-sky-500 cursor-pointer"
            >
              <div>
                <div className={`inline-flex rounded-2xl p-3 ring-1 ring-inset transition-colors duration-300 mb-5 ${domain.colorClasses}`}>
                  <domain.icon className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-bold text-slate-900 mb-2 transition-colors group-hover:text-sky-700">{domain.title}</h3>
                <p className="text-[13px] leading-relaxed text-slate-500">
                  {domain.desc}
                </p>
              </div>
              
              {/* Animated Button at bottom of card */}
              <div className="mt-8 pt-4 border-t border-slate-100 flex items-center justify-between transition-all duration-300">
                <span className="text-sm font-semibold text-slate-400 group-hover:text-sky-600 transition-colors">자세히 보기</span>
                <div className="h-8 w-8 rounded-full bg-slate-50 flex items-center justify-center group-hover:bg-sky-100 group-hover:text-sky-600 transition-colors">
                  <ChevronRight className="h-4 w-4" />
                </div>
              </div>

              {/* Decorative gradient blur showing on hover */}
              <div className="absolute -right-4 -top-4 -z-10 h-24 w-24 rounded-full bg-sky-200/50 opacity-0 blur-2xl transition-opacity duration-300 group-hover:opacity-100" />
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}
