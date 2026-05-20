/**
 * 공통 API 타입 정의
 */

// ===== 환자 =====
export interface Patient {
  id: string;
  name: string;
  age: number;
  gender: 'M' | 'F';
  diagnosisCode: string;
  riskLevel: 'low' | 'medium' | 'high' | 'critical';
  lastVisit: string;
  assignedDoctor: string;
}

// ===== 바이오마커 =====
export interface Biomarker {
  id: string;
  name: string;
  value: number;
  unit: string;
  normalMin: number;
  normalMax: number;
  category: string;
  measuredAt: string;
}

// ===== 센서 =====
export interface SensorDevice {
  id: string;
  name: string;
  type: string;
  status: 'online' | 'offline' | 'warning';
  battery: number;
  signal: number;
  firmware: string;
  lastSeen: string;
}

// ===== 생체신호 =====
export interface VitalSign {
  timestamp: string;
  heartRate: number;
  spo2: number;
  temperature: number;
  bloodPressureSys: number;
  bloodPressureDia: number;
  respiratoryRate: number;
}

// ===== AI 에이전트 =====
export interface AIAgent {
  id: string;
  name: string;
  model: string;
  status: 'running' | 'idle' | 'error';
  accuracy: number;
  latency: number;
  totalRequests: number;
}

// ===== 리워드 =====
export interface RewardBalance {
  userId: string;
  balance: number;
  currency: 'MPS';
  totalEarned: number;
  weeklyChange: number;
}

export interface RewardContribution {
  id: string;
  type: string;
  points: number;
  quality: string;
  date: string;
}

// ===== 질환 예측 =====
export interface DiseasePrediction {
  disease: string;
  riskScore: number;
  trend: 'up' | 'down' | 'stable';
  contributingFactors: string[];
  modelVersion: string;
}

// ===== 파트너 =====
export interface PartnerOrg {
  id: string;
  name: string;
  protocol: string;
  status: 'active' | 'syncing' | 'inactive';
  totalExchanges: number;
  lastSync: string;
}

// ===== GxP =====
export interface AuditLog {
  id: string;
  timestamp: string;
  user: string;
  action: string;
  batchId: string;
  hasSignature: boolean;
}

export interface ComplianceRule {
  id: string;
  rule: string;
  totalItems: number;
  passedItems: number;
  status: 'pass' | 'warning' | 'fail';
}
