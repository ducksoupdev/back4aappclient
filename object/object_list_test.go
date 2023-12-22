package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestList(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results":[{"item":"item"},{"item2":"item2"}]}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, _ := c.List("className")
	assert.NotNil(t, list)
	assert.Len(t, list["results"], 2)
}

func TestListWithOptions(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "5", r.URL.Query().Get("count"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "10", r.URL.Query().Get("skip"))
		assert.Equal(t, "order", r.URL.Query().Get("order"))
		assert.Equal(t, "distinct", r.URL.Query().Get("distinct"))
		assert.Equal(t, "where", r.URL.Query().Get("where"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results":[{"item":"item"},{"item2":"item2"}]}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, _ := c.List("className", WithCount(5), WithLimit(10), WithSkip(10), WithOrder("order"), WithDistinct("distinct"), WithConstraints("where"))
	assert.NotNil(t, list)
	assert.Len(t, list["results"], 2)
}

func TestListError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, err := c.List("className")
	assert.Nil(t, list)
	assert.Error(t, err)
	assert.Equal(t, "unable to list objects: 400", err.Error())
}

func TestListHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":1, "error":"error"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, err := c.List("className")
	assert.Nil(t, list)
	assert.Error(t, err)
	assert.Equal(t, "error: 400", err.Error())
}
