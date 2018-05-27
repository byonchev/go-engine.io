package utils_test

import (
	"net/url"
	"testing"

	"github.com/byonchev/go-engine.io/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateBase64ID(t *testing.T) {
	id := utils.GenerateBase64ID()

	escaped := url.QueryEscape(id)

	assert.Equal(t, id, escaped, "id is not query safe")
	assert.Len(t, id, 22, "id length is not 22 characters")
}
