package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

//go:embed dist
var content embed.FS

var Content http.FileSystem

type onlyFilesFS struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func init() {
	contentDir, err := fs.Sub(content, "dist")
	if err != nil {
		panic(fmt.Sprintf("loggo/web: getting web content subdirectory: %s", err))
	}

	Content = http.FS(contentDir)
}
