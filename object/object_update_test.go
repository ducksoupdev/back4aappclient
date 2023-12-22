package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUpdate(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"updatedAt":"updatedAt"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isUpdated, _ := c.Update("className", "id", make(map[string]interface{}))
	assert.True(t, isUpdated)
}

func TestUpdateError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isUpdated, err := c.Update("className", "id", make(map[string]interface{}))
	assert.False(t, isUpdated)
	assert.Error(t, err)
	assert.Equal(t, "unable to update object: 400", err.Error())
}

func TestUpdateHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":1, "error":"error"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isUpdated, err := c.Update("className", "id", make(map[string]interface{}))
	assert.False(t, isUpdated)
	assert.Error(t, err)
	assert.Equal(t, "error: 400", err.Error())
}
