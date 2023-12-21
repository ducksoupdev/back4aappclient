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

const (
	back4appBaseUrl     = "https://parseapi.back4app.com"
	contentTypeHeader   = "Content-type"
	contentTypeValue    = "application/json"
	applicationIdHeader = "X-Parse-Application-Id"
	restApiKeyHeader    = "X-Parse-REST-API-Key"
	sessionTokenHeader  = "X-Parse-Session-Token"
)

type Error struct {
	StatusCode    int
	HostErrorCode int
	Err           error
}

func (r *Error) Error() string {
	return fmt.Sprintf("%v: %d", r.Err, r.StatusCode)
}

type Object struct {
	httpClient    *http.Client
	baseUrl       *url.URL
	applicationId string
	restApiKey    string
	sessionToken  string
}

type ListOptions struct {
	Count       int
	Limit       int
	Skip        int
	Order       string
	Distinct    string
	Constraints string
}

func getErrorMessage(error string, defaultError string) string {
	if error == "" {
		return defaultError
	}
	return error
}

func NewObject(applicationId string, restApiKey string, sessionToken string, httpClient *http.Client, baseUrl *url.URL) *Object {
	c := &Object{
		httpClient:    httpClient,
		baseUrl:       baseUrl,
		applicationId: applicationId,
		restApiKey:    restApiKey,
		sessionToken:  sessionToken,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	if c.baseUrl == nil {
		c.baseUrl, _ = url.Parse(back4appBaseUrl)
	}
	return c
}

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

	// parse the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{StatusCode: 500, Err: err}
	}

	// check the status code
	if resp.StatusCode != http.StatusCreated {
		message := getErrorMessage(result["error"].(string), "unable to create object")
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(int),
			Err:           errors.New(message),
		}
	}

	return result, nil
}

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
		fmt.Println(err)
		return false, &Error{StatusCode: 500, Err: err}
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		// parse the error result
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return false, &Error{
				StatusCode: 500,
				Err:        err,
			}
		}
		message := getErrorMessage(result["error"].(string), "unable to delete object")
		return false, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(int),
			Err:           errors.New(message),
		}
	}

	return true, nil
}

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
		fmt.Println(err)
		return nil, &Error{StatusCode: 500, Err: err}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// parse the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &Error{
			StatusCode: 500,
			Err:        err,
		}
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		message := getErrorMessage(result["error"].(string), "unable to read object")
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(int),
			Err:           errors.New(message),
		}
	}

	return result, nil
}

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
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, &Error{
				StatusCode: 500,
				Err:        err,
			}
		}
		message := getErrorMessage(result["error"].(string), "unable to list objects")
		return nil, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(int),
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
		fmt.Println(err)
		return false, &Error{StatusCode: 500, Err: err}
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		// parse the error result
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return false, &Error{
				StatusCode: 500,
				Err:        err,
			}
		}
		message := getErrorMessage(result["error"].(string), "unable to update object")
		return false, &Error{
			StatusCode:    resp.StatusCode,
			HostErrorCode: result["code"].(int),
			Err:           errors.New(message),
		}
	}

	return true, nil
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
