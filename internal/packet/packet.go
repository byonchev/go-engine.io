package packet

import (
	"fmt"
)

// Packet is a single encoded message
type Packet struct {
	Binary bool
	Type   Type
	Data   []byte
}

func (packet Packet) String() string {
	return fmt.Sprintf("[%d] %s", packet.Type, string(packet.Data))
}

// NewOpen creates new open packet
func NewOpen(data []byte) Packet {
	return Packet{false, Open, data}
}

// NewClose creates new close packet
func NewClose() Packet {
	return Packet{false, Close, []byte{}}
}

// NewPong creates new pong packet
func NewPong(data []byte) Packet {
	return Packet{false, Pong, data}
}

// NewStringMessage creates new string message packet
func NewStringMessage(data string) Packet {
	return Packet{false, Message, []byte(data)}
}

// NewBinaryMessage creates new binary message packet
func NewBinaryMessage(data []byte) Packet {
	return Packet{true, Message, data}
}

// NewMessage creates new message depending on its type
func NewMessage(binary bool, data []byte) Packet {
	if binary {
		return NewBinaryMessage(data)
	}

	return NewStringMessage(string(data))
}

// NewNOOP creates new NOOP packet
func NewNOOP() Packet {
	return Packet{false, NOOP, []byte{}}
}
