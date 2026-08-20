package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/evolbioinfo/gotree/draw"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var drawNoTipLabels bool
var drawNoBranchLengths bool
var drawInternalNodeLabels bool
var drawSupport bool
var drawSupportCutoff float64
var drawInternalNodeSymbols bool
var drawNodeComment bool
var metadataFile string
var metadataColorsFile string

// drawCmd represents the draw command
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Draw trees",
	Long:  `Draw trees `,
}

func init() {
	RootCmd.AddCommand(drawCmd)

	drawCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	drawCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	drawCmd.PersistentFlags().BoolVar(&drawNoTipLabels, "no-tip-labels", false, "Draw the tree without tip labels")
	drawCmd.PersistentFlags().BoolVar(&drawNoBranchLengths, "no-branch-lengths", false, "Draw the tree without branch lengths (all the same length)")
	drawCmd.PersistentFlags().BoolVar(&drawInternalNodeLabels, "with-node-labels", false, "Draw the tree with internal node labels")
	drawCmd.PersistentFlags().BoolVar(&drawInternalNodeSymbols, "with-node-symbols", false, "Draw the tree with internal node symbols")
	drawCmd.PersistentFlags().BoolVar(&drawSupport, "with-branch-support", false, "Highlight highly supported branches")
	drawCmd.PersistentFlags().Float64Var(&drawSupportCutoff, "support-cutoff", 0.7, "Cutoff for highlithing supported branches")
	drawCmd.PersistentFlags().BoolVar(&drawNodeComment, "with-node-comments", false, "Draw the tree with internal node comments (if --with-node-labels is not set)")
	drawCmd.PersistentFlags().StringVarP(&metadataFile, "metadata-file", "m", "", "Tab separated metadata file to add colored circles to tip nodes (svg & png): tip name in the first column (header ignored), then one column per metadata field (header = field name). Values are auto-detected as discrete or continuous; colors are auto-assigned unless overridden with --metadata-colors. Empty cells draw an unfilled grey circle.")
	drawCmd.PersistentFlags().StringVar(&metadataColorsFile, "metadata-colors", "", "Optional YAML file overriding the color scheme and/or marker shape of one or more --metadata-file fields (discrete value->color map, or continuous low/high/min/max; shape: circle|square|triangle|diamond|star)")
}

// parseMetadataTSV reads a tab separated metadata file: tip name in the
// first column (header ignored), then one column per metadata field
// (header = field name). Returns field names in column order, tip names
// in row order (for deterministic first-appearance color assignment), and
// values indexed as values[tipName][fieldName].
func parseMetadataTSV(filepath string) (fields []string, tipOrder []string, values map[string]map[string]string, err error) {
	var file *os.File
	if file, err = os.Open(filepath); err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	var header []string
	if header, err = reader.Read(); err != nil {
		return
	}
	if len(header) < 2 {
		err = fmt.Errorf("metadata file must have at least 2 columns (tip name + 1 metadata field), got %d", len(header))
		return
	}
	fields = header[1:]

	values = make(map[string]map[string]string)
	tipOrder = make([]string, 0)

	for {
		var record []string
		record, err = reader.Read()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if len(record) != len(header) {
			err = fmt.Errorf("metadata file: row for tip %q has %d columns, expecting %d", record[0], len(record), len(header))
			return
		}
		tip := record[0]
		if _, ok := values[tip]; !ok {
			tipOrder = append(tipOrder, tip)
			values[tip] = make(map[string]string)
		}
		for i, field := range fields {
			values[tip][field] = record[i+1]
		}
	}

	return
}

// yamlFieldSpec mirrors the per-field entries of a --metadata-colors YAML file.
type yamlFieldSpec struct {
	Type    string            `yaml:"type"`
	Colors  map[string]string `yaml:"colors"`
	Default string            `yaml:"default"`
	Low     string            `yaml:"low"`
	High    string            `yaml:"high"`
	Min     *float64          `yaml:"min"`
	Max     *float64          `yaml:"max"`
	Shape   string            `yaml:"shape"`
}

// parseMetadataColorsYAML reads an optional per-field color scheme override file.
func parseMetadataColorsYAML(filepath string) (specs map[string]draw.FieldColorSpec, err error) {
	var content []byte
	if content, err = os.ReadFile(filepath); err != nil {
		return
	}

	yamlSpecs := make(map[string]yamlFieldSpec)
	if err = yaml.Unmarshal(content, &yamlSpecs); err != nil {
		return
	}

	specs = make(map[string]draw.FieldColorSpec, len(yamlSpecs))
	for field, y := range yamlSpecs {
		spec := draw.FieldColorSpec{
			Discrete: y.Colors,
			Default:  y.Default,
			Low:      y.Low,
			High:     y.High,
			Min:      y.Min,
			Max:      y.Max,
		}
		switch y.Type {
		case "discrete":
			spec.HasType = true
			spec.Type = draw.FieldDiscrete
		case "continuous":
			spec.HasType = true
			spec.Type = draw.FieldContinuous
		case "":
			// auto-detected
		default:
			err = fmt.Errorf("metadata-colors file: field %q has invalid type %q (must be \"discrete\" or \"continuous\")", field, y.Type)
			return
		}
		if y.Shape != "" {
			if spec.Shape, err = draw.ParseShapeName(y.Shape); err != nil {
				err = fmt.Errorf("metadata-colors file: field %q: %w", field, err)
				return
			}
			spec.HasShape = true
		}
		specs[field] = spec
	}
	return
}

// loadTipMetadata reads --metadata-file (and optional --metadata-colors),
// and resolves per-tip, per-field marker colors and per-field marker
// shapes. Returns (nil, nil, nil, nil) when --metadata-file is not set.
func loadTipMetadata() (fields []string, shapes []draw.Shape, values map[string][]draw.TipMetaColor, err error) {
	if metadataFile == "" {
		return nil, nil, nil, nil
	}

	var tipOrder []string
	var raw map[string]map[string]string
	if fields, tipOrder, raw, err = parseMetadataTSV(metadataFile); err != nil {
		return
	}

	overrides := map[string]draw.FieldColorSpec{}
	if metadataColorsFile != "" {
		if overrides, err = parseMetadataColorsYAML(metadataColorsFile); err != nil {
			return
		}
	}

	shapes = draw.ResolveFieldShapes(fields, overrides)
	values, err = draw.ResolveTipMetadata(fields, tipOrder, raw, overrides)
	return
}
