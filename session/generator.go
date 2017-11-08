package session

import "github.com/satori/go.uuid"

type IDGenerator interface {
	Generate() string
}

type UUIDGenerator struct{}

func (generator UUIDGenerator) Generate() string {
	return uuid.NewV4().String()
}
