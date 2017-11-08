package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"."
)

func TestUUIDGenerator(t *testing.T) {
	generator := session.UUIDGenerator{}

	sid1 := generator.Generate()
	sid2 := generator.Generate()

	assert.Len(t, sid1, 36, "uuid v4 length should be 36 characters")
	assert.Len(t, sid2, 36, "uuid v4 length should be 36 characters")
	assert.NotEqual(t, sid1, sid2, "generated session ids should be unique")
}
