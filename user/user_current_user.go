package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const unableToGetCurrentUserMessage = "unable to get current user"

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
		// Parse the response
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return nil, &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(unableToGetCurrentUserMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToGetCurrentUserMessage)
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

	return result, nil
}
