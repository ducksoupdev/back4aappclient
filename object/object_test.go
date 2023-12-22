package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewObject(t *testing.T) {
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, nil)
	assert.NotNil(t, c.httpClient)
	assert.NotNil(t, c.baseUrl)
}

func TestInitialize(t *testing.T) {
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, nil)
	assert.NotNil(t, c.httpClient)
	assert.NotNil(t, c.baseUrl)
	assert.Equal(t, c.applicationId, "applicationId")
	assert.Equal(t, c.restApiKey, "restApiKey")
	assert.Equal(t, c.sessionToken, "sessionToken")
	assert.Equal(t, c.baseUrl.String(), "https://parseapi.back4app.com")
}
