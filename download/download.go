package download

import (
	"github.com/fredericlemoine/gotree/tree"
)

type ImageDownloader interface {
	Download(id string, format int) ([]byte, error) // Downdload a tree image from a server
}

type TreeDownloader interface {
	Download(id string) (*tree.Tree, error) // Download a tree from a server
}
