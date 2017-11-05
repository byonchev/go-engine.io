package codec

import (
	"bytes"
	"encoding/base64"
	"strconv"

	"../packet"
)

// XHR is a codec for encoding messages for standard long polling
type XHR struct{}

// Encode encodes payload
func (xhr XHR) Encode(payload packet.Payload) []byte {
	var buffer bytes.Buffer

	for _, packet := range payload {
		buffer.Write(xhr.encodePacket(packet))
	}

	return buffer.Bytes()
}

func (xhr XHR) encodePacket(packet packet.Packet) []byte {
	var data []byte

	if packet.Binary {
		data = xhr.encodeBinaryData(packet)
	} else {
		data = xhr.encodeStringData(packet)
	}

	var buffer bytes.Buffer

	length := len(data)

	buffer.WriteString(strconv.Itoa(length))
	buffer.WriteRune(':')
	buffer.Write(data)

	return buffer.Bytes()
}

func (xhr XHR) encodeStringData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune(rune(packet.Type))
	buffer.Write(packet.Data)

	return buffer.Bytes()
}

func (xhr XHR) encodeBinaryData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune('b')
	buffer.WriteRune(rune(packet.Type))
	buffer.WriteString(base64.StdEncoding.EncodeToString(packet.Data))

	return buffer.Bytes()
}
