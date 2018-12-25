package config

import (
	"net/http"
	"time"
)

// Config holds the configuration for a single session
type Config struct {
	PingInterval time.Duration
	PingTimeout  time.Duration

	Transports     []string
	AllowUpgrades  bool
	UpgradeTimeout time.Duration

	PollingBufferFlushLimit   int
	PollingBufferReceiveLimit int
	// HTTPCompression bool

	WebsocketReadBufferSize  int
	WebsocketWriteBufferSize int
	PerMessageDeflate        bool

	CheckOrigin func(*http.Request) bool
}
