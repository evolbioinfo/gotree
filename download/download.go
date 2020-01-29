package download

import (
	"github.com/evolbioinfo/gotree/tree"
)

// ImageDownloader defines function to download an image
type ImageDownloader interface {
	Download(id string, format int) ([]byte, error) // Downdload a tree image from a server
}

// TreeDownloader  defines function to download an Tree
type TreeDownloader interface {
	Download(id string) (*tree.Tree, error) // Download a tree from a server
}
