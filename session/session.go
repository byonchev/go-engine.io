package session

import (
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/logger"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
	uuid "github.com/satori/go.uuid"
)

// Session holds information for a single connected client
type Session struct {
	id        string
	config    Config
	transport transport.Transport

	events chan<- interface{}

	handshaked bool
	closed     bool

	sending sync.WaitGroup

	lastPingTime time.Time
}

// NewSession creates a new client session
func NewSession(config Config, events chan<- interface{}) *Session {
	uuid, _ := uuid.NewV4()
	id := base64.URLEncoding.EncodeToString(uuid.Bytes())

	return &Session{
		id:     id,
		config: config,

		events: events,

		handshaked: false,
		closed:     false,
	}
}

// HandleRequest is the bridge between the engine.io endpoint and the selected transport
// TODO: Refactor
func (session *Session) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	query := request.URL.Query()
	requestedTransport := query.Get("transport")

	if session.transport == nil {
		session.transport = transport.NewTransport(requestedTransport)
	}

	if !session.handshaked {
		if session.transport.Type() == transport.PollingType {
			session.handshake()
		} else {
			defer session.handshake()
		}
	}

	if requestedTransport != session.transport.Type() {
		newTransport := transport.NewTransport(requestedTransport)
		session.upgrade(writer, request, newTransport)
	} else {
		session.transport.HandleRequest(writer, request)
	}
}

// Send enqueues packets for sending
func (session *Session) Send(packet packet.Packet) error {
	if session.closed {
		return errors.New("session closed")
	}

	session.sending.Add(1)
	err := session.transport.Send(packet)
	session.sending.Done()

	return err
}

// Close changes the session state and closes the channels
func (session *Session) Close(reason interface{}) {
	if session.closed {
		return
	}

	session.closed = true

	session.sending.Wait()
	session.transport.Shutdown()

	session.debug("Session closed. Reason:", reason)

	session.emit(DisconnectEvent{session.id})
}

// ID returns the session ID
func (session *Session) ID() string {
	return session.id
}

// Expired check if session is closed or last ping was not before (ping interval + ping timeout)
func (session *Session) Expired() bool {
	now := time.Now()
	threshold := session.config.PingInterval + session.config.PingTimeout

	return session.closed || session.lastPingTime.Add(threshold).Before(now)
}

func (session *Session) handshake() {
	packet := createHandshakePacket(session.id, session.transport, session.config)

	err := session.Send(packet)

	if err != nil {
		logger.Error("Handshake error:", err, "for", packet)
		return
	}

	session.debug("Session created")

	session.handshaked = true
	session.ping()

	go session.receivePackets()

	session.emit(ConnectEvent{session.id})
}

func (session *Session) ping() {
	session.lastPingTime = time.Now()
}

func (session *Session) receivePackets() {
	for !session.closed {
		received, err := session.transport.Receive()

		if err != nil {
			continue
		}

		session.ping()

		switch received.Type {
		case packet.Ping:
			session.handlePing(received)
		case packet.Close:
			session.handleClose(received)
		case packet.Message:
			session.handleMessage(received)
		}
	}
}

func (session *Session) handlePing(ping packet.Packet) {
	session.debug("Ping received")
	session.debug("Sending pong")

	session.Send(packet.NewPong(ping.Data))
}

func (session *Session) handleClose(close packet.Packet) {
	session.Close("close packet received")
}

func (session *Session) handleMessage(message packet.Packet) {
	session.debug("Message received:", message.Data)

	event := MessageEvent{
		SessionID: session.id,
		Binary:    message.Binary,
		Data:      message.Data,
	}

	session.emit(event)
}

func (session *Session) upgrade(writer http.ResponseWriter, request *http.Request, target transport.Transport) error {
	// TODO: Error
	target.HandleRequest(writer, request)

	session.debug("Upgrading transport")

	for {
		received, err := target.Receive()

		if err != nil {
			return err
		}

		if received.Type == packet.Ping && string(received.Data) == "probe" {
			err := target.Send(packet.NewPong(received.Data))

			if err != nil {
				return err
			}

			session.debug("Sending pong probe")

			session.transport.Send(packet.NewNOOP())
		} else if received.Type == packet.Upgrade {
			session.debug("Upgrade packet recevied")
			session.transport.Shutdown()
			session.transport = target
			break
		}
	}

	return errors.New("upgrade failed")
}

func (session *Session) emit(event interface{}) {
	go func() {
		session.events <- event
	}()
}

func (session *Session) debug(data ...interface{}) {
	prefix := []interface{}{
		"[",
		session.id,
		"]",
	}

	logger.Debug(append(prefix, data...)...)
}
