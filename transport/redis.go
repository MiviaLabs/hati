package transport

import (
	"context"
	"errors"
	"sync"

	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
	"github.com/adjust/rmq/v5"
	redis "github.com/redis/go-redis/v9"
)

var (
	ErrRedisAlreadySubscribed = func(channel types.Channel) error {
		return errors.New("already subscribed to this channel: " + string(channel))
	}
)

type RedisConfig struct {
	On             bool   `yaml:"on" json:"on"`
	Host           string `yaml:"host" json:"host"`
	Port           string `yaml:"port" json:"port"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	Database       int    `yaml:"database" json:"database"`
	Protocol       int    `yaml:"protocol" json:"protocol"`
	PoolSize       int    `yaml:"pool_size" json:"pool_size"`
	MaxActiveConns int    `yaml:"max_active_conns" json:"max_active_conns"`
}

type Redis struct {
	Transport
	Client           *redis.Client
	config           RedisConfig
	publisherClient  *redis.Client
	subscriberClient *redis.Client
	publishers       map[string]*RedisPublisher
	subscribers      map[string]*RedisSubscriber
	publishersWg     sync.WaitGroup
	subscribersWg    sync.WaitGroup
	// queue         map[string]rmq.Queue
	// queuePrefix   string
	rmqConnection rmq.Connection
}

func NewRedis(config RedisConfig) *Redis {
	return &Redis{
		config:      config,
		publishers:  make(map[string]*RedisPublisher),
		subscribers: make(map[string]*RedisSubscriber),
	}
}

func (r *Redis) Start() error {
	log.Debug("starting redis")

	if r.config.On {
		options := &redis.Options{
			Addr:           r.config.Host + ":" + r.config.Port,
			Password:       r.config.Password,
			DB:             r.config.Database,
			Protocol:       r.config.Protocol, // specify 2 for RESP 2 or 3 for RESP 3
			MaxActiveConns: r.config.MaxActiveConns,
			PoolSize:       r.config.PoolSize,
			// TLSConfig: &tls.Config{},
		}

		// r.Client = redis.NewClient(options)
		r.subscriberClient = redis.NewClient(options)
		r.publisherClient = redis.NewClient(options)
	}

	return nil
}

func (r *Redis) Publish(channel types.Channel, payload []byte) error {
	if r.publishers[string(channel)] == nil {
		log.Debug("starting redis publisher: " + string(channel))

		r.publishers[string(channel)] = NewRedisPublisher(channel, r.publisherClient, &r.publishersWg)

		if err := r.publishers[string(channel)].Start(); err != nil {
			return err
		}
	}

	return r.publishers[string(channel)].Publish(payload)
}

func (r *Redis) Subscribe(channel types.Channel, callback func(payload []byte) (types.Response, error)) error {
	ctx := context.Background()
	sub := r.subscriberClient.Subscribe(ctx, string(channel))

	if r.subscribers[string(channel)] == nil {
		r.subscribers[string(channel)] = NewRedisSubscriber(channel, sub, &r.subscribersWg)
		if err := r.subscribers[string(channel)].Start(callback); err != nil {
			return err
		}

		return nil
	}

	return ErrRedisAlreadySubscribed(channel)
}

func (r *Redis) Stop() error {
	log.Debug("stopping redis")

	if r.config.On {
		log.Debug("  stopping redis publishers")
		// for _, publisher := range r.publishers {
		// 	publisher.Stop()
		// }

		log.Debug("  stopping redis subscribers")
		for _, subscriber := range r.subscribers {
			subscriber.Stop()
		}

		// r.publishersWg.Wait()
		r.subscribersWg.Wait()
		// c := r.Client.Close()
		// if c.Error() != "" {
		// 	return errors.New(c.Error())
		// }
	}

	return nil
}
