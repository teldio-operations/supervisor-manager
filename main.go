package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/teldio-operations/supervisor-manager/manager"
)

var configsEnvVarName = "SUPERVISOR_CONFIGS_PATH"
var repositoryEnvVarName = "SUPERVISOR_MODULES_PATH"

var configsPath = flag.String(
	"configs-path",
	os.Getenv(configsEnvVarName),
	fmt.Sprintf("The path where module configs are stored (may also be set by env var %s)", configsEnvVarName))

var repositoryPath = flag.String(
	"repository-path",
	os.Getenv(repositoryEnvVarName),
	fmt.Sprintf("The path where modules executables are stored (may also be set by env var %s)", repositoryEnvVarName))

func ensureFlag[T comparable](value *T) {
	if value == nil {
		flag.Usage()
		os.Exit(1)
	}
	var empty T
	if *value == empty {
		flag.Usage()
		os.Exit(1)
	}
}

func isExecutable(mode fs.FileMode) bool {
	return !mode.IsDir() && mode&0111 != 0
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	ensureFlag(configsPath)
	ensureFlag(repositoryPath)

	configs, err := os.ReadDir(*configsPath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read configs path: %s", err))
		return
	}

	var modules []manager.Module

	modulesRepo, err := manager.GetRegisteredModules(*repositoryPath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieve module repository: %s", err))
		return
	}

	slog.Debug("all modules registered")

	for _, configFile := range configs {
		path := filepath.Join(*configsPath, configFile.Name())
		fileInfo, err := configFile.Info()
		if err != nil {
			continue
		}

		if filepath.Ext(configFile.Name()) == ".json" {
			module, err := modulesRepo.NewExecutableModuleFromConfig(path, "")
			if err != nil {
				slog.Error(fmt.Sprintf("attempted to load config %s but failed: %s", configFile.Name(), err))
				continue
			}
			modules = append(modules, module)
		}

		if fileInfo.IsDir() {
			configBytes, err := os.ReadFile(filepath.Join(path, "config.json"))
			if err != nil {
				continue
			}
			var config map[string]any
			err = json.Unmarshal(configBytes, &config)
			if err != nil {
				continue
			}

			if config["name"] == "webapp" || config["name"] == "web" {
				var webappConfig *manager.WebappConfig
				err = json.Unmarshal(configBytes, &webappConfig)
				if err != nil {
					continue
				}
				module, err := manager.NewWebappModule(path, webappConfig)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to register webapp module: %s", err))
					continue
				}
				modules = append(modules, module)
			} else {
				module, err := modulesRepo.NewExecutableModuleFromConfig(filepath.Join(path, "config.json"), path)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to register executable module: %s", err))
					continue
				}
				modules = append(modules, module)
			}
		}
	}

	var wg sync.WaitGroup

	for _, module := range modules {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := module.Execute()
			if err != nil {
				slog.Error(fmt.Sprintf("module %s quit with error: %s", module.Info().Title, err))
			}
		}()
	}

	wg.Wait()
}
