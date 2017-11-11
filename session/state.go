package session

type state string

const (
	new    state = "NEW"
	active state = "ACTIVE"
	closed state = "CLOSED"
)
