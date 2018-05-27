package config

import (
	"net/http"
	"time"
)

// Config holds the configuration for a single session
type Config struct {
	// Expected interval for receiving ping packets from client
	PingInterval time.Duration

	// After no ping is received for PingInterval + PingTimeout,
	// the session is considered to be expired
	PingTimeout time.Duration

	// List of supported transports
	Transports []string

	// Whether to allow transport upgrades or not
	AllowUpgrades bool

	// Maximum time to wait for uncompleted upgrade
	UpgradeTimeout time.Duration

	// Maximum buffered packets to be flushed
	// on a single read polling request
	PollingBufferFlushLimit int

	// Maximum received packets to be buffered
	// before write polling requests are blocked
	PollingBufferReceiveLimit int

	// Websocket I/O read buffer size
	WebsocketReadBufferSize int

	// Websocket I/O write buffer size
	WebsocketWriteBufferSize int

	// Whether to enable gzip on polling transport or not
	// HTTPCompression bool

	// Whether to enable websocket permessage-deflate extension or not
	PerMessageDeflate bool

	// Function used by transports to validate request
	// and prevent cross-site request forgery
	CheckOrigin func(*http.Request) bool
}
