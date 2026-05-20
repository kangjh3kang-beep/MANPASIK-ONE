import { http, HttpResponse, delay } from 'msw';

export const handlers = [
  // 🏥 Clinical (임상) 도메인 모의 API
  http.get('*/api/v1/patients/*/vitals', async () => {
    return HttpResponse.json({
      data: [
        { timestamp: '08:00', heartRate: 72, bloodPressureSys: 120, spo2: 98 },
        { timestamp: '09:00', heartRate: 75, bloodPressureSys: 122, spo2: 97 },
      ]
    });
  }),

  // 🧪 Predictor (예측) 도메인 모의 API
  http.get('*/api/v1/predictor/risk-scores', async () => {
    return HttpResponse.json({
      data: [
        { disease: '제2형 당뇨', riskScore: 0.82, contributingFactors: ['BMI'], trend: 'up' },
        { disease: '관상동맥 질환', riskScore: 0.15, contributingFactors: ['Cholesterol'], trend: 'stable' },
      ]
    });
  }),

  // 🛡️ GxP (감사) 도메인 모의 API
  http.get('*/api/v1/gxp/audit-logs', async () => {
    return HttpResponse.json({
      data: [
        { id: '1', user: '김품질', action: '의약품 입고 승인', batchId: 'B1', hasSignature: true, timestamp: '14:20' },
      ]
    });
  }),

  http.get('*/api/v1/gxp/compliance', async () => {
    return HttpResponse.json({
      data: [
        { ruleId: 'KGMP 의약품 제조관리', status: 'pass', itemsCount: 41, passedCount: 41 },
      ]
    });
  }),
];
