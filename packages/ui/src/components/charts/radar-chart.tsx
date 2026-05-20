'use client';

import React from 'react';
import { Radar, RadarChart as RechartsRadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, ResponsiveContainer, Tooltip } from 'recharts';

export interface RadarChartProps {
  data: any[];
  subject: string;
  fullMarkKey?: string;
  dataKey: string;
  color?: string;
  height?: number;
}

export function RadarChart({
  data,
  subject,
  fullMarkKey,
  dataKey,
  color = '#10b981', // emerald-500
  height = 300
}: RadarChartProps) {
  return (
    <div style={{ width: '100%', height }} data-testid="radar-chart">
      {/* @ts-expect-error React 19 type mismatch with Recharts ResponsiveContainer children */}
      <ResponsiveContainer width="100%" height="100%">
        <RechartsRadarChart cx="50%" cy="50%" outerRadius="70%" data={data}>
          <PolarGrid stroke="#e2e8f0" />
          <PolarAngleAxis dataKey={subject} tick={{ fill: '#475569', fontSize: 12, fontWeight: 500 }} />
          <PolarRadiusAxis angle={30} domain={['dataMin', 'dataMax']} tick={false} axisLine={false} />
          <Tooltip 
            contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
          />
          <Radar
            name={dataKey}
            dataKey={dataKey}
            stroke={color}
            strokeWidth={2}
            fill={color}
            fillOpacity={0.4}
            animationDuration={1500}
            animationEasing="ease-out"
          />
        </RechartsRadarChart>
      </ResponsiveContainer>
    </div>
  );
}
