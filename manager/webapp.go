package manager

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/teldio-operations/supervisor-go/module"
	"github.com/teldio-operations/supervisor-manager/utils"
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

func (w *WebappModule) Run() error {
	if w.config.Port == 0 {
		return errors.New("webapp did not have a port defined")
	}
	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", w.config.Port),
		Handler: http.HandlerFunc(utils.Spa(w.fs)),
	}
	slog.Info(fmt.Sprintf("serving %s at http://%s", w.config.Title, server.Addr))
	return server.ListenAndServe()
}

func NewWebappModule(path string, config *WebappConfig) (*WebappModule, error) {
	module := WebappModule{config: config, fs: os.DirFS(path)}

	return &module, nil
}
