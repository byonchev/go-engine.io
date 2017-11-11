package packet

// Packet is a single encoded message
type Packet struct {
	Binary bool
	Type   Type
	Data   []byte
}

// NewOpen creates new open packet
func NewOpen(data []byte) Packet {
	return Packet{false, Open, data}
}

// NewClose creates new close packet
func NewClose() Packet {
	return Packet{false, Close, nil}
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