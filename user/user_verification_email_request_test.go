package user

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestVerificationEmailRequest(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.VerificationEmailRequest("email")
	assert.Nil(t, err)
}

func TestVerificationEmailRequestError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.VerificationEmailRequest("email")
	assert.Error(t, err)
	assert.Equal(t, "verify email request failed: 400", err.Error())
}

func TestVerificationEmailRequestHostError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":205,"error":"invalid login parameters"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	s := NewUser("applicationId", "restApiKey", nil, b)
	err := s.VerificationEmailRequest("email")
	assert.Error(t, err)
	assert.Equal(t, "invalid login parameters: 400", err.Error())
}
