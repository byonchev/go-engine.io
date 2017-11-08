package codec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
)

func TestWebSocketEncodeStringPacket(t *testing.T) {
	codec := codec.WebSocket{}

	p1 := packet.NewStringMessage("hi")

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "4hi"

	assert.Equal(t, actual, expected, "string packet was not encoded properly")
}

func TestWebSocketEncodeBinaryPacket(t *testing.T) {
	codec := codec.WebSocket{}

	p1 := packet.NewBinaryMessage([]byte{0, 1})

	payload := packet.Payload{p1}

	actual := codec.Encode(payload)
	expected := []byte{'4', 0, 1}

	assert.Equal(t, actual, expected, "binary packet was not encoded properly")
}

func TestWebSocketEncodeOnlyFirstPacket(t *testing.T) {
	codec := codec.WebSocket{}

	p1 := packet.NewStringMessage("hello")
	p2 := packet.NewStringMessage("world")

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "4hello"

	assert.Equal(t, actual, expected, "payload was not encoded properly")
}

func TestWebSocketEncodeEmpty(t *testing.T) {
	codec := codec.WebSocket{}

	payload := packet.Payload{}

	actual := codec.Encode(payload)
	expected := []byte{}

	assert.Equal(t, actual, expected, "empty payload should not be encoded")
}

func TestWebSocketDecodeString(t *testing.T) {
	codec := codec.WebSocket{}

	data := []byte("4hello")

	actual, err := codec.Decode(data)
	expected := packet.Payload{packet.NewBinaryMessage([]byte("hello"))}

	assert.Nil(t, err, "error while decoding string payload")
	assert.Equal(t, expected, actual, "string payload was not decoded properly")
}

func TestWebSocketDecodeBinary(t *testing.T) {
	codec := codec.WebSocket{}

	data := []byte{'4', 0, 1}

	actual, err := codec.Decode(data)
	expected := packet.Payload{packet.NewBinaryMessage([]byte{0, 1})}

	assert.Nil(t, err, "binary while decoding string payload")
	assert.Equal(t, expected, actual, "binary payload was not decoded properly")
}

func TestWebSocketDecodeError(t *testing.T) {
	codec := codec.WebSocket{}

	data := []byte{}

	_, err := codec.Decode(data)

	assert.Error(t, err)
}
