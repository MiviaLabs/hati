package transport

import (
	"sync"

	"github.com/MiviaLabs/hati/common/types"
	"github.com/MiviaLabs/hati/log"
	redis "github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	channel              types.Channel
	sub                  *redis.PubSub
	subCloseChan         chan bool
	subCallbackCloseChan chan bool
	outChan              chan string
	wg                   *sync.WaitGroup
}

func NewRedisSubscriber(channel types.Channel, sub *redis.PubSub, wg *sync.WaitGroup) *RedisSubscriber {
	return &RedisSubscriber{
		channel:              channel,
		sub:                  sub,
		subCloseChan:         make(chan bool),
		subCallbackCloseChan: make(chan bool),
		outChan:              make(chan string, 100),
		wg:                   wg,
	}
}

func (rs *RedisSubscriber) Start(callback func(payload []byte) error) error {
	log.Debug("    starting redis subscriber: " + string(rs.channel))

	rs.wg.Add(1)
	go func(s *RedisSubscriber) {

		defer s.sub.Close()
		defer close(s.subCloseChan)
		defer close(s.outChan)
		defer s.wg.Done()
	Loop:
		for {
			select {
			case msg := <-s.sub.Channel():
				{
					if msg.Channel == string(s.channel) {
						s.outChan <- msg.Payload
					}
				}
			case <-s.subCloseChan:
				log.Debug("    stopping redis subscriber: " + string(s.channel))
				break Loop
			}
		}
	}(rs)

	rs.wg.Add(1)
	go func(s *RedisSubscriber) {
		defer s.wg.Done()

	Loop:
		for {
			select {
			case payload := <-s.outChan:
				callback([]byte(payload))
			case <-s.subCallbackCloseChan:
				break Loop
			}
		}
	}(rs)

	return nil
}

func (rs *RedisSubscriber) Stop() error {
	rs.subCloseChan <- true
	rs.subCallbackCloseChan <- true

	return nil
}
