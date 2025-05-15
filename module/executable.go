package module

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

type ExecutableModule struct {
	info *Info

	directoryPath  string
	executablePath string
}

func (w *ExecutableModule) Info() *Info {
	return w.info
}

func (w *ExecutableModule) Execute() error {
	command := exec.Cmd{
		Path: w.executablePath,
		Args: []string{"module.json"},
		Dir:  w.directoryPath,
	}
	return command.Run()
}

func NewExecutableModuleFromDirectory(directoryPath string) (*ExecutableModule, error) {
	openapiBytes, err := os.ReadFile(filepath.Join(directoryPath, "openapi.json"))
	if err != nil {
		return nil, err
	}

	module := ExecutableModule{directoryPath: directoryPath}

	err = json.Unmarshal(openapiBytes, &module.info)
	if err != nil {
		return nil, err
	}

	module.executablePath = filepath.Join(directoryPath, module.info.Name)

	return &module, err
}
