package codec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"."
	"../packet"
)

func TestXHREncodeSingleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("utf8 string")

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "12:4utf8 string"

	assert.Equal(t, expected, actual, "single string payload was not encoded propery")
}

func TestXHREncodeMultipleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("hello")
	p2 := packet.NewPong()

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "6:4hello1:3"

	assert.Equal(t, expected, actual, "multiple strings payload was not encoded properly")
}

func TestXHREncodeSingleBinaryPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewBinaryMessage([]byte{2, 4, 8})

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "6:b4AgQI"

	assert.Equal(t, expected, actual, "single binary payload was not encoded properly")
}

func TestXHREncodeMixedPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewPong()
	p2 := packet.NewBinaryMessage([]byte{42})

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "1:36:b4Kg=="

	assert.Equal(t, expected, actual, "Mixed payload was not encoded properly")
}

func TestXHRDecodeSingleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:4hello")

	actual, err := codec.Decode(data)
	expected := packet.Payload{packet.NewStringMessage("hello")}

	assert.Nil(t, err, "error while decoding single string payload")
	assert.Equal(t, expected, actual, "single string payload was not decoded properly")
}

func TestXHRDecodeMultipleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:4hello1:3")

	actual, err := codec.Decode(data)
	expected := packet.Payload{
		packet.NewStringMessage("hello"),
		packet.NewPong(),
	}

	assert.Nil(t, err, "error while decoding multiple strings payload")
	assert.Equal(t, expected, actual, "multiple strings payload was not decode properly")
}

func TestXHRDecodeSingleBinaryPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:b4AgQI")

	actual, err := codec.Decode(data)
	expected := packet.Payload{packet.NewBinaryMessage([]byte{2, 4, 8})}

	assert.Nil(t, err, "error while decoding binary payload")
	assert.Equal(t, expected, actual, "single binary payload was not decoded properly")
}

func TestXHRDecodeMixedPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("1:36:b4Kg==")

	actual, err := codec.Decode(data)
	expected := packet.Payload{
		packet.NewPong(),
		packet.NewBinaryMessage([]byte{42}),
	}

	assert.Nil(t, err, "error while decoding mixed payload")
	assert.Equal(t, expected, actual, "Mixed payload was not decoded properly")
}
