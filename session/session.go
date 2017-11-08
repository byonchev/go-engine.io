package session

import (
	"net/http"

	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
)

// Session holds information for a single connected client
type Session struct {
	id        string
	config    Config
	transport transport.Transport

	established bool

	sendChannel    chan packet.Packet
	receiveChannel chan packet.Packet
}

// New creates a new client session
func New(id string, config Config) *Session {
	sendChannel := make(chan packet.Packet)
	receiveChannel := make(chan packet.Packet)

	return &Session{
		id:     id,
		config: config,

		transport:   transport.NewXHR(sendChannel, receiveChannel),
		established: false,

		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
	}
}

// HandleRequest is the bridge between the engine.io endpoint and the selected transport
func (session *Session) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if !session.established {
		session.handshake()
	}

	session.transport.HandleRequest(writer, request)
}

// Send enqueues packets for sending (non-blocking)
func (session *Session) Send(packet packet.Packet) {
	go func() { session.sendChannel <- packet }()
}

func (session *Session) handshake() {
	packet := createHandshakePacket(session.id, session.config)

	session.Send(packet)

	session.established = true

	go session.receiveLoop()
}

func (session *Session) receiveLoop() {
	for {
		received := <-session.receiveChannel

		if received.Type == packet.Ping {
			session.Send(packet.NewPong())
		}
	}
}
