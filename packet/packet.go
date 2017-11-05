package packet

// Packet is a single encoded message
type Packet struct {
	Binary bool
	Type   Type
	Data   []byte
}

// NewClose creates new close packet
func NewClose() Packet {
	return Packet{false, Close, nil}
}

// NewPong creates new pong packet
func NewPong() Packet {
	return Packet{false, Pong, nil}
}

// NewStringMessage creates new string message packet
func NewStringMessage(data string) Packet {
	return Packet{false, Message, []byte(data)}
}

// NewBinaryMessage creates new binary message packet
func NewBinaryMessage(data []byte) Packet {
	return Packet{true, Message, data}
}
