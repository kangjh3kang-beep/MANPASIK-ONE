#!/usr/bin/env python3
"""
WP-4 전수 수정 스크립트
- 7개 도메인 앱에 루트 layout.tsx 생성
- 각 앱 package.json에 고유 dev 포트 할당
- locale layout에서 중복 html/body 태그 제거
"""
import os, json

ROOT = os.path.expanduser("~/Manpasik")

# 포트 매핑
PORT_MAP = {
    "web": 3000,
    "clinical": 3001,
    "agents-hub": 3002,
    "predictor": 3003,
    "reward": 3004,
    "partner": 3005,
    "gxp": 3006,
    "dev-portal": 3007,
    "app": 3008,
}

# 앱별 메타데이터
APP_META = {
    "web": ("MMUP 만파식 생태계", "MMUP 통합 의료 생태계 플랫폼"),
    "clinical": ("MMUP 임상 데이터 콘솔", "실시간 패킷 검증 및 해시 체인 무결성 대시보드"),
    "agents-hub": ("MMUP AI 에이전트 허브", "의료 AI 모델 오케스트레이션 및 상태 관리"),
    "predictor": ("MMUP 생체 지표 예측", "GNN 기반 다중 질환 시계열 예측 시스템"),
    "reward": ("MMUP 환자 리워드 풀", "개인 데이터 주권 및 토큰 보상"),
    "partner": ("MMUP 파트너 통합 연동", "HL7/FHIR 국제 표준 의료 데이터 파이프라인"),
    "gxp": ("MMUP 의약품 GxP 준수", "cGMP 추적성 확보 및 전자 서명 워크플로우"),
    "dev-portal": ("MMUP 개발자 포털", "오픈 API 명세서 및 실시간 샌드박스"),
    "app": ("MMUP 모바일 하이브리드", "Flutter Native 통신 연동용 PWA Fallback"),
}

fixed_count = 0

for app_name, port in PORT_MAP.items():
    app_dir = os.path.join(ROOT, "apps", app_name)
    if not os.path.isdir(app_dir):
        continue

    # 1. package.json에 포트 설정
    pkg_path = os.path.join(app_dir, "package.json")
    if os.path.exists(pkg_path):
        with open(pkg_path, "r") as f:
            pkg = json.load(f)
        if "scripts" in pkg:
            pkg["scripts"]["dev"] = f"next dev -p {port}"
        with open(pkg_path, "w") as f:
            json.dump(pkg, f, indent=2, ensure_ascii=False)
            f.write("\n")
        fixed_count += 1
        print(f"✅ {app_name}/package.json → dev port {port}")

    # 2. 루트 layout.tsx 생성 (web, clinical은 이미 있음)
    root_layout_path = os.path.join(app_dir, "app", "layout.tsx")
    if not os.path.exists(root_layout_path):
        title, desc = APP_META.get(app_name, ("MMUP", "MMUP Platform"))
        content = f'''import './globals.css';

export const metadata = {{
  title: '{title}',
  description: '{desc}',
}};

export default function RootLayout({{
  children,
}}: {{
  children: React.ReactNode;
}}) {{
  return (
    <html lang="ko" suppressHydrationWarning>
      <body suppressHydrationWarning className="antialiased">
        {{children}}
      </body>
    </html>
  );
}}
'''
        os.makedirs(os.path.dirname(root_layout_path), exist_ok=True)
        with open(root_layout_path, "w") as f:
            f.write(content)
        fixed_count += 1
        print(f"✅ {app_name}/app/layout.tsx → 루트 레이아웃 생성")

    # 3. [locale]/layout.tsx에서 중복 html/body 제거
    locale_layout = os.path.join(app_dir, "app", "[locale]", "layout.tsx")
    if os.path.exists(locale_layout):
        with open(locale_layout, "r") as f:
            code = f.read()
        if "<html" in code and "import '../globals.css'" in code:
            # 이미 html이 있고 globals.css를 import하는 경우 → 수정 필요
            code = code.replace("import '../globals.css';\n", "")
            code = code.replace("import '../globals.css'\n", "")
            # html/body 태그를 제거하고 content만 남기기
            code = code.replace("<html lang={locale} suppressHydrationWarning>", "")
            code = code.replace("<html lang={locale}>", "")
            code = code.replace("      <body suppressHydrationWarning>", "")
            code = code.replace("      <body>", "")
            code = code.replace("      </body>", "")
            code = code.replace("    </html>", "")
            # 들여쓰기 정리 (8공백 → 4공백)
            lines = code.split('\n')
            new_lines = []
            for line in lines:
                if line.startswith('        '):
                    new_lines.append(line[4:])  # remove 4 spaces indent
                else:
                    new_lines.append(line)
            code = '\n'.join(new_lines)
            with open(locale_layout, "w") as f:
                f.write(code)
            fixed_count += 1
            print(f"✅ {app_name}/app/[locale]/layout.tsx → 중복 html/body 제거")

print(f"\n🎯 총 {fixed_count}개 파일 수정 완료!")
