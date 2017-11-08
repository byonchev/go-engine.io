package engine

import (
	"github.com/byonchev/go-engine.io/config"
	"github.com/byonchev/go-engine.io/session"
)

// Config holds global engine.io server configuration
type Config struct {
	config.PingSettings

	SIDGenerator session.IDGenerator
}
