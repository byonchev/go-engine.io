package codec

import "github.com/byonchev/go-engine.io/packet"

// Codec encodes and decodes packets for transportation
type Codec interface {
	Encode(packet.Payload) []byte
	Decode([]byte) (packet.Payload, error)
}
