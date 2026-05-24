'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Package, ChevronDown, ChevronUp, CheckCircle, Clock, Truck, Box, CreditCard } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type OrderStatus = '준비중' | '배송중' | '배송완료';

interface OrderProduct {
  name: string;
  image: string;
  quantity: number;
  price: number;
}

interface Order {
  id: string;
  orderNumber: string;
  date: string;
  status: OrderStatus;
  total: number;
  products: OrderProduct[];
  trackingNumber?: string;
  deliveryDate?: string;
}

const TIMELINE_STEPS = ['주문접수', '결제완료', '배송준비', '배송중', '배송완료'] as const;

const STATUS_TO_STEP: Record<OrderStatus, number> = {
  '준비중': 2,
  '배송중': 3,
  '배송완료': 4,
};

const STATUS_STYLES: Record<OrderStatus, string> = {
  '준비중': 'bg-amber-50 text-amber-700 border-amber-200',
  '배송중': 'bg-sky-50 text-sky-700 border-sky-200',
  '배송완료': 'bg-emerald-50 text-emerald-700 border-emerald-200',
};

const STEP_ICONS = [CreditCard, CheckCircle, Box, Truck, Package];

const MOCK_ORDERS: Order[] = [
  {
    id: 'ord-1',
    orderNumber: 'MPS-20260520-001',
    date: '2026-05-20',
    status: '배송중',
    total: 64000,
    trackingNumber: 'KR1234567890',
    products: [
      { name: '혈당 카트리지 (10개입)', image: '🩸', quantity: 1, price: 29000 },
      { name: '콜레스테롤 카트리지 (10개입)', image: '💉', quantity: 1, price: 35000 },
    ],
  },
  {
    id: 'ord-2',
    orderNumber: 'MPS-20260515-004',
    date: '2026-05-15',
    status: '배송완료',
    total: 302000,
    deliveryDate: '2026-05-18',
    products: [
      { name: 'MPS 리더기 Pro', image: '📱', quantity: 1, price: 290000 },
      { name: '채혈침 세트 (100개입)', image: '💊', quantity: 1, price: 12000 },
    ],
  },
  {
    id: 'ord-3',
    orderNumber: 'MPS-20260523-002',
    date: '2026-05-23',
    status: '준비중',
    total: 134000,
    products: [
      { name: '종합검사 카트리지 세트', image: '📦', quantity: 1, price: 89000 },
      { name: '당화혈색소 카트리지 (5개입)', image: '🔬', quantity: 1, price: 45000 },
    ],
  },
];

export default function OrdersPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [expandedOrder, setExpandedOrder] = useState<string | null>(MOCK_ORDERS[0].id);

  const toggleExpand = (orderId: string) => {
    setExpandedOrder(prev => (prev === orderId ? null : orderId));
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="주문 내역">
      <DomainHeader
        title="주문 내역"
        subtitle="주문 현황과 배송 상태를 확인하세요"
        icon={Package}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-4xl mx-auto px-6 py-8">
        {/* 주문 요약 */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          {[
            { label: '전체 주문', value: MOCK_ORDERS.length, color: 'bg-slate-900' },
            { label: '배송중', value: MOCK_ORDERS.filter(o => o.status === '배송중').length, color: 'bg-sky-600' },
            { label: '배송완료', value: MOCK_ORDERS.filter(o => o.status === '배송완료').length, color: 'bg-emerald-600' },
          ].map(stat => (
            <div key={stat.label} className="rounded-2xl border border-slate-200 bg-white p-4 text-center">
              <p className="text-xs text-slate-500 font-medium mb-1">{stat.label}</p>
              <p className={`text-2xl font-bold text-slate-900`}>{stat.value}</p>
            </div>
          ))}
        </div>

        {/* 주문 목록 */}
        <div className="space-y-4">
          {MOCK_ORDERS.map(order => {
            const isExpanded = expandedOrder === order.id;
            const currentStep = STATUS_TO_STEP[order.status];

            return (
              <div key={order.id} className="rounded-2xl border border-slate-200 bg-white overflow-hidden transition-shadow hover:shadow-sm">
                {/* 주문 헤더 */}
                <button
                  onClick={() => toggleExpand(order.id)}
                  className="w-full flex items-center justify-between p-5 text-left"
                  aria-expanded={isExpanded}
                  aria-label={`주문 ${order.orderNumber} 상세 보기`}
                >
                  <div className="flex items-center gap-4">
                    <div className="flex -space-x-2">
                      {order.products.slice(0, 2).map((p, i) => (
                        <span key={i} className="h-10 w-10 rounded-xl bg-slate-50 border border-slate-200 flex items-center justify-center text-lg">
                          {p.image}
                        </span>
                      ))}
                      {order.products.length > 2 && (
                        <span className="h-10 w-10 rounded-xl bg-slate-100 border border-slate-200 flex items-center justify-center text-xs font-bold text-slate-500">
                          +{order.products.length - 2}
                        </span>
                      )}
                    </div>
                    <div>
                      <p className="text-sm font-bold text-slate-900">{order.orderNumber}</p>
                      <p className="text-xs text-slate-400">{order.date}</p>
                    </div>
                  </div>

                  <div className="flex items-center gap-4">
                    <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${STATUS_STYLES[order.status]}`}>
                      {order.status}
                    </span>
                    <span className="text-sm font-bold text-slate-900">₩{order.total.toLocaleString()}</span>
                    {isExpanded ? <ChevronUp className="h-4 w-4 text-slate-400" /> : <ChevronDown className="h-4 w-4 text-slate-400" />}
                  </div>
                </button>

                {/* 확장 상세 */}
                {isExpanded && (
                  <div className="border-t border-slate-100 px-5 pb-5">
                    {/* 상태 타임라인 */}
                    <div className="py-6">
                      <div className="flex items-center justify-between relative">
                        {/* 배경 라인 */}
                        <div className="absolute top-4 left-0 right-0 h-0.5 bg-slate-200" />
                        <div
                          className="absolute top-4 left-0 h-0.5 bg-sky-600 transition-all"
                          style={{ width: `${(currentStep / (TIMELINE_STEPS.length - 1)) * 100}%` }}
                        />

                        {TIMELINE_STEPS.map((step, index) => {
                          const StepIcon = STEP_ICONS[index];
                          const isComplete = index <= currentStep;
                          const isCurrent = index === currentStep;

                          return (
                            <div key={step} className="relative flex flex-col items-center z-10">
                              <div className={`h-8 w-8 rounded-full flex items-center justify-center transition-colors ${
                                isCurrent
                                  ? 'bg-sky-600 text-white ring-4 ring-sky-100'
                                  : isComplete
                                    ? 'bg-sky-600 text-white'
                                    : 'bg-slate-100 text-slate-400 border border-slate-200'
                              }`}>
                                <StepIcon className="h-4 w-4" />
                              </div>
                              <span className={`mt-2 text-[11px] font-medium whitespace-nowrap ${
                                isComplete ? 'text-slate-900' : 'text-slate-400'
                              }`}>{step}</span>
                            </div>
                          );
                        })}
                      </div>
                    </div>

                    {/* 주문 상품 목록 */}
                    <div className="space-y-2">
                      {order.products.map((product, idx) => (
                        <div key={idx} className="flex items-center gap-3 p-3 rounded-xl bg-slate-50">
                          <span className="text-2xl">{product.image}</span>
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium text-slate-900 truncate">{product.name}</p>
                            <p className="text-xs text-slate-500">수량: {product.quantity}</p>
                          </div>
                          <span className="text-sm font-bold text-slate-700">₩{(product.price * product.quantity).toLocaleString()}</span>
                        </div>
                      ))}
                    </div>

                    {/* 배송 정보 */}
                    <div className="mt-4 pt-4 border-t border-slate-100 flex items-center justify-between">
                      <div className="text-xs text-slate-500">
                        {order.trackingNumber && (
                          <span>운송장: <span className="font-mono font-semibold text-slate-700">{order.trackingNumber}</span></span>
                        )}
                        {order.deliveryDate && (
                          <span>배송 완료일: <span className="font-semibold text-slate-700">{order.deliveryDate}</span></span>
                        )}
                        {!order.trackingNumber && !order.deliveryDate && (
                          <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            배송 준비 중입니다
                          </span>
                        )}
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-slate-400">총 결제금액</p>
                        <p className="text-base font-bold text-slate-900">₩{order.total.toLocaleString()}</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
