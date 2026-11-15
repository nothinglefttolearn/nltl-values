package values_test

import (
	"math"
	"slices"
	"testing"

	values "github.com/nltl/nltl-values/go"
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
	if all, _ := values.Colors(); len(all) != 2 {
		t.Errorf("Colors() = %v, want 2 tokens", all)
	}
	if _, err := values.Color("nope"); err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestShapes(t *testing.T) {
	names, err := values.ShapeNames()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(names, []string{"face", "space"}) {
		t.Errorf("ShapeNames() = %v, want [face space]", names)
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

// face and space tile the unit square: their areas sum to 1.
func TestFaceAndSpaceFillUnitSquare(t *testing.T) {
	face, _ := values.Shape("face")
	space, _ := values.Shape("space")
	if got := area(face.Points) + area(space.Points); math.Abs(got-1) > 1e-9 {
		t.Errorf("area(face)+area(space) = %v, want 1", got)
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
