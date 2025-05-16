package api

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/teldio-operations/supervisor-manager/utils"
)

func Webui(dist fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "development" {
			httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   "localhost:23497",
			}).ServeHTTP(w, r)
		} else {
			utils.Spa(dist)(w, r)
		}
	}
}
