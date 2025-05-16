package manager

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/teldio-operations/supervisor-go/module"
)

type WebappModule struct {
	config *WebappConfig
	fs     fs.FS
}

type WebappConfig struct {
	Title string
	Port  int
}

func (w *WebappModule) Info() *module.Info {
	return &module.Info{
		BaseInfo: module.BaseInfo{
			Name:  "webapp",
			Title: w.config.Title,
		},
	}
}

func (mod *WebappModule) spa(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Stat(mod.fs, strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			http.ServeFileFS(w, r, mod.fs, "index.html")
			return
		}

		fileServer.ServeHTTP(w, r)
	}
}

func (w *WebappModule) Execute() error {
	if w.config.Port == 0 {
		return errors.New("webapp did not have a port defined")
	}
	fileServer := http.FileServerFS(w.fs)
	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", w.config.Port),
		Handler: http.HandlerFunc(w.spa(fileServer)),
	}
	slog.Info(fmt.Sprintf("serving %s at http://%s", w.config.Title, server.Addr))
	return server.ListenAndServe()
}

func NewWebappModule(path string, config *WebappConfig) (*WebappModule, error) {
	module := WebappModule{config: config, fs: os.DirFS(path)}

	return &module, nil
}
