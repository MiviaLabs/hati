package interfaces

import (
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
)

type Module interface {
	SetTransportManager(tm TransportManager)
	GetTransportManager() TransportManager
	AddAction(name string, handler types.ActionHandler, route *structs.ActionRoute) error
	GetActions() *map[string]*types.Action
	CallAction(name string, payload *structs.Message[[]byte]) (types.Response, error)
	Start() error
	BeforeStart(callback func(m Module))
	BeforeStop(callback func(m Module))
	Stop() error
	GetName() string
}
