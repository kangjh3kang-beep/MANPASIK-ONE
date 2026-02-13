#!/bin/bash

# Flutter SDK 설치 스크립트 (WSL용)

echo "🚀 Flutter SDK 설치를 시작합니다..."

# 1. 필수 패키지 설치
echo "📦 필수 패키지 설치 중..."
sudo apt-get update
sudo apt-get install -y curl git unzip xz-utils zip libglu1-mesa clang cmake ninja-build pkg-config libgtk-3-dev

# 2. Flutter 다운로드
if [ -d "$HOME/flutter" ]; then
    echo "⚠️ Flutter가 이미 $HOME/flutter 에 존재합니다."
else
    echo "⬇️ Flutter SDK 다운로드 중..."
    git clone https://github.com/flutter/flutter.git -b stable $HOME/flutter
fi

# 3. 환경변수 설정 (.bashrc)
if grep -q "flutter/bin" "$HOME/.bashrc"; then
    echo "✅ PATH 설정이 이미 되어있습니다."
else
    echo "🔧 PATH 환경변수 추가 중..."
    echo '' >> $HOME/.bashrc
    echo '# Flutter SDK Path' >> $HOME/.bashrc
    echo 'export PATH="$PATH:$HOME/flutter/bin"' >> $HOME/.bashrc
    echo "✅ .bashrc 업데이트 완료"
fi

# 4. 현재 세션에 경로 적용
export PATH="$PATH:$HOME/flutter/bin"

# 5. 설치 확인
echo "🔍 설치 확인 중..."
flutter --version

echo "🎉 설치가 완료되었습니다!"
echo "⚠️ 중요: 터미널을 닫았다가 다시 열거나 'source ~/.bashrc'를 실행해주세요."
echo "⚠️ 그 후 'flutter doctor'를 실행하여 상태를 점검하세요."
