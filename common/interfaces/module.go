package interfaces

import "github.com/MiviaLabs/hati/common/types"

type Module interface {
	SetTransportManager(tm *TransportManager)
	GetTransportManager() *TransportManager
	AddAction(name string, handler types.ActionHandler) error
	Start() error
	BeforeStart(callback func(m *Module))
	BeforeStop(callback func(m *Module))
	Stop() error
}
