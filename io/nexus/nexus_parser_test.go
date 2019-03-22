package nexus_test

import (
	"strings"
	"testing"

	"github.com/evolbioinfo/gotree/io/nexus"
	"github.com/evolbioinfo/gotree/tree"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseTree(t *testing.T) {
	goodtrees := [...]string{
		`#NEXUS
BEGIN TAXA;
      TaxLabels fish frog snake mouse;
END;

BEGIN CHARACTERS;
      Dimensions NChar=40;
      Format DataType=DNA;
      Matrix
        fish   ACATA GAGGG TACCT CTAAA
        fish   ACATA GAGGG TACCT CTAAG

        frog   ACATA GAGGG TACCT CTAAC
        frog   CCATA GAGGG TACCT CTAAG

        snake  ACATA GAGGG TACCT CTAAG
        snake  GCATA GAGGG TACCT CTAAG

        mouse  ACATA GAGGG TACCT CTAAT
        mouse  TCATA GAGGG TACCT CTAAG
;
END;

BEGIN TREES;
      Tree best1=(fish, (frog, (snake, mouse)));
      Tree best2=(fish, (frog, (snake, mouse)));
      Tree best3=(fish, (frog, (snake, mouse)));
END;
EOF
`,
	}

	for i, intree := range goodtrees {
		nex, err := nexus.NewParser(strings.NewReader(intree)).Parse()
		if err != nil {
			t.Errorf("Tree %d ERROR: %s\n", i, err.Error())
		} else {
			if nex.NTrees() != 3 {
				t.Errorf("There should be 3 trees in the nexus file, and there are %q\n", nex.NTrees())
			} else {
				nex.IterateTrees(func(name string, tr *tree.Tree) {
					if tr.Newick() != "(fish,(frog,(snake,mouse)));" {
						t.Errorf("Tree should be: \"(fish,(frog,(snake,mouse)));\" and is %q\n", tr.Newick())
					}
				})
			}
			if !nex.HasAlignment {
				t.Errorf("Tree should be an alignment in the nexus file\n")
			} else {
				if nex.Alignment().NbSequences() != 4 {
					t.Errorf("Alignment should have 4 sequences but has %d\n", nex.Alignment().NbSequences())
				}
				if nex.Alignment().Length() != 40 {
					t.Errorf("Alignment should be 40 nt long, but is %d\n", nex.Alignment().Length())
				}
			}
		}
	}
}
