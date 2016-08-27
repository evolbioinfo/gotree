package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gotree",
	Short: "gotree: A set of tools to handle phylogenetic trees in go",
	Long: `gotree is a set of tools to handle phylogenetic trees in go.

Different usages are implemented: 
- Generating random trees
- Transforming trees (renaming tips, pruning/removing tips)
- Comparing trees (computing bootstrap supports, counting common edges)
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
