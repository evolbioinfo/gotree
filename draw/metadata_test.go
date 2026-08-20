package draw

import (
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, nil)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, nil)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, nil)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, nil)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, nil)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, overrides)
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

	resolved, err := ResolveTipMetadata(fields, tipOrder, raw, overrides)
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

func mustParseHex(t *testing.T, s string) TipMetaColor {
	t.Helper()
	r, g, b, a, err := parseHexColor(s)
	if err != nil {
		t.Fatalf("failed to parse hex color %q: %v", s, err)
	}
	return TipMetaColor{R: r, G: g, B: b, A: a}
}
