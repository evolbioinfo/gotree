package support_test

import (
	"strings"
	"testing"

	"github.com/evolbioinfo/gotree/io/newick"
	"github.com/evolbioinfo/gotree/support"
	"github.com/evolbioinfo/gotree/tree"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestMinDist(t *testing.T) {
	reftrees := [...]string{
		// Here branch length encodes the expected transfer distance
		// to the given compared tree
		"(t1:0,t2:0,(t3:0,(t4:0,t5:0):0):1);",
		"(t1:0,t2:0,(t3:0,(t4:0,t5:0):0):1);",
		"(t1:0,t2:0,(t3:0,(t4:0,t5:0):0):1);",
		"((1:0,2:0)1:0,(3:0,4:0)2:0,((7:0,8:0)3:0,(5:0,6:0)4:0)5:2);",
		"((1:0,2:0,3:0,4:0,5:0,6:0):0,7:0,(8:0,9:0,10:0,11:0,12:0,13:0,14:0):1);",
		"(((a:0,b:0):0,(c:0,d:0):1):0, (e:0,f:0):1,(g:0,h:0):0);",
		"(((a:0,b:0):0,(c:0,d:0):1):1, (e:0,f:0):1,(g:0,h:0):1);",
	}
	comptrees := [...]string{
		"(t1,t3,(t2,(t4,t5)));",
		"(t4,t5,(t2,(t1,t3)));",
		"((t1,t3),t2,(t4,t5));",
		"((7,8),(3,4),((1,2),(5,6)));",
		"((1,2,3,4,5,6),(7,8,9),(10,11,12,13,14));",
		"(a, b, (c, d, (e, f, (g, h))));",
		"((a,b), d, ((f,c),((g,e),h)));",
	}

	var reftree, comptree *tree.Tree
	var err error
	for i, intree := range reftrees {
		// parse reftree
		if reftree, err = newick.NewParser(strings.NewReader(intree)).Parse(); err != nil {
			t.Errorf("Tree %d ERROR: %s\n", i, err.Error())
			return
		}
		if err = reftree.ReinitIndexes(); err != nil {
			t.Error(err)
			return
		}
		// parse comptree
		if comptree, err = newick.NewParser(strings.NewReader(comptrees[i])).Parse(); err != nil {
			t.Errorf("Tree %d ERROR: %s\n", i, err.Error())
			return
		}
		if err = comptree.ReinitIndexes(); err != nil {
			t.Error(err)
			return
		}

		ntips := len(reftree.Tips())
		compedges := comptree.Edges()
		for _, e1 := range reftree.Edges() {
			dist, minedge, _, _ := support.MinTransferDist(e1, reftree, comptree, ntips, compedges, false)
			p, _ := e1.TopoDepth()
			if dist != int(e1.Length()) {
				bp, _ := minedge.TopoDepth()
				t.Errorf("Tree %d: (p=%d, s=%f, bp=%d) Min dist ERROR: dist is %d, should be %d\n%s\n%s", i, p, e1.Support(), bp, dist, int(e1.Length()), reftree.Newick(), comptree.Newick())
				return
			}
		}
	}
}
