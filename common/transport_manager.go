package common

type TransportManager interface {
	Start() error
	Stop() error
	SetModules(modules map[string]Module)
	Send(transportType TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	SendResponse(transportType TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	Publish(transportType TransportType, channel Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error)
	Subscribe(transportType TransportType, channel Channel, callback func(payload []byte) (Response, error)) error
	ReceiveMessage(payload []byte) (Response, error)
	ReceiveMessageResponse(payload []byte) (Response, error)
	GetRedis() Redis
	// SetHttpRoute(methods) error
}
