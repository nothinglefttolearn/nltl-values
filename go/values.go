// Package values exposes typed, read-only access to the shared nltl values.
//
//	import values "github.com/nothinglefttolearn/nltl-values/go"
//
//	hex, err := values.Color("primary") // "#DFA4E3"
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
	shapes map[string]ShapeDef // name -> shape
}

var load = sync.OnceValues(func() (*data, error) {
	d := &data{colors: map[string]string{}, shapes: map[string]ShapeDef{}}

	raw, err := fs.ReadFile(nltlvalues.FS, "tokens/colors.tokens.json")
	if err != nil {
		return nil, err
	}
	var tree map[string]any
	if err := json.Unmarshal(raw, &tree); err != nil {
		return nil, fmt.Errorf("tokens/colors.tokens.json: %w", err)
	}
	flattenColors(tree, "", d.colors)

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

// flattenColors walks a DTCG group, recording every color token under its
// dotted path. Keys starting with "$" are metadata, not children.
func flattenColors(group map[string]any, prefix string, out map[string]string) {
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
			if hex, ok := val.(string); ok && node["$type"] == "color" {
				out[p] = hex
			}
			continue
		}
		flattenColors(node, p, out)
	}
}

// Color returns the hex value of the color token at a dotted path such as
// "primary".
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
