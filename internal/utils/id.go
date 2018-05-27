package utils

import (
	"encoding/base64"

	"github.com/satori/go.uuid"
)

// GenerateBase64ID generates UUID v4 and encodes it into base64 without padding
func GenerateBase64ID() string {
	uuid, _ := uuid.NewV4()

	return base64.RawURLEncoding.EncodeToString(uuid.Bytes())
}