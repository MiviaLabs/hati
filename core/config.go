package core

import "github.com/MiviaLabs/hati/transport"

type Config struct {
	// Name of this hati instance
	Name string `yaml:"name" json:"name"`

	Transport transport.TransportManagerConfig `yaml:"transport" json:"transport"`
}
