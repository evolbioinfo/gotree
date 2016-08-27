package cmd

import (
	"github.com/spf13/cobra"
	"runtime"
	"strconv"
)

var supportIntree string
var supportBoottrees string
var supportCpus int
var supportOutFile string

// supportCmd represents the support command
var supportCmd = &cobra.Command{
	Use:   "support",
	Short: "Computes different kind of branch supports",
	Long: `Computes different kind of branch supports.

The supports implemented are :
- mast-like support
- parsimony based support
- Classical Felsenstein support

`,
}

func init() {
	computeCmd.AddCommand(supportCmd)
	maxcpus := runtime.NumCPU()

	supportCmd.PersistentFlags().StringVarP(&supportIntree, "reftree", "i", "stdin", "Reference tree input file")
	supportCmd.PersistentFlags().StringVarP(&supportBoottrees, "bootstrap", "b", "none", "Bootstrap trees input file")
	supportCmd.PersistentFlags().StringVarP(&supportOutFile, "out", "o", "stdout", "Output tree file, with supports")
	supportCmd.PersistentFlags().IntVarP(&supportCpus, "threads", "t", 1, "Number of threads (Max "+strconv.Itoa(maxcpus)+")")
}
