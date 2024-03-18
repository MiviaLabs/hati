package common

import (
	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Start() error
	Publish(channel Channel, payload []byte) error
	Subscribe(channel Channel, callback func(payload []byte) (Response, error)) error
	Stop() error
	Client() *redis.Client
	RmqClient() *rmq.Connection
}
