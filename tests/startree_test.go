package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/tree"
	"os"
	"testing"
)

func TestStarTree(t *testing.T) {
	ntips := 100
	tr, err := tree.StarTree(ntips)
	if err != nil {
		t.Error(err)
	}
	if len(tr.Tips()) != ntips {
		t.Error("The star tree should have 100 tips")
	}
	fmt.Fprintf(os.Stdout, "%s\n", tr.Newick())

}
