package config

import "time"

// PingSettings hold configuration for session pings
type PingSettings struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
}
