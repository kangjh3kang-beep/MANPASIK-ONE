import { Page } from '@playwright/test';

export async function setupApiMocks(page: Page) {
  // 브라우저 로그 중계 (디버깅)
  page.on('console', msg => {
    if (msg.type() === 'error') console.log(`[Browser Error] ${msg.text()}`);
  });

  // 모든 API 요청 가로채기 (강력한 정규표현식 매칭)
  await page.route(/\/api\/v1\//, async (route) => {
    const url = route.request().url();
    
    // Clinical Vitals - 현실 데이터 반영
    if (url.includes('vitals')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            { timestamp: '08:00', heartRate: 72, bloodPressureSys: 120, spo2: 98 },
            { timestamp: '09:00', heartRate: 75, bloodPressureSys: 122, spo2: 97 },
          ]
        }),
      });
    } 
    // Predictor Risk Scores - 현실 데이터 반영 (관상동맥 질환)
    else if (url.includes('risk-scores')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            { disease: '제2형 당뇨', riskScore: 0.82, contributingFactors: ['BMI'], trend: 'up' },
            { disease: '관상동맥 질환', riskScore: 0.15, contributingFactors: ['Cholesterol'], trend: 'stable' },
          ]
        }),
      });
    } 
    // GxP Audit Logs - 현실 데이터 반영 (김품질)
    else if (url.includes('audit-logs')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            { id: '1', user: '김품질', action: '의약품 입고 승인', batchId: 'B2026-001', hasSignature: true, timestamp: '14:20' },
          ]
        }),
      });
    } 
    // GxP Compliance
    else if (url.includes('compliance')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            { ruleId: 'KGMP 의약품 제조관리', status: 'pass', itemsCount: 41, passedCount: 41 },
          ]
        }),
      });
    } else {
      await route.continue();
    }
  });
}
