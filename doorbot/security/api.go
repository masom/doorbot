package security

import (
	"github.com/satori/go.uuid"
)

// GenerateAPIToken generates a unique token
func GenerateAPIToken() string {
	return uuid.NewV4().String()
}
