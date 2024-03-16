package core

import (
	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
	"github.com/MiviaLabs/hati/transport"
	"github.com/joho/godotenv"
)

type Hati struct {
	config           Config
	moduleManager    interfaces.ModuleManager
	transportManager interfaces.TransportManager
	stopChan         chan bool
}

func NewHati(config Config) Hati {

	hati := Hati{
		config:   config,
		stopChan: make(chan bool),
	}

	hati.moduleManager = module.NewModuleManager()
	hati.transportManager = transport.NewTransportManager(config.Name, config.Transport, hati.moduleManager)

	return hati
}

func (h Hati) AddModule(modules ...interfaces.Module) error {
	return h.moduleManager.AddModule(modules...)
}

func (h Hati) Start() error {
	log.Default("hati v0.1.0")
	log.Debug("starting hati")

	_ = godotenv.Load()

	if err := h.transportManager.Start(); err != nil {
		return err
	}

	if err := h.moduleManager.Start(); err != nil {
		return err
	}

	return nil
}

func (h Hati) Stop() error {
	log.Debug("stopping hati")

	if err := h.transportManager.Stop(); err != nil {
		return err
	}

	if err := h.moduleManager.Stop(); err != nil {
		return err
	}

	return nil
}

func (h Hati) TransportManager() *interfaces.TransportManager {
	return &h.transportManager
}
