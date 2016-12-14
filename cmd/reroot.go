package cmd

import (
	"github.com/spf13/cobra"
)

var reroottipfile string
var rerootinputfile string
var rerootoutputfile string

// rerootCmd represents the reroot command
var rerootCmd = &cobra.Command{
	Use:   "reroot",
	Short: "Reroot trees using an outgroup or at midpoint",
	Long: `Reroot trees using an outgroup or at midpoint.
`,
}

func init() {
	RootCmd.AddCommand(rerootCmd)
	rerootCmd.PersistentFlags().StringVarP(&rerootinputfile, "input", "i", "stdin", "Input Tree")
	rerootCmd.PersistentFlags().StringVarP(&rerootoutputfile, "output", "o", "stdout", "Rerooted output tree file")
}
