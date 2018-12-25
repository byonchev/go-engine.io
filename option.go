package eio

import (
	"net/http"
	"time"

	"github.com/byonchev/go-engine.io/internal/config"
	"github.com/byonchev/go-engine.io/internal/logger"
)

// Option defines a single configuration change
type Option func(*config.Config)

// PingInterval sets the expected interval for receiving ping packets
func PingInterval(interval time.Duration) Option {
	return func(config *config.Config) { config.PingInterval = interval }
}

// PingTimeout sets the maximum time to wait after the expected ping interval
func PingTimeout(timeout time.Duration) Option {
	return func(config *config.Config) { config.PingTimeout = timeout }
}

// Transports explicitly sets which transports are enabled ("websocket", "polling")
func Transports(transports ...string) Option {
	return func(config *config.Config) { config.Transports = transports }
}

// EnableUpgrades enables transport upgrades
func EnableUpgrades() Option {
	return func(config *config.Config) { config.AllowUpgrades = true }
}

// DisableUpgrades disables transport upgrades
func DisableUpgrades() Option {
	return func(config *config.Config) { config.AllowUpgrades = false }
}

// UpgradeTimeout sets the timeout for the transport upgrade process
func UpgradeTimeout(timeout time.Duration) Option {
	return func(config *config.Config) { config.UpgradeTimeout = timeout }
}

// PollingBufferFlushLimit sets the maximum packets to be flushed on a single poll request
func PollingBufferFlushLimit(limit int) Option {
	return func(config *config.Config) { config.PollingBufferFlushLimit = limit }
}

// PollingBufferReceiveLimit sets the maximum packets to be buffered before write requests are blocked
func PollingBufferReceiveLimit(limit int) Option {
	return func(config *config.Config) { config.PollingBufferReceiveLimit = limit }
}

// WebsocketReadBufferSize sets the websocket I/O read buffer size
func WebsocketReadBufferSize(size int) Option {
	return func(config *config.Config) { config.WebsocketReadBufferSize = size }
}

// WebsocketWriteBufferSize sets the websocket I/O write buffer size
func WebsocketWriteBufferSize(size int) Option {
	return func(config *config.Config) { config.WebsocketWriteBufferSize = size }
}

// EnablePerMessageDeflate enables the websocket per-message-deflate extension
func EnablePerMessageDeflate() Option {
	return func(config *config.Config) { config.PerMessageDeflate = true }
}

// DisablePerMessageDeflate disables the websocket per-message-deflate extension
func DisablePerMessageDeflate() Option {
	return func(config *config.Config) { config.PerMessageDeflate = false }
}

// OriginCheckFunction sets the function to validate incoming requests
func OriginCheckFunction(function func(*http.Request) bool) Option {
	return func(config *config.Config) { config.CheckOrigin = function }
}

// Logger sets the default logger for the library
func Logger(loggerInstance logger.Logger) Option {
	return func(*config.Config) { logger.Init(loggerInstance) }
}
