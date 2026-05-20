$ErrorActionPreference = "Continue"

# Remove stale subst
subst Z: /D 2>$null

# Map fresh
subst Z: "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
Push-Location Z:\

$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH

# Full clean
Write-Host "=== flutter clean ==="
flutter.bat clean 2>&1
Write-Host ""

# Remove all lock/cache artifacts
if (Test-Path "Z:\.dart_tool") {
    Remove-Item -Recurse -Force "Z:\.dart_tool" 2>$null
    Write-Host "Removed .dart_tool"
}

# Pub get fresh
Write-Host "=== flutter pub get ==="
flutter.bat pub get 2>&1
Write-Host ""

# Build web
Write-Host "=== flutter build web ==="
flutter.bat build web --release --no-tree-shake-icons 2>&1
Write-Host ""
Write-Host "=== Exit code: $LASTEXITCODE ==="

if (Test-Path "Z:\build\web\main.dart.js") {
    Write-Host "BUILD SUCCESS - main.dart.js generated"
    $size = (Get-Item "Z:\build\web\main.dart.js").Length
    Write-Host "JS size: $size bytes"
} else {
    Write-Host "BUILD FAILED - main.dart.js not found"
    Write-Host "Trying alternative: flutter run -d chrome"
}

Pop-Location
subst Z: /D 2>$null
