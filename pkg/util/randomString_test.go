package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStr(t *testing.T) {
	length := 32
	result, err := GenerateRandomString(length)
	t.Log("GenerateRandomString: ", result)
	assert.Nil(t, err)
	assert.Equal(t, length, len(result))
}
