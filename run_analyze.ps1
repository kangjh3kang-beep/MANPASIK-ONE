$ErrorActionPreference = "Continue"
subst Z: /D 2>$null
subst Z: "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
Push-Location Z:\
$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH
Write-Host "=== flutter analyze ==="
flutter.bat analyze 2>&1
Pop-Location
subst Z: /D 2>$null
