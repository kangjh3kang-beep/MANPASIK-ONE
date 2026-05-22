/**
 * 도메인별 API 엔드포인트 정의
 */
import { apiClient } from './client';
import type {
  Patient, Biomarker, SensorDevice, VitalSign, AIAgent,
  RewardBalance, RewardContribution, DiseasePrediction,
  PartnerOrg, AuditLog, ComplianceRule,
  AgentPipelineState, DevPortalKey, DevPortalEndpoint,
  HardwareDevice, HardwareDiagnostics, AppMetric,
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

// ===== 에이전트 파이프라인 =====
export const agentPipelineApi = {
  getState: () =>
    apiClient.get<AgentPipelineState>('/api/v1/agents/pipeline'),
};

// ===== 개발자 포털 =====
export const devPortalApi = {
  getKeys: () =>
    apiClient.get<DevPortalKey[]>('/api/v1/dev-portal/keys'),
  getEndpoints: () =>
    apiClient.get<DevPortalEndpoint[]>('/api/v1/dev-portal/endpoints'),
};

// ===== 하드웨어 =====
export const hardwareApi = {
  getDevices: () =>
    apiClient.get<HardwareDevice[]>('/api/v1/hardware/devices'),
  getDiagnostics: () =>
    apiClient.get<HardwareDiagnostics>('/api/v1/hardware/diagnostics'),
};

// ===== 앱 메트릭 =====
export const appMetricsApi = {
  getMetrics: () =>
    apiClient.get<AppMetric[]>('/api/v1/app/metrics'),
};
