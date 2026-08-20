package draw

import (
	"fmt"
	"testing"
)

func TestResolveTipMetadata_DiscreteAutoPalette(t *testing.T) {
	fields := []string{"country"}
	tipOrder := []string{"t1", "t2", "t3", "t4"}
	raw := map[string]map[string]string{
		"t1": {"country": "France"},
		"t2": {"country": "Germany"},
		"t3": {"country": "France"},
		"t4": {"country": "Italy"},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	franceColor := resolved["t1"][0]
	if franceColor.Empty {
		t.Fatalf("expected t1/country to be filled")
	}
	if resolved["t3"][0] != franceColor {
		t.Errorf("expected same color for repeated value France, got %v vs %v", resolved["t3"][0], franceColor)
	}
	if resolved["t2"][0] == franceColor {
		t.Errorf("expected different colors for different discrete values")
	}

	wantFrance := mustParseHex(t, defaultPalette[0])
	if franceColor.R != wantFrance.R || franceColor.G != wantFrance.G || franceColor.B != wantFrance.B {
		t.Errorf("expected first-appearance value to get first palette color, got %+v want %+v", franceColor, wantFrance)
	}
	wantGermany := mustParseHex(t, defaultPalette[1])
	if resolved["t2"][0].R != wantGermany.R || resolved["t2"][0].G != wantGermany.G || resolved["t2"][0].B != wantGermany.B {
		t.Errorf("expected second-appearance value to get second palette color, got %+v want %+v", resolved["t2"][0], wantGermany)
	}
}

func TestResolveTipMetadata_ContinuousAutoRange(t *testing.T) {
	fields := []string{"age"}
	tipOrder := []string{"t1", "t2", "t3"}
	raw := map[string]map[string]string{
		"t1": {"age": "0"},
		"t2": {"age": "50"},
		"t3": {"age": "100"},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	low := mustParseHex(t, defaultContinuousLow)
	high := mustParseHex(t, defaultContinuousHigh)

	minColor := resolved["t1"][0]
	if minColor.R != low.R || minColor.G != low.G || minColor.B != low.B {
		t.Errorf("expected min value to resolve to low color, got %+v want %+v", minColor, low)
	}
	maxColor := resolved["t3"][0]
	if maxColor.R != high.R || maxColor.G != high.G || maxColor.B != high.B {
		t.Errorf("expected max value to resolve to high color, got %+v want %+v", maxColor, high)
	}
	midColor := resolved["t2"][0]
	wantMidR := lerpUint8(low.R, high.R, 0.5)
	if midColor.R != wantMidR {
		t.Errorf("expected midpoint value to interpolate, got R=%d want R=%d", midColor.R, wantMidR)
	}
}

func TestResolveTipMetadata_MixedParseableIsDiscrete(t *testing.T) {
	fields := []string{"mixed"}
	tipOrder := []string{"t1", "t2"}
	raw := map[string]map[string]string{
		"t1": {"mixed": "42"},
		"t2": {"mixed": "not-a-number"},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want42 := mustParseHex(t, defaultPalette[0])
	got := resolved["t1"][0]
	if got.R != want42.R || got.G != want42.G || got.B != want42.B {
		t.Errorf("expected mixed-type field to fall back to discrete palette coloring, got %+v want %+v", got, want42)
	}
}

func TestResolveTipMetadata_EmptyCellIsEmptyPlaceholder(t *testing.T) {
	fields := []string{"country"}
	tipOrder := []string{"t1", "t2"}
	raw := map[string]map[string]string{
		"t1": {"country": ""},
		"t2": {"country": "Germany"},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resolved["t1"][0].Empty {
		t.Errorf("expected blank cell to resolve to Empty:true, got %+v", resolved["t1"][0])
	}
	if resolved["t2"][0].Empty {
		t.Errorf("expected non-blank cell to resolve to Empty:false")
	}
}

func TestResolveTipMetadata_TipAbsentFromFileGetsNothing(t *testing.T) {
	fields := []string{"country"}
	tipOrder := []string{"t1"}
	raw := map[string]map[string]string{
		"t1": {"country": "France"},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := resolved["t2"]; ok {
		t.Errorf("expected no entry for a tip not present in tipOrder")
	}
}

func TestResolveTipMetadata_DiscreteYamlOverrideTakesPrecedence(t *testing.T) {
	fields := []string{"country"}
	tipOrder := []string{"t1", "t2"}
	raw := map[string]map[string]string{
		"t1": {"country": "France"},
		"t2": {"country": "Germany"},
	}
	overrides := map[string]FieldColorSpec{
		"country": {
			HasType:  true,
			Type:     FieldDiscrete,
			Discrete: map[string]string{"France": "#123456"},
			Default:  "#abcdef",
		},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, overrides, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantFrance := mustParseHex(t, "#123456")
	got := resolved["t1"][0]
	if got.R != wantFrance.R || got.G != wantFrance.G || got.B != wantFrance.B {
		t.Errorf("expected explicit override color, got %+v want %+v", got, wantFrance)
	}

	wantDefault := mustParseHex(t, "#abcdef")
	got2 := resolved["t2"][0]
	if got2.R != wantDefault.R || got2.G != wantDefault.G || got2.B != wantDefault.B {
		t.Errorf("expected default fallback color for unlisted value, got %+v want %+v", got2, wantDefault)
	}
}

func TestResolveTipMetadata_ContinuousYamlOverride(t *testing.T) {
	fields := []string{"age"}
	tipOrder := []string{"t1", "t2"}
	raw := map[string]map[string]string{
		"t1": {"age": "10"},
		"t2": {"age": "20"},
	}
	min, max := 0.0, 100.0
	overrides := map[string]FieldColorSpec{
		"age": {
			HasType: true,
			Type:    FieldContinuous,
			Low:     "#000000",
			High:    "#ffffff",
			Min:     &min,
			Max:     &max,
		},
	}

	resolved, _, err := ResolveTipMetadata(fields, tipOrder, raw, overrides, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With overridden min=0, max=100, a value of 10 should map to ~10% brightness.
	got := resolved["t1"][0]
	want := lerpUint8(0, 255, 0.1)
	if got.R != want {
		t.Errorf("expected overridden min/max range to shift interpolation, got R=%d want R=%d", got.R, want)
	}
}

func TestResolveTipMetadata_LegendDiscrete(t *testing.T) {
	fields := []string{"country"}
	tipOrder := []string{"t1", "t2", "t3"}
	raw := map[string]map[string]string{
		"t1": {"country": "France"},
		"t2": {"country": "Germany"},
		"t3": {"country": "France"},
	}
	shapes := []Shape{ShapeDiamond}

	_, legend, err := ResolveTipMetadata(fields, tipOrder, raw, nil, shapes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(legend) != 1 {
		t.Fatalf("expected 1 legend entry, got %d", len(legend))
	}
	e := legend[0]
	if e.Field != "country" {
		t.Errorf("expected legend field name %q, got %q", "country", e.Field)
	}
	if e.Shape != ShapeDiamond {
		t.Errorf("expected legend shape to match field shape, got %v", e.Shape)
	}
	if e.Truncated {
		t.Errorf("did not expect truncation for 2 distinct values")
	}
	if len(e.Values) != 2 {
		t.Fatalf("expected 2 distinct values in legend, got %d", len(e.Values))
	}
	if e.Values[0].Label != "France" || e.Values[1].Label != "Germany" {
		t.Errorf("expected legend values in first-appearance order [France, Germany], got %v", e.Values)
	}
	wantFrance := mustParseHex(t, defaultPalette[0])
	if e.Values[0].R != wantFrance.R || e.Values[0].G != wantFrance.G || e.Values[0].B != wantFrance.B {
		t.Errorf("expected legend color to match resolved marker color, got %+v want %+v", e.Values[0], wantFrance)
	}
}

func TestResolveTipMetadata_LegendDiscreteTruncated(t *testing.T) {
	fields := []string{"id"}
	tipOrder := make([]string, 0)
	raw := make(map[string]map[string]string)
	for i := 0; i < maxLegendDiscreteValues+3; i++ {
		tip := fmt.Sprintf("t%d", i)
		tipOrder = append(tipOrder, tip)
		raw[tip] = map[string]string{"id": fmt.Sprintf("v%d", i)}
	}

	_, legend, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !legend[0].Truncated {
		t.Errorf("expected legend to be marked truncated with %d distinct values", maxLegendDiscreteValues+3)
	}
	if len(legend[0].Values) != maxLegendDiscreteValues {
		t.Errorf("expected legend to cap at %d values, got %d", maxLegendDiscreteValues, len(legend[0].Values))
	}
}

func TestResolveTipMetadata_LegendContinuous(t *testing.T) {
	fields := []string{"age"}
	tipOrder := []string{"t1", "t2", "t3"}
	raw := map[string]map[string]string{
		"t1": {"age": "0"},
		"t2": {"age": "50"},
		"t3": {"age": "100"},
	}

	_, legend, err := ResolveTipMetadata(fields, tipOrder, raw, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(legend[0].Values) != 2 {
		t.Fatalf("expected 2 legend values (min, max) for a continuous field, got %d", len(legend[0].Values))
	}
	if legend[0].Values[0].Label != "0" || legend[0].Values[1].Label != "100" {
		t.Errorf("expected legend labels [0, 100], got [%s, %s]", legend[0].Values[0].Label, legend[0].Values[1].Label)
	}
}

func TestResolveFieldShapes_AutoCyclesByColumn(t *testing.T) {
	fields := []string{"a", "b", "c", "d", "e", "f"}
	shapes := ResolveFieldShapes(fields, nil)
	want := []Shape{ShapeCircle, ShapeSquare, ShapeTriangle, ShapeDiamond, ShapeStar, ShapeCircle}
	if len(shapes) != len(want) {
		t.Fatalf("expected %d shapes, got %d", len(want), len(shapes))
	}
	for i := range want {
		if shapes[i] != want[i] {
			t.Errorf("field %d: expected shape %v, got %v", i, want[i], shapes[i])
		}
	}
}

func TestResolveFieldShapes_OverrideTakesPrecedence(t *testing.T) {
	fields := []string{"a", "b"}
	overrides := map[string]FieldColorSpec{
		"a": {HasShape: true, Shape: ShapeStar},
	}
	shapes := ResolveFieldShapes(fields, overrides)
	if shapes[0] != ShapeStar {
		t.Errorf("expected overridden shape ShapeStar for field a, got %v", shapes[0])
	}
	if shapes[1] != ShapeSquare {
		t.Errorf("expected auto-cycled shape ShapeSquare for field b, got %v", shapes[1])
	}
}

func TestParseShapeName(t *testing.T) {
	cases := map[string]Shape{
		"circle":   ShapeCircle,
		"Square":   ShapeSquare,
		"TRIANGLE": ShapeTriangle,
		"diamond":  ShapeDiamond,
		"star":     ShapeStar,
	}
	for name, want := range cases {
		got, err := ParseShapeName(name)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", name, err)
		}
		if got != want {
			t.Errorf("ParseShapeName(%q) = %v, want %v", name, got, want)
		}
	}
	if _, err := ParseShapeName("hexagon"); err == nil {
		t.Errorf("expected error for unknown shape name")
	}
}

func mustParseHex(t *testing.T, s string) TipMetaColor {
	t.Helper()
	r, g, b, a, err := parseHexColor(s)
	if err != nil {
		t.Fatalf("failed to parse hex color %q: %v", s, err)
	}
	return TipMetaColor{R: r, G: g, B: b, A: a}
}
