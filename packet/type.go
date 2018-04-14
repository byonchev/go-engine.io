package packet

// Type defines the packet type
type Type byte

// Supported packet types
const (
	Open    Type = 0
	Close   Type = 1
	Ping    Type = 2
	Pong    Type = 3
	Message Type = 4
	Upgrade Type = 5
	NOOP    Type = 6
)

// Char returns character representation of the type
func (t Type) Char() byte {
	return t.Byte() + '0'
}

// Byte returns numeric representation of the type
func (t Type) Byte() byte {
	return byte(t)
}

// TypeFromChar return packet type from character representation
func TypeFromChar(data byte) Type {
	return Type(data - '0')
}

// TypeFromByte returns packet type from numeric representation
func TypeFromByte(data byte) Type {
	return Type(data)
}
