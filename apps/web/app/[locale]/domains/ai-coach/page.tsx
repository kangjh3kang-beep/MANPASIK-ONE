'use client';

import React, { useState, useCallback } from 'react';
import { DualPersonaChat, ChatMessage, ChatPersona, DomainHeader, NextActionGroup } from '@mmup/ui';
import { Brain, Video, ShoppingCart } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

const INITIAL_MESSAGES: ChatMessage[] = [
  {
    id: '1',
    persona: 'assistant',
    role: 'ai',
    content: '안녕하세요! 건강 생활 도우미입니다. 오늘 건강에 대해 궁금한 점이 있으시면 편하게 물어보세요.',
    timestamp: new Date().toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' }),
  },
];

const DEMO_RESPONSES: Record<string, { content: string; citations?: any[]; isMedicalClaim?: boolean; uncertainty?: string }> = {
  assistant: {
    content: '좋은 질문이에요! 규칙적인 식사와 30분 이상의 가벼운 운동을 추천드려요. 충분한 수면도 건강 유지에 중요합니다.',
  },
  doctor: {
    content: '최근 측정 데이터를 기반으로 분석하면, 현재 건강 지표는 정상 범위 내에 있습니다. 다만 정기적인 모니터링을 권장합니다.',
    citations: [
      { id: 'c1', source: '최근 측정 데이터', confidence: 0.94 },
    ],
    isMedicalClaim: true,
    uncertainty: '이 분석은 AI 참고 정보이며, 정확한 진단은 전문의 상담이 필요합니다',
  },
};

export default function AiCoachPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();
  const [persona, setPersona] = useState<ChatPersona>('assistant');
  const [messages, setMessages] = useState<ChatMessage[]>(INITIAL_MESSAGES);
  const [isStreaming, setIsStreaming] = useState(false);

  const handleSend = useCallback((text: string) => {
    const userMsg: ChatMessage = {
      id: `u-${Date.now()}`,
      persona,
      role: 'user',
      content: text,
      timestamp: new Date().toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' }),
    };
    setMessages(prev => [...prev, userMsg]);
    setIsStreaming(true);

    setTimeout(() => {
      const resp = DEMO_RESPONSES[persona];
      const aiMsg: ChatMessage = {
        id: `a-${Date.now()}`,
        persona,
        role: 'ai',
        content: resp.content,
        citations: resp.citations,
        isMedicalClaim: resp.isMedicalClaim,
        uncertainty: resp.uncertainty,
        timestamp: new Date().toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' }),
      };
      setMessages(prev => [...prev, aiMsg]);
      setIsStreaming(false);
    }, 1500);
  }, [persona]);

  return (
    <main className="min-h-screen bg-slate-50 flex flex-col" aria-label="AI 건강 코치">
      <DomainHeader
        title="AI 건강 코치"
        subtitle="생활비서와 AI주치의가 건강 궁금증에 답변합니다"
        icon={Brain}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'patient' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="flex-1 max-w-4xl mx-auto w-full px-4 py-6">
        <div className="rounded-2xl border border-slate-200 bg-white shadow-sm overflow-hidden h-[600px] flex flex-col" style={{ borderColor: 'var(--mps-border)' }}>
          <DualPersonaChat
            persona={persona}
            onPersonaChange={setPersona}
            messages={messages}
            onSend={handleSend}
            isStreaming={isStreaming}
          />
        </div>

        <p className="text-center text-xs text-slate-400 mt-4">
          AI 코치의 답변은 참고 정보이며, 의학적 진단이나 처방을 대체하지 않습니다.
        </p>

        {/* 다음 단계 추천 */}
        <div className="mt-6">
          <NextActionGroup
            title="다음 단계"
            actions={[
              { id: 'telemedicine', title: '화상 진료 예약', description: '전문의와 직접 상담하세요', icon: <Video className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/telemedicine`, priority: 'high' },
              { id: 'store', title: '건강용품 구매', description: '추천 카트리지와 건강식품을 확인하세요', icon: <ShoppingCart className="h-4 w-4" />, onClick: () => window.location.href = `${localePrefix}/domains/store` },
            ]}
          />
        </div>
      </div>

      {/* 관련 도메인 */}
      <section aria-label="관련 도메인" className="max-w-4xl mx-auto w-full px-4 pb-8">
        <h3 className="text-sm font-bold text-slate-500 mb-3">관련 도메인</h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[
            { name: '화상 진료', path: '/domains/telemedicine', desc: '전문의 상담 연결' },
            { name: '위험 예측', path: '/domains/predictor', desc: '건강 위험도 분석' },
            { name: '건강 데이터', path: '/domains/clinical', desc: '측정 결과 확인' },
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
    </main>
  );
}
