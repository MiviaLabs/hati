package transport

type Transport interface {
	Start() error
	Stop() error
	// Send(channel Channel, payload []byte) error
	// Subscribe(channel Channel, callback func(payload []byte)) error
}

// type Transport struct{}

// func New() {}
