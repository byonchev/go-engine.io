package engine

import (
	"github.com/byonchev/go-engine.io/session"
)

// Listener receives events for server sessions
type Listener interface {
	session.MessageListener

	OnOpen(*session.Session)
	OnClose(*session.Session)
}
