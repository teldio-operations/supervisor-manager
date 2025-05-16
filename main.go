package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/teldio-operations/supervisor-manager/api"
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

//go:embed webui/dist
var dist embed.FS

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	ensureFlag(configsPath)
	ensureFlag(repositoryPath)

	go DetectAndRunModules()

	dist, _ := fs.Sub(dist, "webui/dist")
	api.StartServer(dist)
}
