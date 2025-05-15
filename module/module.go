package module

type Module interface {
	Info() *Info
	Execute() error
}
