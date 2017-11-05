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
