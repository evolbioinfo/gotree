package cmd

import (
	"bufio"
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var gotreeSeed int64
var gotreeNbTips int

// gotreeCmd represents the gotree command
var gotreeCmd = &cobra.Command{
	Use:   "gotree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		seed := gotreeSeed
		nbtips := gotreeNbTips

		fmt.Fprintf(os.Stderr, "Seed   : %d\n", seed)
		fmt.Fprintf(os.Stderr, "Nb Tips: %d\n", nbtips)

		rand.Seed(seed)

		intree := "(t1:1,t2:2,(t3:3,(t4:4,t5:5)0.8:6)0.9:7);"
		t, err2 := newick.NewParser(strings.NewReader(intree)).Parse()
		if err2 != nil {
			io.ExitWithMessage(err2)
		}
		fmt.Println(t.Newick())
		if err3 := t.RemoveTips("t2", "t3"); err3 != nil {
			io.ExitWithMessage(err3)
		}

		fmt.Println(t.Newick())

		fmt.Println("Generating Tree")
		t, err := tree.RandomBinaryTree(nbtips)
		fmt.Println("Done")

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		} else {
			f, err := os.Create("/tmp/tree")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Writing to /tmp/tree")
				f.WriteString(t.Newick() + "\n")
				f.Close()
				fmt.Println("Done")
				// fmt.Println(t.Newick())
				fmt.Println("Reading from /tmp/tree")
				//gotree.FromNewickFile("/tmp/tree")
				fmt.Println("Done")
				fi, err := os.Open("t2.nh")
				if err != nil {
					io.ExitWithMessage(err)
				}
				defer func() {
					if err := fi.Close(); err != nil {
						io.ExitWithMessage(err)
					}
				}()
				r := bufio.NewReader(fi)
				tree, err := newick.NewParser(r).Parse()

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(tree.Newick())
					for _, e := range tree.Edges() {
						fmt.Println(e.DumpBitSet())
					}
					names := tree.AllTipNames()

					for _, name := range names {
						i, err := tree.TipIndex(name)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Printf("name: %s | index: %d\n", name, i)
						}
					}
					sort.Strings(names)
					for _, name := range names {
						i, err := tree.TipIndex(name)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(name + " [" + strconv.Itoa(int(i)) + "]")
						}
					}
					_, common, err := t.CommonEdges(tree, false)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Common edges: " + strconv.Itoa(common))
					}
					_, common, err = tree.CommonEdges(tree, false)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Common edges: " + strconv.Itoa(common))
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(gotreeCmd)
	gotreeCmd.Flags().Int64VarP(&gotreeSeed, "seed", "s", time.Now().UTC().UnixNano(), "Seed (Optional)")
	gotreeCmd.Flags().IntVarP(&gotreeNbTips, "nbtips", "t", 10, "Number of tips")
}
