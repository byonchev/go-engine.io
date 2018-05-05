package session

import (
	"time"
)

// Config holds the configuration for a single session
type Config struct {
	PingInterval time.Duration
	PingTimeout  time.Duration

	Listener Listener
}
