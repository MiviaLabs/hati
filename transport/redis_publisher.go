package transport

import (
	"fmt"

	"github.com/MiviaLabs/hati/common/types"
	redis "github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	channel      types.Channel
	pub          *redis.PubSub
	pubCloseChan chan bool
	inChan       chan string
}

func NewRedisPublisher(channel types.Channel, pub *redis.PubSub) *RedisPublisher {
	return &RedisPublisher{
		channel:      channel,
		pub:          pub,
		pubCloseChan: make(chan bool),
		inChan:       make(chan string, 100),
	}
}

func (rp *RedisPublisher) Start() error {
	go func(p *RedisPublisher) {
		defer p.pub.Close()
		defer close(p.inChan)
		defer close(p.pubCloseChan)

	Loop:
		for {
			select {
			case payload := <-p.inChan:
				fmt.Println(payload)
				// res := p.pub.Channel()
				// // res
			case <-p.pubCloseChan:
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
