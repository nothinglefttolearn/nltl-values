Shared static values for nltl repos: color tokens and named shapes (XY point sets), with thin loaders for Go and Python.

This repo is the single source of truth. Change a value here, tag a release, and bump the pin in each consumer. Do not keep local copies in other repos.

**In scope:** static, versioned constants (colors, shapes, and similar).
**Out of scope:** per-environment config such as URLs, secrets, and feature flags.

## What's inside

| Path | Contents |
| --- | --- |
| `tokens/colors.tokens.json` | Colors in [W3C Design Tokens (DTCG)](https://www.designtokens.org/tr/drafts/format/) format |
| `shapes/<name>.json` | One shape per file: `{"name": "arrow", "closed": true, "points": [[x, y], ...]}` |
| `schema/` | JSON Schemas for both |
| `scripts/validate` | Validates everything; CI runs it on every push and PR |
| `go/` | Go package (data embedded with `go:embed`) |
| `python/` | Python package (data shipped as package data) |

### Colors

Tokens can sit at the top level or nest in groups. A token's path is its group names and its own name joined with dots:

```json
{ "primary": { "$value": "#DFA4E3", "$type": "color" }, "text": { "muted": { "$value": "#666666", "$type": "color" } } }
```

`primary` → `#DFA4E3`, `text.muted` → `#666666`. Values are `#rrggbb` or `#rrggbbaa`. `$description` is allowed on tokens and groups.

### Shapes

- `name` uses lowercase letters, digits, and hyphens. It must be unique and match the file name (`shapes/arrow.json` → `"arrow"`).
- `points` are `[x, y]` number pairs, in drawing order.
- `closed: true` means the last point joins back to the first. A closed shape needs at least 3 points and an open one at least 2.
- Coordinates have no implied unit or origin. Each consumer decides how to scale them.

## Consuming

Pin to a release tag. Never depend on `main`.

**Go**

```sh
go get github.com/nltl/nltl-values@v0.1.0
```

```go
import values "github.com/nltl/nltl-values/go"

hex, err := values.Color("primary")   // "#DFA4E3"
arrow, err := values.Shape("arrow")         // arrow.Points[0].X, .Y
all, err := values.Colors()                 // map[path]hex
names, err := values.ShapeNames()
```

The raw files are also available as an `embed.FS` at `nltlvalues.FS` (import `github.com/nltl/nltl-values`).

**Python** (3.9+, no runtime dependencies)

```sh
pip install "nltl-values @ git+https://github.com/nltl/nltl-values@v0.1.0"
```

```python
import nltl_values

nltl_values.color("primary")   # "#DFA4E3"
nltl_values.shape("arrow").points    # ((0.0, 3.0), (6.0, 3.0), ...)
nltl_values.colors()                 # {path: hex}
nltl_values.shape_names()
```

Unknown names raise `KeyError` in Python and return an error in Go.

## Versioning

Releases are git tags `vX.Y.Z` ([semver](https://semver.org)). The tag is the version for every language.

- **Patch:** a value changes without changing meaning, e.g. a color tweak or a shape that was redrawn.
- **Minor:** something is added, e.g. a new token, a new shape, or a new accessor.
- **Major:** something is removed or renamed, or a file format or API changes. Consumer code may break.

### Releasing

1. Merge changes to `main` with CI green.
2. Set `__version__` in `python/nltl_values/__init__.py` to the new version.
3. `git tag vX.Y.Z && git push origin vX.Y.Z`

## Developing

```sh
scripts/validate                         # needs uv; or: pip install jsonschema && python3 scripts/validate
go test ./...
uv venv && uv pip install . && .venv/bin/python -m unittest discover -s python/tests
```

To add a value, edit or add the JSON, run `scripts/validate`, and open a PR. The loaders read every file in `shapes/` automatically, so a new shape needs no code change.
