package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreate(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"username":"username", "objectId":"objectId"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	obj, _ := c.Create("className", make(map[string]interface{}))
	assert.NotNil(t, obj)
	assert.Equal(t, "username", obj["username"])
}

func TestCreateError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	obj, err := c.Create("className", make(map[string]interface{}))
	assert.Nil(t, obj)
	assert.Error(t, err)
	assert.Equal(t, "unable to create object: 400", err.Error())
}

func TestCreateHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":1, "error":"error"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	obj, err := c.Create("className", make(map[string]interface{}))
	assert.Nil(t, obj)
	assert.Error(t, err)
	assert.Equal(t, "error: 400", err.Error())
}
