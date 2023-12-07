package session

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	back4appBaseUrl     = "https://parseapi.back4app.com"
	contentTypeHeader   = "Content-type"
	contentTypeValue    = "application/json"
	applicationIdHeader = "X-Parse-Application-Id"
	restApiKeyHeader    = "X-Parse-REST-API-Key"
	revocableHeader     = "X-Parse-Revocable-Session"
)

type Session struct {
	client        *http.Client
	baseUrl       *url.URL
	applicationId string
	restApiKey    string
	sessionToken  string
}

func NewSession(httpClient *http.Client) *Session {
	s := &Session{client: httpClient}
	s.initialize()
	return s
}

func (s *Session) initialize() {
	if s.client == nil {
		s.client = &http.Client{}
	}
	if s.baseUrl == nil {
		s.baseUrl, _ = url.Parse(back4appBaseUrl)
	}
}

func (s *Session) login(username string, password string) {
	r, _ := s.client.Get(fmt.Sprintf("%s", s.baseUrl))
	r.Header.Add(contentTypeHeader, contentTypeValue)
	r.Header.Add(applicationIdHeader, s.applicationId)
	r.Header.Add(restApiKeyHeader, s.restApiKey)
	r.Header.Add(revocableHeader, "1")
	// TODO: add JSON encoded query params for username and password
}

func (s *Session) createUser() {

}

func (s *Session) passwordReset(email string) {

}
