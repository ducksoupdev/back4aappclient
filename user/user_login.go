package user

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

const unableToLoginMessage = "unable to login"

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
	if resp.StatusCode != http.StatusOK {
		// Parse the response
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return nil, &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(unableToLoginMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToLoginMessage)
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
