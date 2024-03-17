package common

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

type HatiContext struct {
	RoutingContext   *routing.Context
	TransportManager *TransportManager
}
