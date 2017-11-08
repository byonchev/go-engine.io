package config

import "time"

type PingSettings struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
}
