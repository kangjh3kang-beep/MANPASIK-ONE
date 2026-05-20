# 딥링크 매니페스트 (Phase AC-3)

만파식 앱의 딥링크 / Universal Link / App Link 운영 가이드.

## 지원 URL 형식

| 형식 | 라우팅 | 예시 |
|------|--------|------|
| `manpasik://invite/{token}` | InviteAcceptScreen (initialToken=token) | `manpasik://invite/abc123` |
| `https://app.manpasik.com/invite/{token}` | 동일 (Universal/App Link) | `https://app.manpasik.com/invite/abc123` |
| `https://manpasik.com/invite/{token}` | 동일 | — |

## 도메인 호스팅 (필수)

Universal Link / App Link 가 동작하려면 도메인 루트에 다음 파일을 호스팅해야 함.
HTTPS 필수, MIME `application/json`, 응답 코드 200.

### iOS — `.well-known/apple-app-site-association`

`apple-app-site-association.json` 파일 참고. 파일명은 확장자 없는
`apple-app-site-association` 으로 배포. CDN/Nginx 설정 예:

```nginx
location = /.well-known/apple-app-site-association {
    default_type application/json;
    alias /var/www/manpasik/.well-known/apple-app-site-association.json;
}
```

### Android — `.well-known/assetlinks.json`

`assetlinks.json` 파일 참고. 디버그/릴리스 키스토어 SHA-256 fingerprint 를
실제 값으로 교체 필수. 검증 방법:

```bash
# 디버그 키스토어
keytool -list -v -keystore ~/.android/debug.keystore -alias androiddebugkey -storepass android

# 릴리스 키스토어
keytool -list -v -keystore release.keystore -alias <alias>
```

배포 후 검증:

```bash
# Android 검증
adb shell pm verify-app-links --re-verify com.manpasik.app
adb shell pm get-app-links com.manpasik.app

# iOS 검증 (Xcode + 디바이스)
# Settings → Developer → Universal Links → Diagnostics
```

## 앱 라우터 흐름

1. iOS/Android OS 가 URL 을 앱으로 라우팅
2. Flutter `go_router` (app_router.dart) 가 `/invite/:token` 매칭
3. `InviteAcceptScreen(initialToken=...)` 빌드
4. 미로그인이면 인증 우회 (`isInviteDeepLink` redirect skip), 로그인 후 자동 복귀
5. 사용자가 [가입 수락] 탭 → POST /api/v1/tenancy/invitations/accept

## 보안 고려사항

- 토큰은 32 hex chars (16 bytes random) — 추측 불가
- 7일 만료 + revoke 가능 (admin)
- accept 시 동일 사용자 중복 가입 차단 (이미 active 멤버는 거부)
- HTTPS 강제 + verifying domain — App Link autoVerify 가 자동 검증
