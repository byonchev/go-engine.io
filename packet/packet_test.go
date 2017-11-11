package packet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/byonchev/go-engine.io/packet"
)

func TestNewOpenPacket(t *testing.T) {
	expected := packet.Packet{
		Binary: false,
		Type:   packet.Open,
		Data:   []byte{1, 2, 3},
	}

	actual := packet.NewOpen([]byte{1, 2, 3})

	assert.Equal(t, expected, actual, "invalid open packet")
}

func TestNewClosePacket(t *testing.T) {
	expected := packet.Packet{
		Binary: false,
		Type:   packet.Close,
		Data:   nil,
	}

	actual := packet.NewClose()

	assert.Equal(t, expected, actual, "invalid close packet")
}

func TestNewPongPacket(t *testing.T) {
	expected := packet.Packet{
		Binary: false,
		Type:   packet.Pong,
		Data:   []byte("probe"),
	}

	actual := packet.NewPong([]byte("probe"))

	assert.Equal(t, expected, actual, "invalid pong packet")
}

func TestNewStringPacket(t *testing.T) {
	expected := packet.Packet{
		Binary: false,
		Type:   packet.Message,
		Data:   []byte("hello"),
	}

	actual := packet.NewStringMessage("hello")

	assert.Equal(t, expected, actual, "invalid string message packet")
}

func TestNewBinaryPacket(t *testing.T) {
	expected := packet.Packet{
		Binary: true,
		Type:   packet.Message,
		Data:   []byte{1, 2, 3},
	}

	actual := packet.NewBinaryMessage([]byte{1, 2, 3})

	assert.Equal(t, expected, actual, "invalid binary message packet")
}
