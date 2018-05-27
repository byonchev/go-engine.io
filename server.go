package eio

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/internal/config"
	"github.com/byonchev/go-engine.io/internal/logger"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/byonchev/go-engine.io/internal/transport"
)

// Server defines engine.io http endpoint and holds connected clients
type Server struct {
	config.Config

	sync.RWMutex

	clients map[string]*Session

	events chan interface{}
}

// NewServer creates a new engine server
func NewServer() *Server {
	server := &Server{
		clients: make(map[string]*Session),
		events:  make(chan interface{}),

		Config: config.Config{
			PingInterval:              25 * time.Second,
			PingTimeout:               60 * time.Second,
			Transports:                []string{transport.PollingType, transport.WebsocketType},
			AllowUpgrades:             true,
			UpgradeTimeout:            10 * time.Second,
			PollingBufferFlushLimit:   10,
			PollingBufferReceiveLimit: 10,
			WebsocketReadBufferSize:   1024,
			WebsocketWriteBufferSize:  1024,
			PerMessageDeflate:         true,
			CheckOrigin:               func(*http.Request) bool { return true },
		},
	}

	go server.checkPing()

	return server
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	sessionID := request.URL.Query().Get("sid")

	var client *Session

	if sessionID == "" {
		client = server.createSession(request.URL.Query())
	} else {
		client = server.findSession(sessionID)
	}

	if client == nil {
		logger.Error("Session ", sessionID, " not found")
		return
	}

	client.HandleRequest(writer, request)
}

// Events returns the channel for session events
func (server *Server) Events() <-chan interface{} {
	return server.events
}

// Send sends message to a specific session
func (server *Server) Send(id string, binary bool, data []byte) error {
	server.RLock()
	session := server.findSession(id)
	server.RUnlock()

	if session == nil {
		return errors.New("invalid session")
	}

	return session.Send(packet.NewMessage(binary, data))
}

// SetLogger initializes logging with a specific implementation
func (server *Server) SetLogger(loggerInstance logger.Logger) {
	logger.Init(loggerInstance)
}

func (server *Server) checkPing() {
	interval := server.PingInterval + server.PingTimeout

	for {
		time.Sleep(interval)

		server.Lock()

		for id, session := range server.clients {
			if session.Expired() {
				go session.Close("ping timeout")

				delete(server.clients, id)
			}
		}

		server.Unlock()
	}
}

func (server *Server) createSession(params url.Values) *Session {
	session := NewSession(server.Config, server.events)

	server.Lock()
	defer server.Unlock()

	server.clients[session.ID()] = session

	return session
}

func (server *Server) findSession(id string) *Session {
	server.RLock()
	defer server.RUnlock()

	return server.clients[id]
}
