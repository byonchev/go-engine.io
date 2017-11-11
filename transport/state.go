package transport

type state string

const (
	active   state = "ACTIVE"
	paused   state = "PAUSED"
	shutdown state = "SHUTDOWN"
)
