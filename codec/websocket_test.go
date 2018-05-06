package codec_test

import (
	"bytes"
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
			[]byte{4, 0, 1},
		},
		{
			packet.Payload{
				packet.NewStringMessage("hello"),
				packet.NewStringMessage("world"),
			},
			[]byte("4hello4world"),
		},
		{
			packet.Payload{
				packet.NewClose(),
			},
			[]byte{'1'},
		},
		{
			packet.Payload{},
			nil,
		},
	}

	for _, test := range tests {
		var buffer bytes.Buffer

		err := codec.Encode(test.payload, &buffer)

		actual := buffer.Bytes()
		expected := test.encoded

		assert.Nil(t, err, "error while encoding valid payload")
		assert.Equal(t, actual, expected, "payload was not encoded properly")
	}
}

func TestWebSocketEncodeWriterError(t *testing.T) {
	codec := codec.WebSocket{}

	err := codec.Encode(packet.Payload{packet.NewNOOP()}, errorWriter{})

	assert.Error(t, err, "reader error was expected")
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
				packet.NewStringMessage("hello"),
			},
		},
		{
			[]byte{4, 0, 1},
			packet.Payload{
				packet.NewBinaryMessage([]byte{0, 1}),
			},
		},
	}

	for _, test := range tests {
		buffer := bytes.NewBuffer(test.data)

		actual, err := codec.Decode(buffer)
		expected := test.decoded

		assert.Nil(t, err, "error while decoding valid payload")
		assert.Equal(t, expected, actual, "payload was not decoded properly")
	}
}

func TestWebSocketDecodeErrors(t *testing.T) {
	codec := codec.WebSocket{}

	data := []byte{}

	buffer := bytes.NewBuffer(data)

	payload, err := codec.Decode(buffer)

	assert.Empty(t, payload, "decoded invalid payload was not empty")
	assert.Error(t, err)
}

func TestWebSocketDecodeReaderError(t *testing.T) {
	codec := codec.WebSocket{}

	_, err := codec.Decode(errorReader{})

	assert.Error(t, err, "reader error was expected")
}

func BenchmarkWebSocketEncode(b *testing.B) {
	codec := codec.WebSocket{}

	payload := packet.Payload{
		packet.NewOpen([]byte("hello")),
		packet.NewStringMessage("world"),
		packet.NewBinaryMessage([]byte{'!'}),
	}

	var buffer bytes.Buffer

	for n := 0; n < b.N; n++ {
		codec.Encode(payload, &buffer)
	}
}

func BenchmarkWebSocketDecode(b *testing.B) {
	codec := codec.WebSocket{}

	buffer := bytes.NewBuffer([]byte("0hello"))

	for n := 0; n < b.N; n++ {
		codec.Decode(buffer)

	}
}
