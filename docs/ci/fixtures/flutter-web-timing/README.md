# Flutter Web Timing Fixture Root

This directory holds reproducible timing review fixtures for `flutter-web-timing-report` artifacts.

Each review must follow `docs/ci/flutter-web-timing-review-fixtures.md`:

```text
docs/ci/fixtures/flutter-web-timing/<review-id>/
  manifest.env
  samples/sample-01.env
  samples/sample-02.env
  samples/sample-03.env
  samples/sample-04.env
  samples/sample-05.env
  aggregate.env
```

Only exported GitHub Actions `flutter-web-timing-report` artifacts are allowed here.
Do not store PHI, user identifiers, device identifiers, access tokens, or raw medical measurement data in timing fixtures.
