package manager

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	repository *ModulesRepository
	instances  []Module
	wg         sync.WaitGroup
}

func NewManager(repositoryPath string, configsPath string) (*Manager, error) {
	repository, err := GetRegisteredModules(repositoryPath)
	if err != nil {
		return nil, err
	}
	manager := &Manager{
		repository: &repository,
	}
	configs, err := os.ReadDir(configsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configs path: %s", err)
	}
	for _, configFile := range configs {
		path := filepath.Join(configsPath, configFile.Name())
		fileInfo, err := configFile.Info()
		if err != nil {
			continue
		}

		if filepath.Ext(configFile.Name()) == ".json" {
			module, err := manager.repository.NewExecutableModuleFromConfig(path, "")
			if err != nil {
				slog.Error(fmt.Sprintf("attempted to load config %s but failed: %s", configFile.Name(), err))
				continue
			}
			manager.instances = append(manager.instances, module)
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
				var webappConfig *WebappConfig
				err = json.Unmarshal(configBytes, &webappConfig)
				if err != nil {
					continue
				}
				module, err := NewWebappModule(path, webappConfig)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to register webapp module: %s", err))
					continue
				}
				manager.instances = append(manager.instances, module)
			} else {
				module, err := manager.repository.NewExecutableModuleFromConfig(filepath.Join(path, "config.json"), path)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to register executable module: %s", err))
					continue
				}
				manager.instances = append(manager.instances, module)
			}
		}
	}
	return manager, nil
}

func (manager *Manager) startModule(module Module) {
	manager.wg.Add(1)
	go func() {
		defer manager.wg.Done()
		err := module.Run()
		if err != nil {
			slog.Error(fmt.Sprintf("module %s quit with error: %s", module.Info().Title, err))
		}
	}()

}

func (manager *Manager) StartModules() {
	for _, module := range manager.instances {
		manager.startModule(module)
	}

	manager.wg.Wait()
}

func (manager *Manager) StartModule(config string) error {
	module, err := manager.repository.NewExecutableModule(config, "")
	if err != nil {
		return err
	}
	manager.instances = append(manager.instances, module)
	manager.startModule(module)
	return nil
}
