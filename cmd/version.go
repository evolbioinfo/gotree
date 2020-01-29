package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines the version of gotree
// It is initialized during compilation
// with -ldflags "-X github.com/evolbioinfo/gotree/cmd.Version=Major.Minor..."
var Version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version of gotree",
	Long:  `Displays version of gotree.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
