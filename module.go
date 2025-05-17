package main

import (
	"log/slog"

	"github.com/teldio-operations/supervisor-manager/manager"
)

func DetectAndRunModules() {
	manager, err := manager.NewManager(*repositoryPath, *configsPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	manager.StartModules()
}
