/**
 * @mmup/api-client 배럴 익스포트
 */

// Core client
export { apiClient, MmupApiClient } from './client';
export type { ApiConfig, ApiError, ApiResponse } from './client';

// Types
export type {
  Patient, Biomarker, SensorDevice, VitalSign, AIAgent,
  RewardBalance, RewardContribution, DiseasePrediction,
  PartnerOrg, AuditLog, ComplianceRule,
  AgentPipelineState, DevPortalKey, DevPortalEndpoint,
  HardwareDevice, HardwareDiagnostics, AppMetric,
  MeasurementRecord, TelemedicineDoctor, StoreProduct,
  AdminKPI, AdminEvent,
} from './types';

// Endpoints
export {
  patientsApi, sensorsApi, agentsApi, predictionsApi,
  rewardsApi, partnersApi, gxpApi,
  agentPipelineApi, devPortalApi, hardwareApi, appMetricsApi,
  measurementsApi, telemedicineApi, storeApi, adminApi,
} from './endpoints';

// React Query Hooks
export {
  usePatients, usePatient, usePatientVitals, usePatientBiomarkers,
  useSensors, useAgents, useRunPrediction, usePredictions,
  useRewardBalance, useRewardContributions,
  usePartners, useAuditLogs, useCompliance,
  useAgentPipeline, useDevPortalKeys, useDevPortalEndpoints,
  useHardwareDevices, useHardwareDiagnostics, useAppMetrics,
  useMeasurementHistory, useTelemedicineDoctors, useStoreProducts,
  useAdminKPIs, useAdminEvents,
} from './hooks';

// Mocks
export { initMocks } from './mocks';
