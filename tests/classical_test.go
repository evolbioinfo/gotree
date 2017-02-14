package tests

import (
	"fmt"
	"github.com/fredericlemoine/gotree/support"
	"testing"
)

func TestClassicalSupport(t *testing.T) {
	tr, e := support.Classical("data/rand_tree.nw.gz", "data/rand_tree_boot.nw.gz", 1)
	if e != nil {
		t.Error(e)
	}
	for _, e := range tr.Edges() {
		if !e.Right().Tip() && e.Support() != 0 {
			t.Error(fmt.Sprintf("Non Tip support should be 0.00 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
	tr, e = support.Classical("data/rand_tree.nw.gz", "data/rand_tree_same.nw.gz", 1)
	if e != nil {
		t.Error(e)
	}
	for _, e := range tr.Edges() {
		if !e.Right().Tip() && e.Support() != 1.00 {
			t.Error(fmt.Sprintf("Non Tip support should be 1.00 and is %.2f", e.Support()))
		} else if e.Right().Tip() && e.Support() != -1 {
			t.Error(fmt.Sprintf("Tip support should be -1.00 and is %.2f", e.Support()))
		}
	}
}
