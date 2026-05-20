# Flutter Web Build Script
# Maps WSL path to drive letter to avoid UNC path issues

$ErrorActionPreference = "Stop"

# Remove existing mapping if any
subst Z: /D 2>$null

# Map WSL project path to Z: drive
$wslPath = "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
subst Z: $wslPath
Write-Host "Mapped Z: -> $wslPath"

Push-Location Z:\

try {
    # Set Flutter path
    $env:PATH = "D:\flutter_cache\flutter\bin;$env:PATH"

    Write-Host "=== flutter pub get ==="
    flutter.bat pub get 2>&1

    Write-Host ""
    Write-Host "=== flutter build web ==="
    flutter.bat build web --release 2>&1

    Write-Host ""
    Write-Host "=== Build Complete ==="
    Write-Host "Build output: Z:\build\web\"
}
finally {
    Pop-Location
    subst Z: /D 2>$null
}
