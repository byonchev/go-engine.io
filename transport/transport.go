package transport

import (
	"net/http"

	"github.com/byonchev/go-engine.io/packet"
)

// String identifiers for each supported transport
const (
	WebsocketType = "websocket"
	PollingType   = "polling"
)

// Transport handles the delivery of packets between the client and the server
type Transport interface {
	Type() string
	Upgrades() []string

	HandleRequest(http.ResponseWriter, *http.Request) error

	Send(packet.Packet) error
	Receive() (packet.Packet, error)

	Shutdown()
}

// NewTransport creates a transport of the selected type
func NewTransport(name string) Transport {
	switch name {
	case WebsocketType:
		return NewWebsocket()
	case PollingType:
		return NewPolling(10, 10) // TODO: Configuration
	default:
		return nil
	}
}
