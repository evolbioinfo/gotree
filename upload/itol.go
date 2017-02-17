package upload

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/fredericlemoine/gotree/tree"
)

type ItolUploader struct {
	uploadid        string
	projectname     string
	annotationfiles []string
}

// Initialize a new itoluploader
// if uploadid=="", then tree will be public and deleted
// after 30 days
func NewItolUploader(uploadid, projectname string, annotationfiles ...string) *ItolUploader {
	return &ItolUploader{uploadid, projectname, annotationfiles}
}

func (u *ItolUploader) Upload(name string, t *tree.Tree) (string, string, error) {
	var client http.Client
	var body *bytes.Buffer
	var writer *multipart.Writer
	var err error
	var formwriter io.Writer
	var zipbytes []byte
	var postrequest *http.Request
	var postresponse *http.Response
	var responsebody []byte
	var responsebodystring string

	var responsesplit []string
	var errorRegexp *regexp.Regexp
	var successRegexp *regexp.Regexp

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	if formwriter, err = writer.CreateFormFile("zipFile", "tree.zip"); err != nil {
		return "", "", err
	}

	if zipbytes, err = compressTree(name+".tree", t.Newick(), u.annotationfiles); err != nil {
		return "", "", err
	}
	formwriter.Write(zipbytes)

	if err = writer.WriteField("treeFormat", "newick"); err != nil {
		return "", "", err
	}

	if name != "" {
		if err = writer.WriteField("treeName", name); err != nil {
			return "", "", err
		}
	}

	if u.uploadid != "" && u.projectname != "" {
		if err = writer.WriteField("uploadID", u.uploadid); err != nil {
			return "", "", err
		}
		if err = writer.WriteField("projectName", u.projectname); err != nil {
			return "", "", err
		}
	}

	if err = writer.Close(); err != nil {
		return "", "", err
	}

	if postrequest, err = http.NewRequest("POST", "http://itol.embl.de/batch_uploader.cgi", body); err != nil {
		return "", "", err
	}

	postrequest.Header.Set("Content-Type", writer.FormDataContentType())
	if postresponse, err = client.Do(postrequest); err != nil {
		return "", "", err
	}

	defer postresponse.Body.Close()
	if responsebody, err = ioutil.ReadAll(postresponse.Body); err != nil {
		return "", "", err
	}

	responsebodystring = string(responsebody)

	responsesplit = strings.Split(responsebodystring, "\n")
	if errorRegexp, err = regexp.Compile("^ERR"); err != nil {
		return "", "", err
	}

	if successRegexp, err = regexp.Compile("^SUCCESS: (\\S+)"); err != nil {
		return "", "", err
	}

	sub := successRegexp.FindStringSubmatch(responsesplit[max(len(responsesplit)-2, 0)])
	if len(sub) < 2 {
		sub = successRegexp.FindStringSubmatch(responsesplit[0])
		if len(sub) < 2 {
			if errorRegexp.MatchString(responsesplit[max(len(responsesplit)-2, 0)]) {
				return "", "", errors.New(fmt.Sprintf("Upload failed. iTOL returned the following error message: %s", responsesplit[len(responsesplit)-1]))
			}
		}
	}

	if len(sub) > 1 {
		return "http://itol.embl.de/tree/" + sub[1], responsebodystring, nil
	} else {
		return "", "", errors.New(responsebodystring)
	}
}

// Comprees the given string in file named filename in a zip file
func compressTree(filename, s string, annotationfiles []string) ([]byte, error) {
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	z, e := zw.Create(filename)
	if e != nil {
		return b.Bytes(), e
	}
	if _, err := z.Write([]byte(s)); err != nil {
		return b.Bytes(), err
	}

	for _, ann := range annotationfiles {
		a, e := zw.Create(ann)
		if e != nil {
			return b.Bytes(), e
		}
		dat, err := ioutil.ReadFile(ann)
		if err != nil {
			return b.Bytes(), e
		}
		if _, err := a.Write(dat); err != nil {
			return b.Bytes(), err
		}
	}

	if err := zw.Flush(); err != nil {
		return b.Bytes(), err
	}
	if err := zw.Close(); err != nil {
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}

func max(a, b int) (max int) {
	max = a
	if b > a {
		max = b
	}
	return
}
