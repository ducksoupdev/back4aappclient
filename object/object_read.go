package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const unableToReadObjectMessage = "unable to read object"

func (c *Object) Read(className string, id string) (map[string]interface{}, *Error) {
	// create the URL
	readUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	readClassUrl := c.baseUrl.ResolveReference(readUrl)

	// create the request
	req, _ := http.NewRequest("GET", readClassUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("Error: ", err)
		return nil, &Error{StatusCode: 500, Err: err}
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
				Err:        errors.New(unableToReadObjectMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToReadObjectMessage)
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
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}

	return result, nil
}
