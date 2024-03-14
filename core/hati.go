package core

import (
	"errors"

	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
	"github.com/MiviaLabs/hati/transport"
)

type Hati struct {
	config           Config
	modules          map[string]module.Module
	transportManager transport.TransportManager
	stopChan         chan bool
}

func NewHati(config Config) Hati {
	return Hati{
		config:           config,
		modules:          make(map[string]module.Module),
		transportManager: transport.NewTransportManager(config.Transport),
		stopChan:         make(chan bool),
	}
}

func (h Hati) AddModule(modules ...module.Module) error {
	for _, module := range modules {
		if h.modules[module.Name].Name == "" {
			return errors.New("module " + module.Name + " already exist")
		}

		h.modules[module.Name] = module
	}

	return nil
}

func (h Hati) Start() error {
	log.Default("hati v0.1.0")
	log.Debug("starting hati")

	if err := h.transportManager.Start(); err != nil {
		return err
	}

	return nil
}

func (h Hati) Stop() error {
	log.Debug("stopping hati")

	if err := h.transportManager.Stop(); err != nil {
		return err
	}

	return nil
}
