package codec

import (
	"errors"

	"../packet"
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

	packetType := byte(packet.Type)

	return append([]byte{packetType}, packet.Data...)
}

// Decode decodes single packet from encoded payload
func (WebSocket) Decode(encoded []byte) (packet.Payload, error) {
	size := len(encoded)

	if size == 0 {
		return nil, errors.New("invalid packet type")
	}

	packetType := packet.Type(encoded[0])

	var data []byte

	if size > 1 {
		data = encoded[1:]
	}

	return packet.Payload{packet.Packet{Binary: true, Type: packetType, Data: data}}, nil
}
