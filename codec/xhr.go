package codec

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/byonchev/go-engine.io/packet"
)

var base64Encoding = base64.StdEncoding

// XHR is a codec for encoding messages for standard long polling
type XHR struct{}

// Encode encodes payload of packets for single poll
func (codec XHR) Encode(payload packet.Payload, writer io.Writer) error {
	for _, packet := range payload {
		err := codec.encodePacket(packet, writer)

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

	packets, err := codec.splitPayload(encoded)

	if err != nil {
		return nil, err
	}

	payload := make(packet.Payload, len(packets))

	for i, packet := range packets {
		decoded, err := codec.decodePacket(packet)

		if err != nil {
			return nil, err
		}

		payload[i] = decoded
	}

	return payload, nil
}

func (codec XHR) encodePacket(packet packet.Packet, writer io.Writer) error {
	var data []byte

	if packet.Binary {
		data = codec.encodeBinaryData(packet)
	} else {
		data = codec.encodeStringData(packet)
	}

	length := len(data)

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

func (codec XHR) encodeBinaryData(packet packet.Packet) []byte {
	var buffer bytes.Buffer

	buffer.WriteRune('b')
	buffer.WriteByte(packet.Type.Char())
	buffer.WriteString(base64Encoding.EncodeToString(packet.Data))

	return buffer.Bytes()
}

func (codec XHR) splitPayload(data []byte) ([][]byte, error) {
	var packets [][]byte
	var lengthBytes []byte

	for i := 0; i < len(data); i++ {
		b := data[i]

		if b != ':' {
			lengthBytes = append(lengthBytes, b)
			continue
		}

		length, err := strconv.Atoi(string(lengthBytes))

		if err != nil {
			return nil, errors.New("invalid packet length")
		}

		start := i + 1
		end := start + length

		packets = append(packets, data[start:end])

		lengthBytes = nil
		i = end - 1
	}

	return packets, nil
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
		Type:   packet.TypeFromChar(data[0]),
		Data:   decoded,
	}, nil
}

func (codec XHR) decodeBinaryData(data []byte) (packet.Packet, error) {
	var decoded []byte
	var err error

	if len(data) > 1 {
		decoded, err = base64Encoding.DecodeString(string(data[1:]))

		if err != nil {
			return packet.Packet{}, errors.New("base64 decoding error: " + err.Error())
		}
	}

	return packet.Packet{
		Binary: true,
		Type:   packet.TypeFromChar(data[0]),
		Data:   decoded,
	}, nil
}
