package transport

import "github.com/MiviaLabs/hati/common/types"

type Transport interface {
	Start() error
	Stop() error
	Send(payload []byte) error
	// Send(channel Channel, payload []byte) error
	Subscribe(channel types.Channel, callback func(payload []byte)) error
}

// type Transport struct{}

// func New() {}
