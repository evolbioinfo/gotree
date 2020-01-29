package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	goio "io"
	"os"
	"regexp"
	"strings"

	"github.com/evolbioinfo/gotree/acr"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
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
   d) NONE

Should work on multifurcated trees.

If --random-resolve is given then, during the last pass, each time 
a node with several possible states still exists, one state is chosen 
randomly before going deeper in the tree.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var algo int
		var statemap map[string]string
		var tipstates map[string]string
		var resfile *os.File
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var f *os.File

		switch strings.ToLower(parsimonyAlgo) {
		case "acctran":
			algo = acr.ALGO_ACCTRAN
		case "deltran":
			algo = acr.ALGO_DELTRAN
		case "downpass":
			algo = acr.ALGO_DOWNPASS
		case "none":
			algo = acr.ALGO_NONE
		default:
			io.LogError(fmt.Errorf("Unkown parsimony algorithm: %s", parsimonyAlgo))
			return
		}
		// Reading tip state in an input file
		if tipstates, err = parseTipStates(acrstates); err != nil {
			io.LogError(err)
			return
		}
		// Reading the trees

		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		// Computing parsimony ACR and writing each trees
		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		if outresfile != "none" {
			if resfile, err = openWriteFile(outresfile); err != nil {
				io.LogError(err)
				return
			}
			defer closeWriteFile(resfile, outresfile)
		}
		for t := range treechan {
			statemap, err = acr.ParsimonyAcr(t.Tree, tipstates, algo, acrrandomresolve)
			if err != nil {
				io.LogError(err)
				return
			}
			f.WriteString(t.Tree.Newick() + "\n")
			if outresfile != "none" {
				for k, v := range statemap {
					resfile.WriteString(fmt.Sprintf("%s,%s\n", k, v))
				}
			}
		}
		return
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
}

func parseTipStates(file string) (states map[string]string, err error) {
	var f *os.File
	var r *bufio.Reader
	var gr *gzip.Reader

	states = make(map[string]string)

	if file == "stdin" || file == "-" {
		f = os.Stdin
	} else {
		if f, err = os.Open(file); err != nil {
			return
		}
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err = gzip.NewReader(f); err != nil {
			return
		}
		r = bufio.NewReader(gr)
	} else {
		r = bufio.NewReader(f)
	}
	l, e := Readln(r)
	// Split using either '\t' or ','
	re := regexp.MustCompile("\t|,")
	for e == nil {
		cols := re.Split(l, -1)
		if cols == nil || len(cols) != 2 {
			err = errors.New("Bad format for tip states: Wrong number of columns")
			return
		}
		states[cols[0]] = cols[1]
		l, e = Readln(r)
	}
	return
}
