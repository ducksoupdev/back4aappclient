package object

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const unableToCreateObjectMessage = "unable to create object"

func (c *Object) Create(className string, data map[string]interface{}) (map[string]interface{}, *Error) {
	// create the URL
	createUrl, _ := url.Parse(fmt.Sprintf("/classes/%s", className))
	createClassUrl := c.baseUrl.ResolveReference(createUrl)

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("POST", createClassUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, &Error{StatusCode: 500, Err: err}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
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
				Err:        errors.New(unableToCreateObjectMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToCreateObjectMessage)
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
		}
	}

	// parse the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}

	return result, nil
}
