package session

import (
	"github.com/byonchev/go-engine.io/packet"
)

// Listener is called on session events
type Listener interface {
	OnOpen(*Session)
	OnClose(*Session)
	OnMessage(*Session, packet.Packet)
}
