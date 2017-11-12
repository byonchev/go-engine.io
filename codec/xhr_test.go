package codec_test

import (
	"testing"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestXHREncode(t *testing.T) {
	codec := codec.XHR{}

	tests := []struct {
		payload packet.Payload
		encoded string
	}{
		{
			packet.Payload{
				packet.NewStringMessage("utf八 string"),
			},
			"14:4utf八 string",
		},
		{
			packet.Payload{
				packet.NewOpen([]byte("hello")),
				packet.NewStringMessage("world"),
			},
			"6:0hello6:4world",
		},
		{
			packet.Payload{
				packet.NewBinaryMessage([]byte{2, 4, 8}),
			},
			"6:b4AgQI",
		},
		{
			packet.Payload{
				packet.NewClose(),
				packet.NewBinaryMessage([]byte{42}),
			},
			"1:16:b4Kg==",
		},
	}

	for _, test := range tests {
		actual := string(codec.Encode(test.payload))
		expected := test.encoded

		assert.Equal(t, expected, actual, "payload was not encoded propery")
	}
}

func TestXHRDecode(t *testing.T) {
	codec := codec.XHR{}

	tests := []struct {
		data    []byte
		decoded packet.Payload
	}{
		{
			[]byte("6:4hello"),
			packet.Payload{
				packet.NewStringMessage("hello"),
			},
		},
		{
			[]byte("6:4hello6:4world6:3probe"),
			packet.Payload{
				packet.NewStringMessage("hello"),
				packet.NewStringMessage("world"),
				packet.NewPong([]byte("probe")),
			},
		},
		{
			[]byte("6:b4AgQI"),
			packet.Payload{
				packet.NewBinaryMessage([]byte{2, 4, 8}),
			},
		},
		{
			[]byte("1:16:b4Kg=="),
			packet.Payload{
				packet.NewClose(),
				packet.NewBinaryMessage([]byte{42}),
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

func TestXHRDecodeErrors(t *testing.T) {
	codec := codec.XHR{}

	tests := [][]byte{
		[]byte("INVALID_FORMAT"),
		[]byte("INVALID_LENGTH:3"),
		[]byte("1:30:"),
		[]byte("6:b4AGQI0:"),
		[]byte("8:bINVALID_BASE64"),
	}

	for _, test := range tests {
		payload, err := codec.Decode(test)

		assert.Empty(t, payload, "decoded invalid payload was not empty")
		assert.Error(t, err)
	}
}

func BenchmarkXHREncode(b *testing.B) {
	codec := codec.XHR{}

	payload := packet.Payload{
		packet.NewOpen([]byte("hello")),
		packet.NewStringMessage("world"),
		packet.NewBinaryMessage([]byte{'!'}),
	}

	for n := 0; n < b.N; n++ {
		codec.Encode(payload)
	}
}

func BenchmarkXHRDecode(b *testing.B) {
	codec := codec.XHR{}

	data := []byte("6:0hello6:4world6:b4IQ==")

	for n := 0; n < b.N; n++ {
		codec.Decode(data)
	}
}
