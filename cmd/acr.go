package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"

	"github.com/fredericlemoine/gotree/acr"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var acrstates string

// acrCmd represents the acr command
var acrCmd = &cobra.Command{
	Use:   "acr",
	Short: "Reconstructs most parsimonious ancestral characters",
	Long: `Reconstructs most parsimonious ancestral characters.

It does 2 tree straversal:
1) One postorder
2) One preorder

Works on multifurcated trees, by taking the most frequent state(s).

`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Reading tip state in an input file
		tipstates := parseTipStates(acrstates)
		// Reading the trees
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		// Computing parsimony ACR and writing each trees
		f := openWriteFile(outtreefile)
		for t := range treechan {
			err = acr.ParsimonyAcr(t.Tree, tipstates)
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(acrCmd)
	acrCmd.PersistentFlags().StringVarP(&acrstates, "states", "s", "stdin", "Tip state file (One line per tip, tab separated: tipname\\tstate)")
	acrCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	acrCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
}

func parseTipStates(file string) map[string]string {
	var f *os.File
	var r *bufio.Reader
	states := make(map[string]string)
	var err error

	if file == "stdin" || file == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(file)
		if err != nil {
			io.ExitWithMessage(err)
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err := gzip.NewReader(f); err != nil {
			io.ExitWithMessage(err)
		} else {
			r = bufio.NewReader(gr)
		}
	} else {
		r = bufio.NewReader(f)
	}
	l, e := Readln(r)
	for e == nil {
		cols := strings.Split(l, "\t")
		if cols == nil || len(cols) != 2 {
			io.ExitWithMessage(errors.New("Bad format for tip states: Wrong number of columns"))
		}
		states[cols[0]] = cols[1]
		l, e = Readln(r)
	}
	return states
}
