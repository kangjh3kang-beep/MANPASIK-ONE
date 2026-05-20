# Flutter Web Release Gate Policy

```yaml
release_target: js_web
required_on_pull_request: true
release_branch_required: true
ci_execution_mode: pull_request_and_release_branch
nightly_split: false
cost_review_required_before_relaxing: true
blocking_build_command: flutter build web --no-pub
timing_capture: true
timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
wasm_dry_run_blocking: false
gate_script: scripts/flutter_web_release_gate.sh
```

## Policy

The current release target for the Flutter client is the JavaScript web artifact produced by:

```bash
flutter build web --no-pub
```

This build must pass on pull requests and release branches. A missing `build/web` success marker is release blocking.

The gate remains PR-required even though it has non-trivial CI cost, because JS web is the current release artifact. Moving this build to nightly-only or release-branch-only requires a measured CI cost review and an explicit policy update.

Each gate run must emit `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` so cost review decisions can be based on measured build timing rather than assumption.

Flutter Wasm dry-run compatibility warnings are tracked by `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`, but they are not release blocking while `release_target` remains `js_web`.

## When This Changes

If the JS web build becomes too expensive for pull requests, collect timing data and create a separate CI cost policy plan before setting `nightly_split: true`.

If Wasm becomes a product release target, create a separate strict Wasm compatibility plan before changing this policy. That plan must address the currently reported web dependencies that import unsupported `dart:html`, `dart:js`, or JS interop APIs in the Wasm dry run.
