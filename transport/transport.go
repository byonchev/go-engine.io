package transport

import (
	"net/http"

	"github.com/byonchev/go-engine.io/packet"
)

// Upgrades holds the possible upgrades for each transport
var Upgrades = map[Type][]Type{
	PollingType:   []Type{WebSocketType},
	WebSocketType: []Type{},
}

// Transport handles the delivery of packets between the client and the server
type Transport interface {
	Type() Type

	HandleRequest(http.ResponseWriter, *http.Request)

	Send(packet.Packet) error
	Receive() (packet.Packet, error)

	Shutdown()
}

// NewTransport creates a transport of the selected type
func NewTransport(transportType Type) Transport {
	switch transportType {
	case WebSocketType:
		return NewWebSocket()
	case PollingType:
		return NewXHR()
	default:
		return nil
	}
}
