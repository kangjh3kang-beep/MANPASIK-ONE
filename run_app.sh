#!/bin/bash
# ManPaSik Flutter 앱 실행 스크립트

echo "🚀 만파식(ManPaSik) Flutter 앱을 시작합니다..."

# 절대 경로 사용으로 더 안정적으로 이동
PROJECT_ROOT="$HOME/Manpasik"
FLUTTER_APP_DIR="$PROJECT_ROOT/frontend/flutter-app"

# Flutter SDK 경로 설정 (WSL Native 설치 기준)
FLUTTER_BIN="$HOME/flutter/bin/flutter"

if [ ! -f "$FLUTTER_BIN" ]; then
    echo "⚠️  '$FLUTTER_BIN' 을 찾을 수 없습니다. PATH에서 flutter를 찾습니다."
    if ! command -v flutter &> /dev/null; then
        echo "❌ 에러: flutter 명령어를 찾을 수 없습니다. SDK가 설치되었는지 확인해주세요."
        exit 1
    else
        FLUTTER_BIN="flutter"
    fi
fi

if [ -d "$FLUTTER_APP_DIR" ]; then
    echo "📂 이동: $FLUTTER_APP_DIR"
    cd "$FLUTTER_APP_DIR"
else
    echo "❌ 에러: $FLUTTER_APP_DIR 디렉토리를 찾을 수 없습니다."
    exit 1
fi

if [ ! -f "pubspec.yaml" ]; then
    echo "❌ 에러: pubspec.yaml이 없습니다. 경로를 확인해주세요."
    exit 1
fi

echo "✨ Flutter Linux 데스크톱 앱 빌드 및 실행 중..."
echo "----------------------------------------------------------------"
echo "  [INFO] 앱 빌드를 시작합니다. 잠시만 기다려주세요..."
echo "  [INFO] 초기 빌드는 1~3분 정도 소요될 수 있습니다. 멈춘 것이 아니니 기다려주세요!"
echo "  [INFO] 빌드가 완료되면 GUI 창이 뜹니다."
echo "----------------------------------------------------------------"

"$FLUTTER_BIN" run -d linux

echo "----------------------------------------------------------------"
echo "✅ 앱 실행이 종료되었습니다."
echo "----------------------------------------------------------------"
