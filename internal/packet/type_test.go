package packet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/byonchev/go-engine.io/internal/packet"
)

func TestPacketTypeToChar(t *testing.T) {
	packet := packet.Packet{Type: 9}

	assert.Equal(t, byte('9'), packet.Type.Char(), "invalid packet type character representation")
}

func TestPacketTypeFromChar(t *testing.T) {
	packetType := packet.TypeFromChar('6')

	assert.Equal(t, packet.NOOP, packetType, "invalid packet type parsing from character")
}

func TestPacketTypeToByte(t *testing.T) {
	packet := packet.Packet{Type: 9}

	assert.Equal(t, byte(9), packet.Type.Byte(), "invalid packet type numeric representation")
}

func TestPacketTypeFromByte(t *testing.T) {
	packetType := packet.TypeFromByte(6)

	assert.Equal(t, packet.NOOP, packetType, "invalid packet type parsing from numeric")
}
