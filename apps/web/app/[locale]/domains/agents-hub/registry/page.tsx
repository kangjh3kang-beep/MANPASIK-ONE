'use client';

import React, { useState } from 'react';
import { DomainHeader } from '@mmup/ui';
import { Database, CheckCircle, Clock, AlertTriangle, TrendingUp, Code, FileText } from 'lucide-react';
import { useSession, signOut } from 'next-auth/react';
import { usePathname } from 'next/navigation';

function useLocalePrefix(): string {
  const pathname = usePathname();
  const match = pathname.match(/^\/(ko|en|ja|zh)/);
  return match ? `/${match[1]}` : '/ko';
}

type DeployStatus = '배포됨' | '스테이징' | '비활성';
type Framework = 'PyTorch' | 'XGBoost' | 'ONNX';

interface ModelVersion {
  version: string;
  date: string;
  accuracy: number;
}

interface RegisteredModel {
  id: string;
  name: string;
  version: string;
  framework: Framework;
  accuracy: number;
  size: string;
  deployStatus: DeployStatus;
  description: string;
  trainingDataset: string;
  inputSchema: string;
  outputSchema: string;
  versions: ModelVersion[];
}

const DEPLOY_STYLES: Record<DeployStatus, string> = {
  '배포됨': 'bg-emerald-50 text-emerald-700 border-emerald-200',
  '스테이징': 'bg-amber-50 text-amber-700 border-amber-200',
  '비활성': 'bg-slate-100 text-slate-500 border-slate-200',
};

const FRAMEWORK_STYLES: Record<Framework, string> = {
  'PyTorch': 'bg-orange-50 text-orange-700',
  'XGBoost': 'bg-sky-50 text-sky-700',
  'ONNX': 'bg-violet-50 text-violet-700',
};

const MODELS: RegisteredModel[] = [
  {
    id: 'mdl-1', name: '혈당 예측', version: 'v3.2', framework: 'PyTorch',
    accuracy: 96.8, size: '124MB', deployStatus: '배포됨',
    description: 'LSTM 기반 시계열 혈당 예측 모델',
    trainingDataset: '환자 12,400명 / 측정 데이터 1,200만건',
    inputSchema: '{ patient_id, features[88], timestamp }',
    outputSchema: '{ prediction, confidence, uncertainty }',
    versions: [
      { version: 'v3.2', date: '2026-05-10', accuracy: 96.8 },
      { version: 'v3.1', date: '2026-04-01', accuracy: 95.2 },
      { version: 'v3.0', date: '2026-02-15', accuracy: 93.7 },
    ],
  },
  {
    id: 'mdl-2', name: '심전도 분석', version: 'v2.1', framework: 'PyTorch',
    accuracy: 94.5, size: '256MB', deployStatus: '배포됨',
    description: 'CNN-Transformer 기반 12리드 심전도 분류 모델',
    trainingDataset: '심전도 기록 85,000건 / PTB-XL 포함',
    inputSchema: '{ ecg_signal[12][5000], sampling_rate }',
    outputSchema: '{ classification, probabilities[], anomaly_score }',
    versions: [
      { version: 'v2.1', date: '2026-04-20', accuracy: 94.5 },
      { version: 'v2.0', date: '2026-03-05', accuracy: 92.1 },
      { version: 'v1.9', date: '2026-01-18', accuracy: 89.8 },
    ],
  },
  {
    id: 'mdl-3', name: '영양 분석', version: 'v1.8', framework: 'XGBoost',
    accuracy: 91.2, size: '42MB', deployStatus: '스테이징',
    description: 'XGBoost 앙상블 기반 영양소 섭취 분석 및 권장량 계산',
    trainingDataset: '식이 기록 320,000건 / 영양성분 DB 연동',
    inputSchema: '{ food_items[], portions[], user_profile }',
    outputSchema: '{ nutrients{}, recommendations[], score }',
    versions: [
      { version: 'v1.8', date: '2026-05-18', accuracy: 91.2 },
      { version: 'v1.7', date: '2026-04-10', accuracy: 89.5 },
      { version: 'v1.6', date: '2026-03-01', accuracy: 87.3 },
    ],
  },
  {
    id: 'mdl-4', name: '위험 예측', version: 'v4.0', framework: 'PyTorch',
    accuracy: 97.1, size: '198MB', deployStatus: '배포됨',
    description: 'GNN 기반 다차원 건강 위험도 예측 모델',
    trainingDataset: '환자 28,000명 / 종합검진 데이터 5년분',
    inputSchema: '{ patient_id, measurements[448], demographics }',
    outputSchema: '{ risk_scores{}, alerts[], timeline }',
    versions: [
      { version: 'v4.0', date: '2026-05-05', accuracy: 97.1 },
      { version: 'v3.9', date: '2026-03-20', accuracy: 96.3 },
      { version: 'v3.8', date: '2026-02-08', accuracy: 95.0 },
    ],
  },
  {
    id: 'mdl-5', name: 'OCR 처리', version: 'v2.5', framework: 'ONNX',
    accuracy: 88.4, size: '78MB', deployStatus: '비활성',
    description: 'YOLOv8 + CRNN 기반 처방전/검사지 문자 인식',
    trainingDataset: '처방전 이미지 150,000건 / 한글 OCR 코퍼스',
    inputSchema: '{ image_base64, document_type }',
    outputSchema: '{ text_blocks[], structured_data{}, confidence }',
    versions: [
      { version: 'v2.5', date: '2026-04-25', accuracy: 88.4 },
      { version: 'v2.4', date: '2026-03-12', accuracy: 86.1 },
      { version: 'v2.3', date: '2026-01-30', accuracy: 84.5 },
    ],
  },
];

export default function RegistryPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const localePrefix = useLocalePrefix();

  const [expandedModel, setExpandedModel] = useState<string | null>(null);

  const toggleExpand = (id: string) => {
    setExpandedModel(prev => (prev === id ? null : id));
  };

  const accuracyBadge = (acc: number) =>
    acc >= 95 ? 'bg-emerald-50 text-emerald-700 border-emerald-200' :
    acc >= 90 ? 'bg-amber-50 text-amber-700 border-amber-200' :
    'bg-red-50 text-red-700 border-red-200';

  return (
    <main className="min-h-screen bg-slate-50" aria-label="모델 레지스트리">
      <DomainHeader
        title="모델 레지스트리"
        subtitle="등록된 AI 모델을 관리하고 버전 이력을 확인합니다"
        icon={Database}
        user={user ? { name: user.name || '', email: user.email || '', persona: user.persona || 'researcher' } : undefined}
        onLogout={() => signOut({ callbackUrl: '/login' })}
      />

      <div className="max-w-5xl mx-auto px-6 py-8 space-y-6">
        {/* 요약 */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <span className="text-xs font-bold text-slate-500">등록 모델 {MODELS.length}개</span>
            <span className="text-xs text-slate-400">|</span>
            <span className="text-xs text-emerald-600 font-semibold">배포됨 {MODELS.filter(m => m.deployStatus === '배포됨').length}</span>
            <span className="text-xs text-amber-600 font-semibold">스테이징 {MODELS.filter(m => m.deployStatus === '스테이징').length}</span>
            <span className="text-xs text-slate-400 font-semibold">비활성 {MODELS.filter(m => m.deployStatus === '비활성').length}</span>
          </div>
          <button className="px-4 py-2 rounded-xl bg-sky-600 text-white text-xs font-bold hover:bg-sky-700 transition-colors">
            + 새 모델 등록
          </button>
        </div>

        {/* 모델 카드 목록 */}
        <div className="space-y-4">
          {MODELS.map(model => {
            const isExpanded = expandedModel === model.id;
            return (
              <div key={model.id} className="rounded-2xl border border-slate-200 bg-white overflow-hidden transition-shadow hover:shadow-sm">
                <button
                  onClick={() => toggleExpand(model.id)}
                  className="w-full p-5 text-left"
                  aria-expanded={isExpanded}
                  aria-label={`${model.name} 상세 보기`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="h-10 w-10 rounded-xl bg-slate-900 flex items-center justify-center">
                        <Database className="h-5 w-5 text-white" />
                      </div>
                      <div>
                        <h3 className="text-sm font-bold text-slate-900">{model.name}</h3>
                        <p className="text-[11px] text-slate-400">{model.description}</p>
                      </div>
                    </div>

                    <div className="flex items-center gap-3">
                      <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${FRAMEWORK_STYLES[model.framework]}`}>
                        {model.framework}
                      </span>
                      <span className="text-xs font-mono font-bold text-slate-600">{model.version}</span>
                      <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold border ${accuracyBadge(model.accuracy)}`}>
                        {model.accuracy}%
                      </span>
                      <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold border ${DEPLOY_STYLES[model.deployStatus]}`}>
                        {model.deployStatus}
                      </span>
                      <span className="text-xs text-slate-400">{model.size}</span>
                    </div>
                  </div>
                </button>

                {isExpanded && (
                  <div className="border-t border-slate-100 px-5 pb-5 space-y-4">
                    {/* 버전 이력 */}
                    <div className="pt-4">
                      <h4 className="text-xs font-bold text-slate-500 mb-2 flex items-center gap-1">
                        <Clock className="h-3.5 w-3.5" /> 버전 이력
                      </h4>
                      <div className="space-y-1">
                        {model.versions.map((v, i) => (
                          <div key={v.version} className={`flex items-center gap-3 rounded-xl px-3 py-2 ${i === 0 ? 'bg-sky-50 border border-sky-100' : 'bg-slate-50'}`}>
                            <span className="text-xs font-mono font-bold text-slate-700 w-12">{v.version}</span>
                            <span className="text-[11px] text-slate-500">{v.date}</span>
                            <span className={`ml-auto text-[10px] font-bold ${v.accuracy >= 95 ? 'text-emerald-600' : v.accuracy >= 90 ? 'text-amber-600' : 'text-red-600'}`}>
                              정확도 {v.accuracy}%
                            </span>
                            {i === 0 && <CheckCircle className="h-3.5 w-3.5 text-sky-600" />}
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* 학습 데이터 */}
                    <div>
                      <h4 className="text-xs font-bold text-slate-500 mb-1 flex items-center gap-1">
                        <FileText className="h-3.5 w-3.5" /> 학습 데이터셋
                      </h4>
                      <p className="text-xs text-slate-600 bg-slate-50 rounded-xl px-3 py-2">{model.trainingDataset}</p>
                    </div>

                    {/* 입출력 스키마 */}
                    <div className="grid grid-cols-2 gap-3">
                      <div>
                        <h4 className="text-xs font-bold text-slate-500 mb-1 flex items-center gap-1">
                          <Code className="h-3.5 w-3.5" /> 입력 스키마
                        </h4>
                        <code className="block text-[11px] font-mono text-slate-600 bg-slate-900 text-emerald-400 rounded-xl px-3 py-2">
                          {model.inputSchema}
                        </code>
                      </div>
                      <div>
                        <h4 className="text-xs font-bold text-slate-500 mb-1 flex items-center gap-1">
                          <Code className="h-3.5 w-3.5" /> 출력 스키마
                        </h4>
                        <code className="block text-[11px] font-mono text-slate-600 bg-slate-900 text-emerald-400 rounded-xl px-3 py-2">
                          {model.outputSchema}
                        </code>
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
