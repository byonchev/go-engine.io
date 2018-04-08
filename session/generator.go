package session

import "github.com/satori/go.uuid"

// IDGenerator is an interface for generating session IDs
type IDGenerator interface {
	Generate() string
}

// UUIDGenerator generates UUIDv4 session Ids
type UUIDGenerator struct{}

// Generate returns UUIDv4 string
func (generator UUIDGenerator) Generate() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}
