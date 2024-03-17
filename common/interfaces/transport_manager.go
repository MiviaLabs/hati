package interfaces

import "github.com/MiviaLabs/hati/common/types"

type TransportManager interface {
	Start() error
	Stop() error
	SetModules(modules map[string]Module)
	Send(transportType types.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	SendResponse(transportType types.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	Publish(transportType types.TransportType, channel types.Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	Subscribe(transportType types.TransportType, channel types.Channel, callback func(payload []byte) (types.Response, error)) error
	ReceiveMessage(payload []byte) (types.Response, error)
	ReceiveMessageResponse(payload []byte) (types.Response, error)
	GetRedis() Redis
	// SetHttpRoute(methods) error
}
