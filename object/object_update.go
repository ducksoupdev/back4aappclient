package object

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const unableToUpdateObjectMessage = "unable to update object"

func (c *Object) Update(className string, id string, data map[string]interface{}) (bool, *Error) {
	// create the URL
	updateUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	updateClassUrl := c.baseUrl.ResolveReference(updateUrl)

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("PUT", updateClassUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("Error: ", err)
		return false, &Error{StatusCode: 500, Err: err}
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		// parse the error result
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return false, &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(unableToUpdateObjectMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToUpdateObjectMessage)
		return false, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
		}
	}

	return true, nil
}
