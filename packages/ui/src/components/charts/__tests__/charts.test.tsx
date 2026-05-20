import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/react';
import { AreaChart } from '../area-chart';
import { BarChart } from '../bar-chart';
import { RadarChart } from '../radar-chart';

// Recharts relies on DOM APIs (ResizeObserver) not available in jsdom by default
// Mocking ResizeObserver for tests
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
})) as any;

describe('Chart Components', () => {
  const commonData = [
    { name: 'A', value: 400 },
    { name: 'B', value: 300 },
  ];

  it('AreaChart가 에러 없이 렌더링되어야 한다', () => {
    const { container } = render(
      <AreaChart data={commonData} xKey="name" yKey="value" />
    );
    // ResponsiveContainer creates a wrapper div. We check if the custom chart wrapper exists.
    expect(container.querySelector('[data-testid="area-chart"]')).toBeTruthy();
  });

  it('BarChart가 커스텀 색상과 함께 렌더링되어야 한다', () => {
    const { container } = render(
      <BarChart data={commonData} xKey="name" yKey="value" color="#ff0000" />
    );
    expect(container.querySelector('[data-testid="bar-chart"]')).toBeTruthy();
  });

  it('RadarChart가 다양한 차원 데이터를 렌더링해야 한다', () => {
    const radarData = [
      { subject: 'Math', A: 120, fullMark: 150 },
      { subject: 'Chinese', A: 98, fullMark: 150 },
      { subject: 'English', A: 86, fullMark: 150 },
    ];
    
    const { container } = render(
      <RadarChart data={radarData} subject="subject" dataKey="A" />
    );
    expect(container.querySelector('[data-testid="radar-chart"]')).toBeTruthy();
  });
});
