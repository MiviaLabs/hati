package transport

import (
	"sync"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/log"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type HttpConfig struct {
	On   bool   `yaml:"on" json:"on"`
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
}

type HttpServer struct {
	interfaces.HttpServer
	config           HttpConfig
	transportManager interfaces.TransportManager
	router           *routing.Router
	stopChan         chan bool
	stopWg           sync.WaitGroup
}

func NewHttpServer(config HttpConfig) *HttpServer {
	return &HttpServer{
		config:   config,
		router:   routing.New(),
		stopChan: make(chan bool),
	}
}
func (s *HttpServer) SetTransportManager(transportManager interfaces.TransportManager) {
	s.transportManager = transportManager
}

func (s *HttpServer) Start() error {
	if !s.config.On {
		return nil
	}

	log.Debug("http server is running at: http://" + s.config.Host + ":" + s.config.Port)

	s.stopWg.Add(1)

	go func(w *sync.WaitGroup, stop chan bool) {
		defer close(stop)
		defer w.Done()
	Loop:
		for {
			select {
			default:
				if err := fasthttp.ListenAndServe(s.config.Host+":"+s.config.Port, s.router.HandleRequest); err != nil {
					log.Error(err.Error())
				}
				break Loop
			case <-stop:
				w.Done()
				break Loop
			}
		}

	}(&s.stopWg, s.stopChan)

	return nil
}

func (s *HttpServer) Stop() error {
	log.Debug("stopping http server")

	// s.stopChan <- true

	// s.stopWg.Wait()

	return nil
}
