package cmd

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
	"os"
)

// difftipsCmd represents the difftips command
var difftipsCmd = &cobra.Command{
	Use:   "difftips",
	Short: "Print diff between tip names of two trees",
	Long: `Print diff between tip names of two trees.

For example:
t1.nh : (t1,t2,(t3,t4));
t2.nh : (t10,t2,(t3,t4));

gotree difftips -i t1.nh -c t2.nh

should produce the following output:
< t1
> t10
= 3

`,
	Run: func(cmd *cobra.Command, args []string) {
		if intree2file == "none" {
			io.ExitWithMessage(errors.New("Compare tree file must be provided with -c"))
		}
		eq := 0

		refTree := readTree(intreefile)
		compTree := readTree(intree2file)

		for _, t := range refTree.Tips() {
			if ok, err3 := compTree.ExistsTip(t.Name()); err3 != nil {
				io.ExitWithMessage(err3)
			} else {
				if !ok {
					fmt.Fprintf(os.Stdout, "< %s\n", t.Name())
				} else {
					eq++
				}
			}
		}
		for _, t := range compTree.Tips() {
			if ok, err4 := refTree.ExistsTip(t.Name()); err4 != nil {
				io.ExitWithMessage(err4)
			} else {
				if !ok {
					fmt.Fprintf(os.Stdout, "> %s\n", t.Name())
				}
			}
		}
		fmt.Fprintf(os.Stdout, "= %d\n", eq)
	},
}

func init() {
	RootCmd.AddCommand(difftipsCmd)
	difftipsCmd.Flags().StringVarP(&intreefile, "reftree", "i", "stdin", "Reference tree input file")
	difftipsCmd.Flags().StringVarP(&intree2file, "compared", "c", "none", "Other tree file to compare with")
}
