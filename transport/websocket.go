package transport

import (
	"errors"
	"net/http"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/logger"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	running bool

	socket *websocket.Conn

	codec codec.Codec
}

// NewWebSocket creates new WebSocket transport
func NewWebSocket() *WebSocket {
	transport := &WebSocket{
		running: false,
		codec:   codec.WebSocket{},
	}

	return transport
}

// HandleRequest handles WebSocket upgrade request
func (transport *WebSocket) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024, // TODO: Configuration!
		WriteBufferSize: 1024,
	}

	socket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		logger.Error("Protocol upgrade error:", err)
		return
	}

	transport.socket = socket
	transport.running = true
}

// Shutdown closes the client socket
func (transport *WebSocket) Shutdown() {
	transport.running = false
	transport.socket.Close()
}

// Send writes packet to the client socket
func (transport *WebSocket) Send(message packet.Packet) error {
	if !transport.running {
		return errors.New("transport not running")
	}

	var messageType int

	if message.Binary {
		messageType = websocket.BinaryMessage
	} else {
		messageType = websocket.TextMessage
	}

	writer, err := transport.socket.NextWriter(messageType)

	if err != nil {
		return err
	}

	payload := packet.Payload{message}

	err = transport.codec.Encode(payload, writer)

	if err != nil {
		return err
	}

	return writer.Close()
}

// Receive receives the next packet from the client socket
func (transport *WebSocket) Receive() (packet.Packet, error) {
	if !transport.running {
		return packet.Packet{}, errors.New("transport not running")
	}

	_, reader, err := transport.socket.NextReader()

	if err != nil {
		transport.Shutdown()
		return packet.Packet{}, err
	}

	payload, err := transport.codec.Decode(reader)

	if err != nil {
		return packet.Packet{}, err
	}

	count := len(payload)

	if count == 0 {
		return packet.Packet{}, errors.New("empty payload received")
	} else if count > 1 {
		return packet.Packet{}, errors.New("multiple packets received on single websocket frame")
	}

	return payload[0], nil
}

// Type returns the transport type
func (transport *WebSocket) Type() Type {
	return WebSocketType
}
