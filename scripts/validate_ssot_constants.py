#!/usr/bin/env python3
from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def read(path: str) -> str:
    return (ROOT / path).read_text(encoding="utf-8")


def find_float(name: str, text: str) -> float:
    pattern = rf"{re.escape(name)}\s*:\s*f64\s*=\s*([0-9]+(?:\.[0-9]+)?)"
    match = re.search(pattern, text)
    if not match:
        raise AssertionError(f"{name} not found")
    return float(match.group(1))


def find_usize(name: str, text: str) -> int:
    pattern = rf"{re.escape(name)}\s*:\s*usize\s*=\s*([0-9]+)"
    match = re.search(pattern, text)
    if not match:
        raise AssertionError(f"{name} not found")
    return int(match.group(1))


def require(condition: bool, message: str, failures: list[str]) -> None:
    if not condition:
        failures.append(message)


def main() -> int:
    spec = read("ManPaSik_Tech_Spec_v2.4.3.md")
    lib = read("rust-core/manpasik-engine/src/lib.rs")
    agents = read("AGENTS.md")

    expected_alpha = find_float("ALPHA_DEFAULT", spec)
    expected_max_channels = find_usize("FINGERPRINT_DIM_MAX", spec)

    failures: list[str] = []
    require(
        find_float("DEFAULT_ALPHA", lib) == expected_alpha,
        f"rust DEFAULT_ALPHA must equal spec ALPHA_DEFAULT {expected_alpha}",
        failures,
    )
    require(
        find_usize("MAX_CHANNELS", lib) == expected_max_channels,
        f"rust MAX_CHANNELS must equal spec FINGERPRINT_DIM_MAX {expected_max_channels}",
        failures,
    )
    require(
        f"기본값 = {expected_alpha:.2f}" in agents,
        f"AGENTS.md must document alpha default {expected_alpha:.2f}",
        failures,
    )
    require(
        f"MAX_CHANNELS = {expected_max_channels}" in agents,
        f"AGENTS.md must document MAX_CHANNELS = {expected_max_channels}",
        failures,
    )
    require("기본값 = 0.95" not in agents, "AGENTS.md still contains obsolete alpha 0.95", failures)
    require("MAX_CHANNELS = 896" not in agents, "AGENTS.md still contains obsolete MAX_CHANNELS = 896", failures)

    if failures:
        for failure in failures:
            print(f"SSOT_CHECK_FAIL: {failure}", file=sys.stderr)
        return 1

    print(
        "SSOT_CHECK_PASS "
        f"alpha={expected_alpha:.2f} "
        f"max_channels={expected_max_channels}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
