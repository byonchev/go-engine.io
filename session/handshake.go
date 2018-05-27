package session

import (
	"encoding/json"
	"time"

	"github.com/byonchev/go-engine.io/config"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/byonchev/go-engine.io/transport"
)

type handshakeMessage struct {
	SessionID    string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int64    `json:"pingTimeout"`
	PingInterval int64    `json:"pingInterval"`
}

func createHandshakePacket(sid string, transport transport.Transport, config config.Config) packet.Packet {
	handshake := handshakeMessage{
		SessionID:    sid,
		PingInterval: int64(config.PingInterval / time.Millisecond),
		PingTimeout:  int64(config.PingTimeout / time.Millisecond),
		Upgrades:     transport.Upgrades(),
	}

	json, _ := json.Marshal(handshake)

	return packet.NewOpen(json)
}
