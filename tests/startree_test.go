package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree"
	"github.com/fredericlemoine/gotree/io/newick"
	"os"
	"strings"
	"testing"
)

func TestStarTree(t *testing.T) {
	ntips := 100
	t, err := tree.StarTree(ntips)
	if err != nil {
		t.Error(err)
	}
	if len(t.Tips()) != ntips {
		t.Error("The star tree should have 100 tips")
	}
	fmt.Fprintf(os.Stdout, "%s\n", t.Newick())
}
