package engine

import (
	"net/http"
	"net/url"
	"sync"

	"../session"
)

// Server defines engine.io http endpoint and holds connected clients
type Server struct {
	sync.RWMutex

	config  Config
	clients map[string]*session.Session
}

// NewServer creates a new engine server
func NewServer(config Config) *Server {
	return &Server{
		config:  config,
		clients: make(map[string]*session.Session),
	}
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	sessionID := request.URL.Query().Get("sid")

	var client *session.Session

	if sessionID == "" {
		client = server.createSession(request.URL.Query())
	} else {
		client = server.clients[sessionID]
	}

	if client == nil {
		// TODO: no session error
		return
	}

	client.HandleRequest(writer, request)
}

func (server *Server) createSession(params url.Values) *session.Session {
	sid := server.config.SIDGenerator.Generate()

	config := session.Config{
		server.config.PingSettings,
	}

	session := session.New(sid, config)

	server.Lock()
	defer server.Unlock()

	server.clients[sid] = session

	return session
}
