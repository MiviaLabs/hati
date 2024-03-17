package interfaces

type HttpManager interface {
	Start() error
	Stop() error
	SetTransportManager(transportManager TransportManager)
}
