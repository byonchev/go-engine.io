package codec

import (
	"bytes"
	"encoding/base64"
	"strconv"

	"../packet"
)

// XHR is a codec for encoding messages for standard long polling
type XHR struct{}

// Encode encodes payload
func (xhr XHR) Encode(payload packet.Payload) []byte {
	var buffer bytes.Buffer

	for _, packet := range payload {
		buffer.Write(xhr.encodePacket(packet))
	}

	return buffer.Bytes()
}

func (xhr XHR) encodePacket(packet packet.Packet) []byte {
	var data string

	packetType := strconv.Itoa(int(packet.Type))

	if packet.Binary {
		data = "b" + packetType + base64.StdEncoding.EncodeToString(packet.Data)
	} else {
		data = packetType + string(packet.Data)
	}

	length := len(data)

	result := strconv.Itoa(length) + ":" + data

	return []byte(result)
}
