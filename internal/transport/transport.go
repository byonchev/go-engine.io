package transport

import (
	"net/http"

	"github.com/byonchev/go-engine.io/internal/config"
	"github.com/byonchev/go-engine.io/internal/packet"
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

	HandleRequest(http.ResponseWriter, *http.Request)

	Send(packet.Packet) error
	Receive() (packet.Packet, error)

	Shutdown()
	Running() bool
}

// NewTransport creates a transport of the selected type
func NewTransport(name string, config config.Config) Transport {
	originCheck := config.CheckOrigin

	switch name {
	case WebsocketType:
		readBufferSize := config.WebsocketReadBufferSize
		writeBufferSize := config.WebsocketWriteBufferSize
		enableCompression := config.PerMessageDeflate

		return NewWebsocket(readBufferSize, writeBufferSize, enableCompression, originCheck)
	case PollingType:
		flushLimit := config.PollingBufferFlushLimit
		receiveLimit := config.PollingBufferReceiveLimit

		return NewPolling(flushLimit, receiveLimit, originCheck)
	default:
		return nil
	}
}
