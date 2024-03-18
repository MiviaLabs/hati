package transport

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/MiviaLabs/hati/common"
	"github.com/MiviaLabs/hati/log"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type HttpConfig struct {
	On   bool      `yaml:"on" json:"on"`
	Host string    `yaml:"host" json:"host"`
	Port string    `yaml:"port" json:"port"`
	Cors *HttpCors `yaml:"cors" json:"cors"`
}

type HttpCors struct {
	AllowHeaders     string `yaml:"allow_headers" json:"allow_headers"`
	AllowCredentials string `yaml:"allow_credentials" json:"allow_credentials"`
	AllowOrigin      string `yaml:"allow_origin" json:"allow_origin"`
	AllowMethod      string `yaml:"allow_method" json:"allow_method"`
}

type HttpServer struct {
	common.HttpServer
	config           HttpConfig
	transportManager common.TransportManager
	modules          *map[string]common.Module
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
func (s *HttpServer) SetTransportManager(transportManager common.TransportManager) {
	s.transportManager = transportManager
}

func (s *HttpServer) AddRoute() error {
	return nil
}

func (s *HttpServer) SetModules(modules *map[string]common.Module) {
	s.modules = modules
}

func (s *HttpServer) Start() error {
	if !s.config.On {
		return nil
	}

	s.router.Use(func(c *routing.Context) error {
		c.Response.Header.Set("Content-Type", "application/json")

		return c.Next()
	})

	if s.config.Cors != nil {
		s.router.Use(func(c *routing.Context) error {

			if s.config.Cors.AllowCredentials != "" {
				c.Response.Header.Set("Access-Control-Allow-Credentials", s.config.Cors.AllowCredentials)
			}

			if s.config.Cors.AllowHeaders != "" {
				c.Response.Header.Set("Access-Control-Allow-Headers", s.config.Cors.AllowHeaders)
			}

			if s.config.Cors.AllowOrigin != "" {
				c.Response.Header.Set("Access-Control-Allow-Origin", s.config.Cors.AllowOrigin)
			}

			if s.config.Cors.AllowMethod != "" {
				c.Response.Header.Set("Access-Control-Allow-Method", s.config.Cors.AllowMethod)
			}

			return c.Next()
		})
	}

	for _, module := range *s.modules {
		moduleActions := module.GetActions()

		for _, action := range *moduleActions {
			if action != nil && action.Route != nil && len(action.Route.Methods) > 0 {
				for _, httpMethod := range action.Route.Methods {

					switch httpMethod {
					case common.GET.String():
						{
							s.router.Get(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(&common.HatiRequest{
									RoutingContext:   c,
									TransportManager: &s.transportManager,
								})

								if err != nil {
									return err
								}

								resBytes, err := json.Marshal(res)
								if err != nil {
									return err
								}

								c.Response.SetBody(resBytes)

								return nil
							})
							break
						}
					case common.POST.String():
						{
							s.router.Post(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(&common.HatiRequest{
									RoutingContext:   c,
									TransportManager: &s.transportManager,
								})

								if err != nil {
									return err
								}

								resBytes, err := json.Marshal(res)
								if err != nil {
									return err
								}

								c.Response.SetBody(resBytes)

								return nil
							})
							break
						}
					case common.PATCH.String():
						{
							s.router.Patch(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(&common.HatiRequest{
									RoutingContext:   c,
									TransportManager: &s.transportManager,
								})

								if err != nil {
									return err
								}

								resBytes, err := json.Marshal(res)
								if err != nil {
									return err
								}

								c.Response.SetBody(resBytes)

								return nil
							})
							break
						}
					case common.PUT.String():
						{
							s.router.Put(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(&common.HatiRequest{
									RoutingContext:   c,
									TransportManager: &s.transportManager,
								})

								if err != nil {
									return err
								}

								resBytes, err := json.Marshal(res)
								if err != nil {
									return err
								}

								c.Response.SetBody(resBytes)

								return nil
							})
							break
						}
					case common.DELETE.String():
						{
							s.router.Delete(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(&common.HatiRequest{
									RoutingContext:   c,
									TransportManager: &s.transportManager,
								})

								if err != nil {
									return err
								}

								resBytes, err := json.Marshal(res)
								if err != nil {
									return err
								}

								c.Response.SetBody(resBytes)

								return nil
							})
							break
						}
					default:
						return errors.New("invalid http method")
					}
				}
			}
		}
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
