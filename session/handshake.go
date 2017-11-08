package session

import (
	"encoding/json"
	"time"

	"../packet"
)

type handshakeMessage struct {
	SessionID    string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int64    `json:"pingTimeout"`
	PingInterval int64    `json:"pingInterval"`
}

func createHandshakePacket(sid string, config Config) packet.Packet {
	handshake := handshakeMessage{
		SessionID:    sid,
		PingInterval: int64(config.PingInterval / time.Millisecond),
		PingTimeout:  int64(config.PingTimeout / time.Millisecond),
		Upgrades:     []string{},
	}

	json, _ := json.Marshal(handshake)

	return packet.NewOpen(json)
}
