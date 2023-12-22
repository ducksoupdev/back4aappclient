package object

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

type Object struct {
	httpClient    *http.Client
	baseUrl       *url.URL
	applicationId string
	restApiKey    string
	sessionToken  string
}

type ListOptions struct {
	Count       int
	Limit       int
	Skip        int
	Order       string
	Distinct    string
	Constraints string
}

func getErrorMessage(error string, defaultError string) string {
	if error == "" {
		return defaultError
	}
	return error
}

func NewObject(applicationId string, restApiKey string, sessionToken string, httpClient *http.Client, baseUrl *url.URL) *Object {
	c := &Object{
		httpClient:    httpClient,
		baseUrl:       baseUrl,
		applicationId: applicationId,
		restApiKey:    restApiKey,
		sessionToken:  sessionToken,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	if c.baseUrl == nil {
		c.baseUrl, _ = url.Parse(back4appBaseUrl)
	}
	return c
}
