package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

const unableToSignUpMessage = "unable to sign up user"

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
		log.Println("Error: ", err)
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error: ", err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusCreated {
		// Parse the response
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return nil, &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(unableToSignUpMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToSignUpMessage)
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
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
