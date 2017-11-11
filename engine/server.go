package engine

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/byonchev/go-engine.io/session"
)

// Server defines engine.io http endpoint and holds connected clients
type Server struct {
	sync.RWMutex

	config  Config
	clients map[string]*session.Session
}

// NewServer creates a new engine server
func NewServer(config Config) *Server {
	server := &Server{
		config:  config,
		clients: make(map[string]*session.Session),
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
		// TODO: no session error
		return
	}

	client.HandleRequest(writer, request)
}

func (server *Server) checkPing() {
	interval := server.config.PingInterval + server.config.PingTimeout

	for {
		time.Sleep(interval)

		server.Lock()

		for id, session := range server.clients {
			if session.Expired() {
				session.Close()

				delete(server.clients, id)
			}
		}

		server.Unlock()
	}
}

func (server *Server) createSession(params url.Values) *session.Session {
	sid := server.config.SIDGenerator.Generate()

	config := session.Config{
		PingSettings: server.config.PingSettings,
	}

	session := session.NewSession(sid, config)

	server.Lock()
	defer server.Unlock()

	server.clients[sid] = session

	return session
}

func (server *Server) findSession(id string) *session.Session {
	server.RLock()
	defer server.RUnlock()

	return server.clients[id]
}