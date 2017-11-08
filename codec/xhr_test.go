package codec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"."
	"../packet"
)

func TestEncodeSingleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("utf8 string")

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "12:4utf8 string"

	assert.Equal(t, expected, actual, "Single string payload was not encoded propery")
}

func TestEncodeMultipleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("hello")
	p2 := packet.NewPong()

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "6:4hello1:3"

	assert.Equal(t, expected, actual, "Multiple strings payload was not encoded properly")
}

func TestEncodeSingleBinaryPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewBinaryMessage([]byte{2, 4, 8})

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "6:b4AgQI"

	assert.Equal(t, expected, actual, "Single binary payload was not encoded properly")
}

func TestEncodeMixedPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewPong()
	p2 := packet.NewBinaryMessage([]byte{42})

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "1:36:b4Kg=="

	assert.Equal(t, expected, actual, "Mixed payload was not encoded properly")
}

func TestDecodeSingleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:4hello")

	actual := codec.Decode(data)
	expected := packet.Payload{packet.NewStringMessage("hello")}

	assert.Equal(t, expected, actual, "Single string payload was not decoded properly")
}

func TestDecodeMultipleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:4hello1:3")

	actual := codec.Decode(data)
	expected := packet.Payload{
		packet.NewStringMessage("hello"),
		packet.NewPong(),
	}

	assert.Equal(t, expected, actual, "Multiple strings payload was not decode properly")
}

func TestDecodeSingleBinaryPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("6:b4AgQI")

	actual := codec.Decode(data)
	expected := packet.Payload{packet.NewBinaryMessage([]byte{2, 4, 8})}

	assert.Equal(t, expected, actual, "Single binary payload was not decoded properly")
}

func TestDecodeMixedPayload(t *testing.T) {
	codec := codec.XHR{}

	data := []byte("1:36:b4Kg==")

	actual := codec.Decode(data)
	expected := packet.Payload{
		packet.NewPong(),
		packet.NewBinaryMessage([]byte{42}),
	}

	assert.Equal(t, expected, actual, "Mixed payload was not decoded properly")
}
