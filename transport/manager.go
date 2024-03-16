package transport

import (
	"fmt"

	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
)

type TransportManager struct {
	config TransportManagerConfig
	redis  *Redis
}

type TransportManagerConfig struct {
	Redis RedisConfig `yaml:"redis" json:"redis"`
}

func NewTransportManager(config TransportManagerConfig) TransportManager {
	return TransportManager{
		config: config,
		redis:  NewRedis(config.Redis),
	}
}

func (tm TransportManager) Start() error {

	if tm.config.Redis.On {
		tm.redis.Start()

		if err := tm.Subscribe(types.CHAN_MESSAGE, tm.ReceiveMessage); err != nil {
			return err
		}

		if err := tm.Subscribe(types.CHAN_MESSAGE_RESPONSE, tm.ReceiveMessageResponse); err != nil {
			return err
		}
	}

	return nil
}

func (tm TransportManager) Stop() error {
	log.Debug("stopping transport manager")

	if tm.config.Redis.On {
		if err := tm.redis.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (tm TransportManager) Subscribe(channel types.Channel, callback func(payload []byte) error) error {
	return tm.redis.Subscribe(channel, callback)
}

func (tm TransportManager) ReceiveMessage(payload []byte) error {
	fmt.Println("---> RECEIVE MESSAGE")
	fmt.Println(string(payload))

	return nil
}

func (tm TransportManager) ReceiveMessageResponse(payload []byte) error {
	fmt.Println("<--- RECEIVE MESSAGE RESPONSE")
	fmt.Println(string(payload))

	return nil
}
