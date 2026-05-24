'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import {
  ShoppingCart,
  ChevronRight,
  ChevronLeft,
  CheckCircle,
  CreditCard,
  MapPin,
  ClipboardList,
  Truck,
} from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

/* ── 목 데이터 ── */
const MOCK_CART = [
  { id: '1', name: '혈당 카트리지 (10개입)', price: 29000, quantity: 2, image: '🩸' },
  { id: '3', name: '당화혈색소 카트리지 (5개입)', price: 45000, quantity: 1, image: '🔬' },
  { id: '6', name: '채혈침 세트 (100개입)', price: 12000, quantity: 1, image: '💊' },
];

const PAYMENT_METHODS = [
  { id: 'card', label: '카드결제', icon: '💳' },
  { id: 'toss', label: '토스페이', icon: '🟦' },
  { id: 'transfer', label: '계좌이체', icon: '🏦' },
  { id: 'point', label: '포인트 사용', icon: '🪙' },
] as const;

const STEPS = [
  { step: 1, label: '배송 정보', icon: MapPin },
  { step: 2, label: '결제 방법', icon: CreditCard },
  { step: 3, label: '주문 확인', icon: ClipboardList },
] as const;

interface ShippingForm {
  name: string;
  phone: string;
  zipcode: string;
  address: string;
  addressDetail: string;
  memo: string;
}

export default function CheckoutPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [currentStep, setCurrentStep] = useState<1 | 2 | 3>(1);
  const [orderComplete, setOrderComplete] = useState(false);
  const [orderNumber, setOrderNumber] = useState('');

  const [shipping, setShipping] = useState<ShippingForm>({
    name: '홍길동',
    phone: '010-1234-5678',
    zipcode: '06234',
    address: '서울특별시 강남구 테헤란로 123',
    addressDetail: 'MPS 타워 4층',
    memo: '',
  });

  const [paymentMethod, setPaymentMethod] = useState('card');
  const [cartItems] = useState(MOCK_CART);

  const subtotal = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
  const shippingFee = subtotal >= 50000 ? 0 : 3000;
  const total = subtotal + shippingFee;

  const handleOrder = () => {
    const num = `MPS-${Date.now().toString(36).toUpperCase()}`;
    setOrderNumber(num);
    setOrderComplete(true);
  };

  const goNext = () => {
    if (currentStep < 3) setCurrentStep((currentStep + 1) as 1 | 2 | 3);
  };
  const goPrev = () => {
    if (currentStep > 1) setCurrentStep((currentStep - 1) as 1 | 2 | 3);
  };

  /* ── 주문 완료 ── */
  if (orderComplete) {
    return (
      <main className="min-h-screen bg-slate-50" aria-label="주문 완료">
        <DomainHeader
          title="주문 완료"
          subtitle="주문이 성공적으로 접수되었습니다"
          icon={ShoppingCart}
          user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
          onLogout={() => signOut({ callbackUrl: '/login' })}
        />
        <div className="max-w-lg mx-auto px-6 py-20 text-center">
          <div className="mb-6 animate-bounce">
            <CheckCircle className="h-20 w-20 text-emerald-500 mx-auto" />
          </div>
          <h2 className="text-2xl font-bold text-slate-900 mb-2">주문이 완료되었습니다</h2>
          <p className="text-sm text-slate-500 mb-1">주문번호</p>
          <p className="text-lg font-mono font-bold text-sky-600 mb-6">{orderNumber}</p>
          <p className="text-sm text-slate-400 mb-8">
            배송 준비가 시작되며, 배송 상태는 주문 내역에서 확인하실 수 있습니다.
          </p>
          <div className="flex flex-col gap-3">
            <Link
              href={`${localePrefix}/domains/store`}
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-sky-600 text-white font-bold hover:bg-sky-700 transition-colors"
            >
              주문 내역 보기
            </Link>
            <Link
              href={`${localePrefix}/domains/store`}
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl border border-slate-200 bg-white text-slate-700 font-semibold hover:bg-slate-50 transition-colors"
            >
              스토어로 돌아가기
            </Link>
          </div>
        </div>
      </main>
    );
  }

  /* ── 스텝 인디케이터 ── */
  const StepIndicator = () => (
    <div className="flex items-center justify-center gap-0 mb-10">
      {STEPS.map((s, idx) => {
        const StepIcon = s.icon;
        const isActive = currentStep === s.step;
        const isDone = currentStep > s.step;
        return (
          <React.Fragment key={s.step}>
            {idx > 0 && (
              <div
                className={`w-16 h-0.5 ${
                  isDone ? 'bg-sky-500' : 'bg-slate-200'
                } transition-colors`}
              />
            )}
            <div className="flex flex-col items-center gap-1.5">
              <div
                className={`h-10 w-10 rounded-full flex items-center justify-center text-sm font-bold transition-all ${
                  isDone
                    ? 'bg-sky-500 text-white'
                    : isActive
                    ? 'bg-sky-600 text-white ring-4 ring-sky-100'
                    : 'bg-slate-100 text-slate-400'
                }`}
              >
                {isDone ? <CheckCircle className="h-5 w-5" /> : <StepIcon className="h-5 w-5" />}
              </div>
              <span
                className={`text-xs font-semibold ${
                  isActive || isDone ? 'text-slate-900' : 'text-slate-400'
                }`}
              >
                {s.label}
              </span>
            </div>
          </React.Fragment>
        );
      })}
    </div>
  );

  /* ── 오른쪽 요약 패널 ── */
  const OrderSummaryPanel = () => (
    <div className="rounded-2xl border border-slate-200 bg-white p-6">
      <h3 className="text-sm font-bold text-slate-900 mb-4">주문 요약</h3>
      <div className="space-y-3 mb-4">
        {cartItems.map((item) => (
          <div key={item.id} className="flex items-center gap-3">
            <span className="text-xl">{item.image}</span>
            <div className="flex-1 min-w-0">
              <p className="text-sm text-slate-700 truncate">{item.name}</p>
              <p className="text-xs text-slate-400">{item.quantity}개</p>
            </div>
            <span className="text-sm font-semibold text-slate-900">
              ₩{(item.price * item.quantity).toLocaleString()}
            </span>
          </div>
        ))}
      </div>
      <div className="border-t border-slate-100 pt-3 space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-slate-500">상품 합계</span>
          <span className="text-slate-700">₩{subtotal.toLocaleString()}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-slate-500">배송비</span>
          <span className={shippingFee === 0 ? 'text-emerald-600 font-semibold' : 'text-slate-700'}>
            {shippingFee === 0 ? '무료' : `₩${shippingFee.toLocaleString()}`}
          </span>
        </div>
        <div className="border-t border-slate-100 pt-2 flex justify-between">
          <span className="text-sm font-bold text-slate-900">총 결제금액</span>
          <span className="text-lg font-bold text-sky-600">₩{total.toLocaleString()}</span>
        </div>
      </div>
    </div>
  );

  return (
    <main className="min-h-screen bg-slate-50" aria-label="주문/결제">
      <DomainHeader
        title="주문/결제"
        subtitle="안전하고 간편한 결제"
        icon={ShoppingCart}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8">
        <StepIndicator />

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 메인 콘텐츠 */}
          <div className="lg:col-span-2">
            {/* ── 1단계: 배송 정보 ── */}
            {currentStep === 1 && (
              <div className="rounded-2xl border border-slate-200 bg-white p-6">
                <div className="flex items-center gap-2 mb-6">
                  <Truck className="h-5 w-5 text-sky-600" />
                  <h2 className="text-lg font-bold text-slate-900">배송 정보</h2>
                </div>
                <div className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                        수령인 <span className="text-rose-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={shipping.name}
                        onChange={(e) => setShipping({ ...shipping, name: e.target.value })}
                        className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                        연락처 <span className="text-rose-500">*</span>
                      </label>
                      <input
                        type="tel"
                        value={shipping.phone}
                        onChange={(e) => setShipping({ ...shipping, phone: e.target.value })}
                        className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                      우편번호 <span className="text-rose-500">*</span>
                    </label>
                    <div className="flex gap-2">
                      <input
                        type="text"
                        value={shipping.zipcode}
                        onChange={(e) => setShipping({ ...shipping, zipcode: e.target.value })}
                        className="w-32 px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                      />
                      <button className="px-4 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-semibold text-slate-600 hover:bg-slate-50">
                        주소 검색
                      </button>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                      주소 <span className="text-rose-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={shipping.address}
                      onChange={(e) => setShipping({ ...shipping, address: e.target.value })}
                      className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                      상세주소
                    </label>
                    <input
                      type="text"
                      value={shipping.addressDetail}
                      onChange={(e) => setShipping({ ...shipping, addressDetail: e.target.value })}
                      className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                      배송 메모
                    </label>
                    <select
                      value={shipping.memo}
                      onChange={(e) => setShipping({ ...shipping, memo: e.target.value })}
                      className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/20 focus:border-sky-500"
                    >
                      <option value="">배송 메모를 선택하세요</option>
                      <option value="문 앞에 놓아주세요">문 앞에 놓아주세요</option>
                      <option value="경비실에 맡겨주세요">경비실에 맡겨주세요</option>
                      <option value="배송 전 연락 부탁드립니다">배송 전 연락 부탁드립니다</option>
                      <option value="부재 시 문 앞에 놓아주세요">부재 시 문 앞에 놓아주세요</option>
                    </select>
                  </div>
                </div>
              </div>
            )}

            {/* ── 2단계: 결제 방법 ── */}
            {currentStep === 2 && (
              <div className="rounded-2xl border border-slate-200 bg-white p-6">
                <div className="flex items-center gap-2 mb-6">
                  <CreditCard className="h-5 w-5 text-sky-600" />
                  <h2 className="text-lg font-bold text-slate-900">결제 방법</h2>
                </div>
                <div className="space-y-3 mb-6">
                  {PAYMENT_METHODS.map((method) => (
                    <label
                      key={method.id}
                      className={`flex items-center gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all ${
                        paymentMethod === method.id
                          ? 'border-sky-500 bg-sky-50'
                          : 'border-slate-200 bg-white hover:border-slate-300'
                      }`}
                    >
                      <input
                        type="radio"
                        name="payment"
                        value={method.id}
                        checked={paymentMethod === method.id}
                        onChange={(e) => setPaymentMethod(e.target.value)}
                        className="sr-only"
                      />
                      <div
                        className={`h-5 w-5 rounded-full border-2 flex items-center justify-center ${
                          paymentMethod === method.id ? 'border-sky-500' : 'border-slate-300'
                        }`}
                      >
                        {paymentMethod === method.id && (
                          <div className="h-2.5 w-2.5 rounded-full bg-sky-500" />
                        )}
                      </div>
                      <span className="text-xl">{method.icon}</span>
                      <span className="text-sm font-semibold text-slate-900">{method.label}</span>
                    </label>
                  ))}
                </div>

                {paymentMethod === 'card' && (
                  <div className="rounded-xl bg-gradient-to-br from-slate-800 to-slate-900 p-6 text-white">
                    <p className="text-xs text-slate-400 mb-4">등록된 카드</p>
                    <p className="text-lg font-mono tracking-widest mb-4">
                      •••• •••• •••• 4242
                    </p>
                    <div className="flex justify-between items-end">
                      <div>
                        <p className="text-[10px] text-slate-400">카드 소유자</p>
                        <p className="text-sm font-semibold">{shipping.name}</p>
                      </div>
                      <div>
                        <p className="text-[10px] text-slate-400">유효기간</p>
                        <p className="text-sm font-semibold">12/28</p>
                      </div>
                    </div>
                  </div>
                )}

                {paymentMethod === 'point' && (
                  <div className="rounded-xl bg-amber-50 border border-amber-200 p-4">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-amber-800 font-semibold">보유 포인트</span>
                      <span className="text-lg font-bold text-amber-700">12,500 P</span>
                    </div>
                    <p className="text-xs text-amber-600 mt-1">
                      부족한 금액은 다른 결제 수단으로 결제됩니다.
                    </p>
                  </div>
                )}
              </div>
            )}

            {/* ── 3단계: 주문 확인 ── */}
            {currentStep === 3 && (
              <div className="space-y-4">
                {/* 상품 목록 */}
                <div className="rounded-2xl border border-slate-200 bg-white p-6">
                  <h3 className="text-sm font-bold text-slate-900 mb-4">주문 상품</h3>
                  <div className="space-y-3">
                    {cartItems.map((item) => (
                      <div
                        key={item.id}
                        className="flex items-center gap-4 p-3 rounded-xl bg-slate-50"
                      >
                        <span className="text-2xl">{item.image}</span>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-slate-900">{item.name}</p>
                          <p className="text-xs text-slate-400">{item.quantity}개</p>
                        </div>
                        <span className="text-sm font-bold text-slate-900">
                          ₩{(item.price * item.quantity).toLocaleString()}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* 배송 정보 확인 */}
                <div className="rounded-2xl border border-slate-200 bg-white p-6">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-sm font-bold text-slate-900">배송 정보</h3>
                    <button
                      onClick={() => setCurrentStep(1)}
                      className="text-xs font-semibold text-sky-600 hover:text-sky-700"
                    >
                      수정
                    </button>
                  </div>
                  <div className="space-y-1.5 text-sm">
                    <p className="text-slate-700">
                      <span className="text-slate-400 inline-block w-16">수령인</span>
                      {shipping.name}
                    </p>
                    <p className="text-slate-700">
                      <span className="text-slate-400 inline-block w-16">연락처</span>
                      {shipping.phone}
                    </p>
                    <p className="text-slate-700">
                      <span className="text-slate-400 inline-block w-16">주소</span>
                      ({shipping.zipcode}) {shipping.address} {shipping.addressDetail}
                    </p>
                    {shipping.memo && (
                      <p className="text-slate-700">
                        <span className="text-slate-400 inline-block w-16">메모</span>
                        {shipping.memo}
                      </p>
                    )}
                  </div>
                </div>

                {/* 결제 수단 확인 */}
                <div className="rounded-2xl border border-slate-200 bg-white p-6">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-sm font-bold text-slate-900">결제 수단</h3>
                    <button
                      onClick={() => setCurrentStep(2)}
                      className="text-xs font-semibold text-sky-600 hover:text-sky-700"
                    >
                      수정
                    </button>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-xl">
                      {PAYMENT_METHODS.find((m) => m.id === paymentMethod)?.icon}
                    </span>
                    <span className="text-sm font-semibold text-slate-700">
                      {PAYMENT_METHODS.find((m) => m.id === paymentMethod)?.label}
                    </span>
                    {paymentMethod === 'card' && (
                      <span className="text-xs text-slate-400 font-mono">•••• 4242</span>
                    )}
                  </div>
                </div>

                {/* 최종 금액 + 주문 버튼 */}
                <div className="rounded-2xl border-2 border-sky-200 bg-sky-50 p-6">
                  <div className="space-y-2 mb-4">
                    <div className="flex justify-between text-sm">
                      <span className="text-slate-600">상품 합계</span>
                      <span className="text-slate-700">₩{subtotal.toLocaleString()}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-slate-600">배송비</span>
                      <span className={shippingFee === 0 ? 'text-emerald-600 font-semibold' : 'text-slate-700'}>
                        {shippingFee === 0 ? '무료' : `₩${shippingFee.toLocaleString()}`}
                      </span>
                    </div>
                    <div className="border-t border-sky-200 pt-2 flex justify-between items-center">
                      <span className="text-base font-bold text-slate-900">총 결제금액</span>
                      <span className="text-2xl font-bold text-sky-600">
                        ₩{total.toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <button
                    onClick={handleOrder}
                    className="w-full py-4 rounded-xl bg-sky-600 text-white text-base font-bold hover:bg-sky-700 transition-colors shadow-lg shadow-sky-200"
                  >
                    주문하기
                  </button>
                  <p className="text-xs text-slate-400 text-center mt-3">
                    주문 완료 시 이용약관 및 개인정보 처리방침에 동의한 것으로 간주합니다.
                  </p>
                </div>
              </div>
            )}

            {/* 네비게이션 버튼 */}
            {currentStep < 3 && (
              <div className="flex items-center justify-between mt-6">
                {currentStep > 1 ? (
                  <button
                    onClick={goPrev}
                    className="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                  >
                    <ChevronLeft className="h-4 w-4" />
                    이전
                  </button>
                ) : (
                  <Link
                    href={`${localePrefix}/domains/store`}
                    className="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                  >
                    <ChevronLeft className="h-4 w-4" />
                    스토어로 돌아가기
                  </Link>
                )}
                <button
                  onClick={goNext}
                  className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-bold hover:bg-sky-700 transition-colors"
                >
                  다음
                  <ChevronRight className="h-4 w-4" />
                </button>
              </div>
            )}
            {currentStep === 3 && (
              <div className="mt-6">
                <button
                  onClick={goPrev}
                  className="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  <ChevronLeft className="h-4 w-4" />
                  이전
                </button>
              </div>
            )}
          </div>

          {/* 오른쪽 요약 패널 */}
          <div className="lg:col-span-1">
            <div className="sticky top-28">
              <OrderSummaryPanel />
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
