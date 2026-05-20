import { test, expect } from '@playwright/test';
import { setupApiMocks } from './support/mock-api';

test.describe('Clinical Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Playwright 레벨에서 API 모킹 설정
    await setupApiMocks(page);
    
    await page.goto('/domains/clinical');
    await page.waitForLoadState('networkidle');
    // 데이터 카드가 나타날 때까지 명시적 대기
    await page.waitForSelector('text=심박수', { timeout: 10000 });
  });

  test('대시보드 헤더와 환자 정보가 올바르게 표시되어야 한다', async ({ page }) => {
    await expect(page.getByRole('heading', { name: '임상 데이터 콘솔' })).toBeVisible();
    await expect(page.getByText('활성 환자')).toBeVisible();
  });

  test('바이탈 사인 카드 3개가 모두 렌더링되어야 한다', async ({ page }) => {
    await expect(page.getByText('심박수')).toBeVisible();
    await expect(page.getByText('산소포화도')).toBeVisible();
    await expect(page.getByText('수축기 혈압')).toBeVisible();
  });

  test('시각화 컴포넌트(스파크라인)가 렌더링되어야 한다', async ({ page }) => {
    // SVG 요소가 존재하는지 확인 (리차트 의존성 제거)
    const svgs = page.locator('svg polyline');
    expect(await svgs.count()).toBeGreaterThanOrEqual(1);
  });
});
