"""Typed, read-only access to the shared nltl values.

    import nltl_values

    nltl_values.color("brand.primary")   # "#DFA4E3"
    nltl_values.font("sans")       # ("Helvetica Neue", ..., "sans-serif")
    nltl_values.shape("face").points     # ((0.0, 0.0), (0.75, 0.25), ...)
"""
from __future__ import annotations

import json
from dataclasses import dataclass
from functools import lru_cache
from importlib import resources
from pathlib import Path
from typing import Any, Dict, Tuple

__version__ = "0.1.0"
__all__ = ["Shape", "color", "colors", "font", "fonts", "shape", "shape_names"]

# (x, y) in screen orientation: origin top-left, y increases downward.
Point = Tuple[float, float]


@dataclass(frozen=True)
class Shape:
    """A named list of points. ``closed`` (default true) means the last point joins the first."""

    name: str
    closed: bool
    points: Tuple[Point, ...]
    description: str = ""


def _data_dir():
    """Packaged data when installed; the repo root when run from a checkout."""
    packaged = resources.files(__name__) / "data"
    if packaged.is_dir():
        return packaged
    return Path(__file__).resolve().parents[2]


def _flatten(group: dict, prefix: str, token_type: str, out: Dict[str, Any]) -> None:
    """Record every token of ``token_type`` in a DTCG group under its dotted path."""
    for key, node in group.items():
        if key.startswith("$") or not isinstance(node, dict):
            continue
        path = f"{prefix}.{key}" if prefix else key
        if "$value" in node:
            if node.get("$type") == token_type:
                out[path] = node["$value"]
        else:
            _flatten(node, path, token_type, out)


def _tokens(file_name: str, token_type: str) -> Dict[str, Any]:
    tree = json.loads((_data_dir() / "tokens" / file_name).read_text("utf-8"))
    out: Dict[str, Any] = {}
    _flatten(tree, "", token_type, out)
    return out


def _resolve(values: Dict[str, str]) -> Dict[str, str]:
    """Replace each alias such as ``"{brand.secondary}"`` with the hex it points at."""
    out: Dict[str, str] = {}
    for path, value in values.items():
        seen = [path]
        while value.startswith("{"):
            target = value[1:-1]
            if target not in values or target in seen:
                raise ValueError(f"nltl-values: bad color alias {value} at {path!r}")
            seen.append(target)
            value = values[target]
        out[path] = value
    return out


@lru_cache(maxsize=None)
def _colors() -> Dict[str, str]:
    return _resolve(_tokens("colors.tokens.json", "color"))


@lru_cache(maxsize=None)
def _fonts() -> Dict[str, Tuple[str, ...]]:
    return {k: tuple(v) for k, v in _tokens("fonts.tokens.json", "fontFamily").items()}


@lru_cache(maxsize=None)
def _shapes() -> Dict[str, Shape]:
    out: Dict[str, Shape] = {}
    for f in (_data_dir() / "shapes").iterdir():
        if not f.name.endswith(".json"):
            continue
        raw = json.loads(f.read_text("utf-8"))
        out[raw["name"]] = Shape(
            name=raw["name"],
            closed=raw.get("closed", True),
            points=tuple((float(x), float(y)) for x, y in raw["points"]),
            description=raw.get("description", ""),
        )
    return out


def color(path: str) -> str:
    """Hex value of the color token at a dotted path, e.g. ``"brand.primary"``.
    Aliases are resolved, so ``color("web.accent")`` gives the hex of ``brand.secondary``."""
    try:
        return _colors()[path]
    except KeyError:
        raise KeyError(f"nltl-values: unknown color {path!r}") from None


def colors() -> Dict[str, str]:
    """Every color token, keyed by dotted path."""
    return dict(_colors())


def font(path: str) -> Tuple[str, ...]:
    """Font stack of the font token at a dotted path, e.g. ``"sans"``: family
    names, most preferred first, ending with a generic family."""
    try:
        return _fonts()[path]
    except KeyError:
        raise KeyError(f"nltl-values: unknown font {path!r}") from None


def fonts() -> Dict[str, Tuple[str, ...]]:
    """Every font token, keyed by dotted path."""
    return dict(_fonts())


def shape(name: str) -> Shape:
    """The shape with the given name."""
    try:
        return _shapes()[name]
    except KeyError:
        raise KeyError(f"nltl-values: unknown shape {name!r}") from None


def shape_names() -> list:
    """Names of all shapes, sorted."""
    return sorted(_shapes())
