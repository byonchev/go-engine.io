package packet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/byonchev/go-engine.io/internal/packet"
)

func TestPacketTypes(t *testing.T) {
	tests := []struct {
		expected packet.Packet
		actual   packet.Packet
	}{
		{
			packet.Packet{
				Binary: false,
				Type:   packet.Open,
				Data:   []byte{1, 2, 3},
			},
			packet.NewOpen([]byte{1, 2, 3}),
		},
		{
			packet.Packet{
				Binary: false,
				Type:   packet.Close,
				Data:   []byte{},
			},
			packet.NewClose(),
		},
		{
			packet.Packet{
				Binary: false,
				Type:   packet.Pong,
				Data:   []byte("probe"),
			},
			packet.NewPong([]byte("probe")),
		},
		{
			packet.Packet{
				Binary: false,
				Type:   packet.Message,
				Data:   []byte("hello"),
			},
			packet.NewStringMessage("hello"),
		},
		{
			packet.Packet{
				Binary: true,
				Type:   packet.Message,
				Data:   []byte{1, 2, 3},
			},
			packet.NewBinaryMessage([]byte{1, 2, 3}),
		},
		{
			packet.Packet{
				Binary: false,
				Type:   packet.Message,
				Data:   []byte("hello"),
			},
			packet.NewMessage(false, []byte("hello")),
		},
		{
			packet.Packet{
				Binary: true,
				Type:   packet.Message,
				Data:   []byte{1, 2, 3},
			},
			packet.NewMessage(true, []byte{1, 2, 3}),
		},
		{
			packet.Packet{
				Binary: false,
				Type:   packet.NOOP,
				Data:   []byte{},
			},
			packet.NewNOOP(),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.actual, "packet was not created properly")
	}
}

func TestPacketToString(t *testing.T) {
	packet := packet.NewStringMessage("hello")

	assert.Equal(t, "[4] hello", packet.String(), "packet was not converted to string properly")
}
