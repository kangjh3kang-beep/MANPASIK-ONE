$ErrorActionPreference = "Continue"
$env:PATH = "D:\flutter_cache\flutter\bin;$env:PATH"

# UNC path를 네트워크 드라이브로 매핑 (Z:)
$uncPath = "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
$driveLetter = "Z:"

# 기존 매핑 제거 후 재매핑
net use $driveLetter /delete 2>$null
subst Z: /d 2>$null
subst Z: $uncPath

Set-Location Z:\
Write-Host "Working Dir: $(Get-Location)"

Write-Host "=== flutter pub get ==="
& flutter.bat pub get

Write-Host "=== flutter run (Chrome) ==="
& flutter.bat run -d chrome --web-port=8080

# 정리
Set-Location C:\
subst Z: /d
