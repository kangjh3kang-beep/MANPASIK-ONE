$ErrorActionPreference = "Continue"
subst Z: /D 2>$null
subst Z: "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
Push-Location Z:\

# Clean stale locks
if (Test-Path "Z:\.dart_tool\hooks_runner") {
    Remove-Item -Recurse -Force "Z:\.dart_tool\hooks_runner" 2>$null
    Write-Host "Cleaned hooks_runner locks"
}

$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH
Write-Host "=== flutter build web ==="
flutter.bat build web --release --no-tree-shake-icons 2>&1
Write-Host ""
Write-Host "=== Exit code: $LASTEXITCODE ==="

if (Test-Path "Z:\build\web\index.html") {
    Write-Host "BUILD SUCCESS - Z:\build\web\index.html exists"
} else {
    Write-Host "BUILD FAILED - no index.html found"
}

Pop-Location
subst Z: /D 2>$null
