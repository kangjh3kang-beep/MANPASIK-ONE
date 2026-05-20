import os

base_dir = r"\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik"

files_to_create = {
    "pnpm-workspace.yaml": """packages:
  - "apps/*"
  - "packages/*"
  - "services/*"
""",
    "turbo.json": """{
  "$schema": "https://turbo.build/schema.json",
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "!.next/cache/**", "dist/**"]
    },
    "lint": {
      "dependsOn": ["^lint"]
    },
    "typecheck": {
      "dependsOn": ["^typecheck"]
    },
    "test": {
      "dependsOn": ["^build"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    }
  }
}
""",
    "packages/tsconfig/base.json": """{
  "$schema": "https://json.schemastore.org/tsconfig",
  "display": "Default",
  "compilerOptions": {
    "composite": false,
    "declaration": true,
    "declarationMap": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "inlineSources": false,
    "isolatedModules": true,
    "moduleResolution": "node",
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "preserveWatchOutput": true,
    "skipLibCheck": true,
    "strict": true,
    "noImplicitAny": true,
    "exactOptionalPropertyTypes": true,
    "noUncheckedIndexedAccess": true
  },
  "exclude": ["node_modules"]
}
""",
    "packages/tsconfig/nextjs.json": """{
  "$schema": "https://json.schemastore.org/tsconfig",
  "display": "Next.js",
  "extends": "./base.json",
  "compilerOptions": {
    "plugins": [{ "name": "next" }],
    "allowJs": true,
    "jsx": "preserve",
    "lib": ["dom", "dom.iterable", "esnext"],
    "module": "esnext",
    "resolveJsonModule": true,
    "target": "es5"
  }
}
""",
    "packages/tsconfig/nestjs.json": """{
  "$schema": "https://json.schemastore.org/tsconfig",
  "display": "NestJS",
  "extends": "./base.json",
  "compilerOptions": {
    "module": "commonjs",
    "target": "es2021",
    "sourceMap": true,
    "outDir": "./dist",
    "baseUrl": "./",
    "experimentalDecorators": true,
    "emitDecoratorMetadata": true
  }
}
""",
    "packages/tsconfig/vitest.json": """{
  "$schema": "https://json.schemastore.org/tsconfig",
  "display": "Vitest",
  "extends": "./base.json",
  "compilerOptions": {
    "lib": ["dom", "dom.iterable", "esnext"],
    "types": ["vitest/globals"]
  }
}
""",
    "packages/eslint-config/base.js": """/**
 * @type {import("eslint").Linter.Config}
 */
module.exports = {
  parser: "@typescript-eslint/parser",
  plugins: ["@typescript-eslint"],
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended"
  ],
  rules: {
    "@typescript-eslint/no-explicit-any": "error", // Use unknown instead of any
    "no-console": ["error", { allow: ["warn", "error"] }]
  }
};
""",
    "packages/eslint-config/nextjs.js": """module.exports = {
  extends: ["next/core-web-vitals", "./base.js"]
};
""",
    "packages/eslint-config/nestjs.js": """module.exports = {
  extends: ["./base.js"],
  env: {
    node: true,
    jest: true
  }
};
""",
    "packages/ssot/src/index.ts": """/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 * 
 * MMUP SSOT v1.6 — 영구 단일출처 (변경 금지·Board 결의 시만 갱신)
 * 만파식: 모든 축 / 단계: 모든 단계의 기준
 */

export const MMUP = {
  fullName_KO: "MPS 만파식 다중측정원리 유니버설 POCT 플랫폼(MMUP)",
  fullName_EN: "MPS Manpasik Multi-Measurement Universal POCT Platform (MMUP)",
  abbr: "MMUP",
} as const;

export const FAMILY = {
  A: { id: "A", appNo: "APP2026-0022KR", title: "디지털 트윈 기반의 적응형 유체장 제어 장치 및 방법" },
  B: { id: "B", appNo: "APP2025-0967KR", title: "차동 측정 기반 범용 분석 장치 및 방법" },
  C: { id: "C", appNo: "APP2025-0968KR", title: "플랫폼 공통층과 결합된 모듈러 측정 장치 및 방법" },
} as const;

export const TRINITY = {
  IP1: { name: "Physical Matrix Removal", primary: "A", secondary: "B", tertiary: "C" },
  IP2: { name: "Differential Measurement", primary: "B", secondary: "A", tertiary: "C" },
  IP3: { name: "AI Correction (distributed)", primary: "C", secondary: "A", tertiary: "B" },
} as const;

export const BUSINESS_SSOT = {
  fxKRWperUSD: 1480,
  skuTotal: 44,
  layerSkus: { L1: 16, L2: 12, L3: 10, L4: 8, L5: 6 },
  clinicalAssays: 169,
  tamUSD_B: 34.47,
  samUSD_B: 8.0,
  somUSD_B: 0.13,
  fundingKRW_B: { seriesA: 6.5, seriesB: 5.0, seriesC: 15.0 },
} as const;

export const MFG_SSOT = {
  connector: "Samtec MECF-08-01-L-DV (16핀, 1.27mm)",
  parallelPlane_mm: { length: 49.7, width: 30, height: 4.3 },
  differentialFormula: "Sdiff_n = S_n - alpha_n * R_n",
} as const;

export const FAMILY_C_PARTS = {
  dockingInterface: 110,
  dataPacketEngine: 130,
  uncertainty: "135-P5",
  updateRollback: 140,
  safetyGuard: 150,
  multimodalScheduler: 160,
  quietWindow: 164,
  sensorFusion: 170,
  hmiSafety: 175,
  securityModule: 190,
  tpm: "191-T",
  pcr: 194,
  measurementEngines: [200, 300, 400],
} as const;

export const DEPRECATED_FOREVER = [
  "9블록", "9-block", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9",
  "Family L",
  "100% 완성", "Zero 오류",
  "소형 정량 면역분석기", "혈액 전용 분석기", "유니버설 POCT 분석장치 (단축)",
] as const;

export const PRIORITY = [
  "환자 안전",
  "규제 준수",
  "SSOT 무결성",
  "사실성",
  "절차",
  "형식·간결성",
] as const;
""",
    "packages/types/src/loop.ts": """/**
 * @mmup-axis 8 이행·검증·보상
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export const LOOP_STAGES = {
  1: "measure",
  2: "diagnose",
  3: "predict",
  4: "prescribe",
  5: "execute",
  6: "verify",
  7: "reward",
  8: "remeasure",
} as const;

export type LoopStage = keyof typeof LOOP_STAGES;

export const PRESCRIPTION_TYPES = [
  "clinic_referral",
  "diet_adjustment",
  "exercise_plan",
  "environment_setup",
  "supplement_recommendation",
] as const;
""",
    "packages/types/src/fhir.ts": """/**
 * @mmup-axis 3 종합 데이터 분석
 * @mmup-stage 2 진단
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */

export interface Observation {
  resourceType: string;
  meta?: Record<string, any>;
  [key: string]: any;
}

export interface MMUPMeasurement extends Observation {
  meta: Observation["meta"] & {
    deviceId: string;
    cartridgeUDI: string;
    skuLayer: "L1" | "L2" | "L3" | "L4" | "L5";
    skuId: string;
    confidence_135P5: number;
    differentialMode: boolean;
    twinModelVersion: string;
    aiCorrectionApplied: boolean;
    hashChain: string;
  };
}
""",
    "packages/types/src/persona.ts": """/**
 * @mmup-axis 5 분야별 전문의 에이전트
 * @mmup-stage 4 처방
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export const PERSONA = {
  P1: "일반 사용자",
  P2: "만성질환자",
  P3: "부모(영유아)",
  P4: "전문의",
  P5: "검사기사·약사",
  P6: "투자자·파트너",
  P7: "규제·감사"
} as const;

export type PersonaType = keyof typeof PERSONA;
""",
    "packages/types/src/family.ts": """/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export type FamilyId = "A" | "B" | "C";

export interface FamilyMapping {
  family: FamilyId;
  description: string;
}
""",
    "packages/i18n/ko.json": """{
  "welcome": "MPS 만파식 다중측정원리 유니버설 POCT 플랫폼에 오신 것을 환영합니다.",
  "disclaimer": "본 정보는 의료 자문 보조이며 진료 행위를 대체하지 않습니다.",
  "emergency": "의료 상담 내역을 바탕으로 즉각적인 응급 처치가 필요합니다. 주치의 상담을 권장합니다.",
  "stages": {
    "measure": "[1 측정] 측정을 진행합니다.",
    "diagnose": "[2 진단] 결과 분석 중...",
    "predict": "[3 예측] 예측 모델 동작 중...",
    "prescribe": "[4 처방] 전문의 에이전트가 결과를 제공합니다.",
    "execute": "[5 이행] 처방 이행 등록",
    "verify": "[6 검증] Family B 차동 비교 완료",
    "reward": "[7 보상] 포인트 수령",
    "remeasure": "[8 재측정] 재측정을 예약합니다."
  }
}
""",
    "packages/i18n/en.json": """{
  "welcome": "Welcome to MPS Manpasik Multi-Measurement Universal POCT Platform.",
  "disclaimer": "This information is intended as a supplementary medical advisory and does not replace medical diagnosis.",
  "emergency": "Based on the consultation, immediate care may be required. Please consult a physician.",
  "stages": {
    "measure": "[1 Measure] Measuring in progress...",
    "diagnose": "[2 Diagnose] Analyzing...",
    "predict": "[3 Predict] Predictive model running...",
    "prescribe": "[4 Prescribe] Specialist Agent recommends:",
    "execute": "[5 Execute] Execute prescription",
    "verify": "[6 Verify] Differential check complete",
    "reward": "[7 Reward] Rewards issued",
    "remeasure": "[8 Remeasure] Scheduling next measurement..."
  }
}
""",
    "packages/i18n/ja.json": """{
  "welcome": "MPS 萬波息 多重測定原理 Universal POCT Platform へようこそ。",
  "disclaimer": "本情報は医療アドバイスの補助であり、診療行為を代替するものではありません。"
}
""",
    "packages/i18n/zh.json": """{
  "welcome": "欢迎使用 MPS 万波息多重测量原理通用 POCT 平台。",
  "disclaimer": "本信息仅作为医疗咨询辅助，不能代替医疗诊断。"
}
""",
    "packages/trinity/src/index.ts": """/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

import { TRINITY } from "@mmup/ssot";

export function getTrinityMapping(ip: keyof typeof TRINITY) {
  return TRINITY[ip];
}

export function isProperlyDistributed(ip: keyof typeof TRINITY): boolean {
  const mapping = TRINITY[ip];
  return mapping.primary !== mapping.secondary && mapping.secondary !== mapping.tertiary;
}
""",
    "packages/test-utils/src/index.ts": """/**
 * @mmup-axis 6 예측·예방
 * @mmup-stage 3 예측
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */

import type { MMUPMeasurement } from "@mmup/types";

export const mockMeasurement: MMUPMeasurement = {
  resourceType: "Observation",
  status: "final",
  code: {
    coding: [{ system: "http://loinc.org", code: "1234-5" }]
  },
  meta: {
    deviceId: "DEV-12345",
    cartridgeUDI: "00123456789012",
    skuLayer: "L1",
    skuId: "SKU-BASIC-01",
    confidence_135P5: 0.95,
    differentialMode: true,
    twinModelVersion: "1.0",
    aiCorrectionApplied: true,
    hashChain: "mock-hash"
  }
};
""",
    ".husky/pre-commit": """#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

pnpm lint || exit 1
pnpm type || exit 1
""",
    "commitlint.config.cjs": """module.exports = {
  extends: ['@commitlint/config-conventional']
};
""",
    "SECURITY.md": """# Security Policy (MMUP Platform)

## Supported Versions
Only the latest version is supported.

## Reporting a Vulnerability
- Send an email to security@mmup.health.
- We implement IEC 81001-5-1 and will respond promptly.
""",
    "CODE_OF_CONDUCT.md": """# Contributor Covenant Code of Conduct

We aim to create an inclusive environment for developing the MMUP Platform. Please respect the Behavioral Ethics and Golden Rules.
""",
    "README.md": """# MPS Manpasik Multi-Measurement Universal POCT Platform (MMUP)

본 저장소는 MMUP 생태계의 모노레포 플랫폼입니다.

## 8축 만파식 생태계
1. 유니버설 측정 (액·기·고)
2. 가정·직장 종합검진 (169 임상검사·44 SKU)
3. 종합 데이터 분석 (패턴·변화·증감)
4. 자가학습 AI 종합병원 (DB + 연합학습)
5. 분야별 전문의 에이전트 (12 도메인)
6. 예측·예방 (Pre-Symptomatic)
7. 5대 처방 프로그램
8. 이행·검증·보상 (Behavioral Engineering)

## 기여 가이드
모든 코드는 `MASTER_SYSTEM_PROMPT.md` 및 `CODING_STANDARDS_V1.md`에 명시된 엄격한 룰과 주석 포맷(`@mmup-axis`, `@mmup-stage`, `@family`)을 준수해야 합니다.
""",
    "packages/ssot/package.json": """{
  "name": "@mmup/ssot",
  "version": "1.0.0",
  "main": "src/index.ts",
  "types": "src/index.ts"
}""",
    "packages/types/package.json": """{
  "name": "@mmup/types",
  "version": "1.0.0",
  "main": "src/index.ts",
  "types": "src/index.ts"
}""",
    "packages/types/src/index.ts": """export * from './loop';\\nexport * from './fhir';\\nexport * from './persona';\\nexport * from './family';""",
    "packages/trinity/package.json": """{
  "name": "@mmup/trinity",
  "version": "1.0.0",
  "main": "src/index.ts",
  "types": "src/index.ts",
  "dependencies": {
    "@mmup/ssot": "workspace:*"
  }
}""",
    "packages/test-utils/package.json": """{
  "name": "@mmup/test-utils",
  "version": "1.0.0",
  "main": "src/index.ts",
  "types": "src/index.ts",
  "dependencies": {
    "@mmup/types": "workspace:*"
  }
}"""
}

for fp, contents in files_to_create.items():
    full_path = os.path.join(base_dir, fp.replace("/", os.sep))
    os.makedirs(os.path.dirname(full_path), exist_ok=True)
    with open(full_path, "w", encoding="utf-8") as f:
        f.write(contents)

print(f"Created {len(files_to_create)} files successfully for Sprint 1 foundation.")
