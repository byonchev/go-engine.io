package utils_test

import (
	"testing"
	"time"

	"github.com/byonchev/go-engine.io/internal/config"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/byonchev/go-engine.io/internal/transport"
	"github.com/byonchev/go-engine.io/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandshakePacket(t *testing.T) {
	config := config.Config{
		PingInterval:  1 * time.Second,
		PingTimeout:   2 * time.Second,
		Transports:    []string{"polling", "websocket"},
		AllowUpgrades: true,
	}

	expected := packet.Packet{
		Binary: false,
		Type:   packet.Open,
		Data:   []byte("{\"sid\":\"100200300\",\"upgrades\":[\"websocket\"],\"pingTimeout\":2000,\"pingInterval\":1000}"),
	}

	actual := utils.CreateHandshakePacket("100200300", &transport.Polling{}, config)

	assert.Equal(t, expected, actual, "handshake packet is invalid")
}

func TestNotUpgradablePacket(t *testing.T) {
	config := config.Config{
		PingInterval:  1 * time.Second,
		PingTimeout:   2 * time.Second,
		Transports:    []string{"polling", "websocket"},
		AllowUpgrades: false,
	}

	expected := packet.Packet{
		Binary: false,
		Type:   packet.Open,
		Data:   []byte("{\"sid\":\"100200300\",\"upgrades\":[],\"pingTimeout\":2000,\"pingInterval\":1000}"),
	}

	actual := utils.CreateHandshakePacket("100200300", &transport.Polling{}, config)

	assert.Equal(t, expected, actual, "handshake packet is invalid")
}
