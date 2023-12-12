package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (s *Session) login(username string, password string) (string, error) {
	// Define the parameters
	params := url.Values{}
	usernameB, _ := json.Marshal(username)
	passwordB, _ := json.Marshal(password)
	params.Add("username", string(usernameB))
	params.Add("password", string(passwordB))

	// Create the URL with the parameters
	loginUrl, _ := url.Parse("/login")
	joinedUrl := s.baseUrl.ResolveReference(loginUrl)
	joinedUrl.RawQuery = params.Encode()

	// create the request
	req, _ := http.NewRequest("GET", joinedUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, s.applicationId)
	req.Header.Add(restApiKeyHeader, s.restApiKey)
	req.Header.Add(revocableHeader, "1")

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	s.sessionToken = result["sessionToken"].(string)
	return s.sessionToken, nil
}

func (s *Session) createUser(data map[string]interface{}) (string, error) {
	usersUrl, _ := url.Parse("/users")
	createUserUrl := s.baseUrl.ResolveReference(usersUrl)

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("POST", createUserUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, s.applicationId)
	req.Header.Add(restApiKeyHeader, s.restApiKey)
	req.Header.Add(revocableHeader, "1")

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	s.sessionToken = result["sessionToken"].(string)
	return s.sessionToken, nil
}

func (s *Session) passwordReset(email string) error {
	resetUrl, _ := url.Parse("/requestPasswordReset")
	requestPasswordResetUrl := s.baseUrl.ResolveReference(resetUrl)

	// create the body
	var jsonBody = []byte(fmt.Sprintf(`{"email":"%s"}`, email))

	// create the request
	req, _ := http.NewRequest("POST", requestPasswordResetUrl.String(), bytes.NewBuffer(jsonBody))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, s.applicationId)
	req.Header.Add(restApiKeyHeader, s.restApiKey)

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		return errors.New("request password reset failed")
	}

	return nil
}
