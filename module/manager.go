package module

import (
	"errors"

	"github.com/MiviaLabs/hati/common/interfaces"
)

type Manager struct {
	interfaces.ModuleManager
	modules          map[string]interfaces.Module
	transportManager interfaces.TransportManager
}

func NewModuleManager() *Manager {
	return &Manager{
		modules: make(map[string]interfaces.Module),
	}
}

func (m *Manager) AddModule(modules ...interfaces.Module) error {
	for _, module := range modules {
		if m.modules[module.GetName()] != nil && interfaces.Module(m.modules[module.GetName()]).GetName() == module.GetName() {
			return errors.New("module " + module.GetName() + " already exist")
		}

		module.SetTransportManager(m.transportManager)

		m.modules[module.GetName()] = module
	}

	return nil
}

func (m *Manager) GetModules() *map[string]interfaces.Module {
	return &m.modules
}

func (m *Manager) SetTransportManager(tm interfaces.TransportManager) {
	m.transportManager = tm
}

func (m *Manager) GetModule(name string) interfaces.Module {
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
