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
func (codec XHR) Encode(payload packet.Payload) []byte {
	var buffer bytes.Buffer

	for _, packet := range payload {
		buffer.Write(codec.encodePacket(packet))
	}

	return buffer.Bytes()
}

// TODO: Error handling
func (codec XHR) Decode(encoded []byte) packet.Payload {
	var payload packet.Payload

	var buffer bytes.Buffer

	for i := 0; i < len(encoded); i++ {
		ch := rune(encoded[i])

		if ch != ':' {
			buffer.WriteRune(ch)
			continue
		}

		length, _ := strconv.Atoi(buffer.String())
		start := i + 1
		end := start + length

		payload = append(payload, codec.decodePacket(encoded[start:end]))

		buffer.Reset()

		i = end - 1
	}

	return payload
}

func (codec XHR) encodePacket(packet packet.Packet) []byte {
	var data []byte

	if packet.Binary {
		data = codec.encodeBinaryData(packet)
	} else {
		data = codec.encodeStringData(packet)
	}

	var buffer bytes.Buffer

	length := len(data)

	buffer.WriteString(strconv.Itoa(length))
	buffer.WriteRune(':')
	buffer.Write(data)

	return buffer.Bytes()
}

func (codec XHR) encodeStringData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune(rune(packet.Type))
	buffer.Write(packet.Data)

	return buffer.Bytes()
}

func (codec XHR) encodeBinaryData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune('b')
	buffer.WriteRune(rune(packet.Type))
	buffer.WriteString(base64.StdEncoding.EncodeToString(packet.Data))

	return buffer.Bytes()
}

func (codec XHR) decodePacket(data []byte) packet.Packet {
	binary := (data[0] == 'b')

	if binary {
		return codec.decodeBinaryData(data[1:])
	}

	return codec.decodeStringData(data)
}

func (codec XHR) decodeStringData(data []byte) packet.Packet {
	var decoded []byte

	if len(data) > 1 {
		decoded = data[1:]
	}

	return packet.Packet{
		Binary: false,
		Type:   packet.Type(data[0]),
		Data:   decoded,
	}
}

func (codec XHR) decodeBinaryData(data []byte) packet.Packet {
	var decoded []byte

	if len(data) > 1 {
		decoded, _ = base64.StdEncoding.DecodeString(string(data[1:]))
	}

	return packet.Packet{
		Binary: true,
		Type:   packet.Type(data[0]),
		Data:   decoded,
	}
}
