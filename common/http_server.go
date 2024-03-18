package common

// HttpServer interface for http server
type HttpServer interface {
	Start() error
	Stop() error
	SetTransportManager(transportManager TransportManager)
	SetModules(modules *map[string]Module)
}
