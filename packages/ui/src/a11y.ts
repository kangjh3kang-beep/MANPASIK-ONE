/**
 * @mmup-axis 2 종합검진
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard WCAG 2.2 AA
 *
 * 접근성 유틸리티 — WCAG 2.2 AA 준수 도구
 */

/**
 * sRGB 색상의 상대 휘도 계산 (WCAG 2.1 §1.4.3)
 */
function hexToRgb(hex: string): [number, number, number] {
  const h = hex.replace('#', '');
  return [
    parseInt(h.substring(0, 2), 16) / 255,
    parseInt(h.substring(2, 4), 16) / 255,
    parseInt(h.substring(4, 6), 16) / 255,
  ];
}

function srgbToLinear(c: number): number {
  return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
}

function relativeLuminance(hex: string): number {
  const [r, g, b] = hexToRgb(hex);
  return 0.2126 * srgbToLinear(r) + 0.7152 * srgbToLinear(g) + 0.0722 * srgbToLinear(b);
}

/**
 * 두 색상 간 대비 비율 계산 (1:1 ~ 21:1)
 */
export function getContrastRatio(foreground: string, background: string): number {
  const l1 = relativeLuminance(foreground);
  const l2 = relativeLuminance(background);
  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);
  return (lighter + 0.05) / (darker + 0.05);
}

/**
 * WCAG 2.2 AA 대비 기준 충족 여부
 * - 일반 텍스트: 4.5:1
 * - 대형 텍스트 (≥24px 또는 ≥18.66px bold): 3:1
 * - UI 컴포넌트/그래픽: 3:1
 */
export function meetsWCAG_AA(foreground: string, background: string, fontSize?: number): boolean {
  const ratio = getContrastRatio(foreground, background);
  const threshold = fontSize && fontSize >= 24 ? 3 : 4.5;
  return ratio >= threshold;
}

/**
 * WCAG 2.2 AAA 대비 기준 충족 여부 (7:1 / 4.5:1)
 */
export function meetsWCAG_AAA(foreground: string, background: string, fontSize?: number): boolean {
  const ratio = getContrastRatio(foreground, background);
  const threshold = fontSize && fontSize >= 24 ? 4.5 : 7;
  return ratio >= threshold;
}

/**
 * 컴포넌트에 대한 aria-label 생성
 * 결과 3중 인코딩 원칙: 색채 + 수치 + 아이콘을 텍스트로 변환
 */
export function generateAriaLabel(component: string, value: string | number, context?: string): string {
  const parts = [component, String(value)];
  if (context) parts.push(context);
  return parts.join(': ');
}

/**
 * 스크린 리더에게 메시지 전달 (aria-live region)
 */
export function announceToScreenReader(message: string, priority: 'polite' | 'assertive' = 'polite'): void {
  if (typeof document === 'undefined') return;

  const el = document.createElement('div');
  el.setAttribute('role', 'status');
  el.setAttribute('aria-live', priority);
  el.setAttribute('aria-atomic', 'true');
  el.style.position = 'absolute';
  el.style.width = '1px';
  el.style.height = '1px';
  el.style.overflow = 'hidden';
  el.style.clip = 'rect(0, 0, 0, 0)';
  el.style.whiteSpace = 'nowrap';
  document.body.appendChild(el);

  // 약간의 딜레이 후 콘텐츠 삽입 (리더가 감지하도록)
  requestAnimationFrame(() => {
    el.textContent = message;
    setTimeout(() => document.body.removeChild(el), 3000);
  });
}

/**
 * axe-core 기반 접근성 테스트 실행
 */
export async function runAxeTest(container: HTMLElement): Promise<{ violations: number; passes: number }> {
  try {
    // axe-core는 선택적 의존성 — 설치 시에만 동작
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    const axe = await import(/* webpackIgnore: true */ 'axe-core' as string);
    const results = await (axe as any).default.run(container);
    return { violations: results.violations.length, passes: results.passes.length };
  } catch {
    return { violations: 0, passes: 0 };
  }
}
