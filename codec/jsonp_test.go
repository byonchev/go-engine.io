package codec_test

import (
	"bytes"
	"testing"

	"github.com/byonchev/go-engine.io/codec"
	"github.com/byonchev/go-engine.io/packet"
	"github.com/stretchr/testify/assert"
)

func TestJSONPEncode(t *testing.T) {
	codec := codec.JSONP{Index: "0"}

	tests := []struct {
		payload packet.Payload
		encoded string
	}{
		{
			packet.Payload{
				packet.NewStringMessage("javascript-safe \n\b\f\r\t \" / \u0015 \u516B\u2028\u2029"),
			},
			`___eio[0]("32:4javascript-safe \n\b\f\r\t \" \/ \u0015 八\u2028\u2029");`,
		},
		{
			packet.Payload{
				packet.NewOpen([]byte("hello")),
				packet.NewStringMessage("world"),
			},
			"___eio[0](\"6:0hello6:4world\");",
		},
		{
			packet.Payload{
				packet.NewBinaryMessage([]byte{2, 4, 8}),
			},
			"___eio[0](\"6:b4AgQI\");",
		},
		{
			packet.Payload{
				packet.NewClose(),
				packet.NewBinaryMessage([]byte{42}),
			},
			"___eio[0](\"1:16:b4Kg==\");",
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

func TestJSONPDecode(t *testing.T) {
	codec := codec.JSONP{}

	tests := []struct {
		data    []byte
		decoded packet.Payload
	}{
		{
			[]byte(`d=4:4\\n\\\\n`),
			packet.Payload{
				packet.NewStringMessage("\n\\n"),
			},
		},
		{
			[]byte("d=6:4hello"),
			packet.Payload{
				packet.NewStringMessage("hello"),
			},
		},
		{
			[]byte("d=6:4hello6:4world6:3probe"),
			packet.Payload{
				packet.NewStringMessage("hello"),
				packet.NewStringMessage("world"),
				packet.NewPong([]byte("probe")),
			},
		},
		{
			[]byte("d=6:b4AgQI"),
			packet.Payload{
				packet.NewBinaryMessage([]byte{2, 4, 8}),
			},
		},
		{
			[]byte("d=1:16:b4Kg=="),
			packet.Payload{
				packet.NewClose(),
				packet.NewBinaryMessage([]byte{42}),
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

func TestJSONPDecodeErrors(t *testing.T) {
	codec := codec.JSONP{}

	tests := [][]byte{
		[]byte("1"),
		[]byte("123%"),
		[]byte("INVALID_LENGTH:3"),
		[]byte("1:30:"),
		[]byte("6:b4AGQI0:"),
		[]byte("d=INVALID_LENGTH:3"),
		[]byte("d=1:30:"),
		[]byte("d=6:b4AGQI0:"),
		[]byte("d=8:bINVALID_BASE64"),
	}

	for _, test := range tests {
		payload, err := codec.Decode(bytes.NewBuffer(test))

		assert.Empty(t, payload, "decoded invalid payload was not empty")
		assert.Error(t, err, "error was expected for decoding "+string(test))
	}
}

func TestJSONPDecodeReaderError(t *testing.T) {
	codec := codec.JSONP{}

	_, err := codec.Decode(errorReader{})

	assert.Error(t, err, "reader error was expected")
}

func BenchmarkJSONPEncode(b *testing.B) {
	codec := codec.JSONP{Index: "0"}

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

func BenchmarkJSONPDecode(b *testing.B) {
	codec := codec.JSONP{}

	data := []byte("d=6:0hello6:4world6:b4IQ==")

	buffer := bytes.NewBuffer(data)

	for n := 0; n < b.N; n++ {
		codec.Decode(buffer)
	}
}
