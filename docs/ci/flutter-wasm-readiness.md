# Flutter Wasm Readiness

```yaml
wasm_release_target: false
current_release_target: js_web
wasm_dry_run_blocking: false
```

## Current Status

Flutter Wasm is not a ManPaSik release target yet. The current supported web artifact is the JavaScript web build produced by `flutter build web --no-pub`.

Wasm dry-run findings are tracked by the Flutter web release gate, but they are non-blocking while `wasm_release_target: false`.

## Current Dry-Run Blockers

The current dry-run output reports unsupported web/Wasm imports through:

- `flutter_secure_storage_web`
- `share_plus`
- `connectivity_plus`
- `package:js`

The observed incompatibility classes include `dart:html`, `dart:js`, and JS interop imports that are not valid for the Wasm dry-run target.

## Promotion Criteria

Before changing `wasm_release_target` to `true`, complete a separate compatibility implementation plan that:

- Replaces or conditionally excludes web plugins that require unsupported browser APIs in Wasm.
- Confirms secure storage behavior for web sessions without weakening authentication or PHI protections.
- Re-runs `scripts/flutter_web_release_gate.sh` and removes the `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN` marker from the build output.
- Updates `docs/ci/flutter-web-release-gate-policy.md` so `wasm_dry_run_blocking: true` only after the strict Wasm target passes.
