package user

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
	sessionTokenHeader  = "X-Parse-Session-Token"
)

type Error struct {
	StatusCode    int
	HostErrorCode float64
	Err           error
}

func (r *Error) Error() string {
	return fmt.Sprintf("%v: %d", r.Err, r.StatusCode)
}

type User struct {
	client        *http.Client
	baseUrl       *url.URL
	applicationId string
	restApiKey    string
	Session       map[string]interface{}
}

func getErrorMessage(error string, defaultError string) string {
	if error == "" {
		return defaultError
	}
	return error
}

func NewUser(applicationId string, restApiKey string, httpClient *http.Client, baseUrl *url.URL) *User {
	s := &User{
		client:        httpClient,
		baseUrl:       baseUrl,
		applicationId: applicationId,
		restApiKey:    restApiKey,
	}
	if s.client == nil {
		s.client = &http.Client{}
	}
	if s.baseUrl == nil {
		s.baseUrl, _ = url.Parse(back4appBaseUrl)
	}
	return s
}
