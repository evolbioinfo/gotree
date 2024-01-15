package cmd

import (
	"bufio"
	"fmt"
	goio "io"
	"os"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/goalign/io/fasta"
	"github.com/evolbioinfo/goalign/io/phylip"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/utils"
	"github.com/evolbioinfo/gotree/mutations"
	"github.com/evolbioinfo/gotree/tree"
	"github.com/spf13/cobra"
)

var mutationsalign string
var mutationsphylip bool
var mutationsinputstrict bool
var mutationseems bool
var outfile string

// mutationsCmd represents the mutations command
var mutationsCmd = &cobra.Command{
	Use:   "mutations",
	Short: "Extract the list of mutations along the branches of the phylogeny.",
	Long: `Extract the list of mutations along the branches of the phylogeny, given 
	the full list of ancestral (and terminal) sequences.

	The input tree must have internal node names specified and must be rooted.
	The input alignment (fasta or phylip only) must specify one sequence per internal 
	node name and tip.

	The output consists of the list of mutations that appear along the branches of the 
	tree, tab separated text file:

	1. Tree index (useful if several trees in the input tree file)
	2. Alignment site index
	3. Branch index
	4. Child node name
	5. Parent character
	6. Child character
	7. Number of descendent tips
	8. Number of descendent tips that have the child character

	If --eems is specified, then it will compute the number of emergences, i.e. the number of occurence of
	each mutation that is still present to at least ont tip. The columns of the output file will then be :
	1. Tree index (useful if several trees in the input tree file)
	2. Alignment site index
	5. Parent character
	6. Child character
	7. Number of emergence
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var align align.Alignment
		var fi goio.Closer
		var r *bufio.Reader
		var treefile goio.Closer
		var treechan <-chan tree.Trees
		var f *os.File
		var muts *mutations.MutationList

		// Reading the alignment
		if fi, r, err = utils.GetReader(mutationsalign); err != nil {
			io.LogError(err)
			return
		}
		if mutationsphylip {
			if align, err = phylip.NewParser(r, mutationsinputstrict).Parse(); err != nil {
				io.LogError(err)
				return
			}
		} else {
			if align, err = fasta.NewParser(r).Parse(); err != nil {
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

		if f, err = openWriteFile(outfile); err != nil {
			io.LogError(err)
			return
		}
		defer closeWriteFile(f, outfile)

		if mutationseems {
			fmt.Fprintf(f, "Tree ID\tSite\tParent Character\tChild Character\tNb EEMs\n")
		} else {
			fmt.Fprintf(f, "Tree ID\tSite\tBranch ID\tNode Name\tParent Character\tChild Character\tTotal tips\tSame Character Tips\n")
		}
		for t := range treechan {
			if mutationseems {
				if muts, err = mutations.CountEEMs(t.Tree, align); err != nil {
					io.LogError(err)
					return
				}
				for _, m := range muts.Mutations {
					fmt.Fprintf(f, "%d\t%d\t%c\t%c\t%d\n", t.Id, m.AlignmentSite, m.ParentCharacter, m.ChildCharacter, m.NumEEM)
				}
			} else {
				if muts, err = mutations.CountMutations(t.Tree, align); err != nil {
					io.LogError(err)
					return
				}
				for _, m := range muts.Mutations {
					fmt.Fprintf(f, "%d\t%d\t%d\t%s\t%c\t%c\t%d\t%d\n", t.Id, m.AlignmentSite, m.BranchIndex, m.ChildNodeName, m.ParentCharacter, m.ChildCharacter, m.NumTips, m.NumTipsWithChildCharacter)
				}
			}
		}
		return
	},
}

func init() {
	computeCmd.AddCommand(mutationsCmd)
	mutationsCmd.PersistentFlags().StringVarP(&mutationsalign, "align", "a", "stdin", "Alignment input file")
	mutationsCmd.PersistentFlags().BoolVarP(&mutationsphylip, "phylip", "p", false, "Alignment is in phylip? default : false (Fasta)")
	mutationsCmd.PersistentFlags().BoolVar(&mutationsinputstrict, "input-strict", false, "Strict phylip input format (only used with -p)")
	mutationsCmd.PersistentFlags().BoolVar(&mutationseems, "eems", false, "If true, extracts mutations that goes to tips, with their number of emergence (see https://doi.org/10.1101/2021.06.30.450558)")
	mutationsCmd.PersistentFlags().StringVarP(&intreefile, "input", "i", "stdin", "Input tree")
	mutationsCmd.PersistentFlags().StringVarP(&outfile, "output", "o", "stdout", "Output file")
}
