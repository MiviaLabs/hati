package common

type Module interface {
	SetTransportManager(tm TransportManager)
	GetTransportManager() TransportManager
	AddAction(name string, handler ActionHandler, route *ActionRoute) error
	GetActions() *map[string]*Action
	CallAction(name string, payload *Message[[]byte]) (Response, error)
	Start() error
	BeforeStart(callback func(m Module))
	BeforeStop(callback func(m Module))
	Stop() error
	GetName() string
}
