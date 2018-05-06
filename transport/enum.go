package transport

type Type string

const (
	WebSocketType Type = "websocket"
	PollingType   Type = "polling"
)
