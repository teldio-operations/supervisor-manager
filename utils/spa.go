package utils

import (
	"io/fs"
	"net/http"
	"strings"
)

func Spa(ffs fs.FS) http.HandlerFunc {
	fileServer := http.FileServerFS(ffs)
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Stat(ffs, strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			http.ServeFileFS(w, r, ffs, "index.html")
			return
		}

		fileServer.ServeHTTP(w, r)
	}
}
