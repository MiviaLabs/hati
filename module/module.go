package module

import (
	"errors"

	"github.com/MiviaLabs/hati/common"
	"github.com/MiviaLabs/hati/log"
)

var (
	ErrModuleExist = errors.New("module exist")
)

type Module struct {
	common.Module
	Name             string
	actions          map[string]*common.Action
	beforeStart      func(m common.Module)
	beforeStop       func(m common.Module)
	transportManager common.TransportManager
}

func New(name string) *Module {
	return &Module{
		Name:    name,
		actions: make(map[string]*common.Action),
	}
}

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) GetActions() *map[string]*common.Action {
	return &m.actions
}
func (m *Module) SetTransportManager(tm common.TransportManager) {
	m.transportManager = tm
}

func (m *Module) GetTransportManager() common.TransportManager {
	return m.transportManager
}

func (m *Module) AddAction(name string, handler common.ActionHandler, route *common.ActionRoute) error {
	if m.actions[name] != nil {
		return ErrModuleExist
	}

	m.actions[name] = &common.Action{
		Name:    name,
		Route:   route,
		Handler: handler,
	}

	return nil
}

func (m *Module) Start() error {
	log.Debug("starting module: " + m.Name)

	if m.beforeStart != nil {
		c := common.Module(m)
		m.beforeStart(c)
	}

	return nil
}

func (m *Module) BeforeStart(callback func(m common.Module)) {
	m.beforeStart = callback
}

func (m *Module) BeforeStop(callback func(m common.Module)) {
	m.beforeStop = callback
}

func (m *Module) Stop() error {
	log.Debug("stopping module: " + m.Name)

	if m.beforeStop != nil {
		c := common.Module(m)
		m.beforeStop(c)
	}

	return nil
}

func (m *Module) CallAction(name string, payload *common.Message[[]byte]) (common.Response, error) {
	if m.actions[name] == nil {
		return nil, errors.New("action does not exist")
	}

	return m.actions[name].Handler(&common.HatiRequest{
		Message: *payload,
	})
}
