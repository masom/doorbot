// +build tests

package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAPIToken(t *testing.T) {
	assert.Len(t, GenerateAPIToken(), 36)
	assert.NotEqual(t, GenerateAPIToken(), GenerateAPIToken())
}
