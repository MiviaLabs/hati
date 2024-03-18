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
	// CreateQueue(name string, consumer interface{}) error
	// FindQueue(name string) error
	// StartQueue(name string) error
	// StopQueue(name string) error
}
