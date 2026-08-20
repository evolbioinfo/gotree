package draw

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldType describes how a metadata field's values are interpreted for coloring.
type FieldType int

const (
	FieldDiscrete FieldType = iota
	FieldContinuous
)

// Shape is the marker shape drawn for a whole metadata field (column).
// Color still varies per value/tip; shape is constant across a field, so it
// helps tell columns apart even without a legend.
type Shape int

const (
	ShapeCircle Shape = iota
	ShapeSquare
	ShapeTriangle
	ShapeDiamond
	ShapeStar
)

// defaultShapeCycle is used to auto-assign a distinct shape to each
// metadata field, in column order, when not overridden.
var defaultShapeCycle = []Shape{ShapeCircle, ShapeSquare, ShapeTriangle, ShapeDiamond, ShapeStar}

// ParseShapeName parses a shape name ("circle", "square", "triangle",
// "diamond", "star", case-insensitive) as used in --metadata-colors YAML.
func ParseShapeName(s string) (Shape, error) {
	switch strings.ToLower(s) {
	case "circle":
		return ShapeCircle, nil
	case "square":
		return ShapeSquare, nil
	case "triangle":
		return ShapeTriangle, nil
	case "diamond":
		return ShapeDiamond, nil
	case "star":
		return ShapeStar, nil
	default:
		return ShapeCircle, fmt.Errorf("unknown shape %q (must be one of: circle, square, triangle, diamond, star)", s)
	}
}

// ResolveFieldShapes returns the shape to draw for each field, in the same
// order as fields: the override's shape if given, otherwise a shape cycled
// from a fixed default sequence by column position.
func ResolveFieldShapes(fields []string, overrides map[string]FieldColorSpec) []Shape {
	shapes := make([]Shape, len(fields))
	for i, fieldName := range fields {
		if override, ok := overrides[fieldName]; ok && override.HasShape {
			shapes[i] = override.Shape
		} else {
			shapes[i] = defaultShapeCycle[i%len(defaultShapeCycle)]
		}
	}
	return shapes
}

// FieldColorSpec is an optional per-field color scheme override (typically
// loaded from a YAML file). Zero value means "fully auto-detect".
type FieldColorSpec struct {
	HasType  bool
	Type     FieldType
	Discrete map[string]string // value -> hex color
	Default  string            // hex color for discrete values not in Discrete
	Low      string            // hex color for continuous minimum
	High     string            // hex color for continuous maximum
	Min      *float64
	Max      *float64
	HasShape bool
	Shape    Shape
}

// TipMetaColor is the resolved color for one (tip, field) cell.
// Empty is true when the tip was present in the metadata file but the cell
// for this field was blank: it should be rendered as an unfilled,
// grey-bordered placeholder circle, regardless of R/G/B/A.
type TipMetaColor struct {
	Empty      bool
	R, G, B, A uint8
}

const (
	defaultContinuousLow  = "#2c7bb6"
	defaultContinuousHigh = "#d7191c"
)

// defaultPalette is a categorical color palette (matplotlib "tab20"),
// cycled in first-appearance order for auto-colored discrete fields.
var defaultPalette = []string{
	"#1f77b4", "#ff7f0e", "#2ca02c", "#d62728", "#9467bd",
	"#8c564b", "#e377c2", "#7f7f7f", "#bcbd22", "#17becf",
	"#aec7e8", "#ffbb78", "#98df8a", "#ff9896", "#c5b0d5",
	"#c49c94", "#f7b6d2", "#c7c7c7", "#dbdb8d", "#9edae5",
}

// ResolveTipMetadata computes, for each tip in tipOrder and each field in
// fields, the color that should be drawn for that (tip, field) cell.
//
// raw is indexed as raw[tipName][fieldName] = rawValue ("" meaning blank).
// overrides is indexed by field name and may be nil or partial: any field
// (or attribute of a field) not present is auto-detected/auto-colored.
func ResolveTipMetadata(fields []string, tipOrder []string, raw map[string]map[string]string, overrides map[string]FieldColorSpec) (map[string][]TipMetaColor, error) {
	result := make(map[string][]TipMetaColor, len(tipOrder))
	for _, tip := range tipOrder {
		result[tip] = make([]TipMetaColor, len(fields))
	}

	for fi, fieldName := range fields {
		override := overrides[fieldName]

		distinctValues := make([]string, 0)
		seen := make(map[string]bool)
		for _, tip := range tipOrder {
			v := raw[tip][fieldName]
			if v == "" {
				continue
			}
			if !seen[v] {
				seen[v] = true
				distinctValues = append(distinctValues, v)
			}
		}

		fieldType := override.Type
		if !override.HasType {
			fieldType = FieldDiscrete
			if len(distinctValues) > 0 {
				allNumeric := true
				for _, v := range distinctValues {
					if _, err := strconv.ParseFloat(v, 64); err != nil {
						allNumeric = false
						break
					}
				}
				if allNumeric {
					fieldType = FieldContinuous
				}
			}
		}

		if fieldType == FieldContinuous {
			if err := resolveContinuousField(fieldName, fi, distinctValues, override, tipOrder, raw, result); err != nil {
				return nil, err
			}
		} else {
			if err := resolveDiscreteField(fieldName, fi, distinctValues, override, tipOrder, raw, result); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func resolveContinuousField(fieldName string, fi int, distinctValues []string, override FieldColorSpec, tipOrder []string, raw map[string]map[string]string, result map[string][]TipMetaColor) error {
	low := defaultContinuousLow
	if override.Low != "" {
		low = override.Low
	}
	high := defaultContinuousHigh
	if override.High != "" {
		high = override.High
	}
	lr, lg, lb, la, err := parseHexColor(low)
	if err != nil {
		return fmt.Errorf("field %q: invalid low color %q: %w", fieldName, low, err)
	}
	hr, hg, hb, ha, err := parseHexColor(high)
	if err != nil {
		return fmt.Errorf("field %q: invalid high color %q: %w", fieldName, high, err)
	}

	min, max := 0.0, 0.0
	if override.Min != nil {
		min = *override.Min
	} else if len(distinctValues) > 0 {
		min, _ = strconv.ParseFloat(distinctValues[0], 64)
	}
	if override.Max != nil {
		max = *override.Max
	} else if len(distinctValues) > 0 {
		max = min
	}
	if override.Min == nil || override.Max == nil {
		for _, v := range distinctValues {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			if override.Min == nil && f < min {
				min = f
			}
			if override.Max == nil && f > max {
				max = f
			}
		}
	}

	for _, tip := range tipOrder {
		v := raw[tip][fieldName]
		if v == "" {
			result[tip][fi] = TipMetaColor{Empty: true}
			continue
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			result[tip][fi] = TipMetaColor{Empty: true}
			continue
		}
		t := 0.5
		if max > min {
			t = (f - min) / (max - min)
		}
		if t < 0 {
			t = 0
		} else if t > 1 {
			t = 1
		}
		result[tip][fi] = TipMetaColor{
			Empty: false,
			R:     lerpUint8(lr, hr, t),
			G:     lerpUint8(lg, hg, t),
			B:     lerpUint8(lb, hb, t),
			A:     lerpUint8(la, ha, t),
		}
	}
	return nil
}

func resolveDiscreteField(fieldName string, fi int, distinctValues []string, override FieldColorSpec, tipOrder []string, raw map[string]map[string]string, result map[string][]TipMetaColor) error {
	type rgba struct{ r, g, b, a uint8 }
	colorFor := make(map[string]rgba, len(distinctValues))
	paletteIdx := 0

	for _, v := range distinctValues {
		var hex string
		switch {
		case override.Discrete != nil && override.Discrete[v] != "":
			hex = override.Discrete[v]
		case override.Default != "":
			hex = override.Default
		default:
			hex = defaultPalette[paletteIdx%len(defaultPalette)]
			paletteIdx++
		}
		r, g, b, a, err := parseHexColor(hex)
		if err != nil {
			return fmt.Errorf("field %q: invalid color %q for value %q: %w", fieldName, hex, v, err)
		}
		colorFor[v] = rgba{r, g, b, a}
	}

	for _, tip := range tipOrder {
		v := raw[tip][fieldName]
		if v == "" {
			result[tip][fi] = TipMetaColor{Empty: true}
			continue
		}
		c := colorFor[v]
		result[tip][fi] = TipMetaColor{Empty: false, R: c.r, G: c.g, B: c.b, A: c.a}
	}
	return nil
}

func lerpUint8(a, b uint8, t float64) uint8 {
	v := float64(a) + t*(float64(b)-float64(a))
	if v < 0 {
		v = 0
	} else if v > 255 {
		v = 255
	}
	return uint8(v + 0.5)
}

// parseHexColor parses a "#rrggbb" or "#rrggbbaa" hex color string.
func parseHexColor(s string) (r, g, b, a uint8, err error) {
	a = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	case 9:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	default:
		err = fmt.Errorf("invalid length (%d) for hex color %q, must be 7 or 9", len(s), s)
	}
	return
}
