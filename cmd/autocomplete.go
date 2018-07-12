package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// itolCmd represents the itol command
var autocompleteCmd = &cobra.Command{
	Use:   "completion SHELL",
	Args:  cobra.ExactArgs(1),
	Short: "Generates auto-completion commands for bash or zsh",
	Long: `Generates auto-completion commands for bash or zsh. 

Examples (Largely inspired from kubectl command):
  # bash completion on macOS using homebrew
  ## If running Bash 3.2 included with macOS
  brew install bash-completion
  ## or, if running Bash 4.1+
  brew install bash-completion@2
  # Then add auto completion commands
  gotree completion bash > $(brew --prefix)/etc/bash_completion.d/gotree

  
  # Installing bash completion on Linux
  ## Load the gotree completion code for bash into the current shell
  source <(gotree completion bash)
  ## Write bash completion code to a file and source if from .bash_profile
  mkdir ~/.gotree
  gotree completion bash > ~/.gotree/completion.bash.inc
  printf "
  # gotree shell completion
  source '$HOME/.gotree/completion.bash.inc'
  " >> $HOME/.bashrc
  source $HOME/.bashrc

  # Load the gotree completion code for zsh[1] into the current shell
  source <(gotree completion zsh)
  # Set the gotree completion code for zsh[1] to autoload on startup
  gotree completion zsh > "${fpath[1]}/_gotree"

`,
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		if shell == "bash" {
			RootCmd.GenBashCompletion(os.Stdout)
		} else if shell == "zsh" {
			RootCmd.GenZshCompletion(os.Stdout)
		}
	},
}

func init() {
	RootCmd.AddCommand(autocompleteCmd)
}
