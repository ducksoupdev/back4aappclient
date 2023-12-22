package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const unableToDeleteObjectMessage = "unable to delete object"

func (c *Object) Delete(className string, id string) (bool, *Error) {
	// create the URL
	deleteUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	deleteClassUrl := c.baseUrl.ResolveReference(deleteUrl)

	// create the request
	req, _ := http.NewRequest("DELETE", deleteClassUrl.String(), nil)
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
				Err:        errors.New(unableToDeleteObjectMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToDeleteObjectMessage)
		return false, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
		}
	}

	return true, nil
}
