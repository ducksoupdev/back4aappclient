package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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
	u, _ := s.SignUp(data)
	assert.NotEmptyf(t, u["sessionToken"], "Expected sessionToken to be initialized")
	assert.NotEmptyf(t, s.Session["sessionToken"], "Expected sessionToken to be initialized")
}

func TestSignUpWithSessionToken(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"objectId":"objectId","createdAt":"createdAt"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	var data = make(map[string]interface{})
	data["username"] = "username"
	data["password"] = "password"
	data["sessionToken"] = "sessionToken"
	u, _ := s.SignUp(data)
	assert.NotEmptyf(t, u["objectId"], "Expected objectId to be initialized")
	assert.NotEmptyf(t, u["createdAt"], "Expected createdAt to be initialized")
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
	u, err := s.SignUp(data)
	assert.Equal(t, nil, u["sessionToken"])
	assert.Error(t, err)
	assert.Equal(t, "unable to sign up user: 400", err.Error())
}

func TestSignUpHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":205,"error":"invalid login parameters"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	var data = make(map[string]interface{})
	data["username"] = "username"
	data["password"] = "password"
	u, err := s.SignUp(data)
	assert.Equal(t, nil, u["sessionToken"])
	assert.Error(t, err)
	assert.Equal(t, "invalid login parameters: 400", err.Error())
}
