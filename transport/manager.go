package transport

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/MiviaLabs/hati/common"
	"github.com/MiviaLabs/hati/log"
)

var (
	ErrInvalidTransportType = errors.New("invalid transport type")
)

const (
	REDIS_TYPE common.TransportType = "redis"
)

type TransportManager struct {
	modules                map[string]common.Module
	serverName             string
	config                 TransportManagerConfig
	redis                  *Redis
	httpServer             *HttpServer
	moduleManager          common.ModuleManager
	waitingForResponse     map[string]chan common.Message[[]byte]
	waitingForResponseLock sync.Mutex
}

type TransportManagerConfig struct {
	Redis RedisConfig `yaml:"redis" json:"redis"`
	Http  HttpConfig  `yaml:"http" json:"http"`
}

func NewTransportManager(serverName string, config TransportManagerConfig, moduleManager common.ModuleManager) *TransportManager {
	tm := &TransportManager{
		serverName:         serverName,
		config:             config,
		redis:              NewRedis(config.Redis),
		httpServer:         NewHttpServer(config.Http),
		moduleManager:      moduleManager,
		waitingForResponse: make(map[string]chan common.Message[[]byte], 10),
	}
	moduleManager.SetTransportManager(tm)

	return tm
}

func (tm *TransportManager) GetRedis() common.Redis {
	return tm.redis
}

func (tm *TransportManager) SetModules(modules map[string]common.Module) {
	tm.modules = modules
}

func (tm *TransportManager) Start() error {
	if tm.config.Redis.On {
		if err := tm.redis.Start(); err != nil {
			return err
		}

		if err := tm.Subscribe(REDIS_TYPE, common.CHAN_MESSAGE, tm.ReceiveMessage); err != nil {
			return err
		}

		if err := tm.Subscribe(REDIS_TYPE, common.CHAN_MESSAGE_RESPONSE, tm.ReceiveMessageResponse); err != nil {
			return err
		}
	}

	if tm.config.Http.On {
		tm.httpServer.SetModules(tm.moduleManager.GetModules())
		tm.httpServer.SetTransportManager(tm)

		if err := tm.httpServer.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (tm *TransportManager) Stop() error {
	log.Debug("stopping transport manager")

	if tm.config.Http.On {
		if err := tm.httpServer.Stop(); err != nil {
			return err
		}
	}

	if tm.config.Redis.On {
		if err := tm.redis.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (tm *TransportManager) Send(transportType common.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error) {
	res, err := tm.Publish(transportType, common.CHAN_MESSAGE, targetServer, targetModule, targetAction, payload, waitForResponse, responseHash)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (tm *TransportManager) SendResponse(transportType common.TransportType, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error) {
	res, err := tm.Publish(transportType, common.CHAN_MESSAGE_RESPONSE, targetServer, targetModule, targetAction, payload, waitForResponse, responseHash)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (tm *TransportManager) Publish(transportType common.TransportType, channel common.Channel, targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (any, error) {
	switch transportType {
	case REDIS_TYPE:
		{
			msg, err := tm.prepareMessage(targetServer, targetModule, targetAction, payload, waitForResponse, responseHash)
			if err != nil {
				return nil, err
			}

			msg.UpdateHash()

			msgBytes, err := msg.MarshalMessage()
			if err != nil {
				return nil, err
			}

			if err := tm.redis.Publish(channel, msgBytes); err != nil {
				return nil, err
			}

			if msg.WaitForResponse {
				timer := time.NewTimer(time.Duration(1) * time.Second)

				tm.waitingForResponseLock.Lock()
				tm.waitingForResponse[msg.Hash] = make(chan common.Message[[]byte])
				tm.waitingForResponseLock.Unlock()

			Loop:
				for {

					select {
					case msg := <-tm.waitingForResponse[msg.Hash]:

						if msg.FromID == "" {
							continue
						}

						tm.waitingForResponseLock.Lock()
						defer tm.waitingForResponseLock.Unlock()

						if tm.waitingForResponse[msg.Hash] != nil {
							close(tm.waitingForResponse[msg.Hash])
						}
						tm.waitingForResponse[msg.Hash] = nil

						return msg, nil
					case <-timer.C:
						timer.Stop()
						tm.waitingForResponseLock.Lock()

						close(tm.waitingForResponse[msg.Hash])
						tm.waitingForResponse[msg.Hash] = nil

						tm.waitingForResponseLock.Unlock()
						break Loop
					}
				}
			}

			return nil, nil
		}
	default:
		return nil, ErrInvalidTransportType
	}
}

func (tm *TransportManager) Subscribe(transportType common.TransportType, channel common.Channel, callback func(payload []byte) (common.Response, error)) error {
	switch transportType {
	case REDIS_TYPE:
		{
			return tm.redis.Subscribe(channel, callback)
		}
	default:
		return ErrInvalidTransportType
	}
}

func (tm *TransportManager) ReceiveMessage(payload []byte) (common.Response, error) {
	var message *common.Message[[]byte] = &common.Message[[]byte]{}

	err := json.Unmarshal(payload, message)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if message.TargetID != tm.serverName {
		return nil, nil
	}

	module := (tm.moduleManager).GetModule(message.TargetAction.Module)
	if module == nil {
		log.Warning("module does not exist")
		return nil, errors.New("module does not exist")
	}

	if !message.WaitForResponse {
		_, err := module.CallAction(message.TargetAction.Action, message)
		if err != nil {
			log.Warning(err.Error())

			return nil, err
		}

		return nil, nil
	}

	response, err := module.CallAction(message.TargetAction.Action, message)
	if err != nil {
		log.Warning(err.Error())

		return nil, err
	}

	s := response.(string)
	tm.SendResponse(REDIS_TYPE, message.FromID, "", "", []byte(s), false, message.Hash)

	return nil, nil
}

func (tm *TransportManager) ReceiveMessageResponse(payload []byte) (common.Response, error) {
	var message *common.Message[[]byte] = &common.Message[[]byte]{}

	err := json.Unmarshal(payload, message)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if message.TargetID != tm.serverName {
		return nil, nil
	}

	tm.waitingForResponseLock.Lock()
	defer tm.waitingForResponseLock.Unlock()

	if tm.waitingForResponse[message.ResponseHash] != nil {
		tm.waitingForResponse[message.ResponseHash] <- *message
	}

	return nil, nil
}

func (tm *TransportManager) prepareMessage(targetServer string, targetModule string, targetAction string, payload []byte, waitForResponse bool, responseHash string) (*common.Message[[]byte], error) {
	msg := &common.Message[[]byte]{
		FromID:   tm.serverName,
		TargetID: targetServer,
		TargetAction: common.TargetAction{
			Module: targetModule,
			Action: targetAction,
		},
		Payload:         payload,
		WaitForResponse: waitForResponse,
		ResponseHash:    responseHash,
	}

	msg.UpdateHash()

	return msg, nil
}
