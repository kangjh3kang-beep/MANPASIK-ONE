'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Zap, CheckCircle, AlertTriangle, Clock, Globe, Play, Settings } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

interface Webhook {
  id: string;
  name: string;
  url: string;
  events: string[];
  active: boolean;
  lastDelivery: { status: '성공' | '실패'; time: string };
}

interface DeliveryLog {
  id: string;
  webhookName: string;
  event: string;
  status: '성공' | '실패';
  responseCode: number;
  retryCount: number;
  timestamp: string;
}

const EVENT_TYPES = [
  { id: 'measurement.completed', label: '측정 완료' },
  { id: 'order.updated', label: '주문 상태 변경' },
  { id: 'prescription.created', label: '처방전 생성' },
  { id: 'alert.triggered', label: '알림 발생' },
];

const MOCK_WEBHOOKS: Webhook[] = [
  {
    id: 'wh-1',
    name: '측정 완료 알림',
    url: 'https://example.com/webhooks/measurement',
    events: ['measurement.completed', 'alert.triggered'],
    active: true,
    lastDelivery: { status: '성공', time: '2026-05-24 09:15:32' },
  },
  {
    id: 'wh-2',
    name: '주문 상태 변경',
    url: 'https://shop.example.com/hooks/order',
    events: ['order.updated'],
    active: false,
    lastDelivery: { status: '실패', time: '2026-05-23 18:42:10' },
  },
];

const MOCK_LOGS: DeliveryLog[] = [
  { id: 'dl-1', webhookName: '측정 완료 알림', event: 'measurement.completed', status: '성공', responseCode: 200, retryCount: 0, timestamp: '2026-05-24 09:15:32' },
  { id: 'dl-2', webhookName: '측정 완료 알림', event: 'alert.triggered', status: '성공', responseCode: 200, retryCount: 0, timestamp: '2026-05-24 08:30:11' },
  { id: 'dl-3', webhookName: '주문 상태 변경', event: 'order.updated', status: '실패', responseCode: 503, retryCount: 3, timestamp: '2026-05-23 18:42:10' },
  { id: 'dl-4', webhookName: '측정 완료 알림', event: 'measurement.completed', status: '성공', responseCode: 200, retryCount: 0, timestamp: '2026-05-23 14:22:05' },
  { id: 'dl-5', webhookName: '주문 상태 변경', event: 'order.updated', status: '성공', responseCode: 200, retryCount: 1, timestamp: '2026-05-22 11:08:44' },
];

export default function WebhooksPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [webhooks, setWebhooks] = useState(MOCK_WEBHOOKS);
  const [showAddForm, setShowAddForm] = useState(false);
  const [newUrl, setNewUrl] = useState('');
  const [newEvents, setNewEvents] = useState<string[]>([]);
  const [testingId, setTestingId] = useState<string | null>(null);

  const toggleWebhook = (id: string) => {
    setWebhooks(prev => prev.map(wh => wh.id === id ? { ...wh, active: !wh.active } : wh));
  };

  const handleTest = (id: string) => {
    setTestingId(id);
    setTimeout(() => setTestingId(null), 1500);
  };

  const handleEventToggle = (eventId: string) => {
    setNewEvents(prev => prev.includes(eventId) ? prev.filter(e => e !== eventId) : [...prev, eventId]);
  };

  return (
    <main className="min-h-screen bg-slate-50" aria-label="웹훅 관리">
      <DomainHeader
        title="웹훅 관리"
        subtitle="이벤트 기반 알림을 설정하고 전송 이력을 확인하세요"
        icon={Zap}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'dev' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8 space-y-6">
        {/* 구성된 웹훅 */}
        <section aria-label="웹훅 목록">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-bold text-slate-900">구성된 웹훅</h2>
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="px-4 py-2 rounded-xl bg-sky-600 text-white text-xs font-bold hover:bg-sky-700 transition-colors"
            >
              + 새 웹훅 추가
            </button>
          </div>

          <div className="space-y-4">
            {webhooks.map(wh => (
              <div key={wh.id} className="rounded-2xl border border-slate-200 bg-white p-5">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className={`h-8 w-8 rounded-lg flex items-center justify-center ${wh.active ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'}`}>
                      <Zap className="h-4 w-4" />
                    </div>
                    <div>
                      <h3 className="text-sm font-bold text-slate-900">{wh.name}</h3>
                      <code className="text-xs font-mono text-slate-500">{wh.url}</code>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => handleTest(wh.id)}
                      className="px-3 py-1.5 rounded-lg border border-slate-200 text-xs font-semibold text-slate-600 hover:bg-slate-50 transition-colors flex items-center gap-1"
                      aria-label={`${wh.name} 테스트 전송`}
                    >
                      {testingId === wh.id ? <CheckCircle className="h-3.5 w-3.5 text-emerald-500" /> : <Play className="h-3.5 w-3.5" />}
                      {testingId === wh.id ? '전송 완료' : '테스트 전송'}
                    </button>
                    <button
                      onClick={() => toggleWebhook(wh.id)}
                      className={`relative w-11 h-6 rounded-full transition-colors ${wh.active ? 'bg-emerald-500' : 'bg-slate-300'}`}
                      aria-label={`${wh.name} ${wh.active ? '비활성화' : '활성화'}`}
                    >
                      <span className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform ${wh.active ? 'left-[22px]' : 'left-0.5'}`} />
                    </button>
                  </div>
                </div>

                <div className="flex items-center gap-2 mb-3">
                  <span className="text-[10px] font-bold text-slate-400">구독 이벤트:</span>
                  {wh.events.map(ev => (
                    <span key={ev} className="px-2 py-0.5 rounded-full bg-sky-50 text-sky-700 text-[10px] font-semibold border border-sky-100">
                      {EVENT_TYPES.find(e => e.id === ev)?.label || ev}
                    </span>
                  ))}
                </div>

                <div className="flex items-center gap-2 text-[11px] text-slate-400">
                  <Clock className="h-3 w-3" />
                  <span>마지막 전송: {wh.lastDelivery.time}</span>
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold ${wh.lastDelivery.status === '성공' ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600'}`}>
                    {wh.lastDelivery.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* 새 웹훅 추가 폼 */}
        {showAddForm && (
          <section aria-label="새 웹훅 추가" className="rounded-2xl border border-sky-200 bg-white p-5">
            <h3 className="text-sm font-bold text-slate-900 mb-4 flex items-center gap-2">
              <Settings className="h-4 w-4 text-sky-600" /> 새 웹훅 설정
            </h3>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-1">Webhook URL</label>
                <div className="flex items-center gap-2">
                  <Globe className="h-4 w-4 text-slate-400" />
                  <input
                    type="url"
                    value={newUrl}
                    onChange={e => setNewUrl(e.target.value)}
                    placeholder="https://example.com/webhooks"
                    className="flex-1 rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm font-mono text-slate-700 focus:outline-none focus:ring-2 focus:ring-sky-500"
                    aria-label="웹훅 URL 입력"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-bold text-slate-500 mb-2">구독 이벤트</label>
                <div className="grid grid-cols-2 gap-2">
                  {EVENT_TYPES.map(evt => (
                    <label key={evt.id} className={`flex items-center gap-2 rounded-xl border p-3 cursor-pointer transition-colors ${newEvents.includes(evt.id) ? 'border-sky-300 bg-sky-50' : 'border-slate-200 bg-slate-50'}`}>
                      <input
                        type="checkbox"
                        checked={newEvents.includes(evt.id)}
                        onChange={() => handleEventToggle(evt.id)}
                        className="rounded border-slate-300 text-sky-600 focus:ring-sky-500"
                        aria-label={evt.label}
                      />
                      <span className="text-xs font-medium text-slate-700">{evt.label}</span>
                      <code className="text-[10px] text-slate-400 font-mono ml-auto">{evt.id}</code>
                    </label>
                  ))}
                </div>
              </div>
              <div className="flex justify-end gap-2">
                <button onClick={() => setShowAddForm(false)} className="px-4 py-2 rounded-xl border border-slate-200 text-xs font-semibold text-slate-600 hover:bg-slate-50">취소</button>
                <button className="px-4 py-2 rounded-xl bg-sky-600 text-white text-xs font-bold hover:bg-sky-700 transition-colors">저장</button>
              </div>
            </div>
          </section>
        )}

        {/* 전송 이력 */}
        <section aria-label="전송 이력">
          <h2 className="text-sm font-bold text-slate-900 mb-4">전송 이력</h2>
          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b border-slate-100 bg-slate-50">
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">시간</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">웹훅</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">이벤트</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">상태</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">응답 코드</th>
                  <th className="px-5 py-3 text-[11px] font-bold text-slate-500">재시도</th>
                </tr>
              </thead>
              <tbody>
                {MOCK_LOGS.map(log => (
                  <tr key={log.id} className="border-b border-slate-50 hover:bg-slate-50 transition-colors">
                    <td className="px-5 py-3 text-xs text-slate-600 font-mono">{log.timestamp}</td>
                    <td className="px-5 py-3 text-xs font-medium text-slate-700">{log.webhookName}</td>
                    <td className="px-5 py-3">
                      <span className="px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 text-[10px] font-mono">{log.event}</span>
                    </td>
                    <td className="px-5 py-3">
                      <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold ${log.status === '성공' ? 'bg-emerald-50 text-emerald-600' : 'bg-red-50 text-red-600'}`}>
                        {log.status === '성공' ? <CheckCircle className="h-3 w-3" /> : <AlertTriangle className="h-3 w-3" />}
                        {log.status}
                      </span>
                    </td>
                    <td className="px-5 py-3">
                      <span className={`text-xs font-mono font-bold ${log.responseCode < 300 ? 'text-emerald-600' : 'text-red-600'}`}>{log.responseCode}</span>
                    </td>
                    <td className="px-5 py-3 text-xs text-slate-500">{log.retryCount}회</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </main>
  );
}
