package transport

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
)

var (
	ErrInvalidTransportType = errors.New("invalid transport type")
)

const (
	REDIS_TYPE types.TransportType = "redis"
)

type TransportManager struct {
	modules    *map[string]interfaces.Module
	serverName string
	config     TransportManagerConfig
	redis      *Redis
}

type TransportManagerConfig struct {
	Redis RedisConfig `yaml:"redis" json:"redis"`
}

func NewTransportManager(serverName string, config TransportManagerConfig) TransportManager {
	return TransportManager{
		serverName: serverName,
		config:     config,
		redis:      NewRedis(config.Redis),
	}
}

func (tm TransportManager) SetModules(modules *map[string]interfaces.Module) {
	tm.modules = modules
}

func (tm TransportManager) Start() error {
	if tm.config.Redis.On {
		tm.redis.Start()

		if err := tm.Subscribe(REDIS_TYPE, types.CHAN_MESSAGE, tm.ReceiveMessage); err != nil {
			return err
		}

		if err := tm.Subscribe(REDIS_TYPE, types.CHAN_MESSAGE_RESPONSE, tm.ReceiveMessageResponse); err != nil {
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

func (tm TransportManager) Send(transportType types.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) (any, error) {
	err := tm.Publish(transportType, types.CHAN_MESSAGE, targetServer, targetModule, targetAction, payload, waitForResponse)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (tm TransportManager) Publish(transportType types.TransportType, channel types.Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) error {
	switch transportType {
	case REDIS_TYPE:
		{
			msg, err := tm.prepareMessage(targetServer, targetModule, targetAction, payload, waitForResponse)
			if err != nil {
				return err
			}

			msgBytes, err := msg.MarshalMessage()
			if err != nil {
				return err
			}

			if err := tm.redis.Publish(channel, msgBytes); err != nil {
				return err
			}
			return nil
		}
	default:
		return ErrInvalidTransportType
	}
}

func (tm TransportManager) Subscribe(transportType types.TransportType, channel types.Channel, callback func(payload []byte) error) error {
	switch transportType {
	case REDIS_TYPE:
		{
			return tm.redis.Subscribe(channel, callback)
		}
	default:
		return ErrInvalidTransportType
	}
}

func (tm TransportManager) ReceiveMessage(payload []byte) error {
	var message *structs.Message[[]byte] = &structs.Message[[]byte]{}

	err := json.Unmarshal(payload, message)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if message.TargetID != tm.serverName {
		return errors.New("i am not the target")
	}

	// tm.

	fmt.Println(string(payload))
	fmt.Println(message)
	fmt.Println(string(message.Payload))
	// out := ""
	// err = message.UnmarshalPayload(&out)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return err
	// }

	// fmt.Println(out)
	return nil
}

func (tm TransportManager) ReceiveMessageResponse(payload []byte) error {
	fmt.Println("<--- RECEIVE MESSAGE RESPONSE")
	fmt.Println(string(payload))

	return nil
}

func (tm TransportManager) prepareMessage(targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool) (*structs.Message[[]byte], error) {
	msg := &structs.Message[[]byte]{
		FromID:   tm.serverName,
		TargetID: targetServer,
		TargetAction: structs.TargetAction{
			Module: targetModule,
			Action: targetAction,
		},
		Payload:         payload,
		WaitForResponse: waitForResponse,
	}

	return msg, nil
}
