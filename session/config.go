package session

import (
	"time"
)

// Config holds the configuration for a single session
type Config struct {
	// Expected interval for receiving ping packets from client
	PingInterval time.Duration

	// After no ping is received for PingInterval + PingTimeout,
	// the session is considered to be expired
	PingTimeout time.Duration

	// Maximum buffered packets to be flushed
	// on a single read polling request
	PollingBufferFlushLimit int

	// Maximum received packets to be buffered
	// before write polling requests are blocked
	PollingBufferReceiveLimit int
}
