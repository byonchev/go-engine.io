package utils

import (
	"encoding/json"
	"time"

	"github.com/byonchev/go-engine.io/internal/config"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/byonchev/go-engine.io/internal/transport"
)

type handshakeMessage struct {
	SessionID    string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int64    `json:"pingTimeout"`
	PingInterval int64    `json:"pingInterval"`
}

// CreateHandshakePacket creates open packet with JSON serialized handshake messag
func CreateHandshakePacket(sid string, transport transport.Transport, config config.Config) packet.Packet {
	handshake := handshakeMessage{
		SessionID:    sid,
		PingInterval: int64(config.PingInterval / time.Millisecond),
		PingTimeout:  int64(config.PingTimeout / time.Millisecond),
		Upgrades:     getSupportedUpgrades(transport, config),
	}

	json, _ := json.Marshal(handshake)

	return packet.NewOpen(json)
}

func getSupportedUpgrades(transport transport.Transport, config config.Config) []string {
	result := []string{}

	if !config.AllowUpgrades {
		return result
	}

	for _, upgrade := range transport.Upgrades() {
		if StringSliceContains(config.Transports, upgrade) {
			result = append(result, upgrade)
		}
	}

	return result
}
