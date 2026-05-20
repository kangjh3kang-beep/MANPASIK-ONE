#!/usr/bin/env python3
"""Apply a generated design concept to active slim asset files."""

from __future__ import annotations

import argparse
import pathlib
import shutil
import sys

ROOT_DIR = pathlib.Path(__file__).resolve().parents[1]
IMAGE_DIR = ROOT_DIR / "frontend" / "flutter-app" / "assets" / "images"
CONCEPT_ROOT = IMAGE_DIR / "concepts"
FILES = (
    "header_3d_frame_slim.png",
    "bottom_3d_dock_slim.png",
    "btn_3d_action_slim.png",
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Apply one concept set to active image assets.")
    parser.add_argument(
        "--concept",
        required=True,
        help="Concept folder name under assets/images/concepts (e.g. royal_orbit)",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    concept_dir = CONCEPT_ROOT / args.concept
    if not concept_dir.is_dir():
        print(f"[ERROR] Concept folder not found: {concept_dir}")
        return 1

    for name in FILES:
        src = concept_dir / name
        dst = IMAGE_DIR / name
        if not src.exists():
            print(f"[ERROR] Missing concept file: {src}")
            return 1
        shutil.copy2(src, dst)
        print(f"[OK] Applied: {src} -> {dst}")

    print("[DONE] Concept assets applied.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
