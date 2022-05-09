package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var drawNoTipLabels bool
var drawNoBranchLengths bool
var drawInternalNodeLabels bool
var drawSupport bool
var drawSupportCutoff float64
var drawInternalNodeSymbols bool
var drawNodeComment bool
var annotFile string

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
	drawCmd.PersistentFlags().StringVarP(&annotFile, "annotation-file", "f", "", "Annotation file to add colored circles to tip nodes (svg & png)\nTab separated, with <tip-name  Red  Green  Blue> or\n<tip-name hex-value> on each line")
}

// Parse tab separated value file to add colored nodes to specific tips
func parseAnnot(filepath string) (map[string][]uint8, error) {

	colors := make(map[string][]uint8)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Tab separated values file

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		colors[record[0]] = make([]uint8, 3)

		if len(record) == 4 { // R,G,B format
			for i, col := range record[1:] {
				comp, err := strconv.ParseUint(col, 10, 8)
				if err != nil {
					return nil, err
				}
				colors[record[0]][i] = uint8(comp)
			}
		} else if len(record) == 2 { // HEX format

			// adapted from stackoverflow.com/a/54200713
			switch len(record[1]) {
			case 7:
				_, err = fmt.Sscanf(record[1], "#%02x%02x%02x", &colors[record[0]][0], &colors[record[0]][1], &colors[record[0]][2])
			case 9: // Discard Alpha channel
				var ignore uint8
				_, err = fmt.Sscanf(record[1], "#%02x%02x%02x%02x", &colors[record[0]][0], &colors[record[0]][1], &colors[record[0]][2], &ignore)
			default:
				err = fmt.Errorf("invalid length (%v) for hex code, must be 7 or 9", len(record[1]))
			}

			if err != nil {
				return nil, err
			}

		} else { // Wrong format
			return nil, fmt.Errorf("annotation file is the wrong format. (Expecting 4 or 2 fields got %d)", len(record))
		}

	}

	return colors, nil
}
