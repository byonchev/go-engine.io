package packet_test

import (
	"testing"

	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/stretchr/testify/assert"
)

func TestPayloadContainsBinary(t *testing.T) {
	tests := []struct {
		Payload  packet.Payload
		Expected bool
	}{
		{
			Payload: packet.Payload{
				packet.NewBinaryMessage(nil),
				packet.NewStringMessage(""),
			},
			Expected: true,
		},
		{
			Payload: packet.Payload{
				packet.NewOpen(nil),
				packet.NewClose(),
			},
			Expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, test.Payload.ContainsBinary())
	}
}
