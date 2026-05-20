import React from 'react';
import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { RiskProgress } from '../risk-progress';

describe('RiskProgress Component', () => {
  it('질환명과 퍼센테이지가 렌더링되어야 한다', () => {
    const { getByText } = render(<RiskProgress diseaseName="제2형 당뇨" riskScore={0.82} />);
    expect(getByText('제2형 당뇨')).toBeTruthy();
    expect(getByText(/82\s*%/)).toBeTruthy();
    expect(getByText('고위험')).toBeTruthy(); // 0.75 이상이므로
  });

  it('요인(factors)이 주어진 경우 렌더링되어야 한다', () => {
    const { getByText } = render(
      <RiskProgress diseaseName="고혈압" riskScore={0.65} factors={['수축기 138', '가족력']} />
    );
    expect(getByText('고혈압')).toBeTruthy();
    expect(getByText('주의')).toBeTruthy(); // 0.5 이상 0.75 미만이므로
    expect(getByText('수축기 138')).toBeTruthy();
    expect(getByText('가족력')).toBeTruthy();
  });

  it('우수 범위의 리스크 스코어(0.2)는 우수를 표시해야 한다', () => {
    const { getByText } = render(<RiskProgress diseaseName="간경변" riskScore={0.2} />);
    expect(getByText('우수')).toBeTruthy();
    expect(getByText(/20\s*%/)).toBeTruthy();
  });
});
