/**
 * @mmup/auth — 페르소나 라우팅 맵 단위 테스트
 */
import { describe, it, expect } from 'vitest';
import { PERSONA_ROUTES, getPersonaRedirectUrl } from '../personas';
import type { Persona } from '../personas';

const ALL_PERSONAS: Persona[] = ['patient', 'doctor', 'researcher', 'pharma', 'partner', 'developer', 'admin'];

describe('Persona System', () => {
  it('7개 페르소나 라우트가 정의되어 있다', () => {
    expect(Object.keys(PERSONA_ROUTES)).toHaveLength(7);
  });

  it('모든 페르소나에 라우팅 정보(label, path, port)가 있다', () => {
    for (const persona of ALL_PERSONAS) {
      const route = PERSONA_ROUTES[persona];
      expect(route).toBeDefined();
      expect(route.label).toBeTruthy();
      expect(route.path).toBeTruthy();
      expect(route.port).toBeGreaterThan(0);
    }
  });

  it('포트 번호가 3000~3008 범위 내에 있다', () => {
    for (const key of Object.keys(PERSONA_ROUTES)) {
      const route = PERSONA_ROUTES[key as Persona];
      expect(route.port).toBeGreaterThanOrEqual(3000);
      expect(route.port).toBeLessThanOrEqual(3008);
    }
  });

  it('포트 번호가 중복되지 않는다', () => {
    const ports = Object.values(PERSONA_ROUTES).map(r => r.port);
    const uniquePorts = new Set(ports);
    expect(uniquePorts.size).toBe(ports.length);
  });

  it('patient 페르소나는 건강 앱(3008)으로 라우팅된다', () => {
    expect(PERSONA_ROUTES.patient.port).toBe(3008);
  });

  it('doctor 페르소나는 임상 콘솔(3001)로 라우팅된다', () => {
    expect(PERSONA_ROUTES.doctor.port).toBe(3001);
  });

  it('getPersonaRedirectUrl이 올바른 URL을 생성한다', () => {
    const url = getPersonaRedirectUrl('doctor');
    expect(url).toBe('http://localhost:3001/ko');
  });
});
