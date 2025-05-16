package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/teldio-operations/supervisor-go/module"
)

type ExecutableModule struct {
	info           *module.Info
	executablePath string
}

type ModulesRepository []*ExecutableModule

func GetRegisteredModules(repositoryPath string) (ModulesRepository, error) {
	files, err := os.ReadDir(repositoryPath)
	if err != nil {
		return nil, err
	}

	var modules []*ExecutableModule
	for _, file := range files {
		path := filepath.Join(repositoryPath, file.Name())

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		command := exec.CommandContext(ctx, path, "-info")
		infoBytes, err := command.Output()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		var info module.Info
		err = json.Unmarshal(infoBytes, &info)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to unmarshal json for %s: %s", file.Name(), err))
			continue
		}

		modules = append(modules, &ExecutableModule{
			info:           &info,
			executablePath: path,
		})
	}
	return modules, nil
}

type ExecutableModuleInstance struct {
	*ExecutableModule
	config        string
	directoryPath string
}

func (w *ExecutableModuleInstance) Info() *module.Info {
	return w.info
}

func (w *ExecutableModuleInstance) Run() error {
	command := exec.Cmd{
		Path: w.executablePath,
		Dir:  w.directoryPath,
	}
	slog.Info(fmt.Sprintf("running module %s", w.info.Name))
	command.Stdin = strings.NewReader(w.config)
	return command.Run()
}

func (repo *ModulesRepository) NewExecutableModuleFromConfig(configPath, directoryPath string) (*ExecutableModuleInstance, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return repo.NewExecutableModule(string(configBytes), directoryPath)
}

func (repo *ModulesRepository) NewExecutableModule(configStr string, directoryPath string) (*ExecutableModuleInstance, error) {
	var config struct {
		Name string `json:"module"`
	}
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return nil, err
	}

	for _, module := range *repo {
		if module.info.Name == config.Name {
			return &ExecutableModuleInstance{
				ExecutableModule: module,
				config:           configStr,
				directoryPath:    directoryPath,
			}, nil
		}
	}

	return nil, errors.New("no module with that name found")
}
