package interfaces

type HttpServer interface {
	Start() error
	Stop() error
	SetTransportManager(transportManager TransportManager)
}
