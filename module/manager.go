package module

import (
	"errors"
	"github.com/MiviaLabs/hati/common"
)

type Manager struct {
	common.ModuleManager
	modules          map[string]common.Module
	transportManager common.TransportManager
}

func NewModuleManager() *Manager {
	return &Manager{
		modules: make(map[string]common.Module),
	}
}

func (m *Manager) AddModule(modules ...common.Module) error {
	for _, module := range modules {
		if m.modules[module.GetName()] != nil && common.Module(m.modules[module.GetName()]).GetName() == module.GetName() {
			return errors.New("module " + module.GetName() + " already exist")
		}

		module.SetTransportManager(m.transportManager)

		m.modules[module.GetName()] = module
	}

	return nil
}

func (m *Manager) GetModules() *map[string]common.Module {
	return &m.modules
}

func (m *Manager) SetTransportManager(tm common.TransportManager) {
	m.transportManager = tm
}

func (m *Manager) GetModule(name string) common.Module {
	return m.modules[name]
}

func (m *Manager) Start() error {
	for _, module := range m.modules {
		if err := module.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Stop() error {
	for _, module := range m.modules {
		if err := module.Stop(); err != nil {
			return err
		}
	}
	return nil
}
