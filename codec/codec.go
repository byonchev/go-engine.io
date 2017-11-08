package codec

import "../packet"

// Codec encodes and decodes packets for transportation
type Codec interface {
	Encode(packet.Payload) []byte
	Decode([]byte) (packet.Payload, error)
}
