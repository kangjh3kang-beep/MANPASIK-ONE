#!/usr/bin/env python3
"""
WP-4.3 도메인 내비게이션 바 통합 스크립트
- 7개 도메인 앱의 page.tsx에 DomainNav 컴포넌트 추가
"""
import os

ROOT = os.path.expanduser("~/Manpasik")

APPS = {
    "agents-hub": "agents-hub",
    "predictor": "predictor",
    "reward": "reward",
    "partner": "partner",
    "gxp": "gxp",
    "dev-portal": "dev-portal",
    "app": "app",
}

fixed = 0

for app_name, domain_id in APPS.items():
    page_path = os.path.join(ROOT, "apps", app_name, "app", "[locale]", "page.tsx")
    if not os.path.exists(page_path):
        continue

    with open(page_path, "r") as f:
        content = f.read()

    # Skip if already has DomainNav
    if "DomainNav" in content:
        print(f"⏭️  {app_name} → 이미 DomainNav 있음")
        continue

    # Add import
    content = content.replace(
        "'use client';",
        "'use client';\n\nimport { DomainNav } from '@mmup/ui';"
    )

    # Wrap the root div with a fragment that includes DomainNav
    # Find the first return ( and add DomainNav after the opening tag
    content = content.replace(
        "return (\n    <div",
        f"return (\n    <>\n      <DomainNav currentDomain=\"{domain_id}\" />\n      <div"
    )

    # Close the fragment before the last closing paren
    # Find the last </div> in the return and add </>
    # We need to find the closing of the outermost div
    lines = content.split('\n')
    # Find the last line that has just "  );" which closes the return
    for i in range(len(lines) - 1, -1, -1):
        stripped = lines[i].strip()
        if stripped == ");":
            # Insert </> before this
            indent = len(lines[i]) - len(lines[i].lstrip())
            lines.insert(i, " " * (indent + 2) + "</>")
            # Find the matching </div> above and keep as is
            break

    content = '\n'.join(lines)

    with open(page_path, "w") as f:
        f.write(content)

    fixed += 1
    print(f"✅ {app_name}/page.tsx → DomainNav 추가 완료")

# Also update clinical page
clinical_page = os.path.join(ROOT, "apps", "clinical", "app", "[locale]", "page.tsx")
if os.path.exists(clinical_page):
    with open(clinical_page, "r") as f:
        content = f.read()
    if "DomainNav" not in content:
        content = content.replace(
            "import { ClinicalShell }",
            "import { DomainNav } from '@mmup/ui';\nimport { ClinicalShell }"
        )
        content = content.replace(
            "<ClinicalShell>",
            "<>\n      <DomainNav currentDomain=\"clinical\" />\n      <ClinicalShell>"
        )
        # Fix closing
        content = content.replace(
            "</ClinicalShell>\n  );",
            "</ClinicalShell>\n    </>\n  );"
        )
        with open(clinical_page, "w") as f:
            f.write(content)
        fixed += 1
        print(f"✅ clinical/page.tsx → DomainNav 추가 완료")

print(f"\n🎯 총 {fixed}개 페이지에 DomainNav 추가 완료!")
