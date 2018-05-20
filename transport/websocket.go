package transport

import (
	"errors"
	"net/http"
	"sync"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/logger"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/gorilla/websocket"
)

// WebSocket handles protocol upgrade and transmission over websockets
type WebSocket struct {
	writeLock sync.Mutex
	readLock  sync.Mutex

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

	transport.lock()

	return transport
}

// HandleRequest handles initial websocket upgrade request
func (transport *WebSocket) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	defer transport.unlock()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024, // TODO: Configuration!
		WriteBufferSize: 1024,
	}

	socket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		logger.Error("WebSocket upgrade failed:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	transport.socket = socket
	transport.running = true
}

// Shutdown closes the client socket
func (transport *WebSocket) Shutdown() {
	transport.lock()
	defer transport.unlock()

	transport.running = false
	transport.socket.Close()
}

// Send writes packet to the client socket
func (transport *WebSocket) Send(message packet.Packet) error {
	transport.writeLock.Lock()
	defer transport.writeLock.Unlock()

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
	transport.readLock.Lock()
	defer transport.readLock.Unlock()

	if !transport.running {
		return packet.Packet{}, errors.New("transport not running")
	}

	_, reader, err := transport.socket.NextReader()

	if err != nil {
		transport.running = false

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

// Type returns the transport identifier
func (transport *WebSocket) Type() string {
	return WebsocketType
}

// Upgrades returns the possible transport upgrades
func (transport *WebSocket) Upgrades() []string {
	return []string{}
}

func (transport *WebSocket) lock() {
	transport.readLock.Lock()
	transport.writeLock.Lock()
}

func (transport *WebSocket) unlock() {
	transport.readLock.Unlock()
	transport.writeLock.Unlock()
}
