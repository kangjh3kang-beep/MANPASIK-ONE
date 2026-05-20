#!/usr/bin/env python3
"""Generate multiple top/bottom/button design concepts for selection."""

from __future__ import annotations

import argparse
import base64
import json
import os
import pathlib
import sys
import urllib.error
import urllib.request

ROOT_DIR = pathlib.Path(__file__).resolve().parents[1]
OUT_DIR = ROOT_DIR / "frontend" / "flutter-app" / "assets" / "images" / "concepts"
DEFAULT_SIZE = "1024x1024"

COMMON_NEGATIVE = (
    "No text, no letters, no numbers, no logo, no watermark, no symbols, no branding, "
    "no UI screenshots, no software windows, no diagrams, no annotations. "
    "No people, no hands, no faces, no animals, no camera lens module, no phone body, "
    "no box product shot, no tabletop scene. Single object only."
)

BASE_STYLE = (
    "Photorealistic premium Korean medical UI hardware, cinematic studio lighting, "
    "deep-sea blue ambient reflection, realistic PBR metal detail."
)

STRICT_GLOBAL = (
    "Strict composition rules: orthographic front view, perfectly horizontal alignment, "
    "exact center placement, no perspective tilt, no rotation, no fisheye, no depth-of-field blur. "
    "The generated object must match the requested part only."
)

STRICT_BY_KIND = {
    "header": (
        "Generate only a top bezel strip component. "
        "Ultra-thin full-width horizontal strip, 8-10% visual thickness of the canvas. "
        "Center area must remain visually open and uncluttered. "
        "No legs, no side blocks, no floating props, no central camera hole."
    ),
    "bottom": (
        "Generate only a bottom dock strip component. "
        "Very thin horizontal dock, 12-15% visual thickness of the canvas. "
        "No support stands, no protruding legs, no deep 3D body blocks, no extra modules."
    ),
    "button": (
        "Generate only one circular action button icon asset. "
        "Centered circular shape, clear silhouette, no rectangular housing, no additional hardware body."
    ),
}

ASSET_TARGETS = (
    ("header_3d_frame_slim.png", "header", "1792x1024"),
    ("bottom_3d_dock_slim.png", "bottom", "1792x1024"),
    ("btn_3d_action_slim.png", "button", DEFAULT_SIZE),
)

CONCEPTS: dict[str, dict[str, str]] = {
    "royal_orbit": {
        "style": (
            f"{BASE_STYLE} Dark brushed titanium + thin Sanggam gold trim (#D4AF37), "
            "clean royal minimalism with precise edge highlights."
        ),
        "header": (
            "Ultra-thin top header bezel strip, full-width horizontal composition, "
            "front orthographic view, center visually open."
        ),
        "bottom": (
            "Matching ultra-thin bottom dock strip, full-width horizontal composition, "
            "subtle engraved Korean geometric pattern, center kept clean for overlay controls."
        ),
        "button": (
            "Matching circular action button, thin gold bezel and cyan jagae glow core, "
            "premium tactile machining, embossed icon recess."
        ),
    },
    "jagae_wave": {
        "style": (
            f"{BASE_STYLE} Deep navy base with iridescent mother-of-pearl (Jagae) wave sheen, "
            "gold trim softened with aqua highlights."
        ),
        "header": (
            "Very thin top frame strip with elegant jagae micro-reflection only on edges, "
            "front view, straight horizontal alignment."
        ),
        "bottom": (
            "Very thin bottom dock strip with flowing jagae wave micro-pattern engraving, "
            "high-end but restrained, no protruding parts."
        ),
        "button": (
            "Round action button with jagae center and slim gold ring, "
            "soft luminous halo, embossed metallic depth."
        ),
    },
    "hanji_gold": {
        "style": (
            f"{BASE_STYLE} Midnight indigo metal blended with subtle Korean Hanji texture, "
            "refined Sanggam gold edge, heritage-modern balance."
        ),
        "header": (
            "Slim top bezel strip with micro Hanji grain and clean gold border line, "
            "front view, symmetrical, transparent feeling center."
        ),
        "bottom": (
            "Slim bottom dock strip with understated traditional lattice engraving inspired by Korean motifs, "
            "no text, no center support leg."
        ),
        "button": (
            "Circular premium action button with gold embossed relief and satin jagae core, "
            "high precision, front view."
        ),
    },
}


def _post_json(url: str, payload: dict, headers: dict[str, str]) -> dict:
    request = urllib.request.Request(
        url,
        data=json.dumps(payload).encode("utf-8"),
        method="POST",
        headers=headers,
    )
    with urllib.request.urlopen(request, timeout=240) as response:
        return json.loads(response.read().decode("utf-8"))


def _download_binary(url: str, headers: dict[str, str] | None = None) -> bytes:
    request = urllib.request.Request(url, headers=headers or {})
    with urllib.request.urlopen(request, timeout=240) as response:
        return response.read()


def _generate_with_openai(prompt: str, api_key: str, size: str) -> bytes:
    api_url = os.getenv("OPENAI_IMAGE_API_URL", "https://api.openai.com/v1/images/generations")
    model = os.getenv("OPENAI_IMAGE_MODEL", "dall-e-3")

    payload: dict[str, object] = {
        "model": model,
        "prompt": prompt,
        "size": size,
        "n": 1,
        "response_format": "b64_json",
    }
    if model == "dall-e-3":
        payload["quality"] = os.getenv("OPENAI_IMAGE_QUALITY", "hd")
        payload["style"] = os.getenv("OPENAI_IMAGE_STYLE", "natural")

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }
    try:
        result = _post_json(api_url, payload, headers)
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", errors="replace")
        if "response_format" in body or "b64_json" in body:
            payload.pop("response_format", None)
            result = _post_json(api_url, payload, headers)
        else:
            raise RuntimeError(f"OpenAI image API HTTP {exc.code}: {body}") from exc

    first = (result.get("data") or [{}])[0]
    b64 = first.get("b64_json")
    if isinstance(b64, str) and b64:
        return base64.b64decode(b64)

    image_url = first.get("url")
    if isinstance(image_url, str) and image_url:
        return _download_binary(image_url)

    raise RuntimeError(f"OpenAI response does not include image payload: {json.dumps(result, ensure_ascii=False)}")


def _generate_image(prompt: str, size: str) -> bytes:
    key = os.getenv("OPENAI_API_KEY", "").strip()
    if not key:
        raise RuntimeError("OPENAI_API_KEY is not set")
    return _generate_with_openai(prompt, key, size)


def _compose_prompt(style: str, body: str, kind: str) -> str:
    strict_by_kind = STRICT_BY_KIND[kind]
    return (
        f"{style} {STRICT_GLOBAL} {strict_by_kind} "
        f"{body} "
        "Isolate object on a clean dark background with generous negative space. "
        f"{COMMON_NEGATIVE}"
    )


def generate_concept(concept: str) -> None:
    spec = CONCEPTS[concept]
    concept_dir = OUT_DIR / concept
    concept_dir.mkdir(parents=True, exist_ok=True)

    print(f"[INFO] Concept: {concept}")
    for filename, kind, size in ASSET_TARGETS:
        prompt = _compose_prompt(spec["style"], spec[kind], kind)
        path = concept_dir / filename
        print(f"[INFO]   generating {filename} ({size})")
        data = _generate_image(prompt, size)
        path.write_bytes(data)
        print(f"[OK]     saved {path}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate design concepts for top/bottom/button assets")
    parser.add_argument(
        "--concept",
        action="append",
        choices=sorted(CONCEPTS.keys()),
        help="Generate only specified concept(s). If omitted, all concepts are generated.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    selected = args.concept or sorted(CONCEPTS.keys())

    if not os.getenv("OPENAI_API_KEY", "").strip():
        print("[ERROR] OPENAI_API_KEY is not set.")
        return 1

    OUT_DIR.mkdir(parents=True, exist_ok=True)
    print(f"[INFO] Output root: {OUT_DIR}")
    print(f"[INFO] Default resolution: {DEFAULT_SIZE}")

    for concept in selected:
        generate_concept(concept)

    print("[DONE] Concept assets generated.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
