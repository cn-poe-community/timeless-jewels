package main

import (
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"path"
	"time"
)

func FileServer(root string) http.Handler {
	return &fileHandler{
		handler: http.FileServer(http.Dir(root)),
		root:    root,
	}
}

type fileHandler struct {
	handler http.Handler
	root    string
}

func (h *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join(h.root, r.URL.Path)
	stat, err := os.Stat(filePath)
	if err == nil {
		if stat.IsDir() {
			indexFile := path.Join(filePath, "index.html")
			stat, err = os.Stat(indexFile)
			if err == nil {
				filePath = indexFile
			}
		}
	} else {
		if errors.Is(err, os.ErrNotExist) {
			htmlFile := filePath + ".html"
			stat, err = os.Stat(htmlFile)
			if err == nil {
				filePath = htmlFile
				r.URL.Path = r.URL.Path + ".html"
			}
		}
	}

	if stat != nil && !stat.IsDir() {
		etag, err := h.getEtag(filePath, stat.ModTime())
		if err == nil {
			w.Header().Set("Etag", etag)
		}
	}

	w.Header().Set("Cache-Control", "no-cache")
	h.handler.ServeHTTP(w, r)
}

func (h *fileHandler) getEtag(filePath string, modtime time.Time) (string, error) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d%d", modtime.Unix(), crc32.ChecksumIEEE(dat)), nil
}
