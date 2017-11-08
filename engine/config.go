package engine

import (
	"github.com/byonchev/go-engine.io/config"
	"github.com/byonchev/go-engine.io/session"
)

type Config struct {
	config.PingSettings

	SIDGenerator session.IDGenerator
}
