package types

import (
	"github.com/MiviaLabs/hati/common/structs"
)

type Channel string
type DeliveryMethod string
type Response any
type TransportType string

const (
	CHAN_PREFIX                   = "hati_"
	CHAN_DISCOVERY        Channel = CHAN_PREFIX + "discovery"
	CHAN_MESSAGE          Channel = CHAN_PREFIX + "msg"
	CHAN_MESSAGE_RESPONSE Channel = CHAN_PREFIX + "msg_res"
	CHAN_STORAGE          Channel = CHAN_PREFIX + "storage"
)

const (
	DELIVERY_RANDOM DeliveryMethod = "random"
	DELIVERY_ALL    DeliveryMethod = "all"
)

type ActionHandler func(payload structs.Message[[]byte]) (Response, error)

const (
	GET HttpMethod = iota
	POST
	PATCH
	PUT
	DELETE
	OPTIONS
)

type HttpMethod uint8

func (m HttpMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PATCH:
		return "PATCH"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case OPTIONS:
		return "OPTIONS"
	default:
		return "Invalid http method"
	}
}

type Action struct {
	Handler ActionHandler
	Name    string
	Route   *structs.ActionRoute
}
