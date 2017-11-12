package codec_test

import (
	"testing"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketEncode(t *testing.T) {
	codec := codec.WebSocket{}

	tests := []struct {
		payload packet.Payload
		encoded []byte
	}{
		{
			packet.Payload{
				packet.NewStringMessage("hello"),
			},
			[]byte("4hello"),
		},
		{
			packet.Payload{
				packet.NewBinaryMessage([]byte{0, 1}),
			},
			[]byte{'4', 0, 1},
		},
		{
			packet.Payload{
				packet.NewStringMessage("hello"),
				packet.NewStringMessage("world"),
			},
			[]byte("4hello"),
		},
		{
			packet.Payload{
				packet.NewClose(),
			},
			[]byte{'1'},
		},
		{
			packet.Payload{},
			[]byte{},
		},
	}

	for _, test := range tests {
		actual := codec.Encode(test.payload)
		expected := test.encoded

		assert.Equal(t, actual, expected, "payload was not encoded properly")
	}
}

func TestWebSocketDecode(t *testing.T) {
	codec := codec.WebSocket{}

	tests := []struct {
		data    []byte
		decoded packet.Payload
	}{
		{
			[]byte("4hello"),
			packet.Payload{
				packet.NewBinaryMessage([]byte("hello")),
			},
		},
		{
			[]byte{'4', 0, 1},
			packet.Payload{
				packet.NewBinaryMessage([]byte{0, 1}),
			},
		},
	}

	for _, test := range tests {
		actual, err := codec.Decode(test.data)
		expected := test.decoded

		assert.Nil(t, err, "error while decoding valid payload")
		assert.Equal(t, expected, actual, "payload was not decoded properly")
	}
}

func TestWebSocketDecodeErrors(t *testing.T) {
	codec := codec.WebSocket{}

	data := []byte{}

	payload, err := codec.Decode(data)

	assert.Empty(t, payload, "decoded invalid payload was not empty")
	assert.Error(t, err)
}

func BenchmarkWebSocketEncode(b *testing.B) {
	codec := codec.WebSocket{}

	payloads := []packet.Payload{
		packet.Payload{packet.NewOpen([]byte("hello"))},
		packet.Payload{packet.NewStringMessage("world")},
		packet.Payload{packet.NewBinaryMessage([]byte{'!'})},
	}

	for n := 0; n < b.N; n++ {
		for _, payload := range payloads {
			codec.Encode(payload)
		}
	}
}

func BenchmarkWebSocketDecode(b *testing.B) {
	codec := codec.WebSocket{}

	packets := [][]byte{
		[]byte("0hello"),
		[]byte("4world"),
		[]byte{'4', '!'},
	}

	for n := 0; n < b.N; n++ {
		for _, data := range packets {
			codec.Decode(data)
		}
	}
}
