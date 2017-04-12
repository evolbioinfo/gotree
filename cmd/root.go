package cmd

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
)

// Variables used in lots of commands
var intreefile, intree2file, outtreefile string
var seed int64
var inputname string
var mapfile string
var revert bool
var transferdist bool
var compareTips bool
var tipfile string
var cutoff float64
var movedtaxa bool

var cfgFile string
var rootCpus int

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
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	maxcpus := runtime.NumCPU()
	RootCmd.PersistentFlags().IntVarP(&rootCpus, "threads", "t", 1, "Number of threads (Max="+strconv.Itoa(maxcpus)+")")

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func openWriteFile(file string) *os.File {
	if file == "stdout" || file == "-" {
		return os.Stdout
	} else {
		f, err := os.Create(file)
		if err != nil {
			io.ExitWithMessage(err)
		}
		return f
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

func readTrees(infile string) <-chan tree.Trees {
	// Read Tree
	intreesChan := make(chan tree.Trees, 15)
	/* Read ref tree(s) */
	go func() {
		if _, err := utils.ReadMultiTreeFile(infile, intreesChan); err != nil {
			io.ExitWithMessage(err)
		}
		close(intreesChan)
	}()
	return intreesChan
}

func readTree(infile string) *tree.Tree {
	var tree *tree.Tree
	var err error
	if infile != "none" {
		// Read comp Tree : Only one tree in input
		tree, err = utils.ReadRefTree(infile)
		if err != nil {
			io.ExitWithMessage(err)
		}
	}
	return tree
}

func parseTipsFile(file string) []string {
	tips := make([]string, 0, 100)
	f, r, err := utils.GetReader(file)
	if err != nil {
		io.ExitWithMessage(err)
	}
	l, e := Readln(r)
	for e == nil {
		for _, name := range strings.Split(l, ",") {
			tips = append(tips, name)
		}
		l, e = Readln(r)
	}
	f.Close()
	return tips
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
	line, e := utils.Readln(reader)
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
		line, e = utils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil
}
