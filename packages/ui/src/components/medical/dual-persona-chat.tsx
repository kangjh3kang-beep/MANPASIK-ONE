'use client';

import React, { useState } from 'react';

export type ChatPersona = 'assistant' | 'doctor';

export interface ChatMessage {
  id: string;
  persona: ChatPersona;
  role: 'user' | 'ai';
  content: string;
  citations?: { id: string; source: string; confidence?: number }[];
  isMedicalClaim?: boolean;
  uncertainty?: string;
  timestamp: string;
}

export interface DualPersonaChatProps {
  persona: ChatPersona;
  onPersonaChange: (persona: ChatPersona) => void;
  messages: ChatMessage[];
  onSend: (text: string) => void;
  isStreaming?: boolean;
  offline?: boolean;
}

const PERSONA_CONFIG = {
  assistant: {
    label: '생활비서',
    icon: '🏠',
    accent: 'bg-sky-500',
    accentLight: 'bg-sky-50 border-sky-200',
    accentText: 'text-sky-700',
    description: '따뜻한 건강 생활 도우미',
  },
  doctor: {
    label: 'AI주치의',
    icon: '🩺',
    accent: 'bg-emerald-600',
    accentLight: 'bg-emerald-50 border-emerald-200',
    accentText: 'text-emerald-700',
    description: '전문의 기반 의학 상담',
  },
};

export function DualPersonaChat({
  persona, onPersonaChange, messages, onSend, isStreaming, offline,
}: DualPersonaChatProps) {
  const [input, setInput] = useState('');
  const config = PERSONA_CONFIG[persona];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim() && !isStreaming) {
      onSend(input.trim());
      setInput('');
    }
  };

  return (
    <div className="flex flex-col h-full" aria-label={`${config.label} 채팅`}>
      {/* Persona toggle — label + icon 구분 (색만 X) */}
      <div className="flex items-center gap-2 p-3 border-b border-slate-200 bg-white">
        {(['assistant', 'doctor'] as ChatPersona[]).map((p) => {
          const cfg = PERSONA_CONFIG[p];
          const isActive = persona === p;
          return (
            <button
              key={p}
              onClick={() => onPersonaChange(p)}
              className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all ${
                isActive
                  ? `${cfg.accentLight} ${cfg.accentText} border`
                  : 'text-slate-400 hover:text-slate-600'
              }`}
              aria-label={`${cfg.label} 모드로 전환`}
              aria-pressed={isActive}
            >
              <span aria-hidden="true">{cfg.icon}</span>
              <span>{cfg.label}</span>
            </button>
          );
        })}
        {offline && (
          <span className="ml-auto text-[10px] px-2 py-0.5 rounded-full bg-slate-100 text-slate-500">
            오프라인
          </span>
        )}
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4" role="log" aria-live="polite">
        {messages.map((msg) => {
          const msgConfig = PERSONA_CONFIG[msg.persona];
          const isUser = msg.role === 'user';

          return (
            <div key={msg.id} className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}>
              <div className={`max-w-[80%] ${isUser ? 'order-1' : ''}`}>
                {/* AI persona label — icon + text (색 단독 구분 금지) */}
                {!isUser && (
                  <div className="flex items-center gap-1 mb-1">
                    <span className="text-xs" aria-hidden="true">{msgConfig.icon}</span>
                    <span className={`text-xs font-medium ${msgConfig.accentText}`}>{msgConfig.label}</span>
                  </div>
                )}

                {/* Bubble */}
                <div
                  className={`rounded-2xl px-4 py-3 text-sm leading-relaxed ${
                    isUser
                      ? 'bg-slate-900 text-white rounded-br-md'
                      : `${msgConfig.accentLight} border rounded-bl-md`
                  }`}
                >
                  {msg.content}

                  {/* Medical claim warning */}
                  {msg.isMedicalClaim && (
                    <div className="mt-2 flex items-center gap-1 text-[10px] text-amber-600" role="note">
                      <span aria-hidden="true">△</span>
                      <span>의학적 주장 — {msg.uncertainty || '전문의 상담을 권장합니다'}</span>
                    </div>
                  )}
                </div>

                {/* Source citations with confidence */}
                {msg.citations && msg.citations.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1">
                    {msg.citations.map((cite) => (
                      <span
                        key={cite.id}
                        className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-slate-100 text-[10px] text-slate-500"
                        aria-label={`출처: ${cite.source}${cite.confidence ? `, 신뢰도 ${Math.round(cite.confidence * 100)}%` : ''}`}
                      >
                        📎 {cite.source}
                        {cite.confidence != null && (
                          <span className="text-slate-400">{Math.round(cite.confidence * 100)}%</span>
                        )}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          );
        })}

        {isStreaming && (
          <div className="flex justify-start">
            <div className={`${config.accentLight} border rounded-2xl rounded-bl-md px-4 py-3`}>
              <span className="animate-pulse text-sm text-slate-400">입력 중...</span>
            </div>
          </div>
        )}
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="p-3 border-t border-slate-200 bg-white">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={persona === 'assistant' ? '건강 궁금한 점을 물어보세요' : '증상이나 검사 결과를 알려주세요'}
            className="flex-1 px-4 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500/30"
            disabled={isStreaming}
            aria-label="메시지 입력"
          />
          <button
            type="submit"
            disabled={!input.trim() || isStreaming}
            className={`px-4 py-2.5 rounded-xl text-sm font-medium text-white ${config.accent} disabled:opacity-40`}
            aria-label="전송"
          >
            전송
          </button>
        </div>
      </form>
    </div>
  );
}
