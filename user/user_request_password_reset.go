package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const requestPasswordResetFailedMessage = "request password reset failed"

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
		log.Println("Error: ", err)
		return &Error{
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

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		// Parse the response
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(requestPasswordResetFailedMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), requestPasswordResetFailedMessage)
		return &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
		}
	}

	return nil
}
