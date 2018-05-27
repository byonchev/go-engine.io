package transport

import (
	"errors"
	"net/http"
	"sync"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/gorilla/websocket"
)

// Websocket handles protocol upgrade and transmission over websockets
type Websocket struct {
	writeLock sync.Mutex
	readLock  sync.Mutex

	running bool

	socket *websocket.Conn

	codec codec.Codec
}

// NewWebsocket creates new Websocket transport
func NewWebsocket() *Websocket {
	transport := &Websocket{
		running: false,
		codec:   codec.Websocket{},
	}

	transport.lock()

	return transport
}

// HandleRequest handles initial websocket upgrade request
func (transport *Websocket) HandleRequest(writer http.ResponseWriter, request *http.Request) error {
	defer transport.unlock()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024, // TODO: Configuration!
		WriteBufferSize: 1024,
	}

	socket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		return err
	}

	transport.socket = socket
	transport.running = true

	return nil
}

// Shutdown closes the client socket
func (transport *Websocket) Shutdown() {
	transport.lock()
	defer transport.unlock()

	transport.running = false
	transport.socket.Close()
}

// Send writes packet to the client socket
func (transport *Websocket) Send(message packet.Packet) error {
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
func (transport *Websocket) Receive() (packet.Packet, error) {
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
func (transport *Websocket) Type() string {
	return WebsocketType
}

// Upgrades returns the possible transport upgrades
func (transport *Websocket) Upgrades() []string {
	return []string{}
}

func (transport *Websocket) lock() {
	transport.readLock.Lock()
	transport.writeLock.Lock()
}

func (transport *Websocket) unlock() {
	transport.readLock.Unlock()
	transport.writeLock.Unlock()
}
