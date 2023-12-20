package user

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
	sessionTokenHeader  = "X-Parse-Session-Token"
)

type Error struct {
	StatusCode int
	Err        error
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

func (s *User) Login(username string, password string) (map[string]interface{}, *Error) {
	// Define the parameters
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

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
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Err:        errors.New("unable to login"),
		}
	}

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}

	// Save the session
	s.Session = result

	return result, nil
}

func (s *User) SignUp(data map[string]interface{}) (map[string]interface{}, *Error) {
	// create the URL
	usersUrl, _ := url.Parse("/users")
	createUserUrl := s.baseUrl.ResolveReference(usersUrl)

	// Set the session token
	var sessionToken = ""
	if data["sessionToken"] != nil {
		sessionToken = data["sessionToken"].(string)
		// remove the session token from the data
		delete(data, "sessionToken")
	}

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("POST", createUserUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, s.applicationId)
	req.Header.Add(restApiKeyHeader, s.restApiKey)
	req.Header.Add(revocableHeader, "1")

	// If we have a session token, add it to the request
	if sessionToken != "" {
		req.Header.Add(sessionTokenHeader, sessionToken)
	}

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusCreated {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Err:        errors.New("unable to sign up user"),
		}
	}

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}

	// Save the session
	s.Session = result
	return result, nil
}

func (s *User) RequestPasswordReset(email string) *Error {
	// create the URL
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
		return &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		return &Error{
			StatusCode: resp.StatusCode,
			Err:        errors.New("request password reset failed"),
		}
	}

	return nil
}

func (s *User) CurrentUser(sessionToken string) (map[string]interface{}, *Error) {
	// Create the URL with the parameters
	userUrl, _ := url.Parse("/users/me")
	joinedUrl := s.baseUrl.ResolveReference(userUrl)

	// create the request
	req, _ := http.NewRequest("GET", joinedUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, s.applicationId)
	req.Header.Add(restApiKeyHeader, s.restApiKey)
	req.Header.Add(sessionTokenHeader, sessionToken)

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Err:        errors.New("unable to get current user"),
		}
	}

	// Parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	return result, nil
}
