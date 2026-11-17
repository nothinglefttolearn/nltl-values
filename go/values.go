// Package values exposes typed, read-only access to the shared nltl values.
//
//	import values "github.com/nothinglefttolearn/nltl-values/go"
//
//	hex, err := values.Color("primary") // "#DFA4E3"
//	stack, err := values.Font("sans")   // ["Helvetica Neue", ..., "sans-serif"]
//	face, err := values.Shape("face")
//
// Every accessor returns a fresh copy, so callers may modify results freely.
package values

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"

	nltlvalues "github.com/nothinglefttolearn/nltl-values"
)

// Point is one XY coordinate of a shape, in screen orientation: the origin
// is top-left and Y increases downward.
type Point struct {
	X, Y float64
}

// UnmarshalJSON decodes a point from its JSON form, a [x, y] pair.
func (p *Point) UnmarshalJSON(b []byte) error {
	var xy []float64
	if err := json.Unmarshal(b, &xy); err != nil {
		return err
	}
	if len(xy) != 2 {
		return fmt.Errorf("point must have exactly 2 coordinates, got %d", len(xy))
	}
	p.X, p.Y = xy[0], xy[1]
	return nil
}

// ShapeDef is a named list of points. Closed means the last point connects
// back to the first; it is true unless the file says otherwise.
type ShapeDef struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Closed      bool    `json:"closed"`
	Points      []Point `json:"points"`
}

func (s ShapeDef) clone() ShapeDef {
	s.Points = append([]Point(nil), s.Points...)
	return s
}

type data struct {
	colors map[string]string   // dotted token path -> hex
	fonts  map[string][]string // dotted token path -> font stack
	shapes map[string]ShapeDef // name -> shape
}

var load = sync.OnceValues(func() (*data, error) {
	d := &data{colors: map[string]string{}, fonts: map[string][]string{}, shapes: map[string]ShapeDef{}}

	colorTokens, err := readTokens("tokens/colors.tokens.json", "color")
	if err != nil {
		return nil, err
	}
	for p, v := range colorTokens {
		hex, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("tokens/colors.tokens.json: %s: value is not a string", p)
		}
		d.colors[p] = hex
	}
	if err := resolveAliases(d.colors); err != nil {
		return nil, fmt.Errorf("tokens/colors.tokens.json: %w", err)
	}

	fontTokens, err := readTokens("tokens/fonts.tokens.json", "fontFamily")
	if err != nil {
		return nil, err
	}
	for p, v := range fontTokens {
		list, _ := v.([]any)
		stack := make([]string, 0, len(list))
		for _, f := range list {
			name, ok := f.(string)
			if !ok {
				return nil, fmt.Errorf("tokens/fonts.tokens.json: %s: family name is not a string", p)
			}
			stack = append(stack, name)
		}
		d.fonts[p] = stack
	}

	files, err := fs.Glob(nltlvalues.FS, "shapes/*.json")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		raw, err := fs.ReadFile(nltlvalues.FS, f)
		if err != nil {
			return nil, err
		}
		s := ShapeDef{Closed: true} // default when "closed" is absent
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, fmt.Errorf("%s: %w", path.Base(f), err)
		}
		d.shapes[s.Name] = s
	}
	return d, nil
})

// readTokens decodes a DTCG token file and returns the $value of every token
// of the given $type, keyed by dotted path.
func readTokens(name, tokenType string) (map[string]any, error) {
	raw, err := fs.ReadFile(nltlvalues.FS, name)
	if err != nil {
		return nil, err
	}
	var tree map[string]any
	if err := json.Unmarshal(raw, &tree); err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	out := map[string]any{}
	flatten(tree, "", tokenType, out)
	return out, nil
}

// flatten walks a DTCG group, recording the value of every token of the
// given type under its dotted path. Keys starting with "$" are metadata, not
// children.
func flatten(group map[string]any, prefix, tokenType string, out map[string]any) {
	for k, v := range group {
		if strings.HasPrefix(k, "$") {
			continue
		}
		node, ok := v.(map[string]any)
		if !ok {
			continue
		}
		p := k
		if prefix != "" {
			p = prefix + "." + k
		}
		if val, isToken := node["$value"]; isToken {
			if node["$type"] == tokenType {
				out[p] = val
			}
			continue
		}
		flatten(node, p, tokenType, out)
	}
}

// resolveAliases replaces, in place, each alias such as "{secondary}" with
// the hex value of the token it points at.
func resolveAliases(colors map[string]string) error {
	resolved := make(map[string]string, len(colors))
	for p, v := range colors {
		seen := map[string]bool{p: true}
		for strings.HasPrefix(v, "{") {
			target := strings.TrimSuffix(strings.TrimPrefix(v, "{"), "}")
			next, ok := colors[target]
			if !ok || seen[target] {
				return fmt.Errorf("bad color alias %s at %q", v, p)
			}
			seen[target] = true
			v = next
		}
		resolved[p] = v
	}
	for p, v := range resolved {
		colors[p] = v
	}
	return nil
}

// Color returns the hex value of the color token at a dotted path such as
// "primary". Aliases are resolved, so Color("accent") gives the hex of
// "secondary".
func Color(path string) (string, error) {
	d, err := load()
	if err != nil {
		return "", err
	}
	hex, ok := d.colors[path]
	if !ok {
		return "", fmt.Errorf("nltl-values: unknown color %q", path)
	}
	return hex, nil
}

// Colors returns every color token, keyed by dotted path.
func Colors() (map[string]string, error) {
	d, err := load()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(d.colors))
	for k, v := range d.colors {
		out[k] = v
	}
	return out, nil
}

// Font returns the font stack of the font token at a dotted path such as
// "sans": family names, most preferred first, ending with a generic family.
func Font(path string) ([]string, error) {
	d, err := load()
	if err != nil {
		return nil, err
	}
	stack, ok := d.fonts[path]
	if !ok {
		return nil, fmt.Errorf("nltl-values: unknown font %q", path)
	}
	return append([]string(nil), stack...), nil
}

// Fonts returns every font token, keyed by dotted path.
func Fonts() (map[string][]string, error) {
	d, err := load()
	if err != nil {
		return nil, err
	}
	out := make(map[string][]string, len(d.fonts))
	for k, v := range d.fonts {
		out[k] = append([]string(nil), v...)
	}
	return out, nil
}

// Shape returns the shape with the given name.
func Shape(name string) (ShapeDef, error) {
	d, err := load()
	if err != nil {
		return ShapeDef{}, err
	}
	s, ok := d.shapes[name]
	if !ok {
		return ShapeDef{}, fmt.Errorf("nltl-values: unknown shape %q", name)
	}
	return s.clone(), nil
}

// ShapeNames returns the names of all shapes, sorted.
func ShapeNames() ([]string, error) {
	d, err := load()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(d.shapes))
	for n := range d.shapes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}
