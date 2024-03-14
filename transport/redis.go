package transport

import "github.com/MiviaLabs/hati/log"

type RedisConfig struct {
	On       bool   `yaml:"on" json:"on"`
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database int    `yaml:"database" json:"database"`
}

type Redis struct {
	Transport
}

func NewRedis(config RedisConfig) Redis {
	return Redis{}
}

func (r Redis) Start() error {
	log.Debug("starting redis")

	return nil
}

func (r Redis) Stop() error {
	log.Debug("stopping redis")

	return nil
}
