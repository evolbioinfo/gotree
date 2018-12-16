package download

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ItolImageDownloader struct {
	config map[string]string
}

func NewItolImageDownloader(config map[string]string) *ItolImageDownloader {
	return &ItolImageDownloader{config}
}

// Down a tree image from ITOL
func (d *ItolImageDownloader) Download(id string, format int) ([]byte, error) {
	posturl := "https://itol.embl.de/batch_downloader.cgi"
	var err error
	var postresponse *http.Response
	var responsebody []byte

	form := url.Values{}
	form.Add("tree", id)

	strformat := StrFormat(format)
	if strformat == "unknown" {
		return nil, errors.New("Output image format unknown")
	}
	form.Add("format", strformat)

	for k, v := range d.config {
		if k != "" && v != "" {
			form.Add(k, v)
		}
	}
	postresponse, err = http.PostForm(posturl, form)

	defer postresponse.Body.Close()
	if responsebody, err = ioutil.ReadAll(postresponse.Body); err != nil {
		return nil, err
	}

	if postresponse.Header.Get("Content-Type") == "text/html" {
		return nil, errors.New(string(responsebody))
	}
	return responsebody, nil
}
