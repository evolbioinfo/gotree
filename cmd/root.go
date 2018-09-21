package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	goio "io"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fredericlemoine/cobrashell"
	"github.com/fredericlemoine/gotree/io/fileutils"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// Variables used in lots of commands
var inalignfile string
var intreefile, intree2file, outtreefile string
var outresfile string
var seed int64 = -1
var inputname string
var mapfile string
var revert bool
var transferdist bool
var deepestedge bool
var edgeformattext bool
var parsimonyAlgo string
var compareTips bool
var tipfile string
var cutoff float64
var replace bool
var treeformat = utils.FORMAT_NEWICK

var cfgFile string
var rootCpus int
var rootInputFormat string
var removeoutgroup bool
var rerootstrict bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gotree",
	Short: "gotree: A set of tools to handle phylogenetic trees in go",
	Long: `gotree is a set of tools to handle phylogenetic trees in go.

Different usages are implemented: 
- Generating random trees
- Transforming trees (renaming tips, pruning/removing tips)
- Comparing trees (computing bootstrap supports, counting common edges)
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch rootInputFormat {
		case "newick":
			treeformat = utils.FORMAT_NEWICK
		case "nexus":
			treeformat = utils.FORMAT_NEXUS
		case "phyloxml":
			treeformat = utils.FORMAT_PHYLOXML
		default:
			treeformat = utils.FORMAT_NEWICK
		}
		if seed == -1 {
			seed = time.Now().UTC().UnixNano()
		}
		rand.Seed(seed)
	},

	Run: func(cmd *cobra.Command, args []string) {
		s := cobrashell.New()
		// display welcome info.
		s.Println(fmt.Sprintf("Welcome to Gotree Console %s", Version))
		s.Println("type \"help\" to get a list of available commands")
		cobrashell.AddCommands(s, cmd.Root(), nil, cmd.Root().Commands()...)
		// We open a gotree console to interactively execute commands
		s.Run()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	maxcpus := runtime.NumCPU()

	RootCmd.PersistentFlags().Int64Var(&seed, "seed", -1, "Random Seed: -1 = nano seconds since 1970/01/01 00:00:00")
	RootCmd.PersistentFlags().IntVarP(&rootCpus, "threads", "t", 1, "Number of threads (Max="+strconv.Itoa(maxcpus)+")")
	RootCmd.PersistentFlags().StringVar(&rootInputFormat, "format", "newick", "Input tree format (newick, nexus, or phyloxml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func openWriteFile(file string) (f *os.File, err error) {
	if file == "stdout" || file == "-" {
		f = os.Stdout
	} else {
		f, err = os.Create(file)
	}
	return
}

func closeWriteFile(f goio.Closer, filename string) {
	if filename != "-" && filename != "stdout" {
		f.Close()
	}
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

/*File in output must be closed by calling function */
func readTrees(infile string) (treefile goio.Closer, treeChannel <-chan tree.Trees, err error) {
	// Read Tree
	var treereader *bufio.Reader

	if treefile, treereader, err = utils.GetReader(infile); err == nil {
		treeChannel = utils.ReadMultiTrees(treereader, treeformat)
	}

	return
}

func readTree(infile string) (t *tree.Tree, err error) {
	if infile != "none" {
		// Read comp Tree : Only one tree in input
		t, err = utils.ReadTree(infile, treeformat)
	} else {
		err = errors.New("Cannot use \"none\" as input file")
	}
	return
}

func parseTipsFile(file string) (tips []string, err error) {

	var treereader *bufio.Reader
	var treefile goio.Closer
	var line string
	var err2 error

	tips = make([]string, 0, 100)

	if treefile, treereader, err = utils.GetReader(file); err == nil {
		line, err2 = Readln(treereader)
		for err2 == nil {
			for _, name := range strings.Split(line, ",") {
				tips = append(tips, name)
			}
			line, err2 = Readln(treereader)
		}
		treefile.Close()
	}
	return
}

func readMapFile(file string, revert bool) (map[string]string, error) {
	outmap := make(map[string]string, 0)
	var mapfile *os.File
	var err error
	var reader *bufio.Reader

	if mapfile, err = os.Open(file); err != nil {
		return outmap, err
	}

	if strings.HasSuffix(file, ".gz") {
		if gr, err2 := gzip.NewReader(mapfile); err2 != nil {
			return outmap, err2
		} else {
			reader = bufio.NewReader(gr)
		}
	} else {
		reader = bufio.NewReader(mapfile)
	}
	line, e := fileutils.Readln(reader)
	nl := 1
	for e == nil {
		cols := strings.Split(line, "\t")
		if len(cols) != 2 {
			return outmap, errors.New("Map file does not have 2 fields at line: " + fmt.Sprintf("%d", nl))
		}
		if revert {
			outmap[cols[1]] = cols[0]
		} else {
			outmap[cols[0]] = cols[1]
		}
		line, e = fileutils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil
}
