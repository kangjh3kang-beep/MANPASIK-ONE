$ErrorActionPreference = "Continue"

$wslSource = "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
$buildDir = "D:\manpasik_build"

# Clean previous build dir
if (Test-Path $buildDir) {
    Remove-Item -Recurse -Force $buildDir 2>$null
}

Write-Host "=== Copying project to $buildDir ==="
# Use robocopy for fast copy (exclude build, .dart_tool)
robocopy $wslSource $buildDir /E /XD build .dart_tool .pub-cache windows linux macos ios android /NP /NFL /NDL /NJH /NJS 2>&1 | Out-Null
Write-Host "Copy done"

Push-Location $buildDir
$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH

Write-Host "=== flutter pub get ==="
flutter.bat pub get 2>&1
Write-Host ""

Write-Host "=== flutter build web ==="
flutter.bat build web --release --no-tree-shake-icons 2>&1
Write-Host ""
Write-Host "=== Exit code: $LASTEXITCODE ==="

if (Test-Path "$buildDir\build\web\main.dart.js") {
    Write-Host "BUILD SUCCESS - main.dart.js generated"
    $size = (Get-Item "$buildDir\build\web\main.dart.js").Length
    Write-Host "JS size: $([math]::Round($size/1024/1024, 2)) MB"

    # Copy build output back to WSL
    $wslBuildDir = "$wslSource\build\web"
    if (-not (Test-Path $wslBuildDir)) {
        New-Item -ItemType Directory -Path $wslBuildDir -Force | Out-Null
    }
    robocopy "$buildDir\build\web" $wslBuildDir /E /NP /NFL /NDL /NJH /NJS 2>&1 | Out-Null
    Write-Host "Build output copied back to WSL"
} else {
    Write-Host "BUILD FAILED - main.dart.js not found"
    # List what's in build/web
    if (Test-Path "$buildDir\build\web") {
        Get-ChildItem "$buildDir\build\web" | ForEach-Object { Write-Host $_.Name }
    }
}

Pop-Location
