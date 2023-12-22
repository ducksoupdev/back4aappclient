package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const unableToListObjectsMessage = "unable to list objects"

func (c *Object) List(className string, option ...ListOptions) (map[string][]map[string]interface{}, *Error) {
	// create the URL
	listUrl, _ := url.Parse(fmt.Sprintf("/classes/%s", className))

	// create the query string parameters
	for _, opt := range option {
		q := listUrl.Query()
		if opt.Count != 0 {
			q.Set("count", fmt.Sprintf("%d", opt.Count))
		}
		if opt.Limit != 0 {
			q.Set("limit", fmt.Sprintf("%d", opt.Limit))
		}
		if opt.Skip != 0 {
			q.Set("skip", fmt.Sprintf("%d", opt.Skip))
		}
		if opt.Order != "" {
			q.Set("order", opt.Order)
		}
		if opt.Distinct != "" {
			q.Set("distinct", opt.Distinct)
		}
		if opt.Constraints != "" {
			q.Set("where", opt.Constraints)
		}
		listUrl.RawQuery = q.Encode()
	}

	// create the URL
	listClassUrl := c.baseUrl.ResolveReference(listUrl)

	// create the request
	req, _ := http.NewRequest("GET", listClassUrl.String(), nil)
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
	if resp.StatusCode != http.StatusOK {
		// parse the error result
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		if result == nil || (result["error"] == nil && result["code"] == nil) {
			return nil, &Error{
				StatusCode: resp.StatusCode,
				Err:        errors.New(unableToListObjectsMessage),
			}
		}
		message := getErrorMessage(result["error"].(string), unableToListObjectsMessage)
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(float64),
			Err:           errors.New(message),
		}
	}

	// parse the result
	var result map[string][]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}

	return result, nil
}

func WithCount(i int) ListOptions {
	return ListOptions{
		Count: i,
	}
}

func WithSkip(i int) ListOptions {
	return ListOptions{
		Skip: i,
	}
}

func WithLimit(i int) ListOptions {
	return ListOptions{
		Limit: i,
	}
}

func WithOrder(s string) ListOptions {
	return ListOptions{
		Order: s,
	}
}

func WithDistinct(s string) ListOptions {
	return ListOptions{
		Distinct: s,
	}
}

func WithConstraints(s string) ListOptions {
	return ListOptions{
		Constraints: s,
	}
}
