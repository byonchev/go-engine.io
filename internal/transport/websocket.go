package transport

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/byonchev/go-engine.io/internal/codec"
	"github.com/byonchev/go-engine.io/internal/logger"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/gorilla/websocket"
)

// Websocket handles protocol upgrade and transmission over websockets
type Websocket struct {
	readBufferSize    int
	writeBufferSize   int
	enableCompression bool
	originCheck       func(*http.Request) bool

	writeLock sync.Mutex
	readLock  sync.Mutex

	running bool

	socket *websocket.Conn

	codec codec.Codec
}

// NewWebsocket creates new Websocket transport
func NewWebsocket(readBufferSize int, writeBufferSize int, enableCompression bool, originCheck func(*http.Request) bool) *Websocket {
	transport := &Websocket{
		readBufferSize:    readBufferSize,
		writeBufferSize:   writeBufferSize,
		enableCompression: enableCompression,
		originCheck:       originCheck,

		running: false,
		codec:   codec.Websocket{},
	}

	transport.lock()

	return transport
}

// HandleRequest handles initial websocket upgrade request
func (transport *Websocket) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	defer transport.unlock()

	upgrader := websocket.Upgrader{
		ReadBufferSize:    transport.readBufferSize,
		WriteBufferSize:   transport.writeBufferSize,
		EnableCompression: transport.enableCompression,
		CheckOrigin:       transport.originCheck,
	}

	socket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		logger.Error("Websocket upgrade failed: ", err)
		return
	}

	transport.socket = socket
	transport.running = true
}

// Shutdown closes the client socket
func (transport *Websocket) Shutdown() {
	transport.lock()
	defer transport.unlock()

	transport.close()
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
		return packet.Packet{}, io.EOF
	}

	_, reader, err := transport.socket.NextReader()

	if err != nil {
		transport.close()

		return packet.Packet{}, io.EOF
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

// Running returns true if the transport is active
func (transport *Websocket) Running() bool {
	return transport.running
}

// Type returns the transport identifier
func (transport *Websocket) Type() string {
	return WebsocketType
}

// Upgrades returns the possible transport upgrades
func (transport *Websocket) Upgrades() []string {
	return []string{}
}

func (transport *Websocket) close() {
	transport.running = false
	transport.socket.Close()
}

func (transport *Websocket) lock() {
	transport.readLock.Lock()
	transport.writeLock.Lock()
}

func (transport *Websocket) unlock() {
	transport.readLock.Unlock()
	transport.writeLock.Unlock()
}
