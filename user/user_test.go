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

func TestRequestPasswordReset(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.RequestPasswordReset("email")
	assert.Nil(t, err)
}

func TestRequestPasswordResetError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.RequestPasswordReset("email")
	assert.Error(t, err)
	assert.Equal(t, "request password reset failed: 400", err.Error())
}

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
