package cmd

import (
	"bufio"
	"fmt"
	goio "io"
	"os"
	"strings"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/goalign/io/fasta"
	"github.com/evolbioinfo/goalign/io/phylip"
	"github.com/evolbioinfo/gotree/asr"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var asralign string
var asrphylip bool
var asrinputstrict bool
var asrrandomresolve bool // Resolve ambiguities randomly in the downpass/deltran/acctran algo
var outlogfile string

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
		var logf *os.File
		var nsteps []int

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
			err = fmt.Errorf("Unkown parsimony algorithm: %s", parsimonyAlgo)
			io.LogError(err)
			return
		}

		// Reading the alignment
		fi, r, err = utils.GetReader(asralign)
		if err != nil {
			io.LogError(err)
			return
		}
		if asrphylip {
			align, err = phylip.NewParser(r, asrinputstrict).Parse()
			if err != nil {
				io.LogError(err)
				return
			}
		} else {
			align, err = fasta.NewParser(r).Parse()
			if err != nil {
				io.LogError(err)
				return
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
		if logf, err = openWriteFile(outlogfile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(logf, outlogfile)

		for t := range treechan {
			nsteps, err = asr.ParsimonyAsr(t.Tree, align, algo, asrrandomresolve)
			if err != nil {
				io.LogError(err)
				return
			}
			fmt.Fprintf(logf, "steps")
			for _, s := range nsteps {
				fmt.Fprintf(logf, " %d", s)
			}
			fmt.Fprintf(logf, "\n")
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
	asrCmd.PersistentFlags().StringVar(&outlogfile, "log", "stdout", "Output log file")
	asrCmd.PersistentFlags().StringVar(&parsimonyAlgo, "algo", "acctran", "Parsimony algorithm for resolving ambiguities: acctran, deltran, or downpass")
	asrCmd.PersistentFlags().BoolVar(&asrrandomresolve, "random-resolve", false, "Random resolve states when several possibilities in: acctran, deltran, or downpass")
}
