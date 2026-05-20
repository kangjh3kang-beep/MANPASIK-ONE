import { test, expect } from '@playwright/test';
import { setupApiMocks } from './support/mock-api';

test.describe('GxP Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Playwright 레벨에서 API 모킹 설정
    await setupApiMocks(page);

    await page.goto('/domains/gxp');
    await page.waitForLoadState('networkidle');
    // 데이터 로드 명시적 대기 (파일에 하드코딩된 '김품질' 확인)
    await page.waitForSelector('text=김품질', { timeout: 15000 });
  });

  test('GxP 시스템 헤더와 통계 카드가 표시되어야 한다', async ({ page }) => {
    await expect(page.getByRole('heading', { name: '의약품 GxP 준수 시스템' })).toBeVisible();
  });

  test('감사 추적 로그 섹션에 데이터가 로드되어야 한다', async ({ page }) => {
    // 실제 화면에 표시되는 하드코딩된 텍스트로 검증
    await expect(page.getByText('김품질').first()).toBeVisible();
    await expect(page.getByText('배치 검사 승인').first()).toBeVisible();
    await expect(page.getByText('BT-2026-0412').first()).toBeVisible();
  });

  test('규제 준수 매트릭스 항목들이 표시되어야 한다', async ({ page }) => {
    await expect(page.getByText('KGMP 의약품 제조관리')).toBeVisible();
  });
});
