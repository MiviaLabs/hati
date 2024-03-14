package transport

import "github.com/MiviaLabs/hati/log"

type TransportManager struct {
	config TransportManagerConfig
	redis  Redis
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
	}

	return nil
}

func (tm TransportManager) Stop() error {
	log.Debug("stopping transport manager")
	return nil
}
