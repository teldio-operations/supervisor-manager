package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/teldio-operations/supervisor-manager/module"
)

var configsEnvVarName = "SUPERVISOR_CONFIGS_PATH"
var repositoryEnvVarName = "SUPERVISOR_REPOSITORY_PATH"

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
	flag.Parse()

	ensureFlag(configsPath)
	ensureFlag(repositoryPath)

	configs, err := os.ReadDir(*configsPath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read configs path: %s", err))
		return
	}

	var modules []module.Module

	for _, configFile := range configs {
		path := filepath.Join(*configsPath, configFile.Name())
		fileInfo, err := configFile.Info()
		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

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
			var webappConfig *module.WebappConfig
			err = json.Unmarshal(configBytes, &webappConfig)
			if err != nil {
				continue
			}
			module, err := module.NewWebappModule(path, webappConfig)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to register webapp module: %s", err))
				continue
			}
			modules = append(modules, module)
		} else {
			module, err := module.NewExecutableModuleFromDirectory(path)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to register executable module: %s", err))
				continue
			}
			modules = append(modules, module)
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
