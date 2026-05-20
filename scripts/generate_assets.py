#!/usr/bin/env python3
"""Generate slim-fit hybrid UI assets via image API.

Priority:
1) OpenAI Images API (requires OPENAI_API_KEY)
2) Optional fallback API (set FALLBACK_IMAGE_API_BASE)

Outputs:
  frontend/flutter-app/assets/images/header_3d_frame_slim.png
  frontend/flutter-app/assets/images/bottom_3d_dock_slim.png
  frontend/flutter-app/assets/images/btn_3d_action_slim.png
"""

from __future__ import annotations

import base64
import json
import os
import pathlib
import sys
import urllib.error
import urllib.parse
import urllib.request

ROOT_DIR = pathlib.Path(__file__).resolve().parents[1]
OUT_DIR = ROOT_DIR / "frontend" / "flutter-app" / "assets" / "images"
SIZE = "1024x1024"

STYLE_LOCK = (
    "Photorealistic premium Korean medical UI hardware family. "
    "Shared materials for all assets: dark brushed titanium body, "
    "thin Sanggam gold trim (#D4AF37), subtle mother-of-pearl (Jagae) cyan iridescence. "
    "Lighting must be consistent across assets: low-key studio, deep-sea blue reflections, "
    "clean cinematic highlights, realistic PBR shading. "
    "Camera must be straight-on orthographic front view, perfectly horizontal, no perspective tilt. "
    "No text, no letters, no numbers, no logo, no watermark, no symbols, no branding."
)

ASSET_SPECS = [
    (
        "header_3d_frame_slim.png",
        (
            f"{STYLE_LOCK} "
            "Create an ultra-thin TOP header bezel strip for a full-width dashboard. "
            "Single centered object, full-width horizontal composition. "
            "The middle area must stay visually open for UI content. "
            "Minimalist, sleek, very thin profile, black background."
        ),
    ),
    (
        "bottom_3d_dock_slim.png",
        (
            f"{STYLE_LOCK} "
            "Create a matching ultra-thin BOTTOM dashboard dock strip from the exact same product line as the header. "
            "Single centered object, full-width horizontal composition. "
            "No stand leg, no camera hole, no center protrusion, no typography. "
            "Include a subtle engraved traditional Korean geometric pattern etched into the surface. "
            "Keep center area clean and minimal for overlays. "
            "Very thin profile, black background."
        ),
    ),
    (
        "btn_3d_action_slim.png",
        (
            f"{STYLE_LOCK} "
            "Create a matching circular action button from the same product line. "
            "Front view, centered. Thin gold bezel with glowing Jagae core. "
            "Embossed gold relief feel for the icon recess, premium tactile machining detail. "
            "Premium tactile look, isolated on black background."
        ),
    ),
]


def _post_json(url: str, payload: dict, headers: dict[str, str]) -> dict:
    request = urllib.request.Request(
        url,
        data=json.dumps(payload).encode("utf-8"),
        method="POST",
        headers=headers,
    )
    with urllib.request.urlopen(request, timeout=180) as response:
        return json.loads(response.read().decode("utf-8"))


def _download_binary(url: str, headers: dict[str, str] | None = None) -> bytes:
    request = urllib.request.Request(url, headers=headers or {})
    with urllib.request.urlopen(request, timeout=180) as response:
        return response.read()


def _generate_with_openai(prompt: str, api_key: str) -> bytes:
    api_url = os.getenv("OPENAI_IMAGE_API_URL", "https://api.openai.com/v1/images/generations")
    model = os.getenv("OPENAI_IMAGE_MODEL", "dall-e-3")

    payload: dict[str, object] = {
        "model": model,
        "prompt": prompt,
        "size": SIZE,
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
        # Some deployments only return URL and do not accept b64_json.
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


def _generate_with_fallback(prompt: str) -> bytes:
    # Example:
    # FALLBACK_IMAGE_API_BASE="https://your-image-api.example/generate?prompt={prompt}&size={size}"
    # The service should return raw image bytes.
    fallback_base = os.getenv("FALLBACK_IMAGE_API_BASE", "").strip()
    if not fallback_base:
        raise RuntimeError("Fallback provider is not configured.")

    url = fallback_base.format(
        prompt=urllib.parse.quote(prompt),
        size=SIZE,
    )
    return _download_binary(url, headers={"User-Agent": "ManPaSik-Asset-Generator/1.0"})


def _generate_image(prompt: str) -> bytes:
    openai_key = os.getenv("OPENAI_API_KEY", "").strip()
    if openai_key:
        return _generate_with_openai(prompt, openai_key)

    # Optional fallback when key is not available.
    return _generate_with_fallback(prompt)


def main() -> int:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    print(f"[INFO] Output directory: {OUT_DIR}")
    print(f"[INFO] Resolution: {SIZE}")

    if not os.getenv("OPENAI_API_KEY", "").strip() and not os.getenv("FALLBACK_IMAGE_API_BASE", "").strip():
        print("[ERROR] OPENAI_API_KEY is not set and FALLBACK_IMAGE_API_BASE is not configured.")
        print("        Set OPENAI_API_KEY or configure FALLBACK_IMAGE_API_BASE, then retry.")
        return 1

    for filename, prompt in ASSET_SPECS:
        target = OUT_DIR / filename
        print(f"[INFO] Generating {filename} ...")
        image_bytes = _generate_image(prompt)
        target.write_bytes(image_bytes)
        print(f"[OK] Saved: {target}")

    print("[DONE] Slim-fit assets generated successfully.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
