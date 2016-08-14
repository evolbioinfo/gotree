package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	gotree "github.com/fredericlemoine/gotree/lib"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "Seed (Optional)")
	nbtips := flag.Int("nbtips", 10, "Number of tips (default 10)")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Seed   : %d\n", *seed)
	fmt.Fprintf(os.Stderr, "Nb Tips: %d\n", *nbtips)

	rand.Seed(*seed)

	intree := "(t1:1,t2:2,(t3:3,(t4:4,t5:5)0.8:6)0.9:7);"
	tree, err2 := newick.NewParser(strings.NewReader(intree)).Parse()
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(tree.Newick())
	if err3 := tree.RemoveTips("t2", "t3"); err3 != nil {
		panic(err3)
	}

	fmt.Println(tree.Newick())

	fmt.Println("Generating Tree")
	t, err := gotree.RandomBinaryTree(*nbtips)
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
				panic(err)
			}
			defer func() {
				if err := fi.Close(); err != nil {
					panic(err)
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
				_, common, _, err := t.CommonEdges(tree, false)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Common edges: " + strconv.Itoa(common))
				}
				_, _, common, err = tree.CommonEdges(tree, false)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Common edges: " + strconv.Itoa(common))
				}
			}
		}
	}
}
