package codec

import "../packet"

type Codec interface {
	Encode(packet.Payload) []byte
	Decode([]byte) packet.Payload
}
