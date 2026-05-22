'use client';

import React from 'react';
import { Network, Database, ShieldCheck, ActivitySquare, Pill, FileCode2, Users2, Stethoscope, ChevronRight, Globe } from 'lucide-react';
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
    id: 'clinical',
    title: '임상 데이터 콘솔',
    desc: '실시간 패킷 검증 및 해시 체인 무결성 시각화 (IEC 62304 B등급)',
    icon: Stethoscope,
    colorClasses: 'bg-sky-50 text-sky-600 ring-sky-500/20 group-hover:bg-sky-500 group-hover:text-white',
    url: '/domains/clinical'
  },
  {
    id: 'agents-hub',
    title: 'AI 에이전트 허브',
    desc: '의료 AI 모델 오케스트레이션 및 상태 관리',
    icon: Network,
    colorClasses: 'bg-indigo-50 text-indigo-600 ring-indigo-500/20 group-hover:bg-indigo-500 group-hover:text-white',
    url: '/domains/agents-hub'
  },
  {
    id: 'predictor',
    title: '생체 지표 예측',
    desc: 'GNN 기반 다중 질환 시계열 예측 시스템',
    icon: ActivitySquare,
    colorClasses: 'bg-emerald-50 text-emerald-600 ring-emerald-500/20 group-hover:bg-emerald-500 group-hover:text-white',
    url: '/domains/predictor'
  },
  {
    id: 'reward',
    title: '환자 리워드 풀',
    desc: '개인 데이터 주권 및 토큰 보상 스마트 컨트랙트',
    icon: Database,
    colorClasses: 'bg-amber-50 text-amber-600 ring-amber-500/20 group-hover:bg-amber-500 group-hover:text-white',
    url: '/domains/reward'
  },
  {
    id: 'partner',
    title: '파트너 통합 연동',
    desc: '제네바 HL7/FHIR 국제 표준 의료 데이터 파이프라인 연계',
    icon: Users2,
    colorClasses: 'bg-rose-50 text-rose-600 ring-rose-500/20 group-hover:bg-rose-500 group-hover:text-white',
    url: '/domains/partner'
  },
  {
    id: 'gxp',
    title: '의약품 GxP 준수',
    desc: 'cGMP 추적성 확보 및 전자 서명 기반 워크플로우',
    icon: Pill,
    colorClasses: 'bg-violet-50 text-violet-600 ring-violet-500/20 group-hover:bg-violet-500 group-hover:text-white',
    url: '/domains/gxp'
  },
  {
    id: 'dev-portal',
    title: '개발자 포털',
    desc: '오픈 API 명세서, 웹훅 연동 및 실시간 샌드박스',
    icon: FileCode2,
    colorClasses: 'bg-slate-50 text-slate-600 ring-slate-500/20 group-hover:bg-slate-700 group-hover:text-white',
    url: '/domains/dev-portal'
  },
  {
    id: 'hardware-core',
    title: '하드웨어 코어',
    desc: '분광 분석 엔진 및 센서 캘리브레이션 시스템 (Optics Control)',
    icon: Globe,
    colorClasses: 'bg-orange-50 text-orange-600 ring-orange-500/20 group-hover:bg-orange-500 group-hover:text-white',
    url: '/domains/hardware-core'
  },
  {
    id: 'app',
    title: '모바일 하이브리드',
    desc: 'Flutter Native 통신 연동용 PWA Fallback 환경',
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
          <h2 className="text-sm font-bold tracking-widest text-sky-600 uppercase">Architecture</h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl">
            초연결 의료 생태계, <span className="text-sky-600">9개의 도메인 허브</span>
          </p>
          <p className="mt-4 max-w-2xl mx-auto text-lg text-slate-600">
            마이크로서비스 아키텍처(MSA)를 통해 분리된 9개의 전문 도메인 앱은 하나의 Turborepo 모노레포 위에서 완벽히 조화롭게 유기적으로 결합합니다.
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
                <span className="text-sm font-semibold text-slate-400 group-hover:text-sky-600 transition-colors">허브 연결 및 이동</span>
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
