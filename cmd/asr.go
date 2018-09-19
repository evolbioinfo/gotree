package cmd

import (
	"bufio"
	"fmt"
	goio "io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fredericlemoine/goalign/align"
	"github.com/fredericlemoine/goalign/io/fasta"
	"github.com/fredericlemoine/goalign/io/phylip"
	"github.com/fredericlemoine/gotree/asr"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

var asralign string
var asrphylip bool
var asrinputstrict bool
var asrrandomresolve bool // Resolve ambiguities randomly in the downpass/deltran/acctran algo

// asrCmd represents the asr command
var asrCmd = &cobra.Command{
	Use:   "asr",
	Short: "Reconstructs most parsimonious ancestral sequences",
	Long: `Reconstructs most parsimonious ancestral sequences.

Depending on the chosen algorithm, it will run:
1) UP-PASS and
2) Either
   a) DOWN-PASS or
   b) DOWN-PASS+DELTRAN or
   c) ACCTRAN
   d) NONE

Should work on multifurcated trees

If --random-resolve is given then, during the last pass, each time 
a node with several possible states still exists, one state is chosen 
randomly before going deeper in the tree.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var align align.Alignment
		var fi goio.Closer
		var r *bufio.Reader
		var algo int
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var f *os.File

		rand.Seed(seed)

		switch strings.ToLower(parsimonyAlgo) {
		case "acctran":
			algo = asr.ALGO_ACCTRAN
		case "deltran":
			algo = asr.ALGO_DELTRAN
		case "downpass":
			algo = asr.ALGO_DOWNPASS
		case "none":
			algo = asr.ALGO_NONE
		default:
			io.ExitWithMessage(fmt.Errorf("Unkown parsimony algorithm: %s", parsimonyAlgo))
		}

		// Reading the alignment
		fi, r, err = utils.GetReader(asralign)
		if err != nil {
			io.ExitWithMessage(err)
		}
		if asrphylip {
			align, err = phylip.NewParser(r, asrinputstrict).Parse()
			if err != nil {
				io.ExitWithMessage(err)
			}
		} else {
			align, err = fasta.NewParser(r).Parse()
			if err != nil {
				io.ExitWithMessage(err)
			}
		}
		fi.Close()

		// Reading the trees
		if treefile, treechan, err = readTrees(intreefile); err != nil {
			io.LogError(err)
			return
		}
		defer treefile.Close()

		// Computing parsimony ASR and writing each trees
		if f, err = openWriteFile(outtreefile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outtreefile)

		for t := range treechan {
			err = asr.ParsimonyAsr(t.Tree, align, algo, asrrandomresolve)
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		return
	},
}

func init() {
	RootCmd.AddCommand(asrCmd)
	asrCmd.PersistentFlags().StringVarP(&asralign, "align", "a", "stdin", "Alignment input file")
	asrCmd.PersistentFlags().BoolVarP(&asrphylip, "phylip", "p", false, "Alignment is in phylip? default : false (Fasta)")
	asrCmd.PersistentFlags().BoolVar(&asrinputstrict, "input-strict", false, "Strict phylip input format (only used with -p)")
	asrCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	asrCmd.PersistentFlags().StringVarP(&outtreefile, "output", "o", "stdout", "Output file")
	asrCmd.PersistentFlags().StringVar(&parsimonyAlgo, "algo", "acctran", "Parsimony algorithm for resolving ambiguities: acctran, deltran, or downpass")
	asrCmd.PersistentFlags().BoolVar(&asrrandomresolve, "random-resolve", false, "Random resolve states when several possibilities in: acctran, deltran, or downpass")
	asrCmd.Flags().Int64VarP(&seed, "seed", "s", time.Now().UTC().UnixNano(), "Initial Random Seed")
}
