'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { ShoppingCart, Package, Star, Truck, Filter, Plus, Minus, X, CheckCircle } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const PRODUCTS = [
  { id: '1', name: '혈당 카트리지 (10개입)', category: '카트리지', price: 29000, rating: 4.8, reviews: 1234, image: '🩸', badge: '인기', subscription: true },
  { id: '2', name: '콜레스테롤 카트리지 (10개입)', category: '카트리지', price: 35000, rating: 4.7, reviews: 892, image: '💉', badge: null, subscription: true },
  { id: '3', name: '당화혈색소 카트리지 (5개입)', category: '카트리지', price: 45000, rating: 4.9, reviews: 567, image: '🔬', badge: '신규', subscription: true },
  { id: '4', name: '종합검사 카트리지 세트', category: '카트리지', price: 89000, rating: 4.8, reviews: 345, image: '📦', badge: '추천', subscription: false },
  { id: '5', name: 'MPS 리더기 Pro', category: '리더기', price: 290000, rating: 4.9, reviews: 2341, image: '📱', badge: '베스트', subscription: false },
  { id: '6', name: '채혈침 세트 (100개입)', category: '소모품', price: 12000, rating: 4.6, reviews: 3456, image: '💊', badge: null, subscription: false },
  { id: '7', name: '알코올 솜 (200매)', category: '소모품', price: 8000, rating: 4.5, reviews: 5678, image: '🧴', badge: null, subscription: false },
  { id: '8', name: '휴대용 파우치', category: '액세서리', price: 25000, rating: 4.7, reviews: 890, image: '👜', badge: null, subscription: false },
];

const CATEGORIES = ['전체', '카트리지', '리더기', '소모품', '액세서리'];

type CartItem = { product: typeof PRODUCTS[0]; quantity: number };

export default function StorePage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [category, setCategory] = useState('전체');
  const [cart, setCart] = useState<CartItem[]>([]);
  const [cartOpen, setCartOpen] = useState(false);
  const [orderComplete, setOrderComplete] = useState(false);

  const filtered = category === '전체' ? PRODUCTS : PRODUCTS.filter(p => p.category === category);

  const addToCart = (product: typeof PRODUCTS[0]) => {
    setCart(prev => {
      const existing = prev.find(item => item.product.id === product.id);
      if (existing) return prev.map(item => item.product.id === product.id ? { ...item, quantity: item.quantity + 1 } : item);
      return [...prev, { product, quantity: 1 }];
    });
  };

  const updateQuantity = (productId: string, delta: number) => {
    setCart(prev => prev.map(item => {
      if (item.product.id !== productId) return item;
      const newQty = item.quantity + delta;
      return newQty <= 0 ? item : { ...item, quantity: newQty };
    }).filter(item => item.quantity > 0));
  };

  const removeFromCart = (productId: string) => {
    setCart(prev => prev.filter(item => item.product.id !== productId));
  };

  const totalPrice = cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0);
  const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);

  const handleOrder = () => {
    setOrderComplete(true);
    setCart([]);
    setTimeout(() => { setOrderComplete(false); setCartOpen(false); }, 3000);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="쇼핑몰">
      <DomainHeader
        title="MPS 스토어"
        subtitle="카트리지와 건강 관리 용품을 구매하세요"
        icon={ShoppingCart}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-6xl mx-auto px-6 py-8">
        {/* 카테고리 필터 + 장바구니 */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <Filter className="h-4 w-4 text-slate-400" />
            {CATEGORIES.map(c => (
              <button key={c} onClick={() => setCategory(c)}
                className={`px-4 py-2 rounded-xl text-sm font-semibold transition-colors ${
                  category === c ? 'bg-sky-600 text-white' : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-50'
                }`}>
                {c}
              </button>
            ))}
          </div>
          <button onClick={() => setCartOpen(true)} className="relative flex items-center gap-2 px-4 py-2 rounded-xl bg-slate-900 text-white text-sm font-semibold hover:bg-slate-800">
            <ShoppingCart className="h-4 w-4" />
            장바구니
            {totalItems > 0 && (
              <span className="absolute -top-1.5 -right-1.5 h-5 w-5 rounded-full bg-sky-500 text-[10px] font-bold flex items-center justify-center">{totalItems}</span>
            )}
          </button>
        </div>

        {/* 상품 그리드 */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {filtered.map(product => (
            <div key={product.id} className="rounded-2xl border border-slate-200 bg-white overflow-hidden hover:shadow-md transition-shadow group">
              <div className="relative h-40 bg-slate-50 flex items-center justify-center">
                <span className="text-5xl">{product.image}</span>
                {product.badge && (
                  <span className={`absolute top-3 left-3 px-2 py-0.5 rounded-full text-[10px] font-bold text-white ${
                    product.badge === '인기' ? 'bg-rose-500' : product.badge === '신규' ? 'bg-emerald-500' : product.badge === '추천' ? 'bg-sky-500' : 'bg-amber-500'
                  }`}>{product.badge}</span>
                )}
                {product.subscription && (
                  <span className="absolute top-3 right-3 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-purple-50 text-purple-600 border border-purple-200">
                    <Truck className="h-3 w-3 inline mr-0.5" />정기배송
                  </span>
                )}
              </div>
              <div className="p-4">
                <h3 className="text-sm font-bold text-slate-900 mb-1 line-clamp-2">{product.name}</h3>
                <div className="flex items-center gap-1 mb-2">
                  <Star className="h-3 w-3 text-amber-500" />
                  <span className="text-xs text-slate-500">{product.rating} ({product.reviews.toLocaleString()})</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-lg font-bold text-slate-900">₩{product.price.toLocaleString()}</span>
                  <button onClick={() => addToCart(product)}
                    className="h-9 w-9 rounded-xl bg-sky-600 text-white flex items-center justify-center hover:bg-sky-700 transition-colors">
                    <Plus className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* 관련 도메인 */}
        <section aria-label="관련 도메인" className="mt-8">
          <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {[
              { name: '건강 측정', path: '/domains/measure', desc: '카트리지로 측정' },
              { name: '보상', path: '/domains/reward', desc: '포인트 사용' },
              { name: '화상 진료', path: '/domains/telemedicine', desc: '전문의 상담' },
            ].map(d => (
              <Link key={d.path} href={`${localePrefix}${d.path}`}
                className="flex items-center gap-3 p-3 rounded-xl border border-slate-200 bg-white hover:border-sky-300 hover:shadow-sm transition-all group">
                <div className="h-8 w-8 rounded-lg bg-sky-50 flex items-center justify-center text-sky-600 group-hover:bg-sky-100">
                  <span className="text-sm">→</span>
                </div>
                <div>
                  <span className="text-sm font-semibold text-slate-700 group-hover:text-sky-600">{d.name}</span>
                  <p className="text-[11px] text-slate-400">{d.desc}</p>
                </div>
              </Link>
            ))}
          </div>
        </section>
      </div>

      {/* 장바구니 사이드 패널 */}
      {cartOpen && (
        <div className="fixed inset-0 z-50" onClick={() => setCartOpen(false)}>
          <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" />
          <div className="absolute right-0 top-0 h-full w-96 max-w-[90vw] bg-white shadow-2xl p-6 overflow-y-auto"
            onClick={e => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-bold text-slate-900">장바구니 ({totalItems})</h2>
              <button onClick={() => setCartOpen(false)} className="p-1 rounded-lg hover:bg-slate-100"><X className="h-5 w-5" /></button>
            </div>

            {orderComplete ? (
              <div className="text-center py-12">
                <CheckCircle className="h-16 w-16 text-emerald-500 mx-auto mb-4" />
                <h3 className="text-lg font-bold text-slate-900 mb-2">주문이 완료되었습니다</h3>
                <p className="text-sm text-slate-500">배송 준비가 시작됩니다</p>
              </div>
            ) : cart.length === 0 ? (
              <div className="text-center py-12">
                <ShoppingCart className="h-12 w-12 text-slate-300 mx-auto mb-3" />
                <p className="text-sm text-slate-500">장바구니가 비어 있습니다</p>
              </div>
            ) : (
              <>
                <div className="space-y-3 mb-6">
                  {cart.map(item => (
                    <div key={item.product.id} className="flex items-center gap-3 p-3 rounded-xl bg-slate-50">
                      <span className="text-2xl">{item.product.image}</span>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-slate-900 truncate">{item.product.name}</p>
                        <p className="text-sm font-bold text-slate-700">₩{(item.product.price * item.quantity).toLocaleString()}</p>
                      </div>
                      <div className="flex items-center gap-1">
                        <button onClick={() => updateQuantity(item.product.id, -1)} className="h-7 w-7 rounded-lg bg-white border flex items-center justify-center"><Minus className="h-3 w-3" /></button>
                        <span className="text-sm font-bold w-6 text-center">{item.quantity}</span>
                        <button onClick={() => updateQuantity(item.product.id, 1)} className="h-7 w-7 rounded-lg bg-white border flex items-center justify-center"><Plus className="h-3 w-3" /></button>
                      </div>
                      <button onClick={() => removeFromCart(item.product.id)} className="text-slate-400 hover:text-red-500"><X className="h-4 w-4" /></button>
                    </div>
                  ))}
                </div>
                <div className="border-t pt-4">
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-sm text-slate-500">합계</span>
                    <span className="text-xl font-bold text-slate-900">₩{totalPrice.toLocaleString()}</span>
                  </div>
                  <button onClick={handleOrder} className="w-full py-3.5 rounded-xl bg-sky-600 text-white font-bold hover:bg-sky-700 transition-colors">
                    주문하기
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </main>
  );
}
