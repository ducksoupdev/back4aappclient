package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLogin(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"sessionToken":"sessionToken"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	u, _ := s.Login("username", "password")
	assert.NotEmptyf(t, u["sessionToken"], "Expected sessionToken to be initialized")
	assert.NotEmptyf(t, s.Session["sessionToken"], "Expected sessionToken to be initialized")
}

func TestLoginError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	u, err := s.Login("username", "password")
	assert.Equal(t, map[string]interface{}(nil), u)
	assert.Error(t, err)
	assert.Equal(t, "unable to login: 400", err.Error())
}

func TestLoginHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":205,"error":"invalid login parameters"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	u, err := s.Login("username", "password")
	assert.Equal(t, map[string]interface{}(nil), u)
	assert.Error(t, err)
	assert.Equal(t, "invalid login parameters: 400", err.Error())
}
