$ErrorActionPreference = "Continue"
$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH

# Clean previous build dir and copy fresh
if (Test-Path "D:\manpasik_build") { Remove-Item -Recurse -Force "D:\manpasik_build" }
robocopy "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app" "D:\manpasik_build" /E /NFL /NDL /NJH /NJS /XD .dart_tool build .pub-cache 2>&1 | Out-Null

Push-Location "D:\manpasik_build"
Write-Host "=== flutter pub get ==="
flutter.bat pub get 2>&1
Write-Host ""
Write-Host "=== flutter build web ==="
flutter.bat build web --release --no-tree-shake-icons 2>&1
Write-Host ""
Write-Host "=== Exit code: $LASTEXITCODE ==="

if (Test-Path "D:\manpasik_build\build\web\main.dart.js") {
    $size = (Get-Item "D:\manpasik_build\build\web\main.dart.js").Length
    Write-Host "BUILD SUCCESS - main.dart.js size: $size bytes"
} elseif (Test-Path "D:\manpasik_build\build\web\flutter_bootstrap.js") {
    Write-Host "flutter_bootstrap.js exists - partial build output"
} else {
    Write-Host "No build output found"
}

Pop-Location
