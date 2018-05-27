package logger

// Logger is an interface for logging messages at different verbosity levels
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

var loggerInstance Logger

// Init sets a logger instance
func Init(logger Logger) {
	loggerInstance = logger
}

// Debug logs message with DEBUG verbosity
func Debug(data ...interface{}) {
	if loggerInstance == nil {
		return
	}

	loggerInstance.Debug(data...)
}

// Info logs message with INFO verbosity
func Info(data ...interface{}) {
	if loggerInstance == nil {
		return
	}

	loggerInstance.Info(data...)
}

// Error logs message with ERROR verbosity
func Error(data ...interface{}) {
	if loggerInstance == nil {
		return
	}

	loggerInstance.Error(data...)
}
