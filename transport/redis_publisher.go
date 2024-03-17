package transport

import (
	"context"
	"sync"

	"github.com/MiviaLabs/hati/common"
	"github.com/MiviaLabs/hati/log"
	redis "github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	channel      common.Channel
	pub          *redis.Client
	pubCloseChan chan bool
	inChan       chan string
	wg           *sync.WaitGroup
}

func NewRedisPublisher(channel common.Channel, pub *redis.Client, wg *sync.WaitGroup) *RedisPublisher {
	return &RedisPublisher{
		channel:      channel,
		pub:          pub,
		pubCloseChan: make(chan bool),
		inChan:       make(chan string, 100),
		wg:           wg,
	}
}

func (rp *RedisPublisher) Start() error {
	log.Debug("starting redis publisher: " + string(rp.channel))

	rp.wg.Add(1)

	go func(p *RedisPublisher) {
		defer p.pub.Close()
		defer close(p.inChan)
		defer close(p.pubCloseChan)
		defer p.wg.Done()

	Loop:
		for {
			select {
			case payload := <-p.inChan:
				ctx := context.Background()

				p.pub.Publish(ctx, string(p.channel), payload)
			case <-p.pubCloseChan:
				log.Debug("    stopping redis publisher: " + string(p.channel))

				break Loop
			}
		}
	}(rp)

	return nil
}

func (rp *RedisPublisher) Stop() error {
	rp.pubCloseChan <- true

	return nil
}

func (rp *RedisPublisher) Publish(payload []byte) error {
	rp.inChan <- string(payload)

	return nil
}
