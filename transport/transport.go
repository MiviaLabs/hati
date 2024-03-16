package transport

import "github.com/MiviaLabs/hati/common/types"

type Transport interface {
	Start() error
	Stop() error
	Send(transportType types.TransportType, target string, payload []byte) error
	Publish(channel types.Channel, payload []byte) error
	Subscribe(channel types.Channel, callback func(payload []byte)) error
}

// type Transport struct{}

// func New() {}
