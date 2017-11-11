package session

import "github.com/byonchev/go-engine.io/config"

// Config holds the configuration for a single session
type Config struct {
	config.PingSettings

	Listener MessageListener
}
