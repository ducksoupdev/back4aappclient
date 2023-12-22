package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRead(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"item":"item"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	item, _ := c.Read("className", "id")
	assert.NotNil(t, item)
	assert.Equal(t, "item", item["item"])
}

func TestReadError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	item, err := c.Read("className", "id")
	assert.Nil(t, item)
	assert.Error(t, err)
	assert.Equal(t, "unable to read object: 400", err.Error())
}

func TestReadHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":1, "error":"error"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	item, err := c.Read("className", "id")
	assert.Nil(t, item)
	assert.Error(t, err)
	assert.Equal(t, "error: 400", err.Error())
}
