# ManPaSik Flutter App

만파식(萬波息) 헬스케어 생태계의 모바일/데스크톱 클라이언트 애플리케이션입니다.

## 🚀 시작하기

### 1. 전제 조건
- Flutter SDK 3.22.0 이상
- Rust Toolchain (Core 연동 시 필요)

### 2. 설정 및 실행

```bash
# 1. 의존성 설치
flutter pub get

# 2. 코드 생성 (Riverpod, JSON Serializable 등)
dart run build_runner build -d

# 3. 앱 실행
flutter run
```

## 📁 프로젝트 구조

- `lib/core`: 라우터, 테마, 상수, 유틸리티
- `lib/features`: 기능별 모듈 (auth, home, measurement 등)
- `lib/shared`: 공통 위젯, 모델, 프로바이더
- `assets`: 이미지, 폰트

## 🛠 기술 스택

- **UI Framework**: Flutter 3
- **State Management**: Riverpod 2 (Code Generation)
- **Routing**: GoRouter
- **Networking**: Dio + Retrofit
- **Local DB**: SQLite + Hive
- **FFI**: flutter_rust_bridge (Rust Core 연동)
