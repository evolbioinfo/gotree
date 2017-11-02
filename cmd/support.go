package cmd

import (
	"os"

	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var supportIntree string
var supportBoottrees string
var supportOutFile string
var supportLogFile string
var movedtaxa bool
var taxperbranches bool     // If we should compute all avg tax transfers per branches
var hightaxperbranches bool // If we should compute all avg tax transfers per branches

// For booster computation : output tree with raw avg distances as supports
// in the form: branchid|avg_distance|depth
var rawSupportOutputFile string
var rawSupportOut *os.File
var supportOut *os.File
var supportLog *os.File
var supportSilent bool

// supportCmd represents the support command
var supportCmd = &cobra.Command{
	Use:   "support",
	Short: "Computes different kind of branch supports",
	Long: `Computes different kind of branch supports.

The supports implemented are :
- booster support
- Classical Felsenstein support

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if supportOutFile != "stdout" {
			supportOut, err = os.Create(supportOutFile)
		} else {
			supportOut = os.Stdout
		}
		if err != nil {
			io.ExitWithMessage(err)
		}
		if supportLogFile != "stderr" {
			supportLog, err = os.Create(supportLogFile)
		} else {
			supportLog = os.Stderr
		}
		if err != nil {
			io.ExitWithMessage(err)
		}
		if rawSupportOutputFile != "none" {
			if rawSupportOutputFile != "stdout" {
				rawSupportOut, err = os.Create(rawSupportOutputFile)
			} else {
				rawSupportOut = os.Stdout
			}
			if err != nil {
				io.ExitWithMessage(err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		supportOut.Close()
		supportLog.Close()
	},
}

func init() {
	computeCmd.AddCommand(supportCmd)

	supportCmd.PersistentFlags().StringVarP(&supportIntree, "reftree", "i", "stdin", "Reference tree input file")
	supportCmd.PersistentFlags().StringVarP(&supportBoottrees, "bootstrap", "b", "none", "Bootstrap trees input file")
	supportCmd.PersistentFlags().StringVarP(&supportOutFile, "out", "o", "stdout", "Output tree file, with supports")
	supportCmd.PersistentFlags().StringVarP(&supportLogFile, "log-file", "l", "stderr", "Output log file")
	supportCmd.PersistentFlags().BoolVar(&supportSilent, "silent", false, "If true, progress messages will not be printed to stderr")
}
