package transport

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/MiviaLabs/hati/common/interfaces"
	"github.com/MiviaLabs/hati/common/structs"
	"github.com/MiviaLabs/hati/common/types"
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
	modules          *map[string]interfaces.Module
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

func (s *HttpServer) AddRoute() error {
	return nil
}

func (s *HttpServer) SetModules(modules *map[string]interfaces.Module) {
	s.modules = modules
}

func (s *HttpServer) Start() error {
	if !s.config.On {
		return nil
	}

	s.router.Use(func(c *routing.Context) error {
		c.Response.Header.Set("Access-Control-Allow-Headers", "*")
		c.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		c.Response.Header.Set("Access-Control-Allow-Method", "*")

		return c.Next()
	})

	for _, module := range *s.modules {
		moduleActions := module.GetActions()

		for _, action := range *moduleActions {
			if action != nil && action.Route != nil && len(action.Route.Methods) > 0 {
				for _, httpMethod := range action.Route.Methods {

					switch httpMethod {
					case types.GET.String():
						{
							s.router.Get(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(structs.Message[[]byte]{
									RoutingContext: c,
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
					case types.POST.String():
						{
							s.router.Post(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(structs.Message[[]byte]{
									RoutingContext: c,
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
					case types.PATCH.String():
						{
							s.router.Patch(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(structs.Message[[]byte]{
									RoutingContext: c,
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
					case types.PUT.String():
						{
							s.router.Put(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(structs.Message[[]byte]{
									RoutingContext: c,
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
					case types.DELETE.String():
						{
							s.router.Delete(action.Route.Path, func(c *routing.Context) error {
								res, err := action.Handler(structs.Message[[]byte]{
									RoutingContext: c,
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
