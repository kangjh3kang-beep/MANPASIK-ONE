#!/usr/bin/env python3
"""
WP-4.2 의존성 보강 스크립트
- 7개 도메인 앱에 lucide-react 추가
- tailwind.config.ts에 components 경로 추가
"""
import os, json

ROOT = os.path.expanduser("~/Manpasik")
APPS = ["agents-hub", "predictor", "reward", "partner", "gxp", "dev-portal", "app"]

fixed = 0

for app in APPS:
    # 1. package.json에 lucide-react 추가
    pkg_path = os.path.join(ROOT, "apps", app, "package.json")
    if os.path.exists(pkg_path):
        with open(pkg_path, "r") as f:
            pkg = json.load(f)
        deps = pkg.get("dependencies", {})
        if "lucide-react" not in deps:
            deps["lucide-react"] = "^0.300.0"
            pkg["dependencies"] = deps
            with open(pkg_path, "w") as f:
                json.dump(pkg, f, indent=2, ensure_ascii=False)
                f.write("\n")
            print(f"✅ {app}/package.json → lucide-react 추가")
            fixed += 1

    # 2. tailwind.config.ts에 components 경로 추가
    tw_path = os.path.join(ROOT, "apps", app, "tailwind.config.ts")
    if os.path.exists(tw_path):
        with open(tw_path, "r") as f:
            content = f.read()
        if "./components/**" not in content:
            content = content.replace(
                '"./app/**/*.{js,ts,jsx,tsx,mdx}",',
                '"./app/**/*.{js,ts,jsx,tsx,mdx}",\n    "./components/**/*.{js,ts,jsx,tsx,mdx}",'
            )
            with open(tw_path, "w") as f:
                f.write(content)
            print(f"✅ {app}/tailwind.config.ts → components 경로 추가")
            fixed += 1

print(f"\n🎯 총 {fixed}개 수정 완료!")
