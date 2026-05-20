/**
 * 도메인별 API 엔드포인트 정의
 */
import { apiClient } from './client';
import type {
  Patient, Biomarker, SensorDevice, VitalSign, AIAgent,
  RewardBalance, RewardContribution, DiseasePrediction,
  PartnerOrg, AuditLog, ComplianceRule,
} from './types';

// ===== 환자 =====
export const patientsApi = {
  list: (params?: Record<string, string>) =>
    apiClient.get<Patient[]>('/api/v1/patients', params),
  getById: (id: string) =>
    apiClient.get<Patient>(`/api/v1/patients/${id}`),
  getVitals: (patientId: string) =>
    apiClient.get<VitalSign[]>(`/api/v1/patients/${patientId}/vitals`),
  getBiomarkers: (patientId: string) =>
    apiClient.get<Biomarker[]>(`/api/v1/patients/${patientId}/biomarkers`),
};

// ===== 센서 =====
export const sensorsApi = {
  list: () =>
    apiClient.get<SensorDevice[]>('/api/v1/sensors'),
  getById: (id: string) =>
    apiClient.get<SensorDevice>(`/api/v1/sensors/${id}`),
};

// ===== AI 에이전트 =====
export const agentsApi = {
  list: () =>
    apiClient.get<AIAgent[]>('/api/v1/agents'),
  getById: (id: string) =>
    apiClient.get<AIAgent>(`/api/v1/agents/${id}`),
  runPrediction: (agentId: string, input: unknown) =>
    apiClient.post<DiseasePrediction>(`/api/v1/agents/${agentId}/predict`, input),
};

// ===== 예측 =====
export const predictionsApi = {
  getRiskScores: () =>
    apiClient.get<DiseasePrediction[]>('/api/v1/predictor/risk-scores'),
  getByPatient: (patientId: string) =>
    apiClient.get<DiseasePrediction[]>(`/api/v1/predictions/patient/${patientId}`),
  run: (data: { patientId: string; models: string[] }) =>
    apiClient.post<DiseasePrediction[]>('/api/v1/predictions/run', data),
};

// ===== 리워드 =====
export const rewardsApi = {
  getBalance: (userId: string) =>
    apiClient.get<RewardBalance>(`/api/v1/rewards/${userId}/balance`),
  getContributions: (userId: string) =>
    apiClient.get<RewardContribution[]>(`/api/v1/rewards/${userId}/contributions`),
};

// ===== 파트너 =====
export const partnersApi = {
  list: () =>
    apiClient.get<PartnerOrg[]>('/api/v1/partners'),
  sendFhirBundle: (bundle: unknown) =>
    apiClient.post('/api/v1/fhir/bundle', bundle),
};

// ===== GxP 감사 =====
export const gxpApi = {
  getAuditLogs: (params?: Record<string, string>) =>
    apiClient.get<AuditLog[]>('/api/v1/gxp/audit-logs', params),
  getCompliance: () =>
    apiClient.get<ComplianceRule[]>('/api/v1/gxp/compliance'),
};
