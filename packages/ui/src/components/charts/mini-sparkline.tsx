'use client';

import React from 'react';

interface MiniSparklineProps {
  data: number[];
  color?: string;
  width?: number;
  height?: number;
}

export function MiniSparkline({ data, color = '#0ea5e9', width = 120, height = 32 }: MiniSparklineProps) {
  if (data.length < 2) return null;
  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;
  const points = data.map((v, i) => `${(i / (data.length - 1)) * width},${height - ((v - min) / range) * height}`).join(' ');

  return (
    <svg width={width} height={height} className="inline-block" aria-hidden="true">
      <polyline fill="none" stroke={color} strokeWidth="2" strokeLinejoin="round" points={points} />
    </svg>
  );
}
