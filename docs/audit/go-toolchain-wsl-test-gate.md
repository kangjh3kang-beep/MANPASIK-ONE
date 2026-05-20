# Go Toolchain WSL Test Gate

**Date**: 2026-05-01  
**Owner**: Codex  
**Scope**: WSL Go installation and Measure/Gateway Go quality gate

## Installed Toolchain

- Version: `go1.26.2 linux/amd64`
- Install path: `/home/kangjh3kang/sdk/go-go1.26.2`
- Download source: `https://go.dev/dl/go1.26.2.linux-amd64.tar.gz`
- Archive SHA256: `990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282`

## Invocation

Use a minimal WSL PATH to avoid Windows `Program Files (...)` path parsing issues:

```bash
cd /home/kangjh3kang/Manpasik
export PATH=$HOME/sdk/go-go1.26.2/bin:/usr/local/bin:/usr/bin:/bin
go test ./backend/services/measurement-service/... ./backend/services/gateway/...
```

## Fixed During Gate

- `backend/services/measurement-service/internal/service/measurement.go`
  - Restored `observations` declaration that had been swallowed by a malformed comment line.
- `backend/services/gateway/internal/handler/e2e_test.go`
  - Gateway E2E JSON key assertion now accepts the snake_case proto JSON contract.
  - Added `/api/v1/measurements/process` contract test that verifies REST payload forwarding into gRPC `StreamMeasurement`.

## Verification

- `go test ./backend/services/measurement-service/... ./backend/services/gateway/...`: PASS

