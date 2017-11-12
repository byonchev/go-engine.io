package codec

import (
	"io"

	"github.com/byonchev/go-engine.io/packet"
)

// Codec encodes and decodes packets for transportation
type Codec interface {
	Encode(packet.Payload) []byte
	Decode(io.Reader) (packet.Payload, error)
}
