import { test, expect } from '@playwright/test';
import { setupApiMocks } from './support/mock-api';

test.describe('Predictor Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Playwright 레벨에서 API 모킹 설정
    await setupApiMocks(page);

    await page.goto('/domains/predictor');
    await page.waitForLoadState('networkidle');
    // 데이터가 로드될 때까지 명시적 대기
    await page.waitForSelector('text=관상동맥 질환', { timeout: 10000 });
  });

  test('AI 예측 분석 헤더가 표시되어야 한다', async ({ page }) => {
    await expect(page.getByRole('heading', { name: '생체 지표 예측 엔진' })).toBeVisible();
  });

  test('질환 위험도 카드가 표시되어야 한다', async ({ page }) => {
    // 모킹 데이터 기반 확인
    await expect(page.getByText('제2형 당뇨')).toBeVisible();
    await expect(page.getByText('관상동맥 질환')).toBeVisible();
  });
});
