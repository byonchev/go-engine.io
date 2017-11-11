package session

import (
	"net/http"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
)

// Session holds information for a single connected client
type Session struct {
	id        string
	config    Config
	transport transport.Transport

	state state

	sending sync.WaitGroup

	sendChannel    chan packet.Packet
	receiveChannel chan packet.Packet

	lastPingTime time.Time
}

// NewSession creates a new client session
func NewSession(id string, config Config) *Session {
	sendChannel := make(chan packet.Packet, 10)
	receiveChannel := make(chan packet.Packet)

	return &Session{
		id:     id,
		config: config,

		transport: transport.NewXHR(sendChannel, receiveChannel),
		state:     new,

		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
	}
}

// HandleRequest is the bridge between the engine.io endpoint and the selected transport
func (session *Session) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if session.state == new {
		session.handshake()
	}

	session.transport.HandleRequest(writer, request)
}

// Send enqueues packets for sending (non-blocking)
func (session *Session) Send(packet packet.Packet) {
	if session.state == closed {
		return
	}

	session.sending.Add(1)
	session.sendChannel <- packet
	session.sending.Done()
}

// Close changes the session state and closes the channels
func (session *Session) Close() {
	session.state = closed

	session.transport.Shutdown()
	close(session.receiveChannel)

	session.sending.Wait()
	close(session.sendChannel)
}

func (session *Session) handshake() {
	packet := createHandshakePacket(session.id, session.config)

	session.Send(packet)

	session.state = active
	session.ping()

	go session.receivePackets()
}

// Expired check if session is closed or last ping was not before (ping interval + ping timeout)
func (session *Session) Expired() bool {
	now := time.Now()
	threshold := session.config.PingInterval + session.config.PingTimeout

	return session.state == closed || session.lastPingTime.Add(threshold).Before(now)
}

func (session *Session) ping() {
	session.lastPingTime = time.Now()
}

func (session *Session) receivePackets() {
	for {
		received, ok := <-session.receiveChannel

		if !ok {
			return
		}

		session.ping()

		switch received.Type {
		case packet.Ping:
			session.Send(packet.NewPong(nil))

		case packet.Close:
			session.Close()
		}
	}
}
