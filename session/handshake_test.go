package session

import (
	"testing"
	"time"

	"github.com/byonchev/go-engine.io/config"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestHandshakePacket(t *testing.T) {
	config := Config{
		PingSettings: config.PingSettings{
			PingInterval: 1 * time.Second,
			PingTimeout:  2 * time.Second,
		},
	}

	expected := packet.Packet{
		Binary: false,
		Type:   packet.Open,
		Data:   []byte("{\"sid\":\"100200300\",\"upgrades\":[],\"pingTimeout\":2000,\"pingInterval\":1000}"),
	}

	actual := createHandshakePacket("100200300", config)

	assert.Equal(t, expected, actual, "handshake packet is invalid")
}