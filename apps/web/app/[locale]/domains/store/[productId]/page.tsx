'use client';

import React, { useState } from 'react';
import { Star, Minus, Plus, ArrowLeft, ShoppingCart, Truck, CheckCircle, ChevronRight } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { use } from 'react';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const PRODUCTS = [
  { id: '1', name: '혈당 카트리지 (10개입)', category: '카트리지', price: 29000, rating: 4.8, reviews: 1234, image: '🩸', badge: '인기', subscription: true, description: 'MPS 리더기와 호환되는 혈당 측정 전용 카트리지입니다. 미세혈액 샘플로 빠르고 정확한 혈당 수치를 제공합니다. 14핀 CSI v1.0 인터페이스 기반으로 안정적인 측정이 가능합니다.' },
  { id: '2', name: '콜레스테롤 카트리지 (10개입)', category: '카트리지', price: 35000, rating: 4.7, reviews: 892, image: '💉', badge: null, subscription: true, description: '총콜레스테롤, HDL, LDL을 동시에 측정할 수 있는 고정밀 카트리지입니다. 차동측정 원리를 적용하여 환경 노이즈를 제거하고 신뢰도 높은 결과를 제공합니다.' },
  { id: '3', name: '당화혈색소 카트리지 (5개입)', category: '카트리지', price: 45000, rating: 4.9, reviews: 567, image: '🔬', badge: '신규', subscription: true, description: 'HbA1c 수치를 가정에서 간편하게 확인할 수 있는 카트리지입니다. 3개월간의 혈당 관리 상태를 한눈에 파악하세요. 의료기관 수준의 정확도를 지원합니다.' },
  { id: '4', name: '종합검사 카트리지 세트', category: '카트리지', price: 89000, rating: 4.8, reviews: 345, image: '📦', badge: '추천', subscription: false, description: '혈당, 콜레스테롤, 당화혈색소를 포함한 종합 건강 체크 세트입니다. 개별 구매 대비 할인된 가격으로 전체적인 건강 상태를 모니터링할 수 있습니다.' },
  { id: '5', name: 'MPS 리더기 Pro', category: '리더기', price: 290000, rating: 4.9, reviews: 2341, image: '📱', badge: '베스트', subscription: false, description: 'MMUP 기반 다중측정원리 유니버설 POCT 리더기입니다. BLE 5.0 연결, NFC 카트리지 인식, 오프라인 측정을 지원하며 Flutter 앱과 실시간 연동됩니다.' },
  { id: '6', name: '채혈침 세트 (100개입)', category: '소모품', price: 12000, rating: 4.6, reviews: 3456, image: '💊', badge: null, subscription: false, description: '안전하고 위생적인 일회용 채혈침입니다. 미세 바늘로 통증을 최소화하며, 사용 후 자동 잠금 기능으로 안전하게 폐기할 수 있습니다.' },
  { id: '7', name: '알코올 솜 (200매)', category: '소모품', price: 8000, rating: 4.5, reviews: 5678, image: '🧴', badge: null, subscription: false, description: '측정 전 피부 소독을 위한 의료용 알코올 솜입니다. 개별 포장으로 위생적이며, 적절한 알코올 함량으로 빠르게 건조됩니다.' },
  { id: '8', name: '휴대용 파우치', category: '액세서리', price: 25000, rating: 4.7, reviews: 890, image: '👜', badge: null, subscription: false, description: 'MPS 리더기와 카트리지를 안전하게 보관하고 휴대할 수 있는 전용 파우치입니다. 충격 방지 내부 구조와 카트리지 수납 포켓이 포함되어 있습니다.' },
];

const MOCK_REVIEWS = [
  { id: 'r1', author: '김**', rating: 5, date: '2026-05-10', text: '정확도가 정말 뛰어납니다. 병원에서 측정한 수치와 거의 동일하게 나와서 놀랐습니다. 정기배송으로 이용 중인데 편리합니다.', verified: true },
  { id: 'r2', author: '이**', rating: 4, date: '2026-04-28', text: '사용법이 간단하고 결과가 빨리 나옵니다. 앱과 자동으로 연동되어 기록 관리가 편합니다. 가격이 조금만 더 저렴하면 좋겠습니다.', verified: true },
  { id: 'r3', author: '박**', rating: 5, date: '2026-04-15', text: '가족 모두 사용하고 있습니다. AI 코칭과 연결해서 건강 관리를 하니 훨씬 체계적으로 관리할 수 있어요. 강력 추천합니다.', verified: false },
];

function getRelatedProducts(currentId: string) {
  return PRODUCTS.filter(p => p.id !== currentId).slice(0, 3);
}

export default function ProductDetailPage({ params }: { params: Promise<{ productId: string }> }) {
  const { productId } = use(params);
  const localePrefix = useLocalePrefix();
  const [quantity, setQuantity] = useState(1);
  const [subscriptionEnabled, setSubscriptionEnabled] = useState(false);
  const [addedToCart, setAddedToCart] = useState(false);

  const product = PRODUCTS.find(p => p.id === productId);

  if (!product) {
    return (
      <main className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-6xl mb-4">🔍</p>
          <h2 className="text-xl font-bold text-slate-900 mb-2">상품을 찾을 수 없습니다</h2>
          <p className="text-sm text-slate-500 mb-6">요청하신 상품이 존재하지 않거나 삭제되었습니다.</p>
          <Link href={`${localePrefix}/domains/store`} className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-sky-600 text-white text-sm font-semibold hover:bg-sky-700 transition-colors">
            <ArrowLeft className="h-4 w-4" />
            스토어로 돌아가기
          </Link>
        </div>
      </main>
    );
  }

  const related = getRelatedProducts(product.id);

  const handleAddToCart = () => {
    setAddedToCart(true);
    setTimeout(() => setAddedToCart(false), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="상품 상세">
      <div className="max-w-6xl mx-auto px-6 py-8">
        {/* 뒤로가기 */}
        <Link href={`${localePrefix}/domains/store`} className="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-sky-600 transition-colors mb-6">
          <ArrowLeft className="h-4 w-4" />
          스토어로 돌아가기
        </Link>

        {/* 상품 메인 영역 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
          {/* 상품 이미지 */}
          <div className="rounded-2xl border border-slate-200 bg-white p-8 flex items-center justify-center min-h-[360px] relative">
            <span className="text-[120px] leading-none">{product.image}</span>
            {product.badge && (
              <span className={`absolute top-4 left-4 px-3 py-1 rounded-full text-xs font-bold text-white ${
                product.badge === '인기' ? 'bg-rose-500' : product.badge === '신규' ? 'bg-emerald-500' : product.badge === '추천' ? 'bg-sky-500' : 'bg-amber-500'
              }`}>{product.badge}</span>
            )}
          </div>

          {/* 상품 정보 */}
          <div className="flex flex-col">
            <span className="text-xs font-semibold text-sky-600 uppercase tracking-wider mb-2">{product.category}</span>
            <h1 className="text-2xl font-bold text-slate-900 mb-3">{product.name}</h1>

            <div className="flex items-center gap-2 mb-4">
              <div className="flex items-center gap-0.5">
                {[1, 2, 3, 4, 5].map(star => (
                  <Star key={star} className={`h-4 w-4 ${star <= Math.round(product.rating) ? 'text-amber-500 fill-amber-500' : 'text-slate-200'}`} />
                ))}
              </div>
              <span className="text-sm font-medium text-slate-700">{product.rating}</span>
              <span className="text-sm text-slate-400">({product.reviews.toLocaleString()}개 리뷰)</span>
            </div>

            <p className="text-slate-600 text-sm leading-relaxed mb-6">{product.description}</p>

            <div className="text-3xl font-bold text-slate-900 mb-6">
              ₩{product.price.toLocaleString()}
            </div>

            {/* 정기배송 옵션 */}
            {product.subscription && (
              <div className="rounded-xl border border-slate-200 bg-slate-50 p-4 mb-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Truck className="h-5 w-5 text-purple-600" />
                    <div>
                      <p className="text-sm font-semibold text-slate-900">정기배송으로 받기</p>
                      <p className="text-xs text-slate-500">매월 자동 배송 | 10% 할인 적용</p>
                    </div>
                  </div>
                  <button
                    onClick={() => setSubscriptionEnabled(!subscriptionEnabled)}
                    role="switch"
                    aria-checked={subscriptionEnabled}
                    aria-label="정기배송 설정"
                    className={`relative w-12 h-7 rounded-full transition-colors ${subscriptionEnabled ? 'bg-purple-600' : 'bg-slate-300'}`}
                  >
                    <span className={`absolute top-0.5 left-0.5 h-6 w-6 rounded-full bg-white shadow transition-transform ${subscriptionEnabled ? 'translate-x-5' : 'translate-x-0'}`} />
                  </button>
                </div>
                {subscriptionEnabled && (
                  <div className="mt-3 pt-3 border-t border-slate-200">
                    <p className="text-sm text-purple-700 font-semibold">
                      정기배송 가격: ₩{Math.round(product.price * 0.9).toLocaleString()} <span className="text-xs text-purple-500 font-normal">(10% 할인)</span>
                    </p>
                  </div>
                )}
              </div>
            )}

            {/* 수량 선택 + 장바구니 */}
            <div className="flex items-center gap-4 mt-auto">
              <div className="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-2 py-1.5">
                <button
                  onClick={() => setQuantity(q => Math.max(1, q - 1))}
                  aria-label="수량 감소"
                  className="h-8 w-8 rounded-lg hover:bg-slate-100 flex items-center justify-center transition-colors"
                >
                  <Minus className="h-4 w-4 text-slate-600" />
                </button>
                <span className="w-10 text-center text-sm font-bold text-slate-900">{quantity}</span>
                <button
                  onClick={() => setQuantity(q => q + 1)}
                  aria-label="수량 증가"
                  className="h-8 w-8 rounded-lg hover:bg-slate-100 flex items-center justify-center transition-colors"
                >
                  <Plus className="h-4 w-4 text-slate-600" />
                </button>
              </div>

              <button
                onClick={handleAddToCart}
                className={`flex-1 py-3.5 rounded-xl font-bold text-sm transition-colors flex items-center justify-center gap-2 ${
                  addedToCart
                    ? 'bg-emerald-600 text-white'
                    : 'bg-sky-600 text-white hover:bg-sky-700'
                }`}
              >
                {addedToCart ? (
                  <>
                    <CheckCircle className="h-5 w-5" />
                    담기 완료!
                  </>
                ) : (
                  <>
                    <ShoppingCart className="h-5 w-5" />
                    장바구니 담기
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        {/* 리뷰 섹션 */}
        <section aria-label="고객 리뷰" className="mb-12">
          <h2 className="text-lg font-bold text-slate-900 mb-4">고객 리뷰</h2>
          <div className="space-y-4">
            {MOCK_REVIEWS.map(review => (
              <div key={review.id} className="rounded-2xl border border-slate-200 bg-white p-5">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="flex items-center gap-0.5">
                      {[1, 2, 3, 4, 5].map(star => (
                        <Star key={star} className={`h-3.5 w-3.5 ${star <= review.rating ? 'text-amber-500 fill-amber-500' : 'text-slate-200'}`} />
                      ))}
                    </div>
                    <span className="text-sm font-semibold text-slate-900">{review.author}</span>
                    {review.verified && (
                      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 text-[10px] font-semibold border border-emerald-200">
                        <CheckCircle className="h-3 w-3" />
                        구매 인증
                      </span>
                    )}
                  </div>
                  <span className="text-xs text-slate-400">{review.date}</span>
                </div>
                <p className="text-sm text-slate-600 leading-relaxed">{review.text}</p>
              </div>
            ))}
          </div>
        </section>

        {/* 관련 상품 */}
        <section aria-label="관련 상품" className="mb-12">
          <h2 className="text-lg font-bold text-slate-900 mb-4">관련 상품</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            {related.map(item => (
              <Link
                key={item.id}
                href={`${localePrefix}/domains/store/${item.id}`}
                className="rounded-2xl border border-slate-200 bg-white overflow-hidden hover:shadow-md transition-shadow group"
              >
                <div className="h-32 bg-slate-50 flex items-center justify-center">
                  <span className="text-4xl">{item.image}</span>
                </div>
                <div className="p-4">
                  <h3 className="text-sm font-bold text-slate-900 mb-1 line-clamp-1 group-hover:text-sky-600 transition-colors">{item.name}</h3>
                  <div className="flex items-center gap-1 mb-2">
                    <Star className="h-3 w-3 text-amber-500" />
                    <span className="text-xs text-slate-500">{item.rating}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-base font-bold text-slate-900">₩{item.price.toLocaleString()}</span>
                    <ChevronRight className="h-4 w-4 text-slate-400 group-hover:text-sky-600 transition-colors" />
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </section>
      </div>
    </main>
  );
}
