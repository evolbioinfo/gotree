package tests

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/fredericlemoine/gotree/hashmap"
	"github.com/fredericlemoine/gotree/io/newick"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
)

func TestQuartets(t *testing.T) {
	var treefile *os.File
	var treereader *bufio.Reader
	var err error
	var quartet *tree.Tree

	// Parsing single tree newick file
	if treefile, treereader, err = utils.GetReader("data/quartets.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()

	quartet, err = newick.NewParser(treereader).Parse()
	if err != nil {
		t.Error(err)
	}
	quartet.UpdateTipIndex()
	nbspec, nbtotal := 0, 0
	quartet.Quartets(true, func(q *tree.Quartet) {
		nbspec++
	})
	quartet.Quartets(false, func(q *tree.Quartet) {
		nbtotal++
	})

	if nbspec != 4 {
		t.Error(fmt.Sprintf("There should be 4 specific quartets"))
	}
	if nbtotal != 6 {
		t.Error(fmt.Sprintf("There should be 6 total quartets"))
	}
}

func initQuartetData() (equals, conflict, different []*tree.Quartet) {
	// Equals
	equals = []*tree.Quartet{
		&tree.Quartet{1, 2, 3, 4},
		&tree.Quartet{2, 1, 3, 4},
		&tree.Quartet{1, 2, 4, 3},
		&tree.Quartet{2, 1, 4, 3},
		&tree.Quartet{3, 4, 1, 2},
		&tree.Quartet{3, 4, 2, 1},
		&tree.Quartet{4, 3, 1, 2},
		&tree.Quartet{4, 3, 2, 1},
	}

	// Conflict
	conflict = []*tree.Quartet{
		&tree.Quartet{1, 3, 2, 4},
		&tree.Quartet{1, 3, 4, 2},
		&tree.Quartet{3, 1, 2, 4},
		&tree.Quartet{3, 1, 4, 2},
		&tree.Quartet{2, 4, 1, 3},
		&tree.Quartet{2, 4, 3, 1},
		&tree.Quartet{4, 2, 1, 3},
		&tree.Quartet{4, 2, 3, 1},
		&tree.Quartet{1, 3, 2, 4},
		&tree.Quartet{1, 3, 4, 2},
		&tree.Quartet{3, 1, 2, 4},
		&tree.Quartet{3, 1, 4, 2},
		&tree.Quartet{2, 4, 1, 3},
		&tree.Quartet{2, 4, 3, 1},
		&tree.Quartet{4, 2, 1, 3},
		&tree.Quartet{4, 2, 3, 1},
	}

	// Different
	different = []*tree.Quartet{
		&tree.Quartet{1, 3, 5, 4},
		&tree.Quartet{1, 8, 4, 2},
		&tree.Quartet{3, 1, 2, 10},
		&tree.Quartet{10, 1, 4, 2},
		&tree.Quartet{7, 8, 1, 3},
		&tree.Quartet{2, 4, 7, 8},
		&tree.Quartet{7, 8, 9, 3},
		&tree.Quartet{4, 7, 8, 9},
		&tree.Quartet{7, 8, 9, 10},
	}
	return equals, conflict, different
}

func TestCompareQuartets(t *testing.T) {
	// reference quartet
	ref := &tree.Quartet{1, 2, 3, 4}
	equals, conflict, different := initQuartetData()

	for i, comp := range equals {
		if ref.Compare(comp) != tree.QUARTET_EQUALS {
			t.Error(fmt.Sprintf("Quartets ref and %d should be equal", i))
		}
	}

	for i, comp := range conflict {
		if ref.Compare(comp) != tree.QUARTET_CONFLICT {
			t.Error(fmt.Sprintf("Quartets ref and %d should be in conflict", i))
		}
	}

	for i, comp := range different {
		if ref.Compare(comp) != tree.QUARTET_DIFF {
			t.Error(fmt.Sprintf("Quartets ref and %d should be different", i))
		}
	}
}
func TestIndexQuartet(t *testing.T) {
	// Quartet index
	index := hashmap.NewHashMap(128, .75)

	// reference quartet
	ref := &tree.Quartet{1, 2, 3, 4}

	equals, conflict, different := initQuartetData()

	// We put the reference quartet
	index.PutValue(ref, ref)

	for i, q := range equals {
		_, ok := index.Value(q)
		if !ok {
			t.Error(fmt.Sprintf("Equal Quartet %d should be in the map", i))
		}
	}

	for i, q := range conflict {
		_, ok := index.Value(q)
		if !ok {
			t.Error(fmt.Sprintf("Conflicting Quartet %d should be in the map", i))
		}
	}

	for i, q := range different {
		_, ok := index.Value(q)
		if ok {
			t.Error(fmt.Sprintf("Different Quartet %d should be in the map", i))
		} else {
			index.PutValue(q, q)
			_, ok := index.Value(q)
			if !ok {
				t.Error(fmt.Sprintf("Newly inserted Quartet %d should be in the map", i))
			}
		}
	}

}

func TestIndexQuartets2(t *testing.T) {
	// Quartet index
	index := hashmap.NewHashMap(128, .75)

	// reference quartet
	ref := &tree.Quartet{1, 2, 3, 4}

	equals, conflict, different := initQuartetData()

	// We put the reference quartet
	index.PutValue(ref, ref)

	for _, q := range equals {
		index.PutValue(q, q)
	}

	for _, q := range conflict {
		index.PutValue(q, q)
	}

	for _, q := range different {
		index.PutValue(q, q)
	}

	l := len(index.Keys())
	if l != 10 {
		t.Error(fmt.Sprintf("There should be 10 elements in the quartet HashMap, but: %d", l))
	}
}

/**
Test with real data
*/
func TestIndexQuartets3(t *testing.T) {
	var treefile *os.File
	var treereader *bufio.Reader
	var err error
	var quartet *tree.Tree

	// Parsing single tree newick file
	if treefile, treereader, err = utils.GetReader("data/quartets.nw.gz"); err != nil {
		t.Error(err)
	}
	defer treefile.Close()

	quartet, err = newick.NewParser(treereader).Parse()
	if err != nil {
		t.Error(err)
	}
	quartet.UpdateTipIndex()

	index := quartet.IndexQuartets(false)

	l := len(index.Keys())
	if l != 5 {
		for _, h := range index.Keys() {
			q := h.(*tree.Quartet)
			fmt.Printf("%d - %d | %d - %d\n", q.T1, q.T2, q.T3, q.T4)
		}
		t.Error(fmt.Sprintf("There should be 5 quartets in the index, but: %d", l))
	}
}
