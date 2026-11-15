package values_test

import (
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

func TestShape(t *testing.T) {
	arrow, err := values.Shape("arrow")
	if err != nil {
		t.Fatal(err)
	}
	if arrow.Name != "arrow" || !arrow.Closed {
		t.Errorf("got name=%q closed=%v", arrow.Name, arrow.Closed)
	}
	if len(arrow.Points) != 7 || arrow.Points[0] != (values.Point{X: 0, Y: 3}) {
		t.Errorf("unexpected points: %v", arrow.Points)
	}

	// Results are copies: mutating one must not leak into the next call.
	arrow.Points[0].X = 99
	again, _ := values.Shape("arrow")
	if again.Points[0].X != 0 {
		t.Error("Shape returned shared backing data")
	}

	names, err := values.ShapeNames()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 || names[0] != "arrow" {
		t.Errorf("ShapeNames() = %v", names)
	}
	if _, err := values.Shape("nope"); err == nil {
		t.Error("expected error for unknown shape")
	}
}
