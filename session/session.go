package session

import (
	"errors"
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

	lastPingTime time.Time

	listener MessageListener
}

// NewSession creates a new client session
func NewSession(id string, config Config) *Session {
	return &Session{
		id:     id,
		config: config,

		transport: transport.NewXHR(),
		state:     new,
	}
}

// HandleRequest is the bridge between the engine.io endpoint and the selected transport
func (session *Session) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if session.state == new {
		session.handshake()
	}

	defer request.Body.Close()

	session.transport.HandleRequest(writer, request)
}

// Send enqueues packets for sending (non-blocking)
func (session *Session) Send(packet packet.Packet) error {
	if session.state == closed {
		return errors.New("send on closed session")
	}

	session.sending.Add(1)
	err := session.transport.Send(packet)
	session.sending.Done()

	return err
}

// Close changes the session state and closes the channels
func (session *Session) Close() {
	if session.state == closed {
		return
	}

	session.state = closed

	session.sending.Wait()
	session.transport.Shutdown()
}

// AttachListener sets listener for received packets
func (session *Session) AttachListener(listener MessageListener) {
	session.config.Listener = listener
}

// ID returns the session ID
func (session *Session) ID() string {
	return session.id
}

// Expired check if session is closed or last ping was not before (ping interval + ping timeout)
func (session *Session) Expired() bool {
	now := time.Now()
	threshold := session.config.PingInterval + session.config.PingTimeout

	return session.state == closed || session.lastPingTime.Add(threshold).Before(now)
}

func (session *Session) handshake() {
	packet := createHandshakePacket(session.id, session.config)

	session.Send(packet)

	session.state = active
	session.ping()

	go session.receivePackets()
}

func (session *Session) ping() {
	session.lastPingTime = time.Now()
}

func (session *Session) receivePackets() {
	for session.state != closed {
		received, err := session.transport.Receive()

		if err != nil {
			session.Close()
			break
		}

		session.ping()

		switch received.Type {
		case packet.Ping:
			session.Send(packet.NewPong(nil))
		case packet.Close:
			session.Close()
		case packet.Message:
			listener := session.config.Listener

			if listener != nil {
				listener.OnMessage(session, received)
			}
		}
	}
}
