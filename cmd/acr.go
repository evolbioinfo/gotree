package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fredericlemoine/gotree/acr"
	"github.com/fredericlemoine/gotree/io"
	"github.com/spf13/cobra"
)

var acrstates string
var acrrandomresolve bool // Resolve ambiguities randomly in the downpass/deltran/acctran algo

// acrCmd represents the acr command
var acrCmd = &cobra.Command{
	Use:   "acr",
	Short: "Reconstructs most parsimonious ancestral characters",
	Long: `Reconstructs most parsimonious ancestral characters.

Depending on the chosen algorithm, it will run:
1) UP-PASS and
2) Either
   a) DOWN-PASS or
   b) DOWN-PASS+DELTRAN or
   c) ACCTRAN

Should work on multifurcated trees.

If --random-resolve is given then, during the last pass, each time 
a node with several possible states still exists, one state is chosen 
randomly before going deeper in the tree.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var algo int
		var statemap map[string]string
		var resfile *os.File
		rand.Seed(seed)

		switch strings.ToLower(parsimonyAlgo) {
		case "acctran":
			algo = acr.ALGO_ACCTRAN
		case "deltran":
			algo = acr.ALGO_DELTRAN
		case "downpass":
			algo = acr.ALGO_DOWNPASS
		default:
			io.ExitWithMessage(fmt.Errorf("Unkown parsimony algorithm: %s", parsimonyAlgo))
		}
		// Reading tip state in an input file
		tipstates := parseTipStates(acrstates)
		// Reading the trees
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		// Computing parsimony ACR and writing each trees
		f := openWriteFile(outtreefile)
		defer f.Close()

		if outresfile != "none" {
			resfile = openWriteFile(outresfile)
			defer resfile.Close()
		}
		for t := range treechan {
			statemap, err = acr.ParsimonyAcr(t.Tree, tipstates, algo, acrrandomresolve)
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
			if outresfile != "none" {
				for k, v := range statemap {
					resfile.WriteString(fmt.Sprintf("%s,%s\n", k, v))
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(acrCmd)
	acrCmd.PersistentFlags().StringVar(&acrstates, "states", "stdin", "Tip state file (One line per tip, tab separated: tipname\\tstate)")
	acrCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	acrCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	acrCmd.PersistentFlags().StringVar(&outresfile, "out-states", "none", "Output mapping file between node names and states")
	acrCmd.PersistentFlags().StringVar(&parsimonyAlgo, "algo", "acctran", "Parsimony algorithm for resolving ambiguities: acctran, deltran, or downpass")
	acrCmd.PersistentFlags().BoolVar(&acrrandomresolve, "random-resolve", false, "Random resolve states when several possibilities in: acctran, deltran, or downpass")
	acrCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")

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
	// Split using either '\t' or ','
	re := regexp.MustCompile("\t|,")
	for e == nil {
		cols := re.Split(l, -1)
		if cols == nil || len(cols) != 2 {
			io.ExitWithMessage(errors.New("Bad format for tip states: Wrong number of columns"))
		}
		states[cols[0]] = cols[1]
		l, e = Readln(r)
	}
	return states
}
