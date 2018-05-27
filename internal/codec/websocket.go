package codec

import (
	"errors"
	"io"
	"io/ioutil"
	"unicode"

	"github.com/byonchev/go-engine.io/internal/packet"
)

// Websocket is a codec for encoding packets for websocket transport
type Websocket struct{}

// Encode encodes a single packet in payload
func (codec Websocket) Encode(payload packet.Payload, writer io.Writer) error {
	if len(payload) == 0 {
		return nil
	}

	for _, packet := range payload {
		err := codec.encodePacket(packet, writer)

		if err != nil {
			return err
		}
	}

	return nil
}

// Decode decodes single packet from encoded payload
func (Websocket) Decode(reader io.Reader) (packet.Payload, error) {
	encoded, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	if len(encoded) == 0 {
		return nil, errors.New("invalid packet type")
	}

	var binary bool
	var packetType packet.Type

	typeByte := encoded[0]

	if unicode.IsNumber(rune(typeByte)) {
		binary = false
		packetType = packet.TypeFromChar(typeByte)
	} else {
		binary = true
		packetType = packet.TypeFromByte(typeByte)
	}

	data := encoded[1:]

	decoded := packet.Packet{
		Binary: binary,
		Type:   packetType,
		Data:   data,
	}

	return packet.Payload{decoded}, nil
}

func (codec Websocket) encodePacket(packet packet.Packet, writer io.Writer) error {
	encoded := make([]byte, len(packet.Data)+1)

	var packetType byte

	if packet.Binary {
		packetType = packet.Type.Byte()
	} else {
		packetType = packet.Type.Char()
	}

	encoded[0] = packetType
	copy(encoded[1:], packet.Data)

	_, err := writer.Write(encoded)

	return err
}
