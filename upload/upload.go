package upload

import (
	"github.com/evolbioinfo/gotree/tree"
)

type Uploader interface {
	Upload(name string, t *tree.Tree) (string, string, error)        // Upload a tree with a given name to a server and returns the url to access the tree, and full server response
	UploadNewick(name string, newick string) (string, string, error) // Upload a newick tree with a given name to a server and returns the url to access the tree, and full server response
}
