package transport

import "github.com/MiviaLabs/hati/common"

type Transport interface {
	Start() error
	Stop() error
	Send(transportType common.TransportType, target string, payload []byte) error
	Publish(channel common.Channel, payload []byte) error
	Subscribe(channel common.Channel, callback func(payload []byte)) error
}
