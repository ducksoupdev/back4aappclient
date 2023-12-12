package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSession(t *testing.T) {
	s := NewSession(nil)
	if s.client == nil {
		t.Error("Expected client to be initialized")
	}
	if s.baseUrl == nil {
		t.Error("Expected baseUrl to be initialized")
	}
}

func TestInitialize(t *testing.T) {
	s := NewSession(nil)
	s.initialize()
	if s.client == nil {
		t.Error("Expected client to be initialized")
	}
	if s.baseUrl == nil {
		t.Error("Expected baseUrl to be initialized")
	}
	if s.baseUrl.String() != "https://parseapi.back4app.com" {
		t.Error("Expected baseUrl to be initialized with back4appBaseUrl")
	}
}

func TestLogin(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sessionToken":"sessionToken"}`))
	}))
	defer svr.Close()
	s := NewSession(nil)
	s.baseUrl, _ = s.baseUrl.Parse(svr.URL)
	s.applicationId = "applicationId"
	s.restApiKey = "restApiKey"
	s.initialize()
	sessionToken, _ := s.login("username", "password")
	if sessionToken == "" {
		t.Error("Expected sessionToken to be initialized")
	}
	if s.sessionToken == "" {
		t.Error("Expected sessionToken to be initialized")
	}
}

func TestCreateUser(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"sessionToken":"sessionToken"}`))
	}))
	defer svr.Close()
	s := NewSession(nil)
	s.baseUrl, _ = s.baseUrl.Parse(svr.URL)
	s.applicationId = "applicationId"
	s.restApiKey = "restApiKey"
	s.initialize()
	var data = make(map[string]interface{})
	data["username"] = "username"
	data["password"] = "password"
	sessionToken, _ := s.createUser(data)
	if sessionToken == "" {
		t.Error("Expected sessionToken to be initialized")
	}
	if s.sessionToken == "" {
		t.Error("Expected sessionToken to be initialized")
	}
}

func TestRequestPasswordReset(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	s := NewSession(nil)
	s.baseUrl, _ = s.baseUrl.Parse(svr.URL)
	s.applicationId = "applicationId"
	s.restApiKey = "restApiKey"
	s.initialize()
	err := s.passwordReset("email")
	if err != nil {
		t.Error("Expected password reset to be successful")
	}
}
