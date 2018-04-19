package transport

import (
	"net/http"

	"github.com/byonchev/go-engine.io/packet"
)

// Transport handles the delivery of packets between the client and the server
type Transport interface {
	HandleRequest(http.ResponseWriter, *http.Request)
	Send(packet.Packet) error
	Receive() (packet.Packet, error)
	Shutdown()
}

// NewTransport creates a transport of the selected type
func NewTransport(transportType Type, sendChannel <-chan packet.Packet, receiveChannel chan<- packet.Packet) Transport {
	switch transportType {
	// case WebSocketType:
	// return NewWebSocket(sendChannel, receiveChannel)
	case PollingType:
		return NewXHR()
	default:
		return nil
	}
}
