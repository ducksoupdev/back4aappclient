package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"item":"item"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	item, _ := s.CurrentUser("sessionToken")
	assert.NotNil(t, item)
	assert.Equal(t, "item", item["item"])
}

func TestCurrentUserError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	u, err := s.CurrentUser("sessionToken")
	assert.Equal(t, map[string]interface{}(nil), u)
	assert.Error(t, err)
	assert.Equal(t, "unable to get current user: 400", err.Error())
}

func TestCurrentUserHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":205,"error":"invalid login parameters"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	u, err := s.CurrentUser("sessionToken")
	assert.Equal(t, map[string]interface{}(nil), u)
	assert.Error(t, err)
	assert.Equal(t, "invalid login parameters: 400", err.Error())
}
