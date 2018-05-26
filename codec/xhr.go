package codec

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"unicode/utf8"

	"github.com/byonchev/go-engine.io/packet"
)

var base64Encoding = base64.StdEncoding

// XHR is a codec for encoding messages for standard long polling
type XHR struct {
	ForceBase64 bool
}

// Encode encodes payload of packets for single poll
func (codec XHR) Encode(payload packet.Payload, writer io.Writer) error {
	if len(payload) == 0 {
		return nil
	}

	binary := !codec.ForceBase64 && payload.ContainsBinary()

	for _, packet := range payload {
		var err error

		if binary {
			err = codec.encodeBinaryPacket(packet, writer)
		} else {
			err = codec.encodeStringPacket(packet, writer)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Decode decodes payload of packets
func (codec XHR) Decode(reader io.Reader) (packet.Payload, error) {
	encoded, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	if len(encoded) == 0 {
		return nil, errors.New("payload is empty")
	}

	if encoded[0] <= 1 {
		return codec.decodeBinaryPayload(encoded)
	}

	return codec.decodeStringPayload(encoded)
}

func (codec XHR) encodeStringPacket(packet packet.Packet, writer io.Writer) error {
	var data []byte
	var length int

	if packet.Binary {
		data = codec.encodeBase64Data(packet)
		length = len(data)
	} else {
		data = codec.encodeStringData(packet)
		length = utf8.RuneCount(data)
	}

	var encoded []byte

	encoded = append(encoded, []byte(strconv.Itoa(length))...)
	encoded = append(encoded, ':')
	encoded = append(encoded, data...)

	_, err := writer.Write(encoded)

	return err
}

func (codec XHR) encodeStringData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteByte(packet.Type.Char())
	buffer.Write(packet.Data)

	return buffer.Bytes()
}

func (codec XHR) encodeBase64Data(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune('b')
	buffer.WriteByte(packet.Type.Char())
	buffer.WriteString(base64Encoding.EncodeToString(packet.Data))

	return buffer.Bytes()
}

func (codec XHR) encodeBinaryPacket(packet packet.Packet, writer io.Writer) error {
	var messageType byte
	var packetType byte

	if packet.Binary {
		messageType = 1
		packetType = packet.Type.Byte()
	} else {
		messageType = 0
		packetType = packet.Type.Char()
	}

	var lengthBytes []byte

	for length := len(packet.Data) + 1; length > 0; length /= 10 {
		digit := byte(length % 10)

		lengthBytes = append([]byte{digit}, lengthBytes...)
	}

	var encoded []byte

	encoded = append(encoded, messageType)
	encoded = append(encoded, lengthBytes...)
	encoded = append(encoded, 255)
	encoded = append(encoded, packetType)
	encoded = append(encoded, packet.Data...)

	_, err := writer.Write(encoded)

	return err
}

func (codec XHR) decodeStringPayload(data []byte) (packet.Payload, error) {
	var payload packet.Payload
	var lengthRunes []rune

	runes := []rune(string(data))
	total := len(runes)

	for i := 0; i < total; i++ {
		r := runes[i]

		if r != ':' {
			lengthRunes = append(lengthRunes, r)
			continue
		}

		length, err := strconv.Atoi(string(lengthRunes))

		if err != nil {
			return nil, errors.New("invalid packet length")
		}

		start := i + 1
		end := start + length

		if end > total {
			return nil, errors.New("packet length overflow")
		}

		packet, err := codec.decodeStringPacket([]byte(string(runes[start:end])))

		if err != nil {
			return nil, err
		}

		payload = append(payload, packet)

		lengthRunes = nil
		i = end - 1
	}

	return payload, nil
}

func (codec XHR) decodeBinaryPayload(data []byte) (packet.Payload, error) {
	var payload packet.Payload

	total := len(data)

	for offset := 0; offset < total-1; {
		messageType := data[offset]

		offset++

		length := 0

		for offset < total && data[offset] != 255 {
			length = length*10 + int(data[offset])
			offset++
		}

		offset++

		start := offset
		end := offset + length

		if end > total {
			return nil, errors.New("packet length overflow")
		}

		packet, err := codec.decodeBinaryPacket(messageType, data[start:end])

		if err != nil {
			return nil, err
		}

		payload = append(payload, packet)
		offset += length
	}

	return payload, nil
}

func (codec XHR) decodeBinaryPacket(messageType byte, data []byte) (packet.Packet, error) {
	if len(data) < 1 {
		return packet.Packet{}, errors.New("invalid packet")
	}

	var binary bool
	var packetType packet.Type

	if messageType == 1 {
		binary = true
	}

	if binary {
		packetType = packet.TypeFromByte(data[0])
	} else {
		packetType = packet.TypeFromChar(data[0])
	}

	return packet.Packet{
		Binary: binary,
		Data:   data[1:],
		Type:   packetType,
	}, nil
}

func (codec XHR) decodeStringPacket(data []byte) (packet.Packet, error) {
	if len(data) < 1 {
		return packet.Packet{}, errors.New("invalid packet")
	}

	if data[0] == 'b' {
		return codec.decodeBase64Packet(data[1:])
	}

	return packet.Packet{
		Binary: false,
		Type:   packet.TypeFromChar(data[0]),
		Data:   data[1:],
	}, nil
}

func (codec XHR) decodeBase64Packet(data []byte) (packet.Packet, error) {
	var decoded []byte
	var err error

	if len(data) < 1 {
		return packet.Packet{}, errors.New("invalid packet")
	}

	decoded, err = base64Encoding.DecodeString(string(data[1:]))

	if err != nil {
		return packet.Packet{}, errors.New("base64 decoding error: " + err.Error())
	}

	return packet.Packet{
		Binary: true,
		Type:   packet.TypeFromChar(data[0]),
		Data:   decoded,
	}, nil
}
