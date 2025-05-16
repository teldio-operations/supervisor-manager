package main

import (
	"log/slog"

	"github.com/teldio-operations/supervisor-manager/manager"
)

func DetectAndRunModules() {
	slog.Debug("all modules registered")

	manager, err := manager.NewManager(*repositoryPath, *configsPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	manager.StartModules()
}
