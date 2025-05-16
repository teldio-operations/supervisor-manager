package api

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func StartServer(dist fs.FS) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Webui(dist))
	mux.HandleFunc("/*", Webui(dist))

	api := humago.New(mux, huma.DefaultConfig("Manager", "0.1.0"))

	type B struct {
		Hello string
	}
	type O struct {
		Body B
	}
	huma.Get(api, "/jesus", func(ctx context.Context, i *B) (*O, error) {
		return &O{Body: B{"world"}}, nil
	})

	slog.Info("Serving manager at http://localhost:30605")
	slog.Error(http.ListenAndServe(":30605", mux).Error())
}
