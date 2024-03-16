package module

import (
	"errors"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
)

var (
	ErrModuleExist = errors.New("module exist")
)

type Module struct {
	interfaces.Module
	Name             string
	actions          map[string]types.ActionHandler
	beforeStart      func(m interfaces.Module)
	beforeStop       func(m interfaces.Module)
	transportManager interfaces.TransportManager
}

func New(name string) *Module {
	return &Module{
		Name:    name,
		actions: make(map[string]types.ActionHandler),
	}
}

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) SetTransportManager(tm interfaces.TransportManager) {
	m.transportManager = tm
}

func (m *Module) GetTransportManager() interfaces.TransportManager {
	return m.transportManager
}

func (m *Module) AddAction(name string, handler types.ActionHandler) error {
	if m.actions[name] != nil {
		return ErrModuleExist
	}

	m.actions[name] = handler

	return nil
}

func (m *Module) Start() error {
	log.Debug("starting module: " + m.Name)

	if m.beforeStart != nil {
		c := interfaces.Module(m)
		m.beforeStart(c)
	}

	return nil
}

func (m *Module) BeforeStart(callback func(m interfaces.Module)) {
	m.beforeStart = callback
}

func (m *Module) BeforeStop(callback func(m interfaces.Module)) {
	m.beforeStop = callback
}

func (m *Module) Stop() error {
	log.Debug("stopping module: " + m.Name)

	if m.beforeStop != nil {
		c := interfaces.Module(m)
		m.beforeStop(c)
	}

	return nil
}

func (m *Module) CallAction(name string, payload *structs.Message[[]byte]) (types.Response, error) {
	if m.actions[name] == nil {
		return nil, errors.New("action does not exist")
	}

	return m.actions[name](*payload)
}
