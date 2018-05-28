package utils_test

import (
	"testing"

	"github.com/byonchev/go-engine.io/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringSliceContains(t *testing.T) {
	haystack := []string{"Hello"}

	assert.True(t, utils.StringSliceContains(haystack, "Hello"))
	assert.False(t, utils.StringSliceContains(haystack, "World"))
}
