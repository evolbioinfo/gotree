package cmd

import (
	"bufio"
	"fmt"
	goio "io"
	"strings"

	"github.com/fredericlemoine/goalign/align"
	"github.com/fredericlemoine/goalign/io/fasta"
	"github.com/fredericlemoine/goalign/io/phylip"
	"github.com/fredericlemoine/gotree/asr"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/spf13/cobra"
)

var asralign string
var asrphylip bool
var asrinputstrict bool

// asrCmd represents the asr command
var asrCmd = &cobra.Command{
	Use:   "asr",
	Short: "Reconstructs most parsimonious ancestral sequences",
	Long: `Reconstructs most parsimonious ancestral sequences.

It does 2 tree straversal:
1) One postorder
2) One preorder

Works on multifurcated trees, by taking the most frequent state(s).

`,
	Run: func(cmd *cobra.Command, args []string) {
		var align align.Alignment
		var fi goio.Closer
		var r *bufio.Reader
		var err error
		var algo int

		switch strings.ToLower(parsimonyAlgo) {
		case "acctran":
			algo = asr.ALGO_ACCTRAN
		case "deltran":
			algo = asr.ALGO_DELTRAN
		case "downpass":
			algo = asr.ALGO_DOWNPASS
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
		treefile, treechan := readTrees(intreefile)
		defer treefile.Close()

		// Computing parsimony ASR and writing each trees
		f := openWriteFile(outtreefile)
		for t := range treechan {
			err = asr.ParsimonyAsr(t.Tree, align, algo)
			if err != nil {
				io.ExitWithMessage(err)
			}
			f.WriteString(t.Tree.Newick() + "\n")
		}
		f.Close()
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
}
