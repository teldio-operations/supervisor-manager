package manager

import "github.com/teldio-operations/supervisor-go/module"

type Module interface {
	Info() *module.Info
	Run() error
}
