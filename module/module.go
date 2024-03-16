package module

import (
	"errors"

	"github.com/MiviaLabs/hati/log"
)

type ActionHandler func(payload any) (any, error)

var (
	ErrModuleExist = errors.New("module exist")
)

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

func (m Module) AddAction(name string, handler ActionHandler) error {
	if m.actions[name] != nil {
		return ErrModuleExist
	}

	m.actions[name] = handler

	return nil
}

func (m Module) Start() error {
	log.Debug("starting module: " + m.Name)

	return nil
}

func (m Module) Stop() error {
	log.Debug("stopping module: " + m.Name)

	return nil
}
