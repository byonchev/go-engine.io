package engine

import (
	"../config"
	"../session"
)

type Config struct {
	config.PingSettings

	SIDGenerator session.IDGenerator
}
