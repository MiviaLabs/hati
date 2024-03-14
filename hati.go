// Package hati is a framework for building distributed applications
package hati

import (
	"github.com/MiviaLabs/hati/core"
)

func New(config core.Config) core.Hati {
	return core.NewHati(config)
}
