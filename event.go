package eio

// ConnectEvent is emitted on new client connection
type ConnectEvent struct {
	SessionID string
}

// DisconnectEvent is emitted on client connection timeout
type DisconnectEvent struct {
	SessionID string
}

// MessageEvent is emitted on received client message
type MessageEvent struct {
	SessionID string
	Binary    bool
	Data      []byte
}
