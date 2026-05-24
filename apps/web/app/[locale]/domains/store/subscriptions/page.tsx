'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import {
  Truck,
  CalendarClock,
  Pause,
  Play,
  Trash2,
  Plus,
  ChevronDown,
  Package,
  CheckCircle,
  XCircle,
  AlertCircle,
} from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

/* ── 타입 ── */
type DeliveryCycle = '매주' | '격주' | '매월';
type SubscriptionStatus = '활성' | '일시정지' | '해지';

interface Subscription {
  id: string;
  productName: string;
  productImage: string;
  cycle: DeliveryCycle;
  nextDelivery: string;
  status: SubscriptionStatus;
  price: number;
  startDate: string;
  cancelledDate?: string;
}

const CYCLE_OPTIONS: DeliveryCycle[] = ['매주', '격주', '매월'];

const INITIAL_SUBSCRIPTIONS: Subscription[] = [
  {
    id: 'sub-1',
    productName: '혈당 카트리지 (10개입)',
    productImage: '🩸',
    cycle: '격주',
    nextDelivery: '2026-06-07',
    status: '활성',
    price: 29000,
    startDate: '2026-01-15',
  },
  {
    id: 'sub-2',
    productName: '채혈침 세트 (100개입)',
    productImage: '💊',
    cycle: '매월',
    nextDelivery: '2026-06-15',
    status: '일시정지',
    price: 12000,
    startDate: '2026-02-20',
  },
  {
    id: 'sub-3',
    productName: '콜레스테롤 카트리지 (10개입)',
    productImage: '💉',
    cycle: '매월',
    nextDelivery: '',
    status: '해지',
    price: 35000,
    startDate: '2025-11-01',
    cancelledDate: '2026-03-15',
  },
];

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const [y, m, d] = dateStr.split('-');
  return `${y}년 ${Number(m)}월 ${Number(d)}일`;
}

function StatusBadge({ status }: { status: SubscriptionStatus }) {
  const config = {
    '활성': { color: 'bg-emerald-50 text-emerald-700 border-emerald-200', icon: CheckCircle },
    '일시정지': { color: 'bg-amber-50 text-amber-700 border-amber-200', icon: AlertCircle },
    '해지': { color: 'bg-slate-100 text-slate-500 border-slate-200', icon: XCircle },
  };
  const { color, icon: Icon } = config[status];
  return (
    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold border ${color}`}>
      <Icon className="h-3 w-3" />
      {status}
    </span>
  );
}

export default function SubscriptionsPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [subscriptions, setSubscriptions] = useState<Subscription[]>(INITIAL_SUBSCRIPTIONS);
  const [cycleDropdownOpen, setCycleDropdownOpen] = useState<string | null>(null);
  const [confirmAction, setConfirmAction] = useState<{ id: string; type: '일시정지' | '재개' | '해지' } | null>(null);

  const activeSubscriptions = subscriptions.filter((s) => s.status !== '해지');
  const cancelledSubscriptions = subscriptions.filter((s) => s.status === '해지');

  const handleTogglePause = (id: string) => {
    const sub = subscriptions.find((s) => s.id === id);
    if (!sub) return;
    setConfirmAction({
      id,
      type: sub.status === '활성' ? '일시정지' : '재개',
    });
  };

  const handleCancel = (id: string) => {
    setConfirmAction({ id, type: '해지' });
  };

  const executeAction = () => {
    if (!confirmAction) return;
    setSubscriptions((prev) =>
      prev.map((sub) => {
        if (sub.id !== confirmAction.id) return sub;
        if (confirmAction.type === '일시정지') {
          return { ...sub, status: '일시정지' as SubscriptionStatus };
        }
        if (confirmAction.type === '재개') {
          return { ...sub, status: '활성' as SubscriptionStatus };
        }
        if (confirmAction.type === '해지') {
          return {
            ...sub,
            status: '해지' as SubscriptionStatus,
            nextDelivery: '',
            cancelledDate: new Date().toISOString().slice(0, 10),
          };
        }
        return sub;
      })
    );
    setConfirmAction(null);
  };

  const handleCycleChange = (id: string, newCycle: DeliveryCycle) => {
    setSubscriptions((prev) =>
      prev.map((sub) => (sub.id === id ? { ...sub, cycle: newCycle } : sub))
    );
    setCycleDropdownOpen(null);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="정기배송 관리">
      <DomainHeader
        title="정기배송 관리"
        subtitle="정기배송 구독을 관리하세요"
        icon={Truck}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-4xl mx-auto px-6 py-8">
        {/* 활성 구독 목록 */}
        <section aria-label="활성 구독 목록">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-slate-900">활성 구독</h2>
            <span className="text-xs font-semibold text-slate-400">{activeSubscriptions.length}건</span>
          </div>

          {activeSubscriptions.length === 0 ? (
            <div className="rounded-2xl border border-slate-200 bg-white p-12 text-center">
              <Package className="h-12 w-12 text-slate-300 mx-auto mb-3" />
              <p className="text-sm text-slate-500 mb-4">활성 구독이 없습니다.</p>
              <Link
                href={`${localePrefix}/domains/store`}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
              >
                <Plus className="h-4 w-4" />
                새 정기배송 추가
              </Link>
            </div>
          ) : (
            <div className="space-y-4">
              {activeSubscriptions.map((sub) => (
                <div
                  key={sub.id}
                  className="rounded-2xl border border-slate-200 bg-white overflow-hidden hover:shadow-sm transition-shadow"
                >
                  <div className="p-6">
                    <div className="flex items-start gap-4">
                      {/* 상품 이미지 */}
                      <div className="h-16 w-16 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-center flex-shrink-0">
                        <span className="text-3xl">{sub.productImage}</span>
                      </div>

                      {/* 상품 정보 */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <h3 className="text-sm font-bold text-slate-900">{sub.productName}</h3>
                          <StatusBadge status={sub.status} />
                        </div>
                        <p className="text-lg font-bold text-slate-900 mb-2">
                          ₩{sub.price.toLocaleString()}
                          <span className="text-xs font-normal text-slate-400 ml-1">/ 회</span>
                        </p>

                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                          {/* 배송 주기 */}
                          <div>
                            <p className="text-[10px] text-slate-400 font-semibold uppercase mb-1">배송 주기</p>
                            <div className="relative">
                              <button
                                onClick={() =>
                                  setCycleDropdownOpen(cycleDropdownOpen === sub.id ? null : sub.id)
                                }
                                disabled={sub.status === '일시정지'}
                                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg border text-sm font-semibold transition-colors ${
                                  sub.status === '일시정지'
                                    ? 'border-slate-100 bg-slate-50 text-slate-400 cursor-not-allowed'
                                    : 'border-slate-200 bg-white text-slate-700 hover:border-sky-300'
                                }`}
                              >
                                <CalendarClock className="h-3.5 w-3.5" />
                                {sub.cycle}
                                <ChevronDown className="h-3 w-3" />
                              </button>
                              {cycleDropdownOpen === sub.id && (
                                <div className="absolute left-0 top-full mt-1 w-32 bg-white border border-slate-200 rounded-xl shadow-lg z-10 p-1">
                                  {CYCLE_OPTIONS.map((cycle) => (
                                    <button
                                      key={cycle}
                                      onClick={() => handleCycleChange(sub.id, cycle)}
                                      className={`w-full text-left px-3 py-2 text-sm rounded-lg transition-colors ${
                                        sub.cycle === cycle
                                          ? 'bg-sky-50 text-sky-700 font-bold'
                                          : 'text-slate-700 hover:bg-slate-50'
                                      }`}
                                    >
                                      {cycle}
                                    </button>
                                  ))}
                                </div>
                              )}
                            </div>
                          </div>

                          {/* 다음 배송일 */}
                          <div>
                            <p className="text-[10px] text-slate-400 font-semibold uppercase mb-1">다음 배송일</p>
                            <p className="text-sm font-semibold text-slate-700">
                              {sub.status === '일시정지' ? (
                                <span className="text-amber-600">일시정지 중</span>
                              ) : (
                                formatDate(sub.nextDelivery)
                              )}
                            </p>
                          </div>

                          {/* 구독 시작일 */}
                          <div>
                            <p className="text-[10px] text-slate-400 font-semibold uppercase mb-1">구독 시작일</p>
                            <p className="text-sm text-slate-500">{formatDate(sub.startDate)}</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* 액션 버튼 */}
                  <div className="px-6 py-3 bg-slate-50 border-t border-slate-100 flex items-center gap-2">
                    <button
                      onClick={() => handleTogglePause(sub.id)}
                      className={`flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-xs font-bold transition-colors ${
                        sub.status === '활성'
                          ? 'bg-amber-50 text-amber-700 border border-amber-200 hover:bg-amber-100'
                          : 'bg-emerald-50 text-emerald-700 border border-emerald-200 hover:bg-emerald-100'
                      }`}
                    >
                      {sub.status === '활성' ? (
                        <>
                          <Pause className="h-3.5 w-3.5" /> 일시정지
                        </>
                      ) : (
                        <>
                          <Play className="h-3.5 w-3.5" /> 재개
                        </>
                      )}
                    </button>
                    <button
                      onClick={() =>
                        setCycleDropdownOpen(cycleDropdownOpen === sub.id ? null : sub.id)
                      }
                      disabled={sub.status === '일시정지'}
                      className={`flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-xs font-bold transition-colors ${
                        sub.status === '일시정지'
                          ? 'bg-slate-100 text-slate-400 cursor-not-allowed'
                          : 'bg-sky-50 text-sky-700 border border-sky-200 hover:bg-sky-100'
                      }`}
                    >
                      <CalendarClock className="h-3.5 w-3.5" /> 주기 변경
                    </button>
                    <button
                      onClick={() => handleCancel(sub.id)}
                      className="flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-xs font-bold bg-white text-rose-600 border border-rose-200 hover:bg-rose-50 transition-colors ml-auto"
                    >
                      <Trash2 className="h-3.5 w-3.5" /> 해지
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* 해지된 구독 */}
        {cancelledSubscriptions.length > 0 && (
          <section aria-label="해지된 구독 목록" className="mt-10">
            <h2 className="text-base font-bold text-slate-900 mb-4">해지된 구독</h2>
            <div className="space-y-3">
              {cancelledSubscriptions.map((sub) => (
                <div
                  key={sub.id}
                  className="rounded-2xl border border-slate-200 bg-white p-5 opacity-60"
                >
                  <div className="flex items-center gap-4">
                    <div className="h-12 w-12 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-center flex-shrink-0">
                      <span className="text-2xl grayscale">{sub.productImage}</span>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-0.5">
                        <h3 className="text-sm font-bold text-slate-700">{sub.productName}</h3>
                        <StatusBadge status={sub.status} />
                      </div>
                      <p className="text-xs text-slate-400">
                        {formatDate(sub.startDate)} ~ {formatDate(sub.cancelledDate || '')} | {sub.cycle} |{' '}
                        ₩{sub.price.toLocaleString()}/회
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 새 정기배송 추가 */}
        <section className="mt-10">
          <Link
            href={`${localePrefix}/domains/store`}
            className="flex items-center justify-center gap-2 w-full py-4 rounded-2xl border-2 border-dashed border-slate-300 bg-white text-sm font-bold text-slate-500 hover:border-sky-400 hover:text-sky-600 hover:bg-sky-50 transition-all"
          >
            <Plus className="h-5 w-5" />
            새 정기배송 추가
          </Link>
        </section>

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-8">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: 'MPS 스토어', path: '/domains/store', desc: '상품 둘러보기' },
              { name: '건강 측정', path: '/domains/measure', desc: '카트리지로 측정' },
              { name: '보상', path: '/domains/reward', desc: '포인트 사용' },
            ].map((d) => (
              <Link
                key={d.path}
                href={`${localePrefix}${d.path}`}
                className="flex items-center gap-3 p-3 rounded-xl border border-slate-200 bg-white hover:border-sky-300 hover:shadow-sm transition-all group"
              >
                <div className="h-8 w-8 rounded-lg bg-sky-50 flex items-center justify-center text-sky-600 group-hover:bg-sky-100">
                  <span className="text-sm">→</span>
                </div>
                <div>
                  <span className="text-sm font-semibold text-slate-700 group-hover:text-sky-600">
                    {d.name}
                  </span>
                  <p className="text-[11px] text-slate-400">{d.desc}</p>
                </div>
              </Link>
            ))}
          </div>
        </section>
      </div>

      {/* 확인 모달 */}
      {confirmAction && (
        <div className="fixed inset-0 z-50 flex items-center justify-center" onClick={() => setConfirmAction(null)}>
          <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" />
          <div
            className="relative bg-white rounded-2xl shadow-2xl p-6 w-96 max-w-[90vw]"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-lg font-bold text-slate-900 mb-2">
              {confirmAction.type === '해지'
                ? '정기배송을 해지하시겠습니까?'
                : confirmAction.type === '일시정지'
                ? '정기배송을 일시정지하시겠습니까?'
                : '정기배송을 재개하시겠습니까?'}
            </h3>
            <p className="text-sm text-slate-500 mb-6">
              {confirmAction.type === '해지'
                ? '해지 후에는 동일 상품을 다시 구독하셔야 합니다.'
                : confirmAction.type === '일시정지'
                ? '일시정지 중에는 배송이 진행되지 않으며, 언제든 재개할 수 있습니다.'
                : '다음 배송일부터 정기배송이 재개됩니다.'}
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setConfirmAction(null)}
                className="flex-1 py-2.5 rounded-xl border border-slate-200 text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
              >
                취소
              </button>
              <button
                onClick={executeAction}
                className={`flex-1 py-2.5 rounded-xl text-sm font-bold text-white transition-colors ${
                  confirmAction.type === '해지'
                    ? 'bg-rose-600 hover:bg-rose-700'
                    : confirmAction.type === '일시정지'
                    ? 'bg-amber-600 hover:bg-amber-700'
                    : 'bg-emerald-600 hover:bg-emerald-700'
                }`}
              >
                {confirmAction.type === '해지'
                  ? '해지하기'
                  : confirmAction.type === '일시정지'
                  ? '일시정지'
                  : '재개하기'}
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
