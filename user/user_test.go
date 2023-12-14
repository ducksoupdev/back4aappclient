package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewUser(t *testing.T) {
	s := NewUser("applicationId", "restApiKey", nil, nil)
	assert.NotNil(t, s.client)
	assert.NotNil(t, s.baseUrl)
}

func TestLogin(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"sessionToken":"sessionToken"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	sessionToken, _ := s.login("username", "password")
	assert.NotEmptyf(t, sessionToken, "Expected sessionToken to be initialized")
	assert.NotEmptyf(t, s.sessionToken, "Expected sessionToken to be initialized")
}

func TestLoginError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	item, err := s.login("username", "password")
	assert.Equal(t, "", item)
	assert.Error(t, err)
	assert.Equal(t, "unable to login: 400", err.Error())
}

func TestSignUp(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"sessionToken":"sessionToken"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	var data = make(map[string]interface{})
	data["username"] = "username"
	data["password"] = "password"
	sessionToken, _ := s.signUp(data)
	assert.NotEmptyf(t, sessionToken, "Expected sessionToken to be initialized")
	assert.NotEmptyf(t, s.sessionToken, "Expected sessionToken to be initialized")
}

func TestSignUpError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	var data = make(map[string]interface{})
	data["username"] = "username"
	data["password"] = "password"
	sessionToken, err := s.signUp(data)
	assert.Equal(t, "", sessionToken)
	assert.Error(t, err)
	assert.Equal(t, "unable to sign up: 400", err.Error())
}

func TestRequestPasswordReset(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.requestPasswordReset("email")
	assert.Nil(t, err)
}

func TestRequestPasswordResetError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.requestPasswordReset("email")
	assert.Error(t, err)
	assert.Equal(t, "request password reset failed: 400", err.Error())
}
