package download

type ImageDownloader interface {
	Download(id string) ([]byte, error) // Down a tree image from a server
}
