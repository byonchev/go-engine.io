package session

import (
	"github.com/byonchev/go-engine.io/packet"
)

// MessageListener receives events for session messages
type MessageListener interface {
	OnMessage(*Session, packet.Packet)
}
