# Hati

**Build distributed applications with Hati** or don't, choice is yours.

## What is that?

**Hati** gives you - out of the box - layer of communication for your distributed application or api, it has couple of integrations with different protocols which can be used for building your next project.

These integrations are:

- [x] Redis
- [ ] RabbitMQ
- [ ] P2P Protocol

## How does that work?

I'll tell you later :-)

## Examples

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MiviaLabs/hati"
	"github.com/MiviaLabs/hati/core"
	"github.com/MiviaLabs/hati/log"
	"github.com/MiviaLabs/hati/module"
	"github.com/MiviaLabs/hati/transport"
)

func main() {
	hati := hati.New(core.Config{
		Name: "example-app",
		Transport: transport.TransportManagerConfig{
			Redis: transport.RedisConfig{
				On:       true,
				Host:     "localhost",
				Port:     "6379",
				Username: "",
				Password: "",
				Database: 0,
			},
		},
	})

	m := module.New("test")

	hati.AddModule(m)

	if err := hati.Start(); err != nil {
		panic(err)
	}

	var osSignal chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-osSignal:
			log.Warning("shutting down, please wait")

			if err := hati.Stop(); err != nil {
				panic(err)
			}

			os.Exit(0)
		}
	}
}

```

---

### Useful commands

```
go build -gcflags="-m" main.go
```
