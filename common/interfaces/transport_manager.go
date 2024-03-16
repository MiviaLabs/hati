package interfaces

import "github.com/MiviaLabs/hati/common/types"

type TransportManager interface {
	Start() error
	Stop() error
	SetModules(modules *map[string]Module)
	Send(transportType types.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) (any, error)
	Publish(transportType types.TransportType, channel types.Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) error
	Subscribe(transportType types.TransportType, channel types.Channel, callback func(payload []byte) error) error
	ReceiveMessage(payload []byte) error
	ReceiveMessageResponse(payload []byte) error
}
