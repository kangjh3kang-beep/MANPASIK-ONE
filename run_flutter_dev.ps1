$ErrorActionPreference = "Continue"

# Remove stale subst
subst Z: /D 2>$null

# Map fresh
subst Z: "\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app"
Push-Location Z:\

$env:PATH = "D:\flutter_cache\flutter\bin;" + $env:PATH

# Remove stale hooks_runner locks
if (Test-Path "Z:\.dart_tool\hooks_runner") {
    Remove-Item -Recurse -Force "Z:\.dart_tool\hooks_runner" 2>$null
    Write-Host "Cleaned hooks_runner"
}

# Pre-create the lock directory and file to avoid the timeout
$lockDir = "Z:\.dart_tool\hooks_runner\shared\objective_c"
if (-not (Test-Path $lockDir)) {
    New-Item -ItemType Directory -Path $lockDir -Force | Out-Null
}
# Create empty lock file so the framework doesn't block
"" | Out-File -FilePath "$lockDir\.lock" -Force -NoNewline

Write-Host "=== flutter build web (with pre-created lock) ==="
flutter.bat build web --release --no-tree-shake-icons 2>&1
Write-Host ""
Write-Host "=== Exit code: $LASTEXITCODE ==="

if (Test-Path "Z:\build\web\main.dart.js") {
    Write-Host "BUILD SUCCESS - main.dart.js generated"
    $size = (Get-Item "Z:\build\web\main.dart.js").Length
    Write-Host "JS size: $size bytes"
} else {
    Write-Host "main.dart.js not found, checking flutter_bootstrap.js..."
    if (Test-Path "Z:\build\web\flutter_bootstrap.js") {
        Write-Host "flutter_bootstrap.js exists - partial build output"
    }
}

Pop-Location
subst Z: /D 2>$null
