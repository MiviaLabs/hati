package module

type ActionHandler func(payload any) (any, error)

type IModule interface{}

type Module struct {
	IModule
	Name    string
	actions map[string]ActionHandler
}

func New(name string) Module {
	return Module{
		Name:    name,
		actions: make(map[string]ActionHandler),
	}
}
