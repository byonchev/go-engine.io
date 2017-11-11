package codec

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/byonchev/go-engine.io/packet"
)

// XHR is a codec for encoding messages for standard long polling
type XHR struct{}

// Encode encodes payload of packets for single sending.
func (codec XHR) Encode(payload packet.Payload) []byte {
	var buffer bytes.Buffer

	for _, packet := range payload {
		buffer.Write(codec.encodePacket(packet))
	}

	return buffer.Bytes()
}

// Decode decodes payload of packets
func (codec XHR) Decode(encoded []byte) (packet.Payload, error) {
	var buffer bytes.Buffer
	var payload packet.Payload

	for i := 0; i < len(encoded); i++ {
		ch := rune(encoded[i])

		if ch != ':' {
			buffer.WriteRune(ch)
			continue
		}

		length, err := strconv.Atoi(buffer.String())

		if err != nil {
			return nil, errors.New("invalid packet length")
		}

		start := i + 1
		end := start + length

		packet, err := codec.decodePacket(encoded[start:end])

		if err != nil {
			return nil, err
		}

		payload = append(payload, packet)

		buffer.Reset()

		i = end - 1
	}

	return payload, nil
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

func (codec XHR) decodePacket(data []byte) (packet.Packet, error) {
	if len(data) == 0 {
		return packet.Packet{}, errors.New("packet type missing")
	}

	binary := (data[0] == 'b')

	if binary {
		return codec.decodeBinaryData(data[1:])
	}

	return codec.decodeStringData(data)
}

func (codec XHR) decodeStringData(data []byte) (packet.Packet, error) {
	var decoded []byte

	if len(data) > 1 {
		decoded = data[1:]
	}

	return packet.Packet{
		Binary: false,
		Type:   packet.Type(data[0]),
		Data:   decoded,
	}, nil
}

func (codec XHR) decodeBinaryData(data []byte) (packet.Packet, error) {
	var decoded []byte
	var err error

	if len(data) > 1 {
		decoded, err = base64.StdEncoding.DecodeString(string(data[1:]))

		if err != nil {
			return packet.Packet{}, errors.New("base64 decoding error: " + err.Error())
		}
	}

	return packet.Packet{
		Binary: true,
		Type:   packet.Type(data[0]),
		Data:   decoded,
	}, nil
}