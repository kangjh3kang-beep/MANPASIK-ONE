/**
 * @mmup/api-client — 타입 정의 단위 테스트
 */
import { describe, it, expect } from 'vitest';
import type { Patient, Biomarker, AIAgent, VitalSign, RewardBalance, DiseasePrediction, PartnerOrg, AuditLog, ComplianceRule } from '../types';

describe('API Types', () => {
  it('Patient 타입이 올바른 구조를 가진다', () => {
    const patient: Patient = {
      id: 'p-001',
      name: '김환자',
      age: 45,
      gender: 'M',
      diagnosisCode: 'E11.9',
      riskLevel: 'medium',
      lastVisit: '2026-04-26',
      assignedDoctor: '이의사',
    };
    expect(patient.id).toBe('p-001');
    expect(patient.gender).toMatch(/^[MF]$/);
    expect(['low', 'medium', 'high', 'critical']).toContain(patient.riskLevel);
  });

  it('VitalSign 타입이 생체 데이터를 포함한다', () => {
    const vital: VitalSign = {
      timestamp: '2026-04-26T14:30:00Z',
      heartRate: 72,
      spo2: 98,
      temperature: 36.5,
      bloodPressureSys: 120,
      bloodPressureDia: 78,
      respiratoryRate: 16,
    };
    expect(vital.heartRate).toBeGreaterThan(0);
    expect(vital.spo2).toBeLessThanOrEqual(100);
  });

  it('AIAgent 상태가 올바른 enum 값만 허용한다', () => {
    const validStatuses: AIAgent['status'][] = ['running', 'idle', 'error'];
    validStatuses.forEach(status => {
      const agent: AIAgent = {
        id: 'a-001', name: 'Test', model: 'v1',
        status, accuracy: 95, latency: 50, totalRequests: 100,
      };
      expect(validStatuses).toContain(agent.status);
    });
  });

  it('DiseasePrediction 위험도가 0~1 범위다', () => {
    const pred: DiseasePrediction = {
      disease: '제2형 당뇨',
      riskScore: 0.82,
      trend: 'up',
      contributingFactors: ['HbA1c 6.8%'],
      modelVersion: 'v3.2',
    };
    expect(pred.riskScore).toBeGreaterThanOrEqual(0);
    expect(pred.riskScore).toBeLessThanOrEqual(1);
  });

  it('ComplianceRule 준수율이 논리적으로 맞다', () => {
    const rule: ComplianceRule = {
      id: 'r-001', rule: 'cGMP 21 CFR Part 211',
      totalItems: 47, passedItems: 47, status: 'pass',
    };
    expect(rule.passedItems).toBeLessThanOrEqual(rule.totalItems);
  });
});
