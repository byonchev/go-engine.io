package codec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"."
	"../packet"
)

func TestSingleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("utf8 string")

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "12:4utf8 string"

	assert.Equal(t, expected, actual, "Single string payload was not encoded propery")
}

func TestMultipleStringPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewStringMessage("hello")
	p2 := packet.NewPong()

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "6:4hello1:3"

	assert.Equal(t, expected, actual, "Multiple strings payload was not encoded propery")
}

func TestSingleBinaryPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewBinaryMessage([]byte{2, 4, 8})

	payload := packet.Payload{p1}

	actual := string(codec.Encode(payload))
	expected := "6:b4AgQI"

	assert.Equal(t, expected, actual, "Single binary payload was not encoded propery")
}

func TestMixedPayload(t *testing.T) {
	codec := codec.XHR{}

	p1 := packet.NewPong()
	p2 := packet.NewBinaryMessage([]byte{42})

	payload := packet.Payload{p1, p2}

	actual := string(codec.Encode(payload))
	expected := "1:36:b4Kg=="

	assert.Equal(t, expected, actual, "Mixed payload was not encoded propery")
}
