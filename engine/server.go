package engine

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/logger"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/session"
)

// Server defines engine.io http endpoint and holds connected clients
type Server struct {
	sync.RWMutex

	config  session.Config
	clients map[string]*session.Session

	events chan interface{}
}

// NewServer creates a new engine server
func NewServer(config session.Config) *Server {
	server := &Server{
		config:  config,
		clients: make(map[string]*session.Session),
		events:  make(chan interface{}),
	}

	go server.checkPing()

	return server
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	sessionID := request.URL.Query().Get("sid")

	var client *session.Session

	if sessionID == "" {
		client = server.createSession(request.URL.Query())
	} else {
		client = server.findSession(sessionID)
	}

	if client == nil {
		logger.Error("Session", sessionID, "not found")
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

func (server *Server) checkPing() {
	interval := server.config.PingInterval + server.config.PingTimeout

	for {
		time.Sleep(interval)

		server.Lock()

		for id, session := range server.clients {
			if session.Expired() {
				session.Close("ping timeout")

				delete(server.clients, id)
			}
		}

		server.Unlock()
	}
}

func (server *Server) createSession(params url.Values) *session.Session {
	session := session.NewSession(server.config, server.events)

	server.Lock()
	defer server.Unlock()

	server.clients[session.ID()] = session

	return session
}

func (server *Server) findSession(id string) *session.Session {
	server.RLock()
	defer server.RUnlock()

	return server.clients[id]
}
