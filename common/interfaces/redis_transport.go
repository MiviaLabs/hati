package interfaces

import "github.com/MiviaLabs/hati/common/types"

type Redis interface {
	Start() error
	Publish(channel types.Channel, payload []byte) error
	Subscribe(channel types.Channel, callback func(payload []byte) (types.Response, error)) error
	Stop() error
}
