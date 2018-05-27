package codec_test

import (
	"bytes"
	"testing"

	"github.com/byonchev/go-engine.io/internal/codec"
	"github.com/byonchev/go-engine.io/internal/packet"
	"github.com/stretchr/testify/assert"
)

func TestXHREncode(t *testing.T) {
	codec := codec.XHR{}

	tests := []struct {
		payload packet.Payload
		encoded string
	}{
		{
			packet.Payload{},
			"",
		},
		{
			packet.Payload{
				packet.NewStringMessage("utf八 string"),
			},
			"12:4utf八 string",
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
				packet.NewStringMessage("Hello👋"),
			},
			string([]byte{1, 4, 255, 4, 2, 4, 8, 0, 1, 0, 255, '4', 'H', 'e', 'l', 'l', 'o', 0xf0, 0x9f, 0x91, 0x8b}),
		},
	}

	for _, test := range tests {
		var buffer bytes.Buffer

		err := codec.Encode(test.payload, &buffer)

		actual := buffer.String()
		expected := test.encoded

		assert.Nil(t, err, "error while encoding payload")
		assert.Equal(t, expected, actual, "payload was not encoded propery")
	}
}

func TestXHREncodeForceBase64(t *testing.T) {
	codec := codec.XHR{ForceBase64: true}

	payload := packet.Payload{
		packet.NewBinaryMessage([]byte{2, 4, 8}),
	}

	var buffer bytes.Buffer

	err := codec.Encode(payload, &buffer)

	expected := []byte("6:b4AgQI")
	actual := buffer.Bytes()

	assert.Nil(t, err, "error while encoding payload")
	assert.Equal(t, expected, actual, "payload was not encoded propery")

}
func TestXHREncodeWriterError(t *testing.T) {
	codec := codec.XHR{}

	err := codec.Encode(packet.Payload{packet.NewNOOP()}, &errorWriter{})

	assert.Error(t, err, "writer error was expected")
}

func TestXHRDecode(t *testing.T) {
	codec := codec.XHR{}

	tests := []struct {
		data    []byte
		decoded packet.Payload
	}{
		{
			[]byte("8:4hello \u2764"),
			packet.Payload{
				packet.NewStringMessage("hello ❤"),
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
		{
			[]byte{1, 4, 255, 4, 2, 4, 8, 0, 1, 0, 255, '4', 'H', 'e', 'l', 'l', 'o', 0xf0, 0x9f, 0x91, 0x8b},
			packet.Payload{
				packet.NewBinaryMessage([]byte{2, 4, 8}),
				packet.NewStringMessage("Hello👋"),
			},
		},
	}

	for _, test := range tests {
		actual, err := codec.Decode(bytes.NewBuffer(test.data))
		expected := test.decoded

		assert.Nil(t, err, "error while decoding valid payload")
		assert.Equal(t, expected, actual, "payload was not decoded properly")
	}
}

func TestXHRDecodeErrors(t *testing.T) {
	codec := codec.XHR{}

	tests := [][]byte{
		[]byte("INVALID_LENGTH:3"),
		[]byte("1:30:"),
		[]byte("6:b4AGQI0:"),
		[]byte("8:bINVALID_BASE64"),
		[]byte("1:b"),
		[]byte{},
		[]byte{1, 5, 255, 4},
		[]byte{1, 0, 255},
	}

	for _, test := range tests {
		payload, err := codec.Decode(bytes.NewBuffer(test))

		assert.Empty(t, payload, "decoded invalid payload was not empty")
		assert.Error(t, err, "error was expected for decoding "+string(test))
	}
}

func TestXHRDecodeReaderError(t *testing.T) {
	codec := codec.XHR{}

	_, err := codec.Decode(errorReader{})

	assert.Error(t, err, "reader error was expected")
}

func BenchmarkXHREncodeString(b *testing.B) {
	codec := codec.XHR{ForceBase64: true}

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

func BenchmarkXHREncodeBinary(b *testing.B) {
	codec := codec.XHR{}

	payload := packet.Payload{
		packet.NewBinaryMessage([]byte{'!'}),
		packet.NewOpen([]byte("hello")),
		packet.NewStringMessage("world"),
	}

	var buffer bytes.Buffer

	for n := 0; n < b.N; n++ {
		codec.Encode(payload, &buffer)
	}
}

func BenchmarkXHRDecodeString(b *testing.B) {
	codec := codec.XHR{}

	data := []byte("6:0hello6:4world6:b4IQ==")

	buffer := bytes.NewBuffer(data)

	for n := 0; n < b.N; n++ {
		codec.Decode(buffer)
	}
}

func BenchmarkXHRDecodeBinary(b *testing.B) {
	codec := codec.XHR{}

	data := []byte{1, 2, 255, 4, 33, 0, 6, 255, 48, 104, 101, 108, 108, 111, 0, 6, 255, 52, 119, 111, 114, 108, 100}

	buffer := bytes.NewBuffer(data)

	for n := 0; n < b.N; n++ {
		codec.Decode(buffer)
	}
}
