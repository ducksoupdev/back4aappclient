package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	s := NewUser("applicationId", "restApiKey", nil, nil)
	assert.NotNil(t, s.client)
	assert.NotNil(t, s.baseUrl)
}
