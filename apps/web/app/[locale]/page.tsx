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
          {/* Hero Section / Nerve Center */}
          <div className="relative rounded-3xl overflow-hidden bg-slate-900 text-white p-12 shadow-2xl">
            <div className="absolute inset-0 bg-gradient-to-br from-sky-900/50 via-slate-900 to-slate-900 opacity-90" />
            <div className="relative z-10 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <div>
                <span className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-sky-500/20 text-sky-400 text-xs font-bold mb-6 border border-sky-500/30">
                  <Activity className="h-3 w-3" />
                  SYSTEM OVERWATCH ACTIVE
                </span>
                <h1 className="text-4xl lg:text-5xl font-black tracking-tight mb-6 leading-tight">
                  만파식 생태계<br />
                  <span className="text-sky-400 font-extrabold">통합 운영 대시보드</span>
                </h1>
                <p className="text-lg text-slate-300 mb-8 max-w-lg leading-relaxed">
                  9대 핵심 도메인과 AI 에이전트 허브를 실시간으로 모니터링하고 제어합니다.
                  정밀 의학 데이터의 전주기 관리를 한눈에 확인하세요.
                </p>
                <div className="flex flex-wrap gap-4">
                  <Link
                    href={`${localePrefix}/domains/clinical`}
                    className="px-6 py-3.5 bg-sky-600 hover:bg-sky-500 text-white font-bold rounded-2xl transition-all shadow-lg shadow-sky-600/30 flex items-center gap-2"
                  >
                    데모 콘솔 열기
                  </Link>
                  <Link
                    href={`${localePrefix}/domains/app`}
                    className="px-6 py-3.5 bg-white/10 hover:bg-white/20 text-white font-bold rounded-2xl transition-all border border-white/20 backdrop-blur-sm"
                  >
                    건강 앱 관리
                  </Link>
                </div>
              </div>

              <div className="hidden lg:block">
                <EcosystemNerveCenter />
              </div>
            </div>
          </div>

          {/* ===== NEW: Platform Status KPI Cards ===== */}
          <section aria-label="플랫폼 현황">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">플랫폼 실시간 현황</h2>
            <p className="text-slate-500 mb-6">하네스 엔지니어링 H1~H6 기반 6-Layer 통합 생태계</p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <KPICard title="마이크로서비스" value="41" unit="서비스" trend="gRPC + REST Gateway" icon={<Cpu className="h-4 w-4" />} color="text-sky-600 bg-sky-50" />
              <KPICard title="보안 계층" value="6" unit="Layer" trend="AES-256 + 해시체인" icon={<Shield className="h-4 w-4" />} color="text-emerald-600 bg-emerald-50" />
              <KPICard title="AI 모델" value="5" unit="종" trend="BiasDetector 내장" icon={<Brain className="h-4 w-4" />} color="text-purple-600 bg-purple-50" />
              <KPICard title="핑거프린트" value="1,792" unit="차원" trend="88→448→896→1792" icon={<BarChart3 className="h-4 w-4" />} color="text-amber-600 bg-amber-50" />
            </div>
          </section>

          {/* ===== NEW: Measurement Demo with 3-Triple Encoding ===== */}
          <section aria-label="측정 결과 데모">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">차동측정 결과 데모</h2>
            <p className="text-slate-500 mb-6">Sdiff = S - α×R (α=0.98) — 3중 인코딩: 색채 + 수치 + 아이콘</p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <MeasurementResultCard
                biomarker="혈당 (Glucose)"
                value={95}
                unit="mg/dL"
                confidence={0.96}
                uncertainty={{ value: 2.8, ci_lower: 92.2, ci_upper: 97.8 }}
                referenceRange={{ low: 70, high: 140 }}
                status="normal"
                timestamp="14:32"
              />
              <MeasurementResultCard
                biomarker="HbA1c"
                value={6.8}
                unit="%"
                confidence={0.93}
                uncertainty={{ value: 0.3, ci_lower: 6.5, ci_upper: 7.1 }}
                referenceRange={{ low: 4.0, high: 6.5 }}
                status="caution"
                timestamp="14:32"
              />
              <MeasurementResultCard
                biomarker="총 콜레스테롤"
                value={245}
                unit="mg/dL"
                confidence={0.91}
                uncertainty={{ value: 8.5, ci_lower: 236.5, ci_upper: 253.5 }}
                referenceRange={{ low: 0, high: 200 }}
                status="danger"
                timestamp="14:32"
              />
            </div>
          </section>

          {/* ===== NEW: Harness Engineering H1-H6 Status ===== */}
          <section aria-label="하네스 엔지니어링">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">하네스 엔지니어링 준수 현황</h2>
            <p className="text-slate-500 mb-6">의료기기 SW 안전성 6대 원칙</p>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[
                { id: 'H1', title: '모듈 독립성', desc: 'AFE 블록 HAL 추상화, 구체 타입 직접 참조 금지', score: 95, icon: GitBranch },
                { id: 'H2', title: '후방 호환', desc: 'CSI v1.0 ↔ v2.0 어댑터, BLE GATT 버전 협상', score: 95, icon: Zap },
                { id: 'H3', title: '사전인증 플랫폼', desc: '510(k) 모듈 단위 인증, 해시체인 무결성', score: 90, icon: Shield },
                { id: 'H4', title: '점진적 확장', desc: 'Stage-1(전기화학) → Stage-2(광학) → Stage-3(NAAT)', score: 95, icon: Beaker },
                { id: 'H5', title: '인터페이스 계약', desc: 'BLE GATT UUID + NFC 매니페스트 버전화', score: 90, icon: Lock },
                { id: 'H6', title: '실패 격리', desc: '채널별 독립 처리, Result<T,E>, 회로차단기', score: 95, icon: Heart },
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

          {/* ===== NEW: Tech Stack Overview ===== */}
          <section aria-label="기술 스택">
            <div className="rounded-2xl bg-slate-900 text-white p-8">
              <h2 className="text-xl font-bold mb-6">6-Layer 통합 아키텍처</h2>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                {[
                  { layer: 'L6', name: '인프라', tech: 'K8s · Docker', color: 'bg-slate-700' },
                  { layer: 'L5', name: '백엔드', tech: 'Go · gRPC · 41 MSA', color: 'bg-sky-800' },
                  { layer: 'L4', name: '앱', tech: 'Flutter · Riverpod', color: 'bg-cyan-800' },
                  { layer: 'L3', name: '코어', tech: 'Rust · no_std', color: 'bg-emerald-800' },
                  { layer: 'L2', name: 'HW 제어', tech: 'embedded-hal', color: 'bg-amber-800' },
                  { layer: 'L1', name: '하드웨어', tech: 'STM32 · nRF52', color: 'bg-red-800' },
                ].map(l => (
                  <div key={l.layer} className={`${l.color} rounded-xl p-4 text-center`}>
                    <span className="text-xs font-mono text-white/60">{l.layer}</span>
                    <p className="text-sm font-bold mt-1">{l.name}</p>
                    <p className="text-[10px] text-white/70 mt-1">{l.tech}</p>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {/* Infrastructure Nodes */}
          <div>
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-2xl font-bold text-slate-900">도메인 인프라 노드</h2>
                <p className="text-slate-500">생태계의 각 서비스를 독립적으로 관리하고 상태를 확인합니다.</p>
              </div>
              <div className="flex gap-2">
                <div className="flex items-center gap-2 px-3 py-1 bg-white border border-slate-200 rounded-lg text-xs font-semibold text-slate-600">
                  <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  Live Nodes: 9
                </div>
              </div>
            </div>

            <DomainBentoGrid />
          </div>

          {/* ===== NEW: Contracts & Compliance Footer ===== */}
          <section aria-label="계약 및 규제" className="rounded-2xl border border-slate-200 bg-white p-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Database className="h-4 w-4 text-sky-600" />
                  계약 SSOT (contracts/)
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• 표준 측정 패킷 스키마 (JSON Schema)</li>
                  <li>• LOINC/UCUM 매핑 레지스트리 (15종)</li>
                  <li>• Kafka 이벤트 스키마 (8 토픽)</li>
                  <li>• OpenAPI 3.1 (49 엔드포인트)</li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Shield className="h-4 w-4 text-emerald-600" />
                  규제 준수 (19 문서)
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• IEC 62304 Class B (SRS/SDP/SAD)</li>
                  <li>• ISO 14971 FMEA (38 실패 모드)</li>
                  <li>• IEC 62366 HFE 계획</li>
                  <li>• FDA AI/ML 모델카드 (5종)</li>
                </ul>
              </div>
              <div>
                <h3 className="text-sm font-bold text-slate-900 mb-3 flex items-center gap-2">
                  <Beaker className="h-4 w-4 text-purple-600" />
                  SSOT 기준값
                </h3>
                <ul className="space-y-2 text-xs text-slate-500">
                  <li>• 커넥터: CSI v1.0 (16핀, 1.27mm)</li>
                  <li>• 차분식: Sdiff = S - α×R (α=0.98)</li>
                  <li>• 정확도: 92~98% (결합)</li>
                  <li>• 카트리지: 44 SKU (5 Layer)</li>
                </ul>
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}

