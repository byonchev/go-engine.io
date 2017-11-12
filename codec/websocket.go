package codec

import (
	"errors"

	"github.com/byonchev/go-engine.io/packet"
)

// WebSocket is a codec for encoding packets for websocket transport.
// Since the protocol has its own framing mechanism, packets are
// sent and received one by one, instead of being grouped in payloads.
type WebSocket struct{}

// Encode encodes a single packet in payload
func (WebSocket) Encode(payload packet.Payload) []byte {
	if len(payload) == 0 {
		return []byte{}
	}

	packet := payload[0]

	encoded := make([]byte, len(packet.Data)+1)

	encoded[0] = byte(packet.Type)
	copy(encoded[1:], packet.Data)

	return encoded
}

// Decode decodes single packet from encoded payload
func (WebSocket) Decode(encoded []byte) (packet.Payload, error) {
	if len(encoded) == 0 {
		return nil, errors.New("invalid packet type")
	}

	decoded := packet.Packet{
		Binary: true,
		Type:   packet.Type(encoded[0]),
		Data:   encoded[1:],
	}

	return packet.Payload{decoded}, nil
}
