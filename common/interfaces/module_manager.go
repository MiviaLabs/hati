package interfaces

type ModuleManager interface {
	Start() error
	Stop() error
	AddModule(modules ...Module) error
	GetModule(name string) Module
	GetModules() *map[string]Module
	SetTransportManager(tm TransportManager)
}
