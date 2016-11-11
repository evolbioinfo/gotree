package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/utils"
	"os"
	"testing"
)

/*
 Function to test consensus tree generation
 It compares majority and strict consensus to a given already computed
 consensus (from phylip consense)
*/
func TestQuartets(t *testing.T) {
	fmt.Fprintf(os.Stderr, "Started Quartets\n")
	quartet, _ := utils.ReadRefTree("data/quartets.nw.gz")
	quartet.Quartets(func(tb1, tb2, tb3, tb4 uint) {
		fmt.Fprintf(os.Stderr, "(%d,%d)(%d,%d)", tb1, tb2, tb3, tb4)
	})
	fmt.Fprintf(os.Stderr, "End Quartets\n")
}
