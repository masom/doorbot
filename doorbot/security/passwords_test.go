// +build tests

package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPassword(t *testing.T) {

	result, err := PasswordCrypt([]byte("test"))

	assert.Nil(t, err)
	assert.NotEmpty(t, result)

	assert.NoError(t, PasswordCompare(result, []byte("test")))
}
