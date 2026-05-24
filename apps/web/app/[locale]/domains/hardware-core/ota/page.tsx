'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { RefreshCw, CheckCircle, AlertTriangle, Clock, Shield, Loader2, Server } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type UpdateStatus = '완료' | '진행중' | '실패' | '대기';

interface DeviceUpdate {
  id: string;
  model: string;
  currentVersion: string;
  targetVersion: string;
  status: UpdateStatus;
  lastUpdate: string;
}

const STATUS_STYLES: Record<UpdateStatus, string> = {
  '완료': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '진행중': 'bg-sky-50 text-sky-700 border-sky-200',
  '실패': 'bg-red-50 text-red-700 border-red-200',
  '대기': 'bg-slate-50 text-slate-600 border-slate-200',
};

const STATUS_ICONS: Record<UpdateStatus, React.ElementType> = {
  '완료': CheckCircle,
  '진행중': Loader2,
  '실패': AlertTriangle,
  '대기': Clock,
};

const CHANGELOG = [
  '차동측정 보정 알고리즘 개선 — Sdiff 계수 자동 튜닝 정확도 향상',
  'BLE 5.3 연결 안정성 강화 및 재연결 로직 최적화',
  '보안 패치: TLS 1.3 인증서 갱신 및 펌웨어 서명 검증 강화',
];

const MOCK_DEVICES: DeviceUpdate[] = [
  { id: 'MPS-R-001', model: 'MPS Reader Pro', currentVersion: 'v2.5.0', targetVersion: 'v2.5.0', status: '완료', lastUpdate: '2026-05-24 08:30' },
  { id: 'MPS-R-002', model: 'MPS Reader Pro', currentVersion: 'v2.5.0', targetVersion: 'v2.5.0', status: '완료', lastUpdate: '2026-05-24 08:45' },
  { id: 'MPS-R-003', model: 'MPS Reader Lite', currentVersion: 'v2.4.1', targetVersion: 'v2.5.0', status: '진행중', lastUpdate: '2026-05-24 09:10' },
  { id: 'MPS-R-004', model: 'MPS Reader Pro', currentVersion: 'v2.4.1', targetVersion: 'v2.5.0', status: '대기', lastUpdate: '2026-05-23 16:00' },
  { id: 'MPS-R-005', model: 'MPS Reader Lite', currentVersion: 'v2.3.8', targetVersion: 'v2.5.0', status: '실패', lastUpdate: '2026-05-24 07:55' },
  { id: 'MPS-R-006', model: 'MPS Reader Pro', currentVersion: 'v2.4.1', targetVersion: 'v2.5.0', status: '대기', lastUpdate: '2026-05-23 14:20' },
];

const ROLLOUT = {
  total: 1234,
  completed: 892,
  inProgress: 45,
  failed: 3,
  pending: 294,
};

export default function OtaPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [showConfirm, setShowConfirm] = useState<'deploy' | 'rollback' | null>(null);

  const progressPercent = Math.round((ROLLOUT.completed / ROLLOUT.total) * 100);

  return (
    <main className="min-h-screen bg-slate-50" aria-label="펌웨어 업데이트">
      <DomainHeader
        title="펌웨어 업데이트"
        subtitle="OTA 업데이트를 관리하고 기기별 배포 현황을 확인하세요"
        icon={RefreshCw}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'admin' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8 space-y-6">
        {/* 현재 버전 & 신규 버전 */}
        <div className="grid grid-cols-2 gap-4">
          <div className="rounded-2xl border border-slate-200 bg-white p-5">
            <div className="flex items-center gap-3 mb-3">
              <div className="h-10 w-10 rounded-xl bg-emerald-50 flex items-center justify-center">
                <Shield className="h-5 w-5 text-emerald-600" />
              </div>
              <div>
                <p className="text-xs font-bold text-slate-400">현재 안정 버전</p>
                <p className="text-lg font-black text-slate-900">v2.4.1</p>
              </div>
            </div>
            <p className="text-xs text-slate-500">배포일: 2026-04-15 | 적용 기기: 1,189대</p>
          </div>

          <div className="rounded-2xl border border-sky-200 bg-sky-50/50 p-5">
            <div className="flex items-center gap-3 mb-3">
              <div className="h-10 w-10 rounded-xl bg-sky-100 flex items-center justify-center">
                <RefreshCw className="h-5 w-5 text-sky-600" />
              </div>
              <div>
                <p className="text-xs font-bold text-sky-500">신규 업데이트</p>
                <p className="text-lg font-black text-slate-900">v2.5.0</p>
              </div>
            </div>
            <p className="text-xs text-slate-500">배포 시작: 2026-05-22</p>
          </div>
        </div>

        {/* 변경 사항 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <h2 className="text-sm font-bold text-slate-900 mb-3">v2.5.0 변경 사항</h2>
          <ul className="space-y-2">
            {CHANGELOG.map((item, i) => (
              <li key={i} className="flex items-start gap-2 text-xs text-slate-600">
                <CheckCircle className="h-4 w-4 text-emerald-500 mt-0.5 flex-shrink-0" />
                <span>{item}</span>
              </li>
            ))}
          </ul>
        </div>

        {/* 롤아웃 현황 */}
        <div className="rounded-2xl border border-slate-200 bg-white p-5">
          <h2 className="text-sm font-bold text-slate-900 mb-4">배포 현황</h2>
          <div className="grid grid-cols-5 gap-3 mb-4">
            {[
              { label: '대상 기기', value: ROLLOUT.total.toLocaleString(), color: 'text-slate-700' },
              { label: '완료', value: ROLLOUT.completed.toLocaleString(), color: 'text-emerald-600' },
              { label: '진행중', value: ROLLOUT.inProgress.toString(), color: 'text-sky-600' },
              { label: '실패', value: ROLLOUT.failed.toString(), color: 'text-red-600' },
              { label: '대기', value: ROLLOUT.pending.toString(), color: 'text-slate-500' },
            ].map(stat => (
              <div key={stat.label} className="rounded-xl border border-slate-100 bg-slate-50 p-3 text-center">
                <p className="text-[10px] font-bold text-slate-400 mb-1">{stat.label}</p>
                <p className={`text-xl font-black ${stat.color}`}>{stat.value}</p>
              </div>
            ))}
          </div>

          {/* 프로그레스 바 */}
          <div className="mb-2">
            <div className="flex items-center justify-between mb-1">
              <span className="text-xs font-semibold text-slate-600">전체 진행률</span>
              <span className="text-xs font-bold text-sky-600">{progressPercent}%</span>
            </div>
            <div className="h-3 bg-slate-100 rounded-full overflow-hidden">
              <div
                className="h-full bg-sky-600 rounded-full transition-all"
                style={{ width: `${progressPercent}%` }}
                role="progressbar"
                aria-valuenow={progressPercent}
                aria-valuemin={0}
                aria-valuemax={100}
                aria-label={`업데이트 진행률 ${progressPercent}%`}
              />
            </div>
          </div>
        </div>

        {/* 기기별 상태 테이블 */}
        <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
          <div className="px-5 py-4 border-b border-slate-100">
            <h2 className="text-sm font-bold text-slate-900">기기별 업데이트 상태</h2>
          </div>
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-slate-100 bg-slate-50">
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">기기 ID</th>
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">모델</th>
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">현재 버전</th>
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">대상 버전</th>
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">상태</th>
                <th className="px-5 py-3 text-[11px] font-bold text-slate-500">최종 업데이트</th>
              </tr>
            </thead>
            <tbody>
              {MOCK_DEVICES.map(device => {
                const StatusIcon = STATUS_ICONS[device.status];
                return (
                  <tr key={device.id} className="border-b border-slate-50 hover:bg-slate-50 transition-colors">
                    <td className="px-5 py-3 text-xs font-mono font-semibold text-slate-700">{device.id}</td>
                    <td className="px-5 py-3 text-xs text-slate-600">{device.model}</td>
                    <td className="px-5 py-3 text-xs font-mono text-slate-600">{device.currentVersion}</td>
                    <td className="px-5 py-3 text-xs font-mono text-slate-600">{device.targetVersion}</td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold border ${STATUS_STYLES[device.status]}`}>
                        <StatusIcon className={`h-3 w-3 ${device.status === '진행중' ? 'animate-spin' : ''}`} />
                        {device.status}
                      </span>
                    </td>
                    <td className="px-5 py-3 text-xs text-slate-500 font-mono">{device.lastUpdate}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {/* 액션 버튼 */}
        <div className="flex items-center gap-3">
          <button
            onClick={() => setShowConfirm('deploy')}
            className="flex-1 flex items-center justify-center gap-2 rounded-xl bg-sky-600 px-6 py-3 text-sm font-bold text-white hover:bg-sky-700 transition-colors"
          >
            <Server className="h-4 w-4" /> 전체 배포 시작
          </button>
          <button
            onClick={() => setShowConfirm('rollback')}
            className="flex-1 flex items-center justify-center gap-2 rounded-xl border border-red-200 bg-red-50 px-6 py-3 text-sm font-bold text-red-700 hover:bg-red-100 transition-colors"
          >
            <RefreshCw className="h-4 w-4" /> 롤백
          </button>
        </div>

        {/* 확인 다이얼로그 */}
        {showConfirm && (
          <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" role="dialog" aria-modal="true">
            <div className="rounded-2xl bg-white p-6 max-w-md w-full mx-4 shadow-xl">
              <div className="flex items-center gap-3 mb-4">
                <div className={`h-10 w-10 rounded-xl flex items-center justify-center ${showConfirm === 'deploy' ? 'bg-sky-50 text-sky-600' : 'bg-red-50 text-red-600'}`}>
                  {showConfirm === 'deploy' ? <Server className="h-5 w-5" /> : <AlertTriangle className="h-5 w-5" />}
                </div>
                <h3 className="text-base font-bold text-slate-900">
                  {showConfirm === 'deploy' ? '전체 배포를 시작하시겠습니까?' : '이전 버전으로 롤백하시겠습니까?'}
                </h3>
              </div>
              <p className="text-xs text-slate-500 mb-6">
                {showConfirm === 'deploy'
                  ? `대기 중인 ${ROLLOUT.pending}대 기기에 v2.5.0 펌웨어를 배포합니다. 이 작업은 취소할 수 없습니다.`
                  : '모든 기기를 v2.4.1로 롤백합니다. 진행 중인 업데이트가 중단됩니다.'}
              </p>
              <div className="flex justify-end gap-2">
                <button onClick={() => setShowConfirm(null)} className="px-4 py-2 rounded-xl border border-slate-200 text-xs font-semibold text-slate-600 hover:bg-slate-50">취소</button>
                <button
                  onClick={() => setShowConfirm(null)}
                  className={`px-4 py-2 rounded-xl text-xs font-bold text-white ${showConfirm === 'deploy' ? 'bg-sky-600 hover:bg-sky-700' : 'bg-red-600 hover:bg-red-700'}`}
                >
                  {showConfirm === 'deploy' ? '배포 시작' : '롤백 실행'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
