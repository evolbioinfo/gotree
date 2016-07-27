package main

import (
	"flag"
	"fmt"
	gotree "github.com/fredericlemoine/gotree/lib"
	"math/rand"
	"os"
	"time"
)

func main() {
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "Seed (Optional)")
	nbtips := flag.Int("nbtips", 10, "Number of tips (default 10)")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Seed   : %d\n", *seed)
	fmt.Fprintf(os.Stderr, "Nb Tips: %d\n", *nbtips)

	rand.Seed(*seed)
	t, err := gotree.RandomBinaryTree(*nbtips)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
	} else {
		fmt.Println(t.Newick())
		// fmt.Println(t.Root())
	}
}
