/**
 * React Query 기반 데이터 훅
 * 각 도메인의 API 호출을 선언적 훅으로 래핑
 */
'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { patientsApi, sensorsApi, agentsApi, predictionsApi, rewardsApi, partnersApi, gxpApi } from './endpoints';

// ===== 환자 =====
export function usePatients(params?: Record<string, string>) {
  return useQuery({
    queryKey: ['patients', params],
    queryFn: () => patientsApi.list(params),
    staleTime: 30_000,
  });
}

export function usePatient(id: string) {
  return useQuery({
    queryKey: ['patient', id],
    queryFn: () => patientsApi.getById(id),
    enabled: !!id,
  });
}

export function usePatientVitals(patientId: string) {
  return useQuery({
    queryKey: ['patient-vitals', patientId],
    queryFn: () => patientsApi.getVitals(patientId),
    refetchInterval: 5_000, // 5초마다 자동 갱신
    enabled: !!patientId,
  });
}

export function usePatientBiomarkers(patientId: string) {
  return useQuery({
    queryKey: ['patient-biomarkers', patientId],
    queryFn: () => patientsApi.getBiomarkers(patientId),
    staleTime: 60_000,
    enabled: !!patientId,
  });
}

// ===== 센서 =====
export function useSensors() {
  return useQuery({
    queryKey: ['sensors'],
    queryFn: () => sensorsApi.list(),
    refetchInterval: 10_000, // 10초마다 자동 갱신
  });
}

// ===== AI 에이전트 =====
export function useAgents() {
  return useQuery({
    queryKey: ['agents'],
    queryFn: () => agentsApi.list(),
    refetchInterval: 15_000,
  });
}

export function useRunPrediction() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ agentId, input }: { agentId: string; input: unknown }) =>
      agentsApi.runPrediction(agentId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['predictions'] });
    },
  });
}

// ===== 예측 =====
export function usePredictions(patientId?: string) {
  return useQuery({
    queryKey: ['predictions', patientId || 'all'],
    queryFn: () => patientId 
      ? predictionsApi.getByPatient(patientId) 
      : predictionsApi.getRiskScores(),
  });
}

// ===== 리워드 =====
export function useRewardBalance(userId: string) {
  return useQuery({
    queryKey: ['reward-balance', userId],
    queryFn: () => rewardsApi.getBalance(userId),
    staleTime: 60_000,
    enabled: !!userId,
  });
}

export function useRewardContributions(userId: string) {
  return useQuery({
    queryKey: ['reward-contributions', userId],
    queryFn: () => rewardsApi.getContributions(userId),
    enabled: !!userId,
  });
}

// ===== 파트너 =====
export function usePartners() {
  return useQuery({
    queryKey: ['partners'],
    queryFn: () => partnersApi.list(),
    staleTime: 120_000,
  });
}

// ===== GxP =====
export function useAuditLogs(params?: Record<string, string>) {
  return useQuery({
    queryKey: ['audit-logs', params],
    queryFn: () => gxpApi.getAuditLogs(params),
  });
}

export function useCompliance() {
  return useQuery({
    queryKey: ['compliance'],
    queryFn: () => gxpApi.getCompliance(),
    staleTime: 300_000,
  });
}
