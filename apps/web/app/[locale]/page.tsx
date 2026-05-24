'use client';
import React from 'react';
import { SiteHeader } from '@/components/site-header';
import { DomainBentoGrid } from '@/components/domain-bento-grid';
import { EcosystemNerveCenter } from '@/components/dashboard/ecosystem-nerve-center';
import { Activity, Shield, Cpu, Database, Zap, GitBranch, Lock, BarChart3, Heart, Beaker, Brain } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { KPICard, MeasurementResultCard } from '@mmup/ui';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

export default function UnifiedDashboard() {
  const localePrefix = useLocalePrefix();

  return (
    <div className="min-h-screen bg-slate-50">
      <SiteHeader />

      <main className="pt-24 pb-20">
        <div className="max-w-7xl mx-auto px-6 lg:px-8 space-y-12">
          {/* 히어로 섹션 */}
          <div className="relative rounded-3xl overflow-hidden bg-slate-900 text-white p-8 md:p-12 shadow-2xl">
            <div className="absolute inset-0 bg-gradient-to-br from-sky-900/50 via-slate-900 to-slate-900 opacity-90" />
            <div className="relative z-10 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <div>
                <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-sky-500/20 text-sky-400 text-xs font-bold mb-6 border border-sky-500/30">
                  <Activity className="h-3 w-3" />
                  시스템 정상 운영 중
                </span>
                <h1 className="text-3xl lg:text-5xl font-black tracking-tight mb-6 leading-tight">
                  만파식 건강 생태계<br />
                  <span className="text-sky-400 font-extrabold">통합 관리 센터</span>
                </h1>
                <p className="text-base lg:text-lg text-slate-300 mb-8 max-w-lg leading-relaxed">
                  9개 핵심 서비스와 인공지능 분석 허브를 한눈에 모니터링합니다.
                  나의 건강 데이터를 안전하게 관리하세요.
                </p>
                <div className="flex flex-wrap gap-4">
                  <Link
                    href={`${localePrefix}/domains/measure`}
                    className="px-6 py-3.5 bg-sky-600 hover:bg-sky-500 text-white font-bold rounded-2xl transition-all shadow-lg shadow-sky-600/30 flex items-center gap-2"
                  >
                    측정 시작하기
                  </Link>
                  <Link
                    href={`${localePrefix}/domains/health-records`}
                    className="px-6 py-3.5 bg-white/10 hover:bg-white/20 text-white font-bold rounded-2xl transition-all border border-white/20 backdrop-blur-sm"
                  >
                    건강 기록 보기
                  </Link>
                </div>
              </div>

              <div className="hidden lg:block">
                <EcosystemNerveCenter />
              </div>
            </div>
          </div>

          {/* 의료 인증 신뢰 배지 */}
          <section className="py-6 border-y border-slate-100 bg-slate-50/50 -mx-6 lg:-mx-8 px-6 lg:px-8 rounded-2xl">
            <div className="max-w-7xl mx-auto">
              <div className="flex items-center justify-center gap-8 flex-wrap">
                {[
                  { label: 'MFDS 인증', desc: '식품의약품안전처 의료기기 인증' },
                  { label: 'ISO 13485', desc: '의료기기 품질경영시스템' },
                  { label: 'IEC 62304', desc: '의료기기 소프트웨어 수명주기' },
                  { label: 'HIPAA', desc: '미국 건강정보 보호법 준수' },
                  { label: 'WCAG 2.2 AA', desc: '웹 접근성 지침 준수' },
                ].map(badge => (
                  <div key={badge.label} className="flex items-center gap-2 text-slate-500 hover:text-trust-blue transition-colors group" title={badge.desc}>
                    <Shield className="h-4 w-4 text-trust-blue/60 group-hover:text-trust-blue" />
                    <span className="text-xs font-semibold tracking-wide">{badge.label}</span>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {/* 플랫폼 현황 */}
          <section aria-label="플랫폼 현황">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">플랫폼 현황 한눈에 보기</h2>
            <p className="text-slate-500 mb-6">만파식 생태계의 핵심 수치를 실시간으로 확인합니다</p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <KPICard title="연결된 서비스" value="41" unit="개" trend="안정적으로 운영 중" icon={<Cpu className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
              <KPICard title="보안 등급" value="6" unit="단계" trend="군사급 암호화 적용" icon={<Shield className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
              <KPICard title="AI 분석 모델" value="5" unit="종" trend="편향 감지 기능 내장" icon={<Brain className="h-4 w-4" />} color="text-purple-600 bg-purple-50" />
              <KPICard title="측정 정밀도" value="1,792" unit="차원" trend="단계적 정밀도 향상" icon={<BarChart3 className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            </div>
          </section>

          {/* 측정 결과 예시 */}
          <section aria-label="측정 결과 예시">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">건강 측정 결과 예시</h2>
            <p className="text-slate-500 mb-6">색상 + 숫자 + 아이콘으로 결과를 쉽게 이해할 수 있습니다</p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <MeasurementResultCard
                biomarker="혈당"
                value={95}
                unit="mg/dL"
                confidence={0.96}
                uncertainty={{ value: 2.8, ci_lower: 92.2, ci_upper: 97.8 }}
                referenceRange={{ low: 70, high: 140 }}
                status="normal"
                timestamp="오후 2:32"
              />
              <MeasurementResultCard
                biomarker="당화혈색소"
                value={6.8}
                unit="%"
                confidence={0.93}
                uncertainty={{ value: 0.3, ci_lower: 6.5, ci_upper: 7.1 }}
                referenceRange={{ low: 4.0, high: 6.5 }}
                status="caution"
                timestamp="오후 2:32"
              />
              <MeasurementResultCard
                biomarker="총 콜레스테롤"
                value={245}
                unit="mg/dL"
                confidence={0.91}
                uncertainty={{ value: 8.5, ci_lower: 236.5, ci_upper: 253.5 }}
                referenceRange={{ low: 0, high: 200 }}
                status="danger"
                timestamp="오후 2:32"
              />
            </div>
          </section>

          {/* 안전 설계 원칙 */}
          <section aria-label="안전 설계 원칙">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">의료기기 안전 설계 원칙</h2>
            <p className="text-slate-500 mb-6">6가지 핵심 원칙으로 안전한 의료 소프트웨어를 만듭니다</p>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[
                { id: 'H1', title: '부품 독립 설계', desc: '각 부품이 독립적으로 작동하여 한 곳의 문제가 전체에 영향을 주지 않습니다', score: 95, icon: GitBranch },
                { id: 'H2', title: '이전 버전 호환', desc: '새 버전이 나와도 기존 장비와 카트리지가 그대로 작동합니다', score: 95, icon: Zap },
                { id: 'H3', title: '인증 간소화 설계', desc: '의료기기 인증을 모듈별로 받을 수 있도록 설계했습니다', score: 90, icon: Shield },
                { id: 'H4', title: '단계적 기능 확장', desc: '혈액검사에서 시작해 광학·유전자 검사까지 확장 가능합니다', score: 95, icon: Beaker },
                { id: 'H5', title: '통신 규약 명확화', desc: '기기 간 통신 방식을 버전별로 명확히 정의합니다', score: 90, icon: Lock },
                { id: 'H6', title: '오류 자동 격리', desc: '센서 하나가 고장나도 나머지는 정상 작동합니다', score: 95, icon: Heart },
              ].map(h => (
                <div key={h.id} className="rounded-2xl border border-slate-200 bg-white p-5 hover:shadow-md transition-shadow">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-sky-50 text-sky-600">
                      <h.icon className="h-5 w-5" />
                    </div>
                    <div>
                      <span className="text-xs font-bold text-sky-600">{h.id}</span>
                      <h3 className="text-sm font-bold text-slate-900">{h.title}</h3>
                    </div>
                    <div className="ml-auto">
                      <span className="text-lg font-black text-emerald-600">{h.score}%</span>
                    </div>
                  </div>
                  <p className="text-xs text-slate-500 leading-relaxed">{h.desc}</p>
                  <div className="mt-3 h-1.5 bg-slate-100 rounded-full overflow-hidden">
                    <div className="h-full bg-emerald-500 rounded-full transition-all" style={{ width: `${h.score}%` }} />
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* 기술 구조 */}
          <section aria-label="기술 구조">
            <div className="rounded-2xl bg-slate-900 text-white p-8">
              <h2 className="text-xl font-bold mb-6">6단계 통합 기술 구조</h2>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                {[
                  { layer: '6층', name: '클라우드', tech: '서버 관리', color: 'bg-slate-700' },
                  { layer: '5층', name: '백엔드', tech: '41개 서비스', color: 'bg-sky-800' },
                  { layer: '4층', name: '모바일 앱', tech: '건강 관리 앱', color: 'bg-cyan-800' },
                  { layer: '3층', name: '분석 엔진', tech: '측정 데이터 처리', color: 'bg-emerald-800' },
                  { layer: '2층', name: '장비 제어', tech: '센서 연결', color: 'bg-amber-800' },
                  { layer: '1층', name: '측정 장비', tech: '리더기 하드웨어', color: 'bg-red-800' },
                ].map(l => (
                  <div key={l.layer} className={`${l.color} rounded-xl p-4 text-center`}>
                    <span className="text-xs text-white/60">{l.layer}</span>
                    <p className="text-sm font-bold mt-1">{l.name}</p>
                    <p className="text-[10px] text-white/70 mt-1">{l.tech}</p>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {/* 서비스 목록 */}
          <div>
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-2xl font-bold text-slate-900">핵심 서비스 목록</h2>
                <p className="text-slate-500">각 서비스를 클릭하면 상세 정보를 확인할 수 있습니다</p>
              </div>
              <div className="flex gap-2">
                <div className="flex items-center gap-2 px-3 py-1 bg-white border border-slate-200 rounded-lg text-xs font-semibold text-slate-600">
                  <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  정상 운영: 9개
                </div>
              </div>
            </div>

            <DomainBentoGrid />
          </div>

          {/* 플랫폼 요약 정보 */}
          <section aria-label="플랫폼 요약" className="rounded-2xl border border-slate-200 bg-white p-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Database className="h-4 w-4 text-sky-600" />
                  데이터 표준 계약
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• 측정 데이터 표준 형식 정의</li>
                  <li>• 의료 코드 연동 (15종 검사항목)</li>
                  <li>• 실시간 이벤트 연동 (8종)</li>
                  <li>• 49개 공개 API 제공</li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Shield className="h-4 w-4 text-emerald-600" />
                  의료기기 인증 준비 (19건)
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• 의료 소프트웨어 안전 등급 (IEC 62304)</li>
                  <li>• 위험 분석 38건 완료 (ISO 14971)</li>
                  <li>• 사용 편의성 시험 계획 수립</li>
                  <li>• AI 모델 공정성 평가 (5종)</li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Beaker className="h-4 w-4 text-purple-600" />
                  핵심 기준값
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• 카트리지 연결: 16핀 표준 커넥터</li>
                  <li>• 측정 보정: 자동 차동 보정 (α=0.98)</li>
                  <li>• 측정 정확도: 92~98%</li>
                  <li>• 검사 항목: 44종 카트리지 지원</li>
                </ul>
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}
