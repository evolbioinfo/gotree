package draw

/* Package intended to draw phylogenetic trees on different devices :
- Terminal,
- Images (svg, png)
- ...
And with different drawing algorithms. So far, only ASCII form in terminal.
*/

import (
	"github.com/fredericlemoine/gotree/tree"
)

type TreeDrawer interface {
	DrawTree(t *tree.Tree) error
}
