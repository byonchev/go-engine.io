package codec

import (
	"io"

	"github.com/byonchev/go-engine.io/internal/packet"
)

// Codec encodes and decodes packets for transportation
type Codec interface {
	Encode(packet.Payload, io.Writer) error
	Decode(io.Reader) (packet.Payload, error)
}
