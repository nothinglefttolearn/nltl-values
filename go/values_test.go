package values_test

import (
	"math"
	"regexp"
	"slices"
	"testing"

	values "github.com/nothinglefttolearn/nltl-values/go"
)

func TestColor(t *testing.T) {
	got, err := values.Color("primary")
	if err != nil {
		t.Fatal(err)
	}
	if got != "#DFA4E3" {
		t.Errorf("primary = %q, want #DFA4E3", got)
	}
	if got, _ := values.Color("secondary"); got != "#879B87" {
		t.Errorf("secondary = %q, want #879B87", got)
	}
	if got, _ := values.Color("bg"); got != "#141414" {
		t.Errorf("bg = %q, want #141414", got)
	}
	if got, _ := values.Color("accent"); got != "#879B87" {
		t.Errorf("accent = %q, want secondary's #879B87", got)
	}
	all, _ := values.Colors()
	hex := regexp.MustCompile(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	for p, v := range all {
		if !hex.MatchString(v) {
			t.Errorf("Colors()[%q] = %q, want a hex color", p, v)
		}
	}
	if _, err := values.Color("nope"); err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestFont(t *testing.T) {
	sans, err := values.Font("sans")
	if err != nil {
		t.Fatal(err)
	}
	if sans[0] != "Helvetica Neue" || sans[len(sans)-1] != "sans-serif" {
		t.Errorf("sans = %v", sans)
	}
	sans[0] = "changed"
	if again, _ := values.Font("sans"); again[0] != "Helvetica Neue" {
		t.Error("Font returned shared backing data")
	}
	if all, _ := values.Fonts(); len(all) != 2 {
		t.Errorf("Fonts() = %v, want sans and mono", all)
	}
	if _, err := values.Font("nope"); err == nil {
		t.Error("expected error for unknown font")
	}
}

func TestShapes(t *testing.T) {
	names, err := values.ShapeNames()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(names, []string{"face", "face2", "face3", "space", "space2", "space3"}) {
		t.Errorf("ShapeNames() = %v, want [face face2 face3 space space2 space3]", names)
	}

	face, err := values.Shape("face")
	if err != nil {
		t.Fatal(err)
	}
	want := []values.Point{{0, 0}, {0.75, 0.25}, {0.75, 1}, {0, 1}}
	if !face.Closed || face.Description == "" || !slices.Equal(face.Points, want) {
		t.Errorf("face = %+v", face)
	}

	// Results are copies: mutating one must not leak into the next call.
	face.Points[0].X = 99
	again, _ := values.Shape("face")
	if again.Points[0].X != 0 {
		t.Error("Shape returned shared backing data")
	}

	if _, err := values.Shape("nope"); err == nil {
		t.Error("expected error for unknown shape")
	}
}

// Each face/space pair tiles the unit square: their areas sum to 1.
func TestFaceAndSpaceFillUnitSquare(t *testing.T) {
	for _, suffix := range []string{"", "2", "3"} {
		face, err := values.Shape("face" + suffix)
		if err != nil {
			t.Fatal(err)
		}
		space, err := values.Shape("space" + suffix)
		if err != nil {
			t.Fatal(err)
		}
		if got := area(face.Points) + area(space.Points); math.Abs(got-1) > 1e-9 {
			t.Errorf("area(face%[1]s)+area(space%[1]s) = %[2]v, want 1", suffix, got)
		}
	}
}

func area(ps []values.Point) float64 {
	var s float64
	for i, p := range ps {
		q := ps[(i+1)%len(ps)]
		s += p.X*q.Y - q.X*p.Y
	}
	return math.Abs(s) / 2
}
