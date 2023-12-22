package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDelete(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isDeleted, _ := c.Delete("className", "id")
	assert.True(t, isDeleted)
}

func TestDeleteError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isDeleted, err := c.Delete("className", "id")
	assert.False(t, isDeleted)
	assert.Error(t, err)
	assert.Equal(t, "unable to delete object: 400", err.Error())
}

func TestDeleteHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":1, "error":"error"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isDeleted, err := c.Delete("className", "id")
	assert.False(t, isDeleted)
	assert.Error(t, err)
	assert.Equal(t, "error: 400", err.Error())
}
