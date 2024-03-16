package types

import "github.com/MiviaLabs/hati/common/structs"

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
