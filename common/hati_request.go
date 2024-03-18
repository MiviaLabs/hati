package common

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

// HatiRequest is being attached to action handler
// If call to method is coming from HTTP transport, it will contain RoutingContext which is Context from FastHttp.
// Otherwise, Message will be provided with payload and other details for messages sent over by internal transports.
type HatiRequest struct {
	RoutingContext   *routing.Context
	TransportManager *TransportManager
	Message          Message[[]byte]
}
